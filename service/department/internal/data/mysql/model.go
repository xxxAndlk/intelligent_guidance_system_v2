package mysql

import "time"

type DepartmentPO struct {
	ID                  int64     `gorm:"primaryKey;autoIncrement"`
	Name                string    `gorm:"column:name;size:100;not null"`
	Introduction        string    `gorm:"column:introduction;size:500"`
	PersonInChargeID    int64     `gorm:"column:person_in_charge_id"`
	Phone               string    `gorm:"column:phone;size:20"`
	Address             string    `gorm:"column:address;size:200"`
	Type                string    `gorm:"column:type;size:20;not null"`
	StaffCount          int       `gorm:"column:staff_count"`
	Status              string    `gorm:"column:status;size:20;not null"`
	RegistrationFee     int64     `gorm:"column:registration_fee"`
	ExpertRegistrationFee int64   `gorm:"column:expert_registration_fee"`
	DutyDoctorID        int64     `gorm:"column:duty_doctor_id"`
	CreatedAt           time.Time `gorm:"column:created_at"`
	UpdatedAt           time.Time `gorm:"column:updated_at"`
}

func (DepartmentPO) TableName() string {
	return "departments"
}

type DoctorAssignmentPO struct {
	ID           int64     `gorm:"primaryKey;autoIncrement"`
	DepartmentID int64     `gorm:"column:department_id;index"`
	DoctorID     int64     `gorm:"column:doctor_id;index"`
	AssignedAt   time.Time `gorm:"column:assigned_at"`
	IsDuty       bool      `gorm:"column:is_duty"`
	IsActive     bool      `gorm:"column:is_active"`
}

func (DoctorAssignmentPO) TableName() string {
	return "doctor_assignments"
}