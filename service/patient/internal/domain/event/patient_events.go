// service/patient/internal/domain/event/patient_events.go
package event

import (
	"time"

	"intelligent-guidance-system/service/patient/internal/domain/aggregate"
)

// DomainEvent 领域事件接口
type DomainEvent interface {
	GetAggregateID() aggregate.PatientID
	GetOccurredAt() time.Time
	GetEventType() string
}

// PatientRegisteredEvent 患者注册事件
type PatientRegisteredEvent struct {
	AggregateID aggregate.PatientID
	OccurredAt  time.Time
	EventType   string
}

func NewPatientRegisteredEvent(patientID aggregate.PatientID) PatientRegisteredEvent {
	return PatientRegisteredEvent{
		AggregateID: patientID,
		OccurredAt:  time.Now(),
		EventType:   "PatientRegistered",
	}
}

func (e PatientRegisteredEvent) GetAggregateID() aggregate.PatientID {
	return e.AggregateID
}

func (e PatientRegisteredEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

func (e PatientRegisteredEvent) GetEventType() string {
	return e.EventType
}

// PatientProfileUpdatedEvent 患者信息更新事件
type PatientProfileUpdatedEvent struct {
	AggregateID aggregate.PatientID
	OccurredAt  time.Time
	EventType   string
}

func NewPatientProfileUpdatedEvent(patientID aggregate.PatientID) PatientProfileUpdatedEvent {
	return PatientProfileUpdatedEvent{
		AggregateID: patientID,
		OccurredAt:  time.Now(),
		EventType:   "PatientProfileUpdated",
	}
}

func (e PatientProfileUpdatedEvent) GetAggregateID() aggregate.PatientID {
	return e.AggregateID
}

func (e PatientProfileUpdatedEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

func (e PatientProfileUpdatedEvent) GetEventType() string {
	return e.EventType
}

// PatientAllergyAddedEvent 添加过敏史事件
type PatientAllergyAddedEvent struct {
	AggregateID  aggregate.PatientID
	AllergyID    int64
	OccurredAt   time.Time
	EventType    string
}

func NewPatientAllergyAddedEvent(patientID aggregate.PatientID, allergyID int64) PatientAllergyAddedEvent {
	return PatientAllergyAddedEvent{
		AggregateID: patientID,
		AllergyID:   allergyID,
		OccurredAt:  time.Now(),
		EventType:   "PatientAllergyAdded",
	}
}

func (e PatientAllergyAddedEvent) GetAggregateID() aggregate.PatientID {
	return e.AggregateID
}

func (e PatientAllergyAddedEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

func (e PatientAllergyAddedEvent) GetEventType() string {
	return e.EventType
}

// PatientMedicalHistoryAddedEvent 添加病史事件
type PatientMedicalHistoryAddedEvent struct {
	AggregateID aggregate.PatientID
	HistoryID   int64
	OccurredAt  time.Time
	EventType   string
}

func NewPatientMedicalHistoryAddedEvent(patientID aggregate.PatientID, historyID int64) PatientMedicalHistoryAddedEvent {
	return PatientMedicalHistoryAddedEvent{
		AggregateID: patientID,
		HistoryID:   historyID,
		OccurredAt:  time.Now(),
		EventType:   "PatientMedicalHistoryAdded",
	}
}

func (e PatientMedicalHistoryAddedEvent) GetAggregateID() aggregate.PatientID {
	return e.AggregateID
}

func (e PatientMedicalHistoryAddedEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

func (e PatientMedicalHistoryAddedEvent) GetEventType() string {
	return e.EventType
}

// PatientDeletedEvent 患者删除事件
type PatientDeletedEvent struct {
	AggregateID aggregate.PatientID
	OccurredAt  time.Time
	EventType   string
}

func NewPatientDeletedEvent(patientID aggregate.PatientID) PatientDeletedEvent {
	return PatientDeletedEvent{
		AggregateID: patientID,
		OccurredAt:  time.Now(),
		EventType:   "PatientDeleted",
	}
}

func (e PatientDeletedEvent) GetAggregateID() aggregate.PatientID {
	return e.AggregateID
}

func (e PatientDeletedEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

func (e PatientDeletedEvent) GetEventType() string {
	return e.EventType
}

// PatientLoggedInEvent 患者登录事件
type PatientLoggedInEvent struct {
	AggregateID aggregate.PatientID
	OccurredAt  time.Time
	EventType   string
	LoginIP     string
}

func NewPatientLoggedInEvent(patientID aggregate.PatientID, loginIP string) PatientLoggedInEvent {
	return PatientLoggedInEvent{
		AggregateID: patientID,
		OccurredAt:  time.Now(),
		EventType:   "PatientLoggedIn",
		LoginIP:     loginIP,
	}
}

func (e PatientLoggedInEvent) GetAggregateID() aggregate.PatientID {
	return e.AggregateID
}

func (e PatientLoggedInEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

func (e PatientLoggedInEvent) GetEventType() string {
	return e.EventType
}