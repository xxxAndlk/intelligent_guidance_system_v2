package entity

import "errors"

var (
	ErrInvalidSurgicalID   = errors.New("invalid surgical ID")
	ErrInvalidMedicalID    = errors.New("invalid medical ID")
	ErrInvalidPatientID    = errors.New("invalid patient ID")
	ErrInvalidDepartmentID = errors.New("invalid department ID")
	ErrInvalidSurgicalType = errors.New("invalid surgical type")
	ErrInvalidSurgicalStatus = errors.New("invalid surgical status")
)

type SurgicalType int

const (
	SurgicalTypeUnknown SurgicalType = iota
	SurgicalTypeMinor
	SurgicalTypeMajor
	SurgicalTypeEmergency
)

func (t SurgicalType) String() string {
	switch t {
	case SurgicalTypeMinor:
		return "小手术"
	case SurgicalTypeMajor:
		return "大手术"
	case SurgicalTypeEmergency:
		return "急诊手术"
	default:
		return "未知"
	}
}

func SurgicalTypeFromCode(code string) SurgicalType {
	switch code {
	case "MINOR":
		return SurgicalTypeMinor
	case "MAJOR":
		return SurgicalTypeMajor
	case "EMERGENCY":
		return SurgicalTypeEmergency
	default:
		return SurgicalTypeUnknown
	}
}

func (t SurgicalType) Code() string {
	switch t {
	case SurgicalTypeMinor:
		return "MINOR"
	case SurgicalTypeMajor:
		return "MAJOR"
	case SurgicalTypeEmergency:
		return "EMERGENCY"
	default:
		return "UNKNOWN"
	}
}

type SurgicalStatus int

const (
	SurgicalStatusUnknown SurgicalStatus = iota
	SurgicalStatusScheduled
	SurgicalStatusPreparing
	SurgicalStatusInProgress
	SurgicalStatusCompleted
	SurgicalStatusCancelled
)

func (s SurgicalStatus) String() string {
	switch s {
	case SurgicalStatusScheduled:
		return "已排期"
	case SurgicalStatusPreparing:
		return "准备中"
	case SurgicalStatusInProgress:
		return "进行中"
	case SurgicalStatusCompleted:
		return "已完成"
	case SurgicalStatusCancelled:
		return "已取消"
	default:
		return "未知"
	}
}

func SurgicalStatusFromCode(code string) SurgicalStatus {
	switch code {
	case "SCHEDULED":
		return SurgicalStatusScheduled
	case "PREPARING":
		return SurgicalStatusPreparing
	case "IN_PROGRESS":
		return SurgicalStatusInProgress
	case "COMPLETED":
		return SurgicalStatusCompleted
	case "CANCELLED":
		return SurgicalStatusCancelled
	default:
		return SurgicalStatusUnknown
	}
}

func (s SurgicalStatus) Code() string {
	switch s {
	case SurgicalStatusScheduled:
		return "SCHEDULED"
	case SurgicalStatusPreparing:
		return "PREPARING"
	case SurgicalStatusInProgress:
		return "IN_PROGRESS"
	case SurgicalStatusCompleted:
		return "COMPLETED"
	case SurgicalStatusCancelled:
		return "CANCELLED"
	default:
		return "UNKNOWN"
	}
}

func (s SurgicalStatus) CanStart() bool {
	return s == SurgicalStatusScheduled || s == SurgicalStatusPreparing
}

func (s SurgicalStatus) CanCancel() bool {
	return s == SurgicalStatusScheduled || s == SurgicalStatusPreparing
}

func (s SurgicalStatus) IsTerminal() bool {
	return s == SurgicalStatusCompleted || s == SurgicalStatusCancelled
}

type FlowStepStatus int

const (
	FlowStepStatusUnknown FlowStepStatus = iota
	FlowStepStatusPending
	FlowStepStatusInProgress
	FlowStepStatusCompleted
	FlowStepStatusSkipped
)

func (s FlowStepStatus) String() string {
	switch s {
	case FlowStepStatusPending:
		return "待执行"
	case FlowStepStatusInProgress:
		return "执行中"
	case FlowStepStatusCompleted:
		return "已完成"
	case FlowStepStatusSkipped:
		return "已跳过"
	default:
		return "未知"
	}
}

func FlowStepStatusFromCode(code string) FlowStepStatus {
	switch code {
	case "PENDING":
		return FlowStepStatusPending
	case "IN_PROGRESS":
		return FlowStepStatusInProgress
	case "COMPLETED":
		return FlowStepStatusCompleted
	case "SKIPPED":
		return FlowStepStatusSkipped
	default:
		return FlowStepStatusUnknown
	}
}

func (s FlowStepStatus) Code() string {
	switch s {
	case FlowStepStatusPending:
		return "PENDING"
	case FlowStepStatusInProgress:
		return "IN_PROGRESS"
	case FlowStepStatusCompleted:
		return "COMPLETED"
	case FlowStepStatusSkipped:
		return "SKIPPED"
	default:
		return "UNKNOWN"
	}
}

func (s FlowStepStatus) CanStart() bool {
	return s == FlowStepStatusPending
}

func (s FlowStepStatus) IsTerminal() bool {
	return s == FlowStepStatusCompleted || s == FlowStepStatusSkipped
}