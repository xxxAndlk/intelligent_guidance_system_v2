// service/patient/internal/data/mysql/patient_repo_impl.go
package mysql

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"

	"intelligent-guidance-system/service/patient/internal/domain/aggregate"
	"intelligent-guidance-system/service/patient/internal/domain/entity"
	"intelligent-guidance-system/service/patient/internal/domain/repository"
	"intelligent-guidance-system/service/patient/internal/domain/vo"
)

// PatientRepositoryImpl 患者仓储实现
type PatientRepositoryImpl struct {
	db  *gorm.DB
	log *log.Helper
}

// NewPatientRepository 创建患者仓储
func NewPatientRepository(db *gorm.DB, logger log.Logger) repository.PatientRepository {
	return &PatientRepositoryImpl{
		db:  db,
		log: log.NewHelper(logger),
	}
}

// Save 保存患者
func (r *PatientRepositoryImpl) Save(ctx context.Context, patient *aggregate.Patient) error {
	po := r.toPO(patient)

	result := r.db.WithContext(ctx).Create(&po)
	if result.Error != nil {
		r.log.Errorf("保存患者失败: %v", result.Error)
		return result.Error
	}

	patient.ID = aggregate.PatientID(po.ID)

	if len(patient.Allergies) > 0 {
		allergyPOs := r.toAllergyPOs(patient.ID, patient.Allergies)
		if err := r.db.WithContext(ctx).Create(&allergyPOs).Error; err != nil {
			r.log.Errorf("保存过敏史失败: %v", err)
			return err
		}
	}

	if len(patient.MedicalHistory) > 0 {
		historyPOs := r.toMedicalHistoryPOs(patient.ID, patient.MedicalHistory)
		if err := r.db.WithContext(ctx).Create(&historyPOs).Error; err != nil {
			r.log.Errorf("保存病史失败: %v", err)
			return err
		}
	}

	r.log.Infof("保存患者成功: id=%d", patient.ID)
	return nil
}

// Update 更新患者
func (r *PatientRepositoryImpl) Update(ctx context.Context, patient *aggregate.Patient) error {
	po := r.toPO(patient)

	result := r.db.WithContext(ctx).Model(&PatientPO{}).Where("id = ?", patient.ID).Updates(po)
	if result.Error != nil {
		r.log.Errorf("更新患者失败: %v", result.Error)
		return result.Error
	}

	for _, allergy := range patient.Allergies {
		allergyPO := r.toAllergyPO(patient.ID, allergy)
		if int64(allergy.ID) == 0 {
			r.db.WithContext(ctx).Create(&allergyPO)
		} else {
			r.db.WithContext(ctx).Model(&AllergyPO{}).Where("id = ?", allergy.ID).Updates(allergyPO)
		}
	}

	for _, history := range patient.MedicalHistory {
		historyPO := r.toMedicalHistoryPO(patient.ID, history)
		if int64(history.ID) == 0 {
			r.db.WithContext(ctx).Create(&historyPO)
		} else {
			r.db.WithContext(ctx).Model(&MedicalHistoryPO{}).Where("id = ?", history.ID).Updates(historyPO)
		}
	}

	return nil
}

// GetByID 根据ID获取患者
func (r *PatientRepositoryImpl) GetByID(ctx context.Context, id aggregate.PatientID) (*aggregate.Patient, error) {
	var po PatientPO
	result := r.db.WithContext(ctx).Preload("Allergies").Preload("MedicalHistories").First(&po, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("patient not found")
		}
		return nil, result.Error
	}

	return r.toAggregate(&po), nil
}

// GetByOpenID 根据OpenID获取患者
func (r *PatientRepositoryImpl) GetByOpenID(ctx context.Context, openID string) (*aggregate.Patient, error) {
	var po PatientPO
	result := r.db.WithContext(ctx).Preload("Allergies").Preload("MedicalHistories").Where("open_id = ?", openID).First(&po)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("patient not found")
		}
		return nil, result.Error
	}

	return r.toAggregate(&po), nil
}

// GetByPhone 根据手机号获取患者
func (r *PatientRepositoryImpl) GetByPhone(ctx context.Context, phone string) (*aggregate.Patient, error) {
	var po PatientPO
	result := r.db.WithContext(ctx).Preload("Allergies").Preload("MedicalHistories").Where("phone = ?", phone).First(&po)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("patient not found")
		}
		return nil, result.Error
	}

	return r.toAggregate(&po), nil
}

// List 获取患者列表
func (r *PatientRepositoryImpl) List(ctx context.Context, page, pageSize int32) ([]*aggregate.Patient, int64, error) {
	var pos []PatientPO
	var total int64

	offset := (page - 1) * pageSize

	r.db.WithContext(ctx).Model(&PatientPO{}).Count(&total)

	result := r.db.WithContext(ctx).Offset(int(offset)).Limit(int(pageSize)).Find(&pos)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	patients := make([]*aggregate.Patient, len(pos))
	for i, po := range pos {
		patients[i] = r.toAggregate(&po)
	}

	return patients, total, nil
}

// Delete 删除患者
func (r *PatientRepositoryImpl) Delete(ctx context.Context, id aggregate.PatientID) error {
	result := r.db.WithContext(ctx).Delete(&PatientPO{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// ExistsByPhone 检查手机号是否已存在
func (r *PatientRepositoryImpl) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&PatientPO{}).Where("phone = ?", phone).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

// toPO 转换聚合根为持久化对象
func (r *PatientRepositoryImpl) toPO(patient *aggregate.Patient) PatientPO {
	return PatientPO{
		ID:               int64(patient.ID),
		OpenID:           patient.OpenID,
		Username:         patient.Username,
		Password:         patient.Password,
		Name:             patient.BasicInfo.Name,
		Age:              patient.BasicInfo.Age,
		Gender:           int32(patient.BasicInfo.Gender),
		IDCard:           patient.BasicInfo.IDCard.String(),
		Phone:            patient.Contact.Phone.String(),
		Email:            patient.Contact.Email,
		Address:          patient.Contact.Address,
		EmergencyName:    patient.Contact.EmergencyContact.Name,
		EmergencyPhone:   patient.Contact.EmergencyContact.Phone.String(),
		EmergencyRelation: patient.Contact.EmergencyContact.Relation,
		Status:           int32(patient.Status),
		CreatedAt:        patient.CreatedAt,
		UpdatedAt:        patient.UpdatedAt,
	}
}

// toAggregate 转换持久化对象为聚合根
func (r *PatientRepositoryImpl) toAggregate(po *PatientPO) *aggregate.Patient {
	phone, _ := vo.NewPhone(po.Phone)
	idCard, _ := vo.NewIDCard(po.IDCard)

	patient := &aggregate.Patient{
		ID:        aggregate.PatientID(po.ID),
		OpenID:    po.OpenID,
		Username:  po.Username,
		Password:  po.Password,
		BasicInfo: vo.NewPatientBasicInfo(po.Name, idCard),
		Contact:   vo.NewContactInfo(phone, po.Email, po.Address),
		Status:    aggregate.PatientStatus(po.Status),
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,
		Allergies:      []entity.Allergy{},
		MedicalHistory: []entity.MedicalHistory{},
		Events:         []aggregate.DomainEvent{},
	}

	for _, allergyPO := range po.Allergies {
		patient.Allergies = append(patient.Allergies, entity.Allergy{
			ID:        entity.AllergyID(allergyPO.ID),
			Allergen:  allergyPO.Allergen,
			Severity:  entity.AllergySeverity(allergyPO.Severity),
			Reaction:  allergyPO.Reaction,
			CreatedAt: allergyPO.CreatedAt,
			UpdatedAt: allergyPO.UpdatedAt,
		})
	}

	for _, historyPO := range po.MedicalHistories {
		patient.MedicalHistory = append(patient.MedicalHistory, entity.MedicalHistory{
			ID:           entity.HistoryID(historyPO.ID),
			DiseaseName:  historyPO.DiseaseName,
			DiagnoseDate: historyPO.DiagnoseDate,
			Treatment:    historyPO.Treatment,
			Status:       entity.HistoryStatus(historyPO.Status),
			CreatedAt:    historyPO.CreatedAt,
			UpdatedAt:    historyPO.UpdatedAt,
		})
	}

	return patient
}

// toAllergyPO 转换过敏史实体为持久化对象
func (r *PatientRepositoryImpl) toAllergyPO(patientID aggregate.PatientID, allergy entity.Allergy) AllergyPO {
	return AllergyPO{
		ID:        int64(allergy.ID),
		PatientID: int64(patientID),
		Allergen:  allergy.Allergen,
		Severity:  int32(allergy.Severity),
		Reaction:  allergy.Reaction,
		CreatedAt: allergy.CreatedAt,
		UpdatedAt: allergy.UpdatedAt,
	}
}

// toAllergyPOs 批量转换过敏史
func (r *PatientRepositoryImpl) toAllergyPOs(patientID aggregate.PatientID, allergies []entity.Allergy) []AllergyPO {
	pos := make([]AllergyPO, len(allergies))
	for i, a := range allergies {
		pos[i] = r.toAllergyPO(patientID, a)
	}
	return pos
}

// toMedicalHistoryPO 转换病史实体为持久化对象
func (r *PatientRepositoryImpl) toMedicalHistoryPO(patientID aggregate.PatientID, history entity.MedicalHistory) MedicalHistoryPO {
	return MedicalHistoryPO{
		ID:           int64(history.ID),
		PatientID:    int64(patientID),
		DiseaseName:  history.DiseaseName,
		DiagnoseDate: history.DiagnoseDate,
		Treatment:    history.Treatment,
		Status:       int32(history.Status),
		CreatedAt:    history.CreatedAt,
		UpdatedAt:    history.UpdatedAt,
	}
}

// toMedicalHistoryPOs 批量转换病史
func (r *PatientRepositoryImpl) toMedicalHistoryPOs(patientID aggregate.PatientID, histories []entity.MedicalHistory) []MedicalHistoryPO {
	pos := make([]MedicalHistoryPO, len(histories))
	for i, h := range histories {
		pos[i] = r.toMedicalHistoryPO(patientID, h)
	}
	return pos
}