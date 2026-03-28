package aggregate

import (
	"errors"
	"time"

	"intelligent-guidance-system/service/payment/internal/domain/entity"
	"intelligent-guidance-system/service/payment/internal/domain/event"
)

var (
	ErrPaymentNotFound      = errors.New("payment not found")
	ErrPaymentAlreadyPaid   = errors.New("payment already paid")
	ErrPaymentAlreadyRefunded = errors.New("payment already refunded")
	ErrPaymentCannotRefund = errors.New("payment cannot be refunded")
	ErrInvalidAmount        = errors.New("invalid amount")
)

type Payment struct {
	id            int64
	medicalID     int64
	patientID     int64
	amount        int64
	method        entity.PaymentMethod
	status        entity.PaymentStatus
	transactionID string
	refunds       []*entity.Refund
	events        []*event.PaymentEvent
	createdAt     time.Time
	updatedAt     time.Time
}

func NewPayment(
	medicalID int64,
	patientID int64,
	amount int64,
	method entity.PaymentMethod,
) (*Payment, error) {
	if medicalID <= 0 {
		return nil, errors.New("invalid medical ID")
	}
	if patientID <= 0 {
		return nil, errors.New("invalid patient ID")
	}
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	now := time.Now()
	payment := &Payment{
		id:            0,
		medicalID:     medicalID,
		patientID:     patientID,
		amount:        amount,
		method:        method,
		status:        entity.PaymentStatusPending,
		transactionID: "",
		refunds:       make([]*entity.Refund, 0),
		events:        make([]*event.PaymentEvent, 0),
		createdAt:     now,
		updatedAt:     now,
	}

	payment.events = append(payment.events, event.PaymentCreatedEvent(0, medicalID, amount))
	return payment, nil
}

func ReconstructPayment(
	id int64,
	medicalID int64,
	patientID int64,
	amount int64,
	method entity.PaymentMethod,
	status entity.PaymentStatus,
	transactionID string,
	refunds []*entity.Refund,
	createdAt time.Time,
	updatedAt time.Time,
) *Payment {
	return &Payment{
		id:            id,
		medicalID:     medicalID,
		patientID:     patientID,
		amount:        amount,
		method:        method,
		status:        status,
		transactionID: transactionID,
		refunds:       refunds,
		events:        make([]*event.PaymentEvent, 0),
		createdAt:     createdAt,
		updatedAt:     updatedAt,
	}
}

func (p *Payment) ID() int64 { return p.id }
func (p *Payment) MedicalID() int64 { return p.medicalID }
func (p *Payment) PatientID() int64 { return p.patientID }
func (p *Payment) Amount() int64 { return p.amount }
func (p *Payment) Method() entity.PaymentMethod { return p.method }
func (p *Payment) Status() entity.PaymentStatus { return p.status }
func (p *Payment) TransactionID() string { return p.transactionID }
func (p *Payment) Refunds() []*entity.Refund { return p.refunds }
func (p *Payment) Events() []*event.PaymentEvent { return p.events }
func (p *Payment) CreatedAt() time.Time { return p.createdAt }
func (p *Payment) UpdatedAt() time.Time { return p.updatedAt }

func (p *Payment) SetID(id int64) { p.id = id }

func (p *Payment) ProcessPayment(transactionID string) error {
	if p.status.IsTerminal() {
		if p.status == entity.PaymentStatusSuccess {
			return ErrPaymentAlreadyPaid
		}
		return ErrPaymentNotFound
	}

	p.status = entity.PaymentStatusSuccess
	p.transactionID = transactionID
	p.updatedAt = time.Now()
	p.events = append(p.events, event.PaymentSuccessEvent(p.id, transactionID))
	return nil
}

func (p *Payment) FailPayment(reason string) error {
	if p.status.IsTerminal() {
		return ErrPaymentNotFound
	}

	p.status = entity.PaymentStatusFailed
	p.updatedAt = time.Now()
	p.events = append(p.events, event.PaymentFailedEvent(p.id, reason))
	return nil
}

func (p *Payment) Refund(reason string) (*entity.Refund, error) {
	if !p.status.CanRefund() {
		return nil, ErrPaymentCannotRefund
	}

	refund, err := entity.NewRefund(p.id, p.amount, reason)
	if err != nil {
		return nil, err
	}

	p.refunds = append(p.refunds, refund)
	p.status = entity.PaymentStatusRefunded
	p.updatedAt = time.Now()

	return refund, nil
}

func (p *Payment) ClearEvents() {
	p.events = make([]*event.PaymentEvent, 0)
}

func (p *Payment) AmountYuan() float64 {
	return float64(p.amount) / 100.0
}

func (p *Payment) IsPaid() bool {
	return p.status == entity.PaymentStatusSuccess
}

func (p *Payment) IsRefunded() bool {
	return p.status == entity.PaymentStatusRefunded
}