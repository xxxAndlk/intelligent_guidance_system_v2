package entity

import (
	"errors"
	"time"
)

var (
	ErrInvalidFlowStep     = errors.New("invalid flow step")
	ErrFlowStepNotFound    = errors.New("flow step not found")
	ErrFlowStepAlreadyDone = errors.New("flow step already completed")
	ErrFlowStepNotStarted  = errors.New("flow step not started")
)

type SurgicalFlow struct {
	id           int64
	surgicalID   int64
	step         int
	name         string
	status       FlowStepStatus
	startTime    time.Time
	endTime      time.Time
	createdAt    time.Time
	updatedAt    time.Time
}

func NewSurgicalFlow(surgicalID int64, step int, name string) (*SurgicalFlow, error) {
	if surgicalID <= 0 {
		return nil, ErrInvalidSurgicalID
	}
	if step <= 0 {
		return nil, ErrInvalidFlowStep
	}
	if name == "" {
		return nil, errors.New("flow step name cannot be empty")
	}

	now := time.Now()
	return &SurgicalFlow{
		id:         0,
		surgicalID: surgicalID,
		step:       step,
		name:       name,
		status:     FlowStepStatusPending,
		createdAt:  now,
		updatedAt:  now,
	}, nil
}

func ReconstructSurgicalFlow(id, surgicalID int64, step int, name string, status FlowStepStatus, startTime, endTime, createdAt, updatedAt time.Time) *SurgicalFlow {
	return &SurgicalFlow{
		id:         id,
		surgicalID: surgicalID,
		step:       step,
		name:       name,
		status:     status,
		startTime:  startTime,
		endTime:    endTime,
		createdAt:  createdAt,
		updatedAt:  updatedAt,
	}
}

func (f *SurgicalFlow) ID() int64 { return f.id }
func (f *SurgicalFlow) SurgicalID() int64 { return f.surgicalID }
func (f *SurgicalFlow) Step() int { return f.step }
func (f *SurgicalFlow) Name() string { return f.name }
func (f *SurgicalFlow) Status() FlowStepStatus { return f.status }
func (f *SurgicalFlow) StartTime() time.Time { return f.startTime }
func (f *SurgicalFlow) EndTime() time.Time { return f.endTime }
func (f *SurgicalFlow) CreatedAt() time.Time { return f.createdAt }
func (f *SurgicalFlow) UpdatedAt() time.Time { return f.updatedAt }

func (f *SurgicalFlow) SetID(id int64) { f.id = id }

func (f *SurgicalFlow) Start() error {
	if !f.status.CanStart() {
		return ErrFlowStepAlreadyDone
	}
	f.status = FlowStepStatusInProgress
	f.startTime = time.Now()
	f.updatedAt = time.Now()
	return nil
}

func (f *SurgicalFlow) Complete() error {
	if f.status.IsTerminal() {
		return ErrFlowStepAlreadyDone
	}
	f.status = FlowStepStatusCompleted
	f.endTime = time.Now()
	f.updatedAt = time.Now()
	return nil
}

func (f *SurgicalFlow) Skip() error {
	if f.status.IsTerminal() {
		return ErrFlowStepAlreadyDone
	}
	f.status = FlowStepStatusSkipped
	f.updatedAt = time.Now()
	return nil
}

func (f *SurgicalFlow) Duration() time.Duration {
	if f.startTime.IsZero() || f.endTime.IsZero() {
		return 0
	}
	return f.endTime.Sub(f.startTime)
}