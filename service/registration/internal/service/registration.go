package service

import (
	"context"

	"intelligent-guidance-system/service/registration/internal/biz/dto"
	"intelligent-guidance-system/service/registration/internal/biz/usecase"
	"intelligent-guidance-system/service/registration/internal/domain/aggregate"
	"intelligent-guidance-system/service/registration/internal/domain/entity"
)

type RegistrationService struct {
	usecase *usecase.RegistrationUsecase
}

func NewRegistrationService(usecase *usecase.RegistrationUsecase) *RegistrationService {
	return &RegistrationService{usecase: usecase}
}

func (s *RegistrationService) CreateRegistration(ctx context.Context, req *dto.CreateRegistrationRequest, feeAmount int64) (*dto.RegistrationResponse, error) {
	regType := entity.RegistrationTypeFromCode(req.Type)
	
	reg, err := s.usecase.CreateRegistration(ctx,
		req.PatientID,
		req.DoctorID,
		req.DepartmentID,
		regType,
		req.AppointmentTime,
		feeAmount,
	)
	if err != nil {
		return nil, err
	}

	return s.toResponse(reg), nil
}

func (s *RegistrationService) GetRegistration(ctx context.Context, id int64) (*dto.RegistrationResponse, error) {
	reg, err := s.usecase.GetRegistration(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(reg), nil
}

func (s *RegistrationService) CancelRegistration(ctx context.Context, id int64, req *dto.CancelRegistrationRequest) error {
	return s.usecase.CancelRegistration(ctx, id, req.OperatorID, req.Reason)
}

func (s *RegistrationService) ConfirmRegistration(ctx context.Context, id int64, req *dto.ConfirmRegistrationRequest) error {
	return s.usecase.ConfirmRegistration(ctx, id, req.OperatorID)
}

func (s *RegistrationService) CompleteRegistration(ctx context.Context, id int64, req *dto.CompleteRegistrationRequest) error {
	return s.usecase.CompleteRegistration(ctx, id, req.OperatorID)
}

func (s *RegistrationService) GetPatientRegistrations(ctx context.Context, patientID int64) ([]*dto.RegistrationResponse, error) {
	regs, err := s.usecase.GetPatientRegistrations(ctx, patientID)
	if err != nil {
		return nil, err
	}
	return s.toResponses(regs), nil
}

func (s *RegistrationService) GetDoctorRegistrations(ctx context.Context, doctorID int64) ([]*dto.RegistrationResponse, error) {
	regs, err := s.usecase.GetDoctorRegistrations(ctx, doctorID)
	if err != nil {
		return nil, err
	}
	return s.toResponses(regs), nil
}

func (s *RegistrationService) ListRegistrations(ctx context.Context, page, pageSize int) (*dto.RegistrationListResponse, error) {
	regs, total, err := s.usecase.ListRegistrations(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	return &dto.RegistrationListResponse{
		Total:        total,
		Page:         page,
		PageSize:     pageSize,
		Registrations: s.toResponseList(regs),
	}, nil
}

func (s *RegistrationService) DeleteRegistration(ctx context.Context, id int64) error {
	return s.usecase.DeleteRegistration(ctx, id)
}

func (s *RegistrationService) toResponse(reg *aggregate.Registration) *dto.RegistrationResponse {
	return &dto.RegistrationResponse{
		ID:              reg.ID(),
		PatientID:       reg.PatientID(),
		DoctorID:        reg.DoctorID(),
		DepartmentID:    reg.DepartmentID(),
		Type:            reg.Type().Code(),
		TypeName:        reg.Type().String(),
		Status:          reg.Status().Code(),
		StatusName:      reg.Status().String(),
		AppointmentTime: reg.AppointmentTime(),
		Fee:             reg.Fee().AmountYuan(),
		CreatedAt:       reg.CreatedAt(),
		UpdatedAt:       reg.UpdatedAt(),
	}
}

func (s *RegistrationService) toResponses(regs []*aggregate.Registration) []*dto.RegistrationResponse {
	responses := make([]*dto.RegistrationResponse, 0, len(regs))
	for _, reg := range regs {
		responses = append(responses, s.toResponse(reg))
	}
	return responses
}

func (s *RegistrationService) toResponseList(regs []*aggregate.Registration) []dto.RegistrationResponse {
	responses := make([]dto.RegistrationResponse, 0, len(regs))
	for _, reg := range regs {
		responses = append(responses, *s.toResponse(reg))
	}
	return responses
}