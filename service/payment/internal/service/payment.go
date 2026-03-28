package service

import (
	"context"

	"intelligent-guidance-system/service/payment/internal/biz/dto"
	"intelligent-guidance-system/service/payment/internal/biz/usecase"
	"intelligent-guidance-system/service/payment/internal/domain/aggregate"
	"intelligent-guidance-system/service/payment/internal/domain/entity"
)

type PaymentService struct {
	usecase *usecase.PaymentUsecase
}

func NewPaymentService(usecase *usecase.PaymentUsecase) *PaymentService {
	return &PaymentService{usecase: usecase}
}

func (s *PaymentService) CreatePayment(ctx context.Context, req *dto.CreatePaymentRequest) (*dto.PaymentResponse, error) {
	method := entity.PaymentMethodFromCode(req.Method)

	payment, err := s.usecase.CreatePayment(ctx, req.MedicalID, req.PatientID, req.Amount, method)
	if err != nil {
		return nil, err
	}

	return s.toResponse(payment), nil
}

func (s *PaymentService) GetPayment(ctx context.Context, id int64) (*dto.PaymentResponse, error) {
	payment, err := s.usecase.GetPayment(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(payment), nil
}

func (s *PaymentService) ProcessPayment(ctx context.Context, id int64, req *dto.ProcessPaymentRequest) error {
	return s.usecase.ProcessPayment(ctx, id, req.TransactionID)
}

func (s *PaymentService) QueryStatus(ctx context.Context, id int64) (string, error) {
	status, err := s.usecase.QueryStatus(ctx, id)
	if err != nil {
		return "", err
	}
	return status.Code(), nil
}

func (s *PaymentService) Refund(ctx context.Context, paymentID int64, req *dto.RefundRequest) (*dto.RefundResponse, error) {
	refund, err := s.usecase.Refund(ctx, paymentID, req.Reason)
	if err != nil {
		return nil, err
	}

	return &dto.RefundResponse{
		ID:          refund.ID(),
		PaymentID:   refund.PaymentID(),
		Amount:      float64(refund.Amount()) / 100.0,
		Reason:      refund.Reason(),
		Status:      refund.Status().Code(),
		StatusName:  refund.Status().String(),
		ProcessedAt: refund.ProcessedAt(),
		CreatedAt:   refund.CreatedAt(),
	}, nil
}

func (s *PaymentService) GetPatientPayments(ctx context.Context, patientID int64) ([]*dto.PaymentResponse, error) {
	payments, err := s.usecase.GetPatientPayments(ctx, patientID)
	if err != nil {
		return nil, err
	}
	return s.toResponses(payments), nil
}

func (s *PaymentService) ListPayments(ctx context.Context, page, pageSize int) (*dto.PaymentListResponse, error) {
	payments, total, err := s.usecase.ListPayments(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	return &dto.PaymentListResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Payments: s.toResponseList(payments),
	}, nil
}

func (s *PaymentService) DeletePayment(ctx context.Context, id int64) error {
	return s.usecase.DeletePayment(ctx, id)
}

func (s *PaymentService) toResponse(payment *aggregate.Payment) *dto.PaymentResponse {
	return &dto.PaymentResponse{
		ID:            payment.ID(),
		MedicalID:     payment.MedicalID(),
		PatientID:     payment.PatientID(),
		Amount:        payment.AmountYuan(),
		Method:        payment.Method().Code(),
		MethodName:    payment.Method().String(),
		Status:        payment.Status().Code(),
		StatusName:    payment.Status().String(),
		TransactionID: payment.TransactionID(),
		CreatedAt:     payment.CreatedAt(),
		UpdatedAt:     payment.UpdatedAt(),
	}
}

func (s *PaymentService) toResponses(payments []*aggregate.Payment) []*dto.PaymentResponse {
	responses := make([]*dto.PaymentResponse, 0, len(payments))
	for _, payment := range payments {
		responses = append(responses, s.toResponse(payment))
	}
	return responses
}

func (s *PaymentService) toResponseList(payments []*aggregate.Payment) []dto.PaymentResponse {
	responses := make([]dto.PaymentResponse, 0, len(payments))
	for _, payment := range payments {
		responses = append(responses, *s.toResponse(payment))
	}
	return responses
}