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
	"intelligent-guidance-system/service/doctor/internal/domain/vo"

	"github.com/go-kratos/kratos/v2/log"
)

var (
	ErrDoctorAlreadyExists = errors.New("doctor already exists")
	ErrLicenseAlreadyUsed  = errors.New("license number already registered")
	ErrDoctorNotActive     = errors.New("doctor is not active")
)

type DoctorUseCase struct {
	doctorRepo  repository.DoctorRepository
	roleRepo    repository.RoleRepository
	assignmentRepo repository.DoctorRoleAssignmentRepository
	eventRepo   repository.DoctorEventRepository
	log         *log.Helper
}

func NewDoctorUseCase(
	doctorRepo repository.DoctorRepository,
	roleRepo repository.RoleRepository,
	assignmentRepo repository.DoctorRoleAssignmentRepository,
	eventRepo repository.DoctorEventRepository,
	logger log.Logger,
) *DoctorUseCase {
	return &DoctorUseCase{
		doctorRepo:     doctorRepo,
		roleRepo:       roleRepo,
		assignmentRepo: assignmentRepo,
		eventRepo:      eventRepo,
		log:            log.NewHelper(logger),
	}
}

func (uc *DoctorUseCase) CreateDoctor(ctx context.Context, req *dto.DoctorCreateRequest) (*dto.DoctorDTO, error) {
	exists, err := uc.doctorRepo.ExistsByEmployeeID(ctx, req.EmployeeID)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("check employee ID exists failed: %v", err)
		return nil, err
	}
	if exists {
		return nil, ErrDoctorAlreadyExists
	}

	exists, err = uc.doctorRepo.ExistsByLicenseNumber(ctx, req.LicenseNumber)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("check license number exists failed: %v", err)
		return nil, err
	}
	if exists {
		return nil, ErrLicenseAlreadyUsed
	}

	basicInfo, err := vo.NewDoctorBasicInfo(
		req.Name, req.Gender, req.BirthDate,
		req.Phone, req.Email, req.IDCard,
		req.Address, req.AvatarURL,
	)
	if err != nil {
		return nil, err
	}

	professional, err := vo.NewProfessionalInfo(
		req.LicenseNumber,
		vo.PositionFromCode(req.Position),
		vo.TitleFromCode(req.Title),
		req.Specialties,
		vo.EducationLevelFromCode(req.Education),
		req.GraduationSchool,
		req.GraduationYear,
		req.YearsOfExp,
		req.Certifications,
		req.Introduction,
	)
	if err != nil {
		return nil, err
	}

	doctor, err := aggregate.NewDoctor(req.EmployeeID, basicInfo, professional, req.DepartmentID)
	if err != nil {
		return nil, err
	}

	if err := uc.doctorRepo.Save(ctx, doctor); err != nil {
		uc.log.WithContext(ctx).Errorf("save doctor failed: %v", err)
		return nil, err
	}

	evt := event.NewDoctorCreatedEvent(
		doctor.ID(), doctor.EmployeeID(),
		doctor.FullName(), doctor.Position().String(),
		doctor.DepartmentID(),
	)
	uc.eventRepo.Save(ctx, evt)

	return uc.toDoctorDTO(doctor), nil
}

func (uc *DoctorUseCase) GetDoctor(ctx context.Context, doctorID string) (*dto.DoctorDetailDTO, error) {
	doctor, err := uc.doctorRepo.FindByID(ctx, doctorID)
	if err != nil {
		return nil, err
	}
	if doctor == nil {
		return nil, aggregate.ErrDoctorNotFound
	}

	detail := &dto.DoctorDetailDTO{
		DoctorDTO: uc.toDoctorDTO(doctor),
		Schedules: uc.toScheduleDTOs(doctor.Schedules()),
		Roles:     uc.toRoleAssignmentDTOs(doctor.Roles()),
	}
	return detail, nil
}

func (uc *DoctorUseCase) GetDoctorByEmployeeID(ctx context.Context, employeeID string) (*dto.DoctorDetailDTO, error) {
	doctor, err := uc.doctorRepo.FindByEmployeeID(ctx, employeeID)
	if err != nil {
		return nil, err
	}
	if doctor == nil {
		return nil, aggregate.ErrDoctorNotFound
	}

	detail := &dto.DoctorDetailDTO{
		DoctorDTO: uc.toDoctorDTO(doctor),
		Schedules: uc.toScheduleDTOs(doctor.Schedules()),
		Roles:     uc.toRoleAssignmentDTOs(doctor.Roles()),
	}
	return detail, nil
}

func (uc *DoctorUseCase) UpdateDoctor(ctx context.Context, req *dto.DoctorUpdateRequest) (*dto.DoctorDTO, error) {
	doctor, err := uc.doctorRepo.FindByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if doctor == nil {
		return nil, aggregate.ErrDoctorNotFound
	}

	basicInfo := doctor.BasicInfo()
	if req.Name != "" || req.Phone != "" || req.Email != "" || req.Address != "" {
		name := req.Name
		if name == "" {
			name = basicInfo.Name()
		}
		phone := req.Phone
		if phone == "" {
			phone = basicInfo.Phone()
		}
		email := req.Email
		if email == "" {
			email = basicInfo.Email()
		}
		address := req.Address
		if address == "" {
			address = basicInfo.Address()
		}

		newBasic, err := vo.NewDoctorBasicInfo(
			name, req.Gender, req.BirthDate,
			phone, email, req.IDCard,
			address, req.AvatarURL,
		)
		if err != nil {
			return nil, err
		}
		if err := doctor.UpdateBasicInfo(newBasic); err != nil {
			return nil, err
		}
	}

	professional := doctor.Professional()
	if req.Position != "" || req.Title != "" || len(req.Specialties) > 0 {
		position := vo.PositionFromCode(req.Position)
		if req.Position == "" {
			position = professional.Position()
		}
		title := vo.TitleFromCode(req.Title)
		if req.Title == "" {
			title = professional.Title()
		}
		specialties := req.Specialties
		if len(specialties) == 0 {
			specialties = professional.Specialties()
		}
		education := vo.EducationLevelFromCode(req.Education)
		if req.Education == "" {
			education = professional.Education()
		}
		school := req.GraduationSchool
		if school == "" {
			school = professional.GraduationSchool()
		}
		year := req.GraduationYear
		if year == 0 {
			year = professional.GraduationYear()
		}
		years := req.YearsOfExp
		if years == 0 {
			years = professional.YearsOfExperience()
		}
		certs := req.Certifications
		if len(certs) == 0 {
			certs = professional.Certifications()
		}
		intro := req.Introduction
		if intro == "" {
			intro = professional.Introduction()
		}

		newProf, err := vo.NewProfessionalInfo(
			professional.LicenseNumber(),
			position, title, specialties,
			education, school, year, years, certs, intro,
		)
		if err != nil {
			return nil, err
		}
		if err := doctor.UpdateProfessional(newProf); err != nil {
			return nil, err
		}
	}

	if err := uc.doctorRepo.Update(ctx, doctor); err != nil {
		uc.log.WithContext(ctx).Errorf("update doctor failed: %v", err)
		return nil, err
	}

	evt := event.NewDoctorEvent(event.DoctorEventUpdated, doctor.ID(), doctor.EmployeeID(), nil)
	uc.eventRepo.Save(ctx, evt)

	return uc.toDoctorDTO(doctor), nil
}

func (uc *DoctorUseCase) ChangeDepartment(ctx context.Context, req *dto.DoctorDepartmentChangeRequest) error {
	doctor, err := uc.doctorRepo.FindByID(ctx, req.DoctorID)
	if err != nil {
		return err
	}
	if doctor == nil {
		return aggregate.ErrDoctorNotFound
	}

	oldDeptID := doctor.DepartmentID()
	if err := doctor.ChangeDepartment(req.NewDeptID); err != nil {
		return err
	}

	if err := uc.doctorRepo.Update(ctx, doctor); err != nil {
		uc.log.WithContext(ctx).Errorf("update doctor department failed: %v", err)
		return err
	}

	evt := event.NewDoctorDepartmentChangedEvent(
		doctor.ID(), doctor.EmployeeID(),
		oldDeptID, req.NewDeptID, req.Reason,
	)
	uc.eventRepo.Save(ctx, evt)

	return nil
}

func (uc *DoctorUseCase) ChangeStatus(ctx context.Context, req *dto.DoctorStatusChangeRequest) error {
	doctor, err := uc.doctorRepo.FindByID(ctx, req.DoctorID)
	if err != nil {
		return err
	}
	if doctor == nil {
		return aggregate.ErrDoctorNotFound
	}

	oldStatus := doctor.Status().String()
	switch req.Status {
	case "ACTIVE":
		doctor.Activate()
	case "INACTIVE":
		doctor.Deactivate()
	case "ON_LEAVE":
		doctor.SetOnLeave()
	case "SUSPENDED":
		doctor.Suspend()
	case "RETIRED":
		doctor.Retire()
	default:
		return errors.New("invalid status")
	}

	if err := uc.doctorRepo.Update(ctx, doctor); err != nil {
		uc.log.WithContext(ctx).Errorf("update doctor status failed: %v", err)
		return err
	}

	evt := event.NewDoctorStatusChangedEvent(
		doctor.ID(), doctor.EmployeeID(),
		oldStatus, req.Status, req.Reason,
	)
	uc.eventRepo.Save(ctx, evt)

	return nil
}

func (uc *DoctorUseCase) AssignRole(ctx context.Context, req *dto.DoctorRoleAssignRequest) error {
	doctor, err := uc.doctorRepo.FindByID(ctx, req.DoctorID)
	if err != nil {
		return err
	}
	if doctor == nil {
		return aggregate.ErrDoctorNotFound
	}

	role, err := uc.roleRepo.FindByID(ctx, req.RoleID)
	if err != nil {
		return err
	}
	if role == nil {
		return entity.ErrRoleNotFound
	}

	deptID := req.DepartmentID
	if deptID == "" {
		deptID = doctor.DepartmentID()
	}

	if err := doctor.AssignRole(req.RoleID, deptID, req.AssignedBy, req.ExpiresAt); err != nil {
		return err
	}

	if err := uc.doctorRepo.Update(ctx, doctor); err != nil {
		uc.log.WithContext(ctx).Errorf("update doctor roles failed: %v", err)
		return err
	}

	evt := event.NewDoctorRoleAssignedEvent(
		doctor.ID(), doctor.EmployeeID(),
		req.RoleID, role.Name(), req.AssignedBy, req.ExpiresAt,
	)
	uc.eventRepo.Save(ctx, evt)

	return nil
}

func (uc *DoctorUseCase) RemoveRole(ctx context.Context, req *dto.DoctorRoleRemoveRequest) error {
	doctor, err := uc.doctorRepo.FindByID(ctx, req.DoctorID)
	if err != nil {
		return err
	}
	if doctor == nil {
		return aggregate.ErrDoctorNotFound
	}

	if err := doctor.RemoveRole(req.RoleID); err != nil {
		return err
	}

	if err := uc.doctorRepo.Update(ctx, doctor); err != nil {
		uc.log.WithContext(ctx).Errorf("update doctor roles failed: %v", err)
		return err
	}

	evt := event.NewDoctorEvent(event.DoctorEventRoleRemoved, doctor.ID(), doctor.EmployeeID(), map[string]interface{}{
		"role_id": req.RoleID,
	})
	uc.eventRepo.Save(ctx, evt)

	return nil
}

func (uc *DoctorUseCase) ListDoctors(ctx context.Context, req *dto.DoctorFilterRequest) (*dto.DoctorListResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}

	filter := &repository.DoctorFilter{
		DepartmentID: req.DepartmentID,
		Status:       req.Status,
		Position:     req.Position,
		Specialty:    req.Specialty,
		Name:         req.Name,
		EmployeeID:   req.EmployeeID,
		IsExpert:     req.IsExpert,
	}

	doctors, total, err := uc.doctorRepo.List(ctx, filter, req.Page, req.PageSize)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("list doctors failed: %v", err)
		return nil, err
	}

	items := make([]*dto.DoctorDTO, 0, len(doctors))
	for _, d := range doctors {
		items = append(items, uc.toDoctorDTO(d))
	}

	return &dto.DoctorListResponse{
		Items:    items,
		Total:    total,
		Page:     req.Page,
	PageSize: req.PageSize,
	}, nil
}

func (uc *DoctorUseCase) ListDoctorsByDepartment(ctx context.Context, departmentID string) ([]*dto.DoctorDTO, error) {
	doctors, err := uc.doctorRepo.FindByDepartment(ctx, departmentID)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("find doctors by department failed: %v", err)
		return nil, err
	}

	items := make([]*dto.DoctorDTO, 0, len(doctors))
	for _, d := range doctors {
		items = append(items, uc.toDoctorDTO(d))
	}
	return items, nil
}

func (uc *DoctorUseCase) ListExpertDoctors(ctx context.Context) ([]*dto.DoctorDTO, error) {
	doctors, err := uc.doctorRepo.FindExperts(ctx)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("find expert doctors failed: %v", err)
		return nil, err
	}

	items := make([]*dto.DoctorDTO, 0, len(doctors))
	for _, d := range doctors {
		items = append(items, uc.toDoctorDTO(d))
	}
	return items, nil
}

func (uc *DoctorUseCase) ListDoctorsBySpecialty(ctx context.Context, specialty string) ([]*dto.DoctorDTO, error) {
	doctors, err := uc.doctorRepo.FindBySpecialty(ctx, specialty)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("find doctors by specialty failed: %v", err)
		return nil, err
	}

	items := make([]*dto.DoctorDTO, 0, len(doctors))
	for _, d := range doctors {
		items = append(items, uc.toDoctorDTO(d))
	}
	return items, nil
}

func (uc *DoctorUseCase) DeleteDoctor(ctx context.Context, doctorID string) error {
	doctor, err := uc.doctorRepo.FindByID(ctx, doctorID)
	if err != nil {
		return err
	}
	if doctor == nil {
		return aggregate.ErrDoctorNotFound
	}

	return uc.doctorRepo.Delete(ctx, doctorID)
}

func (uc *DoctorUseCase) toDoctorDTO(d *aggregate.Doctor) *dto.DoctorDTO {
	basicInfo := &dto.DoctorBasicInfoDTO{
		Name:      d.BasicInfo().Name(),
		Gender:    d.BasicInfo().Gender(),
		BirthDate: d.BasicInfo().BirthDate(),
		Age:       d.BasicInfo().Age(),
		Phone:     d.BasicInfo().Phone(),
		Email:     d.BasicInfo().Email(),
		IDCard:    d.BasicInfo().IDCard(),
		Address:   d.BasicInfo().Address(),
		AvatarURL: d.BasicInfo().AvatarURL(),
	}

	professional := &dto.ProfessionalInfoDTO{
		LicenseNumber:     d.Professional().LicenseNumber(),
		Position:          d.Professional().Position().Code(),
		PositionDisplay:   d.Professional().Position().String(),
		Title:             d.Professional().Title().Code(),
		TitleDisplay:      d.Professional().Title().String(),
		Specialties:       d.Professional().Specialties(),
		Education:         d.Professional().Education().Code(),
		EducationDisplay:  d.Professional().Education().String(),
		GraduationSchool:  d.Professional().GraduationSchool(),
		GraduationYear:    d.Professional().GraduationYear(),
		YearsOfExperience: d.Professional().YearsOfExperience(),
		Certifications:    d.Professional().Certifications(),
		Introduction:      d.Professional().Introduction(),
		IsExpert:          d.IsExpert(),
		CanTreat:          d.CanTreat(),
	}

	return &dto.DoctorDTO{
		ID:            d.ID(),
		EmployeeID:    d.EmployeeID(),
		BasicInfo:     basicInfo,
		Professional:  professional,
		DepartmentID:  d.DepartmentID(),
		Status:        d.Status().Code(),
		StatusDisplay: d.Status().String(),
		IsExpert:      d.IsExpert(),
		CanTreat:      d.CanTreat(),
		CreatedAt:     d.CreatedAt(),
		UpdatedAt:     d.UpdatedAt(),
	}
}

func (uc *DoctorUseCase) toScheduleDTOs(schedules []*entity.Schedule) []*dto.ScheduleDTO {
	result := make([]*dto.ScheduleDTO, 0, len(schedules))
	for _, s := range schedules {
		result = append(result, uc.toScheduleDTO(s))
	}
	return result
}

func (uc *DoctorUseCase) toScheduleDTO(s *entity.Schedule) *dto.ScheduleDTO {
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

func (uc *DoctorUseCase) toRoleAssignmentDTOs(assignments []*entity.DoctorRoleAssignment) []*dto.RoleAssignmentDTO {
	result := make([]*dto.RoleAssignmentDTO, 0, len(assignments))
	for _, a := range assignments {
		result = append(result, &dto.RoleAssignmentDTO{
			ID:           a.ID(),
			RoleID:       a.RoleID(),
			DepartmentID: a.DepartmentID(),
			AssignedBy:   a.AssignedBy(),
			AssignedAt:   a.AssignedAt(),
			ExpiresAt:    a.ExpiresAt(),
			IsActive:     a.IsActive(),
			IsExpired:    a.IsExpired(),
			IsValid:      a.IsValid(),
		})
	}
	return result
}