package mysql

import "time"

type RegistrationPO struct {
	ID              int64     `gorm:"primaryKey;autoIncrement"`
	PatientID       int64     `gorm:"column:patient_id;index"`
	DoctorID        int64     `gorm:"column:doctor_id;index"`
	DepartmentID    int64     `gorm:"column:department_id;index"`
	Type            string    `gorm:"column:type;size:20;not null"`
	Status          string    `gorm:"column:status;size:20;not null"`
	AppointmentTime time.Time `gorm:"column:appointment_time"`
	Fee             int64     `gorm:"column:fee"`
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
}

func (RegistrationPO) TableName() string {
	return "registrations"
}