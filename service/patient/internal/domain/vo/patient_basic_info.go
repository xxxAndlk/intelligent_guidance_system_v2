// service/patient/internal/domain/vo/patient_basic_info.go
package vo

import (
	"errors"
	"time"
)

// Gender 性别枚举
type Gender int32

const (
	GenderUnknown Gender = 0
	GenderMale    Gender = 1
	GenderFemale  Gender = 2
)

func (g Gender) String() string {
	switch g {
	case GenderMale:
		return "男"
	case GenderFemale:
		return "女"
	default:
		return "未知"
	}
}

// PatientBasicInfo 患者基本信息值对象
type PatientBasicInfo struct {
	Name      string
	IDCard    IDCard
	Age       int32
	Gender    Gender
	BirthDate time.Time
}

// NewPatientBasicInfo 创建患者基本信息
func NewPatientBasicInfo(name string, idCard IDCard) PatientBasicInfo {
	// 从身份证自动提取性别和出生日期
	gender, _ := idCard.GetGender()
	birthDate, _ := idCard.GetBirthDate()
	age, _ := idCard.GetAge()

	return PatientBasicInfo{
		Name:      name,
		IDCard:    idCard,
		Age:       age,
		Gender:    gender,
		BirthDate: birthDate,
	}
}

// Validate 验证基本信息
func (p PatientBasicInfo) Validate() error {
	if p.Name == "" {
		return errors.New("name is required")
	}
	if len(p.Name) > 50 {
		return errors.New("name length exceeds limit")
	}
	if p.Age < 0 || p.Age > 150 {
		return errors.New("invalid age range")
	}
	if err := p.IDCard.Validate(); err != nil {
		return err
	}
	return nil
}

// UpdateName 更新姓名
func (p *PatientBasicInfo) UpdateName(name string) error {
	if name == "" {
		return errors.New("name is required")
	}
	if len(name) > 50 {
		return errors.New("name length exceeds limit")
	}
	p.Name = name
	return nil
}

// GetDisplayName 获取显示名称（用于前端展示）
func (p PatientBasicInfo) GetDisplayName() string {
	return p.Name
}

// GetMaskedIDCard 获取脱敏身份证
func (p PatientBasicInfo) GetMaskedIDCard() string {
	return p.IDCard.Masked()
}