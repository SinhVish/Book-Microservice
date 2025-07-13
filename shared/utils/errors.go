package utils

import (
	"fmt"
	"net/http"
)

// CustomError represents a custom error with HTTP status code and message
type CustomError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error implements the error interface
func (e *CustomError) Error() string {
	return e.Message
}

// HTTP status code getter
func (e *CustomError) StatusCode() int {
	return e.Code
}

// Predefined error constructors - similar to Django's custom exceptions

// BadRequest creates a 400 Bad Request error
func BadRequest(message string) *CustomError {
	return &CustomError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

// Unauthorized creates a 401 Unauthorized error
func Unauthorized(message string) *CustomError {
	return &CustomError{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}

// Forbidden creates a 403 Forbidden error
func Forbidden(message string) *CustomError {
	return &CustomError{
		Code:    http.StatusForbidden,
		Message: message,
	}
}

// NotFound creates a 404 Not Found error
func NotFound(message string) *CustomError {
	return &CustomError{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

// Conflict creates a 409 Conflict error
func Conflict(message string) *CustomError {
	return &CustomError{
		Code:    http.StatusConflict,
		Message: message,
	}
}

// InternalServerError creates a 500 Internal Server Error
func InternalServerError(message string) *CustomError {
	return &CustomError{
		Code:    http.StatusInternalServerError,
		Message: message,
	}
}

// Generic custom error with specific code and message
func NewCustomError(code int, message string) *CustomError {
	return &CustomError{
		Code:    code,
		Message: message,
	}
}

// Helper function to create formatted error messages
func NewCustomErrorf(code int, format string, args ...interface{}) *CustomError {
	return &CustomError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}
