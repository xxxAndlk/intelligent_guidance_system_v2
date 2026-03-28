package usecase

import (
	"context"
	"errors"

	"intelligent-guidance-system/service/notification/internal/domain/entity"
	"intelligent-guidance-system/service/notification/internal/domain/repository"
)

var (
	ErrNotificationNotFound = errors.New("notification not found")
)

type NotificationUsecase struct {
	repo repository.NotificationRepository
}

func NewNotificationUsecase(repo repository.NotificationRepository) *NotificationUsecase {
	return &NotificationUsecase{repo: repo}
}

func (u *NotificationUsecase) SendNotification(ctx context.Context,
	userID int64,
	notiType entity.NotificationType,
	title string,
	content string,
) (*entity.Notification, error) {
	notification, err := entity.NewNotification(userID, notiType, title, content)
	if err != nil {
		return nil, err
	}

	if err := u.repo.Save(ctx, notification); err != nil {
		return nil, err
	}

	return notification, nil
}

func (u *NotificationUsecase) GetNotification(ctx context.Context, id int64) (*entity.Notification, error) {
	notification, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if notification == nil {
		return nil, ErrNotificationNotFound
	}
	return notification, nil
}

func (u *NotificationUsecase) GetNotifications(ctx context.Context, userID int64) ([]*entity.Notification, error) {
	return u.repo.FindByUserID(ctx, userID)
}

func (u *NotificationUsecase) GetUnreadNotifications(ctx context.Context, userID int64) ([]*entity.Notification, error) {
	return u.repo.FindUnreadByUserID(ctx, userID)
}

func (u *NotificationUsecase) MarkAsRead(ctx context.Context, id int64) error {
	notification, err := u.GetNotification(ctx, id)
	if err != nil {
		return err
	}

	notification.MarkAsRead()
	return u.repo.Save(ctx, notification)
}

func (u *NotificationUsecase) GetUnreadCount(ctx context.Context, userID int64) (int64, error) {
	return u.repo.CountUnreadByUserID(ctx, userID)
}

func (u *NotificationUsecase) ListNotifications(ctx context.Context, page, pageSize int) ([]*entity.Notification, int64, error) {
	return u.repo.FindAll(ctx, page, pageSize)
}

func (u *NotificationUsecase) DeleteNotification(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}