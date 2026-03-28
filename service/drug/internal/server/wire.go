package server

import (
	"github.com/google/wire"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"intelligent-guidance-system/service/drug/internal/biz/usecase"
	"intelligent-guidance-system/service/drug/internal/data/mysql"
	"intelligent-guidance-system/service/drug/internal/domain/repository"
	"intelligent-guidance-system/service/drug/internal/service"
)

var ProviderSet = wire.NewSet(
	NewConfig,
	NewLogger,
	NewGRPCServer,
	NewHTTPServer,
	NewApp,
	NewDB,
	NewDrugRepository,
	NewDrugStockRepository,
	NewDrugUsecase,
	NewDrugService,
)

func NewDB(cfg *Config) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(cfg.Database.Url), &gorm.Config{})
}

func NewDrugRepository(db *gorm.DB) repository.DrugRepository {
	return mysql.NewDrugRepoImpl(db)
}

func NewDrugStockRepository(db *gorm.DB) repository.DrugStockRepository {
	return mysql.NewDrugStockRepoImpl(db)
}

func NewDrugUsecase(drugRepo repository.DrugRepository, stockRepo repository.DrugStockRepository) *usecase.DrugUsecase {
	return usecase.NewDrugUsecase(drugRepo, stockRepo)
}

func NewDrugService(uc *usecase.DrugUsecase) *service.DrugService {
	return service.NewDrugService(uc)
}