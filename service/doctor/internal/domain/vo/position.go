package vo

// Position represents the professional position level of a doctor
type Position int

const (
	PositionUnknown Position = iota
	PositionChiefPhysician      // 主任医师
	PositionAssociateChief      // 副主任医师
	PositionAttending           // 主治医师
	PositionResident            // 住院医师
	PositionAssistant           // 医师助理
)

// String returns the Chinese display name for the position
func (p Position) String() string {
	switch p {
	case PositionChiefPhysician:
		return "主任医师"
	case PositionAssociateChief:
		return "副主任医师"
	case PositionAttending:
		return "主治医师"
	case PositionResident:
		return "住院医师"
	case PositionAssistant:
		return "医师助理"
	default:
		return "未知"
	}
}

// Code returns the code for the position
func (p Position) Code() string {
	switch p {
	case PositionChiefPhysician:
		return "CHIEF"
	case PositionAssociateChief:
		return "ASSOCIATE_CHIEF"
	case PositionAttending:
		return "ATTENDING"
	case PositionResident:
		return "RESIDENT"
	case PositionAssistant:
		return "ASSISTANT"
	default:
		return "UNKNOWN"
	}
}

// PositionFromCode creates Position from code string
func PositionFromCode(code string) Position {
	switch code {
	case "CHIEF":
		return PositionChiefPhysician
	case "ASSOCIATE_CHIEF":
		return PositionAssociateChief
	case "ATTENDING":
		return PositionAttending
	case "RESIDENT":
		return PositionResident
	case "ASSISTANT":
		return PositionAssistant
	default:
		return PositionUnknown
	}
}

// IsExpert returns true if the position is at expert level (chief or associate chief)
func (p Position) IsExpert() bool {
	return p == PositionChiefPhysician || p == PositionAssociateChief
}

// CanTreat returns true if the position can independently treat patients
func (p Position) CanTreat() bool {
	return p >= PositionResident && p <= PositionChiefPhysician
}