package mysql

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"intelligent-guidance-system/service/registration/internal/domain/aggregate"
	"intelligent-guidance-system/service/registration/internal/domain/entity"
	"intelligent-guidance-system/service/registration/internal/domain/repository"
	"intelligent-guidance-system/service/registration/internal/domain/vo"
)

type RegistrationRepoImpl struct {
	db *gorm.DB
}

func NewRegistrationRepoImpl(db *gorm.DB) repository.RegistrationRepository {
	return &RegistrationRepoImpl{db: db}
}

func (r *RegistrationRepoImpl) Save(ctx context.Context, registration *aggregate.Registration) error {
	po := r.toPO(registration)

	if po.ID == 0 {
		if err := r.db.WithContext(ctx).Create(&po).Error; err != nil {
			return err
		}
		registration.SetID(po.ID)
	} else {
		if err := r.db.WithContext(ctx).Save(&po).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *RegistrationRepoImpl) FindByID(ctx context.Context, id int64) (*aggregate.Registration, error) {
	var po RegistrationPO
	if err := r.db.WithContext(ctx).First(&po, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.toAggregate(&po), nil
}

func (r *RegistrationRepoImpl) FindByPatientID(ctx context.Context, patientID int64) ([]*aggregate.Registration, error) {
	var pos []RegistrationPO
	if err := r.db.WithContext(ctx).Where("patient_id = ?", patientID).Find(&pos).Error; err != nil {
		return nil, err
	}
	return r.toAggregates(pos), nil
}

func (r *RegistrationRepoImpl) FindByDoctorID(ctx context.Context, doctorID int64) ([]*aggregate.Registration, error) {
	var pos []RegistrationPO
	if err := r.db.WithContext(ctx).Where("doctor_id = ?", doctorID).Find(&pos).Error; err != nil {
		return nil, err
	}
	return r.toAggregates(pos), nil
}

func (r *RegistrationRepoImpl) FindByDepartmentID(ctx context.Context, departmentID int64) ([]*aggregate.Registration, error) {
	var pos []RegistrationPO
	if err := r.db.WithContext(ctx).Where("department_id = ?", departmentID).Find(&pos).Error; err != nil {
		return nil, err
	}
	return r.toAggregates(pos), nil
}

func (r *RegistrationRepoImpl) FindByStatus(ctx context.Context, status int) ([]*aggregate.Registration, error) {
	var pos []RegistrationPO
	statusCode := entity.RegistrationStatus(status).Code()
	if err := r.db.WithContext(ctx).Where("status = ?", statusCode).Find(&pos).Error; err != nil {
		return nil, err
	}
	return r.toAggregates(pos), nil
}

func (r *RegistrationRepoImpl) FindByAppointmentTime(ctx context.Context, startTime, endTime time.Time) ([]*aggregate.Registration, error) {
	var pos []RegistrationPO
	if err := r.db.WithContext(ctx).Where("appointment_time BETWEEN ? AND ?", startTime, endTime).Find(&pos).Error; err != nil {
		return nil, err
	}
	return r.toAggregates(pos), nil
}

func (r *RegistrationRepoImpl) FindAll(ctx context.Context, page, pageSize int) ([]*aggregate.Registration, int64, error) {
	var pos []RegistrationPO
	var total int64

	if err := r.db.WithContext(ctx).Model(&RegistrationPO{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	return r.toAggregates(pos), total, nil
}

func (r *RegistrationRepoImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&RegistrationPO{}, id).Error
}

func (r *RegistrationRepoImpl) toPO(reg *aggregate.Registration) *RegistrationPO {
	return &RegistrationPO{
		ID:              reg.ID(),
		PatientID:       reg.PatientID(),
		DoctorID:        reg.DoctorID(),
		DepartmentID:    reg.DepartmentID(),
		Type:            reg.Type().Code(),
		Status:          reg.Status().Code(),
		AppointmentTime: reg.AppointmentTime(),
		Fee:             reg.Fee().Amount(),
		CreatedAt:       reg.CreatedAt(),
		UpdatedAt:       reg.UpdatedAt(),
	}
}

func (r *RegistrationRepoImpl) toAggregate(po *RegistrationPO) *aggregate.Registration {
	fee, _ := vo.NewRegistrationFee(po.Fee, "CNY")
	return aggregate.ReconstructRegistration(
		po.ID,
		po.PatientID,
		po.DoctorID,
		po.DepartmentID,
		entity.RegistrationTypeFromCode(po.Type),
		entity.RegistrationStatusFromCode(po.Status),
		po.AppointmentTime,
		fee,
		po.CreatedAt,
		po.UpdatedAt,
	)
}

func (r *RegistrationRepoImpl) toAggregates(pos []RegistrationPO) []*aggregate.Registration {
	regs := make([]*aggregate.Registration, 0, len(pos))
	for _, po := range pos {
		regs = append(regs, r.toAggregate(&po))
	}
	return regs
}