package usecase

import (
	"context"
	"errors"

	"intelligent-guidance-system/service/drug/internal/domain/aggregate"
	"intelligent-guidance-system/service/drug/internal/domain/entity"
	"intelligent-guidance-system/service/drug/internal/domain/repository"
)

var (
	ErrDrugNotFound = errors.New("drug not found")
)

type DrugUsecase struct {
	drugRepo  repository.DrugRepository
	stockRepo repository.DrugStockRepository
}

func NewDrugUsecase(drugRepo repository.DrugRepository, stockRepo repository.DrugStockRepository) *DrugUsecase {
	return &DrugUsecase{drugRepo: drugRepo, stockRepo: stockRepo}
}

func (u *DrugUsecase) CreateDrug(ctx context.Context,
	name string,
	description string,
	price int64,
	category entity.DrugCategory,
	useMethod string,
) (*aggregate.Drug, error) {
	exists, err := u.drugRepo.ExistsByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("drug name already exists")
	}

	drug, err := aggregate.NewDrug(name, description, price, category, useMethod)
	if err != nil {
		return nil, err
	}

	if err := u.drugRepo.Save(ctx, drug); err != nil {
		return nil, err
	}

	return drug, nil
}

func (u *DrugUsecase) GetDrug(ctx context.Context, id int64) (*aggregate.Drug, error) {
	drug, err := u.drugRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if drug == nil {
		return nil, ErrDrugNotFound
	}
	return drug, nil
}

func (u *DrugUsecase) UpdateDrug(ctx context.Context, id int64, description string, price int64, useMethod string) error {
	drug, err := u.GetDrug(ctx, id)
	if err != nil {
		return err
	}

	drug.UpdateDescription(description)
	drug.UpdatePrice(price)
	drug.useMethod = useMethod
	drug.updatedAt = drug.UpdatedAt()

	return u.drugRepo.Save(ctx, drug)
}

func (u *DrugUsecase) UpdateStock(ctx context.Context, id int64, quantity int, operatorID int64) error {
	drug, err := u.GetDrug(ctx, id)
	if err != nil {
		return err
	}

	return drug.UpdateStock(quantity, operatorID)
}

func (u *DrugUsecase) AddStock(ctx context.Context, id int64, quantity int, operatorID int64) error {
	drug, err := u.GetDrug(ctx, id)
	if err != nil {
		return err
	}

	return drug.AddStock(quantity, operatorID)
}

func (u *DrugUsecase) ReduceStock(ctx context.Context, id int64, quantity int, operatorID int64) error {
	drug, err := u.GetDrug(ctx, id)
	if err != nil {
		return err
	}

	return drug.ReduceStock(quantity, operatorID)
}

func (u *DrugUsecase) CheckAvailability(ctx context.Context, id int64, quantity int) (bool, error) {
	drug, err := u.GetDrug(ctx, id)
	if err != nil {
		return false, err
	}

	return drug.CheckAvailability(quantity), nil
}

func (u *DrugUsecase) ListDrugs(ctx context.Context, page, pageSize int) ([]*aggregate.Drug, int64, error) {
	return u.drugRepo.FindAll(ctx, page, pageSize)
}

func (u *DrugUsecase) SearchDrugs(ctx context.Context, keyword string) ([]*aggregate.Drug, error) {
	return u.drugRepo.SearchDrugs(ctx, keyword)
}

func (u *DrugUsecase) DiscontinueDrug(ctx context.Context, id int64) error {
	drug, err := u.GetDrug(ctx, id)
	if err != nil {
		return err
	}

	drug.Discontinue()
	return u.drugRepo.Save(ctx, drug)
}

func (u *DrugUsecase) ActivateDrug(ctx context.Context, id int64) error {
	drug, err := u.GetDrug(ctx, id)
	if err != nil {
		return err
	}

	drug.Activate()
	return u.drugRepo.Save(ctx, drug)
}

func (u *DrugUsecase) DeleteDrug(ctx context.Context, id int64) error {
	return u.drugRepo.Delete(ctx, id)
}