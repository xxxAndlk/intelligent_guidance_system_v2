package dto

import "time"

type CreateRegistrationRequest struct {
	PatientID       int64     `json:"patient_id" binding:"required"`
	DoctorID        int64     `json:"doctor_id" binding:"required"`
	DepartmentID    int64     `json:"department_id" binding:"required"`
	Type            string    `json:"type" binding:"required"`
	AppointmentTime time.Time `json:"appointment_time" binding:"required"`
}

type CancelRegistrationRequest struct {
	Reason     string `json:"reason"`
	OperatorID int64  `json:"operator_id"`
}

type ConfirmRegistrationRequest struct {
	OperatorID int64 `json:"operator_id"`
}

type CompleteRegistrationRequest struct {
	OperatorID int64 `json:"operator_id"`
}

type RegistrationResponse struct {
	ID              int64     `json:"id"`
	PatientID       int64     `json:"patient_id"`
	DoctorID        int64     `json:"doctor_id"`
	DepartmentID    int64     `json:"department_id"`
	Type            string    `json:"type"`
	TypeName        string    `json:"type_name"`
	Status          string    `json:"status"`
	StatusName      string    `json:"status_name"`
	AppointmentTime time.Time `json:"appointment_time"`
	Fee             float64   `json:"fee"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type RegistrationListResponse struct {
	Total        int64                 `json:"total"`
	Page         int                   `json:"page"`
	PageSize     int                   `json:"page_size"`
	Registrations []RegistrationResponse `json:"registrations"`
}