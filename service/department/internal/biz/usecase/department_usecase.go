package usecase

import (
	"context"
	"errors"

	"intelligent-guidance-system/service/department/internal/domain/aggregate"
	"intelligent-guidance-system/service/department/internal/domain/entity"
	"intelligent-guidance-system/service/department/internal/domain/repository"
)

var (
	ErrDepartmentNotFound = errors.New("department not found")
	ErrDepartmentDisabled = errors.New("department is disabled")
	ErrDoctorNotFound     = errors.New("doctor not found in department")
)

type DepartmentUsecase struct {
	repo repository.DepartmentRepository
}

func NewDepartmentUsecase(repo repository.DepartmentRepository) *DepartmentUsecase {
	return &DepartmentUsecase{repo: repo}
}

func (u *DepartmentUsecase) CreateDepartment(ctx context.Context,
	name string,
	introduction string,
	personInChargeID int64,
	phone string,
	address string,
	deptType entity.DepartmentType,
	registrationFee int64,
	expertRegistrationFee int64,
) (*aggregate.Department, error) {
	exists, err := u.repo.ExistsByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("department name already exists")
	}

	dept, err := aggregate.NewDepartment(name, introduction, personInChargeID, phone, address, deptType, registrationFee, expertRegistrationFee)
	if err != nil {
		return nil, err
	}

	if err := u.repo.Save(ctx, dept); err != nil {
		return nil, err
	}

	return dept, nil
}

func (u *DepartmentUsecase) GetDepartment(ctx context.Context, id int64) (*aggregate.Department, error) {
	dept, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if dept == nil {
		return nil, ErrDepartmentNotFound
	}
	return dept, nil
}

func (u *DepartmentUsecase) ListDepartments(ctx context.Context, page, pageSize int) ([]*aggregate.Department, int64, error) {
	return u.repo.FindAll(ctx, page, pageSize)
}

func (u *DepartmentUsecase) ListActiveDepartments(ctx context.Context) ([]*aggregate.Department, error) {
	return u.repo.FindActiveDepartments(ctx)
}

func (u *DepartmentUsecase) ListDepartmentsByType(ctx context.Context, deptType entity.DepartmentType) ([]*aggregate.Department, error) {
	return u.repo.FindByType(ctx, int(deptType))
}

func (u *DepartmentUsecase) AssignDoctor(ctx context.Context, departmentID, doctorID, operatorID int64) error {
	dept, err := u.repo.FindByID(ctx, departmentID)
	if err != nil {
		return err
	}
	if dept == nil {
		return ErrDepartmentNotFound
	}

	if err := dept.AssignDoctor(doctorID, operatorID); err != nil {
		return err
	}

	return u.repo.Save(ctx, dept)
}

func (u *DepartmentUsecase) RemoveDoctor(ctx context.Context, departmentID, doctorID, operatorID int64) error {
	dept, err := u.repo.FindByID(ctx, departmentID)
	if err != nil {
		return err
	}
	if dept == nil {
		return ErrDepartmentNotFound
	}

	if err := dept.RemoveDoctor(doctorID, operatorID); err != nil {
		return err
	}

	return u.repo.Save(ctx, dept)
}

func (u *DepartmentUsecase) SetDutyDoctor(ctx context.Context, departmentID, doctorID, operatorID int64) error {
	dept, err := u.repo.FindByID(ctx, departmentID)
	if err != nil {
		return err
	}
	if dept == nil {
		return ErrDepartmentNotFound
	}

	if err := dept.SetDutyDoctor(doctorID, operatorID); err != nil {
		return err
	}

	return u.repo.Save(ctx, dept)
}

func (u *DepartmentUsecase) ActivateDepartment(ctx context.Context, departmentID, operatorID int64) error {
	dept, err := u.repo.FindByID(ctx, departmentID)
	if err != nil {
		return err
	}
	if dept == nil {
		return ErrDepartmentNotFound
	}

	return dept.Activate(operatorID)
}

func (u *DepartmentUsecase) RestDepartment(ctx context.Context, departmentID, operatorID int64) error {
	dept, err := u.repo.FindByID(ctx, departmentID)
	if err != nil {
		return err
	}
	if dept == nil {
		return ErrDepartmentNotFound
	}

	return dept.Rest(operatorID)
}

func (u *DepartmentUsecase) DisableDepartment(ctx context.Context, departmentID, operatorID int64) error {
	dept, err := u.repo.FindByID(ctx, departmentID)
	if err != nil {
		return err
	}
	if dept == nil {
		return ErrDepartmentNotFound
	}

	if err := dept.Disable(operatorID); err != nil {
		return err
	}

	return u.repo.Save(ctx, dept)
}

func (u *DepartmentUsecase) UpdateDepartment(ctx context.Context,
	departmentID int64,
	introduction string,
	phone string,
	address string,
	registrationFee int64,
	expertRegistrationFee int64,
) error {
	dept, err := u.repo.FindByID(ctx, departmentID)
	if err != nil {
		return err
	}
	if dept == nil {
		return ErrDepartmentNotFound
	}

	dept.UpdateIntroduction(introduction)
	dept.UpdatePhone(phone)
	dept.UpdateAddress(address)
	dept.UpdateFees(registrationFee, expertRegistrationFee)

	return u.repo.Save(ctx, dept)
}

func (u *DepartmentUsecase) GetDepartmentsByDoctorID(ctx context.Context, doctorID int64) ([]*aggregate.Department, error) {
	return u.repo.FindByDoctorID(ctx, doctorID)
}

func (u *DepartmentUsecase) DeleteDepartment(ctx context.Context, departmentID int64) error {
	return u.repo.Delete(ctx, departmentID)
}