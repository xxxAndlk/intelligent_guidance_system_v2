package service

import (
	"context"

	pb "intelligent_guidance_system_v2/service/ai/api"

	"intelligent_guidance_system_v2/service/ai/internal/biz/dto"
	"intelligent_guidance_system_v2/service/ai/internal/biz/usecase"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AIService struct {
	pb.UnimplementedAIServiceServer

	diagnosisUC *usecase.DiagnosisUseCase
	logger      *log.Helper
}

func NewAIService(diagnosisUC *usecase.DiagnosisUseCase, logger log.Logger) *AIService {
	return &AIService{
		diagnosisUC: diagnosisUC,
		logger:      log.NewHelper(logger),
	}
}

func (s *AIService) CreateSession(ctx context.Context, req *pb.CreateSessionRequest) (*pb.CreateSessionResponse, error) {
	dtoReq := &dto.CreateSessionRequest{
		PatientID: req.PatientId,
	}

	resp, err := s.diagnosisUC.CreateSession(ctx, dtoReq)
	if err != nil {
		s.logger.Errorf("Failed to create session: %v", err)
		return nil, err
	}

	return &pb.CreateSessionResponse{
		SessionId: resp.SessionID,
		Status:    resp.Status,
		CreatedAt: timestamppb.New(resp.CreatedAt),
	}, nil
}

func (s *AIService) AddSymptom(ctx context.Context, req *pb.AddSymptomRequest) (*pb.AddSymptomResponse, error) {
	dtoReq := &dto.AddSymptomRequest{
		SessionID:   req.SessionId,
		BodyPart:    req.BodyPart,
		Description: req.Description,
		Severity:    int(req.Severity),
		Duration:    req.Duration,
	}

	resp, err := s.diagnosisUC.AddSymptom(ctx, dtoReq)
	if err != nil {
		s.logger.Errorf("Failed to add symptom: %v", err)
		return nil, err
	}

	return &pb.AddSymptomResponse{
		SymptomId: resp.SymptomID,
		CreatedAt: timestamppb.New(resp.CreatedAt),
	}, nil
}

func (s *AIService) AddConversation(ctx context.Context, req *pb.AddConversationRequest) (*pb.AddConversationResponse, error) {
	dtoReq := &dto.AddConversationRequest{
		SessionID: req.SessionId,
		Role:      req.Role,
		Content:   req.Content,
		Type:      req.Type,
	}

	resp, err := s.diagnosisUC.AddConversation(ctx, dtoReq)
	if err != nil {
		s.logger.Errorf("Failed to add conversation: %v", err)
		return nil, err
	}

	return &pb.AddConversationResponse{
		MessageId: resp.MessageID,
		Sequence:  int32(resp.Sequence),
		CreatedAt: timestamppb.New(resp.CreatedAt),
	}, nil
}

func (s *AIService) GenerateRecommendation(ctx context.Context, req *pb.GenerateRecommendationRequest) (*pb.GenerateRecommendationResponse, error) {
	dtoReq := &dto.GenerateRecommendationRequest{
		SessionID:      req.SessionId,
		TopK:           int(req.TopK),
		ScoreThreshold: req.ScoreThreshold,
		IncludeDrugs:    req.IncludeDrugs,
		IncludeDoctors: req.IncludeDoctors,
	}

	resp, err := s.diagnosisUC.GenerateRecommendation(ctx, dtoReq)
	if err != nil {
		s.logger.Errorf("Failed to generate recommendation: %v", err)
		return nil, err
	}

	return &pb.GenerateRecommendationResponse{
		SessionId:   resp.SessionID,
		Diseases:    mapRecommendationsToProto(resp.Diseases),
		Doctors:     mapRecommendationsToProto(resp.Doctors),
		Drugs:       mapRecommendationsToProto(resp.Drugs),
		Departments: mapRecommendationsToProto(resp.Departments),
		GeneratedAt: timestamppb.New(resp.GeneratedAt),
	}, nil
}

func (s *AIService) CompleteSession(ctx context.Context, req *pb.CompleteSessionRequest) (*pb.CompleteSessionResponse, error) {
	resp, err := s.diagnosisUC.CompleteSession(ctx, req.SessionId)
	if err != nil {
		s.logger.Errorf("Failed to complete session: %v", err)
		return nil, err
	}

	return &pb.CompleteSessionResponse{
		SessionId:   resp.SessionID,
		Status:      resp.Status,
		CompletedAt: timestamppb.New(resp.CompletedAt),
	}, nil
}

func (s *AIService) GetSession(ctx context.Context, req *pb.GetSessionRequest) (*pb.GetSessionResponse, error) {
	resp, err := s.diagnosisUC.GetSession(ctx, req.SessionId)
	if err != nil {
		s.logger.Errorf("Failed to get session: %v", err)
		return nil, err
	}

	return &pb.GetSessionResponse{
		SessionId:     resp.SessionID,
		PatientId:     resp.PatientID,
		Status:        resp.Status,
		Symptoms:      mapSymptomsToProto(resp.Symptoms),
		Conversations: mapConversationsToProto(resp.Conversations),
		CreatedAt:     timestamppb.New(resp.CreatedAt),
		UpdatedAt:     timestamppb.New(resp.UpdatedAt),
	}, nil
}

func (s *AIService) ListSessions(ctx context.Context, req *pb.ListSessionsRequest) (*pb.ListSessionsResponse, error) {
	dtoReq := &dto.ListSessionsRequest{
		PatientID: req.PatientId,
		Limit:     int(req.Limit),
		Offset:    int(req.Offset),
	}

	resp, err := s.diagnosisUC.ListSessions(ctx, dtoReq)
	if err != nil {
		s.logger.Errorf("Failed to list sessions: %v", err)
		return nil, err
	}

	summaries := make([]*pb.SessionSummary, len(resp.Sessions))
	for i, s := range resp.Sessions {
		summary := &pb.SessionSummary{
			SessionId:    s.SessionID,
			PatientId:    s.PatientID,
			Status:       s.Status,
			SymptomCount: int32(s.SymptomCount),
			CreatedAt:    timestamppb.New(s.CreatedAt),
		}
		if s.CompletedAt != nil {
			summary.CompletedAt = timestamppb.New(*s.CompletedAt)
		}
		summaries[i] = summary
	}

	return &pb.ListSessionsResponse{
		Sessions: summaries,
		Total:    int32(resp.Total),
	}, nil
}

func mapRecommendationsToProto(recs []dto.RecommendationDTO) []*pb.Recommendation {
	result := make([]*pb.Recommendation, len(recs))
	for i, r := range recs {
		result[i] = &pb.Recommendation{
			Id:              r.ID,
			Type:            r.Type,
			Title:           r.Title,
			Description:     r.Description,
			ConfidenceScore: r.ConfidenceScore,
			ConfidenceLevel: r.ConfidenceLevel,
			Reasoning:       r.Reasoning,
		}
	}
	return result
}

func mapSymptomsToProto(symptoms []dto.SymptomDTO) []*pb.Symptom {
	result := make([]*pb.Symptom, len(symptoms))
	for i, s := range symptoms {
		result[i] = &pb.Symptom{
			Id:          s.ID,
			BodyPart:    s.BodyPart,
			Description: s.Description,
			Severity:    int32(s.Severity),
			Duration:    s.Duration,
			CreatedAt:   timestamppb.New(s.CreatedAt),
		}
	}
	return result
}

func mapConversationsToProto(conversations []dto.ConversationDTO) []*pb.Conversation {
	result := make([]*pb.Conversation, len(conversations))
	for i, c := range conversations {
		result[i] = &pb.Conversation{
			Id:        c.ID,
			Role:      c.Role,
			Content:   c.Content,
			Type:      c.Type,
			Sequence:  int32(c.Sequence),
			CreatedAt: timestamppb.New(c.CreatedAt),
		}
	}
	return result
}