package acl

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/log"
)

var (
	ErrDoctorServiceUnavailable = errors.New("doctor service unavailable")
	ErrDoctorNotFound           = errors.New("doctor not found")
)

type DoctorInfo struct {
	ID           string
	Name         string
	Specialty    string
	Hospital     string
	Department   string
	Title        string
	Experience   int
	Rating       float32
	Availability string
}

type DoctorACL struct {
	doctorServiceClient DoctorServiceClient
	logger              *log.Helper
}

type DoctorServiceClient interface {
	GetDoctorByID(ctx context.Context, doctorID string) (*DoctorServiceResponse, error)
	SearchDoctors(ctx context.Context, specialty string, limit int) ([]*DoctorServiceResponse, error)
	GetDoctorAvailability(ctx context.Context, doctorID string) (*AvailabilityResponse, error)
}

type DoctorServiceResponse struct {
	ID         string
	Name       string
	Specialty  string
	Hospital   string
	Department string
	Title      string
	Experience int
	Rating     float32
}

type AvailabilityResponse struct {
	AvailableSlots []string
	Status         string
}

func NewDoctorACL(client DoctorServiceClient, logger log.Logger) *DoctorACL {
	return &DoctorACL{
		doctorServiceClient: client,
		logger:              log.NewHelper(logger),
	}
}

func (acl *DoctorACL) GetDoctorInfo(ctx context.Context, doctorID string) (*DoctorInfo, error) {
	resp, err := acl.doctorServiceClient.GetDoctorByID(ctx, doctorID)
	if err != nil {
		acl.logger.Warnf("Failed to get doctor %s: %v", doctorID, err)
		return nil, ErrDoctorServiceUnavailable
	}

	if resp == nil {
		return nil, ErrDoctorNotFound
	}

	availability, err := acl.doctorServiceClient.GetDoctorAvailability(ctx, doctorID)
	if err != nil {
		acl.logger.Warnf("Failed to get availability for doctor %s: %v", doctorID, err)
	}

	availStatus := "Unknown"
	if availability != nil {
		availStatus = availability.Status
	}

	return &DoctorInfo{
		ID:           resp.ID,
		Name:         resp.Name,
		Specialty:    resp.Specialty,
		Hospital:     resp.Hospital,
		Department:   resp.Department,
		Title:        resp.Title,
		Experience:   resp.Experience,
		Rating:       resp.Rating,
		Availability: availStatus,
	}, nil
}

func (acl *DoctorACL) SearchDoctorsBySpecialty(ctx context.Context, specialty string, limit int) ([]*DoctorInfo, error) {
	resps, err := acl.doctorServiceClient.SearchDoctors(ctx, specialty, limit)
	if err != nil {
		acl.logger.Warnf("Failed to search doctors by specialty %s: %v", specialty, err)
		return nil, ErrDoctorServiceUnavailable
	}

	doctors := make([]*DoctorInfo, len(resps))
	for i, resp := range resps {
		doctors[i] = &DoctorInfo{
			ID:         resp.ID,
			Name:       resp.Name,
			Specialty:  resp.Specialty,
			Hospital:   resp.Hospital,
			Department: resp.Department,
			Title:      resp.Title,
			Experience: resp.Experience,
			Rating:     resp.Rating,
		}
	}

	return doctors, nil
}

func (acl *DoctorACL) FormatDoctorRecommendation(doctor *DoctorInfo) string {
	return doctor.Name + " - " + doctor.Specialty + " (" + doctor.Hospital + ", " + doctor.Title + ")"
}

type MockDoctorServiceClient struct{}

func (m *MockDoctorServiceClient) GetDoctorByID(ctx context.Context, doctorID string) (*DoctorServiceResponse, error) {
	return &DoctorServiceResponse{
		ID:         doctorID,
		Name:       "Dr. Mock",
		Specialty:  "Internal Medicine",
		Hospital:   "Mock Hospital",
		Department: "Internal Medicine",
		Title:      "Senior Physician",
		Experience: 10,
		Rating:     4.5,
	}, nil
}

func (m *MockDoctorServiceClient) SearchDoctors(ctx context.Context, specialty string, limit int) ([]*DoctorServiceResponse, error) {
	return []*DoctorServiceResponse{
		{
			ID:         "doctor-1",
			Name:       "Dr. Smith",
			Specialty:  specialty,
			Hospital:   "City Hospital",
			Department: specialty,
			Title:      "Chief Physician",
			Experience: 15,
			Rating:     4.8,
		},
	}, nil
}

func (m *MockDoctorServiceClient) GetDoctorAvailability(ctx context.Context, doctorID string) (*AvailabilityResponse, error) {
	return &AvailabilityResponse{
		AvailableSlots: []string{"Monday 10:00", "Tuesday 14:00"},
		Status:         "Available",
	}, nil
}