package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidTimeSlot     = errors.New("invalid time slot")
	ErrTimeSlotOverlap     = errors.New("time slots overlap")
	ErrScheduleNotFound    = errors.New("schedule not found")
	ErrScheduleExpired     = errors.New("schedule has expired")
	ErrInvalidScheduleDate = errors.New("invalid schedule date")
)

// TimeSlot represents a single time slot within a schedule
type TimeSlot struct {
	id          string
	startTime   time.Time
	endTime     time.Time
	isAvailable bool
	maxPatients int
	bookedCount int
}

// NewTimeSlot creates a new time slot
func NewTimeSlot(startTime, endTime time.Time, maxPatients int) (*TimeSlot, error) {
	if startTime.After(endTime) || startTime.Equal(endTime) {
		return nil, ErrInvalidTimeSlot
	}
	if maxPatients <= 0 {
		maxPatients = 1
	}

	return &TimeSlot{
		id:          uuid.New().String(),
		startTime:   startTime,
		endTime:     endTime,
		isAvailable: true,
		maxPatients: maxPatients,
		bookedCount: 0,
	}, nil
}

// ReconstructTimeSlot reconstructs a time slot from persistence
func ReconstructTimeSlot(
	id string,
	startTime, endTime time.Time,
	isAvailable bool,
	maxPatients, bookedCount int,
) *TimeSlot {
	return &TimeSlot{
		id:          id,
		startTime:   startTime,
		endTime:     endTime,
		isAvailable: isAvailable,
		maxPatients: maxPatients,
		bookedCount: bookedCount,
	}
}

func (t *TimeSlot) ID() string          { return t.id }
func (t *TimeSlot) StartTime() time.Time { return t.startTime }
func (t *TimeSlot) EndTime() time.Time   { return t.endTime }
func (t *TimeSlot) IsAvailable() bool    { return t.isAvailable && t.bookedCount < t.maxPatients }
func (t *TimeSlot) MaxPatients() int      { return t.maxPatients }
func (t *TimeSlot) BookedCount() int      { return t.bookedCount }
func (t *TimeSlot) AvailableSlots() int    { return t.maxPatients - t.bookedCount }

// Book books a slot, returns error if not available
func (t *TimeSlot) Book() error {
	if !t.IsAvailable() {
		return ErrTimeSlotOverlap
	}
	t.bookedCount++
	if t.bookedCount >= t.maxPatients {
		t.isAvailable = false
	}
	return nil
}

// Cancel cancels a booking
func (t *TimeSlot) Cancel() {
	if t.bookedCount > 0 {
		t.bookedCount--
	}
	t.isAvailable = true
}

// SetMaxPatients sets the maximum patients for this slot
func (t *TimeSlot) SetMaxPatients(max int) {
	t.maxPatients = max
	if t.bookedCount < t.maxPatients {
		t.isAvailable = true
	}
}

// Duration returns the duration of the time slot
func (t *TimeSlot) Duration() time.Duration {
	return t.endTime.Sub(t.startTime)
}

// Contains checks if the time slot contains a given time
func (t *TimeSlot) Contains(tm time.Time) bool {
	return (tm.Equal(t.startTime) || tm.After(t.startTime)) &&
		(tm.Equal(t.endTime) || tm.Before(t.endTime))
}

// Overlaps checks if two time slots overlap
func (t *TimeSlot) Overlaps(other *TimeSlot) bool {
	return t.startTime.Before(other.endTime) && t.endTime.After(other.startTime)
}

// Schedule represents a doctor's schedule for a specific date
type Schedule struct {
	id        string
	doctorID  string
	date      time.Time
	timeSlots []*TimeSlot
	status    ScheduleStatus
	createdAt time.Time
	updatedAt time.Time
}

type ScheduleStatus int

const (
	ScheduleStatusDraft ScheduleStatus = iota
	ScheduleStatusPublished
	ScheduleStatusCancelled
)

func (s ScheduleStatus) String() string {
	switch s {
	case ScheduleStatusDraft:
		return "草稿"
	case ScheduleStatusPublished:
		return "已发布"
	case ScheduleStatusCancelled:
		return "已取消"
	default:
		return "未知"
	}
}

// NewSchedule creates a new schedule
func NewSchedule(doctorID string, date time.Time) (*Schedule, error) {
	if date.Before(time.Now().Truncate(24 * time.Hour)) {
		return nil, ErrInvalidScheduleDate
	}

	return &Schedule{
		id:        uuid.New().String(),
		doctorID:  doctorID,
		date:      date.Truncate(24 * time.Hour),
		timeSlots: make([]*TimeSlot, 0),
		status:    ScheduleStatusDraft,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}, nil
}

// ReconstructSchedule reconstructs a schedule from persistence
func ReconstructSchedule(
	id, doctorID string,
	date time.Time,
	timeSlots []*TimeSlot,
	status ScheduleStatus,
	createdAt, updatedAt time.Time,
) *Schedule {
	return &Schedule{
		id:        id,
		doctorID:  doctorID,
		date:      date,
		timeSlots: timeSlots,
		status:    status,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (s *Schedule) ID() string            { return s.id }
func (s *Schedule) DoctorID() string      { return s.doctorID }
func (s *Schedule) Date() time.Time       { return s.date }
func (s *Schedule) TimeSlots() []*TimeSlot { return s.timeSlots }
func (s *Schedule) Status() ScheduleStatus { return s.status }
func (s *Schedule) CreatedAt() time.Time  { return s.createdAt }
func (s *Schedule) UpdatedAt() time.Time  { return s.updatedAt }

// AddTimeSlot adds a new time slot to the schedule
func (s *Schedule) AddTimeSlot(startTime, endTime time.Time, maxPatients int) error {
	newSlot, err := NewTimeSlot(startTime, endTime, maxPatients)
	if err != nil {
		return err
	}

	for _, existing := range s.timeSlots {
		if newSlot.Overlaps(existing) {
			return ErrTimeSlotOverlap
		}
	}

	s.timeSlots = append(s.timeSlots, newSlot)
	s.updatedAt = time.Now()
	return nil
}

// RemoveTimeSlot removes a time slot by ID
func (s *Schedule) RemoveTimeSlot(slotID string) error {
	for i, slot := range s.timeSlots {
		if slot.ID() == slotID {
			s.timeSlots = append(s.timeSlots[:i], s.timeSlots[i+1:]...)
			s.updatedAt = time.Now()
			return nil
		}
	}
	return ErrScheduleNotFound
}

// GetTimeSlot gets a time slot by ID
func (s *Schedule) GetTimeSlot(slotID string) (*TimeSlot, error) {
	for _, slot := range s.timeSlots {
		if slot.ID() == slotID {
			return slot, nil
		}
	}
	return nil, ErrScheduleNotFound
}

// BookSlot books a specific time slot
func (s *Schedule) BookSlot(slotID string) error {
	slot, err := s.GetTimeSlot(slotID)
	if err != nil {
		return err
	}
	return slot.Book()
}

// CancelSlotBooking cancels a booking for a specific time slot
func (s *Schedule) CancelSlotBooking(slotID string) error {
	slot, err := s.GetTimeSlot(slotID)
	if err != nil {
		return err
	}
	slot.Cancel()
	s.updatedAt = time.Now()
	return nil
}

// Publish publishes the schedule
func (s *Schedule) Publish() {
	s.status = ScheduleStatusPublished
	s.updatedAt = time.Now()
}

// Cancel cancels the schedule
func (s *Schedule) Cancel() {
	s.status = ScheduleStatusCancelled
	s.updatedAt = time.Now()
}

// IsExpired checks if the schedule date has passed
func (s *Schedule) IsExpired() bool {
	return s.date.Before(time.Now().Truncate(24 * time.Hour))
}

// AvailableSlots returns all available time slots
func (s *Schedule) AvailableSlots() []*TimeSlot {
	result := make([]*TimeSlot, 0)
	for _, slot := range s.timeSlots {
		if slot.IsAvailable() {
			result = append(result, slot)
		}
	}
	return result
}

// TotalCapacity returns the total patient capacity
func (s *Schedule) TotalCapacity() int {
	total := 0
	for _, slot := range s.timeSlots {
		total += slot.MaxPatients()
	}
	return total
}

// TotalBooked returns the total booked count
func (s *Schedule) TotalBooked() int {
	total := 0
	for _, slot := range s.timeSlots {
		total += slot.BookedCount()
	}
	return total
}

// AvailableCapacity returns the remaining available capacity
func (s *Schedule) AvailableCapacity() int {
	return s.TotalCapacity() - s.TotalBooked()
}

// SetTimeSlots replaces all time slots
func (s *Schedule) SetTimeSlots(slots []struct {
	StartTime   time.Time
	EndTime     time.Time
	MaxPatients int
}) error {
	newSlots := make([]*TimeSlot, 0, len(slots))
	for i, slot := range slots {
		newSlot, err := NewTimeSlot(slot.StartTime, slot.EndTime, slot.MaxPatients)
		if err != nil {
			return err
		}
		for j, existing := range newSlots {
			if newSlot.Overlaps(existing) {
				return ErrTimeSlotOverlap
			}
			if j >= i {
				break
			}
		}
		newSlots = append(newSlots, newSlot)
	}
	s.timeSlots = newSlots
	s.updatedAt = time.Now()
	return nil
}