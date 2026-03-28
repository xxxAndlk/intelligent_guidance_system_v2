package aggregate

import (
	"errors"
	"time"

	"intelligent-guidance-system/service/surgical/internal/domain/entity"
	"intelligent-guidance-system/service/surgical/internal/domain/event"
)

var (
	ErrSurgicalNotFound        = errors.New("surgical not found")
	ErrSurgicalAlreadyStarted  = errors.New("surgical already started")
	ErrSurgicalAlreadyCompleted = errors.New("surgical already completed")
	ErrSurgicalAlreadyCancelled = errors.New("surgical already cancelled")
	ErrCannotStartSurgical     = errors.New("cannot start surgical")
	ErrCannotCancelSurgical    = errors.New("cannot cancel surgical")
	ErrCannotCompleteSurgical  = errors.New("cannot complete surgical")
	ErrFlowStepNotFound        = errors.New("flow step not found")
)

type Surgical struct {
	id            int64
	medicalID     int64
	patientID     int64
	doctorIDs     []int64
	departmentID  int64
	surgicalType  entity.SurgicalType
	status        entity.SurgicalStatus
	scheduledTime time.Time
	duration      time.Duration
	flows         []*entity.SurgicalFlow
	events        []*event.SurgicalEvent
	createdAt     time.Time
	updatedAt     time.Time
}

func NewSurgical(
	medicalID int64,
	patientID int64,
	doctorIDs []int64,
	departmentID int64,
	surgicalType entity.SurgicalType,
	scheduledTime time.Time,
) (*Surgical, error) {
	if medicalID <= 0 {
		return nil, entity.ErrInvalidMedicalID
	}
	if patientID <= 0 {
		return nil, entity.ErrInvalidPatientID
	}
	if len(doctorIDs) == 0 {
		return nil, errors.New("at least one doctor is required")
	}
	if departmentID <= 0 {
		return nil, entity.ErrInvalidDepartmentID
	}
	if surgicalType == entity.SurgicalTypeUnknown {
		return nil, entity.ErrInvalidSurgicalType
	}

	now := time.Now()
	surg := &Surgical{
		id:            0,
		medicalID:     medicalID,
		patientID:     patientID,
		doctorIDs:     doctorIDs,
		departmentID:  departmentID,
		surgicalType:  surgicalType,
		status:        entity.SurgicalStatusScheduled,
		scheduledTime: scheduledTime,
		duration:      0,
		flows:         make([]*entity.SurgicalFlow, 0),
		events:        make([]*event.SurgicalEvent, 0),
		createdAt:     now,
		updatedAt:     now,
	}

	surg.events = append(surg.events, event.SurgicalCreatedEvent(0, patientID))
	return surg, nil
}

func ReconstructSurgical(
	id int64,
	medicalID int64,
	patientID int64,
	doctorIDs []int64,
	departmentID int64,
	surgicalType entity.SurgicalType,
	status entity.SurgicalStatus,
	scheduledTime time.Time,
	duration time.Duration,
	flows []*entity.SurgicalFlow,
	createdAt time.Time,
	updatedAt time.Time,
) *Surgical {
	return &Surgical{
		id:            id,
		medicalID:     medicalID,
		patientID:     patientID,
		doctorIDs:     doctorIDs,
		departmentID:  departmentID,
		surgicalType:  surgicalType,
		status:        status,
		scheduledTime: scheduledTime,
		duration:      duration,
		flows:         flows,
		events:        make([]*event.SurgicalEvent, 0),
		createdAt:     createdAt,
		updatedAt:     updatedAt,
	}
}

func (s *Surgical) ID() int64 { return s.id }
func (s *Surgical) MedicalID() int64 { return s.medicalID }
func (s *Surgical) PatientID() int64 { return s.patientID }
func (s *Surgical) DoctorIDs() []int64 { return s.doctorIDs }
func (s *Surgical) DepartmentID() int64 { return s.departmentID }
func (s *Surgical) Type() entity.SurgicalType { return s.surgicalType }
func (s *Surgical) Status() entity.SurgicalStatus { return s.status }
func (s *Surgical) ScheduledTime() time.Time { return s.scheduledTime }
func (s *Surgical) Duration() time.Duration { return s.duration }
func (s *Surgical) Flows() []*entity.SurgicalFlow { return s.flows }
func (s *Surgical) Events() []*event.SurgicalEvent { return s.events }
func (s *Surgical) CreatedAt() time.Time { return s.createdAt }
func (s *Surgical) UpdatedAt() time.Time { return s.updatedAt }

func (s *Surgical) SetID(id int64) { s.id = id }

func (s *Surgical) StartSurgical(operatorID int64) error {
	if !s.status.CanStart() {
		if s.status.IsTerminal() {
			if s.status == entity.SurgicalStatusInProgress {
				return ErrSurgicalAlreadyStarted
			}
			if s.status == entity.SurgicalStatusCompleted {
				return ErrSurgicalAlreadyCompleted
			}
			return ErrSurgicalAlreadyCancelled
		}
		return ErrCannotStartSurgical
	}

	s.status = entity.SurgicalStatusInProgress
	s.updatedAt = time.Now()
	s.events = append(s.events, event.SurgicalStartedEvent(s.id, operatorID))
	return nil
}

func (s *Surgical) CompleteStep(step int, operatorID int64) error {
	for _, flow := range s.flows {
		if flow.Step() == step {
			if err := flow.Complete(); err != nil {
				return err
			}
			s.updatedAt = time.Now()
			s.events = append(s.events, event.SurgicalStepCompletedEvent(s.id, operatorID, step))
			return nil
		}
	}
	return ErrFlowStepNotFound
}

func (s *Surgical) CompleteSurgical(operatorID int64) error {
	if s.status != entity.SurgicalStatusInProgress {
		if s.status.IsTerminal() {
			if s.status == entity.SurgicalStatusCompleted {
				return ErrSurgicalAlreadyCompleted
			}
			return ErrSurgicalAlreadyCancelled
		}
		return ErrCannotCompleteSurgical
	}

	s.status = entity.SurgicalStatusCompleted
	s.updatedAt = time.Now()
	s.events = append(s.events, event.SurgicalCompletedEvent(s.id, operatorID))
	return nil
}

func (s *Surgical) CancelSurgical(operatorID int64, reason string) error {
	if !s.status.CanCancel() {
		if s.status.IsTerminal() {
			if s.status == entity.SurgicalStatusCancelled {
				return ErrSurgicalAlreadyCancelled
			}
			return ErrSurgicalAlreadyCompleted
		}
		return ErrCannotCancelSurgical
	}

	s.status = entity.SurgicalStatusCancelled
	s.updatedAt = time.Now()
	s.events = append(s.events, event.SurgicalCancelledEvent(s.id, operatorID, reason))
	return nil
}

func (s *Surgical) AddFlowStep(step int, name string) error {
	flow, err := entity.NewSurgicalFlow(s.id, step, name)
	if err != nil {
		return err
	}
	s.flows = append(s.flows, flow)
	s.updatedAt = time.Now()
	return nil
}

func (s *Surgical) ClearEvents() {
	s.events = make([]*event.SurgicalEvent, 0)
}

func (s *Surgical) IsTerminal() bool {
	return s.status.IsTerminal()
}

func (s *Surgical) IsInProgress() bool {
	return s.status == entity.SurgicalStatusInProgress
}

func (s *Surgical) TotalDuration() time.Duration {
	total := time.Duration(0)
	for _, flow := range s.flows {
		total += flow.Duration()
	}
	return total
}