// Package entity provides domain entities for the AI diagnosis service.
package entity

import (
	"time"

	"github.com/google/uuid"
)

// VectorCollection represents the type of vector collection.
type VectorCollection string

const (
	// CollectionDisease represents disease vectors
	CollectionDisease VectorCollection = "diseases"
	
	// CollectionDoctor represents doctor vectors
	CollectionDoctor VectorCollection = "doctors"
	
	// CollectionDrug represents drug vectors
	CollectionDrug VectorCollection = "drugs"
	
	// CollectionSymptom represents symptom vectors
	CollectionSymptom VectorCollection = "symptoms"
)

// VectorQuery represents a query for vector similarity search.
// It encapsulates the query parameters for vector database operations.
type VectorQuery struct {
	// ID is the unique identifier for the query
	ID string `json:"id"`
	
	// SessionID references the diagnosis session
	SessionID string `json:"session_id"`
	
	// QueryText is the original text query
	QueryText string `json:"query_text"`
	
	// QueryVector is the embedded vector representation
	QueryVector []float32 `json:"query_vector"`
	
	// Collection specifies which vector collection to search
	Collection VectorCollection `json:"collection"`
	
	// TopK is the number of results to return
	TopK int `json:"top_k"`
	
	// ScoreThreshold is the minimum similarity score (0-1)
	ScoreThreshold float32 `json:"score_threshold"`
	
	// FilterConditions are optional metadata filters
	FilterConditions map[string]interface{} `json:"filter_conditions,omitempty"`
	
	// Results stores the search results
	Results []VectorSearchResult `json:"results,omitempty"`
	
	// CreatedAt is when the query was executed
	CreatedAt time.Time `json:"created_at"`
}

// VectorSearchResult represents a single vector search result.
type VectorSearchResult struct {
	// ID is the vector ID in the collection
	ID string `json:"id"`
	
	// Score is the similarity score
	Score float32 `json:"score"`
	
	// Payload contains the associated data
	Payload map[string]interface{} `json:"payload"`
}

// NewVectorQuery creates a new VectorQuery entity.
func NewVectorQuery(sessionID, queryText string, collection VectorCollection, topK int) *VectorQuery {
	return &VectorQuery{
		ID:               uuid.New().String(),
		SessionID:        sessionID,
		QueryText:        queryText,
		Collection:       collection,
		TopK:             topK,
		ScoreThreshold:   0.7, // Default threshold
		FilterConditions: make(map[string]interface{}),
		Results:          make([]VectorSearchResult, 0),
		CreatedAt:        time.Now(),
	}
}

// SetVector sets the query vector embedding.
func (v *VectorQuery) SetVector(vector []float32) {
	v.QueryVector = vector
}

// SetScoreThreshold sets the minimum similarity score.
func (v *VectorQuery) SetScoreThreshold(threshold float32) {
	if threshold >= 0 && threshold <= 1 {
		v.ScoreThreshold = threshold
	}
}

// AddFilter adds a filter condition to the query.
func (v *VectorQuery) AddFilter(key string, value interface{}) {
	if v.FilterConditions == nil {
		v.FilterConditions = make(map[string]interface{})
	}
	v.FilterConditions[key] = value
}

// AddResult adds a search result to the query.
func (v *VectorQuery) AddResult(result VectorSearchResult) {
	v.Results = append(v.Results, result)
}

// GetTopResults returns the top N results by score.
func (v *VectorQuery) GetTopResults(n int) []VectorSearchResult {
	if n >= len(v.Results) {
		return v.Results
	}
	return v.Results[:n]
}

// FilterByScore filters results by minimum score threshold.
func (v *VectorQuery) FilterByScore() {
	filtered := make([]VectorSearchResult, 0)
	for _, r := range v.Results {
		if r.Score >= v.ScoreThreshold {
			filtered = append(filtered, r)
		}
	}
	v.Results = filtered
}

// HasResults returns true if there are search results.
func (v *VectorQuery) HasResults() bool {
	return len(v.Results) > 0
}

// BestResult returns the highest-scoring result.
func (v *VectorQuery) BestResult() *VectorSearchResult {
	if len(v.Results) == 0 {
		return nil
	}
	best := v.Results[0]
	for _, r := range v.Results[1:] {
		if r.Score > best.Score {
			best = r
		}
	}
	return &best
}