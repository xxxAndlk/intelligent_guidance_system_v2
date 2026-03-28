// service/patient/internal/biz/usecase/profile_usecase.go
package usecase

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/log"

	"intelligent-guidance-system/service/patient/internal/domain/aggregate"
	"intelligent-guidance-system/service/patient/internal/domain/entity"
	"intelligent-guidance-system/service/patient/internal/domain/repository"
	"intelligent-guidance-system/service/patient/internal/domain/vo"
	"intelligent-guidance-system/service/patient/internal/biz/dto"
)

// ProfileUsecase 患者信息用例
type ProfileUsecase struct {
	repo repository.PatientRepository
	log  *log.Helper
}

// NewProfileUsecase 创建患者信息用例
func NewProfileUsecase(repo repository.PatientRepository, logger log.Logger) *ProfileUsecase {
	return &ProfileUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// GetPatient 获取患者信息
func (uc *ProfileUsecase) GetPatient(ctx context.Context, patientID int64) (*dto.PatientDetailDTO, error) {
	uc.log.Infof("获取患者信息: id=%d", patientID)

	patient, err := uc.repo.GetByID(ctx, aggregate.PatientID(patientID))
	if err != nil {
		uc.log.Errorf("获取患者失败: %v", err)
		return nil, err
	}

	if patient.IsDeleted() {
		return nil, errors.New("patient not found")
	}

	return uc.toDetailDTO(patient), nil
}

// UpdatePatient 更新患者信息
func (uc *ProfileUsecase) UpdatePatient(ctx context.Context, patientID int64, req dto.PatientUpdateRequest) (*dto.PatientDTO, error) {
	uc.log.Infof("更新患者信息: id=%d", patientID)

	patient, err := uc.repo.GetByID(ctx, aggregate.PatientID(patientID))
	if err != nil {
		return nil, err
	}

	if !patient.IsNormal() {
		return nil, errors.New("patient account is not available")
	}

	// 更新基本信息
	if req.Name != "" {
		patient.BasicInfo.UpdateName(req.Name)
	}

	// 更新联系信息
	if req.Email != "" {
		patient.Contact.UpdateEmail(req.Email)
	}
	if req.Address != "" {
		patient.Contact.UpdateAddress(req.Address)
	}

	// 保存
	if err := uc.repo.Update(ctx, patient); err != nil {
		uc.log.Errorf("更新患者失败: %v", err)
		return nil, err
	}

	return uc.toDTO(patient), nil
}

// AddMedicalHistory 添加病史
func (uc *ProfileUsecase) AddMedicalHistory(ctx context.Context, req dto.AddMedicalHistoryRequest) (*dto.MedicalHistoryDTO, error) {
	uc.log.Infof("添加病史: patient_id=%d, disease=%s", req.PatientID, req.DiseaseName)

	patient, err := uc.repo.GetByID(ctx, aggregate.PatientID(req.PatientID))
	if err != nil {
		return nil, err
	}

	if !patient.IsNormal() {
		return nil, errors.New("patient account is not available")
	}

	// 调用领域方法添加病史
	if err := patient.AddMedicalHistory(req.DiseaseName, req.DiagnoseDate, req.Treatment); err != nil {
		uc.log.Errorf("添加病史失败: %v", err)
		return nil, err
	}

	// 保存
	if err := uc.repo.Update(ctx, patient); err != nil {
		return nil, err
	}

	// 获取最新添加的病史
	lastHistory := patient.MedicalHistory[len(patient.MedicalHistory)-1]

	return dto.ToMedicalHistoryDTO(lastHistory), nil
}

// AddAllergy 添加过敏史
func (uc *ProfileUsecase) AddAllergy(ctx context.Context, req dto.AddAllergyRequest) (*dto.AllergyDTO, error) {
	uc.log.Infof("添加过敏史: patient_id=%d, allergen=%s", req.PatientID, req.Allergen)

	patient, err := uc.repo.GetByID(ctx, aggregate.PatientID(req.PatientID))
	if err != nil {
		return nil, err
	}

	if !patient.IsNormal() {
		return nil, errors.New("patient account is not available")
	}

	// 调用领域方法添加过敏史
	severity := entity.AllergySeverity(req.Severity)
	if err := patient.AddAllergy(req.Allergen, severity, req.Reaction); err != nil {
		uc.log.Errorf("添加过敏史失败: %v", err)
		return nil, err
	}

	// 保存
	if err := uc.repo.Update(ctx, patient); err != nil {
		return nil, err
	}

	// 获取最新添加的过敏史
	lastAllergy := patient.Allergies[len(patient.Allergies)-1]

	return dto.ToAllergyDTO(lastAllergy), nil
}

// ListPatients 患者列表
func (uc *ProfileUsecase) ListPatients(ctx context.Context, req dto.PatientListRequest) (*dto.PatientListResponse, error) {
	uc.log.Infof("获取患者列表: page=%d, pageSize=%d", req.Page, req.PageSize)

	patients, total, err := uc.repo.List(ctx, req.Page, req.PageSize)
	if err != nil {
		uc.log.Errorf("获取患者列表失败: %v", err)
		return nil, err
	}

	list := make([]dto.PatientDTO, len(patients))
	for i, p := range patients {
		list[i] = *uc.toDTO(p)
	}

	return &dto.PatientListResponse{
		Total: total,
		List:  list,
	}, nil
}

// DeletePatient 删除患者
func (uc *ProfileUsecase) DeletePatient(ctx context.Context, patientID int64) error {
	uc.log.Infof("删除患者: id=%d", patientID)

	return uc.repo.Delete(ctx, aggregate.PatientID(patientID))
}

// toDTO 转换聚合根为DTO
func (uc *ProfileUsecase) toDTO(patient *aggregate.Patient) *dto.PatientDTO {
	return &dto.PatientDTO{
		ID:           int64(patient.ID),
		OpenID:       patient.OpenID,
		Username:     patient.Username,
		Name:         patient.BasicInfo.Name,
		Age:          patient.BasicInfo.Age,
		Gender:       int32(patient.BasicInfo.Gender),
		IDCard:       patient.BasicInfo.IDCard.String(),
		Phone:        patient.Contact.Phone.String(),
		Email:        patient.Contact.Email,
		Address:      patient.Contact.Address,
		Status:       int32(patient.Status),
		CreatedAt:    patient.CreatedAt,
		UpdatedAt:    patient.UpdatedAt,
		MedicalCount: int32(len(patient.MedicalHistory)),
		AllergyCount: int32(len(patient.Allergies)),
	}
}

// toDetailDTO 转换聚合根为详情DTO
func (uc *ProfileUsecase) toDetailDTO(patient *aggregate.Patient) *dto.PatientDetailDTO {
	medicalHistory := make([]dto.MedicalHistoryDTO, len(patient.MedicalHistory))
	for i, h := range patient.MedicalHistory {
		medicalHistory[i] = dto.ToMedicalHistoryDTO(h)
	}

	allergies := make([]dto.AllergyDTO, len(patient.Allergies))
	for i, a := range patient.Allergies {
		allergies[i] = dto.ToAllergyDTO(a)
	}

	return &dto.PatientDetailDTO{
		ID:           int64(patient.ID),
		OpenID:       patient.OpenID,
		Username:     patient.Username,
		Name:         patient.BasicInfo.Name,
		Age:          patient.BasicInfo.Age,
		Gender:       int32(patient.BasicInfo.Gender),
		IDCardMasked: patient.BasicInfo.GetMaskedIDCard(),
		PhoneMasked:  patient.Contact.GetMaskedPhone(),
		Email:        patient.Contact.Email,
		Address:      patient.Contact.Address,
		Status:       int32(patient.Status),
		StatusName:   patient.Status.String(),
		CreatedAt:    patient.CreatedAt,
		UpdatedAt:    patient.UpdatedAt,
		MedicalHistory: medicalHistory,
		Allergies:      allergies,
	}
}