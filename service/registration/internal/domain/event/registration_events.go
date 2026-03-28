package event

import "time"

type RegistrationEvent struct {
	eventType     string
	registrationID int64
	operatorID    int64
	data          map[string]interface{}
	occurredAt    time.Time
}

func NewRegistrationEvent(eventType string, registrationID, operatorID int64, data map[string]interface{}) *RegistrationEvent {
	return &RegistrationEvent{
		eventType:     eventType,
		registrationID: registrationID,
		operatorID:    operatorID,
		data:          data,
		occurredAt:    time.Now(),
	}
}

func (e *RegistrationEvent) EventType() string { return e.eventType }
func (e *RegistrationEvent) RegistrationID() int64 { return e.registrationID }
func (e *RegistrationEvent) OperatorID() int64 { return e.operatorID }
func (e *RegistrationEvent) Data() map[string]interface{} { return e.data }
func (e *RegistrationEvent) OccurredAt() time.Time { return e.occurredAt }

const (
	EventTypeRegistrationCreated   = "registration.created"
	EventTypeRegistrationConfirmed  = "registration.confirmed"
	EventTypeRegistrationCancelled  = "registration.cancelled"
	EventTypeRegistrationCompleted  = "registration.completed"
)

func RegistrationCreatedEvent(registrationID, patientID, doctorID int64) *RegistrationEvent {
	return NewRegistrationEvent(EventTypeRegistrationCreated, registrationID, patientID, map[string]interface{}{
		"patient_id": patientID,
		"doctor_id":  doctorID,
	})
}

func RegistrationConfirmedEvent(registrationID, operatorID int64) *RegistrationEvent {
	return NewRegistrationEvent(EventTypeRegistrationConfirmed, registrationID, operatorID, nil)
}

func RegistrationCancelledEvent(registrationID, operatorID int64, reason string) *RegistrationEvent {
	return NewRegistrationEvent(EventTypeRegistrationCancelled, registrationID, operatorID, map[string]interface{}{
		"reason": reason,
	})
}

func RegistrationCompletedEvent(registrationID, operatorID int64) *RegistrationEvent {
	return NewRegistrationEvent(EventTypeRegistrationCompleted, registrationID, operatorID, nil)
}