package dto

import "time"

type SendNotificationRequest struct {
	UserID  int64  `json:"user_id" binding:"required"`
	Type    string `json:"type" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
}

type MarkAsReadRequest struct {
	ID int64 `json:"id" binding:"required"`
}

type NotificationResponse struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Type      string    `json:"type"`
	TypeName  string    `json:"type_name"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Status    string    `json:"status"`
	StatusName string   `json:"status_name"`
	CreatedAt time.Time `json:"created_at"`
	ReadAt    time.Time `json:"read_at"`
}

type NotificationListResponse struct {
	Total        int64                 `json:"total"`
	Page         int                   `json:"page"`
	PageSize     int                   `json:"page_size"`
	Notifications []NotificationResponse `json:"notifications"`
}

type UnreadCountResponse struct {
	Count int64 `json:"count"`
}