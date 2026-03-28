// Package dto provides Data Transfer Objects for the business layer.
package dto

import (
	"time"
)

// CreateSessionRequest is the request to create a new diagnosis session.
type CreateSessionRequest struct {
	PatientID string `json:"patient_id"`
}

// CreateSessionResponse is the response after creating a session.
type CreateSessionResponse struct {
	SessionID string    `json:"session_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// AddSymptomRequest is the request to add a symptom.
type AddSymptomRequest struct {
	SessionID   string `json:"session_id"`
	BodyPart    string `json:"body_part"`
	Description string `json:"description"`
	Severity    int    `json:"severity"`
	Duration    string `json:"duration,omitempty"`
}

// AddSymptomResponse is the response after adding a symptom.
type AddSymptomResponse struct {
	SymptomID string    `json:"symptom_id"`
	CreatedAt time.Time `json:"created_at"`
}

// AddConversationRequest is the request to add a conversation message.
type AddConversationRequest struct {
	SessionID string `json:"session_id"`
	Role      string `json:"role"`
	Content   string `json:"content"`
	Type      string `json:"type"`
}

// AddConversationResponse is the response after adding a conversation.
type AddConversationResponse struct {
	MessageID string    `json:"message_id"`
	Sequence  int       `json:"sequence"`
	CreatedAt time.Time `json:"created_at"`
}

// GenerateRecommendationRequest is the request to generate recommendations.
type GenerateRecommendationRequest struct {
	SessionID       string  `json:"session_id"`
	TopK            int     `json:"top_k"`
	ScoreThreshold  float32 `json:"score_threshold"`
	IncludeDrugs     bool    `json:"include_drugs"`
	IncludeDoctors   bool    `json:"include_doctors"`
}

// GenerateRecommendationResponse is the response with recommendations.
type GenerateRecommendationResponse struct {
	SessionID      string                       `json:"session_id"`
	Diseases       []RecommendationDTO          `json:"diseases"`
	Doctors        []RecommendationDTO          `json:"doctors"`
	Drugs          []RecommendationDTO          `json:"drugs"`
	Departments    []RecommendationDTO          `json:"departments"`
	GeneratedAt    time.Time                    `json:"generated_at"`
}

// RecommendationDTO is the DTO for a recommendation.
type RecommendationDTO struct {
	ID              string                 `json:"id"`
	Type            string                 `json:"type"`
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	ConfidenceScore float32                `json:"confidence_score"`
	ConfidenceLevel string                 `json:"confidence_level"`
	Reasoning       string                 `json:"reasoning,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// GetSessionResponse is the response for getting session details.
type GetSessionResponse struct {
	SessionID      string                   `json:"session_id"`
	PatientID      string                   `json:"patient_id"`
	Status         string                   `json:"status"`
	Symptoms       []SymptomDTO             `json:"symptoms"`
	Conversations  []ConversationDTO        `json:"conversations"`
	Recommendations *RecommendationSetDTO   `json:"recommendations,omitempty"`
	CreatedAt      time.Time                `json:"created_at"`
	UpdatedAt      time.Time                 `json:"updated_at"`
	CompletedAt    *time.Time               `json:"completed_at,omitempty"`
}

// SymptomDTO is the DTO for a symptom.
type SymptomDTO struct {
	ID          string                 `json:"id"`
	BodyPart    string                 `json:"body_part"`
	Description string                 `json:"description"`
	Severity    int                    `json:"severity"`
	Duration    string                 `json:"duration,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
}

// ConversationDTO is the DTO for a conversation message.
type ConversationDTO struct {
	ID        string    `json:"id"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Type      string    `json:"type"`
	Sequence  int       `json:"sequence"`
	CreatedAt time.Time `json:"created_at"`
}

// RecommendationSetDTO is the DTO for a set of recommendations.
type RecommendationSetDTO struct {
	Diseases    []RecommendationDTO `json:"diseases"`
	Doctors     []RecommendationDTO `json:"doctors"`
	Drugs       []RecommendationDTO `json:"drugs"`
	Departments []RecommendationDTO `json:"departments"`
}

// CompleteSessionResponse is the response after completing a session.
type CompleteSessionResponse struct {
	SessionID   string     `json:"session_id"`
	Status      string     `json:"status"`
	CompletedAt time.Time  `json:"completed_at"`
}

// ListSessionsRequest is the request to list sessions.
type ListSessionsRequest struct {
	PatientID string `json:"patient_id"`
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
}

// ListSessionsResponse is the response for listing sessions.
type ListSessionsResponse struct {
	Sessions []SessionSummaryDTO `json:"sessions"`
	Total    int                 `json:"total"`
}

// SessionSummaryDTO is a summary of a session.
type SessionSummaryDTO struct {
	SessionID       string     `json:"session_id"`
	PatientID       string     `json:"patient_id"`
	Status          string     `json:"status"`
	SymptomCount    int        `json:"symptom_count"`
	CreatedAt       time.Time  `json:"created_at"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
}

// AIPromptRequest is the request for AI prompt processing.
type AIPromptRequest struct {
	SessionID string `json:"session_id"`
	Prompt    string `json:"prompt"`
}

// AIPromptResponse is the response from AI prompt processing.
type AIPromptResponse struct {
	SessionID   string             `json:"session_id"`
	Response    string             `json:"response"`
	Recommendations *RecommendationSetDTO `json:"recommendations,omitempty"`
}

// VectorSearchRequest is the request for vector search.
type VectorSearchRequest struct {
	Collection     string                 `json:"collection"`
	QueryText      string                 `json:"query_text"`
	TopK           int                    `json:"top_k"`
	ScoreThreshold float32                `json:"score_threshold"`
	Filters        map[string]interface{} `json:"filters,omitempty"`
}

// VectorSearchResponse is the response from vector search.
type VectorSearchResponse struct {
	Results []VectorSearchResultDTO `json:"results"`
	QueryID string                 `json:"query_id"`
}

// VectorSearchResultDTO is the DTO for a vector search result.
type VectorSearchResultDTO struct {
	ID      string                 `json:"id"`
	Score   float32                `json:"score"`
	Payload map[string]interface{} `json:"payload"`
}