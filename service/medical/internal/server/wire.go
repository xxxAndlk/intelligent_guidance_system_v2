package server

import (
	"github.com/google/wire"

	"gorm.io/gorm"
	"gorm.io/driver/mysql"

	"intelligent_guidance_system_v2/service/medical/internal/biz/usecase"
	"intelligent_guidance_system_v2/service/medical/internal/data/mysql"
	"intelligent_guidance_system_v2/service/medical/internal/acl"
	"intelligent_guidance_system_v2/service/medical/internal/service"
)

var ProviderSet = wire.NewSet(
	NewGRPCServer,
	NewHTTPServerDirect,
	service.NewMedicalService,
	usecase.NewMedicalUseCase,
	mysql.NewMedicalRecordRepository,
	mysql.NewEventRepository,
	NewDB,
	NewPatientACLClient,
	NewDoctorACLClient,
	NewDepartmentACLClient,
	NewPaymentACLClient,
)

func NewDB(cfg *Config) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(cfg.Database.URL), &gorm.Config{})
}

func NewPatientACLClient(cfg *Config) usecase.PatientACL {
	return acl.NewPatientACLClient(cfg.Services.PatientURL)
}

func NewDoctorACLClient(cfg *Config) usecase.DoctorACL {
	return acl.NewDoctorACLClient(cfg.Services.DoctorURL)
}

func NewDepartmentACLClient(cfg *Config) usecase.DepartmentACL {
	return acl.NewDepartmentACLClient(cfg.Services.DepartmentURL)
}

func NewPaymentACLClient(cfg *Config) usecase.PaymentACL {
	return acl.NewPaymentACLClient(cfg.Services.PaymentURL)
}

type Config struct {
	Server struct {
		HTTPPort int
		GRPCPort int
		Name     string
	}
	Database struct {
		URL string
	}
	Services struct {
		PatientURL     string
		DoctorURL      string
		DepartmentURL  string
		PaymentURL     string
	}
}