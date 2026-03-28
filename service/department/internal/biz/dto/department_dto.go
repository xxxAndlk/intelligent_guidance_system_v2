package dto

import "time"

type CreateDepartmentRequest struct {
	Name                string `json:"name" binding:"required"`
	Introduction        string `json:"introduction"`
	PersonInChargeID    int64  `json:"person_in_charge_id"`
	Phone               string `json:"phone"`
	Address             string `json:"address"`
	Type                string `json:"type" binding:"required"`
	RegistrationFee     int64  `json:"registration_fee"`
	ExpertRegistrationFee int64 `json:"expert_registration_fee"`
}

type UpdateDepartmentRequest struct {
	Introduction        string `json:"introduction"`
	PersonInChargeID    int64  `json:"person_in_charge_id"`
	Phone               string `json:"phone"`
	Address             string `json:"address"`
	RegistrationFee     int64  `json:"registration_fee"`
	ExpertRegistrationFee int64 `json:"expert_registration_fee"`
}

type AssignDoctorRequest struct {
	DoctorID   int64 `json:"doctor_id" binding:"required"`
	OperatorID int64 `json:"operator_id"`
}

type RemoveDoctorRequest struct {
	DoctorID   int64 `json:"doctor_id" binding:"required"`
	OperatorID int64 `json:"operator_id"`
}

type SetDutyDoctorRequest struct {
	DoctorID   int64 `json:"doctor_id" binding:"required"`
	OperatorID int64 `json:"operator_id"`
}

type ChangeStatusRequest struct {
	Status     string `json:"status" binding:"required"`
	OperatorID int64  `json:"operator_id"`
}

type DepartmentResponse struct {
	ID                  int64     `json:"id"`
	Name                string    `json:"name"`
	Introduction        string    `json:"introduction"`
	PersonInChargeID    int64     `json:"person_in_charge_id"`
	Phone               string    `json:"phone"`
	Address             string    `json:"address"`
	Type                string    `json:"type"`
	TypeName            string    `json:"type_name"`
	StaffCount          int       `json:"staff_count"`
	Status              string    `json:"status"`
	StatusName          string    `json:"status_name"`
	RegistrationFee     float64   `json:"registration_fee"`
	ExpertRegistrationFee float64 `json:"expert_registration_fee"`
	DutyDoctorID        int64     `json:"duty_doctor_id"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type DepartmentListResponse struct {
	Total      int64                `json:"total"`
	Page       int                  `json:"page"`
	PageSize   int                  `json:"page_size"`
	Departments []DepartmentResponse `json:"departments"`
}

type DoctorInDepartmentResponse struct {
	DoctorID   int64     `json:"doctor_id"`
	IsDuty     bool      `json:"is_duty"`
	AssignedAt time.Time `json:"assigned_at"`
}