package acl

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/log"
)

var (
	ErrDiseaseServiceUnavailable = errors.New("disease service unavailable")
	ErrDiseaseNotFound           = errors.New("disease not found")
)

type DiseaseInfo struct {
	ID              string
	Name            string
	Category        string
	Symptoms        []string
	Description     string
	Treatment       string
	Severity        string
	CommonAgeGroup  string
	RelatedDiseases []string
}

type DiseaseACL struct {
	diseaseServiceClient DiseaseServiceClient
	logger               *log.Helper
}

type DiseaseServiceClient interface {
	GetDiseaseByID(ctx context.Context, diseaseID string) (*DiseaseServiceResponse, error)
	SearchDiseases(ctx context.Context, symptoms []string, limit int) ([]*DiseaseServiceResponse, error)
	GetDiseaseDetails(ctx context.Context, diseaseID string) (*DiseaseDetailResponse, error)
}

type DiseaseServiceResponse struct {
	ID       string
	Name     string
	Category string
}

type DiseaseDetailResponse struct {
	ID              string
	Name            string
	Category        string
	Symptoms        []string
	Description     string
	Treatment       string
	Severity        string
	CommonAgeGroup  string
	RelatedDiseases []string
}

func NewDiseaseACL(client DiseaseServiceClient, logger log.Logger) *DiseaseACL {
	return &DiseaseACL{
		diseaseServiceClient: client,
		logger:               log.NewHelper(logger),
	}
}

func (acl *DiseaseACL) GetDiseaseInfo(ctx context.Context, diseaseID string) (*DiseaseInfo, error) {
	resp, err := acl.diseaseServiceClient.GetDiseaseDetails(ctx, diseaseID)
	if err != nil {
		acl.logger.Warnf("Failed to get disease %s: %v", diseaseID, err)
		return nil, ErrDiseaseServiceUnavailable
	}

	if resp == nil {
		return nil, ErrDiseaseNotFound
	}

	return &DiseaseInfo{
		ID:              resp.ID,
		Name:            resp.Name,
		Category:        resp.Category,
		Symptoms:        resp.Symptoms,
		Description:     resp.Description,
		Treatment:       resp.Treatment,
		Severity:        resp.Severity,
		CommonAgeGroup:  resp.CommonAgeGroup,
		RelatedDiseases: resp.RelatedDiseases,
	}, nil
}

func (acl *DiseaseACL) SearchDiseasesBySymptoms(ctx context.Context, symptoms []string, limit int) ([]*DiseaseInfo, error) {
	resps, err := acl.diseaseServiceClient.SearchDiseases(ctx, symptoms, limit)
	if err != nil {
		acl.logger.Warnf("Failed to search diseases: %v", err)
		return nil, ErrDiseaseServiceUnavailable
	}

	diseases := make([]*DiseaseInfo, len(resps))
	for i, resp := range resps {
		detail, err := acl.diseaseServiceClient.GetDiseaseDetails(ctx, resp.ID)
		if err != nil {
			acl.logger.Warnf("Failed to get details for disease %s: %v", resp.ID, err)
			diseases[i] = &DiseaseInfo{
				ID:       resp.ID,
				Name:     resp.Name,
				Category: resp.Category,
			}
			continue
		}

		diseases[i] = &DiseaseInfo{
			ID:              detail.ID,
			Name:            detail.Name,
			Category:        detail.Category,
			Symptoms:        detail.Symptoms,
			Description:     detail.Description,
			Treatment:       detail.Treatment,
			Severity:        detail.Severity,
			CommonAgeGroup:  detail.CommonAgeGroup,
			RelatedDiseases: detail.RelatedDiseases,
		}
	}

	return diseases, nil
}

func (acl *DiseaseACL) FormatDiseaseRecommendation(disease *DiseaseInfo) string {
	return disease.Name + " (" + disease.Category + ") - Severity: " + disease.Severity
}

func (acl *DiseaseACL) GetRelatedDiseaseContext(ctx context.Context, diseaseID string) ([]string, error) {
	info, err := acl.GetDiseaseInfo(ctx, diseaseID)
	if err != nil {
		return nil, err
	}

	return info.RelatedDiseases, nil
}

type MockDiseaseServiceClient struct{}

func (m *MockDiseaseServiceClient) GetDiseaseByID(ctx context.Context, diseaseID string) (*DiseaseServiceResponse, error) {
	return &DiseaseServiceResponse{
		ID:       diseaseID,
		Name:     "Mock Disease",
		Category: "General",
	}, nil
}

func (m *MockDiseaseServiceClient) SearchDiseases(ctx context.Context, symptoms []string, limit int) ([]*DiseaseServiceResponse, error) {
	return []*DiseaseServiceResponse{
		{
			ID:       "disease-1",
			Name:     "Common Cold",
			Category: "Respiratory",
		},
	}, nil
}

func (m *MockDiseaseServiceClient) GetDiseaseDetails(ctx context.Context, diseaseID string) (*DiseaseDetailResponse, error) {
	return &DiseaseDetailResponse{
		ID:              diseaseID,
		Name:            "Common Cold",
		Category:        "Respiratory",
		Symptoms:        []string{"cough", "runny nose", "sore throat"},
		Description:     "A viral infection of the upper respiratory tract",
		Treatment:       "Rest, hydration, over-the-counter medications",
		Severity:        "Low",
		CommonAgeGroup:  "All ages",
		RelatedDiseases: []string{"Flu", "Bronchitis"},
	}, nil
}