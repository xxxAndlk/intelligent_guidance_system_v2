package skills

import (
	"context"
	"fmt"

	"intelligent_guidance_system_v2/service/ai/internal/acl"
	"intelligent_guidance_system_v2/service/ai/internal/agent"

	"github.com/go-kratos/kratos/v2/log"
)

type DoctorRecommenderSkill struct {
	vectorACL *acl.VectorDBACL
	doctorACL *acl.DoctorACL
	logger    *log.Helper
}

func NewDoctorRecommenderSkill(vectorACL *acl.VectorDBACL, doctorACL *acl.DoctorACL, logger log.Logger) *DoctorRecommenderSkill {
	return &DoctorRecommenderSkill{
		vectorACL: vectorACL,
		doctorACL: doctorACL,
		logger:    log.NewHelper(logger),
	}
}

func (s *DoctorRecommenderSkill) Name() string {
	return "search_doctor"
}

func (s *DoctorRecommenderSkill) Description() string {
	return "Finds suitable doctors based on symptoms and specialty using vector search and doctor service. Input should include specialty or symptom context."
}

func (s *DoctorRecommenderSkill) Execute(ctx context.Context, input map[string]interface{}) (*agent.SkillResult, error) {
	queryText := s.buildQueryText(input)
	limit := 5
	threshold := float32(0.6)

	if l, ok := input["limit"].(int); ok && l > 0 {
		limit = l
	}

	results, err := s.vectorACL.SearchDoctorsByText(ctx, queryText, limit, threshold)
	if err != nil {
		s.logger.Warnf("Doctor vector search failed: %v", err)
		results = []*acl.VectorSearchResult{}
	}

	specialty, _ := input["specialty"].(string)
	if specialty == "" && len(results) > 0 {
		if spec, ok := results[0].Payload["specialty"].(string); ok {
			specialty = spec
		}
	}

	doctorInfos := make([]*acl.DoctorInfo, 0)

	for _, result := range results {
		doctorID, _ := result.ID
		info, err := s.doctorACL.GetDoctorInfo(ctx, doctorID)
		if err != nil {
			s.logger.Warnf("Failed to get doctor info for %s: %v", doctorID, err)
			continue
		}
		doctorInfos = append(doctorInfos, info)
	}

	if specialty != "" && len(doctorInfos) < limit {
		additionalDoctors, err := s.doctorACL.SearchDoctorsBySpecialty(ctx, specialty, limit-len(doctorInfos))
		if err == nil {
			doctorInfos = append(doctorInfos, additionalDoctors...)
		}
	}

	if len(doctorInfos) == 0 {
		return &agent.SkillResult{
			Success: true,
			Output:  "No matching doctors found. Please consult your local hospital or general practitioner.",
			Data:    map[string]interface{}{"doctors": []interface{}{}},
		}, nil
	}

	output := s.formatDoctorResults(doctorInfos)
	doctorData := s.extractDoctorData(doctorInfos)

	return &agent.SkillResult{
		Success: true,
		Output:  output,
		Data: map[string]interface{}{
			"doctors":      doctorData,
			"result_count": len(doctorInfos),
		},
	}, nil
}

func (s *DoctorRecommenderSkill) buildQueryText(input map[string]interface{}) string {
	queryText := ""

	if symptoms, ok := input["symptoms"].([]agent.SymptomInput); ok {
		for _, sym := range symptoms {
			queryText += sym.BodyPart + " "
		}
	}

	if specialty, ok := input["specialty"].(string); ok {
		queryText += specialty + " "
	}

	if query, ok := input["query"].(string); ok {
		queryText += query
	}

	return queryText
}

func (s *DoctorRecommenderSkill) formatDoctorResults(doctors []*acl.DoctorInfo) string {
	output := "Recommended doctors:\n\n"

	for i, doctor := range doctors {
		output += fmt.Sprintf("%d. %s\n", i+1, doctor.Name)
		output += fmt.Sprintf("   Specialty: %s\n", doctor.Specialty)
		output += fmt.Sprintf("   Hospital: %s\n", doctor.Hospital)
		output += fmt.Sprintf("   Title: %s\n", doctor.Title)
		output += fmt.Sprintf("   Experience: %d years\n", doctor.Experience)
		output += fmt.Sprintf("   Rating: %.1f/5.0\n", doctor.Rating)
		output += fmt.Sprintf("   Availability: %s\n", doctor.Availability)
		output += "\n"
	}

	return output
}

func (s *DoctorRecommenderSkill) extractDoctorData(doctors []*acl.DoctorInfo) []map[string]interface{} {
	data := make([]map[string]interface{}, len(doctors))

	for i, d := range doctors {
		data[i] = map[string]interface{}{
			"id":           d.ID,
			"name":         d.Name,
			"specialty":    d.Specialty,
			"hospital":     d.Hospital,
			"title":        d.Title,
			"experience":   d.Experience,
			"rating":       d.Rating,
			"availability": d.Availability,
		}
	}

	return data
}