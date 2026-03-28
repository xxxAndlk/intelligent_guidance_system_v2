package agent

import (
	"context"
)

type AgentResult struct {
	Success     bool
	Output      string
	Data        map[string]interface{}
	Steps       []AgentStep
	Error       string
	TokensUsed  int
}

type AgentStep struct {
	Type        string
	Thought     string
	Action      string
	ActionInput string
	Observation string
}

type AgentContext struct {
	SessionID   string
	PatientID   string
	Symptoms    []SymptomInput
	History     []ConversationEntry
	Metadata    map[string]interface{}
}

type SymptomInput struct {
	BodyPart    string
	Description string
	Severity    int
}

type ConversationEntry struct {
	Role    string
	Content string
}

type Agent interface {
	Execute(ctx context.Context, task string, agentCtx *AgentContext) (*AgentResult, error)
	GetAvailableSkills() []string
	AddSkill(skill Skill)
}

type Skill interface {
	Name() string
	Description() string
	Execute(ctx context.Context, input map[string]interface{}) (*SkillResult, error)
}

type SkillResult struct {
	Success bool
	Output  string
	Data    map[string]interface{}
}

type AgentConfig struct {
	MaxIterations    int
	Temperature      float32
	SystemPrompt     string
	Model            string
}

type AgentState string

const (
	AgentStateIdle       AgentState = "idle"
	AgentStateThinking   AgentState = "thinking"
	AgentStateActing     AgentState = "acting"
	AgentStateObserving  AgentState = "observing"
	AgentStateCompleted  AgentState = "completed"
	AgentStateFailed     AgentState = "failed"
)