package skills

import (
	"context"
	"fmt"

	"intelligent_guidance_system_v2/service/ai/internal/acl"
	"intelligent_guidance_system_v2/service/ai/internal/agent"

	"github.com/go-kratos/kratos/v2/log"
)

type DrugRecommenderSkill struct {
	vectorACL *acl.VectorDBACL
	logger    *log.Helper
}

func NewDrugRecommenderSkill(vectorACL *acl.VectorDBACL, logger log.Logger) *DrugRecommenderSkill {
	return &DrugRecommenderSkill{
		vectorACL: vectorACL,
		logger:    log.NewHelper(logger),
	}
}

func (s *DrugRecommenderSkill) Name() string {
	return "search_drug"
}

func (s *DrugRecommenderSkill) Description() string {
	return "Searches for relevant medications based on symptoms or disease context. Returns drug information with usage guidelines."
}

func (s *DrugRecommenderSkill) Execute(ctx context.Context, input map[string]interface{}) (*agent.SkillResult, error) {
	queryText := s.buildQueryText(input)
	limit := 5
	threshold := float32(0.6)

	if l, ok := input["limit"].(int); ok && l > 0 {
		limit = l
	}

	results, err := s.vectorACL.SearchDrugsByText(ctx, queryText, limit, threshold)
	if err != nil {
		s.logger.Warnf("Drug search failed: %v", err)
		return &agent.SkillResult{
			Success: false,
			Output:  fmt.Sprintf("Failed to search drugs: %v", err),
		}, nil
	}

	if len(results) == 0 {
		return &agent.SkillResult{
			Success: true,
			Output:  "No matching medications found. A doctor can recommend appropriate treatments after diagnosis.",
			Data:    map[string]interface{}{"drugs": []interface{}{}},
		}, nil
	}

	output := s.formatDrugResults(results)
	drugData := s.extractDrugData(results)

	return &agent.SkillResult{
		Success: true,
		Output:  output,
		Data: map[string]interface{}{
			"drugs":        drugData,
			"result_count": len(results),
		},
	}, nil
}

func (s *DrugRecommenderSkill) buildQueryText(input map[string]interface{}) string {
	queryText := ""

	if symptoms, ok := input["symptoms"].([]agent.SymptomInput); ok {
		for _, sym := range symptoms {
			queryText += sym.Description + " "
		}
	}

	if disease, ok := input["disease"].(string); ok {
		queryText += disease + " treatment "
	}

	if query, ok := input["query"].(string); ok {
		queryText += query
	}

	return queryText
}

func (s *DrugRecommenderSkill) formatDrugResults(results []*acl.VectorSearchResult) string {
	output := "Relevant medications (for reference only):\n\n"
	output += "⚠️ IMPORTANT: These are informational suggestions. Always consult a doctor or pharmacist before taking any medication.\n\n"

	for i, result := range results {
		name, _ := result.Payload["name"].(string)
		indication, _ := result.Payload["indication"].(string)
		contraindications := s.extractContraindications(result.Payload)

		output += fmt.Sprintf("%d. %s\n", i+1, name)
		if indication != "" {
			output += fmt.Sprintf("   Used for: %s\n", indication)
		}
		if len(contraindications) > 0 {
			output += fmt.Sprintf("   Contraindications: %s\n", contraindications)
		}
		output += "\n"
	}

	return output
}

func (s *DrugRecommenderSkill) extractContraindications(payload map[string]interface{}) string {
	if ci, ok := payload["contraindications"].([]string); ok {
		return fmt.Sprintf("%v", ci)
	}
	if ci, ok := payload["contraindications"].([]interface{}); ok {
		strs := make([]string, len(ci))
		for i, v := range ci {
			strs[i] = fmt.Sprintf("%v", v)
		}
		return fmt.Sprintf("%v", strs)
	}
	return ""
}

func (s *DrugRecommenderSkill) extractDrugData(results []*acl.VectorSearchResult) []map[string]interface{} {
	drugs := make([]map[string]interface{}, len(results))

	for i, result := range results {
		drugs[i] = map[string]interface{}{
			"id":       result.ID,
			"score":    result.Score,
			"name":     result.Payload["name"],
			"indication": result.Payload["indication"],
		}
	}

	return drugs
}