package repository

import (
	"context"

	"intelligent_guidance_system_v2/service/medical/internal/domain/aggregate"
)

// MedicalRecordRepository defines the repository interface for MedicalRecord aggregate
type MedicalRecordRepository interface {
	// Save persists a MedicalRecord aggregate
	Save(ctx context.Context, record *aggregate.MedicalRecord) error

	// FindByID retrieves a MedicalRecord by its ID
	FindByID(ctx context.Context, id int64) (*aggregate.MedicalRecord, error)

	// FindByMedicalNumber retrieves a MedicalRecord by its medical number
	FindByMedicalNumber(ctx context.Context, medicalNumber string) (*aggregate.MedicalRecord, error)

	// FindByPatientID retrieves MedicalRecords for a patient
	FindByPatientID(ctx context.Context, patientID int64, limit, offset int) ([]*aggregate.MedicalRecord, int64, error)

	// FindByDoctorID retrieves MedicalRecords for a doctor
	FindByDoctorID(ctx context.Context, doctorID int64, limit, offset int) ([]*aggregate.MedicalRecord, int64, error)

	// FindByDepartmentID retrieves MedicalRecords for a department
	FindByDepartmentID(ctx context.Context, departmentID int64, limit, offset int) ([]*aggregate.MedicalRecord, int64, error)

	// FindByStatus retrieves MedicalRecords by status
	FindByStatus(ctx context.Context, status int, limit, offset int) ([]*aggregate.MedicalRecord, int64, error)

	// List retrieves MedicalRecords with pagination and filters
	List(ctx context.Context, filter *MedicalRecordFilter) ([]*aggregate.MedicalRecord, int64, error)

	// Delete removes a MedicalRecord (soft delete)
	Delete(ctx context.Context, id int64) error

	// Update updates an existing MedicalRecord
	Update(ctx context.Context, record *aggregate.MedicalRecord) error
}

// MedicalRecordFilter defines filter criteria for listing medical records
type MedicalRecordFilter struct {
	PatientID     int64
	DoctorID      int64
	DepartmentID  int64
	Status        int
	MedicalNumber string
	StartDate     string
	EndDate       string
Page          int
	PageSize      int
}

// EventRepository defines the repository interface for domain events
type EventRepository interface {
	// Save persists domain events
	Save(ctx context.Context, events ...interface{}) error

	// FindByRecordID retrieves events for a medical record
	FindByRecordID(ctx context.Context, recordID int64) ([]interface{}, error)
}