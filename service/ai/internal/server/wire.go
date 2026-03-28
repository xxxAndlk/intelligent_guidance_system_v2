//go:build wireinject

package server

import (
	"intelligent_guidance_system_v2/service/ai/internal/acl"
	"intelligent_guidance_system_v2/service/ai/internal/agent"
	"intelligent_guidance_system_v2/service/ai/internal/agent/skills"
	"intelligent_guidance_system_v2/service/ai/internal/biz/usecase"
	"intelligent_guidance_system_v2/service/ai/internal/data/llm"
	"intelligent_guidance_system_v2/service/ai/internal/data/mysql"
	"intelligent_guidance_system_v2/service/ai/internal/data/vector"
	"intelligent_guidance_system_v2/service/ai/internal/domain/repository"
	"intelligent_guidance_system_v2/service/ai/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitApp(cfg *Config, db *gorm.DB, logger log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		provideDB,
		provideLogger,
		provideDiagnosisRepo,
		provideVectorRepo,
		provideLLMClient,
		provideEmbeddingRepo,
		providePatientACL,
		provideDoctorACL,
		provideDiseaseACL,
		provideVectorDBACL,
		provideReActAgent,
		provideSymptomCollectorSkill,
		provideDiseaseRecommenderSkill,
		provideDoctorRecommenderSkill,
		provideDrugRecommenderSkill,
		provideDiagnosisUseCase,
		provideAIService,
		provideAIServiceHTTP,
		provideGRPCServer,
		provideHTTPServer,
		provideApp,
	))
}

func provideDB(cfg *Config) *gorm.DB {
	dsn := cfg.Database.DSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func provideLogger() log.Logger {
	return log.DefaultLogger
}

func provideDiagnosisRepo(db *gorm.DB, logger log.Logger) repository.DiagnosisRepository {
	return mysql.NewDiagnosisRepoImpl(db, logger)
}

func provideVectorRepo(cfg *Config, logger log.Logger) repository.VectorRepository {
	qdrantCfg := vector.QdrantConfig{
		Host:            cfg.Qdrant.Host,
		Port:            cfg.Qdrant.Port,
		DiseaseCollection: cfg.Qdrant.DiseaseCollection,
		DoctorCollection:  cfg.Qdrant.DoctorCollection,
		DrugCollection:    cfg.Qdrant.DrugCollection,
	}
	repo, err := vector.NewQdrantRepo(qdrantCfg, logger)
	if err != nil {
		panic(err)
	}
	return repo
}

func provideLLMClient(cfg *Config, logger log.Logger) *llm.DeepSeekClient {
	llmCfg := llm.DeepSeekConfig{
		BaseURL:     cfg.LLM.BaseURL,
		APIKey:      cfg.LLM.APIKey,
		Model:       cfg.LLM.Model,
		EmbedModel:  cfg.LLM.EmbedModel,
		MaxTokens:   cfg.LLM.MaxTokens,
		Temperature: cfg.LLM.Temperature,
		Timeout:     cfg.LLM.Timeout,
	}
	client, err := llm.NewDeepSeekClient(llmCfg, logger)
	if err != nil {
		panic(err)
	}
	return client
}

func provideEmbeddingRepo(llmClient *llm.DeepSeekClient) repository.EmbeddingRepository {
	return llm.NewEmbeddingRepoAdapter(llmClient)
}

func providePatientACL(logger log.Logger) *acl.PatientACL {
	return acl.NewPatientACL(&acl.MockPatientServiceClient{}, logger)
}

func provideDoctorACL(logger log.Logger) *acl.DoctorACL {
	return acl.NewDoctorACL(&acl.MockDoctorServiceClient{}, logger)
}

func provideDiseaseACL(logger log.Logger) *acl.DiseaseACL {
	return acl.NewDiseaseACL(&acl.MockDiseaseServiceClient{}, logger)
}

func provideVectorDBACL(
	vectorRepo repository.VectorRepository,
	embeddingRepo repository.EmbeddingRepository,
	logger log.Logger,
) *acl.VectorDBACL {
	return acl.NewVectorDBACL(vectorRepo, embeddingRepo, logger)
}

func provideSymptomCollectorSkill(logger log.Logger) *skills.SymptomCollectorSkill {
	return skills.NewSymptomCollectorSkill(logger)
}

func provideDiseaseRecommenderSkill(vectorACL *acl.VectorDBACL, logger log.Logger) *skills.DiseaseRecommenderSkill {
	return skills.NewDiseaseRecommenderSkill(vectorACL, logger)
}

func provideDoctorRecommenderSkill(vectorACL *acl.VectorDBACL, doctorACL *acl.DoctorACL, logger log.Logger) *skills.DoctorRecommenderSkill {
	return skills.NewDoctorRecommenderSkill(vectorACL, doctorACL, logger)
}

func provideDrugRecommenderSkill(vectorACL *acl.VectorDBACL, logger log.Logger) *skills.DrugRecommenderSkill {
	return skills.NewDrugRecommenderSkill(vectorACL, logger)
}

func provideReActAgent(
	llmClient *llm.DeepSeekClient,
	symptomSkill *skills.SymptomCollectorSkill,
	diseaseSkill *skills.DiseaseRecommenderSkill,
	doctorSkill *skills.DoctorRecommenderSkill,
	drugSkill *skills.DrugRecommenderSkill,
	logger log.Logger,
) agent.Agent {
	agentCfg := agent.AgentConfig{
		MaxIterations: 10,
		Temperature:   0.7,
		Model:         "deepseek-chat",
	}

	reactAgent := agent.NewReActAgent(llmClient, agentCfg, logger)
	reactAgent.AddSkill(symptomSkill)
	reactAgent.AddSkill(diseaseSkill)
	reactAgent.AddSkill(doctorSkill)
	reactAgent.AddSkill(drugSkill)

	return reactAgent
}

func provideDiagnosisUseCase(
	sessionRepo repository.DiagnosisRepository,
	vectorRepo repository.VectorRepository,
	embeddingRepo repository.EmbeddingRepository,
	logger log.Logger,
) *usecase.DiagnosisUseCase {
	return usecase.NewDiagnosisUseCase(sessionRepo, vectorRepo, embeddingRepo, logger)
}

func provideAIService(diagnosisUC *usecase.DiagnosisUseCase, logger log.Logger) *service.AIService {
	return service.NewAIService(diagnosisUC, logger)
}

func provideAIServiceHTTP(diagnosisUC *usecase.DiagnosisUseCase, reactAgent agent.Agent, logger log.Logger) *service.AIServiceHTTP {
	return service.NewAIServiceHTTP(diagnosisUC, reactAgent, logger)
}

func provideGRPCServer(cfg *Config, aiSvc *service.AIService, logger log.Logger) *grpc.Server {
	return NewGRPCServer(cfg, aiSvc, logger)
}

func provideHTTPServer(cfg *Config, aiSvcHTTP *service.AIServiceHTTP, logger log.Logger) *http.Server {
	return NewHTTPServer(cfg, aiSvcHTTP, logger)
}

func provideApp(grpcSrv *grpc.Server, httpSrv *http.Server, logger log.Logger) *kratos.App {
	return kratos.New(
		kratos.Server(grpcSrv, httpSrv),
		kratos.Logger(logger),
	)
}

func provideConfig() *Config {
	cfg, _ := LoadConfig("configs/config.yaml")
	return cfg
}