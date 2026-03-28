package aggregate

import (
	"errors"
	"time"

	"intelligent-guidance-system/service/doctor/internal/domain/entity"
	"intelligent-guidance-system/service/doctor/internal/domain/vo"

	"github.com/google/uuid"
)

var (
	ErrDoctorNotFound      = errors.New("doctor not found")
	ErrDoctorInactive      = errors.New("doctor is inactive")
	ErrInvalidDepartment   = errors.New("invalid department ID")
	ErrScheduleConflict    = errors.New("schedule conflict")
	ErrRoleAlreadyAssigned = errors.New("role already assigned")
	ErrRoleNotAssigned     = errors.New("role not assigned")
	ErrUnauthorizedAction  = errors.New("unauthorized action")
)

type DoctorStatus int

const (
	DoctorStatusUnknown DoctorStatus = iota
	DoctorStatusActive
	DoctorStatusInactive
	DoctorStatusOnLeave
	DoctorStatusSuspended
	DoctorStatusRetired
)

func (s DoctorStatus) String() string {
	switch s {
	case DoctorStatusActive:
		return "在职"
	case DoctorStatusInactive:
		return "未激活"
	case DoctorStatusOnLeave:
		return "休假"
	case DoctorStatusSuspended:
		return "停职"
	case DoctorStatusRetired:
		return "退休"
	default:
		return "未知"
	}
}

func DoctorStatusFromCode(code string) DoctorStatus {
	switch code {
	case "ACTIVE":
		return DoctorStatusActive
	case "INACTIVE":
		return DoctorStatusInactive
	case "ON_LEAVE":
		return DoctorStatusOnLeave
	case "SUSPENDED":
		return DoctorStatusSuspended
	case "RETIRED":
		return DoctorStatusRetired
	default:
		return DoctorStatusUnknown
	}
}

func (s DoctorStatus) Code() string {
	switch s {
	case DoctorStatusActive:
		return "ACTIVE"
	case DoctorStatusInactive:
		return "INACTIVE"
	case DoctorStatusOnLeave:
		return "ON_LEAVE"
	case DoctorStatusSuspended:
		return "SUSPENDED"
	case DoctorStatusRetired:
		return "RETIRED"
	default:
		return "UNKNOWN"
	}
}

func (s DoctorStatus) CanWork() bool {
	return s == DoctorStatusActive
}

type Doctor struct {
	id             string
	employeeID     string
	basicInfo      *vo.DoctorBasicInfo
	professional   *vo.ProfessionalInfo
	departmentID   string
	status         DoctorStatus
	schedules      []*entity.Schedule
	roles          []*entity.DoctorRoleAssignment
	createdAt      time.Time
	updatedAt      time.Time
	version        int
}

func NewDoctor(
	employeeID string,
	basicInfo *vo.DoctorBasicInfo,
	professional *vo.ProfessionalInfo,
	departmentID string,
) (*Doctor, error) {
	if employeeID == "" {
		return nil, errors.New("employee ID is required")
	}
	if basicInfo == nil {
		return nil, errors.New("basic info is required")
	}
	if professional == nil {
		return nil, errors.New("professional info is required")
	}
	if departmentID == "" {
		return nil, ErrInvalidDepartment
	}

	now := time.Now()
	return &Doctor{
		id:           uuid.New().String(),
		employeeID:   employeeID,
		basicInfo:    basicInfo,
		professional: professional,
		departmentID: departmentID,
		status:       DoctorStatusInactive,
		schedules:    make([]*entity.Schedule, 0),
		roles:        make([]*entity.DoctorRoleAssignment, 0),
		createdAt:    now,
		updatedAt:    now,
		version:      1,
	}, nil
}

func ReconstructDoctor(
	id, employeeID string,
	basicInfo *vo.DoctorBasicInfo,
	professional *vo.ProfessionalInfo,
	departmentID string,
	status DoctorStatus,
	schedules []*entity.Schedule,
	roles []*entity.DoctorRoleAssignment,
	createdAt, updatedAt time.Time,
	version int,
) *Doctor {
	return &Doctor{
		id:           id,
		employeeID:   employeeID,
		basicInfo:    basicInfo,
		professional: professional,
		departmentID: departmentID,
		status:       status,
		schedules:    schedules,
		roles:        roles,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
		version:      version,
	}
}

func (d *Doctor) ID() string                         { return d.id }
func (d *Doctor) EmployeeID() string                 { return d.employeeID }
func (d *Doctor) BasicInfo() *vo.DoctorBasicInfo     { return d.basicInfo }
func (d *Doctor) Professional() *vo.ProfessionalInfo { return d.professional }
func (d *Doctor) DepartmentID() string               { return d.departmentID }
func (d *Doctor) Status() DoctorStatus               { return d.status }
func (d *Doctor) Schedules() []*entity.Schedule       { return d.schedules }
func (d *Doctor) Roles() []*entity.DoctorRoleAssignment { return d.roles }
func (d *Doctor) CreatedAt() time.Time               { return d.createdAt }
func (d *Doctor) UpdatedAt() time.Time               { return d.updatedAt }
func (d *Doctor) Version() int                        { return d.version }

func (d *Doctor) Activate() error {
	d.status = DoctorStatusActive
	d.updatedAt = time.Now()
	d.version++
	return nil
}

func (d *Doctor) Deactivate() {
	d.status = DoctorStatusInactive
	d.updatedAt = time.Now()
	d.version++
}

func (d *Doctor) SetOnLeave() {
	d.status = DoctorStatusOnLeave
	d.updatedAt = time.Now()
	d.version++
}

func (d *Doctor) Suspend() {
	d.status = DoctorStatusSuspended
	d.updatedAt = time.Now()
	d.version++
}

func (d *Doctor) Retire() {
	d.status = DoctorStatusRetired
	d.updatedAt = time.Now()
	d.version++
}

func (d *Doctor) ChangeDepartment(newDeptID string) error {
	if newDeptID == "" {
		return ErrInvalidDepartment
	}
	d.departmentID = newDeptID
	d.updatedAt = time.Now()
	d.version++
	return nil
}

func (d *Doctor) UpdateBasicInfo(basicInfo *vo.DoctorBasicInfo) error {
	if basicInfo == nil {
		return errors.New("basic info cannot be nil")
	}
	d.basicInfo = basicInfo
	d.updatedAt = time.Now()
	d.version++
	return nil
}

func (d *Doctor) UpdateProfessional(professional *vo.ProfessionalInfo) error {
	if professional == nil {
		return errors.New("professional info cannot be nil")
	}
	d.professional = professional
	d.updatedAt = time.Now()
	d.version++
	return nil
}

func (d *Doctor) SetSchedule(date time.Time, timeSlots []struct {
	StartTime   time.Time
	EndTime     time.Time
	MaxPatients int
}) error {
	if !d.status.CanWork() {
		return ErrDoctorInactive
	}

	date = date.Truncate(24 * time.Hour)
	for _, schedule := range d.schedules {
		if schedule.Date().Equal(date) {
			return schedule.SetTimeSlots(timeSlots)
		}
	}

	schedule, err := entity.NewSchedule(d.id, date)
	if err != nil {
		return err
	}
	if err := schedule.SetTimeSlots(timeSlots); err != nil {
		return err
	}
	d.schedules = append(d.schedules, schedule)
	d.updatedAt = time.Now()
	d.version++
	return nil
}

func (d *Doctor) AddScheduleSlot(date time.Time, startTime, endTime time.Time, maxPatients int) error {
	if !d.status.CanWork() {
		return ErrDoctorInactive
	}

	date = date.Truncate(24 * time.Hour)
	for _, schedule := range d.schedules {
		if schedule.Date().Equal(date) {
			return schedule.AddTimeSlot(startTime, endTime, maxPatients)
		}
	}

	schedule, err := entity.NewSchedule(d.id, date)
	if err != nil {
		return err
	}
	if err := schedule.AddTimeSlot(startTime, endTime, maxPatients); err != nil {
		return err
	}
	d.schedules = append(d.schedules, schedule)
	d.updatedAt = time.Now()
	d.version++
	return nil
}

func (d *Doctor) GetSchedule(date time.Time) (*entity.Schedule, error) {
	date = date.Truncate(24 * time.Hour)
	for _, schedule := range d.schedules {
		if schedule.Date().Equal(date) {
			return schedule, nil
		}
	}
	return nil, entity.ErrScheduleNotFound
}

func (d *Doctor) GetActiveSchedules() []*entity.Schedule {
	result := make([]*entity.Schedule, 0)
	now := time.Now().Truncate(24 * time.Hour)
	for _, schedule := range d.schedules {
		if !schedule.IsExpired() && schedule.Status() == entity.ScheduleStatusPublished {
			result = append(result, schedule)
		}
	}
	return result
}

func (d *Doctor) PublishSchedule(date time.Time) error {
	schedule, err := d.GetSchedule(date)
	if err != nil {
		return err
	}
	schedule.Publish()
	d.updatedAt = time.Now()
	d.version++
	return nil
}

func (d *Doctor) CancelSchedule(date time.Time) error {
	schedule, err := d.GetSchedule(date)
	if err != nil {
		return err
	}
	schedule.Cancel()
	d.updatedAt = time.Now()
	d.version++
	return nil
}

func (d *Doctor) AssignRole(roleID, departmentID, assignedBy string, expiresAt *time.Time) error {
	if roleID == "" {
		return errors.New("role ID is required")
	}

	for _, r := range d.roles {
		if r.RoleID() == roleID && r.IsActive() {
			return ErrRoleAlreadyAssigned
		}
	}

	assignment := entity.NewDoctorRoleAssignment(d.id, roleID, departmentID, assignedBy, expiresAt)
	d.roles = append(d.roles, assignment)
	d.updatedAt = time.Now()
	d.version++
	return nil
}

func (d *Doctor) RemoveRole(roleID string) error {
	for _, r := range d.roles {
		if r.RoleID() == roleID && r.IsActive() {
			r.Deactivate()
			d.updatedAt = time.Now()
			d.version++
			return nil
		}
	}
	return ErrRoleNotAssigned
}

func (d *Doctor) GetActiveRoles() []*entity.DoctorRoleAssignment {
	result := make([]*entity.DoctorRoleAssignment, 0)
	for _, r := range d.roles {
		if r.IsValid() {
			result = append(result, r)
		}
	}
	return result
}

func (d *Doctor) HasRole(roleID string) bool {
	for _, r := range d.roles {
		if r.RoleID() == roleID && r.IsValid() {
			return true
		}
	}
	return false
}

func (d *Doctor) IsExpert() bool {
	return d.professional != nil && d.professional.IsExpert()
}

func (d *Doctor) CanTreat() bool {
	return d.status.CanWork() && d.professional != nil && d.professional.CanTreat()
}

func (d *Doctor) HasSpecialty(specialty string) bool {
	return d.professional != nil && d.professional.HasSpecialty(specialty)
}

func (d *Doctor) FullName() string {
	if d.basicInfo != nil {
		return d.basicInfo.Name()
	}
	return ""
}

func (d *Doctor) Position() vo.Position {
	if d.professional != nil {
		return d.professional.Position()
	}
	return vo.PositionUnknown
}

func (d *Doctor) Title() vo.Title {
	if d.professional != nil {
		return d.professional.Title()
	}
	return vo.TitleUnknown
}

func (d *Doctor) Specialties() []string {
	if d.professional != nil {
		return d.professional.Specialties()
	}
	return nil
}

func (d *Doctor) AvailableScheduleCapacity(date time.Time) (total, booked, available int) {
	schedule, err := d.GetSchedule(date)
	if err != nil {
		return 0, 0, 0
	}
	return schedule.TotalCapacity(), schedule.TotalBooked(), schedule.AvailableCapacity()
}