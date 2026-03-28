package skills

import (
	"context"
	"errors"
	"fmt"

	"intelligent_guidance_system_v2/service/ai/internal/agent"

	"github.com/go-kratos/kratos/v2/log"
)

var (
	ErrInvalidSymptomInput = errors.New("invalid symptom input")
)

type SymptomCollectorSkill struct {
	logger *log.Helper
}

func NewSymptomCollectorSkill(logger log.Logger) *SymptomCollectorSkill {
	return &SymptomCollectorSkill{
		logger: log.NewHelper(logger),
	}
}

func (s *SymptomCollectorSkill) Name() string {
	return "collect_symptom"
}

func (s *SymptomCollectorSkill) Description() string {
	return "Collects additional symptom information from the patient. Input should include body_part to focus on."
}

func (s *SymptomCollectorSkill) Execute(ctx context.Context, input map[string]interface{}) (*agent.SkillResult, error) {
	symptoms, ok := input["symptoms"].([]agent.SymptomInput)
	if !ok {
		symptoms = make([]agent.SymptomInput, 0)
	}

	bodyPart, _ := input["body_part"].(string)
	question, _ := input["question"].(string)

	if len(symptoms) == 0 {
		return &agent.SkillResult{
			Success: true,
			Output:  "No symptoms recorded yet. Please describe your main symptom: What part of your body is affected and how does it feel?",
			Data: map[string]interface{}{
				"next_question": "What is your main symptom?",
				"status":        "initial",
			},
		}, nil
	}

	if bodyPart != "" {
		followUp := s.generateFollowUpQuestion(bodyPart, symptoms)
		return &agent.SkillResult{
			Success: true,
			Output:  followUp,
			Data: map[string]interface{}{
				"next_question": followUp,
				"focus_area":    bodyPart,
				"status":        "follow_up",
			},
		}, nil
	}

	if question != "" {
		return &agent.SkillResult{
			Success: true,
			Output:  question,
			Data: map[string]interface{}{
				"next_question": question,
				"status":        "custom",
			},
		}, nil
	}

	symptomSummary := s.summarizeSymptoms(symptoms)
	nextQuestion := s.generateGeneralFollowUp(symptoms)

	return &agent.SkillResult{
		Success: true,
		Output:  fmt.Sprintf("Current symptoms recorded:\n%s\n\n%s", symptomSummary, nextQuestion),
		Data: map[string]interface{}{
			"symptom_count": len(symptoms),
			"next_question": nextQuestion,
			"status":        "collecting",
		},
	}, nil
}

func (s *SymptomCollectorSkill) generateFollowUpQuestion(bodyPart string, symptoms []agent.SymptomInput) string {
	for _, sym := range symptoms {
		if sym.BodyPart == bodyPart {
			return fmt.Sprintf("You mentioned %s in your %s. How long have you been experiencing this? Has it gotten worse or stayed the same?",
				sym.Description, bodyPart)
		}
	}
	return fmt.Sprintf("Please describe the symptoms you're experiencing in your %s. How severe is the discomfort (1-10)?", bodyPart)
}

func (s *SymptomCollectorSkill) generateGeneralFollowUp(symptoms []agent.SymptomInput) string {
	hasDuration := false
	for _, sym := range symptoms {
		for _, s := range symptoms {
			if s.Description != "" && len(s.Description) > 20 {
				hasDuration = true
			}
		}
	}

	if !hasDuration {
		return "How long have you been experiencing these symptoms? Have they changed over time?"
	}

	return "Are there any other symptoms you haven't mentioned yet? Any fever, fatigue, or changes in appetite?"
}

func (s *SymptomCollectorSkill) summarizeSymptoms(symptoms []agent.SymptomInput) string {
	summary := ""
	for i, sym := range symptoms {
		summary += fmt.Sprintf("%d. %s: %s (severity: %d)\n", i+1, sym.BodyPart, sym.Description, sym.Severity)
	}
	return summary
}