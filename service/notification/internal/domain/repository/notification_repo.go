package repository

import (
	"context"

	"intelligent-guidance-system/service/notification/internal/domain/entity"
)

type NotificationRepository interface {
	Save(ctx context.Context, notification *entity.Notification) error
	FindByID(ctx context.Context, id int64) (*entity.Notification, error)
	FindByUserID(ctx context.Context, userID int64) ([]*entity.Notification, error)
	FindUnreadByUserID(ctx context.Context, userID int64) ([]*entity.Notification, error)
	FindAll(ctx context.Context, page, pageSize int) ([]*entity.Notification, int64, error)
	Delete(ctx context.Context, id int64) error
	DeleteByUserID(ctx context.Context, userID int64) error
	CountUnreadByUserID(ctx context.Context, userID int64) (int64, error)
}