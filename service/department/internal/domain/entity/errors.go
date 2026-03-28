package entity

import "errors"

var (
	ErrInvalidDepartmentID   = errors.New("invalid department ID")
	ErrEmptyDepartmentName   = errors.New("department name cannot be empty")
	ErrInvalidDepartmentType = errors.New("invalid department type")
	ErrDepartmentDisabled    = errors.New("department is disabled")
	ErrDoctorAlreadyAssigned = errors.New("doctor already assigned to department")
	ErrDoctorNotInDepartment = errors.New("doctor not in department")
	ErrInvalidDoctorID       = errors.New("invalid doctor ID")
)