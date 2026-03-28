// service/patient/internal/biz/usecase/register_usecase.go
package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"intelligent-guidance-system/service/patient/internal/domain/aggregate"
	"intelligent-guidance-system/service/patient/internal/domain/event"
	"intelligent-guidance-system/service/patient/internal/domain/repository"
	"intelligent-guidance-system/service/patient/internal/domain/vo"
	"intelligent-guidance-system/service/patient/internal/biz/dto"
)

// RegisterUsecase 注册用例
type RegisterUsecase struct {
	repo     repository.PatientRepository
	jwtKey   string
	log      *log.Helper
}

// NewRegisterUsecase 创建注册用例
func NewRegisterUsecase(repo repository.PatientRepository, jwtKey string, logger log.Logger) *RegisterUsecase {
	return &RegisterUsecase{
		repo:   repo,
		jwtKey: jwtKey,
		log:    log.NewHelper(logger),
	}
}

// RegisterPatient 注册患者
func (uc *RegisterUsecase) RegisterPatient(ctx context.Context, req dto.PatientRegisterRequest) (*dto.PatientDTO, error) {
	uc.log.Infof("注册患者: username=%s, phone=%s", req.Username, req.Phone)

	phone, err := vo.NewPhone(req.Phone)
	if err != nil {
		uc.log.Errorf("手机号验证失败: %v", err)
		return nil, errors.New("invalid phone number")
	}

	exists, err := uc.repo.ExistsByPhone(ctx, req.Phone)
	if err != nil {
		uc.log.Errorf("检查手机号失败: %v", err)
		return nil, err
	}
	if exists {
		return nil, errors.New("phone number already registered")
	}

	idCard, err := vo.NewIDCard(req.IDCard)
	if err != nil {
		uc.log.Errorf("身份证验证失败: %v", err)
		return nil, errors.New("invalid ID card")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		uc.log.Errorf("密码加密失败: %v", err)
		return nil, errors.New("password encryption failed")
	}

	patient := aggregate.NewPatient(req.Username, string(hashedPassword), phone, req.Name, idCard)

	patient.Events = append(patient.Events, event.NewPatientRegisteredEvent(patient.ID))

	if err := uc.repo.Save(ctx, patient); err != nil {
		uc.log.Errorf("保存患者失败: %v", err)
		return nil, err
	}

	uc.log.Infof("患者注册成功: id=%d, username=%s", patient.ID, patient.Username)

	return uc.toDTO(patient), nil
}

// LoginPatient 患者登录
func (uc *RegisterUsecase) LoginPatient(ctx context.Context, req dto.PatientLoginRequest) (*dto.PatientLoginResponse, error) {
	uc.log.Infof("患者登录: phone=%s", req.Phone)

	phone, err := vo.NewPhone(req.Phone)
	if err != nil {
		return nil, errors.New("invalid phone number")
	}

	patient, err := uc.repo.GetByPhone(ctx, phone.String())
	if err != nil {
		uc.log.Errorf("查找患者失败: %v", err)
		return nil, errors.New("patient not found")
	}

	if !patient.IsNormal() {
		return nil, errors.New("patient account is not available")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(patient.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid password")
	}

	token, err := uc.generateToken(patient.ID)
	if err != nil {
		uc.log.Errorf("生成Token失败: %v", err)
		return nil, errors.New("token generation failed")
	}

	uc.log.Infof("患者登录成功: id=%d", patient.ID)

	return &dto.PatientLoginResponse{
		Token:   token,
		Patient: *uc.toDTO(patient),
	}, nil
}

// generateToken 生成JWT Token
func (uc *RegisterUsecase) generateToken(patientID aggregate.PatientID) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"patient_id": int64(patientID),
		"exp":        jwt.NewNumericDate(now.Add(24 * time.Hour)),
		"iat":        jwt.NewNumericDate(now),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.jwtKey))
}

// toDTO 转换聚合根为DTO
func (uc *RegisterUsecase) toDTO(patient *aggregate.Patient) *dto.PatientDTO {
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