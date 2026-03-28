package repository

import (
	"context"
	"time"

	"intelligent-guidance-system/service/doctor/internal/domain/aggregate"
	"intelligent-guidance-system/service/doctor/internal/domain/entity"
)

type DoctorRepository interface {
	Save(ctx context.Context, doctor *aggregate.Doctor) error
	Update(ctx context.Context, doctor *aggregate.Doctor) error
	Delete(ctx context.Context, doctorID string) error
	FindByID(ctx context.Context, doctorID string) (*aggregate.Doctor, error)
	FindByEmployeeID(ctx context.Context, employeeID string) (*aggregate.Doctor, error)
	FindByDepartment(ctx context.Context, departmentID string) ([]*aggregate.Doctor, error)
	FindActive(ctx context.Context) ([]*aggregate.Doctor, error)
	FindExperts(ctx context.Context) ([]*aggregate.Doctor, error)
	FindBySpecialty(ctx context.Context, specialty string) ([]*aggregate.Doctor, error)
	FindByPosition(ctx context.Context, positionCode string) ([]*aggregate.Doctor, error)
	List(ctx context.Context, filters *DoctorFilter, page, pageSize int) ([]*aggregate.Doctor, int64, error)
	ExistsByEmployeeID(ctx context.Context, employeeID string) (bool, error)
	ExistsByLicenseNumber(ctx context.Context, licenseNumber string) (bool, error)
}

type DoctorFilter struct {
	DepartmentID string
	Status       string
	Position     string
	Specialty    string
	Name         string
	EmployeeID   string
	IsExpert     bool
}

type ScheduleRepository interface {
	Save(ctx context.Context, schedule *entity.Schedule) error
	Update(ctx context.Context, schedule *entity.Schedule) error
	Delete(ctx context.Context, scheduleID string) error
	FindByID(ctx context.Context, scheduleID string) (*entity.Schedule, error)
	FindByDoctorAndDate(ctx context.Context, doctorID string, date time.Time) (*entity.Schedule, error)
	FindByDoctor(ctx context.Context, doctorID string, startDate, endDate time.Time) ([]*entity.Schedule, error)
	FindPublishedByDate(ctx context.Context, date time.Time) ([]*entity.Schedule, error)
	FindByDepartmentAndDate(ctx context.Context, departmentID string, date time.Time) ([]*entity.Schedule, error)
	BookSlot(ctx context.Context, scheduleID, slotID string) error
	CancelSlotBooking(ctx context.Context, scheduleID, slotID string) error
	ListAvailableSlots(ctx context.Context, doctorID string, date time.Time) ([]*entity.TimeSlot, error)
}

type RoleRepository interface {
	Save(ctx context.Context, role *entity.Role) error
	Update(ctx context.Context, role *entity.Role) error
	Delete(ctx context.Context, roleID string) error
	FindByID(ctx context.Context, roleID string) (*entity.Role, error)
	FindByCode(ctx context.Context, code string) (*entity.Role, error)
	FindByType(ctx context.Context, roleType entity.RoleType) ([]*entity.Role, error)
	FindAll(ctx context.Context) ([]*entity.Role, error)
	List(ctx context.Context, page, pageSize int) ([]*entity.Role, int64, error)
}

type DoctorRoleAssignmentRepository interface {
	Save(ctx context.Context, assignment *entity.DoctorRoleAssignment) error
	Update(ctx context.Context, assignment *entity.DoctorRoleAssignment) error
	Delete(ctx context.Context, assignmentID string) error
	FindByID(ctx context.Context, assignmentID string) (*entity.DoctorRoleAssignment, error)
	FindByDoctor(ctx context.Context, doctorID string) ([]*entity.DoctorRoleAssignment, error)
	FindActiveByDoctor(ctx context.Context, doctorID string) ([]*entity.DoctorRoleAssignment, error)
	FindByRole(ctx context.Context, roleID string) ([]*entity.DoctorRoleAssignment, error)
	FindByDoctorAndRole(ctx context.Context, doctorID, roleID string) (*entity.DoctorRoleAssignment, error)
	Deactivate(ctx context.Context, assignmentID string) error
}

type DoctorEventRepository interface {
	Save(ctx context.Context, event interface{}) error
	FindByDoctor(ctx context.Context, doctorID string, limit int) ([]interface{}, error)
	FindByType(ctx context.Context, eventType string, limit int) ([]interface{}, error)
	FindRecent(ctx context.Context, limit int) ([]interface{}, error)
}