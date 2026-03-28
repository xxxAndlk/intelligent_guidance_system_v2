package vo

// Title represents the administrative title of a doctor
type Title int

const (
	TitleUnknown Title = iota
	TitleDean          // 院长
	TitleViceDean      // 副院长
	TitleDeptHead      // 科室主任
	TitleMedicalLeader // 医疗组长
	TitleMedicalMember // 医疗组成员
)

// String returns the Chinese display name for the title
func (t Title) String() string {
	switch t {
	case TitleDean:
		return "院长"
	case TitleViceDean:
		return "副院长"
	case TitleDeptHead:
		return "科室主任"
	case TitleMedicalLeader:
		return "医疗组长"
	case TitleMedicalMember:
		return "医疗组成员"
	default:
		return "未知"
	}
}

// Code returns the code for the title
func (t Title) Code() string {
	switch t {
	case TitleDean:
		return "DEAN"
	case TitleViceDean:
		return "VICE_DEAN"
	case TitleDeptHead:
		return "DEPT_HEAD"
	case TitleMedicalLeader:
		return "MEDICAL_LEADER"
	case TitleMedicalMember:
		return "MEDICAL_MEMBER"
	default:
		return "UNKNOWN"
	}
}

// TitleFromCode creates Title from code string
func TitleFromCode(code string) Title {
	switch code {
	case "DEAN":
		return TitleDean
	case "VICE_DEAN":
		return TitleViceDean
	case "DEPT_HEAD":
		return TitleDeptHead
	case "MEDICAL_LEADER":
		return TitleMedicalLeader
	case "MEDICAL_MEMBER":
		return TitleMedicalMember
	default:
		return TitleUnknown
	}
}

// IsManagement returns true if the title is a management position
func (t Title) IsManagement() bool {
	return t == TitleDean || t == TitleViceDean || t == TitleDeptHead
}