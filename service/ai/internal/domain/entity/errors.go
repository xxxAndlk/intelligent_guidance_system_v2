// Package entity provides domain entities for the AI diagnosis service.
package entity

import (
	"errors"
)

// Error definitions for Symptom entity
var (
	ErrInvalidBodyPart    = errors.New("invalid body part: body part cannot be empty")
	ErrInvalidDescription = errors.New("invalid description: description cannot be empty")
	ErrInvalidSeverity    = errors.New("invalid severity: must be between 1 and 10")
)