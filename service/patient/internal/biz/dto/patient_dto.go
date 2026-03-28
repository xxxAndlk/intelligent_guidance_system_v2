// service/patient/internal/biz/dto/patient_dto.go
package dto

import (
	"time"

	"intelligent-guidance-system/service/patient/internal/domain/entity"
)

// PatientDTO 患者数据传输对象
type PatientDTO struct {
	ID           int64
	OpenID       string
	Username     string
	Name         string
	Age          int32
	Gender       int32
	IDCard       string
	Phone        string
	Email        string
	Address      string
	Status       int32
	CreatedAt    time.Time
	UpdatedAt    time.Time
	MedicalCount int32
	AllergyCount int32
}

// PatientRegisterRequest 患者注册请求
type PatientRegisterRequest struct {
	Username string
	Password string
	Phone    string
	Name     string
	IDCard   string
}

// PatientLoginRequest 患者登录请求
type PatientLoginRequest struct {
	Phone    string
	Password string
}

// PatientLoginResponse 患者登录响应
type PatientLoginResponse struct {
	Token    string
	Patient  PatientDTO
}

// PatientUpdateRequest 患者更新请求
type PatientUpdateRequest struct {
	Name    string
	Email   string
	Address string
}

// PatientDetailDTO 患者详情DTO
type PatientDetailDTO struct {
	ID              int64
	OpenID          string
	Username        string
	Name            string
	Age             int32
	Gender          int32
	IDCardMasked    string
	PhoneMasked     string
	Email           string
	Address         string
	Status          int32
	StatusName      string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	MedicalHistory  []MedicalHistoryDTO
	Allergies       []AllergyDTO
}

// MedicalHistoryDTO 病史DTO
type MedicalHistoryDTO struct {
	ID           int64
	DiseaseName  string
	DiagnoseDate time.Time
	Treatment    string
	Status       int32
	StatusName   string
	CreatedAt    time.Time
}

// AllergyDTO 过敏史DTO
type AllergyDTO struct {
	ID        int64
	Allergen  string
	Severity  int32
	SeverityName string
	Reaction  string
	CreatedAt time.Time
}

// AddMedicalHistoryRequest 添加病史请求
type AddMedicalHistoryRequest struct {
	PatientID    int64
	DiseaseName  string
	DiagnoseDate time.Time
	Treatment    string
}

// AddAllergyRequest 添加过敏史请求
type AddAllergyRequest struct {
	PatientID int64
	Allergen  string
	Severity  int32
	Reaction  string
}

// PatientListRequest 患者列表请求
type PatientListRequest struct {
	Page     int32
	PageSize int32
}

// PatientListResponse 患者列表响应
type PatientListResponse struct {
	Total int64
	List  []PatientDTO
}

// ToMedicalHistoryDTO 转换病史实体为DTO
func ToMedicalHistoryDTO(history entity.MedicalHistory) MedicalHistoryDTO {
	return MedicalHistoryDTO{
		ID:           int64(history.ID),
		DiseaseName:  history.DiseaseName,
		DiagnoseDate: history.DiagnoseDate,
		Treatment:    history.Treatment,
		Status:       int32(history.Status),
		StatusName:   history.Status.String(),
		CreatedAt:    history.CreatedAt,
	}
}

// ToAllergyDTO 转换过敏史实体为DTO
func ToAllergyDTO(allergy entity.Allergy) AllergyDTO {
	return AllergyDTO{
		ID:           int64(allergy.ID),
		Allergen:     allergy.Allergen,
		Severity:     int32(allergy.Severity),
		SeverityName: allergy.Severity.String(),
		Reaction:     allergy.Reaction,
		CreatedAt:    allergy.CreatedAt,
	}
}