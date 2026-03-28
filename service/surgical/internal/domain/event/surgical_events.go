package event

import "time"

type SurgicalEvent struct {
	eventType   string
	surgicalID  int64
	operatorID  int64
	data        map[string]interface{}
	occurredAt  time.Time
}

func NewSurgicalEvent(eventType string, surgicalID, operatorID int64, data map[string]interface{}) *SurgicalEvent {
	return &SurgicalEvent{
		eventType:  eventType,
		surgicalID: surgicalID,
		operatorID: operatorID,
		data:       data,
		occurredAt: time.Now(),
	}
}

func (e *SurgicalEvent) EventType() string { return e.eventType }
func (e *SurgicalEvent) SurgicalID() int64 { return e.surgicalID }
func (e *SurgicalEvent) OperatorID() int64 { return e.operatorID }
func (e *SurgicalEvent) Data() map[string]interface{} { return e.data }
func (e *SurgicalEvent) OccurredAt() time.Time { return e.occurredAt }

const (
	EventTypeSurgicalCreated    = "surgical.created"
	EventTypeSurgicalStarted     = "surgical.started"
	EventTypeSurgicalStepCompleted = "surgical.step_completed"
	EventTypeSurgicalCompleted   = "surgical.completed"
	EventTypeSurgicalCancelled   = "surgical.cancelled"
)

func SurgicalCreatedEvent(surgicalID, patientID int64) *SurgicalEvent {
	return NewSurgicalEvent(EventTypeSurgicalCreated, surgicalID, patientID, nil)
}

func SurgicalStartedEvent(surgicalID, operatorID int64) *SurgicalEvent {
	return NewSurgicalEvent(EventTypeSurgicalStarted, surgicalID, operatorID, nil)
}

func SurgicalStepCompletedEvent(surgicalID, operatorID int64, step int) *SurgicalEvent {
	return NewSurgicalEvent(EventTypeSurgicalStepCompleted, surgicalID, operatorID, map[string]interface{}{
		"step": step,
	})
}

func SurgicalCompletedEvent(surgicalID, operatorID int64) *SurgicalEvent {
	return NewSurgicalEvent(EventTypeSurgicalCompleted, surgicalID, operatorID, nil)
}

func SurgicalCancelledEvent(surgicalID, operatorID int64, reason string) *SurgicalEvent {
	return NewSurgicalEvent(EventTypeSurgicalCancelled, surgicalID, operatorID, map[string]interface{}{
		"reason": reason,
	})
}