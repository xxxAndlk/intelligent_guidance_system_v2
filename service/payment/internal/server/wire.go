package server

import (
	"github.com/google/wire"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"intelligent-guidance-system/service/payment/internal/biz/usecase"
	"intelligent-guidance-system/service/payment/internal/data/mysql"
	"intelligent-guidance-system/service/payment/internal/domain/repository"
	"intelligent-guidance-system/service/payment/internal/service"
)

var ProviderSet = wire.NewSet(
	NewConfig,
	NewLogger,
	NewGRPCServer,
	NewHTTPServer,
	NewApp,
	NewDB,
	NewPaymentRepository,
	NewRefundRepository,
	NewPaymentUsecase,
	NewPaymentService,
)

func NewDB(cfg *Config) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(cfg.Database.Url), &gorm.Config{})
}

func NewPaymentRepository(db *gorm.DB) repository.PaymentRepository {
	return mysql.NewPaymentRepoImpl(db)
}

func NewRefundRepository(db *gorm.DB) repository.RefundRepository {
	return mysql.NewRefundRepoImpl(db)
}

func NewPaymentUsecase(paymentRepo repository.PaymentRepository, refundRepo repository.RefundRepository) *usecase.PaymentUsecase {
	return usecase.NewPaymentUsecase(paymentRepo, refundRepo)
}

func NewPaymentService(uc *usecase.PaymentUsecase) *service.PaymentService {
	return service.NewPaymentService(uc)
}