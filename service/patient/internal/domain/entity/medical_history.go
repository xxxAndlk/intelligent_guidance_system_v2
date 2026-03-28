// service/patient/internal/domain/entity/medical_history.go
package entity

import "time"

// HistoryID 病史ID
type HistoryID int64

// MedicalHistory 病史实体
type MedicalHistory struct {
	ID           HistoryID
	DiseaseName  string    // 疾病名称
	DiagnoseDate time.Time // 诊断日期
	Treatment    string    // 治疗方案
	Status       HistoryStatus // 病史状态
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// HistoryStatus 病史状态
type HistoryStatus int32

const (
	HistoryStatusActive    HistoryStatus = 1 // 进行中
	HistoryStatusRecovered HistoryStatus = 2 // 已康复
	HistoryStatusChronic   HistoryStatus = 3 // 慢性病
)

func (s HistoryStatus) String() string {
	switch s {
	case HistoryStatusActive:
		return "进行中"
	case HistoryStatusRecovered:
		return "已康复"
	case HistoryStatusChronic:
		return "慢性病"
	default:
		return "未知"
	}
}

// NewMedicalHistory 创建病史
func NewMedicalHistory(diseaseName string, diagnoseDate time.Time, treatment string) MedicalHistory {
	now := time.Now()
	return MedicalHistory{
		ID:           HistoryID(now.UnixNano()),
		DiseaseName:  diseaseName,
		DiagnoseDate: diagnoseDate,
		Treatment:    treatment,
		Status:       HistoryStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// UpdateTreatment 更新治疗方案
func (m *MedicalHistory) UpdateTreatment(treatment string) {
	m.Treatment = treatment
	m.UpdatedAt = time.Now()
}

// MarkRecovered 标记为已康复
func (m *MedicalHistory) MarkRecovered() {
	m.Status = HistoryStatusRecovered
	m.UpdatedAt = time.Now()
}

// MarkChronic 标记为慢性病
func (m *MedicalHistory) MarkChronic() {
	m.Status = HistoryStatusChronic
	m.UpdatedAt = time.Now()
}

// IsRecovered 是否已康复
func (m MedicalHistory) IsRecovered() bool {
	return m.Status == HistoryStatusRecovered
}

// IsChronic 是否为慢性病
func (m MedicalHistory) IsChronic() bool {
	return m.Status == HistoryStatusChronic
}