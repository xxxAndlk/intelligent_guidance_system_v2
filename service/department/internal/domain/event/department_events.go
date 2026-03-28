package event

import "time"

// DepartmentEvent represents a domain event for department
type DepartmentEvent struct {
	eventType   string
	departmentID int64
	operatorID  int64
	data        map[string]interface{}
	occurredAt  time.Time
}

func NewDepartmentEvent(eventType string, departmentID, operatorID int64, data map[string]interface{}) *DepartmentEvent {
	return &DepartmentEvent{
		eventType:   eventType,
		departmentID: departmentID,
		operatorID:  operatorID,
		data:        data,
		occurredAt:  time.Now(),
	}
}

func (e *DepartmentEvent) EventType() string {
	return e.eventType
}

func (e *DepartmentEvent) DepartmentID() int64 {
	return e.departmentID
}

func (e *DepartmentEvent) OperatorID() int64 {
	return e.operatorID
}

func (e *DepartmentEvent) Data() map[string]interface{} {
	return e.data
}

func (e *DepartmentEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// Event type constants
const (
	EventTypeDepartmentCreated     = "department.created"
	EventTypeDepartmentUpdated     = "department.updated"
	EventTypeDepartmentDisabled    = "department.disabled"
	EventTypeDepartmentActivated   = "department.activated"
	EventTypeDoctorAssigned        = "doctor.assigned"
	EventTypeDoctorRemoved         = "doctor.removed"
	EventTypeDutyDoctorSet         = "duty_doctor.set"
	EventTypeStatusChanged         = "department.status_changed"
)

// Helper functions to create specific events
func DepartmentCreatedEvent(departmentID, operatorID int64, name string) *DepartmentEvent {
	return NewDepartmentEvent(EventTypeDepartmentCreated, departmentID, operatorID, map[string]interface{}{
		"name": name,
	})
}

func DoctorAssignedEvent(departmentID, operatorID, doctorID int64) *DepartmentEvent {
	return NewDepartmentEvent(EventTypeDoctorAssigned, departmentID, operatorID, map[string]interface{}{
		"doctor_id": doctorID,
	})
}

func DoctorRemovedEvent(departmentID, operatorID, doctorID int64) *DepartmentEvent {
	return NewDepartmentEvent(EventTypeDoctorRemoved, departmentID, operatorID, map[string]interface{}{
		"doctor_id": doctorID,
	})
}

func DutyDoctorSetEvent(departmentID, operatorID, doctorID int64) *DepartmentEvent {
	return NewDepartmentEvent(EventTypeDutyDoctorSet, departmentID, operatorID, map[string]interface{}{
		"doctor_id": doctorID,
	})
}

func StatusChangedEvent(departmentID, operatorID int64, oldStatus, newStatus string) *DepartmentEvent {
	return NewDepartmentEvent(EventTypeStatusChanged, departmentID, operatorID, map[string]interface{}{
		"old_status": oldStatus,
		"new_status": newStatus,
	})
}