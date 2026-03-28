package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrEmptyDiseaseName = errors.New("disease name cannot be empty")
	ErrEmptySymptoms    = errors.New("symptoms cannot be empty")
)

// DiseaseItem represents a disease diagnosis item entity
type DiseaseItem struct {
	id           string    // unique identifier for this disease item
	diseaseID    int64     // reference to disease catalog
	diseaseName  string    // disease name
	symptoms     string    // symptoms description
	diagnosis    string    // diagnosis explanation
	createdAt    time.Time // creation timestamp
}

// NewDiseaseItem creates a new DiseaseItem entity
func NewDiseaseItem(diseaseID int64, diseaseName, symptoms, diagnosis string) (*DiseaseItem, error) {
	if diseaseName == "" {
		return nil, ErrEmptyDiseaseName
	}
	if symptoms == "" {
		return nil, ErrEmptySymptoms
	}
	return &DiseaseItem{
		id:          uuid.New().String(),
		diseaseID:   diseaseID,
		diseaseName: diseaseName,
		symptoms:    symptoms,
		diagnosis:   diagnosis,
		createdAt:   time.Now(),
	}, nil
}

// ReconstructDiseaseItem reconstructs a DiseaseItem from persistence
func ReconstructDiseaseItem(id string, diseaseID int64, diseaseName, symptoms, diagnosis string, createdAt time.Time) *DiseaseItem {
	return &DiseaseItem{
		id:          id,
		diseaseID:   diseaseID,
		diseaseName: diseaseName,
		symptoms:    symptoms,
		diagnosis:   diagnosis,
		createdAt:   createdAt,
	}
}

// ID returns the disease item ID
func (d *DiseaseItem) ID() string {
	return d.id
}

// DiseaseID returns the disease catalog ID
func (d *DiseaseItem) DiseaseID() int64 {
	return d.diseaseID
}

// DiseaseName returns the disease name
func (d *DiseaseItem) DiseaseName() string {
	return d.diseaseName
}

// Symptoms returns the symptoms
func (d *DiseaseItem) Symptoms() string {
	return d.symptoms
}

// Diagnosis returns the diagnosis
func (d *DiseaseItem) Diagnosis() string {
	return d.diagnosis
}

// CreatedAt returns the creation timestamp
func (d *DiseaseItem) CreatedAt() time.Time {
	return d.createdAt
}

// UpdateDiagnosis updates the diagnosis explanation
func (d *DiseaseItem) UpdateDiagnosis(diagnosis string) {
	d.diagnosis = diagnosis
}

// UpdateSymptoms updates the symptoms description
func (d *DiseaseItem) UpdateSymptoms(symptoms string) error {
	if symptoms == "" {
		return ErrEmptySymptoms
	}
	d.symptoms = symptoms
	return nil
}