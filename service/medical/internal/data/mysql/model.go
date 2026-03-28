package mysql

import (
	"time"

	"gorm.io/gorm"
)

type MedicalRecordPO struct {
	ID            int64          `gorm:"primaryKey;autoIncrement"`
	MedicalNumber string         `gorm:"column:medical_number;type:varchar(64);uniqueIndex;not null"`
	PatientID     int64          `gorm:"column:patient_id;not null;index"`
	DoctorID      int64          `gorm:"column:doctor_id;not null;index"`
	DepartmentID  int64          `gorm:"column:department_id;not null;index"`
	Symptom       string         `gorm:"column:symptom;type:text"`
	Diagnosis     string         `gorm:"column:diagnosis;type:text"`
	TreatmentPlan string         `gorm:"column:treatment_plan;type:text"`
	Status        int            `gorm:"column:status;not null;default:1;index"`
	PaymentAmount int64          `gorm:"column:payment_amount;not null;default:0"`
	PaymentMethod string         `gorm:"column:payment_method;type:varchar(32);default:'unpaid'"`
	PaymentStatus string         `gorm:"column:payment_status;type:varchar(32);default:'unpaid'"`
	PaymentTime   *time.Time     `gorm:"column:payment_time"`
	CreatedAt     time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;index"`

	Diseases      []DiseaseItemPO      `gorm:"foreignKey:MedicalRecordID;references:ID"`
	Prescriptions []PrescriptionItemPO `gorm:"foreignKey:MedicalRecordID;references:ID"`
	StatusChanges []StatusChangePO     `gorm:"foreignKey:MedicalRecordID;references:ID"`
}

func (MedicalRecordPO) TableName() string {
	return "medical_records"
}

type DiseaseItemPO struct {
	ID              int64     `gorm:"primaryKey;autoIncrement"`
	MedicalRecordID int64     `gorm:"column:medical_record_id;not null;index"`
	ItemID          string    `gorm:"column:item_id;type:varchar(64);not null"`
	DiseaseID       int64     `gorm:"column:disease_id;not null"`
	DiseaseName     string    `gorm:"column:disease_name;type:varchar(128);not null"`
	Symptoms        string    `gorm:"column:symptoms;type:text"`
	Diagnosis       string    `gorm:"column:diagnosis;type:text"`
	CreatedAt       time.Time `gorm:"column:created_at;not null"`
}

func (DiseaseItemPO) TableName() string {
	return "medical_disease_items"
}

type PrescriptionItemPO struct {
	ID              int64     `gorm:"primaryKey;autoIncrement"`
	MedicalRecordID int64     `gorm:"column:medical_record_id;not null;index"`
	ItemID          string    `gorm:"column:item_id;type:varchar(64);not null"`
	DrugID          int64     `gorm:"column:drug_id;not null"`
	DrugName        string    `gorm:"column:drug_name;type:varchar(128);not null"`
	Quantity        int32     `gorm:"column:quantity;not null"`
	Price           int64     `gorm:"column:price;not null"`
	Usage           string    `gorm:"column:usage;type:text"`
	CreatedAt       time.Time `gorm:"column:created_at;not null"`
}

func (PrescriptionItemPO) TableName() string {
	return "medical_prescription_items"
}

type StatusChangePO struct {
	ID              int64     `gorm:"primaryKey;autoIncrement"`
	MedicalRecordID int64     `gorm:"column:medical_record_id;not null;index"`
	FromStatus      int       `gorm:"column:from_status;not null"`
	ToStatus        int       `gorm:"column:to_status;not null"`
	ChangeTime      time.Time `gorm:"column:change_time;not null"`
	OperatorID      int64     `gorm:"column:operator_id;not null"`
	Reason          string    `gorm:"column:reason;type:text"`
}

func (StatusChangePO) TableName() string {
	return "medical_status_changes"
}

type MedicalEventPO struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	RecordID    int64     `gorm:"column:record_id;not null;index"`
	EventType   string    `gorm:"column:event_type;type:varchar(64);not null"`
	OccurredAt  time.Time `gorm:"column:occurred_at;not null"`
	OperatorID  int64     `gorm:"column:operator_id;not null"`
	Payload     string    `gorm:"column:payload;type:text"`
}

func (MedicalEventPO) TableName() string {
	return "medical_events"
}