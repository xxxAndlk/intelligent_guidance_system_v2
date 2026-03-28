package service

import (
	"context"

	"intelligent-guidance-system/service/notification/internal/biz/dto"
	"intelligent-guidance-system/service/notification/internal/biz/usecase"
	"intelligent-guidance-system/service/notification/internal/domain/entity"
)

type NotificationService struct {
	usecase *usecase.NotificationUsecase
}

func NewNotificationService(usecase *usecase.NotificationUsecase) *NotificationService {
	return &NotificationService{usecase: usecase}
}

func (s *NotificationService) SendNotification(ctx context.Context, req *dto.SendNotificationRequest) (*dto.NotificationResponse, error) {
	notiType := entity.NotificationTypeFromCode(req.Type)
	
	notification, err := s.usecase.SendNotification(ctx, req.UserID, notiType, req.Title, req.Content)
	if err != nil {
		return nil, err
	}

	return s.toResponse(notification), nil
}

func (s *NotificationService) GetNotification(ctx context.Context, id int64) (*dto.NotificationResponse, error) {
	notification, err := s.usecase.GetNotification(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(notification), nil
}

func (s *NotificationService) GetNotifications(ctx context.Context, userID int64) ([]*dto.NotificationResponse, error) {
	notifications, err := s.usecase.GetNotifications(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.toResponses(notifications), nil
}

func (s *NotificationService) GetUnreadNotifications(ctx context.Context, userID int64) ([]*dto.NotificationResponse, error) {
	notifications, err := s.usecase.GetUnreadNotifications(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.toResponses(notifications), nil
}

func (s *NotificationService) MarkAsRead(ctx context.Context, id int64) error {
	return s.usecase.MarkAsRead(ctx, id)
}

func (s *NotificationService) GetUnreadCount(ctx context.Context, userID int64) (*dto.UnreadCountResponse, error) {
	count, err := s.usecase.GetUnreadCount(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &dto.UnreadCountResponse{Count: count}, nil
}

func (s *NotificationService) ListNotifications(ctx context.Context, page, pageSize int) (*dto.NotificationListResponse, error) {
	notifications, total, err := s.usecase.ListNotifications(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	return &dto.NotificationListResponse{
		Total:         total,
		Page:          page,
		PageSize:      pageSize,
		Notifications: s.toResponseList(notifications),
	}, nil
}

func (s *NotificationService) DeleteNotification(ctx context.Context, id int64) error {
	return s.usecase.DeleteNotification(ctx, id)
}

func (s *NotificationService) toResponse(n *entity.Notification) *dto.NotificationResponse {
	return &dto.NotificationResponse{
		ID:         n.ID(),
		UserID:     n.UserID(),
		Type:       n.Type().Code(),
		TypeName:   n.Type().String(),
		Title:      n.Title(),
		Content:    n.Content(),
		Status:     n.Status().Code(),
		StatusName: n.Status().String(),
		CreatedAt:  n.CreatedAt(),
		ReadAt:     n.ReadAt(),
	}
}

func (s *NotificationService) toResponses(notifications []*entity.Notification) []*dto.NotificationResponse {
	responses := make([]*dto.NotificationResponse, 0, len(notifications))
	for _, n := range notifications {
		responses = append(responses, s.toResponse(n))
	}
	return responses
}

func (s *NotificationService) toResponseList(notifications []*entity.Notification) []dto.NotificationResponse {
	responses := make([]dto.NotificationResponse, 0, len(notifications))
	for _, n := range notifications {
		responses = append(responses, *s.toResponse(n))
	}
	return responses
}