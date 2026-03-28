package server

import (
	"github.com/google/wire"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"intelligent-guidance-system/service/registration/internal/biz/usecase"
	"intelligent-guidance-system/service/registration/internal/data/mysql"
	"intelligent-guidance-system/service/registration/internal/domain/repository"
	"intelligent-guidance-system/service/registration/internal/service"
)

var ProviderSet = wire.NewSet(
	NewConfig,
	NewLogger,
	NewGRPCServer,
	NewHTTPServer,
	NewApp,
	NewDB,
	NewRegistrationRepository,
	NewRegistrationUsecase,
	NewRegistrationService,
)

func NewDB(cfg *Config) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(cfg.Database.Url), &gorm.Config{})
}

func NewRegistrationRepository(db *gorm.DB) repository.RegistrationRepository {
	return mysql.NewRegistrationRepoImpl(db)
}

func NewRegistrationUsecase(repo repository.RegistrationRepository) *usecase.RegistrationUsecase {
	return usecase.NewRegistrationUsecase(repo)
}

func NewRegistrationService(uc *usecase.RegistrationUsecase) *service.RegistrationService {
	return service.NewRegistrationService(uc)
}