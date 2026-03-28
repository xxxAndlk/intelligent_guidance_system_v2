package usecase

import (
	"context"
	"errors"
	"time"

	"intelligent-guidance-system/service/surgical/internal/domain/aggregate"
	"intelligent-guidance-system/service/surgical/internal/domain/entity"
	"intelligent-guidance-system/service/surgical/internal/domain/repository"
)

var (
	ErrSurgicalNotFound = errors.New("surgical not found")
)

type SurgicalUsecase struct {
	surgicalRepo repository.SurgicalRepository
	flowRepo     repository.SurgicalFlowRepository
}

func NewSurgicalUsecase(surgicalRepo repository.SurgicalRepository, flowRepo repository.SurgicalFlowRepository) *SurgicalUsecase {
	return &SurgicalUsecase{
		surgicalRepo: surgicalRepo,
		flowRepo:     flowRepo,
	}
}

func (u *SurgicalUsecase) CreateSurgical(ctx context.Context,
	medicalID int64,
	patientID int64,
	doctorIDs []int64,
	departmentID int64,
	surgicalType entity.SurgicalType,
	scheduledTime time.Time,
	flowSteps []struct{ Step int; Name string },
) (*aggregate.Surgical, error) {
	surg, err := aggregate.NewSurgical(medicalID, patientID, doctorIDs, departmentID, surgicalType, scheduledTime)
	if err != nil {
		return nil, err
	}

	if err := u.surgicalRepo.Save(ctx, surg); err != nil {
		return nil, err
	}

	for _, step := range flowSteps {
		if err := surg.AddFlowStep(step.Step, step.Name); err != nil {
			return nil, err
		}
	}

	for _, flow := range surg.Flows() {
		if err := u.flowRepo.Save(ctx, flow); err != nil {
			return nil, err
		}
	}

	return surg, nil
}

func (u *SurgicalUsecase) GetSurgical(ctx context.Context, id int64) (*aggregate.Surgical, error) {
	surg, err := u.surgicalRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if surg == nil {
		return nil, ErrSurgicalNotFound
	}

	flows, err := u.flowRepo.FindBySurgicalID(ctx, id)
	if err != nil {
		return nil, err
	}

	return ReconstructSurgicalWithFlows(surg, flows), nil
}

func (u *SurgicalUsecase) StartSurgical(ctx context.Context, id int64, operatorID int64) error {
	surg, err := u.GetSurgical(ctx, id)
	if err != nil {
		return err
	}

	if err := surg.StartSurgical(operatorID); err != nil {
		return err
	}

	return u.surgicalRepo.Save(ctx, surg)
}

func (u *SurgicalUsecase) CompleteStep(ctx context.Context, id int64, step int, operatorID int64) error {
	surg, err := u.GetSurgical(ctx, id)
	if err != nil {
		return err
	}

	if err := surg.CompleteStep(step, operatorID); err != nil {
		return err
	}

	for _, flow := range surg.Flows() {
		if flow.Step() == step {
			if err := u.flowRepo.Save(ctx, flow); err != nil {
				return err
			}
			break
		}
	}

	return u.surgicalRepo.Save(ctx, surg)
}

func (u *SurgicalUsecase) CompleteSurgical(ctx context.Context, id int64, operatorID int64) error {
	surg, err := u.GetSurgical(ctx, id)
	if err != nil {
		return err
	}

	if err := surg.CompleteSurgical(operatorID); err != nil {
		return err
	}

	return u.surgicalRepo.Save(ctx, surg)
}

func (u *SurgicalUsecase) CancelSurgical(ctx context.Context, id int64, operatorID int64, reason string) error {
	surg, err := u.GetSurgical(ctx, id)
	if err != nil {
		return err
	}

	if err := surg.CancelSurgical(operatorID, reason); err != nil {
		return err
	}

	return u.surgicalRepo.Save(ctx, surg)
}

func (u *SurgicalUsecase) GetPatientSurgicals(ctx context.Context, patientID int64) ([]*aggregate.Surgical, error) {
	return u.surgicalRepo.FindByPatientID(ctx, patientID)
}

func (u *SurgicalUsecase) ListSurgicals(ctx context.Context, page, pageSize int) ([]*aggregate.Surgical, int64, error) {
	return u.surgicalRepo.FindAll(ctx, page, pageSize)
}

func (u *SurgicalUsecase) DeleteSurgical(ctx context.Context, id int64) error {
	if err := u.flowRepo.DeleteBySurgicalID(ctx, id); err != nil {
		return err
	}
	return u.surgicalRepo.Delete(ctx, id)
}

func ReconstructSurgicalWithFlows(surg *aggregate.Surgical, flows []*entity.SurgicalFlow) *aggregate.Surgical {
	return aggregate.ReconstructSurgical(
		surg.ID(),
		surg.MedicalID(),
		surg.PatientID(),
		surg.DoctorIDs(),
		surg.DepartmentID(),
		surg.Type(),
		surg.Status(),
		surg.ScheduledTime(),
		surg.Duration(),
		flows,
		surg.CreatedAt(),
		surg.UpdatedAt(),
	)
}