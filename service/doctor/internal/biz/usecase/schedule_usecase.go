package usecase

import (
	"context"
	"errors"
	"time"

	"intelligent-guidance-system/service/doctor/internal/biz/dto"
	"intelligent-guidance-system/service/doctor/internal/domain/aggregate"
	"intelligent-guidance-system/service/doctor/internal/domain/entity"
	"intelligent-guidance-system/service/doctor/internal/domain/event"
	"intelligent-guidance-system/service/doctor/internal/domain/repository"

	"github.com/go-kratos/kratos/v2/log"
)

var (
	ErrScheduleAlreadyExists = errors.New("schedule already exists for this date")
	ErrCannotSetSchedule     = errors.New("cannot set schedule for inactive doctor")
)

type ScheduleUseCase struct {
	doctorRepo    repository.DoctorRepository
	scheduleRepo  repository.ScheduleRepository
	eventRepo     repository.DoctorEventRepository
	log           *log.Helper
}

func NewScheduleUseCase(
	doctorRepo repository.DoctorRepository,
	scheduleRepo repository.ScheduleRepository,
	eventRepo repository.DoctorEventRepository,
	logger log.Logger,
) *ScheduleUseCase {
	return &ScheduleUseCase{
		doctorRepo:   doctorRepo,
		scheduleRepo: scheduleRepo,
		eventRepo:    eventRepo,
		log:          log.NewHelper(logger),
	}
}

func (uc *ScheduleUseCase) SetSchedule(ctx context.Context, req *dto.ScheduleSetRequest) (*dto.ScheduleDTO, error) {
	doctor, err := uc.doctorRepo.FindByID(ctx, req.DoctorID)
	if err != nil {
		return nil, err
	}
	if doctor == nil {
		return nil, aggregate.ErrDoctorNotFound
	}

	if !doctor.Status().CanWork() {
		return nil, ErrCannotSetSchedule
	}

	timeSlots := make([]struct {
		StartTime   time.Time
		EndTime     time.Time
		MaxPatients int
	}, 0, len(req.TimeSlots))
	for _, ts := range req.TimeSlots {
		timeSlots = append(timeSlots, struct {
			StartTime   time.Time
			EndTime     time.Time
			MaxPatients int
		}{
			StartTime:   ts.StartTime,
			EndTime:     ts.EndTime,
			MaxPatients: ts.MaxPatients,
		})
	}

	if err := doctor.SetSchedule(req.Date, timeSlots); err != nil {
		return nil, err
	}

	schedule, err := doctor.GetSchedule(req.Date)
	if err != nil {
		return nil, err
	}

	if err := uc.doctorRepo.Update(ctx, doctor); err != nil {
		uc.log.WithContext(ctx).Errorf("update doctor schedule failed: %v", err)
		return nil, err
	}

	evt := event.NewDoctorScheduleSetEvent(
		doctor.ID(), doctor.EmployeeID(),
		req.Date.Format("2006-01-02"),
		len(schedule.TimeSlots()),
		schedule.TotalCapacity(),
	)
	uc.eventRepo.Save(ctx, evt)

	return uc.toScheduleDTO(schedule), nil
}

func (uc *ScheduleUseCase) AddScheduleSlot(ctx context.Context, req *dto.ScheduleAddSlotRequest) (*dto.ScheduleDTO, error) {
	doctor, err := uc.doctorRepo.FindByID(ctx, req.DoctorID)
	if err != nil {
		return nil, err
	}
	if doctor == nil {
		return nil, aggregate.ErrDoctorNotFound
	}

	if !doctor.Status().CanWork() {
		return nil, ErrCannotSetSchedule
	}

	if err := doctor.AddScheduleSlot(req.Date, req.StartTime, req.EndTime, req.MaxPatients); err != nil {
		return nil, err
	}

	schedule, err := doctor.GetSchedule(req.Date)
	if err != nil {
		return nil, err
	}

	if err := uc.doctorRepo.Update(ctx, doctor); err != nil {
		uc.log.WithContext(ctx).Errorf("update doctor schedule failed: %v", err)
		return nil, err
	}

	return uc.toScheduleDTO(schedule), nil
}

func (uc *ScheduleUseCase) GetSchedule(ctx context.Context, doctorID string, date time.Time) (*dto.ScheduleDTO, error) {
	doctor, err := uc.doctorRepo.FindByID(ctx, doctorID)
	if err != nil {
		return nil, err
	}
	if doctor == nil {
		return nil, aggregate.ErrDoctorNotFound
	}

	schedule, err := doctor.GetSchedule(date)
	if err != nil {
		return nil, err
	}

	return uc.toScheduleDTO(schedule), nil
}

func (uc *ScheduleUseCase) GetSchedules(ctx context.Context, doctorID string, startDate, endDate time.Time) (*dto.ScheduleListResponse, error) {
	schedules, err := uc.scheduleRepo.FindByDoctor(ctx, doctorID, startDate, endDate)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("find schedules failed: %v", err)
		return nil, err
	}

	items := make([]*dto.ScheduleDTO, 0, len(schedules))
	for _, s := range schedules {
		items = append(items, uc.toScheduleDTO(s))
	}

	return &dto.ScheduleListResponse{
		Items: items,
		Total: int64(len(items)),
	}, nil
}

func (uc *ScheduleUseCase) PublishSchedule(ctx context.Context, req *dto.SchedulePublishRequest) (*dto.ScheduleDTO, error) {
	doctor, err := uc.doctorRepo.FindByID(ctx, req.DoctorID)
	if err != nil {
		return nil, err
	}
	if doctor == nil {
		return nil, aggregate.ErrDoctorNotFound
	}

	if err := doctor.PublishSchedule(req.Date); err != nil {
		return nil, err
	}

	schedule, err := doctor.GetSchedule(req.Date)
	if err != nil {
		return nil, err
	}

	if err := uc.doctorRepo.Update(ctx, doctor); err != nil {
		uc.log.WithContext(ctx).Errorf("publish schedule failed: %v", err)
		return nil, err
	}

	evt := event.NewDoctorSchedulePublishedEvent(
		doctor.ID(), doctor.EmployeeID(),
		req.Date.Format("2006-01-02"),
		schedule.AvailableCapacity(),
	)
	uc.eventRepo.Save(ctx, evt)

	return uc.toScheduleDTO(schedule), nil
}

func (uc *ScheduleUseCase) CancelSchedule(ctx context.Context, req *dto.ScheduleCancelRequest) error {
	doctor, err := uc.doctorRepo.FindByID(ctx, req.DoctorID)
	if err != nil {
		return err
	}
	if doctor == nil {
		return aggregate.ErrDoctorNotFound
	}

	if err := doctor.CancelSchedule(req.Date); err != nil {
		return err
	}

	if err := uc.doctorRepo.Update(ctx, doctor); err != nil {
		uc.log.WithContext(ctx).Errorf("cancel schedule failed: %v", err)
		return err
	}

	evt := event.NewDoctorEvent(event.DoctorEventScheduleCancelled, doctor.ID(), doctor.EmployeeID(), map[string]interface{}{
		"date": req.Date.Format("2006-01-02"),
	})
	uc.eventRepo.Save(ctx, evt)

	return nil
}

func (uc *ScheduleUseCase) BookSlot(ctx context.Context, scheduleID, slotID string) error {
	schedule, err := uc.scheduleRepo.FindByID(ctx, scheduleID)
	if err != nil {
		return err
	}
	if schedule == nil {
		return entity.ErrScheduleNotFound
	}

	if schedule.Status() != entity.ScheduleStatusPublished {
		return errors.New("schedule is not published")
	}

	if err := schedule.BookSlot(slotID); err != nil {
		return err
	}

	return uc.scheduleRepo.Update(ctx, schedule)
}

func (uc *ScheduleUseCase) CancelBooking(ctx context.Context, scheduleID, slotID string) error {
	schedule, err := uc.scheduleRepo.FindByID(ctx, scheduleID)
	if err != nil {
		return err
	}
	if schedule == nil {
		return entity.ErrScheduleNotFound
	}

	if err := schedule.CancelSlotBooking(slotID); err != nil {
		return err
	}

	return uc.scheduleRepo.Update(ctx, schedule)
}

func (uc *ScheduleUseCase) GetAvailableSlots(ctx context.Context, doctorID string, date time.Time) (*dto.AvailableSlotsResponse, error) {
	slots, err := uc.scheduleRepo.ListAvailableSlots(ctx, doctorID, date)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("find available slots failed: %v", err)
		return nil, err
	}

	slotDTOs := make([]*dto.TimeSlotDTO, 0, len(slots))
	for _, s := range slots {
		slotDTOs = append(slotDTOs, &dto.TimeSlotDTO{
			ID:            s.ID(),
			StartTime:     s.StartTime(),
			EndTime:       s.EndTime(),
			IsAvailable:   s.IsAvailable(),
			MaxPatients:   s.MaxPatients(),
			BookedCount:   s.BookedCount(),
			AvailableSlots: s.AvailableSlots(),
			Duration:      int(s.Duration().Minutes()),
		})
	}

	return &dto.AvailableSlotsResponse{
		DoctorID: doctorID,
		Date:     date,
		Slots:    slotDTOs,
	}, nil
}

func (uc *ScheduleUseCase) GetPublishedSchedulesByDate(ctx context.Context, date time.Time) ([]*dto.ScheduleDTO, error) {
	schedules, err := uc.scheduleRepo.FindPublishedByDate(ctx, date)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("find published schedules failed: %v", err)
		return nil, err
	}

	items := make([]*dto.ScheduleDTO, 0, len(schedules))
	for _, s := range schedules {
		items = append(items, uc.toScheduleDTO(s))
	}
	return items, nil
}

func (uc *ScheduleUseCase) GetSchedulesByDepartment(ctx context.Context, departmentID string, date time.Time) ([]*dto.ScheduleDTO, error) {
	schedules, err := uc.scheduleRepo.FindByDepartmentAndDate(ctx, departmentID, date)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("find department schedules failed: %v", err)
		return nil, err
	}

	items := make([]*dto.ScheduleDTO, 0, len(schedules))
	for _, s := range schedules {
		items = append(items, uc.toScheduleDTO(s))
	}
	return items, nil
}

func (uc *ScheduleUseCase) RemoveScheduleSlot(ctx context.Context, doctorID string, date time.Time, slotID string) error {
	doctor, err := uc.doctorRepo.FindByID(ctx, doctorID)
	if err != nil {
		return err
	}
	if doctor == nil {
		return aggregate.ErrDoctorNotFound
	}

	schedule, err := doctor.GetSchedule(date)
	if err != nil {
		return err
	}

	if err := schedule.RemoveTimeSlot(slotID); err != nil {
		return err
	}

	return uc.doctorRepo.Update(ctx, doctor)
}

func (uc *ScheduleUseCase) toScheduleDTO(s *entity.Schedule) *dto.ScheduleDTO {
	slots := make([]*dto.TimeSlotDTO, 0, len(s.TimeSlots()))
	for _, ts := range s.TimeSlots() {
		slots = append(slots, &dto.TimeSlotDTO{
			ID:            ts.ID(),
			StartTime:     ts.StartTime(),
			EndTime:       ts.EndTime(),
			IsAvailable:   ts.IsAvailable(),
			MaxPatients:   ts.MaxPatients(),
			BookedCount:   ts.BookedCount(),
			AvailableSlots: ts.AvailableSlots(),
			Duration:      int(ts.Duration().Minutes()),
		})
	}

	statusCode := "DRAFT"
	switch s.Status() {
	case entity.ScheduleStatusDraft:
		statusCode = "DRAFT"
	case entity.ScheduleStatusPublished:
		statusCode = "PUBLISHED"
	case entity.ScheduleStatusCancelled:
		statusCode = "CANCELLED"
	}

	return &dto.ScheduleDTO{
		ID:              s.ID(),
		DoctorID:        s.DoctorID(),
		Date:            s.Date(),
		TimeSlots:       slots,
		Status:          statusCode,
		StatusDisplay:   s.Status().String(),
		TotalCapacity:   s.TotalCapacity(),
		TotalBooked:     s.TotalBooked(),
		AvailableCapacity: s.AvailableCapacity(),
		CreatedAt:       s.CreatedAt(),
		UpdatedAt:       s.UpdatedAt(),
	}
}