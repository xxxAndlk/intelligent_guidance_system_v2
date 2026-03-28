package server

import (
	"github.com/google/wire"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"intelligent-guidance-system/service/notification/internal/biz/usecase"
	"intelligent-guidance-system/service/notification/internal/data/mysql"
	"intelligent-guidance-system/service/notification/internal/domain/repository"
	"intelligent-guidance-system/service/notification/internal/service"
)

var ProviderSet = wire.NewSet(
	NewConfig,
	NewLogger,
	NewGRPCServer,
	NewHTTPServer,
	NewApp,
	NewDB,
	NewNotificationRepository,
	NewNotificationUsecase,
	NewNotificationService,
)

func NewDB(cfg *Config) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(cfg.Database.Url), &gorm.Config{})
}

func NewNotificationRepository(db *gorm.DB) repository.NotificationRepository {
	return mysql.NewNotificationRepoImpl(db)
}

func NewNotificationUsecase(repo repository.NotificationRepository) *usecase.NotificationUsecase {
	return usecase.NewNotificationUsecase(repo)
}

func NewNotificationService(uc *usecase.NotificationUsecase) *service.NotificationService {
	return service.NewNotificationService(uc)
}