// Package usecase provides business use cases for the AI diagnosis service.
package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"intelligent_guidance_system_v2/service/ai/internal/domain/aggregate"
	"intelligent_guidance_system_v2/service/ai/internal/domain/entity"
	"intelligent_guidance_system_v2/service/ai/internal/domain/repository"
	"intelligent_guidance_system_v2/service/ai/internal/domain/vo"
	"intelligent_guidance_system_v2/service/ai/internal/biz/dto"

	"github.com/go-kratos/kratos/v2/log"
)

var (
	ErrSessionNotFound      = errors.New("session not found")
	ErrSessionNotActive     = errors.New("session is not active")
	ErrInvalidPatientID     = errors.New("invalid patient ID")
	ErrInvalidSymptomData   = errors.New("invalid symptom data")
	ErrVectorSearchFailed   = errors.New("vector search failed")
	ErrEmbeddingFailed      = errors.New("embedding generation failed")
	ErrLLMGenerationFailed  = errors.New("LLM generation failed")
)

// DiagnosisUseCase implements diagnosis-related business logic.
type DiagnosisUseCase struct {
	sessionRepo   repository.DiagnosisRepository
	vectorRepo    repository.VectorRepository
	embeddingRepo repository.EmbeddingRepository
	logger        *log.Helper
}

// NewDiagnosisUseCase creates a new DiagnosisUseCase.
func NewDiagnosisUseCase(
	sessionRepo repository.DiagnosisRepository,
	vectorRepo repository.VectorRepository,
	embeddingRepo repository.EmbeddingRepository,
	logger log.Logger,
) *DiagnosisUseCase {
	return &DiagnosisUseCase{
		sessionRepo:   sessionRepo,
		vectorRepo:    vectorRepo,
		embeddingRepo: embeddingRepo,
		logger:        log.NewHelper(logger),
	}
}

// CreateSession creates a new diagnosis session for a patient.
func (uc *DiagnosisUseCase) CreateSession(ctx context.Context, req *dto.CreateSessionRequest) (*dto.CreateSessionResponse, error) {
	if req.PatientID == "" {
		return nil, ErrInvalidPatientID
	}

	session, err := aggregate.NewDiagnosisSession(req.PatientID)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	if err := uc.sessionRepo.Save(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to save session: %w", err)
	}

	uc.logger.Infof("Created diagnosis session %s for patient %s", session.ID(), req.PatientID)

	return &dto.CreateSessionResponse{
		SessionID: session.ID(),
		Status:    string(session.Status()),
		CreatedAt: session.CreatedAt(),
	}, nil
}

// AddSymptom adds a symptom to an existing session.
func (uc *DiagnosisUseCase) AddSymptom(ctx context.Context, req *dto.AddSymptomRequest) (*dto.AddSymptomResponse, error) {
	if req.SessionID == "" {
		return nil, ErrSessionNotFound
	}

	if req.BodyPart == "" || req.Description == "" {
		return nil, ErrInvalidSymptomData
	}

	session, err := uc.sessionRepo.FindByID(ctx, req.SessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find session: %w", err)
	}
	if session == nil {
		return nil, ErrSessionNotFound
	}

	if err := session.AddSymptom(req.BodyPart, req.Description, req.Severity); err != nil {
		return nil, fmt.Errorf("failed to add symptom: %w", err)
	}

	if req.Duration != "" {
		symptoms := session.Symptoms()
		if len(symptoms) > 0 {
			symptoms[len(symptoms)-1].SetDuration(req.Duration)
		}
	}

	if err := uc.sessionRepo.Save(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to save session: %w", err)
	}

	symptoms := session.Symptoms()
	latestSymptom := symptoms[len(symptoms)-1]

	uc.logger.Infof("Added symptom %s to session %s", latestSymptom.ID, session.ID())

	return &dto.AddSymptomResponse{
		SymptomID: latestSymptom.ID,
		CreatedAt: latestSymptom.CreatedAt,
	}, nil
}

// AddConversation adds a conversation message to a session.
func (uc *DiagnosisUseCase) AddConversation(ctx context.Context, req *dto.AddConversationRequest) (*dto.AddConversationResponse, error) {
	if req.SessionID == "" {
		return nil, ErrSessionNotFound
	}

	session, err := uc.sessionRepo.FindByID(ctx, req.SessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find session: %w", err)
	}
	if session == nil {
		return nil, ErrSessionNotFound
	}

	role := entity.ConversationRole(req.Role)
	msgType := entity.MessageType(req.Type)

	if err := session.AddConversation(role, req.Content, msgType); err != nil {
		return nil, fmt.Errorf("failed to add conversation: %w", err)
	}

	if err := uc.sessionRepo.Save(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to save session: %w", err)
	}

	conversations := session.Conversations()
	latestMsg := conversations[len(conversations)-1]

	return &dto.AddConversationResponse{
		MessageID: latestMsg.ID,
		Sequence:  latestMsg.Sequence,
		CreatedAt: latestMsg.CreatedAt,
	}, nil
}

// GenerateRecommendation generates AI-powered recommendations for a session.
func (uc *DiagnosisUseCase) GenerateRecommendation(ctx context.Context, req *dto.GenerateRecommendationRequest) (*dto.GenerateRecommendationResponse, error) {
	if req.SessionID == "" {
		return nil, ErrSessionNotFound
	}

	session, err := uc.sessionRepo.FindByID(ctx, req.SessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find session: %w", err)
	}
	if session == nil {
		return nil, ErrSessionNotFound
	}

	if !session.IsReadyForAnalysis() {
		return nil, ErrSessionNotActive
	}

	if err := session.StartAnalysis(); err != nil {
		return nil, fmt.Errorf("failed to start analysis: %w", err)
	}

	contextStr := session.GetContextForRAG()
	embedding, err := uc.embeddingRepo.Embed(ctx, contextStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrEmbeddingFailed, err)
	}

	topK := req.TopK
	if topK <= 0 {
		topK = 5
	}
	threshold := req.ScoreThreshold
	if threshold <= 0 {
		threshold = 0.7
	}

	diseaseResults, err := uc.vectorRepo.SearchDiseases(ctx, embedding, topK, threshold)
	if err != nil {
		uc.logger.Warnf("Disease search failed: %v", err)
	}

	for _, result := range diseaseResults {
		rec, err := uc.createRecommendationFromResult(result, vo.RecommendationTypeDisease)
		if err == nil {
			session.GenerateRecommendation(rec)
		}
	}

	if req.IncludeDoctors {
		doctorResults, err := uc.vectorRepo.SearchDoctors(ctx, embedding, topK, threshold)
		if err != nil {
			uc.logger.Warnf("Doctor search failed: %v", err)
		}
		for _, result := range doctorResults {
			rec, err := uc.createRecommendationFromResult(result, vo.RecommendationTypeDoctor)
			if err == nil {
				session.GenerateRecommendation(rec)
			}
		}
	}

	if req.IncludeDrugs {
		drugResults, err := uc.vectorRepo.SearchDrugs(ctx, embedding, topK, threshold)
		if err != nil {
			uc.logger.Warnf("Drug search failed: %v", err)
		}
		for _, result := range drugResults {
			rec, err := uc.createRecommendationFromResult(result, vo.RecommendationTypeDrug)
			if err == nil {
				session.GenerateRecommendation(rec)
			}
		}
	}

	session.ClearDomainEvents()

	if err := uc.sessionRepo.Save(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to save session: %w", err)
	}

	return uc.buildRecommendationResponse(session), nil
}

// CompleteSession marks a session as completed.
func (uc *DiagnosisUseCase) CompleteSession(ctx context.Context, sessionID string) (*dto.CompleteSessionResponse, error) {
	if sessionID == "" {
		return nil, ErrSessionNotFound
	}

	session, err := uc.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find session: %w", err)
	}
	if session == nil {
		return nil, ErrSessionNotFound
	}

	if err := session.Complete(); err != nil {
		return nil, fmt.Errorf("failed to complete session: %w", err)
	}

	if err := uc.sessionRepo.Save(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to save session: %w", err)
	}

	uc.logger.Infof("Completed session %s", sessionID)

	return &dto.CompleteSessionResponse{
		SessionID:   session.ID(),
		Status:      string(session.Status()),
		CompletedAt: *session.CompletedAt(),
	}, nil
}

// GetSession retrieves session details.
func (uc *DiagnosisUseCase) GetSession(ctx context.Context, sessionID string) (*dto.GetSessionResponse, error) {
	if sessionID == "" {
		return nil, ErrSessionNotFound
	}

	session, err := uc.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find session: %w", err)
	}
	if session == nil {
		return nil, ErrSessionNotFound
	}

	return uc.buildSessionResponse(session), nil
}

// ListSessions lists sessions for a patient.
func (uc *DiagnosisUseCase) ListSessions(ctx context.Context, req *dto.ListSessionsRequest) (*dto.ListSessionsResponse, error) {
	if req.PatientID == "" {
		return nil, ErrInvalidPatientID
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	sessions, err := uc.sessionRepo.FindByPatientID(ctx, req.PatientID, limit, req.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}

	summaries := make([]dto.SessionSummaryDTO, len(sessions))
	for i, s := range sessions {
		summaries[i] = dto.SessionSummaryDTO{
			SessionID:    s.ID(),
			PatientID:    s.PatientID(),
			Status:       string(s.Status()),
			SymptomCount: len(s.Symptoms()),
			CreatedAt:    s.CreatedAt(),
			CompletedAt:  s.CompletedAt(),
		}
	}

	return &dto.ListSessionsResponse{
		Sessions: summaries,
		Total:    len(summaries),
	}, nil
}

// GetActiveSession retrieves the active session for a patient.
func (uc *DiagnosisUseCase) GetActiveSession(ctx context.Context, patientID string) (*dto.GetSessionResponse, error) {
	if patientID == "" {
		return nil, ErrInvalidPatientID
	}

	session, err := uc.sessionRepo.FindActiveByPatientID(ctx, patientID)
	if err != nil {
		return nil, fmt.Errorf("failed to find active session: %w", err)
	}
	if session == nil {
		return nil, ErrSessionNotFound
	}

	return uc.buildSessionResponse(session), nil
}

func (uc *DiagnosisUseCase) createRecommendationFromResult(result *entity.VectorSearchResult, recType vo.RecommendationType) (*vo.Recommendation, error) {
	title, _ := result.Payload["name"].(string)
	if title == "" {
		title, _ = result.Payload["title"].(string)
	}
	
	description, _ := result.Payload["description"].(string)
	if description == "" {
		description, _ = result.Payload["content"].(string)
	}

	rec, err := vo.NewRecommendation(recType, title, description, result.Score)
	if err != nil {
		return nil, err
	}

	if id, ok := result.Payload["id"].(string); ok {
		rec = rec.WithSourceIDs([]string{id})
	}

	return rec, nil
}

func (uc *DiagnosisUseCase) buildRecommendationResponse(session aggregate.DiagnosisSessionInterface) *dto.GenerateRecommendationResponse {
	rec := session.Recommendations()
	
	return &dto.GenerateRecommendationResponse{
		SessionID:   session.ID(),
		Diseases:    uc.mapRecommendations(rec.Diseases()),
		Doctors:     uc.mapRecommendations(rec.Doctors()),
		Drugs:       uc.mapRecommendations(rec.Drugs()),
		Departments: uc.mapRecommendations(rec.Departments()),
		GeneratedAt: session.UpdatedAt(),
	}
}

func (uc *DiagnosisUseCase) mapRecommendations(recs []*vo.Recommendation) []dto.RecommendationDTO {
	result := make([]dto.RecommendationDTO, len(recs))
	for i, r := range recs {
		result[i] = dto.RecommendationDTO{
			ID:              r.ID(),
			Type:            string(r.Type()),
			Title:           r.Title(),
			Description:     r.Description(),
			ConfidenceScore: r.ConfidenceScore(),
			ConfidenceLevel: string(r.ConfidenceLevel()),
			Reasoning:       r.Reasoning(),
			Metadata:        r.Metadata(),
		}
	}
	return result
}

func (uc *DiagnosisUseCase) buildSessionResponse(session *aggregate.DiagnosisSession) *dto.GetSessionResponse {
	symptoms := make([]dto.SymptomDTO, len(session.Symptoms()))
	for i, s := range session.Symptoms() {
		symptoms[i] = dto.SymptomDTO{
			ID:          s.ID,
			BodyPart:    s.BodyPart,
			Description: s.Description,
			Severity:    s.Severity,
			Duration:    s.Duration,
			Metadata:    s.Metadata,
			CreatedAt:   s.CreatedAt,
		}
	}

	conversations := make([]dto.ConversationDTO, len(session.Conversations()))
	for i, c := range session.Conversations() {
		conversations[i] = dto.ConversationDTO{
			ID:        c.ID,
			Role:      string(c.Role),
			Content:   c.Content,
			Type:      string(c.Type),
			Sequence:  c.Sequence,
			CreatedAt: c.CreatedAt,
		}
	}

	var recSetDTO *dto.RecommendationSetDTO
	if rec := session.Recommendations(); rec != nil && !rec.IsEmpty() {
		recSetDTO = &dto.RecommendationSetDTO{
			Diseases:    uc.mapRecommendations(rec.Diseases()),
			Doctors:     uc.mapRecommendations(rec.Doctors()),
			Drugs:       uc.mapRecommendations(rec.Drugs()),
			Departments: uc.mapRecommendations(rec.Departments()),
		}
	}

	return &dto.GetSessionResponse{
		SessionID:       session.ID(),
		PatientID:       session.PatientID(),
		Status:          string(session.Status()),
		Symptoms:        symptoms,
		Conversations:   conversations,
		Recommendations: recSetDTO,
		CreatedAt:       session.CreatedAt(),
		UpdatedAt:       session.UpdatedAt(),
		CompletedAt:     session.CompletedAt(),
	}
}

type DiagnosisSessionInterface interface {
	ID() string
	Status() aggregate.SessionStatus
	Recommendations() *vo.RecommendationSet
	UpdatedAt() time.Time
	ClearDomainEvents()
}