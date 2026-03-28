package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"intelligent_guidance_system_v2/service/ai/internal/acl"
	"intelligent_guidance_system_v2/service/ai/internal/agent"
	"intelligent_guidance_system_v2/service/ai/internal/agent/skills"
	"intelligent_guidance_system_v2/service/ai/internal/biz/usecase"
	"intelligent_guidance_system_v2/service/ai/internal/data/llm"
	"intelligent_guidance_system_v2/service/ai/internal/data/mysql"
	"intelligent_guidance_system_v2/service/ai/internal/data/vector"
	"intelligent_guidance_system_v2/service/ai/internal/domain/repository"
	"intelligent_guidance_system_v2/service/ai/internal/server"
	"intelligent_guidance_system_v2/service/ai/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	configPath = flag.String("config", "configs/config.yaml", "config file path")
)

func main() {
	flag.Parse()

	logger := log.NewStdLogger(os.Stdout)
	logHelper := log.NewHelper(logger)

	logHelper.Info("Starting AI Diagnosis Service...")

	cfg, err := loadConfig(*configPath)
	if err != nil {
		logHelper.Fatalf("Failed to load config: %v", err)
	}

	db, err := initDatabase(cfg, logger)
	if err != nil {
		logHelper.Fatalf("Failed to init database: %v", err)
	}

	vectorRepo, err := initVectorRepo(cfg, logger)
	if err != nil {
		logHelper.Warnf("Failed to init vector repo: %v", err)
		vectorRepo = nil
	}

	llmClient, err := initLLMClient(cfg, logger)
	if err != nil {
		logHelper.Warnf("Failed to init LLM client: %v", err)
		llmClient = nil
	}

	embeddingRepo := initEmbeddingRepo(llmClient)

	reactAgent := initAgent(llmClient, logger)

	diagnosisUC := initDiagnosisUseCase(db, vectorRepo, embeddingRepo, logger)

	aiSvc := initAIService(diagnosisUC, logger)
	aiSvcHTTP := initAIServiceHTTP(diagnosisUC, reactAgent, logger)

	grpcSrv := server.NewGRPCServer(cfg, aiSvc, logger)
	httpSrv := server.NewHTTPServer(cfg, aiSvcHTTP, logger)

	app := kratos.New(
		kratos.Server(grpcSrv, httpSrv),
		kratos.Logger(logger),
		kratos.BeforeStart(func(ctx context.Context) error {
			logHelper.Info("AI Diagnosis Service is starting...")
			return nil
		}),
		kratos.AfterStart(func(ctx context.Context) error {
			logHelper.Info("AI Diagnosis Service started successfully")
			logHelper.Infof("HTTP server listening on %s", cfg.Server.HTTP.Addr)
			logHelper.Infof("GRPC server listening on %s", cfg.Server.GRPC.Addr)
			return nil
		}),
		kratos.BeforeStop(func(ctx context.Context) error {
			logHelper.Info("AI Diagnosis Service is stopping...")
			return nil
		}),
		kratos.AfterStop(func(ctx context.Context) error {
			logHelper.Info("AI Diagnosis Service stopped")
			return nil
		}),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		logHelper.Infof("Received signal: %v", sig)
		cancel()
	}()

	if err := app.Run(ctx); err != nil {
		logHelper.Fatalf("App run failed: %v", err)
	}
}

func loadConfig(path string) (*server.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := &server.Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, nil
}

func initDatabase(cfg *server.Config, logger log.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	log.NewHelper(logger).Infof("Connected to database: %s", cfg.Database.Database)
	return db, nil
}

func initVectorRepo(cfg *server.Config, logger log.Logger) (repository.VectorRepository, error) {
	qdrantCfg := vector.QdrantConfig{
		Host:              cfg.Qdrant.Host,
		Port:              cfg.Qdrant.Port,
		DiseaseCollection: cfg.Qdrant.DiseaseCollection,
		DoctorCollection:  cfg.Qdrant.DoctorCollection,
		DrugCollection:     cfg.Qdrant.DrugCollection,
	}

	repo, err := vector.NewQdrantRepo(qdrantCfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to connect Qdrant: %w", err)
	}

	log.NewHelper(logger).Infof("Connected to Qdrant at %s:%d", cfg.Qdrant.Host, cfg.Qdrant.Port)
	return repo, nil
}

func initLLMClient(cfg *server.Config, logger log.Logger) (*llm.DeepSeekClient, error) {
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
		return nil, fmt.Errorf("failed to init LLM client: %w", err)
	}

	log.NewHelper(logger).Infof("Initialized LLM client with model: %s", cfg.LLM.Model)
	return client, nil
}

func initEmbeddingRepo(llmClient *llm.DeepSeekClient) repository.EmbeddingRepository {
	if llmClient == nil {
		return nil
	}
	return llm.NewEmbeddingRepoAdapter(llmClient)
}

func initAgent(llmClient *llm.DeepSeekClient, logger log.Logger) agent.Agent {
	agentCfg := agent.AgentConfig{
		MaxIterations: 10,
		Temperature:   0.7,
		Model:         "deepseek-chat",
	}

	if llmClient == nil {
		log.NewHelper(logger).Warn("LLM client not initialized, agent will not function properly")
		return nil
	}

	reactAgent := agent.NewReActAgent(llmClient, agentCfg, logger)

	vectorACL := acl.NewVectorDBACL(nil, nil, logger)
	doctorACL := acl.NewDoctorACL(&acl.MockDoctorServiceClient{}, logger)

	symptomSkill := skills.NewSymptomCollectorSkill(logger)
	diseaseSkill := skills.NewDiseaseRecommenderSkill(vectorACL, logger)
	doctorSkill := skills.NewDoctorRecommenderSkill(vectorACL, doctorACL, logger)
	drugSkill := skills.NewDrugRecommenderSkill(vectorACL, logger)

	reactAgent.AddSkill(symptomSkill)
	reactAgent.AddSkill(diseaseSkill)
	reactAgent.AddSkill(doctorSkill)
	reactAgent.AddSkill(drugSkill)

	log.NewHelper(logger).Info("Initialized ReAct agent with skills")
	return reactAgent
}

func initDiagnosisUseCase(
	db *gorm.DB,
	vectorRepo repository.VectorRepository,
	embeddingRepo repository.EmbeddingRepository,
	logger log.Logger,
) *usecase.DiagnosisUseCase {
	sessionRepo := mysql.NewDiagnosisRepoImpl(db, logger)
	return usecase.NewDiagnosisUseCase(sessionRepo, vectorRepo, embeddingRepo, logger)
}

func initAIService(diagnosisUC *usecase.DiagnosisUseCase, logger log.Logger) *service.AIService {
	return service.NewAIService(diagnosisUC, logger)
}

func initAIServiceHTTP(diagnosisUC *usecase.DiagnosisUseCase, reactAgent agent.Agent, logger log.Logger) *service.AIServiceHTTP {
	return service.NewAIServiceHTTP(diagnosisUC, reactAgent, logger)
}