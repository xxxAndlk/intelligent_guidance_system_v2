package entity

import (
	"errors"
	"time"
)

var (
	ErrInvalidNotificationID = errors.New("invalid notification ID")
	ErrInvalidUserID         = errors.New("invalid user ID")
	ErrNotificationNotFound  = errors.New("notification not found")
)

type NotificationType int

const (
	NotificationTypeUnknown NotificationType = iota
	NotificationTypeSystem
	NotificationTypeAppointment
	NotificationTypePayment
	NotificationTypeMedical
)

func (t NotificationType) String() string {
	switch t {
	case NotificationTypeSystem:
		return "系统通知"
	case NotificationTypeAppointment:
		return "预约通知"
	case NotificationTypePayment:
		return "支付通知"
	case NotificationTypeMedical:
		return "医疗通知"
	default:
		return "未知"
	}
}

func NotificationTypeFromCode(code string) NotificationType {
	switch code {
	case "SYSTEM":
		return NotificationTypeSystem
	case "APPOINTMENT":
		return NotificationTypeAppointment
	case "PAYMENT":
		return NotificationTypePayment
	case "MEDICAL":
		return NotificationTypeMedical
	default:
		return NotificationTypeUnknown
	}
}

func (t NotificationType) Code() string {
	switch t {
	case NotificationTypeSystem:
		return "SYSTEM"
	case NotificationTypeAppointment:
		return "APPOINTMENT"
	case NotificationTypePayment:
		return "PAYMENT"
	case NotificationTypeMedical:
		return "MEDICAL"
	default:
		return "UNKNOWN"
	}
}

type NotificationStatus int

const (
	NotificationStatusUnknown NotificationStatus = iota
	NotificationStatusUnread
	NotificationStatusRead
)

func (s NotificationStatus) String() string {
	switch s {
	case NotificationStatusUnread:
		return "未读"
	case NotificationStatusRead:
		return "已读"
	default:
		return "未知"
	}
}

func NotificationStatusFromCode(code string) NotificationStatus {
	switch code {
	case "UNREAD":
		return NotificationStatusUnread
	case "READ":
		return NotificationStatusRead
	default:
		return NotificationStatusUnknown
	}
}

func (s NotificationStatus) Code() string {
	switch s {
	case NotificationStatusUnread:
		return "UNREAD"
	case NotificationStatusRead:
		return "READ"
	default:
		return "UNKNOWN"
	}
}

type Notification struct {
	id        int64
	userID    int64
	notiType  NotificationType
	title     string
	content   string
	status    NotificationStatus
	createdAt time.Time
	readAt    time.Time
}

func NewNotification(userID int64, notiType NotificationType, title, content string) (*Notification, error) {
	if userID <= 0 {
		return nil, ErrInvalidUserID
	}
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}

	return &Notification{
		id:        0,
		userID:    userID,
		notiType:  notiType,
		title:     title,
		content:   content,
		status:    NotificationStatusUnread,
		createdAt: time.Now(),
	}, nil
}

func ReconstructNotification(id, userID int64, notiType NotificationType, title, content string, status NotificationStatus, createdAt, readAt time.Time) *Notification {
	return &Notification{
		id:        id,
		userID:    userID,
		notiType:  notiType,
		title:     title,
		content:   content,
		status:    status,
		createdAt: createdAt,
		readAt:    readAt,
	}
}

func (n *Notification) ID() int64 { return n.id }
func (n *Notification) UserID() int64 { return n.userID }
func (n *Notification) Type() NotificationType { return n.notiType }
func (n *Notification) Title() string { return n.title }
func (n *Notification) Content() string { return n.content }
func (n *Notification) Status() NotificationStatus { return n.status }
func (n *Notification) CreatedAt() time.Time { return n.createdAt }
func (n *Notification) ReadAt() time.Time { return n.readAt }

func (n *Notification) SetID(id int64) { n.id = id }

func (n *Notification) MarkAsRead() {
	n.status = NotificationStatusRead
	n.readAt = time.Now()
}

func (n *Notification) IsRead() bool {
	return n.status == NotificationStatusRead
}