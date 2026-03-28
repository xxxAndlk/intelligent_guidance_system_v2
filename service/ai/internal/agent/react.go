package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"intelligent_guidance_system_v2/service/ai/internal/data/llm"

	"github.com/go-kratos/kratos/v2/log"
)

var (
	ErrMaxIterationsReached = errors.New("max iterations reached")
	ErrInvalidActionFormat  = errors.New("invalid action format")
	ErrUnknownSkill         = errors.New("unknown skill")
	ErrAgentFailed          = errors.New("agent failed to complete task")
)

type ReActAgent struct {
	llmClient       *llm.DeepSeekClient
	skills          map[string]Skill
	config          AgentConfig
	logger          *log.Helper
	promptBuilder   *llm.DiagnosisPromptBuilder
}

func NewReActAgent(llmClient *llm.DeepSeekClient, config AgentConfig, logger log.Logger) *ReActAgent {
	if config.MaxIterations <= 0 {
		config.MaxIterations = 10
	}
	if config.SystemPrompt == "" {
		config.SystemPrompt = buildDefaultSystemPrompt()
	}

	return &ReActAgent{
		llmClient:     llmClient,
		skills:        make(map[string]Skill),
		config:        config,
		logger:        log.NewHelper(logger),
		promptBuilder: llm.NewDiagnosisPromptBuilder(),
	}
}

func (a *ReActAgent) Execute(ctx context.Context, task string, agentCtx *AgentContext) (*AgentResult, error) {
	result := &AgentResult{
		Steps: make([]AgentStep, 0),
		Data:  make(map[string]interface{}),
	}

	history := a.buildInitialHistory(task, agentCtx)

	for i := 0; i < a.config.MaxIterations; i++ {
		a.logger.Infof("ReAct iteration %d", i+1)

		response, err := a.think(ctx, history)
		if err != nil {
			result.Error = err.Error()
			result.Success = false
			return result, err
		}

		step := a.parseResponse(response)
		result.Steps = append(result.Steps, step)

		if step.Type == "final_answer" {
			result.Success = true
			result.Output = step.Observation
			return result, nil
		}

		if step.Action != "" {
			skillResult, err := a.act(ctx, step.Action, step.ActionInput, agentCtx)
			if err != nil {
				step.Observation = fmt.Sprintf("Error: %v", err)
			} else {
				step.Observation = skillResult.Output
				if skillResult.Data != nil {
					for k, v := range skillResult.Data {
						result.Data[k] = v
					}
				}
			}

			result.Steps[len(result.Steps)-1] = step

			history = append(history, llm.ChatMessage{
				Role:    "assistant",
				Content: response,
			})
			history = append(history, llm.ChatMessage{
				Role:    "user",
				Content: fmt.Sprintf("Observation: %s", step.Observation),
			})
		}
	}

	result.Error = ErrMaxIterationsReached.Error()
	result.Success = false
	return result, ErrMaxIterationsReached
}

func (a *ReActAgent) think(ctx context.Context, history []llm.ChatMessage) (string, error) {
	messages := []llm.ChatMessage{
		{Role: "system", Content: a.config.SystemPrompt},
	}
	messages = append(messages, history...)

	response, err := a.llmClient.Chat(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("LLM call failed: %w", err)
	}

	return response, nil
}

func (a *ReActAgent) act(ctx context.Context, action, actionInput string, agentCtx *AgentContext) (*SkillResult, error) {
	skill, exists := a.skills[action]
	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrUnknownSkill, action)
	}

	input := make(map[string]interface{})
	if actionInput != "" {
		if err := json.Unmarshal([]byte(actionInput), &input); err != nil {
			input["raw_input"] = actionInput
		}
	}

	if agentCtx != nil {
		input["session_id"] = agentCtx.SessionID
		input["patient_id"] = agentCtx.PatientID
		input["symptoms"] = agentCtx.Symptoms
		input["history"] = agentCtx.History
	}

	return skill.Execute(ctx, input)
}

func (a *ReActAgent) parseResponse(response string) AgentStep {
	step := AgentStep{Type: "unknown"}

	response = strings.TrimSpace(response)

	finalAnswerRegex := regexp.MustCompile(`Final Answer:\s*(.+)`)
	if match := finalAnswerRegex.FindStringSubmatch(response); match != nil {
		step.Type = "final_answer"
		step.Observation = strings.TrimSpace(match[1])
		return step
	}

	thoughtRegex := regexp.MustCompile(`Thought:\s*(.+?)(?=Action:|Final Answer:|$)`)
	if match := thoughtRegex.FindStringSubmatch(response); match != nil {
		step.Thought = strings.TrimSpace(match[1])
	}

	actionRegex := regexp.MustCompile(`Action:\s*([a-z_]+)`)
	if match := actionRegex.FindStringSubmatch(response); match != nil {
		step.Action = strings.TrimSpace(match[1])
		step.Type = "action"
	}

	actionInputRegex := regexp.MustCompile(`Action Input:\s*(.+?)(?=Observation:|$)`)
	if match := actionInputRegex.FindStringSubmatch(response); match != nil {
		step.ActionInput = strings.TrimSpace(match[1])
	}

	return step
}

func (a *ReActAgent) buildInitialHistory(task string, agentCtx *AgentContext) []llm.ChatMessage {
	context := fmt.Sprintf("Task: %s\n\n", task)

	if agentCtx != nil {
		if len(agentCtx.Symptoms) > 0 {
			context += "Current Symptoms:\n"
			for _, s := range agentCtx.Symptoms {
				context += fmt.Sprintf("- %s: %s (severity: %d)\n", s.BodyPart, s.Description, s.Severity)
			}
		}

		if len(agentCtx.History) > 0 {
			context += "\nConversation History:\n"
			for _, h := range agentCtx.History {
				context += fmt.Sprintf("%s: %s\n", h.Role, h.Content)
			}
		}
	}

	return []llm.ChatMessage{
		{Role: "user", Content: context},
	}
}

func (a *ReActAgent) GetAvailableSkills() []string {
	skills := make([]string, 0, len(a.skills))
	for name := range a.skills {
		skills = append(skills, name)
	}
	return skills
}

func (a *ReActAgent) AddSkill(skill Skill) {
	a.skills[skill.Name()] = skill
}

func buildDefaultSystemPrompt() string {
	return `You are a medical diagnosis AI assistant using the ReAct (Reasoning + Acting) framework.

Your role is to help diagnose patient symptoms by systematically reasoning and taking actions.

Available actions:
- collect_symptom: Collect additional symptom information from patient
- search_disease: Search for potential diseases based on symptoms
- search_doctor: Find suitable doctors for the patient's condition
- search_drug: Look up relevant medications

For each step, use the following format:

Thought: [Your reasoning about what to do next]
Action: [The action to take]
Action Input: [JSON input for the action]

After receiving an observation, continue reasoning and taking actions until you reach a conclusion.

When you have a final diagnosis recommendation, use:

Final Answer: [Your comprehensive recommendation including potential diseases, recommended doctors, and important notes]

Important guidelines:
- Always recommend consulting a real doctor
- Consider symptom combinations and severity
- Ask clarifying questions when symptoms are ambiguous
- Provide confidence levels for recommendations`
}