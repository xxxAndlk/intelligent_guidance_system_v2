package server

import (
	"github.com/google/wire"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"intelligent-guidance-system/service/department/internal/biz/usecase"
	"intelligent-guidance-system/service/department/internal/data/mysql"
	"intelligent-guidance-system/service/department/internal/domain/repository"
	"intelligent-guidance-system/service/department/internal/service"
)

var ProviderSet = wire.NewSet(
	NewConfig,
	NewLogger,
	NewGRPCServer,
	NewHTTPServer,
	NewApp,
	NewDB,
	NewDepartmentRepository,
	NewDepartmentUsecase,
	NewDepartmentService,
)

func NewDB(cfg *Config) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(cfg.Database.Url), &gorm.Config{})
}

func NewDepartmentRepository(db *gorm.DB) repository.DepartmentRepository {
	return mysql.NewDepartmentRepoImpl(db)
}

func NewDepartmentUsecase(repo repository.DepartmentRepository) *usecase.DepartmentUsecase {
	return usecase.NewDepartmentUsecase(repo)
}

func NewDepartmentService(uc *usecase.DepartmentUsecase) *service.DepartmentService {
	return service.NewDepartmentService(uc)
}