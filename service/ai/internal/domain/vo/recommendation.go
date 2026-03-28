// Package vo provides value objects for the AI diagnosis service.
// Value objects are immutable and defined by their attributes.
package vo

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// RecommendationType represents the type of medical recommendation.
type RecommendationType string

const (
	// RecommendationTypeDisease suggests possible diseases
	RecommendationTypeDisease RecommendationType = "disease"
	
	// RecommendationTypeDoctor suggests consulting specific doctors
	RecommendationTypeDoctor RecommendationType = "doctor"
	
	// RecommendationTypeDrug suggests relevant medications
	RecommendationTypeDrug RecommendationType = "drug"
	
	// RecommendationTypeDepartment suggests relevant departments
	RecommendationTypeDepartment RecommendationType = "department"
)

// ConfidenceLevel represents the confidence of a recommendation.
type ConfidenceLevel string

const (
	// ConfidenceHigh indicates high confidence (>= 0.8)
	ConfidenceHigh ConfidenceLevel = "high"
	
	// ConfidenceMedium indicates medium confidence (0.5 - 0.8)
	ConfidenceMedium ConfidenceLevel = "medium"
	
	// ConfidenceLow indicates low confidence (< 0.5)
	ConfidenceLow ConfidenceLevel = "low"
)

// Recommendation is a value object representing a medical recommendation.
// It is immutable and created fresh for each diagnosis result.
type Recommendation struct {
	// ID is the unique identifier
	id string
	
	// Type indicates what kind of recommendation this is
	recType RecommendationType
	
	// Title is the main recommendation text
	title string
	
	// Description provides detailed explanation
	description string
	
	// ConfidenceScore is the AI confidence (0-1)
	confidenceScore float32
	
	// SourceIDs are the vector search source IDs
	sourceIDs []string
	
	// Reasoning explains why this recommendation was made
	reasoning string
	
	// Metadata contains additional structured data
	metadata map[string]interface{}
	
	// CreatedAt is when the recommendation was generated
	createdAt time.Time
}

// Error definitions for Recommendation VO
var (
	ErrInvalidRecommendationType = errors.New("invalid recommendation type")
	ErrInvalidConfidenceScore   = errors.New("confidence score must be between 0 and 1")
	ErrEmptyTitle                = errors.New("recommendation title cannot be empty")
)

// NewRecommendation creates a new Recommendation value object.
func NewRecommendation(
	recType RecommendationType,
	title string,
	description string,
	confidenceScore float32,
) (*Recommendation, error) {
	if !isValidRecommendationType(recType) {
		return nil, ErrInvalidRecommendationType
	}
	if title == "" {
		return nil, ErrEmptyTitle
	}
	if confidenceScore < 0 || confidenceScore > 1 {
		return nil, ErrInvalidConfidenceScore
	}

	return &Recommendation{
		id:              uuid.New().String(),
		recType:         recType,
		title:           title,
		description:     description,
		confidenceScore: confidenceScore,
		sourceIDs:       make([]string, 0),
		metadata:        make(map[string]interface{}),
		createdAt:       time.Now(),
	}, nil
}

func isValidRecommendationType(t RecommendationType) bool {
	switch t {
	case RecommendationTypeDisease, RecommendationTypeDoctor,
		RecommendationTypeDrug, RecommendationTypeDepartment:
		return true
	default:
		return false
	}
}

// ID returns the recommendation ID.
func (r *Recommendation) ID() string {
	return r.id
}

// Type returns the recommendation type.
func (r *Recommendation) Type() RecommendationType {
	return r.recType
}

// Title returns the recommendation title.
func (r *Recommendation) Title() string {
	return r.title
}

// Description returns the recommendation description.
func (r *Recommendation) Description() string {
	return r.description
}

// ConfidenceScore returns the confidence score.
func (r *Recommendation) ConfidenceScore() float32 {
	return r.confidenceScore
}

// ConfidenceLevel returns the confidence level based on score.
func (r *Recommendation) ConfidenceLevel() ConfidenceLevel {
	if r.confidenceScore >= 0.8 {
		return ConfidenceHigh
	}
	if r.confidenceScore >= 0.5 {
		return ConfidenceMedium
	}
	return ConfidenceLow
}

// SourceIDs returns a copy of source IDs.
func (r *Recommendation) SourceIDs() []string {
	ids := make([]string, len(r.sourceIDs))
	copy(ids, r.sourceIDs)
	return ids
}

// WithSourceIDs returns a new Recommendation with source IDs set.
func (r *Recommendation) WithSourceIDs(ids []string) *Recommendation {
	newIDs := make([]string, len(ids))
	copy(newIDs, ids)
	
	return &Recommendation{
		id:              r.id,
		recType:         r.recType,
		title:           r.title,
		description:     r.description,
		confidenceScore: r.confidenceScore,
		sourceIDs:       newIDs,
		reasoning:       r.reasoning,
		metadata:        copyMetadata(r.metadata),
		createdAt:       r.createdAt,
	}
}

// Reasoning returns the reasoning text.
func (r *Recommendation) Reasoning() string {
	return r.reasoning
}

// WithReasoning returns a new Recommendation with reasoning set.
func (r *Recommendation) WithReasoning(reasoning string) *Recommendation {
	return &Recommendation{
		id:              r.id,
		recType:         r.recType,
		title:           r.title,
		description:     r.description,
		confidenceScore: r.confidenceScore,
		sourceIDs:       r.sourceIDs,
		reasoning:       reasoning,
		metadata:        copyMetadata(r.metadata),
		createdAt:       r.createdAt,
	}
}

// Metadata returns a copy of the metadata.
func (r *Recommendation) Metadata() map[string]interface{} {
	return copyMetadata(r.metadata)
}

// WithMetadata returns a new Recommendation with metadata added.
func (r *Recommendation) WithMetadata(key string, value interface{}) *Recommendation {
	newMetadata := copyMetadata(r.metadata)
	newMetadata[key] = value
	
	return &Recommendation{
		id:              r.id,
		recType:         r.recType,
		title:           r.title,
		description:     r.description,
		confidenceScore: r.confidenceScore,
		sourceIDs:       r.sourceIDs,
		reasoning:       r.reasoning,
		metadata:        newMetadata,
		createdAt:       r.createdAt,
	}
}

// CreatedAt returns the creation timestamp.
func (r *Recommendation) CreatedAt() time.Time {
	return r.createdAt
}

// copyMetadata creates a deep copy of metadata map.
func copyMetadata(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(m))
	for k, v := range m {
		result[k] = v
	}
	return result
}

// RecommendationSet is a collection of recommendations for a diagnosis session.
type RecommendationSet struct {
	diseaseRecommendations     []*Recommendation
	doctorRecommendations      []*Recommendation
	drugRecommendations        []*Recommendation
	departmentRecommendations []*Recommendation
}

// NewRecommendationSet creates an empty recommendation set.
func NewRecommendationSet() *RecommendationSet {
	return &RecommendationSet{
		diseaseRecommendations:     make([]*Recommendation, 0),
		doctorRecommendations:      make([]*Recommendation, 0),
		drugRecommendations:        make([]*Recommendation, 0),
		departmentRecommendations: make([]*Recommendation, 0),
	}
}

// AddDisease adds a disease recommendation.
func (rs *RecommendationSet) AddDisease(r *Recommendation) {
	rs.diseaseRecommendations = append(rs.diseaseRecommendations, r)
}

// AddDoctor adds a doctor recommendation.
func (rs *RecommendationSet) AddDoctor(r *Recommendation) {
	rs.doctorRecommendations = append(rs.doctorRecommendations, r)
}

// AddDrug adds a drug recommendation.
func (rs *RecommendationSet) AddDrug(r *Recommendation) {
	rs.drugRecommendations = append(rs.drugRecommendations, r)
}

// AddDepartment adds a department recommendation.
func (rs *RecommendationSet) AddDepartment(r *Recommendation) {
	rs.departmentRecommendations = append(rs.departmentRecommendations, r)
}

// Diseases returns disease recommendations.
func (rs *RecommendationSet) Diseases() []*Recommendation {
	return rs.diseaseRecommendations
}

// Doctors returns doctor recommendations.
func (rs *RecommendationSet) Doctors() []*Recommendation {
	return rs.doctorRecommendations
}

// Drugs returns drug recommendations.
func (rs *RecommendationSet) Drugs() []*Recommendation {
	return rs.drugRecommendations
}

// Departments returns department recommendations.
func (rs *RecommendationSet) Departments() []*Recommendation {
	return rs.departmentRecommendations
}

// IsEmpty returns true if there are no recommendations.
func (rs *RecommendationSet) IsEmpty() bool {
	return len(rs.diseaseRecommendations) == 0 &&
		len(rs.doctorRecommendations) == 0 &&
		len(rs.drugRecommendations) == 0 &&
		len(rs.departmentRecommendations) == 0
}

// TopDiseases returns the top N disease recommendations by confidence.
func (rs *RecommendationSet) TopDiseases(n int) []*Recommendation {
	return topByConfidence(rs.diseaseRecommendations, n)
}

// TopDoctors returns the top N doctor recommendations by confidence.
func (rs *RecommendationSet) TopDoctors(n int) []*Recommendation {
	return topByConfidence(rs.doctorRecommendations, n)
}

func topByConfidence(recs []*Recommendation, n int) []*Recommendation {
	if n >= len(recs) {
		return recs
	}
	// Simple selection - in production would use proper sorting
	result := make([]*Recommendation, n)
	copy(result, recs[:n])
	return result
}