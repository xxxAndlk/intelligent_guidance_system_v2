// Package repository defines repository interfaces for the domain layer.
// Repositories provide abstract access to persistent storage.
package repository

import (
	"context"

	"intelligent_guidance_system_v2/service/ai/internal/domain/aggregate"
	"intelligent_guidance_system_v2/service/ai/internal/domain/entity"
)

// DiagnosisRepository defines the interface for diagnosis session persistence.
// This is the aggregate repository for the DiagnosisSession aggregate root.
type DiagnosisRepository interface {
	// Save persists a diagnosis session (create or update)
	Save(ctx context.Context, session *aggregate.DiagnosisSession) error
	
	// FindByID retrieves a session by its ID
	FindByID(ctx context.Context, sessionID string) (*aggregate.DiagnosisSession, error)
	
	// FindByPatientID retrieves all sessions for a patient
	FindByPatientID(ctx context.Context, patientID string, limit, offset int) ([]*aggregate.DiagnosisSession, error)
	
	// FindActiveByPatientID retrieves the active session for a patient
	FindActiveByPatientID(ctx context.Context, patientID string) (*aggregate.DiagnosisSession, error)
	
	// Delete removes a session
	Delete(ctx context.Context, sessionID string) error
	
	// UpdateStatus updates the session status
	UpdateStatus(ctx context.Context, sessionID string, status aggregate.SessionStatus) error
}

// SymptomRepository defines the interface for symptom persistence.
type SymptomRepository interface {
	// Save persists a symptom
	Save(ctx context.Context, symptom *entity.Symptom) error
	
	// FindByID retrieves a symptom by ID
	FindByID(ctx context.Context, symptomID string) (*entity.Symptom, error)
	
	// FindBySessionID retrieves all symptoms for a session
	FindBySessionID(ctx context.Context, sessionID string) ([]*entity.Symptom, error)
	
	// Delete removes a symptom
	Delete(ctx context.Context, symptomID string) error
}

// ConversationRepository defines the interface for conversation persistence.
type ConversationRepository interface {
	// Save persists a conversation message
	Save(ctx context.Context, msg *entity.Conversation) error
	
	// FindBySessionID retrieves all messages for a session in order
	FindBySessionID(ctx context.Context, sessionID string) ([]*entity.Conversation, error)
	
	// FindBySessionIDPaginated retrieves messages with pagination
	FindBySessionIDPaginated(ctx context.Context, sessionID string, limit, offset int) ([]*entity.Conversation, error)
	
	// GetLatestBySessionID retrieves the latest message for a session
	GetLatestBySessionID(ctx context.Context, sessionID string) (*entity.Conversation, error)
	
	// Delete removes a message
	Delete(ctx context.Context, messageID string) error
}

// VectorRepository defines the interface for vector database operations.
type VectorRepository interface {
	// SearchDiseases performs vector similarity search for diseases
	SearchDiseases(ctx context.Context, queryVector []float32, limit int, threshold float32) ([]*entity.VectorSearchResult, error)
	
	// SearchDoctors performs vector similarity search for doctors
	SearchDoctors(ctx context.Context, queryVector []float32, limit int, threshold float32) ([]*entity.VectorSearchResult, error)
	
	// SearchDrugs performs vector similarity search for drugs
	SearchDrugs(ctx context.Context, queryVector []float32, limit int, threshold float32) ([]*entity.VectorSearchResult, error)
	
	// UpsertVector inserts or updates a vector in the specified collection
	UpsertVector(ctx context.Context, collection entity.VectorCollection, id string, vector []float32, payload map[string]interface{}) error
	
	// DeleteVector removes a vector from the specified collection
	DeleteVector(ctx context.Context, collection entity.VectorCollection, id string) error
}

// EmbeddingRepository defines the interface for embedding generation.
type EmbeddingRepository interface {
	// Embed generates embeddings for the given text
	Embed(ctx context.Context, text string) ([]float32, error)
	
	// EmbedBatch generates embeddings for multiple texts
	EmbedBatch(ctx context.Context, texts []string) ([][]float32, error)
}