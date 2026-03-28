// Package entity provides domain entities for the AI diagnosis service.
package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ConversationRole represents the role of a message sender.
type ConversationRole string

const (
	// RolePatient indicates message is from the patient
	RolePatient ConversationRole = "patient"
	
	// RoleAI indicates message is from the AI system
	RoleAI ConversationRole = "ai"
	
	// RoleDoctor indicates message is from a doctor
	RoleDoctor ConversationRole = "doctor"
)

// MessageType represents the type of conversation message.
type MessageType string

const (
	// MessageTypeSymptom is a symptom-related message
	MessageTypeSymptom MessageType = "symptom"
	
	// MessageTypeQuestion is a question from AI
	MessageTypeQuestion MessageType = "question"
	
	// MessageTypeAnswer is an answer from patient
	MessageTypeAnswer MessageType = "answer"
	
	// MessageTypeRecommendation is a diagnosis recommendation
	MessageTypeRecommendation MessageType = "recommendation"
	
	// MessageTypeClarification is a clarification request
	MessageTypeClarification MessageType = "clarification"
)

// Conversation represents a single message in the diagnosis dialogue.
// It tracks the multi-turn conversation between patient and AI.
type Conversation struct {
	// ID is the unique identifier for the conversation message
	ID string `json:"id" gorm:"primaryKey;type:varchar(36)"`
	
	// SessionID references the parent diagnosis session
	SessionID string `json:"session_id" gorm:"index;type:varchar(36)"`
	
	// Role indicates who sent the message
	Role ConversationRole `json:"role" gorm:"type:varchar(20)"`
	
	// Content is the message text
	Content string `json:"content" gorm:"type:text"`
	
	// Type categorizes the message purpose
	Type MessageType `json:"type" gorm:"type:varchar(30)"`
	
	// Sequence indicates the order of message in conversation
	Sequence int `json:"sequence" gorm:"type:int"`
	
	// Metadata stores additional message attributes
	Metadata map[string]interface{} `json:"metadata" gorm:"type:json"`
	
	// CreatedAt is the timestamp when message was created
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// Error definitions for Conversation entity
var (
	ErrInvalidRole      = errors.New("invalid conversation role")
	ErrInvalidContent   = errors.New("invalid conversation content")
	ErrInvalidMsgType   = errors.New("invalid message type")
)

// NewConversation creates a new Conversation entity.
func NewConversation(sessionID string, role ConversationRole, content string, msgType MessageType, sequence int) (*Conversation, error) {
	if !isValidRole(role) {
		return nil, ErrInvalidRole
	}
	if content == "" {
		return nil, ErrInvalidContent
	}
	if !isValidMsgType(msgType) {
		return nil, ErrInvalidMsgType
	}

	return &Conversation{
		ID:        uuid.New().String(),
		SessionID: sessionID,
		Role:      role,
		Content:   content,
		Type:      msgType,
		Sequence:  sequence,
		Metadata:  make(map[string]interface{}),
		CreatedAt: time.Now(),
	}, nil
}

// isValidRole validates the conversation role.
func isValidRole(role ConversationRole) bool {
	switch role {
	case RolePatient, RoleAI, RoleDoctor:
		return true
	default:
		return false
	}
}

// isValidMsgType validates the message type.
func isValidMsgType(msgType MessageType) bool {
	switch msgType {
	case MessageTypeSymptom, MessageTypeQuestion, MessageTypeAnswer,
		MessageTypeRecommendation, MessageTypeClarification:
		return true
	default:
		return false
	}
}

// AddMetadata adds metadata to the conversation.
func (c *Conversation) AddMetadata(key string, value interface{}) {
	if c.Metadata == nil {
		c.Metadata = make(map[string]interface{})
	}
	c.Metadata[key] = value
}

// IsPatientMessage returns true if message is from patient.
func (c *Conversation) IsPatientMessage() bool {
	return c.Role == RolePatient
}

// IsAIMessage returns true if message is from AI.
func (c *Conversation) IsAIMessage() bool {
	return c.Role == RoleAI
}

// TableName returns the database table name for GORM.
func (Conversation) TableName() string {
	return "ai_conversations"
}