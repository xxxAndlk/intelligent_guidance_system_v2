// Package entity provides domain entities for the AI diagnosis service.
// Entities have identity and lifecycle, representing core domain concepts.
package entity

import (
	"time"

	"github.com/google/uuid"
)

// Symptom represents a patient's symptom entity within a diagnosis session.
// It captures the symptom details including location, description, and severity.
type Symptom struct {
	// ID is the unique identifier for the symptom
	ID string `json:"id" gorm:"primaryKey;type:varchar(36)"`
	
	// SessionID references the parent diagnosis session
	SessionID string `json:"session_id" gorm:"index;type:varchar(36)"`
	
	// BodyPart indicates the anatomical location of the symptom
	BodyPart string `json:"body_part" gorm:"type:varchar(100)"`
	
	// Description provides detailed symptom description from patient
	Description string `json:"description" gorm:"type:text"`
	
	// Severity indicates symptom intensity (1-10 scale)
	Severity int `json:"severity" gorm:"type:int"`
	
	// Duration indicates how long the symptom has persisted
	Duration string `json:"duration" gorm:"type:varchar(50)"`
	
	// Metadata stores additional symptom attributes in JSON format
	Metadata map[string]interface{} `json:"metadata" gorm:"type:json"`
	
	// CreatedAt is the timestamp when symptom was recorded
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	
	// UpdatedAt is the timestamp of last modification
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// NewSymptom creates a new Symptom entity with validation.
// bodyPart: the anatomical location
// description: detailed symptom description
// severity: intensity level from 1 to 10
// Returns error if validation fails.
func NewSymptom(sessionID, bodyPart, description string, severity int) (*Symptom, error) {
	if bodyPart == "" {
		return nil, ErrInvalidBodyPart
	}
	if description == "" {
		return nil, ErrInvalidDescription
	}
	if severity < 1 || severity > 10 {
		return nil, ErrInvalidSeverity
	}

	return &Symptom{
		ID:          uuid.New().String(),
		SessionID:   sessionID,
		BodyPart:    bodyPart,
		Description: description,
		Severity:    severity,
		Metadata:    make(map[string]interface{}),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

// SetDuration sets the symptom duration.
func (s *Symptom) SetDuration(duration string) {
	s.Duration = duration
	s.UpdatedAt = time.Now()
}

// AddMetadata adds a key-value pair to the symptom metadata.
func (s *Symptom) AddMetadata(key string, value interface{}) {
	if s.Metadata == nil {
		s.Metadata = make(map[string]interface{})
	}
	s.Metadata[key] = value
	s.UpdatedAt = time.Now()
}

// ToVector returns the symptom as a vector-compatible string representation.
func (s *Symptom) ToVector() string {
	return s.BodyPart + " " + s.Description + " severity:" + string(rune(s.Severity))
}

// TableName returns the database table name for GORM.
func (Symptom) TableName() string {
	return "ai_symptoms"
}