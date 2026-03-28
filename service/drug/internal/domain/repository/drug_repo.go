package repository

import (
	"context"

	"intelligent-guidance-system/service/drug/internal/domain/aggregate"
	"intelligent-guidance-system/service/drug/internal/domain/entity"
)

type DrugRepository interface {
	Save(ctx context.Context, drug *aggregate.Drug) error
	FindByID(ctx context.Context, id int64) (*aggregate.Drug, error)
	FindByName(ctx context.Context, name string) (*aggregate.Drug, error)
	FindByCategory(ctx context.Context, category int) ([]*aggregate.Drug, error)
	FindByStatus(ctx context.Context, status int) ([]*aggregate.Drug, error)
	FindAll(ctx context.Context, page, pageSize int) ([]*aggregate.Drug, int64, error)
	SearchDrugs(ctx context.Context, keyword string) ([]*aggregate.Drug, error)
	Delete(ctx context.Context, id int64) error
	ExistsByName(ctx context.Context, name string) (bool, error)
}

type DrugStockRepository interface {
	Save(ctx context.Context, stock *entity.DrugStock) error
	FindByID(ctx context.Context, id int64) (*entity.DrugStock, error)
	FindByDrugID(ctx context.Context, drugID int64) ([]*entity.DrugStock, error)
	FindAvailableStocks(ctx context.Context, drugID int64) ([]*entity.DrugStock, error)
	FindExpiringSoon(ctx context.Context, days int) ([]*entity.DrugStock, error)
	Delete(ctx context.Context, id int64) error
}