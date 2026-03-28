package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// Common error definitions
var (
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrNotFound          = errors.New("resource not found")
	ErrAlreadyExists     = errors.New("resource already exists")
	ErrInvalidArgument   = errors.New("invalid argument")
	ErrInternal          = errors.New("internal error")
	ErrConflict          = errors.New("resource conflict")
	ErrTooManyRequests   = errors.New("too many requests")
	ErrServiceUnavailable = errors.New("service unavailable")
	ErrTimeout           = errors.New("request timeout")
	ErrCanceled          = errors.New("request canceled")
	ErrNotImplemented    = errors.New("not implemented")
	ErrBadGateway        = errors.New("bad gateway")
	ErrGatewayTimeout    = errors.New("gateway timeout")
)

// AppError represents an application error with additional context
type AppError struct {
	Code       int                    `json:"code"`
	Message    string                 `json:"message"`
	Detail     string                 `json:"detail,omitempty"`
	Internal   error                  `json:"-"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	HTTPStatus int                    `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Internal)
	}
	if e.Detail != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Detail)
	}
	return e.Message
}

// Unwrap implements the errors.Unwrap interface
func (e *AppError) Unwrap() error {
	return e.Internal
}

// NewAppError creates a new application error
func NewAppError(code int, message string, opts ...ErrorOption) *AppError {
	appErr := &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: http.StatusOK,
	}

	for _, opt := range opts {
		opt(appErr)
	}

	return appErr
}

// ErrorOption is a functional option for AppError
type ErrorOption func(*AppError)

// WithDetail sets the detail message
func WithDetail(detail string) ErrorOption {
	return func(e *AppError) {
		e.Detail = detail
	}
}

// WithInternal sets the internal error
func WithInternal(err error) ErrorOption {
	return func(e *AppError) {
		e.Internal = err
	}
}

// WithHTTPStatus sets the HTTP status code
func WithHTTPStatus(status int) ErrorOption {
	return func(e *AppError) {
		e.HTTPStatus = status
	}
}

// WithMetadata sets the metadata
func WithMetadata(metadata map[string]interface{}) ErrorOption {
	return func(e *AppError) {
		e.Metadata = metadata
	}
}

// Error codes
const (
	CodeSuccess           = 0
	CodeUnknown           = 1
	CodeInvalidArgument   = 2
	CodeNotFound          = 3
	CodeAlreadyExists     = 4
	CodePermissionDenied  = 5
	CodeUnauthenticated   = 6
	CodeResourceExhausted = 7
	CodeCancelled         = 8
	CodeInternal          = 9
	CodeUnavailable       = 10
	CodeTimeout           = 11
	CodeConflict          = 12
)

// Predefined error constructors
func NewUnauthorized(message string, opts ...ErrorOption) *AppError {
	return NewAppError(CodeUnauthenticated, message,
		append(opts, WithHTTPStatus(http.StatusUnauthorized))...)
}

func NewForbidden(message string, opts ...ErrorOption) *AppError {
	return NewAppError(CodePermissionDenied, message,
		append(opts, WithHTTPStatus(http.StatusForbidden))...)
}

func NewNotFound(resource string, opts ...ErrorOption) *AppError {
	return NewAppError(CodeNotFound, fmt.Sprintf("%s not found", resource),
		append(opts, WithHTTPStatus(http.StatusNotFound))...)
}

func NewAlreadyExists(resource string, opts ...ErrorOption) *AppError {
	return NewAppError(CodeAlreadyExists, fmt.Sprintf("%s already exists", resource),
		append(opts, WithHTTPStatus(http.StatusConflict))...)
}

func NewInvalidArgument(field, reason string, opts ...ErrorOption) *AppError {
	return NewAppError(CodeInvalidArgument, fmt.Sprintf("invalid %s: %s", field, reason),
		append(opts, WithHTTPStatus(http.StatusBadRequest))...)
}

func NewInternal(message string, opts ...ErrorOption) *AppError {
	return NewAppError(CodeInternal, message,
		append(opts, WithHTTPStatus(http.StatusInternalServerError))...)
}

func NewConflict(message string, opts ...ErrorOption) *AppError {
	return NewAppError(CodeConflict, message,
		append(opts, WithHTTPStatus(http.StatusConflict))...)
}

func NewTimeout(message string, opts ...ErrorOption) *AppError {
	return NewAppError(CodeTimeout, message,
		append(opts, WithHTTPStatus(http.StatusRequestTimeout))...)
}

func NewUnavailable(message string, opts ...ErrorOption) *AppError {
	return NewAppError(CodeUnavailable, message,
		append(opts, WithHTTPStatus(http.StatusServiceUnavailable))...)
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// GetAppError extracts AppError from error or creates a new one
func GetAppError(err error) *AppError {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	// Wrap standard errors
	switch {
	case errors.Is(err, ErrUnauthorized):
		return NewUnauthorized(err.Error())
	case errors.Is(err, ErrForbidden):
		return NewForbidden(err.Error())
	case errors.Is(err, ErrNotFound):
		return NewNotFound("resource")
	case errors.Is(err, ErrAlreadyExists):
		return NewAlreadyExists("resource")
	case errors.Is(err, ErrInvalidArgument):
		return NewInvalidArgument("", err.Error())
	case errors.Is(err, ErrInternal):
		return NewInternal(err.Error())
	case errors.Is(err, ErrTimeout):
		return NewTimeout(err.Error())
	default:
		return NewInternal(err.Error(), WithInternal(err))
	}
}

// FromError creates an AppError from a standard error
func FromError(err error) *AppError {
	return GetAppError(err)
}

// Wrap wraps an error with additional context
func Wrap(err error, message string) *AppError {
	if err == nil {
		return nil
	}

	appErr := GetAppError(err)
	appErr.Message = message
	appErr.Internal = err
	return appErr
}

// Wrapf wraps an error with formatted message
func Wrapf(err error, format string, args ...interface{}) *AppError {
	return Wrap(err, fmt.Sprintf(format, args...))
}