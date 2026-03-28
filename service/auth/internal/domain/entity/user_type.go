package entity

type UserType int

const (
	UserTypeUnknown UserType = iota
	UserTypePatient
	UserTypeDoctor
	UserTypeAdmin
)

func (t UserType) String() string {
	switch t {
	case UserTypePatient:
		return "患者"
	case UserTypeDoctor:
		return "医生"
	case UserTypeAdmin:
		return "管理员"
	default:
		return "未知"
	}
}

func UserTypeFromCode(code string) UserType {
	switch code {
	case "PATIENT":
		return UserTypePatient
	case "DOCTOR":
		return UserTypeDoctor
	case "ADMIN":
		return UserTypeAdmin
	default:
		return UserTypeUnknown
	}
}

func (t UserType) Code() string {
	switch t {
	case UserTypePatient:
		return "PATIENT"
	case UserTypeDoctor:
		return "DOCTOR"
	case UserTypeAdmin:
		return "ADMIN"
	default:
		return "UNKNOWN"
	}
}

type DataScope int

const (
	DataScopeUnknown DataScope = iota
	DataScopeAll
	DataScopeDept
	DataScopeDeptAndSub
	DataScopeOwn
)

func (s DataScope) String() string {
	switch s {
	case DataScopeAll:
		return "全部数据"
	case DataScopeDept:
		return "本部门数据"
	case DataScopeDeptAndSub:
		return "本部门及子部门数据"
	case DataScopeOwn:
		return "仅本人数据"
	default:
		return "未知"
	}
}

func DataScopeFromCode(code string) DataScope {
	switch code {
	case "ALL":
		return DataScopeAll
	case "DEPT":
		return DataScopeDept
	case "DEPT_AND_SUB":
		return DataScopeDeptAndSub
	case "OWN":
		return DataScopeOwn
	default:
		return DataScopeUnknown
	}
}

func (s DataScope) Code() string {
	switch s {
	case DataScopeAll:
		return "ALL"
	case DataScopeDept:
		return "DEPT"
	case DataScopeDeptAndSub:
		return "DEPT_AND_SUB"
	case DataScopeOwn:
		return "OWN"
	default:
		return "UNKNOWN"
	}
}

type UserStatus int

const (
	UserStatusUnknown UserStatus = iota
	UserStatusActive
	UserStatusInactive
	UserStatusLocked
	UserStatusDisabled
)

func (s UserStatus) String() string {
	switch s {
	case UserStatusActive:
		return "正常"
	case UserStatusInactive:
		return "未激活"
	case UserStatusLocked:
		return "锁定"
	case UserStatusDisabled:
		return "禁用"
	default:
		return "未知"
	}
}

func UserStatusFromCode(code string) UserStatus {
	switch code {
	case "ACTIVE":
		return UserStatusActive
	case "INACTIVE":
		return UserStatusInactive
	case "LOCKED":
		return UserStatusLocked
	case "DISABLED":
		return UserStatusDisabled
	default:
		return UserStatusUnknown
	}
}

func (s UserStatus) Code() string {
	switch s {
	case UserStatusActive:
		return "ACTIVE"
	case UserStatusInactive:
		return "INACTIVE"
	case UserStatusLocked:
		return "LOCKED"
	case UserStatusDisabled:
		return "DISABLED"
	default:
		return "UNKNOWN"
	}
}

func (s UserStatus) CanLogin() bool {
	return s == UserStatusActive
}

func (s UserStatus) IsTerminal() bool {
	return s == UserStatusDisabled
}