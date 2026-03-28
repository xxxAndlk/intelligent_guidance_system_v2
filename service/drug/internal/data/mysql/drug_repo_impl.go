package mysql

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"intelligent-guidance-system/service/drug/internal/domain/aggregate"
	"intelligent-guidance-system/service/drug/internal/domain/entity"
	"intelligent-guidance-system/service/drug/internal/domain/repository"
)

type DrugRepoImpl struct {
	db *gorm.DB
}

func NewDrugRepoImpl(db *gorm.DB) repository.DrugRepository {
	return &DrugRepoImpl{db: db}
}

func (r *DrugRepoImpl) Save(ctx context.Context, drug *aggregate.Drug) error {
	po := r.toPO(drug)

	if po.ID == 0 {
		if err := r.db.WithContext(ctx).Create(&po).Error; err != nil {
			return err
		}
		drug.SetID(po.ID)
	} else {
		if err := r.db.WithContext(ctx).Save(&po).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *DrugRepoImpl) FindByID(ctx context.Context, id int64) (*aggregate.Drug, error) {
	var po DrugPO
	if err := r.db.WithContext(ctx).First(&po, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.toAggregate(&po), nil
}

func (r *DrugRepoImpl) FindByName(ctx context.Context, name string) (*aggregate.Drug, error) {
	var po DrugPO
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.toAggregate(&po), nil
}

func (r *DrugRepoImpl) FindByCategory(ctx context.Context, category int) ([]*aggregate.Drug, error) {
	var pos []DrugPO
	categoryCode := entity.DrugCategory(category).Code()
	if err := r.db.WithContext(ctx).Where("category = ?", categoryCode).Find(&pos).Error; err != nil {
		return nil, err
	}
	return r.toAggregates(pos), nil
}

func (r *DrugRepoImpl) FindByStatus(ctx context.Context, status int) ([]*aggregate.Drug, error) {
	var pos []DrugPO
	statusCode := entity.DrugStatus(status).Code()
	if err := r.db.WithContext(ctx).Where("status = ?", statusCode).Find(&pos).Error; err != nil {
		return nil, err
	}
	return r.toAggregates(pos), nil
}

func (r *DrugRepoImpl) FindAll(ctx context.Context, page, pageSize int) ([]*aggregate.Drug, int64, error) {
	var pos []DrugPO
	var total int64

	if err := r.db.WithContext(ctx).Model(&DrugPO{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	return r.toAggregates(pos), total, nil
}

func (r *DrugRepoImpl) SearchDrugs(ctx context.Context, keyword string) ([]*aggregate.Drug, error) {
	var pos []DrugPO
	if err := r.db.WithContext(ctx).Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%").Find(&pos).Error; err != nil {
		return nil, err
	}
	return r.toAggregates(pos), nil
}

func (r *DrugRepoImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&DrugPO{}, id).Error
}

func (r *DrugRepoImpl) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&DrugPO{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *DrugRepoImpl) toPO(drug *aggregate.Drug) *DrugPO {
	return &DrugPO{
		ID:              drug.ID(),
		Name:            drug.Name(),
		Description:     drug.Description(),
		Price:           drug.Price(),
		QuantityInStock: drug.QuantityInStock(),
		Category:        drug.Category().Code(),
		UseMethod:       drug.UseMethod(),
		ExpirationDate:  drug.ExpirationDate(),
		Status:          drug.Status().Code(),
		CreatedAt:       drug.CreatedAt(),
		UpdatedAt:       drug.UpdatedAt(),
	}
}

func (r *DrugRepoImpl) toAggregate(po *DrugPO) *aggregate.Drug {
	return aggregate.ReconstructDrug(
		po.ID,
		po.Name,
		po.Description,
		po.Price,
		po.QuantityInStock,
		entity.DrugCategoryFromCode(po.Category),
		po.UseMethod,
		po.ExpirationDate,
		entity.DrugStatusFromCode(po.Status),
		nil,
		po.CreatedAt,
		po.UpdatedAt,
	)
}

func (r *DrugRepoImpl) toAggregates(pos []DrugPO) []*aggregate.Drug {
	drugs := make([]*aggregate.Drug, 0, len(pos))
	for _, po := range pos {
		drugs = append(drugs, r.toAggregate(&po))
	}
	return drugs
}

type DrugStockRepoImpl struct {
	db *gorm.DB
}

func NewDrugStockRepoImpl(db *gorm.DB) repository.DrugStockRepository {
	return &DrugStockRepoImpl{db: db}
}

func (r *DrugStockRepoImpl) Save(ctx context.Context, stock *entity.DrugStock) error {
	po := DrugStockPO{
		ID:             stock.ID(),
		DrugID:         stock.DrugID(),
		BatchNumber:    stock.BatchNumber(),
		Quantity:       stock.Quantity(),
		ExpirationDate: stock.ExpirationDate(),
		PurchasePrice:  stock.PurchasePrice(),
		SellingPrice:   stock.SellingPrice(),
		Supplier:       stock.Supplier(),
		CreatedAt:      stock.CreatedAt(),
		UpdatedAt:      stock.UpdatedAt(),
	}

	if po.ID == 0 {
		if err := r.db.WithContext(ctx).Create(&po).Error; err != nil {
			return err
		}
		stock.SetID(po.ID)
	} else {
		if err := r.db.WithContext(ctx).Save(&po).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *DrugStockRepoImpl) FindByID(ctx context.Context, id int64) (*entity.DrugStock, error) {
	var po DrugStockPO
	if err := r.db.WithContext(ctx).First(&po, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return entity.ReconstructDrugStock(
		po.ID, po.DrugID, po.BatchNumber, po.Quantity,
		po.ExpirationDate, po.PurchasePrice, po.SellingPrice, po.Supplier,
		po.CreatedAt, po.UpdatedAt,
	), nil
}

func (r *DrugStockRepoImpl) FindByDrugID(ctx context.Context, drugID int64) ([]*entity.DrugStock, error) {
	var pos []DrugStockPO
	if err := r.db.WithContext(ctx).Where("drug_id = ?", drugID).Find(&pos).Error; err != nil {
		return nil, err
	}

	stocks := make([]*entity.DrugStock, 0, len(pos))
	for _, po := range pos {
		stocks = append(stocks, entity.ReconstructDrugStock(
			po.ID, po.DrugID, po.BatchNumber, po.Quantity,
			po.ExpirationDate, po.PurchasePrice, po.SellingPrice, po.Supplier,
			po.CreatedAt, po.UpdatedAt,
		))
	}

	return stocks, nil
}

func (r *DrugStockRepoImpl) FindAvailableStocks(ctx context.Context, drugID int64) ([]*entity.DrugStock, error) {
	var pos []DrugStockPO
	if err := r.db.WithContext(ctx).
		Where("drug_id = ? AND quantity > 0 AND expiration_date > NOW()", drugID).
		Order("expiration_date ASC").
		Find(&pos).Error; err != nil {
		return nil, err
	}

	stocks := make([]*entity.DrugStock, 0, len(pos))
	for _, po := range pos {
		stocks = append(stocks, entity.ReconstructDrugStock(
			po.ID, po.DrugID, po.BatchNumber, po.Quantity,
			po.ExpirationDate, po.PurchasePrice, po.SellingPrice, po.Supplier,
			po.CreatedAt, po.UpdatedAt,
		))
	}

	return stocks, nil
}

func (r *DrugStockRepoImpl) FindExpiringSoon(ctx context.Context, days int) ([]*entity.DrugStock, error) {
	var pos []DrugStockPO
	if err := r.db.WithContext(ctx).
		Where("expiration_date BETWEEN NOW() AND DATE_ADD(NOW(), INTERVAL ? DAY)", days).
		Find(&pos).Error; err != nil {
		return nil, err
	}

	stocks := make([]*entity.DrugStock, 0, len(pos))
	for _, po := range pos {
		stocks = append(stocks, entity.ReconstructDrugStock(
			po.ID, po.DrugID, po.BatchNumber, po.Quantity,
			po.ExpirationDate, po.PurchasePrice, po.SellingPrice, po.Supplier,
			po.CreatedAt, po.UpdatedAt,
		))
	}

	return stocks, nil
}

func (r *DrugStockRepoImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&DrugStockPO{}, id).Error
}