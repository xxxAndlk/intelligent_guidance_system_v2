// Package aggregate provides aggregate roots for the AI diagnosis service.
// Aggregates are clusters of domain objects treated as a single unit.
package aggregate

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"intelligent_guidance_system_v2/service/ai/internal/domain/entity"
	"intelligent_guidance_system_v2/service/ai/internal/domain/event"
	"intelligent_guidance_system_v2/service/ai/internal/domain/vo"

	"github.com/google/uuid"
)

// SessionStatus represents the current state of a diagnosis session.
type SessionStatus string

const (
	// StatusActive indicates the session is actively collecting symptoms
	StatusActive SessionStatus = "active"
	
	// StatusAnalyzing indicates the session is being analyzed by AI
	StatusAnalyzing SessionStatus = "analyzing"
	
	// StatusCompleted indicates the session has been completed
	StatusCompleted SessionStatus = "completed"
	
	// StatusArchived indicates the session has been archived
	StatusArchived SessionStatus = "archived"
)

// Error definitions for DiagnosisSession aggregate
var (
	ErrSessionNotActive      = errors.New("session is not active")
	ErrSessionAlreadyComplete = errors.New("session already completed")
	ErrSessionNotFound       = errors.New("session not found")
	ErrNoSymptoms            = errors.New("no symptoms recorded")
	ErrInvalidPatientID      = errors.New("invalid patient ID")
)

// DiagnosisSession is the aggregate root for the diagnosis domain.
// It manages the complete lifecycle of a medical diagnosis session including
// symptoms, conversations, and recommendations.
type DiagnosisSession struct {
	// ID is the unique session identifier
	id string
	
	// PatientID references the patient
	patientID string
	
	// Status is the current session status
	status SessionStatus
	
	// Symptoms is the collection of reported symptoms
	symptoms []*entity.Symptom
	
	// Conversations is the dialogue history
	conversations []*entity.Conversation
	
	// Recommendations is the generated recommendation set
	recommendations *vo.RecommendationSet
	
	// domainEvents stores uncommitted domain events
	domainEvents []event.DomainEvent
	
	// createdAt is when the session was created
	createdAt time.Time
	
	// updatedAt is when the session was last modified
	updatedAt time.Time
	
	// completedAt is when the session was completed
	completedAt *time.Time
}

// NewDiagnosisSession creates a new diagnosis session.
func NewDiagnosisSession(patientID string) (*DiagnosisSession, error) {
	if patientID == "" {
		return nil, ErrInvalidPatientID
	}

	now := time.Now()
	session := &DiagnosisSession{
		id:             uuid.New().String(),
		patientID:      patientID,
		status:         StatusActive,
		symptoms:       make([]*entity.Symptom, 0),
		conversations:  make([]*entity.Conversation, 0),
		recommendations: vo.NewRecommendationSet(),
		domainEvents:   make([]event.DomainEvent, 0),
		createdAt:      now,
		updatedAt:      now,
	}

	session.addEvent(event.NewSessionCreatedEvent(session.id, patientID))
	
	return session, nil
}

// Reconstruct creates a DiagnosisSession from persisted data.
// Used by repositories to reconstruct aggregates from storage.
func Reconstruct(
	id string,
	patientID string,
	status SessionStatus,
	symptoms []*entity.Symptom,
	conversations []*entity.Conversation,
	recommendations *vo.RecommendationSet,
	createdAt time.Time,
	updatedAt time.Time,
	completedAt *time.Time,
) *DiagnosisSession {
	return &DiagnosisSession{
		id:             id,
		patientID:      patientID,
		status:         status,
		symptoms:       symptoms,
		conversations:  conversations,
		recommendations: recommendations,
		domainEvents:   make([]event.DomainEvent, 0),
		createdAt:      createdAt,
		updatedAt:      updatedAt,
		completedAt:    completedAt,
	}
}

// ID returns the session ID.
func (s *DiagnosisSession) ID() string {
	return s.id
}

// PatientID returns the patient ID.
func (s *DiagnosisSession) PatientID() string {
	return s.patientID
}

// Status returns the current status.
func (s *DiagnosisSession) Status() SessionStatus {
	return s.status
}

// Symptoms returns a copy of the symptoms slice.
func (s *DiagnosisSession) Symptoms() []*entity.Symptom {
	result := make([]*entity.Symptom, len(s.symptoms))
	copy(result, s.symptoms)
	return result
}

// Conversations returns a copy of the conversations slice.
func (s *DiagnosisSession) Conversations() []*entity.Conversation {
	result := make([]*entity.Conversation, len(s.conversations))
	copy(result, s.conversations)
	return result
}

// Recommendations returns the recommendation set.
func (s *DiagnosisSession) Recommendations() *vo.RecommendationSet {
	return s.recommendations
}

// CreatedAt returns the creation timestamp.
func (s *DiagnosisSession) CreatedAt() time.Time {
	return s.createdAt
}

// UpdatedAt returns the last update timestamp.
func (s *DiagnosisSession) UpdatedAt() time.Time {
	return s.updatedAt
}

// CompletedAt returns the completion timestamp (may be nil).
func (s *DiagnosisSession) CompletedAt() *time.Time {
	return s.completedAt
}

// AddSymptom adds a symptom to the session.
func (s *DiagnosisSession) AddSymptom(bodyPart, description string, severity int) error {
	if s.status != StatusActive {
		return fmt.Errorf("%w: cannot add symptom to session with status %s", ErrSessionNotActive, s.status)
	}

	symptom, err := entity.NewSymptom(s.id, bodyPart, description, severity)
	if err != nil {
		return fmt.Errorf("failed to create symptom: %w", err)
	}

	s.symptoms = append(s.symptoms, symptom)
	s.updatedAt = time.Now()
	
	s.addEvent(event.NewSymptomAddedEvent(s.id, symptom.ID, bodyPart, description, severity))
	
	return nil
}

// AddConversation adds a message to the conversation history.
func (s *DiagnosisSession) AddConversation(role entity.ConversationRole, content string, msgType entity.MessageType) error {
	if s.status == StatusCompleted {
		return fmt.Errorf("%w: cannot add conversation to completed session", ErrSessionAlreadyComplete)
	}

	sequence := len(s.conversations)
	msg, err := entity.NewConversation(s.id, role, content, msgType, sequence)
	if err != nil {
		return fmt.Errorf("failed to create conversation: %w", err)
	}

	s.conversations = append(s.conversations, msg)
	s.updatedAt = time.Now()
	
	s.addEvent(event.NewConversationAddedEvent(s.id, msg.ID, string(role), content, string(msgType)))
	
	return nil
}

// GenerateRecommendation generates and stores a recommendation.
func (s *DiagnosisSession) GenerateRecommendation(rec *vo.Recommendation) error {
	if s.status == StatusCompleted {
		return ErrSessionAlreadyComplete
	}

	switch rec.Type() {
	case vo.RecommendationTypeDisease:
		s.recommendations.AddDisease(rec)
	case vo.RecommendationTypeDoctor:
		s.recommendations.AddDoctor(rec)
	case vo.RecommendationTypeDrug:
		s.recommendations.AddDrug(rec)
	case vo.RecommendationTypeDepartment:
		s.recommendations.AddDepartment(rec)
	}

	s.updatedAt = time.Now()
	
	s.addEvent(event.NewRecommendationGeneratedEvent(
		s.id,
		len(s.recommendations.Diseases()),
		len(s.recommendations.Doctors()),
		len(s.recommendations.Drugs()),
	))
	
	return nil
}

// Complete marks the session as completed.
func (s *DiagnosisSession) Complete() error {
	if s.status == StatusCompleted {
		return ErrSessionAlreadyComplete
	}

	now := time.Now()
	s.status = StatusCompleted
	s.completedAt = &now
	s.updatedAt = now
	
	s.addEvent(event.NewSessionCompletedEvent(
		s.id,
		len(s.symptoms),
		len(s.conversations),
		countAllRecommendations(s.recommendations),
	))
	
	return nil
}

// StartAnalysis transitions session to analyzing state.
func (s *DiagnosisSession) StartAnalysis() error {
	if s.status != StatusActive {
		return fmt.Errorf("%w: can only start analysis from active state", ErrSessionNotActive)
	}
	
	s.status = StatusAnalyzing
	s.updatedAt = time.Now()
	
	s.addEvent(event.NewAnalysisStartedEvent(s.id, s.GetContextForRAG()))
	
	return nil
}

// IsReadyForAnalysis checks if session has enough data for analysis.
func (s *DiagnosisSession) IsReadyForAnalysis() bool {
	return len(s.symptoms) > 0 && s.status == StatusActive
}

// GetContextForRAG generates a context string for RAG-based analysis.
func (s *DiagnosisSession) GetContextForRAG() string {
	var builder strings.Builder
	
	builder.WriteString(fmt.Sprintf("Patient ID: %s\n", s.patientID))
	builder.WriteString("Symptoms:\n")
	
	for i, symptom := range s.symptoms {
		builder.WriteString(fmt.Sprintf("  %d. %s: %s (severity: %d/10)\n",
			i+1, symptom.BodyPart, symptom.Description, symptom.Severity))
		if symptom.Duration != "" {
			builder.WriteString(fmt.Sprintf("     Duration: %s\n", symptom.Duration))
		}
	}
	
	if len(s.conversations) > 0 {
		builder.WriteString("\nConversation History:\n")
		for _, msg := range s.conversations {
			builder.WriteString(fmt.Sprintf("  %s: %s\n", msg.Role, msg.Content))
		}
	}
	
	return builder.String()
}

// GetSymptomSummary returns a summary of symptoms for analysis.
func (s *DiagnosisSession) GetSymptomSummary() string {
	var parts []string
	for _, symptom := range s.symptoms {
		parts = append(parts, fmt.Sprintf("%s (%d/10)", symptom.Description, symptom.Severity))
	}
	return strings.Join(parts, "; ")
}

// GetLastPatientMessage returns the last message from the patient.
func (s *DiagnosisSession) GetLastPatientMessage() *entity.Conversation {
	for i := len(s.conversations) - 1; i >= 0; i-- {
		if s.conversations[i].IsPatientMessage() {
			return s.conversations[i]
		}
	}
	return nil
}

// GetLastAIMessage returns the last message from the AI.
func (s *DiagnosisSession) GetLastAIMessage() *entity.Conversation {
	for i := len(s.conversations) - 1; i >= 0; i-- {
		if s.conversations[i].IsAIMessage() {
			return s.conversations[i]
		}
	}
	return nil
}

// DomainEvents returns uncommitted domain events.
func (s *DiagnosisSession) DomainEvents() []event.DomainEvent {
	return s.domainEvents
}

// ClearDomainEvents clears all uncommitted events after persistence.
func (s *DiagnosisSession) ClearDomainEvents() {
	s.domainEvents = make([]event.DomainEvent, 0)
}

// addEvent adds a domain event to the aggregate.
func (s *DiagnosisSession) addEvent(e event.DomainEvent) {
	s.domainEvents = append(s.domainEvents, e)
}

// countAllRecommendations counts all recommendations in the set.
func countAllRecommendations(rs *vo.RecommendationSet) int {
	return len(rs.Diseases()) + len(rs.Doctors()) + len(rs.Drugs()) + len(rs.Departments())
}