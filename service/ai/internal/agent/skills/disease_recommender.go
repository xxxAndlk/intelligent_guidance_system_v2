package skills

import (
	"context"
	"fmt"

	"intelligent_guidance_system_v2/service/ai/internal/acl"
	"intelligent_guidance_system_v2/service/ai/internal/agent"

	"github.com/go-kratos/kratos/v2/log"
)

type DiseaseRecommenderSkill struct {
	vectorACL *acl.VectorDBACL
	logger    *log.Helper
}

func NewDiseaseRecommenderSkill(vectorACL *acl.VectorDBACL, logger log.Logger) *DiseaseRecommenderSkill {
	return &DiseaseRecommenderSkill{
		vectorACL: vectorACL,
		logger:    log.NewHelper(logger),
	}
}

func (s *DiseaseRecommenderSkill) Name() string {
	return "search_disease"
}

func (s *DiseaseRecommenderSkill) Description() string {
	return "Searches for potential diseases based on symptoms using vector similarity search. Input should include symptom descriptions."
}

func (s *DiseaseRecommenderSkill) Execute(ctx context.Context, input map[string]interface{}) (*agent.SkillResult, error) {
	queryText := s.buildQueryText(input)
	limit := 5
	threshold := float32(0.7)

	if l, ok := input["limit"].(int); ok && l > 0 {
		limit = l
	}
	if t, ok := input["threshold"].(float32); ok && t > 0 {
		threshold = t
	}

	results, err := s.vectorACL.SearchDiseasesByText(ctx, queryText, limit, threshold)
	if err != nil {
		s.logger.Warnf("Disease search failed: %v", err)
		return &agent.SkillResult{
			Success: false,
			Output:  fmt.Sprintf("Failed to search diseases: %v", err),
		}, nil
	}

	if len(results) == 0 {
		return &agent.SkillResult{
			Success: true,
			Output:  "No matching diseases found in database. Consider expanding symptoms or consulting a doctor directly.",
			Data: map[string]interface{}{
				"results": []interface{}{},
			},
		}, nil
	}

	output := s.formatResults(results)
	diseaseData := s.extractDiseaseData(results)

	return &agent.SkillResult{
		Success: true,
		Output:  output,
		Data: map[string]interface{}{
			"diseases":  diseaseData,
			"result_count": len(results),
		},
	}, nil
}

func (s *DiseaseRecommenderSkill) buildQueryText(input map[string]interface{}) string {
	queryText := ""

	if symptoms, ok := input["symptoms"].([]agent.SymptomInput); ok {
		for _, sym := range symptoms {
			queryText += sym.BodyPart + " " + sym.Description + " "
		}
	}

	if query, ok := input["query"].(string); ok {
		queryText += query
	}

	return queryText
}

func (s *DiseaseRecommenderSkill) formatResults(results []*acl.VectorSearchResult) string {
	output := "Potential diseases based on symptoms:\n\n"

	for i, result := range results {
		name, _ := result.Payload["name"].(string)
		description, _ := result.Payload["description"].(string)
		category, _ := result.Payload["category"].(string)

		output += fmt.Sprintf("%d. %s (confidence: %.2f)\n", i+1, name, result.Score)
		if category != "" {
			output += fmt.Sprintf("   Category: %s\n", category)
		}
		if description != "" {
			output += fmt.Sprintf("   Description: %s\n", description)
		}
		output += "\n"
	}

	output += "Note: These are AI suggestions. Please consult a doctor for definitive diagnosis."

	return output
}

func (s *DiseaseRecommenderSkill) extractDiseaseData(results []*acl.VectorSearchResult) []map[string]interface{} {
	diseases := make([]map[string]interface{}, len(results))

	for i, result := range results {
		diseases[i] = map[string]interface{}{
			"id":       result.ID,
			"score":    result.Score,
			"name":     result.Payload["name"],
			"category": result.Payload["category"],
			"description": result.Payload["description"],
		}
	}

	return diseases
}