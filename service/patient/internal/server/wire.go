// service/patient/internal/server/wire.go
package server

import (
	"time"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	"gorm.io/gorm"

	"github.com/redis/go-redis/v9"

	"intelligent-guidance-system/service/patient/internal/biz/usecase"
	"intelligent-guidance-system/service/patient/internal/data/mysql"
	"intelligent-guidance-system/service/patient/internal/data/redis"
	"intelligent-guidance-system/service/patient/internal/domain/repository"
	"intelligent-guidance-system/service/patient/internal/service"
)

type Config struct {
	GRPC struct {
		Network string
		Addr    string
		Timeout time.Duration
	}
	HTTP struct {
		Network string
		Addr    string
		Timeout time.Duration
	}
	JWT struct {
		Key string
	}
}

func NewConfig() *Config {
	return &Config{
		GRPC: struct {
			Network string
			Addr    string
			Timeout time.Duration
		}{
			Network: "tcp",
			Addr:    ":9000",
			Timeout: 30 * time.Second,
		},
		HTTP: struct {
			Network string
			Addr    string
			Timeout time.Duration
		}{
			Network: "tcp",
			Addr:    ":8000",
			Timeout: 30 * time.Second,
		},
		JWT: struct {
			Key string
		}{
			Key: "patient-service-secret-key",
		},
	}
}

var ProviderSet = wire.NewSet(
	NewConfig,
	NewGRPCServer,
	NewHTTPServer,
	NewPatientRepository,
	NewPatientCache,
	NewRegisterUsecase,
	NewProfileUsecase,
	NewPatientService,
	NewApp,
)

func NewPatientRepository(db *gorm.DB, logger log.Logger) repository.PatientRepository {
	return mysql.NewPatientRepository(db, logger)
}

func NewPatientCache(client *redis.Client, logger log.Logger) *redis.PatientCache {
	return redis.NewPatientCache(client, logger)
}

func NewRegisterUsecase(repo repository.PatientRepository, cfg *Config, logger log.Logger) *usecase.RegisterUsecase {
	return usecase.NewRegisterUsecase(repo, cfg.JWT.Key, logger)
}

func NewProfileUsecase(repo repository.PatientRepository, logger log.Logger) *usecase.ProfileUsecase {
	return usecase.NewProfileUsecase(repo, logger)
}

func NewPatientService(registerUC *usecase.RegisterUsecase, profileUC *usecase.ProfileUsecase, logger log.Logger) *service.PatientService {
	return service.NewPatientService(registerUC, profileUC, logger)
}

func NewApp(gs *grpc.Server, hs *http.Server, logger log.Logger) *kratos.App {
	return kratos.New(
		kratos.ID("patient-service"),
		kratos.Name("patient-service"),
		kratos.Version("v1.0.0"),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func InitApp(db *gorm.DB, redisClient *redis.Client, logger log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(ProviderSet))
}

type App struct {
	kratos *kratos.App
}

func (a *App) Run() error {
	return a.kratos.Run()
}

func (a *App) Stop() error {
	return a.kratos.Stop()
}

func (a *App) Servers() []transport.Server {
	return a.kratos.Servers()
}