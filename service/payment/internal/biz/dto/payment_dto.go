package dto

import "time"

type CreatePaymentRequest struct {
	MedicalID int64  `json:"medical_id" binding:"required"`
	PatientID int64  `json:"patient_id" binding:"required"`
	Amount    int64  `json:"amount" binding:"required"`
	Method    string `json:"method" binding:"required"`
}

type ProcessPaymentRequest struct {
	TransactionID string `json:"transaction_id" binding:"required"`
}

type RefundRequest struct {
	Reason string `json:"reason"`
}

type PaymentResponse struct {
	ID            int64     `json:"id"`
	MedicalID     int64     `json:"medical_id"`
	PatientID     int64     `json:"patient_id"`
	Amount        float64   `json:"amount"`
	Method        string    `json:"method"`
	MethodName    string    `json:"method_name"`
	Status        string    `json:"status"`
	StatusName    string    `json:"status_name"`
	TransactionID string    `json:"transaction_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type RefundResponse struct {
	ID          int64     `json:"id"`
	PaymentID   int64     `json:"payment_id"`
	Amount      float64   `json:"amount"`
	Reason      string    `json:"reason"`
	Status      string    `json:"status"`
	StatusName  string    `json:"status_name"`
	ProcessedAt time.Time `json:"processed_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type PaymentListResponse struct {
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
	Payments []PaymentResponse `json:"payments"`
}