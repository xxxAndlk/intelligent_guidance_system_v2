package entity

import (
	"time"

	"intelligent_guidance_system_v2/service/medical/internal/domain/vo"
)

// StatusChange represents a status change record entity
type StatusChange struct {
	fromStatus  MedicalStatus // original status
	toStatus    MedicalStatus // new status
	changeTime  time.Time     // time of change
	operatorID  int64         // operator (doctor/staff) ID
	reason      string        // reason for change
}

// NewStatusChange creates a new StatusChange entity
func NewStatusChange(fromStatus, toStatus MedicalStatus, operatorID int64, reason string) *StatusChange {
	return &StatusChange{
		fromStatus: fromStatus,
		toStatus:   toStatus,
		changeTime: time.Now(),
		operatorID: operatorID,
		reason:     reason,
	}
}

// ReconstructStatusChange reconstructs a StatusChange from persistence
func ReconstructStatusChange(fromStatus, toStatus MedicalStatus, changeTime time.Time, operatorID int64, reason string) *StatusChange {
	return &StatusChange{
		fromStatus: fromStatus,
		toStatus:   toStatus,
		changeTime: changeTime,
		operatorID: operatorID,
		reason:     reason,
	}
}

// FromStatus returns the original status
func (s *StatusChange) FromStatus() MedicalStatus {
	return s.fromStatus
}

// ToStatus returns the new status
func (s *StatusChange) ToStatus() MedicalStatus {
	return s.toStatus
}

// ChangeTime returns the change time
func (s *StatusChange) ChangeTime() time.Time {
	return s.changeTime
}

// OperatorID returns the operator ID
func (s *StatusChange) OperatorID() int64 {
	return s.operatorID
}

// Reason returns the reason for change
func (s *StatusChange) Reason() string {
	return s.reason
}

// MedicalStatus represents the status of a medical record
type MedicalStatus int

const (
	StatusUnspecified MedicalStatus = 0
	StatusRegistering MedicalStatus = 1 // 挂号中
	StatusQueuing    MedicalStatus = 2 // 排队中
	StatusDiagnosing MedicalStatus = 3 // 诊断中
	StatusOperating  MedicalStatus = 4 // 手术中
	StatusCompleted  MedicalStatus = 5 // 已完成
	StatusCancelled  MedicalStatus = 6 // 已取消
)

// String returns the string representation of MedicalStatus
func (s MedicalStatus) String() string {
	switch s {
	case StatusUnspecified:
		return "unspecified"
	case StatusRegistering:
		return "registering"
	case StatusQueuing:
		return "queuing"
	case StatusDiagnosing:
		return "diagnosing"
	case StatusOperating:
		return "operating"
	case StatusCompleted:
		return "completed"
	case StatusCancelled:
		return "cancelled"
	default:
		return "unknown"
	}
}

// ChineseName returns the Chinese name of MedicalStatus
func (s MedicalStatus) ChineseName() string {
	switch s {
	case StatusUnspecified:
		return "未知"
	case StatusRegistering:
		return "挂号中"
	case StatusQueuing:
		return "排队中"
	case StatusDiagnosing:
		return "诊断中"
	case StatusOperating:
		return "手术中"
	case StatusCompleted:
		return "已完成"
	case StatusCancelled:
		return "已取消"
	default:
		return "未知"
	}
}

// IsTerminal checks if the status is terminal (cannot change further)
func (s MedicalStatus) IsTerminal() bool {
	return s == StatusCompleted || s == StatusCancelled
}

// IsValid checks if the status is valid
func (s MedicalStatus) IsValid() bool {
	return s >= StatusUnspecified && s <= StatusCancelled
}

// Status transitions map: from -> allowed to
var validTransitions = map[MedicalStatus][]MedicalStatus{
	StatusRegistering: {StatusQueuing, StatusDiagnosing, StatusCancelled},
	StatusQueuing:     {StatusDiagnosing, StatusOperating, StatusCancelled},
	StatusDiagnosing:  {StatusOperating, StatusCompleted, StatusCancelled},
	StatusOperating:   {StatusDiagnosing, StatusCompleted, StatusCancelled},
	StatusCompleted:   {}, // Terminal state
	StatusCancelled:   {}, // Terminal state
}

// IsValidTransition checks if transition from one status to another is valid
func IsValidTransition(from, to MedicalStatus) bool {
	if from == to {
		return true // Same status is allowed (no change)
	}
	allowedTargets, exists := validTransitions[from]
	if !exists {
		return false
	}
	for _, target := range allowedTargets {
		if target == to {
			return true
		}
	}
	return false
}

// GetValidTransitions returns all valid target statuses from a given status
func GetValidTransitions(from MedicalStatus) []MedicalStatus {
	if targets, exists := validTransitions[from]; exists {
		return targets
	}
	return []MedicalStatus{}
}

// PaymentStatus constants for consistency with vo.PaymentStatus
const (
	PaymentStatusUnpaid    = vo.PaymentStatusUnpaid
	PaymentStatusPending   = vo.PaymentStatusPending
	PaymentStatusPaid      = vo.PaymentStatusPaid
	PaymentStatusFailed    = vo.PaymentStatusFailed
	PaymentStatusRefunded  = vo.PaymentStatusRefunded
)

// PaymentMethod constants for consistency with vo.PaymentMethod
const (
	PaymentMethodCash      = vo.PaymentMethodCash
	PaymentMethodCard      = vo.PaymentMethodCard
	PaymentMethodWechat    = vo.PaymentMethodWechat
	PaymentMethodAlipay    = vo.PaymentMethodAlipay
	PaymentMethodInsurance = vo.PaymentMethodInsurance
)