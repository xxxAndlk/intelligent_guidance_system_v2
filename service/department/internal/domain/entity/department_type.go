package entity

// DepartmentType represents the type of department
type DepartmentType int

const (
	DepartmentTypeUnknown DepartmentType = iota
	DepartmentTypeOutpatient
	DepartmentTypeInpatient
)

func (t DepartmentType) String() string {
	switch t {
	case DepartmentTypeOutpatient:
		return "门诊"
	case DepartmentTypeInpatient:
		return "住院"
	default:
		return "未知"
	}
}

func DepartmentTypeFromCode(code string) DepartmentType {
	switch code {
	case "OUTPATIENT":
		return DepartmentTypeOutpatient
	case "INPATIENT":
		return DepartmentTypeInpatient
	default:
		return DepartmentTypeUnknown
	}
}

func (t DepartmentType) Code() string {
	switch t {
	case DepartmentTypeOutpatient:
		return "OUTPATIENT"
	case DepartmentTypeInpatient:
		return "INPATIENT"
	default:
		return "UNKNOWN"
	}
}

// DepartmentStatus represents the status of department
type DepartmentStatus int

const (
	DepartmentStatusUnknown DepartmentStatus = iota
	DepartmentStatusActive
	DepartmentStatusResting
	DepartmentStatusDisabled
)

func (s DepartmentStatus) String() string {
	switch s {
	case DepartmentStatusActive:
		return "正常"
	case DepartmentStatusResting:
		return "休息"
	case DepartmentStatusDisabled:
		return "停用"
	default:
		return "未知"
	}
}

func DepartmentStatusFromCode(code string) DepartmentStatus {
	switch code {
	case "ACTIVE":
		return DepartmentStatusActive
	case "RESTING":
		return DepartmentStatusResting
	case "DISABLED":
		return DepartmentStatusDisabled
	default:
		return DepartmentStatusUnknown
	}
}

func (s DepartmentStatus) Code() string {
	switch s {
	case DepartmentStatusActive:
		return "ACTIVE"
	case DepartmentStatusResting:
		return "RESTING"
	case DepartmentStatusDisabled:
		return "DISABLED"
	default:
		return "UNKNOWN"
	}
}

func (s DepartmentStatus) IsActive() bool {
	return s == DepartmentStatusActive
}

func (s DepartmentStatus) CanRegister() bool {
	return s == DepartmentStatusActive
}

func (s DepartmentStatus) IsTerminal() bool {
	return s == DepartmentStatusDisabled
}