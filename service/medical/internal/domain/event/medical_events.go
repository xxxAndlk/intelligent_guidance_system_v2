package event

import (
	"time"

	"intelligent_guidance_system_v2/service/medical/internal/domain/entity"
)

// EventType represents the type of medical domain event
type EventType string

const (
	EventMedicalRecordCreated     EventType = "medical_record.created"
	EventMedicalRecordUpdated     EventType = "medical_record.updated"
	EventMedicalRecordCompleted   EventType = "medical_record.completed"
	EventMedicalRecordCancelled   EventType = "medical_record.cancelled"
	EventDiseaseAdded             EventType = "medical_record.disease_added"
	EventPrescriptionAdded        EventType = "medical_record.prescription_added"
	EventStatusChanged            EventType = "medical_record.status_changed"
	EventPaymentUpdated           EventType = "medical_record.payment_updated"
)

// MedicalEvent represents a domain event for medical records
type MedicalEvent struct {
	eventType    EventType
	recordID     int64
	occurredAt   time.Time
	operatorID   int64
	payload      map[string]interface{}
}

// NewMedicalEvent creates a new medical domain event
func NewMedicalEvent(eventType EventType, recordID, operatorID int64, payload map[string]interface{}) *MedicalEvent {
	return &MedicalEvent{
		eventType:  eventType,
		recordID:   recordID,
		occurredAt: time.Now(),
		operatorID: operatorID,
		payload:    payload,
	}
}

// EventType returns the event type
func (e *MedicalEvent) EventType() EventType {
	return e.eventType
}

// RecordID returns the medical record ID
func (e *MedicalEvent) RecordID() int64 {
	return e.recordID
}

// OccurredAt returns the event occurrence time
func (e *MedicalEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// OperatorID returns the operator ID
func (e *MedicalEvent) OperatorID() int64 {
	return e.operatorID
}

// Payload returns the event payload
func (e *MedicalEvent) Payload() map[string]interface{} {
	return e.payload
}

// MedicalRecordCreatedEvent creates a medical record created event
func MedicalRecordCreatedEvent(recordID, patientID, doctorID, departmentID int64, medicalNumber string, operatorID int64) *MedicalEvent {
	return NewMedicalEvent(
		EventMedicalRecordCreated,
		recordID,
		operatorID,
		map[string]interface{}{
			"patient_id":     patientID,
			"doctor_id":      doctorID,
			"department_id":  departmentID,
			"medical_number": medicalNumber,
		},
	)
}

// MedicalRecordCompletedEvent creates a medical record completed event
func MedicalRecordCompletedEvent(recordID, operatorID int64, totalAmount float64) *MedicalEvent {
	return NewMedicalEvent(
		EventMedicalRecordCompleted,
		recordID,
		operatorID,
		map[string]interface{}{
			"total_amount": totalAmount,
		},
	)
}

// MedicalRecordCancelledEvent creates a medical record cancelled event
func MedicalRecordCancelledEvent(recordID, operatorID int64, reason string) *MedicalEvent {
	return NewMedicalEvent(
		EventMedicalRecordCancelled,
		recordID,
		operatorID,
		map[string]interface{}{
			"reason": reason,
		},
	)
}

// DiseaseAddedEvent creates a disease added event
func DiseaseAddedEvent(recordID, operatorID int64, diseaseID int64, diseaseName string) *MedicalEvent {
	return NewMedicalEvent(
		EventDiseaseAdded,
		recordID,
		operatorID,
		map[string]interface{}{
			"disease_id":   diseaseID,
			"disease_name": diseaseName,
		},
	)
}

// PrescriptionAddedEvent creates a prescription added event
func PrescriptionAddedEvent(recordID, operatorID int64, drugID int64, drugName string, quantity int32, price float64) *MedicalEvent {
	return NewMedicalEvent(
		EventPrescriptionAdded,
		recordID,
		operatorID,
		map[string]interface{}{
			"drug_id":   drugID,
			"drug_name": drugName,
			"quantity":  quantity,
			"price":     price,
		},
	)
}

// StatusChangedEvent creates a status changed event
func StatusChangedEvent(recordID, operatorID int64, fromStatus, toStatus entity.MedicalStatus, reason string) *MedicalEvent {
	return NewMedicalEvent(
		EventStatusChanged,
		recordID,
		operatorID,
		map[string]interface{}{
			"from_status": fromStatus.String(),
			"to_status":   toStatus.String(),
			"reason":      reason,
		},
	)
}

// PaymentUpdatedEvent creates a payment updated event
func PaymentUpdatedEvent(recordID, operatorID int64, amount float64, method, status string) *MedicalEvent {
	return NewMedicalEvent(
		EventPaymentUpdated,
		recordID,
		operatorID,
		map[string]interface{}{
			"amount": amount,
			"method": method,
			"status": status,
		},
	)
}