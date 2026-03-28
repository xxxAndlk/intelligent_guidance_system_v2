package usecase

import (
	"context"
	"errors"
	"time"

	"intelligent-guidance-system/service/registration/internal/domain/aggregate"
	"intelligent-guidance-system/service/registration/internal/domain/entity"
	"intelligent-guidance-system/service/registration/internal/domain/repository"
)

var (
	ErrRegistrationNotFound = errors.New("registration not found")
)

type RegistrationUsecase struct {
	repo repository.RegistrationRepository
}

func NewRegistrationUsecase(repo repository.RegistrationRepository) *RegistrationUsecase {
	return &RegistrationUsecase{repo: repo}
}

func (u *RegistrationUsecase) CreateRegistration(ctx context.Context,
	patientID int64,
	doctorID int64,
	departmentID int64,
	regType entity.RegistrationType,
	appointmentTime time.Time,
	feeAmount int64,
) (*aggregate.Registration, error) {
	reg, err := aggregate.NewRegistration(patientID, doctorID, departmentID, regType, appointmentTime, feeAmount)
	if err != nil {
		return nil, err
	}

	if err := u.repo.Save(ctx, reg); err != nil {
		return nil, err
	}

	return reg, nil
}

func (u *RegistrationUsecase) GetRegistration(ctx context.Context, id int64) (*aggregate.Registration, error) {
	reg, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if reg == nil {
		return nil, ErrRegistrationNotFound
	}
	return reg, nil
}

func (u *RegistrationUsecase) CancelRegistration(ctx context.Context, id int64, operatorID int64, reason string) error {
	reg, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if reg == nil {
		return ErrRegistrationNotFound
	}

	if err := reg.Cancel(operatorID, reason); err != nil {
		return err
	}

	return u.repo.Save(ctx, reg)
}

func (u *RegistrationUsecase) ConfirmRegistration(ctx context.Context, id int64, operatorID int64) error {
	reg, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if reg == nil {
		return ErrRegistrationNotFound
	}

	if err := reg.Confirm(operatorID); err != nil {
		return err
	}

	return u.repo.Save(ctx, reg)
}

func (u *RegistrationUsecase) CompleteRegistration(ctx context.Context, id int64, operatorID int64) error {
	reg, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if reg == nil {
		return ErrRegistrationNotFound
	}

	if err := reg.Complete(operatorID); err != nil {
		return err
	}

	return u.repo.Save(ctx, reg)
}

func (u *RegistrationUsecase) GetPatientRegistrations(ctx context.Context, patientID int64) ([]*aggregate.Registration, error) {
	return u.repo.FindByPatientID(ctx, patientID)
}

func (u *RegistrationUsecase) GetDoctorRegistrations(ctx context.Context, doctorID int64) ([]*aggregate.Registration, error) {
	return u.repo.FindByDoctorID(ctx, doctorID)
}

func (u *RegistrationUsecase) GetDepartmentRegistrations(ctx context.Context, departmentID int64) ([]*aggregate.Registration, error) {
	return u.repo.FindByDepartmentID(ctx, departmentID)
}

func (u *RegistrationUsecase) GetRegistrationsByStatus(ctx context.Context, status entity.RegistrationStatus) ([]*aggregate.Registration, error) {
	return u.repo.FindByStatus(ctx, int(status))
}

func (u *RegistrationUsecase) GetRegistrationsByTimeRange(ctx context.Context, startTime, endTime time.Time) ([]*aggregate.Registration, error) {
	return u.repo.FindByAppointmentTime(ctx, startTime, endTime)
}

func (u *RegistrationUsecase) ListRegistrations(ctx context.Context, page, pageSize int) ([]*aggregate.Registration, int64, error) {
	return u.repo.FindAll(ctx, page, pageSize)
}

func (u *RegistrationUsecase) DeleteRegistration(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}