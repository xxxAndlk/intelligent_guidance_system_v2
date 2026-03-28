package service

import (
	"context"

	v1 "intelligent_guidance_system_v2/api/medical/v1"
	"intelligent_guidance_system_v2/service/medical/internal/biz/dto"
	"intelligent_guidance_system_v2/service/medical/internal/biz/usecase"
)

type MedicalService struct {
	v1.UnimplementedMedicalServiceServer

	uc *usecase.MedicalUseCase
}

func NewMedicalService(uc *usecase.MedicalUseCase) *MedicalService {
	return &MedicalService{uc: uc}
}

func (s *MedicalService) CreateMedicalRecord(ctx context.Context, req *v1.CreateMedicalRecordRequest) (*v1.CreateMedicalRecordReply, error) {
	dtoReq := &dto.CreateMedicalRecordRequest{
		MedicalNumber: req.MedicalNumber,
		PatientID:     req.PatientId,
		DoctorID:      req.DoctorId,
		DepartmentID:  req.DepartmentId,
		Symptom:       req.Symptom,
	}

	result, err := s.uc.CreateMedicalRecord(ctx, dtoReq)
	if err != nil {
		return nil, err
	}

	return &v1.CreateMedicalRecordReply{
		MedicalRecord: s.toProto(result),
	}, nil
}

func (s *MedicalService) GetMedicalRecord(ctx context.Context, req *v1.GetMedicalRecordRequest) (*v1.GetMedicalRecordReply, error) {
	dtoReq := &dto.GetMedicalRecordRequest{ID: req.Id}

	result, err := s.uc.GetMedicalRecord(ctx, dtoReq)
	if err != nil {
		return nil, err
	}

	return &v1.GetMedicalRecordReply{
		MedicalRecord: s.toProto(result),
	}, nil
}

func (s *MedicalService) UpdateMedicalRecord(ctx context.Context, req *v1.UpdateMedicalRecordRequest) (*v1.UpdateMedicalRecordReply, error) {
	dtoReq := &dto.UpdateMedicalRecordRequest{
		ID:            req.Id,
		Symptom:       req.Symptom,
		Diagnosis:     req.Diagnosis,
		TreatmentPlan: req.TreatmentPlan,
	}

	result, err := s.uc.UpdateMedicalRecord(ctx, dtoReq)
	if err != nil {
		return nil, err
	}

	return &v1.UpdateMedicalRecordReply{
		MedicalRecord: s.toProto(result),
	}, nil
}

func (s *MedicalService) AddDisease(ctx context.Context, req *v1.AddDiseaseRequest) (*v1.AddDiseaseReply, error) {
	dtoReq := &dto.AddDiseaseRequest{
		MedicalRecordID: req.MedicalRecordId,
		DiseaseID:       req.DiseaseId,
		DiseaseName:     req.DiseaseName,
		Symptoms:        req.Symptoms,
		Diagnosis:       req.Diagnosis,
	}

	result, err := s.uc.AddDisease(ctx, dtoReq, 0)
	if err != nil {
		return nil, err
	}

	return &v1.AddDiseaseReply{
		MedicalRecord: s.toProto(result),
	}, nil
}

func (s *MedicalService) AddPrescription(ctx context.Context, req *v1.AddPrescriptionRequest) (*v1.AddPrescriptionReply, error) {
	dtoReq := &dto.AddPrescriptionRequest{
		MedicalRecordID: req.MedicalRecordId,
		DrugID:          req.DrugId,
		DrugName:        req.DrugName,
		Quantity:        req.Quantity,
		Price:           req.Price,
		Usage:           req.Usage,
	}

	result, err := s.uc.AddPrescription(ctx, dtoReq, 0)
	if err != nil {
		return nil, err
	}

	return &v1.AddPrescriptionReply{
		MedicalRecord: s.toProto(result),
	}, nil
}

func (s *MedicalService) UpdateStatus(ctx context.Context, req *v1.UpdateStatusRequest) (*v1.UpdateStatusReply, error) {
	dtoReq := &dto.UpdateStatusRequest{
		ID:         req.Id,
		Status:     int(req.Status),
		OperatorID: req.OperatorId,
		Reason:     req.Reason,
	}

	result, err := s.uc.UpdateStatus(ctx, dtoReq)
	if err != nil {
		return nil, err
	}

	return &v1.UpdateStatusReply{
		MedicalRecord: s.toProto(result),
	}, nil
}

func (s *MedicalService) CompleteMedical(ctx context.Context, req *v1.CompleteMedicalRequest) (*v1.CompleteMedicalReply, error) {
	dtoReq := &dto.CompleteMedicalRequest{
		ID:         req.Id,
		OperatorID: req.OperatorId,
	}

	result, err := s.uc.CompleteMedical(ctx, dtoReq)
	if err != nil {
		return nil, err
	}

	return &v1.CompleteMedicalReply{
		MedicalRecord: s.toProto(result),
	}, nil
}

func (s *MedicalService) CancelMedical(ctx context.Context, req *v1.CancelMedicalRequest) (*v1.CancelMedicalReply, error) {
	dtoReq := &dto.CancelMedicalRequest{
		ID:         req.Id,
		OperatorID: req.OperatorId,
		Reason:     req.Reason,
	}

	result, err := s.uc.CancelMedical(ctx, dtoReq)
	if err != nil {
		return nil, err
	}

	return &v1.CancelMedicalReply{
		MedicalRecord: s.toProto(result),
	}, nil
}

func (s *MedicalService) ListMedicalRecords(ctx context.Context, req *v1.ListMedicalRecordsRequest) (*v1.ListMedicalRecordsReply, error) {
	dtoReq := &dto.ListMedicalRecordsRequest{
		PatientID:     req.PatientId,
		DoctorID:      req.DoctorId,
		DepartmentID:  req.DepartmentId,
		Status:        int(req.Status),
		Page:          int(req.Page),
PageSize:     int(req.PageSize),
	}

	if dtoReq.Page <= 0 {
		dtoReq.Page = 1
	}
	if dtoReq.PageSize <= 0 {
		dtoReq.PageSize = 10
	}

	result, err := s.uc.ListMedicalRecords(ctx, dtoReq)
	if err != nil {
		return nil, err
	}

	records := make([]*v1.MedicalRecord, len(result.MedicalRecords))
	for i, r := range result.MedicalRecords {
		records[i] = s.toProto(&r)
	}

	return &v1.ListMedicalRecordsReply{
		MedicalRecords: records,
		Total:          result.Total,
	}, nil
}

func (s *MedicalService) toProto(d *dto.MedicalRecordDTO) *v1.MedicalRecord {
	record := &v1.MedicalRecord{
		Id:            d.ID,
		MedicalNumber: d.MedicalNumber,
		PatientId:     d.PatientID,
		DoctorId:      d.DoctorID,
		DepartmentId:  d.DepartmentID,
		Status:        v1.MedicalStatus(d.Status),
		Symptom:       d.Symptom,
		Diagnosis:     d.Diagnosis,
		TreatmentPlan: d.TreatmentPlan,
	}

	diseases := make([]*v1.DiseaseItem, len(d.Diseases))
	for i, di := range d.Diseases {
		diseases[i] = &v1.DiseaseItem{
			DiseaseId:   di.DiseaseID,
			DiseaseName: di.DiseaseName,
			Symptoms:    di.Symptoms,
			Diagnosis:   di.Diagnosis,
		}
	}
	record.Diseases = diseases

	prescriptions := make([]*v1.PrescriptionItem, len(d.Prescriptions))
	for i, pi := range d.Prescriptions {
		prescriptions[i] = &v1.PrescriptionItem{
			DrugId:   pi.DrugID,
			DrugName: pi.DrugName,
			Quantity: pi.Quantity,
			Price:    pi.Price,
			Usage:    pi.Usage,
		}
	}
	record.Prescriptions = prescriptions

	if d.Payment.Amount > 0 {
		record.Payment = &v1.PaymentInfo{
			Amount: d.Payment.Amount,
			Method: d.Payment.Method,
			Status: d.Payment.Status,
		}
	}

	statusHistory := make([]*v1.StatusChange, len(d.Timeline))
	for i, sc := range d.Timeline {
		statusHistory[i] = &v1.StatusChange{
			FromStatus: v1.MedicalStatus(sc.FromStatus),
			ToStatus:   v1.MedicalStatus(sc.ToStatus),
			OperatorId: sc.OperatorID,
			Reason:     sc.Reason,
		}
	}
	record.StatusHistory = statusHistory

	return record
}