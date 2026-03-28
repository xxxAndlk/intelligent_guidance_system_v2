package repository

import (
	"context"

	"intelligent-guidance-system/service/surgical/internal/domain/aggregate"
	"intelligent-guidance-system/service/surgical/internal/domain/entity"
)

type SurgicalRepository interface {
	Save(ctx context.Context, surgical *aggregate.Surgical) error
	FindByID(ctx context.Context, id int64) (*aggregate.Surgical, error)
	FindByMedicalID(ctx context.Context, medicalID int64) (*aggregate.Surgical, error)
	FindByPatientID(ctx context.Context, patientID int64) ([]*aggregate.Surgical, error)
	FindByStatus(ctx context.Context, status int) ([]*aggregate.Surgical, error)
	FindAll(ctx context.Context, page, pageSize int) ([]*aggregate.Surgical, int64, error)
	Delete(ctx context.Context, id int64) error
}

type SurgicalFlowRepository interface {
	Save(ctx context.Context, flow *entity.SurgicalFlow) error
	FindBySurgicalID(ctx context.Context, surgicalID int64) ([]*entity.SurgicalFlow, error)
	FindByID(ctx context.Context, id int64) (*entity.SurgicalFlow, error)
	DeleteBySurgicalID(ctx context.Context, surgicalID int64) error
}