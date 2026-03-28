// Package mysql provides MySQL-based repository implementations.
package mysql

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"intelligent_guidance_system_v2/service/ai/internal/domain/aggregate"
	"intelligent_guidance_system_v2/service/ai/internal/domain/entity"
	"intelligent_guidance_system_v2/service/ai/internal/domain/repository"
	"intelligent_guidance_system_v2/service/ai/internal/domain/vo"

	"gorm.io/gorm"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	ErrDatabaseConnection = errors.New("database connection error")
	ErrQueryFailed        = errors.New("query failed")
	ErrRecordNotFound     = errors.New("record not found")
)

type DiagnosisSessionModel struct {
	ID             string    `gorm:"primaryKey;type:varchar(36)"`
	PatientID      string    `gorm:"index:idx_patient_status;type:varchar(36)"`
	Status         string    `gorm:"index:idx_patient_status;type:varchar(20)"`
	SymptomData    string    `gorm:"type:json"`
	ConversationData string `gorm:"type:json"`
	RecommendationData string `gorm:"type:json"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
	CompletedAt    *time.Time
}

func (DiagnosisSessionModel) TableName() string {
	return "ai_diagnosis_sessions"
}

type DiagnosisRepoImpl struct {
	db     *gorm.DB
	logger *log.Helper
}

func NewDiagnosisRepoImpl(db *gorm.DB, logger log.Logger) repository.DiagnosisRepository {
	return &DiagnosisRepoImpl{
		db:     db,
		logger: log.NewHelper(logger),
	}
}

func (r *DiagnosisRepoImpl) Save(ctx context.Context, session *aggregate.DiagnosisSession) error {
	model := r.toModel(session)
	
	result := r.db.WithContext(ctx).Save(&model)
	if result.Error != nil {
		r.logger.Errorf("Failed to save session %s: %v", session.ID(), result.Error)
		return fmt.Errorf("%w: %v", ErrQueryFailed, result.Error)
	}
	
	return nil
}

func (r *DiagnosisRepoImpl) FindByID(ctx context.Context, sessionID string) (*aggregate.DiagnosisSession, error) {
	var model DiagnosisSessionModel
	
	result := r.db.WithContext(ctx).Where("id = ?", sessionID).First(&model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w: %v", ErrQueryFailed, result.Error)
	}
	
	return r.toAggregate(&model)
}

func (r *DiagnosisRepoImpl) FindByPatientID(ctx context.Context, patientID string, limit, offset int) ([]*aggregate.DiagnosisSession, error) {
	var models []DiagnosisSessionModel
	
	result := r.db.WithContext(ctx).
		Where("patient_id = ?", patientID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&models)
	
	if result.Error != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryFailed, result.Error)
	}
	
	sessions := make([]*aggregate.DiagnosisSession, len(models))
	for i, m := range models {
		session, err := r.toAggregate(&m)
		if err != nil {
			r.logger.Warnf("Failed to convert model to aggregate: %v", err)
			continue
		}
		sessions[i] = session
	}
	
	return sessions, nil
}

func (r *DiagnosisRepoImpl) FindActiveByPatientID(ctx context.Context, patientID string) (*aggregate.DiagnosisSession, error) {
	var model DiagnosisSessionModel
	
	result := r.db.WithContext(ctx).
		Where("patient_id = ? AND status = ?", patientID, aggregate.StatusActive).
		Order("created_at DESC").
		First(&model)
	
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w: %v", ErrQueryFailed, result.Error)
	}
	
	return r.toAggregate(&model)
}

func (r *DiagnosisRepoImpl) Delete(ctx context.Context, sessionID string) error {
	result := r.db.WithContext(ctx).Delete(&DiagnosisSessionModel{}, "id = ?", sessionID)
	if result.Error != nil {
		return fmt.Errorf("%w: %v", ErrQueryFailed, result.Error)
	}
	return nil
}

func (r *DiagnosisRepoImpl) UpdateStatus(ctx context.Context, sessionID string, status aggregate.SessionStatus) error {
	result := r.db.WithContext(ctx).
		Model(&DiagnosisSessionModel{}).
		Where("id = ?", sessionID).
		Update("status", string(status))
	
	if result.Error != nil {
		return fmt.Errorf("%w: %v", ErrQueryFailed, result.Error)
	}
	
	return nil
}

func (r *DiagnosisRepoImpl) toModel(session *aggregate.DiagnosisSession) DiagnosisSessionModel {
	symptomData, _ := json.Marshal(session.Symptoms())
	conversationData, _ := json.Marshal(session.Conversations())
	recommendationData := r.serializeRecommendations(session.Recommendations())
	
	return DiagnosisSessionModel{
		ID:               session.ID(),
		PatientID:        session.PatientID(),
		Status:           string(session.Status()),
		SymptomData:      string(symptomData),
		ConversationData: string(conversationData),
		RecommendationData: string(recommendationData),
		CreatedAt:        session.CreatedAt(),
		UpdatedAt:        session.UpdatedAt(),
		CompletedAt:      session.CompletedAt(),
	}
}

func (r *DiagnosisRepoImpl) toAggregate(model *DiagnosisSessionModel) (*aggregate.DiagnosisSession, error) {
	symptoms := make([]*entity.Symptom, 0)
	if model.SymptomData != "" {
		if err := json.Unmarshal([]byte(model.SymptomData), &symptoms); err != nil {
			r.logger.Warnf("Failed to unmarshal symptoms: %v", err)
		}
	}
	
	conversations := make([]*entity.Conversation, 0)
	if model.ConversationData != "" {
		if err := json.Unmarshal([]byte(model.ConversationData), &conversations); err != nil {
			r.logger.Warnf("Failed to unmarshal conversations: %v", err)
		}
	}
	
	recommendations := r.deserializeRecommendations(model.RecommendationData)
	
	return aggregate.Reconstruct(
		model.ID,
		model.PatientID,
		aggregate.SessionStatus(model.Status),
		symptoms,
		conversations,
		recommendations,
		model.CreatedAt,
		model.UpdatedAt,
		model.CompletedAt,
	), nil
}

func (r *DiagnosisRepoImpl) serializeRecommendations(rs *vo.RecommendationSet) []byte {
	data := struct {
		Diseases    []recommendationJSON `json:"diseases"`
		Doctors     []recommendationJSON `json:"doctors"`
		Drugs       []recommendationJSON `json:"drugs"`
		Departments []recommendationJSON `json:"departments"`
	}{
		Diseases:    r.mapRecsToJSON(rs.Diseases()),
		Doctors:     r.mapRecsToJSON(rs.Doctors()),
		Drugs:       r.mapRecsToJSON(rs.Drugs()),
		Departments: r.mapRecsToJSON(rs.Departments()),
	}
	
	result, _ := json.Marshal(data)
	return result
}

type recommendationJSON struct {
	ID              string                 `json:"id"`
	Type            string                 `json:"type"`
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	ConfidenceScore float32                `json:"confidence_score"`
	SourceIDs       []string               `json:"source_ids"`
	Reasoning       string                 `json:"reasoning"`
	Metadata        map[string]interface{} `json:"metadata"`
	CreatedAt       time.Time              `json:"created_at"`
}

func (r *DiagnosisRepoImpl) mapRecsToJSON(recs []*vo.Recommendation) []recommendationJSON {
	result := make([]recommendationJSON, len(recs))
	for i, rec := range recs {
		result[i] = recommendationJSON{
			ID:              rec.ID(),
			Type:            string(rec.Type()),
			Title:           rec.Title(),
			Description:     rec.Description(),
			ConfidenceScore: rec.ConfidenceScore(),
			SourceIDs:       rec.SourceIDs(),
			Reasoning:       rec.Reasoning(),
			Metadata:        rec.Metadata(),
			CreatedAt:       rec.CreatedAt(),
		}
	}
	return result
}

func (r *DiagnosisRepoImpl) deserializeRecommendations(data string) *vo.RecommendationSet {
	rs := vo.NewRecommendationSet()
	
	if data == "" {
		return rs
	}
	
	var parsed struct {
		Diseases    []recommendationJSON `json:"diseases"`
		Doctors     []recommendationJSON `json:"doctors"`
		Drugs       []recommendationJSON `json:"drugs"`
		Departments []recommendationJSON `json:"departments"`
	}
	
	if err := json.Unmarshal([]byte(data), &parsed); err != nil {
		r.logger.Warnf("Failed to unmarshal recommendations: %v", err)
		return rs
	}
	
	for _, d := range parsed.Diseases {
		rec, err := vo.NewRecommendation(vo.RecommendationTypeDisease, d.Title, d.Description, d.ConfidenceScore)
		if err == nil {
			rec = rec.WithSourceIDs(d.SourceIDs).WithReasoning(d.Reasoning)
			rs.AddDisease(rec)
		}
	}
	
	for _, d := range parsed.Doctors {
		rec, err := vo.NewRecommendation(vo.RecommendationTypeDoctor, d.Title, d.Description, d.ConfidenceScore)
		if err == nil {
			rec = rec.WithSourceIDs(d.SourceIDs).WithReasoning(d.Reasoning)
			rs.AddDoctor(rec)
		}
	}
	
	for _, d := range parsed.Drugs {
		rec, err := vo.NewRecommendation(vo.RecommendationTypeDrug, d.Title, d.Description, d.ConfidenceScore)
		if err == nil {
			rec = rec.WithSourceIDs(d.SourceIDs).WithReasoning(d.Reasoning)
			rs.AddDrug(rec)
		}
	}
	
	for _, d := range parsed.Departments {
		rec, err := vo.NewRecommendation(vo.RecommendationTypeDepartment, d.Title, d.Description, d.ConfidenceScore)
		if err == nil {
			rec = rec.WithSourceIDs(d.SourceIDs).WithReasoning(d.Reasoning)
			rs.AddDepartment(rec)
		}
	}
	
	return rs
}