// Package event provides domain events for the AI diagnosis service.
// Domain events represent significant occurrences within the domain.
package event

import (
	"time"

	"github.com/google/uuid"
)

// EventType represents the type of domain event.
type EventType string

const (
	// EventSessionCreated is emitted when a new diagnosis session is created
	EventSessionCreated EventType = "diagnosis_session.created"
	
	// EventSymptomAdded is emitted when a symptom is added to a session
	EventSymptomAdded EventType = "diagnosis_session.symptom_added"
	
	// EventConversationAdded is emitted when a conversation message is added
	EventConversationAdded EventType = "diagnosis_session.conversation_added"
	
	// EventRecommendationGenerated is emitted when recommendations are generated
	EventRecommendationGenerated EventType = "diagnosis_session.recommendation_generated"
	
	// EventSessionCompleted is emitted when a session is marked complete
	EventSessionCompleted EventType = "diagnosis_session.completed"
	
	// EventAnalysisStarted is emitted when AI analysis begins
	EventAnalysisStarted EventType = "diagnosis_session.analysis_started"
)

// DomainEvent is the base interface for all domain events.
type DomainEvent interface {
	// EventID returns the unique event identifier
	EventID() string
	
	// EventType returns the type of event
	EventType() EventType
	
	// OccurredAt returns when the event occurred
	OccurredAt() time.Time
	
	// AggregateID returns the aggregate root ID that emitted this event
	AggregateID() string
}

// BaseEvent provides common fields for domain events.
type BaseEvent struct {
	eventID     string
	eventType   EventType
	occurredAt  time.Time
	aggregateID string
}

func newBaseEvent(eventType EventType, aggregateID string) BaseEvent {
	return BaseEvent{
		eventID:     uuid.New().String(),
		eventType:   eventType,
		occurredAt:  time.Now(),
		aggregateID: aggregateID,
	}
}

func (e BaseEvent) EventID() string     { return e.eventID }
func (e BaseEvent) EventType() EventType { return e.eventType }
func (e BaseEvent) OccurredAt() time.Time { return e.occurredAt }
func (e BaseEvent) AggregateID() string { return e.aggregateID }

// SessionCreatedEvent is emitted when a diagnosis session is created.
type SessionCreatedEvent struct {
	BaseEvent
	PatientID string
}

// NewSessionCreatedEvent creates a new SessionCreatedEvent.
func NewSessionCreatedEvent(sessionID, patientID string) *SessionCreatedEvent {
	return &SessionCreatedEvent{
		BaseEvent: newBaseEvent(EventSessionCreated, sessionID),
		PatientID: patientID,
	}
}

// SymptomAddedEvent is emitted when a symptom is added to a session.
type SymptomAddedEvent struct {
	BaseEvent
	SymptomID   string
	BodyPart    string
	Description string
	Severity    int
}

// NewSymptomAddedEvent creates a new SymptomAddedEvent.
func NewSymptomAddedEvent(sessionID, symptomID, bodyPart, description string, severity int) *SymptomAddedEvent {
	return &SymptomAddedEvent{
		BaseEvent:   newBaseEvent(EventSymptomAdded, sessionID),
		SymptomID:   symptomID,
		BodyPart:    bodyPart,
		Description: description,
		Severity:    severity,
	}
}

// ConversationAddedEvent is emitted when a message is added to conversation.
type ConversationAddedEvent struct {
	BaseEvent
	MessageID string
	Role      string
	Content   string
	MsgType   string
}

// NewConversationAddedEvent creates a new ConversationAddedEvent.
func NewConversationAddedEvent(sessionID, messageID, role, content, msgType string) *ConversationAddedEvent {
	return &ConversationAddedEvent{
		BaseEvent: newBaseEvent(EventConversationAdded, sessionID),
		MessageID: messageID,
		Role:      role,
		Content:   content,
		MsgType:   msgType,
	}
}

// RecommendationGeneratedEvent is emitted when AI generates recommendations.
type RecommendationGeneratedEvent struct {
	BaseEvent
	DiseaseCount int
	DoctorCount  int
	DrugCount    int
}

// NewRecommendationGeneratedEvent creates a new RecommendationGeneratedEvent.
func NewRecommendationGeneratedEvent(sessionID string, diseaseCount, doctorCount, drugCount int) *RecommendationGeneratedEvent {
	return &RecommendationGeneratedEvent{
		BaseEvent:    newBaseEvent(EventRecommendationGenerated, sessionID),
		DiseaseCount: diseaseCount,
		DoctorCount:  doctorCount,
		DrugCount:    drugCount,
	}
}

// SessionCompletedEvent is emitted when a session is completed.
type SessionCompletedEvent struct {
	BaseEvent
	SymptomCount        int
	ConversationCount   int
	RecommendationCount int
}

// NewSessionCompletedEvent creates a new SessionCompletedEvent.
func NewSessionCompletedEvent(sessionID string, symptomCount, conversationCount, recommendationCount int) *SessionCompletedEvent {
	return &SessionCompletedEvent{
		BaseEvent:           newBaseEvent(EventSessionCompleted, sessionID),
		SymptomCount:        symptomCount,
		ConversationCount:   conversationCount,
		RecommendationCount: recommendationCount,
	}
}

// AnalysisStartedEvent is emitted when AI analysis begins.
type AnalysisStartedEvent struct {
	BaseEvent
	QueryContext string
}

// NewAnalysisStartedEvent creates a new AnalysisStartedEvent.
func NewAnalysisStartedEvent(sessionID, queryContext string) *AnalysisStartedEvent {
	return &AnalysisStartedEvent{
		BaseEvent:    newBaseEvent(EventAnalysisStarted, sessionID),
		QueryContext: queryContext,
	}
}

// EventBus defines the interface for publishing domain events.
type EventBus interface {
	// Publish publishes one or more domain events
	Publish(events ...DomainEvent) error
	
	// Subscribe registers a handler for specific event types
	Subscribe(eventType EventType, handler EventHandler) error
}

// EventHandler handles domain events.
type EventHandler func(event DomainEvent) error