// service/patient/internal/domain/repository/patient_repo.go
package repository

import (
	"context"

	"intelligent-guidance-system/service/patient/internal/domain/aggregate"
)

// PatientRepository 患者仓储接口
type PatientRepository interface {
	// Save 保存患者
	Save(ctx context.Context, patient *aggregate.Patient) error

	// Update 更新患者
	Update(ctx context.Context, patient *aggregate.Patient) error

	// GetByID 根据ID获取患者
	GetByID(ctx context.Context, id aggregate.PatientID) (*aggregate.Patient, error)

	// GetByOpenID 根据OpenID获取患者
	GetByOpenID(ctx context.Context, openID string) (*aggregate.Patient, error)

	// GetByPhone 根据手机号获取患者
	GetByPhone(ctx context.Context, phone string) (*aggregate.Patient, error)

	// List 获取患者列表
	List(ctx context.Context, page, pageSize int32) ([]*aggregate.Patient, int64, error)

	// Delete 删除患者
	Delete(ctx context.Context, id aggregate.PatientID) error

	// ExistsByPhone 检查手机号是否已存在
	ExistsByPhone(ctx context.Context, phone string) (bool, error)
}