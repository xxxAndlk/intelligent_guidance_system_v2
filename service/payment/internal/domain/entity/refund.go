package entity

import (
	"errors"
	"time"
)

var (
	ErrInvalidRefundID    = errors.New("invalid refund ID")
	ErrRefundNotFound     = errors.New("refund not found")
	ErrRefundAlreadyProcessed = errors.New("refund already processed")
)

type RefundStatus int

const (
	RefundStatusUnknown RefundStatus = iota
	RefundStatusPending
	RefundStatusProcessing
	RefundStatusSuccess
	RefundStatusFailed
)

func (s RefundStatus) String() string {
	switch s {
	case RefundStatusPending:
		return "待退款"
	case RefundStatusProcessing:
		return "处理中"
	case RefundStatusSuccess:
		return "退款成功"
	case RefundStatusFailed:
		return "退款失败"
	default:
		return "未知"
	}
}

func RefundStatusFromCode(code string) RefundStatus {
	switch code {
	case "PENDING":
		return RefundStatusPending
	case "PROCESSING":
		return RefundStatusProcessing
	case "SUCCESS":
		return RefundStatusSuccess
	case "FAILED":
		return RefundStatusFailed
	default:
		return RefundStatusUnknown
	}
}

func (s RefundStatus) Code() string {
	switch s {
	case RefundStatusPending:
		return "PENDING"
	case RefundStatusProcessing:
		return "PROCESSING"
	case RefundStatusSuccess:
		return "SUCCESS"
	case RefundStatusFailed:
		return "FAILED"
	default:
		return "UNKNOWN"
	}
}

type Refund struct {
	id          int64
	paymentID   int64
	amount      int64
	reason      string
	status      RefundStatus
	processedAt time.Time
	createdAt   time.Time
	updatedAt   time.Time
}

func NewRefund(paymentID int64, amount int64, reason string) (*Refund, error) {
	if paymentID <= 0 {
		return nil, ErrInvalidPaymentID
	}
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	now := time.Now()
	return &Refund{
		id:        0,
		paymentID: paymentID,
		amount:    amount,
		reason:    reason,
		status:    RefundStatusPending,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func ReconstructRefund(id, paymentID int64, amount int64, reason string, status RefundStatus, processedAt, createdAt, updatedAt time.Time) *Refund {
	return &Refund{
		id:          id,
		paymentID:   paymentID,
		amount:      amount,
		reason:      reason,
		status:      status,
		processedAt: processedAt,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

func (r *Refund) ID() int64        { return r.id }
func (r *Refund) PaymentID() int64 { return r.paymentID }
func (r *Refund) Amount() int64    { return r.amount }
func (r *Refund) Reason() string   { return r.reason }
func (r *Refund) Status() RefundStatus { return r.status }
func (r *Refund) ProcessedAt() time.Time { return r.processedAt }
func (r *Refund) CreatedAt() time.Time { return r.createdAt }
func (r *Refund) UpdatedAt() time.Time { return r.updatedAt }

func (r *Refund) SetID(id int64) { r.id = id }

func (r *Refund) Process() error {
	if r.status != RefundStatusPending {
		return ErrRefundAlreadyProcessed
	}
	r.status = RefundStatusProcessing
	r.updatedAt = time.Now()
	return nil
}

func (r *Refund) Complete() error {
	r.status = RefundStatusSuccess
	r.processedAt = time.Now()
	r.updatedAt = time.Now()
	return nil
}

func (r *Refund) Fail() error {
	r.status = RefundStatusFailed
	r.updatedAt = time.Now()
	return nil
}