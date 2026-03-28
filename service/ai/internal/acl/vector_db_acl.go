package acl

import (
	"context"
	"errors"

	"intelligent_guidance_system_v2/service/ai/internal/domain/entity"
	"intelligent_guidance_system_v2/service/ai/internal/domain/repository"

	"github.com/go-kratos/kratos/v2/log"
)

var (
	ErrVectorDBUnavailable = errors.New("vector database unavailable")
	ErrVectorSearchFailed  = errors.New("vector search failed")
)

type VectorSearchResult struct {
	ID      string
	Score   float32
	Payload map[string]interface{}
}

type VectorDBACL struct {
	vectorRepo    repository.VectorRepository
	embeddingRepo repository.EmbeddingRepository
	logger        *log.Helper
}

func NewVectorDBACL(
	vectorRepo repository.VectorRepository,
	embeddingRepo repository.EmbeddingRepository,
	logger log.Logger,
) *VectorDBACL {
	return &VectorDBACL{
		vectorRepo:    vectorRepo,
		embeddingRepo: embeddingRepo,
		logger:        log.NewHelper(logger),
	}
}

func (acl *VectorDBACL) SearchDiseasesByText(ctx context.Context, queryText string, limit int, threshold float32) ([]*VectorSearchResult, error) {
	embedding, err := acl.embeddingRepo.Embed(ctx, queryText)
	if err != nil {
		acl.logger.Warnf("Failed to embed query text: %v", err)
		return nil, ErrVectorSearchFailed
	}

	results, err := acl.vectorRepo.SearchDiseases(ctx, embedding, limit, threshold)
	if err != nil {
		acl.logger.Warnf("Failed to search diseases: %v", err)
		return nil, ErrVectorSearchFailed
	}

	return acl.mapResults(results), nil
}

func (acl *VectorDBACL) SearchDoctorsByText(ctx context.Context, queryText string, limit int, threshold float32) ([]*VectorSearchResult, error) {
	embedding, err := acl.embeddingRepo.Embed(ctx, queryText)
	if err != nil {
		acl.logger.Warnf("Failed to embed query text: %v", err)
		return nil, ErrVectorSearchFailed
	}

	results, err := acl.vectorRepo.SearchDoctors(ctx, embedding, limit, threshold)
	if err != nil {
		acl.logger.Warnf("Failed to search doctors: %v", err)
		return nil, ErrVectorSearchFailed
	}

	return acl.mapResults(results), nil
}

func (acl *VectorDBACL) SearchDrugsByText(ctx context.Context, queryText string, limit int, threshold float32) ([]*VectorSearchResult, error) {
	embedding, err := acl.embeddingRepo.Embed(ctx, queryText)
	if err != nil {
		acl.logger.Warnf("Failed to embed query text: %v", err)
		return nil, ErrVectorSearchFailed
	}

	results, err := acl.vectorRepo.SearchDrugs(ctx, embedding, limit, threshold)
	if err != nil {
		acl.logger.Warnf("Failed to search drugs: %v", err)
		return nil, ErrVectorSearchFailed
	}

	return acl.mapResults(results), nil
}

func (acl *VectorDBACL) SearchByVector(ctx context.Context, vector []float32, collection string, limit int, threshold float32) ([]*VectorSearchResult, error) {
	var results []*entity.VectorSearchResult
	var err error

	switch collection {
	case "diseases":
		results, err = acl.vectorRepo.SearchDiseases(ctx, vector, limit, threshold)
	case "doctors":
		results, err = acl.vectorRepo.SearchDoctors(ctx, vector, limit, threshold)
	case "drugs":
		results, err = acl.vectorRepo.SearchDrugs(ctx, vector, limit, threshold)
	default:
		return nil, ErrVectorSearchFailed
	}

	if err != nil {
		return nil, ErrVectorSearchFailed
	}

	return acl.mapResults(results), nil
}

func (acl *VectorDBACL) UpsertDisease(ctx context.Context, id string, name string, description string, symptoms []string) error {
	vectorText := name + " " + description + " " + joinStrings(symptoms)
	
	embedding, err := acl.embeddingRepo.Embed(ctx, vectorText)
	if err != nil {
		return err
	}

	payload := map[string]interface{}{
		"name":        name,
		"description": description,
		"symptoms":    symptoms,
	}

	return acl.vectorRepo.UpsertVector(ctx, entity.CollectionDisease, id, embedding, payload)
}

func (acl *VectorDBACL) UpsertDoctor(ctx context.Context, id string, name string, specialty string, hospital string) error {
	vectorText := name + " " + specialty + " " + hospital
	
	embedding, err := acl.embeddingRepo.Embed(ctx, vectorText)
	if err != nil {
		return err
	}

	payload := map[string]interface{}{
		"name":      name,
		"specialty": specialty,
		"hospital":  hospital,
	}

	return acl.vectorRepo.UpsertVector(ctx, entity.CollectionDoctor, id, embedding, payload)
}

func (acl *VectorDBACL) UpsertDrug(ctx context.Context, id string, name string, indication string, contraindications []string) error {
	vectorText := name + " " + indication + " " + joinStrings(contraindications)
	
	embedding, err := acl.embeddingRepo.Embed(ctx, vectorText)
	if err != nil {
		return err
	}

	payload := map[string]interface{}{
		"name":            name,
		"indication":      indication,
		"contraindications": contraindications,
	}

	return acl.vectorRepo.UpsertVector(ctx, entity.CollectionDrug, id, embedding, payload)
}

func (acl *VectorDBACL) mapResults(results []*entity.VectorSearchResult) []*VectorSearchResult {
	mapped := make([]*VectorSearchResult, len(results))
	for i, r := range results {
		mapped[i] = &VectorSearchResult{
			ID:      r.ID,
			Score:   r.Score,
			Payload: r.Payload,
		}
	}
	return mapped
}

func joinStrings(strs []string) string {
	result := ""
	for _, s := range strs {
		result += s + " "
	}
	return result
}

func (acl *VectorDBACL) GetTopResults(results []*VectorSearchResult, n int) []*VectorSearchResult {
	if n >= len(results) {
		return results
	}
	return results[:n]
}

func (acl *VectorDBACL) FilterByScore(results []*VectorSearchResult, threshold float32) []*VectorSearchResult {
	filtered := make([]*VectorSearchResult, 0)
	for _, r := range results {
		if r.Score >= threshold {
			filtered = append(filtered, r)
		}
	}
	return filtered
}