package dto

import (
	"time"
)

// MedicalRecordDTO represents a medical record for API responses
type MedicalRecordDTO struct {
	ID            int64                `json:"id"`
	MedicalNumber string               `json:"medical_number"`
	PatientID     int64                `json:"patient_id"`
	PatientName   string               `json:"patient_name,omitempty"`
	DoctorID      int64                `json:"doctor_id"`
	DoctorName    string               `json:"doctor_name,omitempty"`
	DepartmentID  int64                `json:"department_id"`
	DepartmentName string              `json:"department_name,omitempty"`
	Symptom       string               `json:"symptom"`
	Diagnosis     string               `json:"diagnosis"`
	TreatmentPlan string               `json:"treatment_plan"`
	Diseases      []DiseaseItemDTO     `json:"diseases"`
	Prescriptions []PrescriptionItemDTO `json:"prescriptions"`
	Payment       PaymentInfoDTO       `json:"payment"`
	Status        int                  `json:"status"`
	StatusName    string               `json:"status_name"`
	Timeline      []StatusChangeDTO    `json:"timeline"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
}

// DiseaseItemDTO represents a disease diagnosis item
type DiseaseItemDTO struct {
	ID          string `json:"id"`
	DiseaseID   int64  `json:"disease_id"`
	DiseaseName string `json:"disease_name"`
	Symptoms    string `json:"symptoms"`
	Diagnosis   string `json:"diagnosis"`
}

// PrescriptionItemDTO represents a prescription drug item
type PrescriptionItemDTO struct {
	ID        string `json:"id"`
	DrugID    int64  `json:"drug_id"`
	DrugName  string `json:"drug_name"`
	Quantity  int32  `json:"quantity"`
	Price     float64 `json:"price"`
	Usage     string `json:"usage"`
	TotalPrice float64 `json:"total_price"`
}

// PaymentInfoDTO represents payment information
type PaymentInfoDTO struct {
	Amount   float64 `json:"amount"`
	Method   string  `json:"method"`
	Status   string  `json:"status"`
	PayTime  *time.Time `json:"pay_time,omitempty"`
}

// StatusChangeDTO represents a status change record
type StatusChangeDTO struct {
	FromStatus int    `json:"from_status"`
	FromStatusName string `json:"from_status_name"`
	ToStatus   int    `json:"to_status"`
	ToStatusName string `json:"to_status_name"`
	ChangeTime time.Time `json:"change_time"`
	OperatorID int64  `json:"operator_id"`
	Reason     string `json:"reason"`
}

// CreateMedicalRecordRequest represents the request for creating a medical record
type CreateMedicalRecordRequest struct {
	MedicalNumber string `json:"medical_number,omitempty"`
	PatientID     int64  `json:"patient_id" validate:"required,min=1"`
	DoctorID      int64  `json:"doctor_id" validate:"required,min=1"`
	DepartmentID  int64  `json:"department_id" validate:"required,min=1"`
	Symptom       string `json:"symptom" validate:"required"`
}

// UpdateMedicalRecordRequest represents the request for updating a medical record
type UpdateMedicalRecordRequest struct {
	ID            int64  `json:"id" validate:"required,min=1"`
	Symptom       string `json:"symptom"`
	Diagnosis     string `json:"diagnosis"`
	TreatmentPlan string `json:"treatment_plan"`
}

// AddDiseaseRequest represents the request for adding a disease
type AddDiseaseRequest struct {
	MedicalRecordID int64  `json:"medical_record_id" validate:"required,min=1"`
	DiseaseID       int64  `json:"disease_id" validate:"required,min=1"`
	DiseaseName     string `json:"disease_name" validate:"required"`
	Symptoms        string `json:"symptoms" validate:"required"`
	Diagnosis       string `json:"diagnosis"`
}

// AddPrescriptionRequest represents the request for adding a prescription
type AddPrescriptionRequest struct {
	MedicalRecordID int64   `json:"medical_record_id" validate:"required,min=1"`
	DrugID          int64   `json:"drug_id" validate:"required,min=1"`
	DrugName        string  `json:"drug_name" validate:"required"`
	Quantity        int32   `json:"quantity" validate:"required,min=1"`
	Price           float64 `json:"price" validate:"min=0"`
	Usage           string  `json:"usage" validate:"required"`
}

// UpdateStatusRequest represents the request for updating status
type UpdateStatusRequest struct {
	ID         int64  `json:"id" validate:"required,min=1"`
	Status     int    `json:"status" validate:"required,min=1,max=6"`
	OperatorID int64  `json:"operator_id" validate:"required,min=1"`
	Reason     string `json:"reason"`
}

// CompleteMedicalRequest represents the request for completing a medical record
type CompleteMedicalRequest struct {
	ID         int64 `json:"id" validate:"required,min=1"`
	OperatorID int64 `json:"operator_id" validate:"required,min=1"`
}

// CancelMedicalRequest represents the request for cancelling a medical record
type CancelMedicalRequest struct {
	ID         int64  `json:"id" validate:"required,min=1"`
	OperatorID int64  `json:"operator_id" validate:"required,min=1"`
	Reason     string `json:"reason" validate:"required"`
}

// ListMedicalRecordsRequest represents the request for listing medical records
type ListMedicalRecordsRequest struct {
	PatientID     int64  `json:"patient_id,omitempty"`
	DoctorID      int64  `json:"doctor_id,omitempty"`
	DepartmentID  int64  `json:"department_id,omitempty"`
	Status        int    `json:"status,omitempty"`
	MedicalNumber string `json:"medical_number,omitempty"`
	StartDate     string `json:"start_date,omitempty"`
	EndDate       string `json:"end_date,omitempty"`
Page          int    `json:"page" validate:"min=1"`
	PageSize      int    `json:"page_size" validate:"min=1,max=100"`
}

// ListMedicalRecordsResponse represents the response for listing medical records
type ListMedicalRecordsResponse struct {
	MedicalRecords []MedicalRecordDTO `json:"medical_records"`
	Total          int64              `json:"total"`
	Page           int                `json:"page"`
	PageSize       int                `json:"page_size"`
}

// GetMedicalRecordRequest represents the request for getting a medical record
type GetMedicalRecordRequest struct {
	ID int64 `json:"id" validate:"required,min=1"`
}