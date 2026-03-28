// service/patient/internal/domain/aggregate/patient.go
package aggregate

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"intelligent-guidance-system/service/patient/internal/domain/entity"
	"intelligent-guidance-system/service/patient/internal/domain/event"
	"intelligent-guidance-system/service/patient/internal/domain/vo"
)

// PatientID 患者ID
type PatientID int64

// PatientStatus 患者状态
type PatientStatus int32

const (
	PatientStatusNormal   PatientStatus = 1 // 正常
	PatientStatusLocked   PatientStatus = 2 // 锁定
	PatientStatusInactive PatientStatus = 3 // 未激活
)

func (s PatientStatus) String() string {
	switch s {
	case PatientStatusNormal:
		return "正常"
	case PatientStatusLocked:
		return "锁定"
	case PatientStatusInactive:
		return "未激活"
	default:
		return "未知"
	}
}

// Patient 患者聚合根
type Patient struct {
	ID        PatientID
	OpenID    string
	Username  string
	Password  string
	BasicInfo vo.PatientBasicInfo
	Contact   vo.ContactInfo
	Status    PatientStatus

	// 关联实体
	MedicalHistory []entity.MedicalHistory
	Allergies      []entity.Allergy

	// 元数据
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	// 领域事件
	Events []event.DomainEvent
}

// 错误定义
var (
	ErrPatientLocked     = errors.New("patient is locked")
	ErrPatientInactive   = errors.New("patient is inactive")
	ErrInvalidInfo       = errors.New("invalid patient info")
	ErrDuplicateAllergy  = errors.New("duplicate allergy")
	ErrDuplicateHistory  = errors.New("duplicate medical history")
)

// NewPatient 创建新患者
func NewPatient(username, password string, phone vo.Phone, name string, idCard vo.IDCard) *Patient {
	now := time.Now()
	return &Patient{
		ID:        PatientID(0), // 由数据库生成
		OpenID:    uuid.New().String(),
		Username:  username,
		Password:  password,
		BasicInfo: vo.NewPatientBasicInfo(name, idCard),
		Contact:   vo.NewContactInfo(phone, "", ""),
		Status:    PatientStatusNormal,
		MedicalHistory: []entity.MedicalHistory{},
		Allergies:      []entity.Allergy{},
		CreatedAt: now,
		UpdatedAt: now,
		Events:    []event.DomainEvent{},
	}
}

// UpdateProfile 更新个人信息
func (p *Patient) UpdateProfile(info vo.PatientBasicInfo, contact vo.ContactInfo) error {
	if p.Status == PatientStatusLocked {
		return ErrPatientLocked
	}

	if err := info.Validate(); err != nil {
		return ErrInvalidInfo
	}

	if err := contact.Validate(); err != nil {
		return err
	}

	p.BasicInfo = info
	p.Contact = contact
	p.UpdatedAt = time.Now()

	p.Events = append(p.Events, event.NewPatientProfileUpdatedEvent(p.ID))

	return nil
}

// AddAllergy 添加过敏史
func (p *Patient) AddAllergy(allergen string, severity entity.AllergySeverity, reaction string) error {
	if p.Status == PatientStatusLocked {
		return ErrPatientLocked
	}

	// 检查是否已存在相同过敏源
	for _, a := range p.Allergies {
		if a.Allergen == allergen {
			return ErrDuplicateAllergy
		}
	}

	allergy := entity.NewAllergy(allergen, severity, reaction)
	p.Allergies = append(p.Allergies, allergy)
	p.UpdatedAt = time.Now()

	p.Events = append(p.Events, event.NewPatientAllergyAddedEvent(p.ID, int64(allergy.ID)))

	return nil
}

// AddMedicalHistory 添加病史
func (p *Patient) AddMedicalHistory(diseaseName string, diagnoseDate time.Time, treatment string) error {
	if p.Status == PatientStatusLocked {
		return ErrPatientLocked
	}

	// 检查是否已存在相同疾病（最近一年内）
	for _, h := range p.MedicalHistory {
		if h.DiseaseName == diseaseName && h.DiagnoseDate.After(time.Now().AddDate(-1, 0, 0)) {
			return ErrDuplicateHistory
		}
	}

	history := entity.NewMedicalHistory(diseaseName, diagnoseDate, treatment)
	p.MedicalHistory = append(p.MedicalHistory, history)
	p.UpdatedAt = time.Now()

	p.Events = append(p.Events, event.NewPatientMedicalHistoryAddedEvent(p.ID, int64(history.ID)))

	return nil
}

// SoftDelete 软删除
func (p *Patient) SoftDelete() {
	now := time.Now()
	p.DeletedAt = &now
	p.Status = PatientStatusInactive
	p.Events = append(p.Events, event.NewPatientDeletedEvent(p.ID))
}

// IsDeleted 是否已删除
func (p *Patient) IsDeleted() bool {
	return p.DeletedAt != nil
}

// IsNormal 是否正常状态
func (p *Patient) IsNormal() bool {
	return p.Status == PatientStatusNormal && !p.IsDeleted()
}

// Lock 锁定患者
func (p *Patient) Lock() {
	p.Status = PatientStatusLocked
	p.UpdatedAt = time.Now()
}

// Unlock 解锁患者
func (p *Patient) Unlock() {
	p.Status = PatientStatusNormal
	p.UpdatedAt = time.Now()
}

// Activate 激活患者
func (p *Patient) Activate() {
	p.Status = PatientStatusNormal
	p.UpdatedAt = time.Now()
}

// ChangePassword 修改密码
func (p *Patient) ChangePassword(newPassword string) error {
	if p.Status == PatientStatusLocked {
		return ErrPatientLocked
	}

	if newPassword == "" {
		return errors.New("password cannot be empty")
	}

	p.Password = newPassword
	p.UpdatedAt = time.Now()

	return nil
}

// ClearEvents 清除领域事件
func (p *Patient) ClearEvents() {
	p.Events = []event.DomainEvent{}
}

// GetMaskedInfo 获取脱敏后的患者信息
func (p Patient) GetMaskedInfo() map[string]string {
	return map[string]string{
		"phone":      p.Contact.GetMaskedPhone(),
		"id_card":    p.BasicInfo.GetMaskedIDCard(),
		"name":       p.BasicInfo.GetDisplayName(),
	}
}

// HasAllergies 是否有过敏史
func (p Patient) HasAllergies() bool {
	return len(p.Allergies) > 0
}

// HasMedicalHistory 是否有病史
func (p Patient) HasMedicalHistory() bool {
	return len(p.MedicalHistory) > 0
}

// GetSevereAllergies 获取重度过敏史
func (p Patient) GetSevereAllergies() []entity.Allergy {
	var severe []entity.Allergy
	for _, a := range p.Allergies {
		if a.Severity == entity.AllergySeveritySevere {
			severe = append(severe, a)
		}
	}
	return severe
}

// GetActiveMedicalHistory 获取进行中的病史
func (p Patient) GetActiveMedicalHistory() []entity.MedicalHistory {
	var active []entity.MedicalHistory
	for _, h := range p.MedicalHistory {
		if h.Status == entity.HistoryStatusActive || h.Status == entity.HistoryStatusChronic {
			active = append(active, h)
		}
	}
	return active
}