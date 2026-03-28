package event

import "time"

type PaymentEvent struct {
	eventType string
	paymentID int64
	operatorID int64
	data      map[string]interface{}
	occurredAt time.Time
}

func NewPaymentEvent(eventType string, paymentID, operatorID int64, data map[string]interface{}) *PaymentEvent {
	return &PaymentEvent{
		eventType:  eventType,
		paymentID:  paymentID,
		operatorID: operatorID,
		data:       data,
		occurredAt: time.Now(),
	}
}

func (e *PaymentEvent) EventType() string { return e.eventType }
func (e *PaymentEvent) PaymentID() int64 { return e.paymentID }
func (e *PaymentEvent) OperatorID() int64 { return e.operatorID }
func (e *PaymentEvent) Data() map[string]interface{} { return e.data }
func (e *PaymentEvent) OccurredAt() time.Time { return e.occurredAt }

const (
	EventTypePaymentCreated  = "payment.created"
	EventTypePaymentSuccess  = "payment.success"
	EventTypePaymentFailed   = "payment.failed"
	EventTypeRefundCreated   = "refund.created"
	EventTypeRefundSuccess   = "refund.success"
)

func PaymentCreatedEvent(paymentID, medicalID int64, amount int64) *PaymentEvent {
	return NewPaymentEvent(EventTypePaymentCreated, paymentID, 0, map[string]interface{}{
		"medical_id": medicalID,
		"amount":     amount,
	})
}

func PaymentSuccessEvent(paymentID int64, transactionID string) *PaymentEvent {
	return NewPaymentEvent(EventTypePaymentSuccess, paymentID, 0, map[string]interface{}{
		"transaction_id": transactionID,
	})
}

func PaymentFailedEvent(paymentID int64, reason string) *PaymentEvent {
	return NewPaymentEvent(EventTypePaymentFailed, paymentID, 0, map[string]interface{}{
		"reason": reason,
	})
}

func RefundSuccessEvent(paymentID, refundID int64) *PaymentEvent {
	return NewPaymentEvent(EventTypeRefundSuccess, paymentID, 0, map[string]interface{}{
		"refund_id": refundID,
	})
}