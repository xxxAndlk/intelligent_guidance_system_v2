package main

import (
	"flag"
	"os"
	"time"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/hashicorp/consul/api"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"intelligent-guidance-system/service/patient/internal/data/mysql"
	"intelligent-guidance-system/service/patient/internal/biz/usecase"
	"intelligent-guidance-system/service/patient/internal/service"
)

var flagConfig = flag.String("c", "./configs/config.yaml", "config path")

type Bootstrap struct {
	Server struct {
		ID       string
		Name     string
		Version  string
		Metadata map[string]string
		GRPC     struct {
			Addr    string
			Timeout time.Duration
		}
		HTTP struct {
			Addr    string
			Timeout time.Duration
		}
	}
	Data struct {
		Database struct {
			Source          string
			MaxIdleConns    int
			MaxOpenConns    int
			ConnMaxLifetime time.Duration
		}
		Redis struct {
			Addr         string
			Password     string
			DB           int
			PoolSize     int
			MinIdleConns int
		}
	}
	JWT struct {
		Key    string
		Expire time.Duration
	}
	Registry struct {
		Consul struct {
			Address string
			Scheme  string
		}
	}
}

func main() {
	flag.Parse()

	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", "patient-service",
		"service.version", "v1.0.0",
		"service.name", "patient-service",
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)

	c := config.New(config.WithSource(file.NewSource(*flagConfig)))
	if err := c.Load(); err != nil {
		log.Errorf("failed to load config: %v", err)
		return
	}

	var bc Bootstrap
	if err := c.Scan(&bc); err != nil {
		log.Errorf("failed to scan config: %v", err)
		return
	}

	db, err := gorm.Open(mysql.Open(bc.Data.Database.Source), &gorm.Config{})
	if err != nil {
		log.Errorf("failed to init database: %v", err)
		return
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(bc.Data.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(bc.Data.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(bc.Data.Database.ConnMaxLifetime)
	db.AutoMigrate(&mysql.PatientPO{}, &mysql.AllergyPO{}, &mysql.MedicalHistoryPO{})

	redisClient := redis.NewClient(&redis.Options{
		Addr:         bc.Data.Redis.Addr,
		Password:     bc.Data.Redis.Password,
		DB:           bc.Data.Redis.DB,
		PoolSize:     bc.Data.Redis.PoolSize,
		MinIdleConns: bc.Data.Redis.MinIdleConns,
	})

	var reg registry.Registrar
	if bc.Registry.Consul.Address != "" {
		consulClient, err := api.NewClient(&api.Config{
			Address: bc.Registry.Consul.Address,
			Scheme:  bc.Registry.Consul.Scheme,
		})
		if err == nil {
			reg = consul.New(consulClient)
		}
	}

	patientRepo := mysql.NewPatientRepository(db, logger)
	registerUC := usecase.NewRegisterUsecase(patientRepo, bc.JWT.Key, logger)
	profileUC := usecase.NewProfileUsecase(patientRepo, logger)
	patientSvc := service.NewPatientService(registerUC, profileUC, logger)

	grpcServer := grpc.NewServer(
		grpc.Address(bc.Server.GRPC.Addr),
		grpc.Timeout(bc.Server.GRPC.Timeout),
		grpc.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			validate.Validator(),
		),
	)

	httpServer := http.NewServer(
		http.Address(bc.Server.HTTP.Addr),
		http.Timeout(bc.Server.HTTP.Timeout),
		http.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			validate.Validator(),
		),
	)

	app := kratos.New(
		kratos.ID(bc.Server.ID),
		kratos.Name(bc.Server.Name),
		kratos.Version(bc.Server.Version),
		kratos.Metadata(bc.Server.Metadata),
		kratos.Logger(logger),
		kratos.Server(grpcServer, httpServer),
	)

	if reg != nil {
		app = kratos.New(
			kratos.ID(bc.Server.ID),
			kratos.Name(bc.Server.Name),
			kratos.Version(bc.Server.Version),
			kratos.Metadata(bc.Server.Metadata),
			kratos.Logger(logger),
			kratos.Registrar(reg),
			kratos.Server(grpcServer, httpServer),
		)
	}

	_ = redisClient

	if err := app.Run(); err != nil {
		log.Errorf("failed to run app: %v", err)
	}
}