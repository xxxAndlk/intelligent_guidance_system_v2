package aggregate

import (
	"errors"
	"time"

	"intelligent-guidance-system/service/registration/internal/domain/entity"
	"intelligent-guidance-system/service/registration/internal/domain/event"
	"intelligent-guidance-system/service/registration/internal/domain/vo"
)

var (
	ErrRegistrationNotFound      = errors.New("registration not found")
	ErrRegistrationAlreadyConfirmed = errors.New("registration already confirmed")
	ErrRegistrationAlreadyCancelled = errors.New("registration already cancelled")
	ErrRegistrationAlreadyCompleted = errors.New("registration already completed")
	ErrCannotCancelRegistration  = errors.New("cannot cancel registration")
	ErrCannotConfirmRegistration = errors.New("cannot confirm registration")
	ErrCannotCompleteRegistration = errors.New("cannot complete registration")
)

type Registration struct {
	id              int64
	patientID       int64
	doctorID        int64
	departmentID    int64
	regType         entity.RegistrationType
	status          entity.RegistrationStatus
	appointmentTime time.Time
	fee             *vo.RegistrationFee
	createdAt       time.Time
	updatedAt       time.Time
	events          []*event.RegistrationEvent
}

func NewRegistration(
	patientID int64,
	doctorID int64,
	departmentID int64,
	regType entity.RegistrationType,
	appointmentTime time.Time,
	feeAmount int64,
) (*Registration, error) {
	if patientID <= 0 {
		return nil, entity.ErrInvalidPatientID
	}
	if doctorID <= 0 {
		return nil, entity.ErrInvalidDoctorID
	}
	if departmentID <= 0 {
		return nil, entity.ErrInvalidDepartmentID
	}
	if regType == entity.RegistrationTypeUnknown {
		return nil, entity.ErrInvalidRegistrationType
	}

	fee, err := vo.NewRegistrationFee(feeAmount, "CNY")
	if err != nil {
		return nil, err
	}

	now := time.Now()
	reg := &Registration{
		id:              0,
		patientID:       patientID,
		doctorID:        doctorID,
		departmentID:    departmentID,
		regType:         regType,
		status:          entity.RegistrationStatusPending,
		appointmentTime: appointmentTime,
		fee:             fee,
		createdAt:       now,
		updatedAt:       now,
		events:          make([]*event.RegistrationEvent, 0),
	}

	reg.events = append(reg.events, event.RegistrationCreatedEvent(0, patientID, doctorID))
	return reg, nil
}

func ReconstructRegistration(
	id int64,
	patientID int64,
	doctorID int64,
	departmentID int64,
	regType entity.RegistrationType,
	status entity.RegistrationStatus,
	appointmentTime time.Time,
	fee *vo.RegistrationFee,
	createdAt time.Time,
	updatedAt time.Time,
) *Registration {
	return &Registration{
		id:              id,
		patientID:       patientID,
		doctorID:        doctorID,
		departmentID:    departmentID,
		regType:         regType,
		status:          status,
		appointmentTime: appointmentTime,
		fee:             fee,
		createdAt:       createdAt,
		updatedAt:       updatedAt,
		events:          make([]*event.RegistrationEvent, 0),
	}
}

func (r *Registration) ID() int64 { return r.id }
func (r *Registration) PatientID() int64 { return r.patientID }
func (r *Registration) DoctorID() int64 { return r.doctorID }
func (r *Registration) DepartmentID() int64 { return r.departmentID }
func (r *Registration) Type() entity.RegistrationType { return r.regType }
func (r *Registration) Status() entity.RegistrationStatus { return r.status }
func (r *Registration) AppointmentTime() time.Time { return r.appointmentTime }
func (r *Registration) Fee() *vo.RegistrationFee { return r.fee }
func (r *Registration) CreatedAt() time.Time { return r.createdAt }
func (r *Registration) UpdatedAt() time.Time { return r.updatedAt }
func (r *Registration) Events() []*event.RegistrationEvent { return r.events }

func (r *Registration) SetID(id int64) { r.id = id }

func (r *Registration) Confirm(operatorID int64) error {
	if !r.status.CanConfirm() {
		if r.status.IsTerminal() {
			if r.status == entity.RegistrationStatusConfirmed {
				return ErrRegistrationAlreadyConfirmed
			}
			if r.status == entity.RegistrationStatusCancelled {
				return ErrRegistrationAlreadyCancelled
			}
			return ErrRegistrationAlreadyCompleted
		}
		return ErrCannotConfirmRegistration
	}

	r.status = entity.RegistrationStatusConfirmed
	r.updatedAt = time.Now()
	r.events = append(r.events, event.RegistrationConfirmedEvent(r.id, operatorID))
	return nil
}

func (r *Registration) Cancel(operatorID int64, reason string) error {
	if !r.status.CanCancel() {
		if r.status.IsTerminal() {
			if r.status == entity.RegistrationStatusCancelled {
				return ErrRegistrationAlreadyCancelled
			}
			return ErrRegistrationAlreadyCompleted
		}
		return ErrCannotCancelRegistration
	}

	r.status = entity.RegistrationStatusCancelled
	r.updatedAt = time.Now()
	r.events = append(r.events, event.RegistrationCancelledEvent(r.id, operatorID, reason))
	return nil
}

func (r *Registration) Complete(operatorID int64) error {
	if !r.status.CanComplete() {
		if r.status.IsTerminal() {
			if r.status == entity.RegistrationStatusCompleted {
				return ErrRegistrationAlreadyCompleted
			}
			if r.status == entity.RegistrationStatusCancelled {
				return ErrRegistrationAlreadyCancelled
			}
		}
		return ErrCannotCompleteRegistration
	}

	r.status = entity.RegistrationStatusCompleted
	r.updatedAt = time.Now()
	r.events = append(r.events, event.RegistrationCompletedEvent(r.id, operatorID))
	return nil
}

func (r *Registration) ClearEvents() {
	r.events = make([]*event.RegistrationEvent, 0)
}

func (r *Registration) IsTerminal() bool {
	return r.status.IsTerminal()
}

func (r *Registration) IsExpert() bool {
	return r.regType == entity.RegistrationTypeExpert
}