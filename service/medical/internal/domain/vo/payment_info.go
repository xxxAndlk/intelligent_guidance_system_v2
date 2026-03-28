package vo

import (
	"errors"
	"time"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusUnpaid    PaymentStatus = "unpaid"
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusPaid      PaymentStatus = "paid"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

// PaymentMethod represents the method of payment
type PaymentMethod string

const (
	PaymentMethodCash      PaymentMethod = "cash"
	PaymentMethodCard      PaymentMethod = "card"
	PaymentMethodWechat    PaymentMethod = "wechat"
	PaymentMethodAlipay    PaymentMethod = "alipay"
	PaymentMethodInsurance PaymentMethod = "insurance"
)

var (
	ErrInvalidPaymentStatus = errors.New("invalid payment status")
	ErrInvalidPaymentMethod = errors.New("invalid payment method")
	ErrPaymentAlreadyPaid   = errors.New("payment already completed")
	ErrPaymentNotPaid       = errors.New("payment not completed")
)

// PaymentInfo value object - represents payment information
type PaymentInfo struct {
	amount   *Money
	method   PaymentMethod
	status   PaymentStatus
	payTime  *time.Time
}

// NewPaymentInfo creates a new PaymentInfo value object
func NewPaymentInfo(amount *Money, method PaymentMethod) (*PaymentInfo, error) {
	if amount == nil {
		amount = Zero("CNY")
	}
	if !isValidPaymentMethod(method) {
		return nil, ErrInvalidPaymentMethod
	}
	return &PaymentInfo{
		amount: amount,
		method: method,
		status: PaymentStatusUnpaid,
		payTime: nil,
	}, nil
}

// Amount returns the payment amount
func (p *PaymentInfo) Amount() *Money {
	return p.amount
}

// Method returns the payment method
func (p *PaymentInfo) Method() PaymentMethod {
	return p.method
}

// Status returns the payment status
func (p *PaymentInfo) Status() PaymentStatus {
	return p.status
}

// PayTime returns the payment time
func (p *PaymentInfo) PayTime() *time.Time {
	return p.payTime
}

// MarkAsPaid marks the payment as paid
func (p *PaymentInfo) MarkAsPaid(payTime time.Time) error {
	if p.status == PaymentStatusPaid {
		return ErrPaymentAlreadyPaid
	}
	p.status = PaymentStatusPaid
	p.payTime = &payTime
	return nil
}

// MarkAsFailed marks the payment as failed
func (p *PaymentInfo) MarkAsFailed() error {
	if p.status == PaymentStatusPaid {
		return ErrPaymentAlreadyPaid
	}
	p.status = PaymentStatusFailed
	return nil
}

// MarkAsRefunded marks the payment as refunded
func (p *PaymentInfo) MarkAsRefunded() error {
	if p.status != PaymentStatusPaid {
		return ErrPaymentNotPaid
	}
	p.status = PaymentStatusRefunded
	return nil
}

// IsPaid checks if the payment is paid
func (p *PaymentInfo) IsPaid() bool {
	return p.status == PaymentStatusPaid
}

// IsPending checks if the payment is pending
func (p *PaymentInfo) IsPending() bool {
	return p.status == PaymentStatusPending
}

// IsUnpaid checks if the payment is unpaid
func (p *PaymentInfo) IsUnpaid() bool {
	return p.status == PaymentStatusUnpaid
}

// SetAmount updates the payment amount
func (p *PaymentInfo) SetAmount(amount *Money) {
	p.amount = amount
}

// SetMethod updates the payment method
func (p *PaymentInfo) SetMethod(method PaymentMethod) error {
	if !isValidPaymentMethod(method) {
		return ErrInvalidPaymentMethod
	}
	p.method = method
	return nil
}

// isValidPaymentMethod validates payment method
func isValidPaymentMethod(method PaymentMethod) bool {
	switch method {
	case PaymentMethodCash, PaymentMethodCard, PaymentMethodWechat,
		PaymentMethodAlipay, PaymentMethodInsurance, "":
		return true
	default:
		return false
	}
}

// isValidPaymentStatus validates payment status
func isValidPaymentStatus(status PaymentStatus) bool {
	switch status {
	case PaymentStatusUnpaid, PaymentStatusPending, PaymentStatusPaid,
		PaymentStatusFailed, PaymentStatusRefunded:
		return true
	default:
		return false
	}
}