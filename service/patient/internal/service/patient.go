// service/patient/internal/service/patient.go
package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "intelligent-guidance-system/api/patient/v1"
	"intelligent-guidance-system/service/patient/internal/biz/dto"
	"intelligent-guidance-system/service/patient/internal/biz/usecase"
)

type PatientService struct {
	pb.UnimplementedPatientServiceServer

	registerUC *usecase.RegisterUsecase
	profileUC  *usecase.ProfileUsecase
	log        *log.Helper
}

func NewPatientService(registerUC *usecase.RegisterUsecase, profileUC *usecase.ProfileUsecase, logger log.Logger) *PatientService {
	return &PatientService{
		registerUC: registerUC,
		profileUC:  profileUC,
		log:        log.NewHelper(logger),
	}
}

func (s *PatientService) RegisterPatient(ctx context.Context, req *pb.RegisterPatientRequest) (*pb.RegisterPatientResponse, error) {
	s.log.Infof("注册患者: username=%s", req.Username)

	dtoReq := dto.PatientRegisterRequest{
		Username: req.Username,
		Password: req.Password,
		Phone:    req.Phone,
		Name:     req.Name,
		IDCard:   req.IdCard,
	}

	patient, err := s.registerUC.RegisterPatient(ctx, dtoReq)
	if err != nil {
		s.log.Errorf("注册患者失败: %v", err)
		return nil, err
	}

	token, err := s.registerUC.LoginPatient(ctx, dto.PatientLoginRequest{
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err != nil {
		s.log.Warnf("自动登录失败: %v", err)
	}

	return &pb.RegisterPatientResponse{
		PatientId: patient.ID,
		Token:     token.Token,
		Patient:   s.toPatientProto(patient),
	}, nil
}

func (s *PatientService) LoginPatient(ctx context.Context, req *pb.LoginPatientRequest) (*pb.LoginPatientResponse, error) {
	s.log.Infof("患者登录: phone=%s", req.Phone)

	dtoReq := dto.PatientLoginRequest{
		Phone:    req.Phone,
		Password: req.Password,
	}

	result, err := s.registerUC.LoginPatient(ctx, dtoReq)
	if err != nil {
		s.log.Errorf("登录失败: %v", err)
		return nil, err
	}

	return &pb.LoginPatientResponse{
		Token:   result.Token,
		Patient: s.toPatientProto(&result.Patient),
	}, nil
}

func (s *PatientService) GetPatient(ctx context.Context, req *pb.GetPatientRequest) (*pb.PatientDetail, error) {
	s.log.Infof("获取患者: id=%d", req.PatientId)

	patient, err := s.profileUC.GetPatient(ctx, req.PatientId)
	if err != nil {
		s.log.Errorf("获取患者失败: %v", err)
		return nil, err
	}

	return s.toPatientDetailProto(patient), nil
}

func (s *PatientService) UpdatePatient(ctx context.Context, req *pb.UpdatePatientRequest) (*pb.Patient, error) {
	s.log.Infof("更新患者: id=%d", req.PatientId)

	dtoReq := dto.PatientUpdateRequest{
		Name:    req.Name,
		Email:   req.Email,
		Address: req.Address,
	}

	patient, err := s.profileUC.UpdatePatient(ctx, req.PatientId, dtoReq)
	if err != nil {
		s.log.Errorf("更新患者失败: %v", err)
		return nil, err
	}

	return s.toPatientProto(patient), nil
}

func (s *PatientService) ListPatients(ctx context.Context, req *pb.ListPatientsRequest) (*pb.ListPatientsResponse, error) {
	s.log.Infof("获取患者列表: page=%d, pageSize=%d", req.Page, req.PageSize)

	dtoReq := dto.PatientListRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	result, err := s.profileUC.ListPatients(ctx, dtoReq)
	if err != nil {
		s.log.Errorf("获取患者列表失败: %v", err)
		return nil, err
	}

	list := make([]*pb.Patient, len(result.List))
	for i, p := range result.List {
		list[i] = s.toPatientProto(&p)
	}

	return &pb.ListPatientsResponse{
		Total: result.Total,
		List:  list,
	}, nil
}

func (s *PatientService) DeletePatient(ctx context.Context, req *pb.DeletePatientRequest) (*pb.DeletePatientResponse, error) {
	s.log.Infof("删除患者: id=%d", req.PatientId)

	err := s.profileUC.DeletePatient(ctx, req.PatientId)
	if err != nil {
		s.log.Errorf("删除患者失败: %v", err)
		return nil, err
	}

	return &pb.DeletePatientResponse{
		Success: true,
	}, nil
}

func (s *PatientService) AddMedicalHistory(ctx context.Context, req *pb.AddMedicalHistoryRequest) (*pb.MedicalHistory, error) {
	s.log.Infof("添加病史: patient_id=%d", req.PatientId)

	dtoReq := dto.AddMedicalHistoryRequest{
		PatientID:    req.PatientId,
		DiseaseName:  req.DiseaseName,
		DiagnoseDate: req.DiagnoseDate.AsTime(),
		Treatment:    req.Treatment,
	}

	history, err := s.profileUC.AddMedicalHistory(ctx, dtoReq)
	if err != nil {
		s.log.Errorf("添加病史失败: %v", err)
		return nil, err
	}

	return s.toMedicalHistoryProto(history), nil
}

func (s *PatientService) AddAllergy(ctx context.Context, req *pb.AddAllergyRequest) (*pb.Allergy, error) {
	s.log.Infof("添加过敏史: patient_id=%d", req.PatientId)

	dtoReq := dto.AddAllergyRequest{
		PatientID: req.PatientId,
		Allergen:  req.Allergen,
		Severity:  req.Severity,
		Reaction:  req.Reaction,
	}

	allergy, err := s.profileUC.AddAllergy(ctx, dtoReq)
	if err != nil {
		s.log.Errorf("添加过敏史失败: %v", err)
		return nil, err
	}

	return s.toAllergyProto(allergy), nil
}

func (s *PatientService) toPatientProto(patient *dto.PatientDTO) *pb.Patient {
	return &pb.Patient{
		Id:           patient.ID,
		OpenId:       patient.OpenID,
		Username:     patient.Username,
		Name:         patient.Name,
		Age:          patient.Age,
		Gender:       patient.Gender,
		IdCard:       patient.IDCard,
		Phone:        patient.Phone,
		Email:        patient.Email,
		Address:      patient.Address,
		Status:       patient.Status,
		CreatedAt:    timestamppb.New(patient.CreatedAt),
		UpdatedAt:    timestamppb.New(patient.UpdatedAt),
		MedicalCount: patient.MedicalCount,
		AllergyCount: patient.AllergyCount,
	}
}

func (s *PatientService) toPatientDetailProto(patient *dto.PatientDetailDTO) *pb.PatientDetail {
	medicalHistory := make([]*pb.MedicalHistory, len(patient.MedicalHistory))
	for i, h := range patient.MedicalHistory {
		medicalHistory[i] = s.toMedicalHistoryProto(&h)
	}

	allergies := make([]*pb.Allergy, len(patient.Allergies))
	for i, a := range patient.Allergies {
		allergies[i] = s.toAllergyProto(&a)
	}

	return &pb.PatientDetail{
		Id:              patient.ID,
		OpenId:          patient.OpenID,
		Username:        patient.Username,
		Name:            patient.Name,
		Age:             patient.Age,
		Gender:          patient.Gender,
		IdCardMasked:    patient.IDCardMasked,
		PhoneMasked:     patient.PhoneMasked,
		Email:           patient.Email,
		Address:         patient.Address,
		Status:          patient.Status,
		StatusName:      patient.StatusName,
		CreatedAt:       timestamppb.New(patient.CreatedAt),
		UpdatedAt:       timestamppb.New(patient.UpdatedAt),
		MedicalHistory:  medicalHistory,
		Allergies:       allergies,
	}
}

func (s *PatientService) toMedicalHistoryProto(history *dto.MedicalHistoryDTO) *pb.MedicalHistory {
	return &pb.MedicalHistory{
		Id:           history.ID,
		DiseaseName:  history.DiseaseName,
		DiagnoseDate: timestamppb.New(history.DiagnoseDate),
		Treatment:    history.Treatment,
		Status:       history.Status,
		StatusName:   history.StatusName,
		CreatedAt:    timestamppb.New(history.CreatedAt),
	}
}

func (s *PatientService) toAllergyProto(allergy *dto.AllergyDTO) *pb.Allergy {
	return &pb.Allergy{
		Id:           allergy.ID,
		Allergen:     allergy.Allergen,
		Severity:     allergy.Severity,
		SeverityName: allergy.SeverityName,
		Reaction:     allergy.Reaction,
		CreatedAt:    timestamppb.New(allergy.CreatedAt),
	}
}