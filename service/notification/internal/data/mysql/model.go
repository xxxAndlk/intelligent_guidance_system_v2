package mysql

import "time"

type NotificationPO struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	UserID    int64     `gorm:"column:user_id;index"`
	Type      string    `gorm:"column:type;size:20"`
	Title     string    `gorm:"column:title;size:200"`
	Content   string    `gorm:"column:content;size:1000"`
	Status    string    `gorm:"column:status;size:20"`
	CreatedAt time.Time `gorm:"column:created_at"`
	ReadAt    time.Time `gorm:"column:read_at"`
}

func (NotificationPO) TableName() string {
	return "notifications"
}