package usecase

import (
	"context"
	"errors"

	"intelligent-guidance-system/service/payment/internal/domain/aggregate"
	"intelligent-guidance-system/service/payment/internal/domain/entity"
	"intelligent-guidance-system/service/payment/internal/domain/repository"
)

var (
	ErrPaymentNotFound = errors.New("payment not found")
)

type PaymentUsecase struct {
	paymentRepo repository.PaymentRepository
	refundRepo  repository.RefundRepository
}

func NewPaymentUsecase(paymentRepo repository.PaymentRepository, refundRepo repository.RefundRepository) *PaymentUsecase {
	return &PaymentUsecase{paymentRepo: paymentRepo, refundRepo: refundRepo}
}

func (u *PaymentUsecase) CreatePayment(ctx context.Context,
	medicalID int64,
	patientID int64,
	amount int64,
	method entity.PaymentMethod,
) (*aggregate.Payment, error) {
	payment, err := aggregate.NewPayment(medicalID, patientID, amount, method)
	if err != nil {
		return nil, err
	}

	if err := u.paymentRepo.Save(ctx, payment); err != nil {
		return nil, err
	}

	return payment, nil
}

func (u *PaymentUsecase) GetPayment(ctx context.Context, id int64) (*aggregate.Payment, error) {
	payment, err := u.paymentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, ErrPaymentNotFound
	}
	return payment, nil
}

func (u *PaymentUsecase) ProcessPayment(ctx context.Context, id int64, transactionID string) error {
	payment, err := u.GetPayment(ctx, id)
	if err != nil {
		return err
	}

	if err := payment.ProcessPayment(transactionID); err != nil {
		return err
	}

	return u.paymentRepo.Save(ctx, payment)
}

func (u *PaymentUsecase) QueryStatus(ctx context.Context, id int64) (entity.PaymentStatus, error) {
	payment, err := u.GetPayment(ctx, id)
	if err != nil {
		return entity.PaymentStatusUnknown, err
	}
	return payment.Status(), nil
}

func (u *PaymentUsecase) Refund(ctx context.Context, paymentID int64, reason string) (*entity.Refund, error) {
	payment, err := u.GetPayment(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	refund, err := payment.Refund(reason)
	if err != nil {
		return nil, err
	}

	if err := u.paymentRepo.Save(ctx, payment); err != nil {
		return nil, err
	}

	if err := u.refundRepo.Save(ctx, refund); err != nil {
		return nil, err
	}

	return refund, nil
}

func (u *PaymentUsecase) GetPatientPayments(ctx context.Context, patientID int64) ([]*aggregate.Payment, error) {
	return u.paymentRepo.FindByPatientID(ctx, patientID)
}

func (u *PaymentUsecase) ListPayments(ctx context.Context, page, pageSize int) ([]*aggregate.Payment, int64, error) {
	return u.paymentRepo.FindAll(ctx, page, pageSize)
}

func (u *PaymentUsecase) DeletePayment(ctx context.Context, id int64) error {
	return u.paymentRepo.Delete(ctx, id)
}