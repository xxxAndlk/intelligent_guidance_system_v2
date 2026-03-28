package repository

import (
	"context"

	"intelligent-guidance-system/service/payment/internal/domain/aggregate"
	"intelligent-guidance-system/service/payment/internal/domain/entity"
)

type PaymentRepository interface {
	Save(ctx context.Context, payment *aggregate.Payment) error
	FindByID(ctx context.Context, id int64) (*aggregate.Payment, error)
	FindByMedicalID(ctx context.Context, medicalID int64) (*aggregate.Payment, error)
	FindByPatientID(ctx context.Context, patientID int64) ([]*aggregate.Payment, error)
	FindByStatus(ctx context.Context, status int) ([]*aggregate.Payment, error)
	FindByTransactionID(ctx context.Context, transactionID string) (*aggregate.Payment, error)
	FindAll(ctx context.Context, page, pageSize int) ([]*aggregate.Payment, int64, error)
	Delete(ctx context.Context, id int64) error
}

type RefundRepository interface {
	Save(ctx context.Context, refund *entity.Refund) error
	FindByID(ctx context.Context, id int64) (*entity.Refund, error)
	FindByPaymentID(ctx context.Context, paymentID int64) ([]*entity.Refund, error)
	Delete(ctx context.Context, id int64) error
}