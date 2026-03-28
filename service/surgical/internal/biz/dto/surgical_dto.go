package dto

import "time"

type CreateSurgicalRequest struct {
	MedicalID     int64      `json:"medical_id" binding:"required"`
	PatientID     int64      `json:"patient_id" binding:"required"`
	DoctorIDs     []int64    `json:"doctor_ids" binding:"required"`
	DepartmentID  int64      `json:"department_id" binding:"required"`
	Type          string     `json:"type" binding:"required"`
	ScheduledTime time.Time  `json:"scheduled_time" binding:"required"`
	FlowSteps     []FlowStep `json:"flow_steps"`
}

type FlowStep struct {
	Step int    `json:"step" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type StartSurgicalRequest struct {
	OperatorID int64 `json:"operator_id"`
}

type CompleteStepRequest struct {
	Step       int   `json:"step" binding:"required"`
	OperatorID int64 `json:"operator_id"`
}

type CompleteSurgicalRequest struct {
	OperatorID int64 `json:"operator_id"`
}

type CancelSurgicalRequest struct {
	Reason     string `json:"reason"`
	OperatorID int64  `json:"operator_id"`
}

type SurgicalResponse struct {
	ID            int64            `json:"id"`
	MedicalID     int64            `json:"medical_id"`
	PatientID     int64            `json:"patient_id"`
	DoctorIDs     []int64          `json:"doctor_ids"`
	DepartmentID  int64            `json:"department_id"`
	Type          string           `json:"type"`
	TypeName      string           `json:"type_name"`
	Status        string           `json:"status"`
	StatusName    string           `json:"status_name"`
	ScheduledTime time.Time        `json:"scheduled_time"`
	Duration      int64            `json:"duration"`
	Flows         []FlowResponse   `json:"flows"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}

type FlowResponse struct {
	ID        int64     `json:"id"`
	Step      int       `json:"step"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type SurgicalListResponse struct {
	Total     int64              `json:"total"`
	Page      int                `json:"page"`
	PageSize  int                `json:"page_size"`
	Surgicals []SurgicalResponse `json:"surgicals"`
}