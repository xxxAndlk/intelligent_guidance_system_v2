package mysql

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"intelligent-guidance-system/service/notification/internal/domain/entity"
	"intelligent-guidance-system/service/notification/internal/domain/repository"
)

type NotificationRepoImpl struct {
	db *gorm.DB
}

func NewNotificationRepoImpl(db *gorm.DB) repository.NotificationRepository {
	return &NotificationRepoImpl{db: db}
}

func (r *NotificationRepoImpl) Save(ctx context.Context, notification *entity.Notification) error {
	po := r.toPO(notification)

	if po.ID == 0 {
		if err := r.db.WithContext(ctx).Create(&po).Error; err != nil {
			return err
		}
		notification.SetID(po.ID)
	} else {
		if err := r.db.WithContext(ctx).Save(&po).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *NotificationRepoImpl) FindByID(ctx context.Context, id int64) (*entity.Notification, error) {
	var po NotificationPO
	if err := r.db.WithContext(ctx).First(&po, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.toEntity(&po), nil
}

func (r *NotificationRepoImpl) FindByUserID(ctx context.Context, userID int64) ([]*entity.Notification, error) {
	var pos []NotificationPO
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&pos).Error; err != nil {
		return nil, err
	}
	return r.toEntities(pos), nil
}

func (r *NotificationRepoImpl) FindUnreadByUserID(ctx context.Context, userID int64) ([]*entity.Notification, error) {
	var pos []NotificationPO
	if err := r.db.WithContext(ctx).Where("user_id = ? AND status = ?", userID, "UNREAD").Order("created_at DESC").Find(&pos).Error; err != nil {
		return nil, err
	}
	return r.toEntities(pos), nil
}

func (r *NotificationRepoImpl) FindAll(ctx context.Context, page, pageSize int) ([]*entity.Notification, int64, error) {
	var pos []NotificationPO
	var total int64

	if err := r.db.WithContext(ctx).Model(&NotificationPO{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	return r.toEntities(pos), total, nil
}

func (r *NotificationRepoImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&NotificationPO{}, id).Error
}

func (r *NotificationRepoImpl) DeleteByUserID(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&NotificationPO{}).Error
}

func (r *NotificationRepoImpl) CountUnreadByUserID(ctx context.Context, userID int64) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&NotificationPO{}).Where("user_id = ? AND status = ?", userID, "UNREAD").Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *NotificationRepoImpl) toPO(n *entity.Notification) *NotificationPO {
	return &NotificationPO{
		ID:        n.ID(),
		UserID:    n.UserID(),
		Type:      n.Type().Code(),
		Title:     n.Title(),
		Content:   n.Content(),
		Status:    n.Status().Code(),
		CreatedAt: n.CreatedAt(),
		ReadAt:    n.ReadAt(),
	}
}

func (r *NotificationRepoImpl) toEntity(po *NotificationPO) *entity.Notification {
	return entity.ReconstructNotification(
		po.ID,
		po.UserID,
		entity.NotificationTypeFromCode(po.Type),
		po.Title,
		po.Content,
		entity.NotificationStatusFromCode(po.Status),
		po.CreatedAt,
		po.ReadAt,
	)
}

func (r *NotificationRepoImpl) toEntities(pos []NotificationPO) []*entity.Notification {
	entities := make([]*entity.Notification, 0, len(pos))
	for _, po := range pos {
		entities = append(entities, r.toEntity(&po))
	}
	return entities
}