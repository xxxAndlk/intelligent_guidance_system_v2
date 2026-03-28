package entity

import "time"

// DoctorAssignment represents a doctor assigned to a department
type DoctorAssignment struct {
	doctorID    int64
	assignedAt  time.Time
	isDuty      bool
	isActive    bool
}

func NewDoctorAssignment(doctorID int64) *DoctorAssignment {
	return &DoctorAssignment{
		doctorID:   doctorID,
		assignedAt: time.Now(),
		isDuty:     false,
		isActive:   true,
	}
}

func ReconstructDoctorAssignment(doctorID int64, assignedAt time.Time, isDuty, isActive bool) *DoctorAssignment {
	return &DoctorAssignment{
		doctorID:   doctorID,
		assignedAt: assignedAt,
		isDuty:     isDuty,
		isActive:   isActive,
	}
}

func (d *DoctorAssignment) DoctorID() int64 {
	return d.doctorID
}

func (d *DoctorAssignment) AssignedAt() time.Time {
	return d.assignedAt
}

func (d *DoctorAssignment) IsDuty() bool {
	return d.isDuty
}

func (d *DoctorAssignment) IsActive() bool {
	return d.isActive
}

func (d *DoctorAssignment) SetDuty(isDuty bool) {
	d.isDuty = isDuty
}

func (d *DoctorAssignment) Deactivate() {
	d.isActive = false
}

func (d *DoctorAssignment) Activate() {
	d.isActive = true
}