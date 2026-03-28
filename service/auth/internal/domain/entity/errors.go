package entity

import "errors"

var (
	ErrInvalidUserID        = errors.New("invalid user ID")
	ErrEmptyUsername        = errors.New("username cannot be empty")
	ErrEmptyPassword        = errors.New("password cannot be empty")
	ErrInvalidUserType      = errors.New("invalid user type")
	ErrInvalidDataScope     = errors.New("invalid data scope")
	ErrRoleAlreadyAssigned  = errors.New("role already assigned")
	ErrRoleNotAssigned      = errors.New("role not assigned")
	ErrPermissionDenied     = errors.New("permission denied")
	ErrInvalidPermission    = errors.New("invalid permission")
	ErrInvalidRole          = errors.New("invalid role")
)