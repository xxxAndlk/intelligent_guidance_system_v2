package service

import (
	"context"

	"intelligent-guidance-system/service/department/internal/biz/dto"
	"intelligent-guidance-system/service/department/internal/biz/usecase"
	"intelligent-guidance-system/service/department/internal/domain/aggregate"
	"intelligent-guidance-system/service/department/internal/domain/entity"
)

type DepartmentService struct {
	usecase *usecase.DepartmentUsecase
}

func NewDepartmentService(usecase *usecase.DepartmentUsecase) *DepartmentService {
	return &DepartmentService{usecase: usecase}
}

func (s *DepartmentService) CreateDepartment(ctx context.Context, req *dto.CreateDepartmentRequest) (*dto.DepartmentResponse, error) {
	deptType := entity.DepartmentTypeFromCode(req.Type)
	
	dept, err := s.usecase.CreateDepartment(ctx,
		req.Name,
		req.Introduction,
		req.PersonInChargeID,
		req.Phone,
		req.Address,
		deptType,
		req.RegistrationFee,
		req.ExpertRegistrationFee,
	)
	if err != nil {
		return nil, err
	}

	return s.toResponse(dept), nil
}

func (s *DepartmentService) GetDepartment(ctx context.Context, id int64) (*dto.DepartmentResponse, error) {
	dept, err := s.usecase.GetDepartment(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(dept), nil
}

func (s *DepartmentService) ListDepartments(ctx context.Context, page, pageSize int) (*dto.DepartmentListResponse, error) {
	departments, total, err := s.usecase.ListDepartments(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.DepartmentResponse, 0, len(departments))
	for _, dept := range departments {
		responses = append(responses, *s.toResponse(dept))
	}

	return &dto.DepartmentListResponse{
		Total:       total,
		Page:        page,
		PageSize:    pageSize,
		Departments: responses,
	}, nil
}

func (s *DepartmentService) ListActiveDepartments(ctx context.Context) ([]*dto.DepartmentResponse, error) {
	departments, err := s.usecase.ListActiveDepartments(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.DepartmentResponse, 0, len(departments))
	for _, dept := range departments {
		responses = append(responses, s.toResponse(dept))
	}
	return responses, nil
}

func (s *DepartmentService) AssignDoctor(ctx context.Context, departmentID int64, req *dto.AssignDoctorRequest) error {
	return s.usecase.AssignDoctor(ctx, departmentID, req.DoctorID, req.OperatorID)
}

func (s *DepartmentService) RemoveDoctor(ctx context.Context, departmentID int64, req *dto.RemoveDoctorRequest) error {
	return s.usecase.RemoveDoctor(ctx, departmentID, req.DoctorID, req.OperatorID)
}

func (s *DepartmentService) SetDutyDoctor(ctx context.Context, departmentID int64, req *dto.SetDutyDoctorRequest) error {
	return s.usecase.SetDutyDoctor(ctx, departmentID, req.DoctorID, req.OperatorID)
}

func (s *DepartmentService) ChangeStatus(ctx context.Context, departmentID int64, req *dto.ChangeStatusRequest) error {
	switch req.Status {
	case "ACTIVE":
		return s.usecase.ActivateDepartment(ctx, departmentID, req.OperatorID)
	case "RESTING":
		return s.usecase.RestDepartment(ctx, departmentID, req.OperatorID)
	case "DISABLED":
		return s.usecase.DisableDepartment(ctx, departmentID, req.OperatorID)
	default:
		return entity.ErrInvalidDepartmentType
	}
}

func (s *DepartmentService) UpdateDepartment(ctx context.Context, departmentID int64, req *dto.UpdateDepartmentRequest) error {
	return s.usecase.UpdateDepartment(ctx,
		departmentID,
		req.Introduction,
		req.Phone,
		req.Address,
		req.RegistrationFee,
		req.ExpertRegistrationFee,
	)
}

func (s *DepartmentService) DeleteDepartment(ctx context.Context, departmentID int64) error {
	return s.usecase.DeleteDepartment(ctx, departmentID)
}

func (s *DepartmentService) GetDepartmentsByDoctor(ctx context.Context, doctorID int64) ([]*dto.DepartmentResponse, error) {
	departments, err := s.usecase.GetDepartmentsByDoctorID(ctx, doctorID)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.DepartmentResponse, 0, len(departments))
	for _, dept := range departments {
		responses = append(responses, s.toResponse(dept))
	}
	return responses, nil
}

func (s *DepartmentService) toResponse(dept *aggregate.Department) *dto.DepartmentResponse {
	return &dto.DepartmentResponse{
		ID:                  dept.ID(),
		Name:                dept.Name(),
		Introduction:        dept.Introduction(),
		PersonInChargeID:    dept.PersonInChargeID(),
		Phone:               dept.Phone(),
		Address:             dept.Address(),
		Type:                dept.Type().Code(),
		TypeName:            dept.Type().String(),
		StaffCount:          dept.StaffCount(),
		Status:              dept.Status().Code(),
		StatusName:          dept.Status().String(),
		RegistrationFee:     dept.Fee().RegistrationFeeYuan(),
		ExpertRegistrationFee: dept.Fee().ExpertRegistrationFeeYuan(),
		DutyDoctorID:        dept.DutyDoctorID(),
		CreatedAt:           dept.CreatedAt(),
		UpdatedAt:           dept.UpdatedAt(),
	}
}