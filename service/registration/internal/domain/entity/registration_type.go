package entity

import "errors"

var (
	ErrInvalidRegistrationID   = errors.New("invalid registration ID")
	ErrInvalidPatientID        = errors.New("invalid patient ID")
	ErrInvalidDoctorID         = errors.New("invalid doctor ID")
	ErrInvalidDepartmentID     = errors.New("invalid department ID")
	ErrInvalidRegistrationType = errors.New("invalid registration type")
	ErrInvalidRegistrationStatus = errors.New("invalid registration status")
)

type RegistrationType int

const (
	RegistrationTypeUnknown RegistrationType = iota
	RegistrationTypeNormal
	RegistrationTypeExpert
)

func (t RegistrationType) String() string {
	switch t {
	case RegistrationTypeNormal:
		return "普通挂号"
	case RegistrationTypeExpert:
		return "专家挂号"
	default:
		return "未知"
	}
}

func RegistrationTypeFromCode(code string) RegistrationType {
	switch code {
	case "NORMAL":
		return RegistrationTypeNormal
	case "EXPERT":
		return RegistrationTypeExpert
	default:
		return RegistrationTypeUnknown
	}
}

func (t RegistrationType) Code() string {
	switch t {
	case RegistrationTypeNormal:
		return "NORMAL"
	case RegistrationTypeExpert:
		return "EXPERT"
	default:
		return "UNKNOWN"
	}
}

type RegistrationStatus int

const (
	RegistrationStatusUnknown RegistrationStatus = iota
	RegistrationStatusPending
	RegistrationStatusConfirmed
	RegistrationStatusCancelled
	RegistrationStatusCompleted
)

func (s RegistrationStatus) String() string {
	switch s {
	case RegistrationStatusPending:
		return "待确认"
	case RegistrationStatusConfirmed:
		return "已确认"
	case RegistrationStatusCancelled:
		return "已取消"
	case RegistrationStatusCompleted:
		return "已完成"
	default:
		return "未知"
	}
}

func RegistrationStatusFromCode(code string) RegistrationStatus {
	switch code {
	case "PENDING":
		return RegistrationStatusPending
	case "CONFIRMED":
		return RegistrationStatusConfirmed
	case "CANCELLED":
		return RegistrationStatusCancelled
	case "COMPLETED":
		return RegistrationStatusCompleted
	default:
		return RegistrationStatusUnknown
	}
}

func (s RegistrationStatus) Code() string {
	switch s {
	case RegistrationStatusPending:
		return "PENDING"
	case RegistrationStatusConfirmed:
		return "CONFIRMED"
	case RegistrationStatusCancelled:
		return "CANCELLED"
	case RegistrationStatusCompleted:
		return "COMPLETED"
	default:
		return "UNKNOWN"
	}
}

func (s RegistrationStatus) CanCancel() bool {
	return s == RegistrationStatusPending || s == RegistrationStatusConfirmed
}

func (s RegistrationStatus) CanConfirm() bool {
	return s == RegistrationStatusPending
}

func (s RegistrationStatus) CanComplete() bool {
	return s == RegistrationStatusConfirmed
}

func (s RegistrationStatus) IsTerminal() bool {
	return s == RegistrationStatusCancelled || s == RegistrationStatusCompleted
}