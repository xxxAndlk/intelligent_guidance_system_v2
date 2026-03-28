// service/patient/internal/data/mysql/model.go
package mysql

import (
	"time"

	"gorm.io/gorm"
)

// PatientPO 患者持久化对象
type PatientPO struct {
	ID           int64          `gorm:"primaryKey;autoIncrement"`
	OpenID       string         `gorm:"column:open_id;type:varchar(64);uniqueIndex"`
	Username     string         `gorm:"column:username;type:varchar(50);uniqueIndex"`
	Password     string         `gorm:"column:password;type:varchar(255);not null"`
	Name         string         `gorm:"column:name;type:varchar(50);not null"`
	Age          int32          `gorm:"column:age;type:int"`
	Gender       int32          `gorm:"column:gender;type:int;default:0"`
	IDCard       string         `gorm:"column:id_card;type:varchar(18);uniqueIndex"`
	Phone        string         `gorm:"column:phone;type:varchar(11);uniqueIndex;not null"`
	Email        string         `gorm:"column:email;type:varchar(100)"`
	Address      string         `gorm:"column:address;type:varchar(200)"`
	EmergencyName   string     `gorm:"column:emergency_name;type:varchar(50)"`
	EmergencyPhone  string     `gorm:"column:emergency_phone;type:varchar(11)"`
	EmergencyRelation string   `gorm:"column:emergency_relation;type:varchar(20)"`
	Status       int32          `gorm:"column:status;type:int;default:1"`
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index"`

	Allergies    []AllergyPO    `gorm:"foreignKey:PatientID;constraint:OnDelete:CASCADE"`
	MedicalHistories []MedicalHistoryPO `gorm:"foreignKey:PatientID;constraint:OnDelete:CASCADE"`
}

// TableName 指定表名
func (PatientPO) TableName() string {
	return "patients"
}

// AllergyPO 过敏史持久化对象
type AllergyPO struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	PatientID int64     `gorm:"column:patient_id;type:bigint;not null;index"`
	Allergen  string    `gorm:"column:allergen;type:varchar(100);not null"`
	Severity  int32     `gorm:"column:severity;type:int;not null"`
	Reaction  string    `gorm:"column:reaction;type:varchar(200)"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

// TableName 指定表名
func (AllergyPO) TableName() string {
	return "patient_allergies"
}

// MedicalHistoryPO 病史持久化对象
type MedicalHistoryPO struct {
	ID           int64     `gorm:"primaryKey;autoIncrement"`
	PatientID    int64     `gorm:"column:patient_id;type:bigint;not null;index"`
	DiseaseName  string    `gorm:"column:disease_name;type:varchar(100);not null"`
	DiagnoseDate time.Time `gorm:"column:diagnose_date;type:date;not null"`
	Treatment    string    `gorm:"column:treatment;type:varchar(500)"`
	Status       int32     `gorm:"column:status;type:int;default:1"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

// TableName 指定表名
func (MedicalHistoryPO) TableName() string {
	return "patient_medical_histories"
}