package usecase

import (
	"context"
	"errors"

	"intelligent_guidance_system_v2/service/medical/internal/domain/aggregate"
	"intelligent_guidance_system_v2/service/medical/internal/domain/entity"
	"intelligent_guidance_system_v2/service/medical/internal/domain/repository"
	"intelligent_guidance_system_v2/service/medical/internal/biz/dto"
)

var (
	ErrMedicalRecordNotFound    = errors.New("medical record not found")
	ErrInvalidInput             = errors.New("invalid input parameters")
	ErrPatientNotFound          = errors.New("patient not found")
	ErrDoctorNotFound           = errors.New("doctor not found")
	ErrDepartmentNotFound       = errors.New("department not found")
	ErrPaymentFailed            = errors.New("payment failed")
)

// MedicalUseCase defines the application use cases for medical records
type MedicalUseCase struct {
	repo         repository.MedicalRecordRepository
	patientACL   PatientACL
	doctorACL    DoctorACL
	departmentACL DepartmentACL
	paymentACL   PaymentACL
	eventRepo    repository.EventRepository
}

// PatientACL defines the anti-corruption layer for patient service
type PatientACL interface {
	GetPatient(ctx context.Context, patientID int64) (*PatientInfo, error)
	ValidatePatient(ctx context.Context, patientID int64) error
}

// PatientInfo represents patient information from ACL
type PatientInfo struct {
	ID       int64
	Name     string
	Phone    string
	Gender   int
	Age      int
}

// DoctorACL defines the anti-corruption layer for doctor service
type DoctorACL interface {
	GetDoctor(ctx context.Context, doctorID int64) (*DoctorInfo, error)
	ValidateDoctor(ctx context.Context, doctorID int64) error
}

// DoctorInfo represents doctor information from ACL
type DoctorInfo struct {
	ID           int64
	Name         string
	Title        string
	DepartmentID int64
}

// DepartmentACL defines the anti-corruption layer for department service
type DepartmentACL interface {
	GetDepartment(ctx context.Context, departmentID int64) (*DepartmentInfo, error)
	ValidateDepartment(ctx context.Context, departmentID int64) error
}

// DepartmentInfo represents department information from ACL
type DepartmentInfo struct {
	ID   int64
	Name string
	Type string
}

// PaymentACL defines the anti-corruption layer for payment service
type PaymentACL interface {
	ProcessPayment(ctx context.Context, recordID int64, amount float64, method string) error
	GetPaymentStatus(ctx context.Context, recordID int64) (*PaymentStatusInfo, error)
}

// PaymentStatusInfo represents payment status from ACL
type PaymentStatusInfo struct {
	RecordID int64
	Amount   float64
	Method   string
	Status   string
	PayTime  string
}

// NewMedicalUseCase creates a new MedicalUseCase
func NewMedicalUseCase(
	repo repository.MedicalRecordRepository,
	patientACL PatientACL,
	doctorACL DoctorACL,
	departmentACL DepartmentACL,
	paymentACL PaymentACL,
	eventRepo repository.EventRepository,
) *MedicalUseCase {
	return &MedicalUseCase{
		repo:         repo,
		patientACL:   patientACL,
		doctorACL:    doctorACL,
		departmentACL: departmentACL,
		paymentACL:   paymentACL,
		eventRepo:    eventRepo,
	}
}

// CreateMedicalRecord creates a new medical record
func (uc *MedicalUseCase) CreateMedicalRecord(ctx context.Context, req *dto.CreateMedicalRecordRequest) (*dto.MedicalRecordDTO, error) {
	// Validate patient exists
	if err := uc.patientACL.ValidatePatient(ctx, req.PatientID); err != nil {
		return nil, ErrPatientNotFound
	}

	// Validate doctor exists
	if err := uc.doctorACL.ValidateDoctor(ctx, req.DoctorID); err != nil {
		return nil, ErrDoctorNotFound
	}

	// Validate department exists
	if err := uc.departmentACL.ValidateDepartment(ctx, req.DepartmentID); err != nil {
		return nil, ErrDepartmentNotFound
	}

	// Create aggregate
	record, err := aggregate.NewMedicalRecord(req.PatientID, req.DoctorID, req.DepartmentID, req.Symptom)
	if err != nil {
		return nil, err
	}

	// Save to repository
	if err := uc.repo.Save(ctx, record); err != nil {
		return nil, err
	}

	// Publish events
	if err := uc.publishEvents(ctx, record); err != nil {
		return nil, err
	}

	return uc.toDTO(ctx, record), nil
}

// GetMedicalRecord retrieves a medical record by ID
func (uc *MedicalUseCase) GetMedicalRecord(ctx context.Context, req *dto.GetMedicalRecordRequest) (*dto.MedicalRecordDTO, error) {
	record, err := uc.repo.FindByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, ErrMedicalRecordNotFound
	}

	return uc.toDTO(ctx, record), nil
}

// UpdateMedicalRecord updates an existing medical record
func (uc *MedicalUseCase) UpdateMedicalRecord(ctx context.Context, req *dto.UpdateMedicalRecordRequest) (*dto.MedicalRecordDTO, error) {
	record, err := uc.repo.FindByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, ErrMedicalRecordNotFound
	}

	// Update fields
	if req.Symptom != "" {
		record.UpdateSymptom(req.Symptom)
	}
	if req.Diagnosis != "" {
		record.UpdateDiagnosis(req.Diagnosis)
	}
	if req.TreatmentPlan != "" {
		record.UpdateTreatmentPlan(req.TreatmentPlan)
	}

	// Save
	if err := uc.repo.Update(ctx, record); err != nil {
		return nil, err
	}

	return uc.toDTO(ctx, record), nil
}

// AddDisease adds a disease to a medical record
func (uc *MedicalUseCase) AddDisease(ctx context.Context, req *dto.AddDiseaseRequest, operatorID int64) (*dto.MedicalRecordDTO, error) {
	record, err := uc.repo.FindByID(ctx, req.MedicalRecordID)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, ErrMedicalRecordNotFound
	}

	// Create disease item
	disease, err := entity.NewDiseaseItem(req.DiseaseID, req.DiseaseName, req.Symptoms, req.Diagnosis)
	if err != nil {
		return nil, err
	}

	// Add to aggregate
	if err := record.AddDisease(disease, operatorID); err != nil {
		return nil, err
	}

	// Save
	if err := uc.repo.Update(ctx, record); err != nil {
		return nil, err
	}

	// Publish events
	if err := uc.publishEvents(ctx, record); err != nil {
		return nil, err
	}

	return uc.toDTO(ctx, record), nil
}

// AddPrescription adds a prescription to a medical record
func (uc *MedicalUseCase) AddPrescription(ctx context.Context, req *dto.AddPrescriptionRequest, operatorID int64) (*dto.MedicalRecordDTO, error) {
	record, err := uc.repo.FindByID(ctx, req.MedicalRecordID)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, ErrMedicalRecordNotFound
	}

	// Convert price to fen
	priceInFen := int64(req.Price * 100)

	// Create prescription item
	prescription, err := entity.NewPrescriptionItem(req.DrugID, req.DrugName, req.Quantity, priceInFen, req.Usage)
	if err != nil {
		return nil, err
	}

	// Add to aggregate
	if err := record.AddPrescription(prescription, operatorID); err != nil {
		return nil, err
	}

	// Save
	if err := uc.repo.Update(ctx, record); err != nil {
		return nil, err
	}

	// Publish events
	if err := uc.publishEvents(ctx, record); err != nil {
		return nil, err
	}

	return uc.toDTO(ctx, record), nil
}

// UpdateStatus updates the status of a medical record
func (uc *MedicalUseCase) UpdateStatus(ctx context.Context, req *dto.UpdateStatusRequest) (*dto.MedicalRecordDTO, error) {
	record, err := uc.repo.FindByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, ErrMedicalRecordNotFound
	}

	// Validate status
	newStatus := entity.MedicalStatus(req.Status)
	if !newStatus.IsValid() {
		return nil, ErrInvalidInput
	}

	// Update status
	if err := record.UpdateStatus(newStatus, req.OperatorID, req.Reason); err != nil {
		return nil, err
	}

	// Save
	if err := uc.repo.Update(ctx, record); err != nil {
		return nil, err
	}

	// Publish events
	if err := uc.publishEvents(ctx, record); err != nil {
		return nil, err
	}

	return uc.toDTO(ctx, record), nil
}

// CompleteMedical marks a medical record as completed
func (uc *MedicalUseCase) CompleteMedical(ctx context.Context, req *dto.CompleteMedicalRequest) (*dto.MedicalRecordDTO, error) {
	record, err := uc.repo.FindByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, ErrMedicalRecordNotFound
	}

	// Complete the record
	if err := record.Complete(req.OperatorID); err != nil {
		return nil, err
	}

	// Save
	if err := uc.repo.Update(ctx, record); err != nil {
		return nil, err
	}

	// Process payment if needed
	if record.Payment() != nil && !record.Payment().IsPaid() && record.CalculateTotal() > 0 {
		if err := uc.paymentACL.ProcessPayment(ctx, record.ID(), record.CalculateTotal(), string(record.Payment().Method())); err != nil {
			return nil, ErrPaymentFailed
		}
	}

	// Publish events
	if err := uc.publishEvents(ctx, record); err != nil {
		return nil, err
	}

	return uc.toDTO(ctx, record), nil
}

// CancelMedical cancels a medical record
func (uc *MedicalUseCase) CancelMedical(ctx context.Context, req *dto.CancelMedicalRequest) (*dto.MedicalRecordDTO, error) {
	record, err := uc.repo.FindByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, ErrMedicalRecordNotFound
	}

	// Cancel the record
	if err := record.Cancel(req.OperatorID, req.Reason); err != nil {
		return nil, err
	}

	// Save
	if err := uc.repo.Update(ctx, record); err != nil {
		return nil, err
	}

	// Publish events
	if err := uc.publishEvents(ctx, record); err != nil {
		return nil, err
	}

	return uc.toDTO(ctx, record), nil
}

// ListMedicalRecords lists medical records with pagination
func (uc *MedicalUseCase) ListMedicalRecords(ctx context.Context, req *dto.ListMedicalRecordsRequest) (*dto.ListMedicalRecordsResponse, error) {
	filter := &repository.MedicalRecordFilter{
		PatientID:     req.PatientID,
		DoctorID:      req.DoctorID,
		DepartmentID:  req.DepartmentID,
		Status:        req.Status,
		MedicalNumber: req.MedicalNumber,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		Page:          req.Page,
PageSize:     req.PageSize,
	}

	records, total, err := uc.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	dtos := make([]dto.MedicalRecordDTO, len(records))
	for i, record := range records {
		dtos[i] = *uc.toDTO(ctx, record)
	}

	return &dto.ListMedicalRecordsResponse{
		MedicalRecords: dtos,
		Total:          total,
		Page:           req.Page,
PageSize:      req.PageSize,
	}, nil
}

// publishEvents publishes domain events
func (uc *MedicalUseCase) publishEvents(ctx context.Context, record *aggregate.MedicalRecord) error {
	if len(record.Events()) == 0 {
		return nil
	}

	events := make([]interface{}, len(record.Events()))
	for i, e := range record.Events() {
		events[i] = e
	}

	if err := uc.eventRepo.Save(ctx, events...); err != nil {
		return err
	}

	record.ClearEvents()
	return nil
}

// toDTO converts aggregate to DTO
func (uc *MedicalUseCase) toDTO(ctx context.Context, record *aggregate.MedicalRecord) *dto.MedicalRecordDTO {
	result := &dto.MedicalRecordDTO{
		ID:            record.ID(),
		MedicalNumber: record.MedicalNumber(),
		PatientID:     record.PatientID(),
		DoctorID:      record.DoctorID(),
		DepartmentID:  record.DepartmentID(),
		Symptom:       record.Symptom(),
		Diagnosis:     record.Diagnosis(),
		TreatmentPlan: record.TreatmentPlan(),
		Status:        int(record.Status()),
		StatusName:    record.Status().String(),
		CreatedAt:     record.CreatedAt(),
		UpdatedAt:     record.UpdatedAt(),
	}

	// Get patient info
	if uc.patientACL != nil {
		patient, err := uc.patientACL.GetPatient(ctx, record.PatientID())
		if err == nil && patient != nil {
			result.PatientName = patient.Name
		}
	}

	// Get doctor info
	if uc.doctorACL != nil {
		doctor, err := uc.doctorACL.GetDoctor(ctx, record.DoctorID())
		if err == nil && doctor != nil {
			result.DoctorName = doctor.Name
		}
	}

	// Get department info
	if uc.departmentACL != nil {
		dept, err := uc.departmentACL.GetDepartment(ctx, record.DepartmentID())
		if err == nil && dept != nil {
			result.DepartmentName = dept.Name
		}
	}

	// Convert diseases
	diseases := make([]dto.DiseaseItemDTO, len(record.Diseases()))
	for i, d := range record.Diseases() {
		diseases[i] = dto.DiseaseItemDTO{
			ID:          d.ID(),
			DiseaseID:   d.DiseaseID(),
			DiseaseName: d.DiseaseName(),
			Symptoms:    d.Symptoms(),
			Diagnosis:   d.Diagnosis(),
		}
	}
	result.Diseases = diseases

	// Convert prescriptions
	prescriptions := make([]dto.PrescriptionItemDTO, len(record.Prescriptions()))
	for i, p := range record.Prescriptions() {
		total, _ := p.TotalPrice()
		prescriptions[i] = dto.PrescriptionItemDTO{
			ID:         p.ID(),
			DrugID:     p.DrugID(),
			DrugName:   p.DrugName(),
			Quantity:   p.Quantity(),
			Price:      p.Price().Yuan(),
			Usage:      p.Usage(),
			TotalPrice: total.Yuan(),
		}
	}
	result.Prescriptions = prescriptions

	// Convert payment
	if record.Payment() != nil {
		result.Payment = dto.PaymentInfoDTO{
			Amount:  record.Payment().Amount().Yuan(),
			Method:  string(record.Payment().Method()),
			Status:  string(record.Payment().Status()),
			PayTime: record.Payment().PayTime(),
		}
	}

	// Convert timeline
	timeline := make([]dto.StatusChangeDTO, len(record.Timeline()))
	for i, t := range record.Timeline() {
		timeline[i] = dto.StatusChangeDTO{
			FromStatus:     int(t.FromStatus()),
			FromStatusName: t.FromStatus().String(),
			ToStatus:       int(t.ToStatus()),
			ToStatusName:   t.ToStatus().String(),
			ChangeTime:     t.ChangeTime(),
			OperatorID:     t.OperatorID(),
			Reason:         t.Reason(),
		}
	}
	result.Timeline = timeline

	return result
}