package service

import (
	"context"
	"time"

	"intelligent-guidance-system/service/doctor/internal/biz/dto"
	"intelligent-guidance-system/service/doctor/internal/biz/usecase"

	"github.com/go-kratos/kratos/v2/log"
)

type DoctorService struct {
	doctorUC   *usecase.DoctorUseCase
	scheduleUC *usecase.ScheduleUseCase
	log        *log.Helper
}

func NewDoctorService(doctorUC *usecase.DoctorUseCase, scheduleUC *usecase.ScheduleUseCase, logger log.Logger) *DoctorService {
	return &DoctorService{
		doctorUC:   doctorUC,
		scheduleUC: scheduleUC,
		log:        log.NewHelper(logger),
	}
}

func (s *DoctorService) CreateDoctor(ctx context.Context, req *dto.DoctorCreateRequest) (*dto.DoctorDTO, error) {
	return s.doctorUC.CreateDoctor(ctx, req)
}

func (s *DoctorService) GetDoctor(ctx context.Context, doctorID string) (*dto.DoctorDetailDTO, error) {
	return s.doctorUC.GetDoctor(ctx, doctorID)
}

func (s *DoctorService) GetDoctorByEmployeeID(ctx context.Context, employeeID string) (*dto.DoctorDetailDTO, error) {
	return s.doctorUC.GetDoctorByEmployeeID(ctx, employeeID)
}

func (s *DoctorService) UpdateDoctor(ctx context.Context, req *dto.DoctorUpdateRequest) (*dto.DoctorDTO, error) {
	return s.doctorUC.UpdateDoctor(ctx, req)
}

func (s *DoctorService) ChangeDepartment(ctx context.Context, req *dto.DoctorDepartmentChangeRequest) error {
	return s.doctorUC.ChangeDepartment(ctx, req)
}

func (s *DoctorService) ChangeStatus(ctx context.Context, req *dto.DoctorStatusChangeRequest) error {
	return s.doctorUC.ChangeStatus(ctx, req)
}

func (s *DoctorService) AssignRole(ctx context.Context, req *dto.DoctorRoleAssignRequest) error {
	return s.doctorUC.AssignRole(ctx, req)
}

func (s *DoctorService) RemoveRole(ctx context.Context, req *dto.DoctorRoleRemoveRequest) error {
	return s.doctorUC.RemoveRole(ctx, req)
}

func (s *DoctorService) ListDoctors(ctx context.Context, req *dto.DoctorFilterRequest) (*dto.DoctorListResponse, error) {
	return s.doctorUC.ListDoctors(ctx, req)
}

func (s *DoctorService) ListDoctorsByDepartment(ctx context.Context, departmentID string) ([]*dto.DoctorDTO, error) {
	return s.doctorUC.ListDoctorsByDepartment(ctx, departmentID)
}

func (s *DoctorService) ListExpertDoctors(ctx context.Context) ([]*dto.DoctorDTO, error) {
	return s.doctorUC.ListExpertDoctors(ctx)
}

func (s *DoctorService) ListDoctorsBySpecialty(ctx context.Context, specialty string) ([]*dto.DoctorDTO, error) {
	return s.doctorUC.ListDoctorsBySpecialty(ctx, specialty)
}

func (s *DoctorService) DeleteDoctor(ctx context.Context, doctorID string) error {
	return s.doctorUC.DeleteDoctor(ctx, doctorID)
}

func (s *DoctorService) SetSchedule(ctx context.Context, req *dto.ScheduleSetRequest) (*dto.ScheduleDTO, error) {
	return s.scheduleUC.SetSchedule(ctx, req)
}

func (s *DoctorService) AddScheduleSlot(ctx context.Context, req *dto.ScheduleAddSlotRequest) (*dto.ScheduleDTO, error) {
	return s.scheduleUC.AddScheduleSlot(ctx, req)
}

func (s *DoctorService) GetSchedule(ctx context.Context, doctorID string, date time.Time) (*dto.ScheduleDTO, error) {
	return s.scheduleUC.GetSchedule(ctx, doctorID, date)
}

func (s *DoctorService) GetSchedules(ctx context.Context, doctorID string, startDate, endDate time.Time) (*dto.ScheduleListResponse, error) {
	return s.scheduleUC.GetSchedules(ctx, doctorID, startDate, endDate)
}

func (s *DoctorService) PublishSchedule(ctx context.Context, req *dto.SchedulePublishRequest) (*dto.ScheduleDTO, error) {
	return s.scheduleUC.PublishSchedule(ctx, req)
}

func (s *DoctorService) CancelSchedule(ctx context.Context, req *dto.ScheduleCancelRequest) error {
	return s.scheduleUC.CancelSchedule(ctx, req)
}

func (s *DoctorService) BookSlot(ctx context.Context, scheduleID, slotID string) error {
	return s.scheduleUC.BookSlot(ctx, scheduleID, slotID)
}

func (s *DoctorService) CancelBooking(ctx context.Context, scheduleID, slotID string) error {
	return s.scheduleUC.CancelBooking(ctx, scheduleID, slotID)
}

func (s *DoctorService) GetAvailableSlots(ctx context.Context, doctorID string, date time.Time) (*dto.AvailableSlotsResponse, error) {
	return s.scheduleUC.GetAvailableSlots(ctx, doctorID, date)
}

func (s *DoctorService) GetPublishedSchedulesByDate(ctx context.Context, date time.Time) ([]*dto.ScheduleDTO, error) {
	return s.scheduleUC.GetPublishedSchedulesByDate(ctx, date)
}

func (s *DoctorService) GetSchedulesByDepartment(ctx context.Context, departmentID string, date time.Time) ([]*dto.ScheduleDTO, error) {
	return s.scheduleUC.GetSchedulesByDepartment(ctx, departmentID, date)
}

func (s *DoctorService) RemoveScheduleSlot(ctx context.Context, doctorID string, date time.Time, slotID string) error {
	return s.scheduleUC.RemoveScheduleSlot(ctx, doctorID, date, slotID)
}