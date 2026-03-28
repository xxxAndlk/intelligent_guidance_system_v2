package mysql

import "time"

type PaymentPO struct {
	ID            int64     `gorm:"primaryKey;autoIncrement"`
	MedicalID     int64     `gorm:"column:medical_id;index"`
	PatientID     int64     `gorm:"column:patient_id;index"`
	Amount        int64     `gorm:"column:amount"`
	Method        string    `gorm:"column:method;size:20"`
	Status        string    `gorm:"column:status;size:20"`
	TransactionID string    `gorm:"column:transaction_id;size:100"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (PaymentPO) TableName() string {
	return "payments"
}

type RefundPO struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	PaymentID   int64     `gorm:"column:payment_id;index"`
	Amount      int64     `gorm:"column:amount"`
	Reason      string    `gorm:"column:reason;size:200"`
	Status      string    `gorm:"column:status;size:20"`
	ProcessedAt time.Time `gorm:"column:processed_at"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (RefundPO) TableName() string {
	return "refunds"
}