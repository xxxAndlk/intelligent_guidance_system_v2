package mysql

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"

	"intelligent-guidance-system/service/surgical/internal/domain/aggregate"
	"intelligent-guidance-system/service/surgical/internal/domain/entity"
	"intelligent-guidance-system/service/surgical/internal/domain/repository"
)

type SurgicalRepoImpl struct {
	db *gorm.DB
}

func NewSurgicalRepoImpl(db *gorm.DB) repository.SurgicalRepository {
	return &SurgicalRepoImpl{db: db}
}

func (r *SurgicalRepoImpl) Save(ctx context.Context, surgical *aggregate.Surgical) error {
	po := r.toPO(surgical)

	if po.ID == 0 {
		if err := r.db.WithContext(ctx).Create(&po).Error; err != nil {
			return err
		}
		surgical.SetID(po.ID)
	} else {
		if err := r.db.WithContext(ctx).Save(&po).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *SurgicalRepoImpl) FindByID(ctx context.Context, id int64) (*aggregate.Surgical, error) {
	var po SurgicalPO
	if err := r.db.WithContext(ctx).First(&po, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.toAggregate(&po), nil
}

func (r *SurgicalRepoImpl) FindByMedicalID(ctx context.Context, medicalID int64) (*aggregate.Surgical, error) {
	var po SurgicalPO
	if err := r.db.WithContext(ctx).Where("medical_id = ?", medicalID).First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.toAggregate(&po), nil
}

func (r *SurgicalRepoImpl) FindByPatientID(ctx context.Context, patientID int64) ([]*aggregate.Surgical, error) {
	var pos []SurgicalPO
	if err := r.db.WithContext(ctx).Where("patient_id = ?", patientID).Find(&pos).Error; err != nil {
		return nil, err
	}
	return r.toAggregates(pos), nil
}

func (r *SurgicalRepoImpl) FindByStatus(ctx context.Context, status int) ([]*aggregate.Surgical, error) {
	var pos []SurgicalPO
	statusCode := entity.SurgicalStatus(status).Code()
	if err := r.db.WithContext(ctx).Where("status = ?", statusCode).Find(&pos).Error; err != nil {
		return nil, err
	}
	return r.toAggregates(pos), nil
}

func (r *SurgicalRepoImpl) FindAll(ctx context.Context, page, pageSize int) ([]*aggregate.Surgical, int64, error) {
	var pos []SurgicalPO
	var total int64

	if err := r.db.WithContext(ctx).Model(&SurgicalPO{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	return r.toAggregates(pos), total, nil
}

func (r *SurgicalRepoImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&SurgicalPO{}, id).Error
}

func (r *SurgicalRepoImpl) toPO(surg *aggregate.Surgical) *SurgicalPO {
	doctorIDsJSON, _ := json.Marshal(surg.DoctorIDs())
	return &SurgicalPO{
		ID:            surg.ID(),
		MedicalID:     surg.MedicalID(),
		PatientID:     surg.PatientID(),
		DoctorIDs:     string(doctorIDsJSON),
		DepartmentID:  surg.DepartmentID(),
		Type:          surg.Type().Code(),
		Status:        surg.Status().Code(),
		ScheduledTime: surg.ScheduledTime(),
		Duration:      int64(surg.Duration()),
		CreatedAt:     surg.CreatedAt(),
		UpdatedAt:     surg.UpdatedAt(),
	}
}

func (r *SurgicalRepoImpl) toAggregate(po *SurgicalPO) *aggregate.Surgical {
	var doctorIDs []int64
	json.Unmarshal([]byte(po.DoctorIDs), &doctorIDs)

	return aggregate.ReconstructSurgical(
		po.ID,
		po.MedicalID,
		po.PatientID,
		doctorIDs,
		po.DepartmentID,
		entity.SurgicalTypeFromCode(po.Type),
		entity.SurgicalStatusFromCode(po.Status),
		po.ScheduledTime,
		time.Duration(po.Duration),
		nil,
		po.CreatedAt,
		po.UpdatedAt,
	)
}

func (r *SurgicalRepoImpl) toAggregates(pos []SurgicalPO) []*aggregate.Surgical {
	surgicals := make([]*aggregate.Surgical, 0, len(pos))
	for _, po := range pos {
		surgicals = append(surgicals, r.toAggregate(&po))
	}
	return surgicals
}

type SurgicalFlowRepoImpl struct {
	db *gorm.DB
}

func NewSurgicalFlowRepoImpl(db *gorm.DB) repository.SurgicalFlowRepository {
	return &SurgicalFlowRepoImpl{db: db}
}

func (r *SurgicalFlowRepoImpl) Save(ctx context.Context, flow *entity.SurgicalFlow) error {
	po := SurgicalFlowPO{
		ID:         flow.ID(),
		SurgicalID: flow.SurgicalID(),
		Step:       flow.Step(),
		Name:       flow.Name(),
		Status:     flow.Status().Code(),
		StartTime:  flow.StartTime(),
		EndTime:    flow.EndTime(),
		CreatedAt:  flow.CreatedAt(),
		UpdatedAt:  flow.UpdatedAt(),
	}

	if po.ID == 0 {
		if err := r.db.WithContext(ctx).Create(&po).Error; err != nil {
			return err
		}
		flow.SetID(po.ID)
	} else {
		if err := r.db.WithContext(ctx).Save(&po).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *SurgicalFlowRepoImpl) FindBySurgicalID(ctx context.Context, surgicalID int64) ([]*entity.SurgicalFlow, error) {
	var pos []SurgicalFlowPO
	if err := r.db.WithContext(ctx).Where("surgical_id = ?", surgicalID).Order("step").Find(&pos).Error; err != nil {
		return nil, err
	}

	flows := make([]*entity.SurgicalFlow, 0, len(pos))
	for _, po := range pos {
		flows = append(flows, entity.ReconstructSurgicalFlow(
			po.ID,
			po.SurgicalID,
			po.Step,
			po.Name,
			entity.FlowStepStatusFromCode(po.Status),
			po.StartTime,
			po.EndTime,
			po.CreatedAt,
			po.UpdatedAt,
		))
	}

	return flows, nil
}

func (r *SurgicalFlowRepoImpl) FindByID(ctx context.Context, id int64) (*entity.SurgicalFlow, error) {
	var po SurgicalFlowPO
	if err := r.db.WithContext(ctx).First(&po, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return entity.ReconstructSurgicalFlow(
		po.ID,
		po.SurgicalID,
		po.Step,
		po.Name,
		entity.FlowStepStatusFromCode(po.Status),
		po.StartTime,
		po.EndTime,
		po.CreatedAt,
		po.UpdatedAt,
	), nil
}

func (r *SurgicalFlowRepoImpl) DeleteBySurgicalID(ctx context.Context, surgicalID int64) error {
	return r.db.WithContext(ctx).Where("surgical_id = ?", surgicalID).Delete(&SurgicalFlowPO{}).Error
}