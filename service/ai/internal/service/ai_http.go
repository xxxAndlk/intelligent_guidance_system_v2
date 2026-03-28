package service

import (
	"context"

	"intelligent_guidance_system_v2/service/ai/internal/agent"
	"intelligent_guidance_system_v2/service/ai/internal/biz/dto"
	"intelligent_guidance_system_v2/service/ai/internal/biz/usecase"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

type AIServiceHTTP struct {
	diagnosisUC *usecase.DiagnosisUseCase
	reactAgent  agent.Agent
	logger      *log.Helper
}

func NewAIServiceHTTP(diagnosisUC *usecase.DiagnosisUseCase, reactAgent agent.Agent, logger log.Logger) *AIServiceHTTP {
	return &AIServiceHTTP{
		diagnosisUC: diagnosisUC,
		reactAgent:  reactAgent,
		logger:      log.NewHelper(logger),
	}
}

func (s *AIServiceHTTP) DiagnosisUC() *usecase.DiagnosisUseCase {
	return s.diagnosisUC
}

func (s *AIServiceHTTP) ProcessAIPrompt(ctx context.Context, req *dto.AIPromptRequest) (*dto.AIPromptResponse, error) {
	if req.SessionID == "" {
		return nil, errors.New(400, "INVALID_REQUEST", "session_id is required")
	}

	session, err := s.diagnosisUC.GetSession(ctx, req.SessionID)
	if err != nil {
		return nil, errors.New(404, "SESSION_NOT_FOUND", "session not found")
	}

	agentCtx := &agent.AgentContext{
		SessionID: req.SessionID,
		PatientID: session.PatientID,
		Symptoms:  mapSymptomsToAgentInput(session.Symptoms),
		History:   mapConversationsToAgentHistory(session.Conversations),
	}

	result, err := s.reactAgent.Execute(ctx, req.Prompt, agentCtx)
	if err != nil {
		s.logger.Errorf("Agent execution failed: %v", err)
		return nil, errors.New(500, "AGENT_FAILED", "AI agent failed to process request")
	}

	resp := &dto.AIPromptResponse{
		SessionID: req.SessionID,
		Response:  result.Output,
	}

	if len(result.Data) > 0 {
		diseases, ok := result.Data["diseases"].([]map[string]interface{})
		if ok {
			resp.Recommendations = &dto.RecommendationSetDTO{
				Diseases: mapAgentDataToRecommendations(diseases, "disease"),
			}
		}
	}

	return resp, nil
}

func (s *AIServiceHTTP) VectorSearch(ctx context.Context, req *dto.VectorSearchRequest) (*dto.VectorSearchResponse, error) {
	return nil, errors.New(501, "NOT_IMPLEMENTED", "vector search not implemented for HTTP")
}

func (s *AIServiceHTTP) HealthCheck(ctx context.Context) (map[string]interface{}, error) {
	return map[string]interface{}{
		"status":    "healthy",
		"service":   "ai-service",
		"timestamp": json.Number(""),
	}, nil
}

func mapSymptomsToAgentInput(symptoms []dto.SymptomDTO) []agent.SymptomInput {
	result := make([]agent.SymptomInput, len(symptoms))
	for i, s := range symptoms {
		result[i] = agent.SymptomInput{
			BodyPart:    s.BodyPart,
			Description: s.Description,
			Severity:    s.Severity,
		}
	}
	return result
}

func mapConversationsToAgentHistory(conversations []dto.ConversationDTO) []agent.ConversationEntry {
	result := make([]agent.ConversationEntry, len(conversations))
	for i, c := range conversations {
		result[i] = agent.ConversationEntry{
			Role:    c.Role,
			Content: c.Content,
		}
	}
	return result
}

func mapAgentDataToRecommendations(data []map[string]interface{}, recType string) []dto.RecommendationDTO {
	result := make([]dto.RecommendationDTO, len(data))
	for i, d := range data {
		name, _ := d["name"].(string)
		score, _ := d["score"].(float32)
		description, _ := d["description"].(string)

		result[i] = dto.RecommendationDTO{
			ID:              d["id"].(string),
			Type:            recType,
			Title:           name,
			Description:     description,
			ConfidenceScore: score,
		}
	}
	return result
}