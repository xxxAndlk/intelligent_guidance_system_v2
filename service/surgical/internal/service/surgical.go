package service

import (
	"context"

	"intelligent-guidance-system/service/surgical/internal/biz/dto"
	"intelligent-guidance-system/service/surgical/internal/biz/usecase"
	"intelligent-guidance-system/service/surgical/internal/domain/aggregate"
	"intelligent-guidance-system/service/surgical/internal/domain/entity"
)

type SurgicalService struct {
	usecase *usecase.SurgicalUsecase
}

func NewSurgicalService(usecase *usecase.SurgicalUsecase) *SurgicalService {
	return &SurgicalService{usecase: usecase}
}

func (s *SurgicalService) CreateSurgical(ctx context.Context, req *dto.CreateSurgicalRequest) (*dto.SurgicalResponse, error) {
	surgicalType := entity.SurgicalTypeFromCode(req.Type)
	
	flowSteps := make([]struct{ Step int; Name string }, 0, len(req.FlowSteps))
	for _, fs := range req.FlowSteps {
		flowSteps = append(flowSteps, struct{ Step int; Name string }{Step: fs.Step, Name: fs.Name})
	}

	surg, err := s.usecase.CreateSurgical(ctx,
		req.MedicalID,
		req.PatientID,
		req.DoctorIDs,
		req.DepartmentID,
		surgicalType,
		req.ScheduledTime,
		flowSteps,
	)
	if err != nil {
		return nil, err
	}

	return s.toResponse(surg), nil
}

func (s *SurgicalService) GetSurgical(ctx context.Context, id int64) (*dto.SurgicalResponse, error) {
	surg, err := s.usecase.GetSurgical(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(surg), nil
}

func (s *SurgicalService) StartSurgical(ctx context.Context, id int64, req *dto.StartSurgicalRequest) error {
	return s.usecase.StartSurgical(ctx, id, req.OperatorID)
}

func (s *SurgicalService) CompleteStep(ctx context.Context, id int64, req *dto.CompleteStepRequest) error {
	return s.usecase.CompleteStep(ctx, id, req.Step, req.OperatorID)
}

func (s *SurgicalService) CompleteSurgical(ctx context.Context, id int64, req *dto.CompleteSurgicalRequest) error {
	return s.usecase.CompleteSurgical(ctx, id, req.OperatorID)
}

func (s *SurgicalService) CancelSurgical(ctx context.Context, id int64, req *dto.CancelSurgicalRequest) error {
	return s.usecase.CancelSurgical(ctx, id, req.OperatorID, req.Reason)
}

func (s *SurgicalService) GetPatientSurgicals(ctx context.Context, patientID int64) ([]*dto.SurgicalResponse, error) {
	surgicals, err := s.usecase.GetPatientSurgicals(ctx, patientID)
	if err != nil {
		return nil, err
	}
	return s.toResponses(surgicals), nil
}

func (s *SurgicalService) ListSurgicals(ctx context.Context, page, pageSize int) (*dto.SurgicalListResponse, error) {
	surgicals, total, err := s.usecase.ListSurgicals(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	return &dto.SurgicalListResponse{
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		Surgicals: s.toResponseList(surgicals),
	}, nil
}

func (s *SurgicalService) DeleteSurgical(ctx context.Context, id int64) error {
	return s.usecase.DeleteSurgical(ctx, id)
}

func (s *SurgicalService) toResponse(surg *aggregate.Surgical) *dto.SurgicalResponse {
	flows := make([]dto.FlowResponse, 0, len(surg.Flows()))
	for _, f := range surg.Flows() {
		flows = append(flows, dto.FlowResponse{
			ID:        f.ID(),
			Step:      f.Step(),
			Name:      f.Name(),
			Status:    f.Status().Code(),
			StartTime: f.StartTime(),
			EndTime:   f.EndTime(),
		})
	}

	return &dto.SurgicalResponse{
		ID:            surg.ID(),
		MedicalID:     surg.MedicalID(),
		PatientID:     surg.PatientID(),
		DoctorIDs:     surg.DoctorIDs(),
		DepartmentID:  surg.DepartmentID(),
		Type:          surg.Type().Code(),
		TypeName:      surg.Type().String(),
		Status:        surg.Status().Code(),
		StatusName:    surg.Status().String(),
		ScheduledTime: surg.ScheduledTime(),
		Duration:      int64(surg.Duration()),
		Flows:         flows,
		CreatedAt:     surg.CreatedAt(),
		UpdatedAt:     surg.UpdatedAt(),
	}
}

func (s *SurgicalService) toResponses(surgicals []*aggregate.Surgical) []*dto.SurgicalResponse {
	responses := make([]*dto.SurgicalResponse, 0, len(surgicals))
	for _, surg := range surgicals {
		responses = append(responses, s.toResponse(surg))
	}
	return responses
}

func (s *SurgicalService) toResponseList(surgicals []*aggregate.Surgical) []dto.SurgicalResponse {
	responses := make([]dto.SurgicalResponse, 0, len(surgicals))
	for _, surg := range surgicals {
		responses = append(responses, *s.toResponse(surg))
	}
	return responses
}