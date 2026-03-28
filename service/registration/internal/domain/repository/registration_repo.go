package repository

import (
	"context"
	"time"

	"intelligent-guidance-system/service/registration/internal/domain/aggregate"
)

type RegistrationRepository interface {
	Save(ctx context.Context, registration *aggregate.Registration) error
	FindByID(ctx context.Context, id int64) (*aggregate.Registration, error)
	FindByPatientID(ctx context.Context, patientID int64) ([]*aggregate.Registration, error)
	FindByDoctorID(ctx context.Context, doctorID int64) ([]*aggregate.Registration, error)
	FindByDepartmentID(ctx context.Context, departmentID int64) ([]*aggregate.Registration, error)
	FindByStatus(ctx context.Context, status int) ([]*aggregate.Registration, error)
	FindByAppointmentTime(ctx context.Context, startTime, endTime time.Time) ([]*aggregate.Registration, error)
	FindAll(ctx context.Context, page, pageSize int) ([]*aggregate.Registration, int64, error)
	Delete(ctx context.Context, id int64) error
}