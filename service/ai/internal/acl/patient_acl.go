package acl

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/log"
)

var (
	ErrPatientServiceUnavailable = errors.New("patient service unavailable")
	ErrPatientNotFound           = errors.New("patient not found")
)

type PatientInfo struct {
	ID          string
	Name        string
	Age         int
	Gender      string
	Phone       string
	MedicalHistory []MedicalHistoryItem
}

type MedicalHistoryItem struct {
	DiseaseName   string
	DiagnosisDate string
	Treatment     string
	Status        string
}

type PatientACL struct {
	patientServiceClient PatientServiceClient
	logger               *log.Helper
}

type PatientServiceClient interface {
	GetPatientByID(ctx context.Context, patientID string) (*PatientServiceResponse, error)
	GetPatientMedicalHistory(ctx context.Context, patientID string) (*MedicalHistoryResponse, error)
}

type PatientServiceResponse struct {
	ID     string
	Name   string
	Age    int
	Gender string
	Phone  string
}

type MedicalHistoryResponse struct {
	Items []MedicalHistoryItem
}

func NewPatientACL(client PatientServiceClient, logger log.Logger) *PatientACL {
	return &PatientACL{
		patientServiceClient: client,
		logger:               log.NewHelper(logger),
	}
}

func (acl *PatientACL) GetPatientInfo(ctx context.Context, patientID string) (*PatientInfo, error) {
	resp, err := acl.patientServiceClient.GetPatientByID(ctx, patientID)
	if err != nil {
		acl.logger.Warnf("Failed to get patient %s: %v", patientID, err)
		return nil, ErrPatientServiceUnavailable
	}

	if resp == nil {
		return nil, ErrPatientNotFound
	}

	historyResp, err := acl.patientServiceClient.GetPatientMedicalHistory(ctx, patientID)
	if err != nil {
		acl.logger.Warnf("Failed to get medical history for patient %s: %v", patientID, err)
	}

	history := make([]MedicalHistoryItem, 0)
	if historyResp != nil {
		history = historyResp.Items
	}

	return &PatientInfo{
		ID:           resp.ID,
		Name:         resp.Name,
		Age:          resp.Age,
		Gender:       resp.Gender,
		Phone:        resp.Phone,
		MedicalHistory: history,
	}, nil
}

func (acl *PatientACL) GetPatientContextForDiagnosis(ctx context.Context, patientID string) (string, error) {
	info, err := acl.GetPatientInfo(ctx, patientID)
	if err != nil {
		return "", err
	}

	context := formatPatientContext(info)
	return context, nil
}

func formatPatientContext(info *PatientInfo) string {
	context := "Patient Information:\n"
	context += "Name: " + info.Name + "\n"
	context += "Age: " + string(rune(info.Age)) + " years\n"
	context += "Gender: " + info.Gender + "\n"

	if len(info.MedicalHistory) > 0 {
		context += "\nMedical History:\n"
		for _, item := range info.MedicalHistory {
			context += "- " + item.DiseaseName + " (" + item.DiagnosisDate + "): " + item.Status + "\n"
		}
	}

	return context
}

type MockPatientServiceClient struct{}

func (m *MockPatientServiceClient) GetPatientByID(ctx context.Context, patientID string) (*PatientServiceResponse, error) {
	return &PatientServiceResponse{
		ID:     patientID,
		Name:   "Mock Patient",
		Age:    30,
		Gender: "Male",
		Phone:  "1234567890",
	}, nil
}

func (m *MockPatientServiceClient) GetPatientMedicalHistory(ctx context.Context, patientID string) (*MedicalHistoryResponse, error) {
	return &MedicalHistoryResponse{
		Items: []MedicalHistoryItem{
			{
				DiseaseName:   "Hypertension",
				DiagnosisDate: "2023-01-15",
				Treatment:     "Medication",
				Status:        "Controlled",
			},
		},
	}, nil
}