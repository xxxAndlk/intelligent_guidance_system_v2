package event

import (
	"time"
)

type DoctorEventType string

const (
	DoctorEventCreated        DoctorEventType = "doctor.created"
	DoctorEventUpdated        DoctorEventType = "doctor.updated"
	DoctorEventActivated      DoctorEventType = "doctor.activated"
	DoctorEventDeactivated    DoctorEventType = "doctor.deactivated"
	DoctorEventDepartmentChanged DoctorEventType = "doctor.department_changed"
	DoctorEventProfessionalUpdated DoctorEventType = "doctor.professional_updated"
	DoctorEventRoleAssigned   DoctorEventType = "doctor.role_assigned"
	DoctorEventRoleRemoved    DoctorEventType = "doctor.role_removed"
	DoctorEventScheduleSet    DoctorEventType = "doctor.schedule_set"
	DoctorEventSchedulePublished DoctorEventType = "doctor.schedule_published"
	DoctorEventScheduleCancelled DoctorEventType = "doctor.schedule_cancelled"
	DoctorEventStatusChanged  DoctorEventType = "doctor.status_changed"
)

type DoctorEvent struct {
	eventID     string
	eventType   DoctorEventType
	doctorID    string
	employeeID  string
	departmentID string
	occurredAt  time.Time
	payload     map[string]interface{}
}

func NewDoctorEvent(eventType DoctorEventType, doctorID, employeeID string, payload map[string]interface{}) *DoctorEvent {
	return &DoctorEvent{
		eventID:    generateEventID(),
		eventType:  eventType,
		doctorID:   doctorID,
		employeeID: employeeID,
		occurredAt: time.Now(),
		payload:    payload,
	}
}

func generateEventID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().Nanosecond()%len(letters)]
	}
	return string(b)
}

func (e *DoctorEvent) EventID() string                  { return e.eventID }
func (e *DoctorEvent) EventType() DoctorEventType       { return e.eventType }
func (e *DoctorEvent) DoctorID() string                 { return e.doctorID }
func (e *DoctorEvent) EmployeeID() string               { return e.employeeID }
func (e *DoctorEvent) DepartmentID() string             { return e.departmentID }
func (e *DoctorEvent) OccurredAt() time.Time            { return e.occurredAt }
func (e *DoctorEvent) Payload() map[string]interface{} { return e.payload }

type DoctorCreatedEvent struct {
	*DoctorEvent
	Name       string
	Position   string
	DepartmentID string
}

type DoctorDepartmentChangedEvent struct {
	*DoctorEvent
	OldDepartmentID string
	NewDepartmentID string
	Reason          string
}

type DoctorStatusChangedEvent struct {
	*DoctorEvent
	OldStatus string
	NewStatus string
	Reason    string
}

type DoctorRoleAssignedEvent struct {
	*DoctorEvent
	RoleID     string
	RoleName   string
	AssignedBy string
	ExpiresAt  *time.Time
}

type DoctorScheduleSetEvent struct {
	*DoctorEvent
	ScheduleDate string
	TotalSlots   int
	TotalCapacity int
}

type DoctorSchedulePublishedEvent struct {
	*DoctorEvent
	ScheduleDate string
	AvailableCapacity int
}

func NewDoctorCreatedEvent(doctorID, employeeID, name, position, deptID string) *DoctorCreatedEvent {
	return &DoctorCreatedEvent{
		DoctorEvent: NewDoctorEvent(DoctorEventCreated, doctorID, employeeID, map[string]interface{}{
			"name":       name,
			"position":   position,
			"department": deptID,
		}),
		Name:         name,
		Position:     position,
		DepartmentID: deptID,
	}
}

func NewDoctorDepartmentChangedEvent(doctorID, employeeID, oldDept, newDept, reason string) *DoctorDepartmentChangedEvent {
	return &DoctorDepartmentChangedEvent{
		DoctorEvent: NewDoctorEvent(DoctorEventDepartmentChanged, doctorID, employeeID, map[string]interface{}{
			"old_department": oldDept,
			"new_department": newDept,
			"reason":         reason,
		}),
		OldDepartmentID: oldDept,
		NewDepartmentID: newDept,
		Reason:          reason,
	}
}

func NewDoctorStatusChangedEvent(doctorID, employeeID, oldStatus, newStatus, reason string) *DoctorStatusChangedEvent {
	return &DoctorStatusChangedEvent{
		DoctorEvent: NewDoctorEvent(DoctorEventStatusChanged, doctorID, employeeID, map[string]interface{}{
			"old_status": oldStatus,
			"new_status": newStatus,
			"reason":     reason,
		}),
		OldStatus: oldStatus,
		NewStatus: newStatus,
		Reason:    reason,
	}
}

func NewDoctorRoleAssignedEvent(doctorID, employeeID, roleID, roleName, assignedBy string, expiresAt *time.Time) *DoctorRoleAssignedEvent {
	payload := map[string]interface{}{
		"role_id":     roleID,
		"role_name":   roleName,
		"assigned_by": assignedBy,
	}
	if expiresAt != nil {
		payload["expires_at"] = expiresAt.Format(time.RFC3339)
	}
	return &DoctorRoleAssignedEvent{
		DoctorEvent: NewDoctorEvent(DoctorEventRoleAssigned, doctorID, employeeID, payload),
		RoleID:      roleID,
		RoleName:    roleName,
		AssignedBy:  assignedBy,
		ExpiresAt:   expiresAt,
	}
}

func NewDoctorScheduleSetEvent(doctorID, employeeID, date string, totalSlots, totalCapacity int) *DoctorScheduleSetEvent {
	return &DoctorScheduleSetEvent{
		DoctorEvent: NewDoctorEvent(DoctorEventScheduleSet, doctorID, employeeID, map[string]interface{}{
			"date":           date,
			"total_slots":    totalSlots,
			"total_capacity": totalCapacity,
		}),
		ScheduleDate:  date,
		TotalSlots:    totalSlots,
		TotalCapacity: totalCapacity,
	}
}

func NewDoctorSchedulePublishedEvent(doctorID, employeeID, date string, availableCapacity int) *DoctorSchedulePublishedEvent {
	return &DoctorSchedulePublishedEvent{
		DoctorEvent: NewDoctorEvent(DoctorEventSchedulePublished, doctorID, employeeID, map[string]interface{}{
			"date":              date,
			"available_capacity": availableCapacity,
		}),
		ScheduleDate:     date,
		AvailableCapacity: availableCapacity,
	}
}