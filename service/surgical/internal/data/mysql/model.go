package mysql

import "time"

type SurgicalPO struct {
	ID            int64     `gorm:"primaryKey;autoIncrement"`
	MedicalID     int64     `gorm:"column:medical_id;index"`
	PatientID     int64     `gorm:"column:patient_id;index"`
	DoctorIDs     string    `gorm:"column:doctor_ids;size:500"`
	DepartmentID  int64     `gorm:"column:department_id;index"`
	Type          string    `gorm:"column:type;size:20;not null"`
	Status        string    `gorm:"column:status;size:20;not null"`
	ScheduledTime time.Time `gorm:"column:scheduled_time"`
	Duration      int64     `gorm:"column:duration"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (SurgicalPO) TableName() string {
	return "surgicals"
}

type SurgicalFlowPO struct {
	ID         int64     `gorm:"primaryKey;autoIncrement"`
	SurgicalID int64     `gorm:"column:surgical_id;index"`
	Step       int       `gorm:"column:step"`
	Name       string    `gorm:"column:name;size:100"`
	Status     string    `gorm:"column:status;size:20"`
	StartTime  time.Time `gorm:"column:start_time"`
	EndTime    time.Time `gorm:"column:end_time"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (SurgicalFlowPO) TableName() string {
	return "surgical_flows"
}