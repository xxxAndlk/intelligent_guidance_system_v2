package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"intelligent-guidance-system/service/doctor/internal/biz/usecase"
	"intelligent-guidance-system/service/doctor/internal/data/mysql"
	"intelligent-guidance-system/service/doctor/internal/server"
	"intelligent-guidance-system/service/doctor/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	flagconf string
	id, _    = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "./configs/config.yaml", "config path, eg: -conf config.yaml")
}

func main() {
	flag.Parse()
	logger := log.NewStdLogger(os.Stdout)
	log.SetLogger(logger)

	db, err := initDB()
	if err != nil {
		log.Errorf("failed to init database: %v", err)
		os.Exit(1)
	}

	doctorRepo := mysql.NewDoctorRepository(db)
	scheduleRepo := mysql.NewScheduleRepository(db)
	roleRepo := mysql.NewRoleRepository(db)
	assignmentRepo := mysql.NewDoctorRoleAssignmentRepository(db)
	eventRepo := mysql.NewDoctorEventRepository(db)

	doctorUC := usecase.NewDoctorUseCase(doctorRepo, roleRepo, assignmentRepo, eventRepo, logger)
	scheduleUC := usecase.NewScheduleUseCase(doctorRepo, scheduleRepo, eventRepo, logger)
	doctorSvc := service.NewDoctorService(doctorUC, scheduleUC, logger)

	httpSrv := server.NewHTTPServer(doctorSvc)
	grpcSrv := server.NewGRPCServer(doctorSvc)

	app := kratos.New(
		kratos.ID(id),
		kratos.Name("doctor-service"),
		kratos.Version("v1.0.0"),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			httpSrv,
			grpcSrv,
		),
	)

	if err := app.Run(); err != nil {
		log.Errorf("app run failed: %v", err)
		os.Exit(1)
	}
}

func initDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		getEnvOrDefault("DB_USER", "root"),
		getEnvOrDefault("DB_PASSWORD", "password"),
		getEnvOrDefault("DB_HOST", "localhost"),
		getEnvOrDefaultInt("DB_PORT", 3306),
		getEnvOrDefault("DB_NAME", "intelligent_guidance"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := autoMigrate(db); err != nil {
		return nil, err
	}

	return db, nil
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&mysql.DoctorPO{},
		&mysql.SchedulePO{},
		&mysql.TimeSlotPO{},
		&mysql.RolePO{},
		&mysql.RolePermissionPO{},
		&mysql.DoctorRoleAssignmentPO{},
		&mysql.DoctorEventPO{},
	)
}

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvOrDefaultInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		var result int
		fmt.Sscanf(val, "%d", &result)
		return result
	}
	return defaultVal
}