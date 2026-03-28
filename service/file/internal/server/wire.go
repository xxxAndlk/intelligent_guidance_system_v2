package server

import (
	"github.com/google/wire"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"intelligent-guidance-system/service/file/internal/biz/usecase"
	"intelligent-guidance-system/service/file/internal/data/mysql"
	"intelligent-guidance-system/service/file/internal/domain/repository"
	"intelligent-guidance-system/service/file/internal/service"
)

var ProviderSet = wire.NewSet(
	NewConfig,
	NewLogger,
	NewGRPCServer,
	NewHTTPServer,
	NewApp,
	NewDB,
	NewFileRepository,
	NewFileUsecase,
	NewFileService,
)

func NewDB(cfg *Config) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(cfg.Database.Url), &gorm.Config{})
}

func NewFileRepository(db *gorm.DB) repository.FileRepository {
	return mysql.NewFileRepoImpl(db)
}

func NewFileUsecase(repo repository.FileRepository) *usecase.FileUsecase {
	return usecase.NewFileUsecase(repo)
}

func NewFileService(uc *usecase.FileUsecase) *service.FileService {
	return service.NewFileService(uc)
}