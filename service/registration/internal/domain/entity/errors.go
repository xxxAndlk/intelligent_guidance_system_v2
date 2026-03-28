package entity

import "errors"

var (
	ErrInvalidRegistrationID   = errors.New("invalid registration ID")
	ErrInvalidPatientID        = errors.New("invalid patient ID")
	ErrInvalidDoctorID         = errors.New("invalid doctor ID")
	ErrInvalidDepartmentID     = errors.New("invalid department ID")
	ErrInvalidRegistrationType = errors.New("invalid registration type")
	ErrInvalidRegistrationStatus = errors.New("invalid registration status")
)