package aggregate

import (
	"errors"
	"time"

	"intelligent-guidance-system/service/department/internal/domain/entity"
	"intelligent-guidance-system/service/department/internal/domain/event"
	"intelligent-guidance-system/service/department/internal/domain/vo"
)

var (
	ErrDepartmentNotFound     = errors.New("department not found")
	ErrDepartmentInactive     = errors.New("department is inactive")
	ErrInvalidDepartmentID    = errors.New("invalid department ID")
	ErrEmptyDepartmentName    = errors.New("department name cannot be empty")
	ErrInvalidPersonInCharge  = errors.New("invalid person in charge ID")
	ErrInvalidDutyDoctor      = errors.New("invalid duty doctor ID")
	ErrDoctorAlreadyAssigned  = errors.New("doctor already assigned")
	ErrDoctorNotInDepartment  = errors.New("doctor not in department")
	ErrCannotAssignToDisabled = errors.New("cannot assign doctor to disabled department")
)

type Department struct {
	id                  int64
	name                string
	introduction        string
	personInChargeID    int64
	phone               string
	address             string
	deptType            entity.DepartmentType
	staffCount          int
	status              entity.DepartmentStatus
	fee                 *vo.DepartmentFee
	dutyDoctorID        int64
	doctorAssignments   []*entity.DoctorAssignment
	events              []*event.DepartmentEvent
	createdAt           time.Time
	updatedAt           time.Time
}

func NewDepartment(
	name string,
	introduction string,
	personInChargeID int64,
	phone string,
	address string,
	deptType entity.DepartmentType,
	registrationFee int64,
	expertRegistrationFee int64,
) (*Department, error) {
	if name == "" {
		return nil, ErrEmptyDepartmentName
	}
	if deptType == entity.DepartmentTypeUnknown {
		return nil, entity.ErrInvalidDepartmentType
	}

	fee, err := vo.NewDepartmentFee(registrationFee, expertRegistrationFee, "CNY")
	if err != nil {
		return nil, err
	}

	now := time.Now()
	dept := &Department{
		id:                  0,
		name:                name,
		introduction:        introduction,
		personInChargeID:    personInChargeID,
		phone:               phone,
		address:             address,
		deptType:            deptType,
		staffCount:          0,
		status:              entity.DepartmentStatusActive,
		fee:                 fee,
		dutyDoctorID:        0,
		doctorAssignments:   make([]*entity.DoctorAssignment, 0),
		events:              make([]*event.DepartmentEvent, 0),
		createdAt:           now,
		updatedAt:           now,
	}

	dept.events = append(dept.events, event.DepartmentCreatedEvent(0, 0, name))
	return dept, nil
}

func ReconstructDepartment(
	id int64,
	name string,
	introduction string,
	personInChargeID int64,
	phone string,
	address string,
	deptType entity.DepartmentType,
	staffCount int,
	status entity.DepartmentStatus,
	fee *vo.DepartmentFee,
	dutyDoctorID int64,
	doctorAssignments []*entity.DoctorAssignment,
	createdAt time.Time,
	updatedAt time.Time,
) *Department {
	return &Department{
		id:                  id,
		name:                name,
		introduction:        introduction,
		personInChargeID:    personInChargeID,
		phone:               phone,
		address:             address,
		deptType:            deptType,
		staffCount:          staffCount,
		status:              status,
		fee:                 fee,
		dutyDoctorID:        dutyDoctorID,
		doctorAssignments:   doctorAssignments,
		events:              make([]*event.DepartmentEvent, 0),
		createdAt:           createdAt,
		updatedAt:           updatedAt,
	}
}

func (d *Department) ID() int64 { return d.id }
func (d *Department) Name() string { return d.name }
func (d *Department) Introduction() string { return d.introduction }
func (d *Department) PersonInChargeID() int64 { return d.personInChargeID }
func (d *Department) Phone() string { return d.phone }
func (d *Department) Address() string { return d.address }
func (d *Department) Type() entity.DepartmentType { return d.deptType }
func (d *Department) StaffCount() int { return d.staffCount }
func (d *Department) Status() entity.DepartmentStatus { return d.status }
func (d *Department) Fee() *vo.DepartmentFee { return d.fee }
func (d *Department) DutyDoctorID() int64 { return d.dutyDoctorID }
func (d *Department) DoctorAssignments() []*entity.DoctorAssignment { return d.doctorAssignments }
func (d *Department) Events() []*event.DepartmentEvent { return d.events }
func (d *Department) CreatedAt() time.Time { return d.createdAt }
func (d *Department) UpdatedAt() time.Time { return d.updatedAt }

func (d *Department) SetID(id int64) { d.id = id }

func (d *Department) AssignDoctor(doctorID int64, operatorID int64) error {
	if d.status == entity.DepartmentStatusDisabled {
		return ErrCannotAssignToDisabled
	}
	if doctorID <= 0 {
		return ErrInvalidDoctorID
	}

	for _, assignment := range d.doctorAssignments {
		if assignment.DoctorID() == doctorID && assignment.IsActive() {
			return ErrDoctorAlreadyAssigned
		}
	}

	assignment := entity.NewDoctorAssignment(doctorID)
	d.doctorAssignments = append(d.doctorAssignments, assignment)
	d.staffCount++
	d.updatedAt = time.Now()

	d.events = append(d.events, event.DoctorAssignedEvent(d.id, operatorID, doctorID))
	return nil
}

func (d *Department) RemoveDoctor(doctorID int64, operatorID int64) error {
	if doctorID <= 0 {
		return ErrInvalidDoctorID
	}

	for _, assignment := range d.doctorAssignments {
		if assignment.DoctorID() == doctorID && assignment.IsActive() {
			assignment.Deactivate()
			d.staffCount--
			if d.dutyDoctorID == doctorID {
				d.dutyDoctorID = 0
			}
			d.updatedAt = time.Now()
			d.events = append(d.events, event.DoctorRemovedEvent(d.id, operatorID, doctorID))
			return nil
		}
	}
	return ErrDoctorNotInDepartment
}

func (d *Department) SetDutyDoctor(doctorID int64, operatorID int64) error {
	if d.status == entity.DepartmentStatusDisabled {
		return ErrCannotAssignToDisabled
	}
	if doctorID <= 0 {
		return ErrInvalidDutyDoctor
	}

	found := false
	for _, assignment := range d.doctorAssignments {
		if assignment.DoctorID() == doctorID && assignment.IsActive() {
			found = true
			assignment.SetDuty(true)
			break
		}
	}

	if !found {
		return ErrDoctorNotInDepartment
	}

	d.dutyDoctorID = doctorID
	d.updatedAt = time.Now()
	d.events = append(d.events, event.DutyDoctorSetEvent(d.id, operatorID, doctorID))
	return nil
}

func (d *Department) Activate(operatorID int64) error {
	if d.status == entity.DepartmentStatusActive {
		return nil
	}
	oldStatus := d.status
	d.status = entity.DepartmentStatusActive
	d.updatedAt = time.Now()
	d.events = append(d.events, event.StatusChangedEvent(d.id, operatorID, oldStatus.Code(), d.status.Code()))
	return nil
}

func (d *Department) Rest(operatorID int64) error {
	if d.status == entity.DepartmentStatusDisabled {
		return ErrDepartmentInactive
	}
	oldStatus := d.status
	d.status = entity.DepartmentStatusResting
	d.updatedAt = time.Now()
	d.events = append(d.events, event.StatusChangedEvent(d.id, operatorID, oldStatus.Code(), d.status.Code()))
	return nil
}

func (d *Department) Disable(operatorID int64) error {
	oldStatus := d.status
	d.status = entity.DepartmentStatusDisabled
	d.updatedAt = time.Now()
	d.events = append(d.events, event.StatusChangedEvent(d.id, operatorID, oldStatus.Code(), d.status.Code()))
	return nil
}

func (d *Department) UpdateIntroduction(introduction string) {
	d.introduction = introduction
	d.updatedAt = time.Now()
}

func (d *Department) UpdatePhone(phone string) {
	d.phone = phone
	d.updatedAt = time.Now()
}

func (d *Department) UpdateAddress(address string) {
	d.address = address
	d.updatedAt = time.Now()
}

func (d *Department) UpdateFees(registrationFee, expertRegistrationFee int64) error {
	if d.fee == nil {
		d.fee = vo.ZeroDepartmentFee("CNY")
	}
	d.fee.SetRegistrationFee(registrationFee)
	d.fee.SetExpertRegistrationFee(expertRegistrationFee)
	d.updatedAt = time.Now()
	return nil
}

func (d *Department) UpdatePersonInCharge(personInChargeID int64) error {
	if personInChargeID <= 0 {
		return ErrInvalidPersonInCharge
	}
	d.personInChargeID = personInChargeID
	d.updatedAt = time.Now()
	return nil
}

func (d *Department) ClearEvents() {
	d.events = make([]*event.DepartmentEvent, 0)
}

func (d *Department) HasDoctor(doctorID int64) bool {
	for _, assignment := range d.doctorAssignments {
		if assignment.DoctorID() == doctorID && assignment.IsActive() {
			return true
		}
	}
	return false
}

func (d *Department) CanRegister() bool {
	return d.status.CanRegister()
}

func (d *Department) IsActive() bool {
	return d.status.IsActive()
}

func (d *Department) ActiveDoctorCount() int {
	count := 0
	for _, assignment := range d.doctorAssignments {
		if assignment.IsActive() {
			count++
		}
	}
	return count
}