// Package vector provides Qdrant vector database repository implementation.
package vector

import (
	"context"
	"errors"
	"fmt"

	"intelligent_guidance_system_v2/service/ai/internal/domain/entity"
	"intelligent_guidance_system_v2/service/ai/internal/domain/repository"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/qdrant/go-client/qdrant"
)

var (
	ErrQdrantConnection  = errors.New("qdrant connection error")
	ErrQdrantSearch      = errors.New("qdrant search failed")
	ErrQdrantUpsert      = errors.New("qdrant upsert failed")
	ErrQdrantDelete      = errors.New("qdrant delete failed")
	ErrCollectionNotSet  = errors.New("collection not configured")
)

type QdrantConfig struct {
	Host            string
	Port            int
	DiseaseCollection string
	DoctorCollection  string
	DrugCollection    string
}

type QdrantRepo struct {
	client     *qdrant.Client
	config     QdrantConfig
	logger     *log.Helper
}

func NewQdrantRepo(cfg QdrantConfig, logger log.Logger) (repository.VectorRepository, error) {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: cfg.Host,
		Port: cfg.Port,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQdrantConnection, err)
	}
	
	return &QdrantRepo{
		client: client,
		config: cfg,
		logger: log.NewHelper(logger),
	}, nil
}

func (r *QdrantRepo) SearchDiseases(ctx context.Context, queryVector []float32, limit int, threshold float32) ([]*entity.VectorSearchResult, error) {
	if r.config.DiseaseCollection == "" {
		return nil, ErrCollectionNotSet
	}
	
	return r.search(ctx, r.config.DiseaseCollection, queryVector, limit, threshold)
}

func (r *QdrantRepo) SearchDoctors(ctx context.Context, queryVector []float32, limit int, threshold float32) ([]*entity.VectorSearchResult, error) {
	if r.config.DoctorCollection == "" {
		return nil, ErrCollectionNotSet
	}
	
	return r.search(ctx, r.config.DoctorCollection, queryVector, limit, threshold)
}

func (r *QdrantRepo) SearchDrugs(ctx context.Context, queryVector []float32, limit int, threshold float32) ([]*entity.VectorSearchResult, error) {
	if r.config.DrugCollection == "" {
		return nil, ErrCollectionNotSet
	}
	
	return r.search(ctx, r.config.DrugCollection, queryVector, limit, threshold)
}

func (r *QdrantRepo) UpsertVector(ctx context.Context, collection entity.VectorCollection, id string, vector []float32, payload map[string]interface{}) error {
	collectionName := r.getCollectionName(collection)
	if collectionName == "" {
		return ErrCollectionNotSet
	}
	
	point := &qdrant.PointStruct{
		Id:      qdrant.NewIDString(id),
		Vector:  qdrant.NewVectorFloat(vector),
		Payload: r.mapPayload(payload),
	}
	
	_, err := r.client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: collectionName,
		Points:         []*qdrant.PointStruct{point},
	})
	
	if err != nil {
		r.logger.Errorf("Failed to upsert vector %s: %v", id, err)
		return fmt.Errorf("%w: %v", ErrQdrantUpsert, err)
	}
	
	return nil
}

func (r *QdrantRepo) DeleteVector(ctx context.Context, collection entity.VectorCollection, id string) error {
	collectionName := r.getCollectionName(collection)
	if collectionName == "" {
		return ErrCollectionNotSet
	}
	
	_, err := r.client.Delete(ctx, &qdrant.DeletePoints{
		CollectionName: collectionName,
		PointsSelector: qdrant.NewPointsSelector(qdrant.NewIDString(id)),
	})
	
	if err != nil {
		r.logger.Errorf("Failed to delete vector %s: %v", id, err)
		return fmt.Errorf("%w: %v", ErrQdrantDelete, err)
	}
	
	return nil
}

func (r *QdrantRepo) search(ctx context.Context, collection string, vector []float32, limit int, threshold float32) ([]*entity.VectorSearchResult, error) {
	results, err := r.client.Query(ctx, &qdrant.QueryPoints{
		CollectionName: collection,
		Vector:         qdrant.NewVectorFloat(vector),
		Limit:          qdrant.Ptr(uint64(limit)),
		ScoreThreshold: qdrant.Ptr(float64(threshold)),
		WithPayload:    qdrant.NewWithPayloadSelector(true),
	})
	
	if err != nil {
		r.logger.Errorf("Failed to search collection %s: %v", collection, err)
		return nil, fmt.Errorf("%w: %v", ErrQdrantSearch, err)
	}
	
	return r.mapResults(results), nil
}

func (r *QdrantRepo) getCollectionName(collection entity.VectorCollection) string {
	switch collection {
	case entity.CollectionDisease:
		return r.config.DiseaseCollection
	case entity.CollectionDoctor:
		return r.config.DoctorCollection
	case entity.CollectionDrug:
		return r.config.DrugCollection
	default:
		return ""
	}
}

func (r *QdrantRepo) mapResults(results []*qdrant.ScoredPoint) []*entity.VectorSearchResult {
	mapped := make([]*entity.VectorSearchResult, len(results))
	
	for i, point := range results {
		id := ""
		if point.Id.GetString() != nil {
			id = *point.Id.GetString()
		}
		
		payload := make(map[string]interface{})
		if point.Payload != nil {
			for k, v := range point.Payload {
				payload[k] = r.extractPayloadValue(v)
			}
		}
		
		mapped[i] = &entity.VectorSearchResult{
			ID:      id,
			Score:   float32(point.Score),
			Payload: payload,
		}
	}
	
	return mapped
}

func (r *QdrantRepo) mapPayload(payload map[string]interface{}) map[string]*qdrant.Value {
	result := make(map[string]*qdrant.Value)
	
	for k, v := range payload {
		result[k] = r.toQdrantValue(v)
	}
	
	return result
}

func (r *QdrantRepo) toQdrantValue(v interface{}) *qdrant.Value {
	switch val := v.(type) {
	case string:
		return qdrant.NewValueString(val)
	case int:
		return qdrant.NewValueInt(int64(val))
	case int64:
		return qdrant.NewValueInt(val)
	case float32:
		return qdrant.NewValueDouble(float64(val))
	case float64:
		return qdrant.NewValueDouble(val)
	case bool:
		return qdrant.NewValueBool(val)
	default:
		return qdrant.NewValueString(fmt.Sprintf("%v", v))
	}
}

func (r *QdrantRepo) extractPayloadValue(v *qdrant.Value) interface{} {
	if v.GetStringValue() != nil {
		return *v.GetStringValue()
	}
	if v.GetIntegerValue() != nil {
		return *v.GetIntegerValue()
	}
	if v.GetDoubleValue() != nil {
		return *v.GetDoubleValue()
	}
	if v.GetBoolValue() != nil {
		return *v.GetBoolValue()
	}
	return nil
}

func (r *QdrantRepo) Close() error {
	return r.client.Close()
}