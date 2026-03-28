// service/patient/internal/domain/entity/allergy.go
package entity

import "time"

// AllergyID 过敏史ID
type AllergyID int64

// AllergySeverity 过敏严重程度
type AllergySeverity int32

const (
	AllergySeverityMild     AllergySeverity = 1 // 轻度
	AllergySeverityModerate AllergySeverity = 2 // 中度
	AllergySeveritySevere   AllergySeverity = 3 // 重度
)

func (s AllergySeverity) String() string {
	switch s {
	case AllergySeverityMild:
		return "轻度"
	case AllergySeverityModerate:
		return "中度"
	case AllergySeveritySevere:
		return "重度"
	default:
		return "未知"
	}
}

// Allergy 过敏史实体
type Allergy struct {
	ID        AllergyID
	Allergen  string        // 过敏源
	Severity  AllergySeverity // 严重程度
	Reaction  string        // 过敏反应描述
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewAllergy 创建过敏史
func NewAllergy(allergen string, severity AllergySeverity, reaction string) Allergy {
	now := time.Now()
	return Allergy{
		ID:        AllergyID(now.UnixNano()),
		Allergen:  allergen,
		Severity:  severity,
		Reaction:  reaction,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// UpdateReaction 更新过敏反应描述
func (a *Allergy) UpdateReaction(reaction string) {
	a.Reaction = reaction
	a.UpdatedAt = time.Now()
}

// UpdateSeverity 更新严重程度
func (a *Allergy) UpdateSeverity(severity AllergySeverity) {
	a.Severity = severity
	a.UpdatedAt = time.Now()
}