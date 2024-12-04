package errors

import "fmt"

// AppError represents an application-specific error
type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// NewAppError creates a new AppError
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Common error codes
const (
	ErrCodeInvalidInput = 400
	ErrCodeUnauthorized = 401
	ErrCodeNotFound     = 404
	ErrCodeInternal     = 500
)
