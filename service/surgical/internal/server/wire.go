package server

import (
	"github.com/google/wire"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"intelligent-guidance-system/service/surgical/internal/biz/usecase"
	"intelligent-guidance-system/service/surgical/internal/data/mysql"
	"intelligent-guidance-system/service/surgical/internal/domain/repository"
	"intelligent-guidance-system/service/surgical/internal/service"
)

var ProviderSet = wire.NewSet(
	NewConfig,
	NewLogger,
	NewGRPCServer,
	NewHTTPServer,
	NewApp,
	NewDB,
	NewSurgicalRepository,
	NewSurgicalFlowRepository,
	NewSurgicalUsecase,
	NewSurgicalService,
)

func NewDB(cfg *Config) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(cfg.Database.Url), &gorm.Config{})
}

func NewSurgicalRepository(db *gorm.DB) repository.SurgicalRepository {
	return mysql.NewSurgicalRepoImpl(db)
}

func NewSurgicalFlowRepository(db *gorm.DB) repository.SurgicalFlowRepository {
	return mysql.NewSurgicalFlowRepoImpl(db)
}

func NewSurgicalUsecase(surgicalRepo repository.SurgicalRepository, flowRepo repository.SurgicalFlowRepository) *usecase.SurgicalUsecase {
	return usecase.NewSurgicalUsecase(surgicalRepo, flowRepo)
}

func NewSurgicalService(uc *usecase.SurgicalUsecase) *service.SurgicalService {
	return service.NewSurgicalService(uc)
}