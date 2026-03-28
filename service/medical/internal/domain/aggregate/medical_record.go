package aggregate

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"intelligent_guidance_system_v2/service/medical/internal/domain/entity"
	"intelligent_guidance_system_v2/service/medical/internal/domain/event"
	"intelligent_guidance_system_v2/service/medical/internal/domain/vo"
)

var (
	ErrInvalidRecordID      = errors.New("invalid medical record ID")
	ErrEmptyMedicalNumber   = errors.New("medical number cannot be empty")
	ErrInvalidPatientID     = errors.New("invalid patient ID")
	ErrInvalidDoctorID      = errors.New("invalid doctor ID")
	ErrInvalidDepartmentID  = errors.New("invalid department ID")
	ErrRecordAlreadyComplete = errors.New("medical record already completed")
	ErrRecordAlreadyCancelled = errors.New("medical record already cancelled")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrCannotAddToComplete  = errors.New("cannot add items to completed record")
	ErrCannotAddToCancelled = errors.New("cannot add items to cancelled record")
	ErrDiseaseAlreadyExists = errors.New("disease already exists in record")
	ErrPrescriptionAlreadyExists = errors.New("prescription already exists in record")
)

// MedicalRecord aggregate root - represents a complete medical record
type MedicalRecord struct {
	id             int64
	medicalNumber  string
	patientID      int64
	doctorID       int64
	departmentID   int64
	symptom        string
	diagnosis      string
	treatmentPlan  string
	diseases       []*entity.DiseaseItem
	prescriptions  []*entity.PrescriptionItem
	payment        *vo.PaymentInfo
	status         entity.MedicalStatus
	timeline       []*entity.StatusChange
	events         []*event.MedicalEvent
	createdAt      time.Time
	updatedAt      time.Time
}

// NewMedicalRecord creates a new MedicalRecord aggregate
func NewMedicalRecord(patientID, doctorID, departmentID int64, symptom string) (*MedicalRecord, error) {
	if patientID <= 0 {
		return nil, ErrInvalidPatientID
	}
	if doctorID <= 0 {
		return nil, ErrInvalidDoctorID
	}
	if departmentID <= 0 {
		return nil, ErrInvalidDepartmentID
	}

	now := time.Now()
	medicalNumber := generateMedicalNumber(now)

	record := &MedicalRecord{
		id:             0, // Will be assigned by persistence
		medicalNumber:  medicalNumber,
		patientID:      patientID,
		doctorID:       doctorID,
		departmentID:   departmentID,
		symptom:        symptom,
		diagnosis:      "",
		treatmentPlan:  "",
		diseases:       make([]*entity.DiseaseItem, 0),
		prescriptions:  make([]*entity.PrescriptionItem, 0),
		payment:        nil,
		status:         entity.StatusRegistering,
		timeline:       make([]*entity.StatusChange, 0),
		events:         make([]*event.MedicalEvent, 0),
		createdAt:      now,
		updatedAt:      now,
	}

	// Initialize payment info
	payment, err := vo.NewPaymentInfo(vo.Zero("CNY"), vo.PaymentMethodUnpaid)
	if err == nil {
		record.payment = payment
	}

	return record, nil
}

// ReconstructMedicalRecord reconstructs a MedicalRecord from persistence
func ReconstructMedicalRecord(
	id int64,
	medicalNumber string,
	patientID, doctorID, departmentID int64,
	symptom, diagnosis, treatmentPlan string,
	diseases []*entity.DiseaseItem,
	prescriptions []*entity.PrescriptionItem,
	payment *vo.PaymentInfo,
	status entity.MedicalStatus,
	timeline []*entity.StatusChange,
	createdAt, updatedAt time.Time,
) *MedicalRecord {
	return &MedicalRecord{
		id:             id,
		medicalNumber:  medicalNumber,
		patientID:      patientID,
		doctorID:       doctorID,
		departmentID:   departmentID,
		symptom:        symptom,
		diagnosis:      diagnosis,
		treatmentPlan:  treatmentPlan,
		diseases:       diseases,
		prescriptions:  prescriptions,
		payment:        payment,
		status:         status,
		timeline:       timeline,
		events:         make([]*event.MedicalEvent, 0),
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}
}

// generateMedicalNumber generates a unique medical number
func generateMedicalNumber(t time.Time) string {
	return fmt.Sprintf("MR%s%s", t.Format("Ymd"), uuid.New().String()[:8])
}

// ID returns the medical record ID
func (m *MedicalRecord) ID() int64 {
	return m.id
}

// MedicalNumber returns the medical number
func (m *MedicalRecord) MedicalNumber() string {
	return m.medicalNumber
}

// PatientID returns the patient ID
func (m *MedicalRecord) PatientID() int64 {
	return m.patientID
}

// DoctorID returns the doctor ID
func (m *MedicalRecord) DoctorID() int64 {
	return m.doctorID
}

// DepartmentID returns the department ID
func (m *MedicalRecord) DepartmentID() int64 {
	return m.departmentID
}

// Symptom returns the symptom description
func (m *MedicalRecord) Symptom() string {
	return m.symptom
}

// Diagnosis returns the diagnosis
func (m *MedicalRecord) Diagnosis() string {
	return m.diagnosis
}

// TreatmentPlan returns the treatment plan
func (m *MedicalRecord) TreatmentPlan() string {
	return m.treatmentPlan
}

// Diseases returns the disease items
func (m *MedicalRecord) Diseases() []*entity.DiseaseItem {
	return m.diseases
}

// Prescriptions returns the prescription items
func (m *MedicalRecord) Prescriptions() []*entity.PrescriptionItem {
	return m.prescriptions
}

// Payment returns the payment info
func (m *MedicalRecord) Payment() *vo.PaymentInfo {
	return m.payment
}

// Status returns the current status
func (m *MedicalRecord) Status() entity.MedicalStatus {
	return m.status
}

// Timeline returns the status change history
func (m *MedicalRecord) Timeline() []*entity.StatusChange {
	return m.timeline
}

// Events returns the pending domain events
func (m *MedicalRecord) Events() []*event.MedicalEvent {
	return m.events
}

// CreatedAt returns the creation timestamp
func (m *MedicalRecord) CreatedAt() time.Time {
	return m.createdAt
}

// UpdatedAt returns the update timestamp
func (m *MedicalRecord) UpdatedAt() time.Time {
	return m.updatedAt
}

// SetID sets the medical record ID (called by repository after save)
func (m *MedicalRecord) SetID(id int64) {
	m.id = id
}

// AddDisease adds a disease item to the medical record
func (m *MedicalRecord) AddDisease(disease *entity.DiseaseItem, operatorID int64) error {
	if !m.CanAddDisease() {
		if m.status.IsTerminal() {
			if m.status == entity.StatusCompleted {
				return ErrCannotAddToComplete
			}
			return ErrCannotAddToCancelled
		}
		return fmt.Errorf("cannot add disease in status %s", m.status.String())
	}

	// Check if disease already exists
	for _, d := range m.diseases {
		if d.DiseaseID() == disease.DiseaseID() {
			return ErrDiseaseAlreadyExists
		}
	}

	m.diseases = append(m.diseases, disease)
	m.updatedAt = time.Now()

	// Add domain event
	m.events = append(m.events, event.DiseaseAddedEvent(
		m.id, operatorID, disease.DiseaseID(), disease.DiseaseName(),
	))

	return nil
}

// AddPrescription adds a prescription item to the medical record
func (m *MedicalRecord) AddPrescription(prescription *entity.PrescriptionItem, operatorID int64) error {
	if !m.CanAddPrescription() {
		if m.status.IsTerminal() {
			if m.status == entity.StatusCompleted {
				return ErrCannotAddToComplete
			}
			return ErrCannotAddToCancelled
		}
		return fmt.Errorf("cannot add prescription in status %s", m.status.String())
	}

	// Check if prescription already exists
	for _, p := range m.prescriptions {
		if p.DrugID() == prescription.DrugID() {
			return ErrPrescriptionAlreadyExists
		}
	}

	m.prescriptions = append(m.prescriptions, prescription)
	m.updatedAt = time.Now()

	// Update payment amount
	m.updatePaymentAmount()

	// Add domain event
	m.events = append(m.events, event.PrescriptionAddedEvent(
		m.id, operatorID, prescription.DrugID(), prescription.DrugName(),
		prescription.Quantity(), prescription.Price().Yuan(),
	))

	return nil
}

// updatePaymentAmount recalculates the total payment amount
func (m *MedicalRecord) updatePaymentAmount() {
	totalFen := int64(0)
	for _, p := range m.prescriptions {
		total, err := p.TotalPrice()
		if err == nil {
			totalFen += total.Amount()
		}
	}

	if m.payment == nil {
		m.payment, _ = vo.NewPaymentInfo(vo.Zero("CNY"), vo.PaymentMethodUnpaid)
	}

	newAmount, err := vo.NewMoney(totalFen, "CNY")
	if err == nil {
		m.payment.SetAmount(newAmount)
	}
}

// UpdateStatus updates the medical record status
func (m *MedicalRecord) UpdateStatus(newStatus entity.MedicalStatus, operatorID int64, reason string) error {
	if !m.IsValidStatusTransition(m.status, newStatus) {
		return ErrInvalidStatusTransition
	}

	oldStatus := m.status
	m.status = newStatus
	m.updatedAt = time.Now()

	// Record status change
	change := entity.NewStatusChange(oldStatus, newStatus, operatorID, reason)
	m.timeline = append(m.timeline, change)

	// Add domain event
	m.events = append(m.events, event.StatusChangedEvent(
		m.id, operatorID, oldStatus, newStatus, reason,
	))

	return nil
}

// Complete marks the medical record as completed
func (m *MedicalRecord) Complete(operatorID int64) error {
	if m.status == entity.StatusCompleted {
		return ErrRecordAlreadyComplete
	}
	if m.status == entity.StatusCancelled {
		return ErrRecordAlreadyCancelled
	}

	// Can only complete from Diagnosing or Operating status
	if m.status != entity.StatusDiagnosing && m.status != entity.StatusOperating {
		return ErrInvalidStatusTransition
	}

	return m.UpdateStatus(entity.StatusCompleted, operatorID, "Medical record completed")
}

// Cancel cancels the medical record
func (m *MedicalRecord) Cancel(operatorID int64, reason string) error {
	if m.status == entity.StatusCompleted {
		return ErrRecordAlreadyComplete
	}
	if m.status == entity.StatusCancelled {
		return ErrRecordAlreadyCancelled
	}

	return m.UpdateStatus(entity.StatusCancelled, operatorID, reason)
}

// CalculateTotal calculates the total amount for all prescriptions
func (m *MedicalRecord) CalculateTotal() float64 {
	totalFen := int64(0)
	for _, p := range m.prescriptions {
		total, err := p.TotalPrice()
		if err == nil {
			totalFen += total.Amount()
		}
	}
	return float64(totalFen) / 100.0
}

// CanAddDisease checks if a disease can be added
func (m *MedicalRecord) CanAddDisease() bool {
	return m.status == entity.StatusDiagnosing || m.status == entity.StatusQueuing
}

// CanAddPrescription checks if a prescription can be added
func (m *MedicalRecord) CanAddPrescription() bool {
	return m.status == entity.StatusDiagnosing || m.status == entity.StatusQueuing
}

// IsValidStatusTransition checks if a status transition is valid
func (m *MedicalRecord) IsValidStatusTransition(from, to entity.MedicalStatus) bool {
	return entity.IsValidTransition(from, to)
}

// UpdateDiagnosis updates the diagnosis
func (m *MedicalRecord) UpdateDiagnosis(diagnosis string) error {
	if m.status.IsTerminal() {
		return fmt.Errorf("cannot update diagnosis in terminal status")
	}
	m.diagnosis = diagnosis
	m.updatedAt = time.Now()
	return nil
}

// UpdateTreatmentPlan updates the treatment plan
func (m *MedicalRecord) UpdateTreatmentPlan(treatmentPlan string) error {
	if m.status.IsTerminal() {
		return fmt.Errorf("cannot update treatment plan in terminal status")
	}
	m.treatmentPlan = treatmentPlan
	m.updatedAt = time.Now()
	return nil
}

// UpdateSymptom updates the symptom description
func (m *MedicalRecord) UpdateSymptom(symptom string) error {
	if m.status.IsTerminal() {
		return fmt.Errorf("cannot update symptom in terminal status")
	}
	m.symptom = symptom
	m.updatedAt = time.Now()
	return nil
}

// ClearEvents clears pending domain events
func (m *MedicalRecord) ClearEvents() {
	m.events = make([]*event.MedicalEvent, 0)
}

// HasDiseases checks if the record has any diseases
func (m *MedicalRecord) HasDiseases() bool {
	return len(m.diseases) > 0
}

// HasPrescriptions checks if the record has any prescriptions
func (m *MedicalRecord) HasPrescriptions() bool {
	return len(m.prescriptions) > 0
}

// DiseaseCount returns the number of diseases
func (m *MedicalRecord) DiseaseCount() int {
	return len(m.diseases)
}

// PrescriptionCount returns the number of prescriptions
func (m *MedicalRecord) PrescriptionCount() int {
	return len(m.prescriptions)
}

// IsTerminal checks if the record is in terminal status
func (m *MedicalRecord) IsTerminal() bool {
	return m.status.IsTerminal()
}