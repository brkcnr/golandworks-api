package apierror

import (
	"net/http"
)

// Base errors.
var (
	ErrDuplicateTodo     = New(http.StatusConflict, "this todo already exists")
	ErrEmptySearchQuery  = New(http.StatusBadRequest, "search query cannot be empty")
	ErrDBConnection      = New(http.StatusServiceUnavailable, "failed to connect to the database")
	ErrDBPing            = New(http.StatusServiceUnavailable, "failed to ping database")
	ErrDBRead            = New(http.StatusInternalServerError, "failed to read from database")
	ErrMissingDBPassword = New(http.StatusBadRequest, "database password is required")
	ErrInvalidDBPort     = New(http.StatusBadRequest, "invalid database port number")
	ErrInvalidRequest    = New(http.StatusBadRequest, "invalid request")
	ErrInternalServer    = New(http.StatusInternalServerError, "internal server error")
	ErrNotFound          = New(http.StatusNotFound, "resource not found")
)

// APIError represents an API error with HTTP status code.
type APIError struct {
	ErrType error  `json:"-"`
	Inner   error  `json:"-"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.Inner != nil {
		return e.Message + ": " + e.Inner.Error()
	}

	return e.Message
}

// Unwrap returns the wrapped error.
func (e *APIError) Unwrap() error {
	return e.Inner
}

// New creates a new APIError.
func New(code int, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
	}
}

// Wrap wraps an error with additional context.
func Wrap(err error, code int, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
		Inner:   err,
	}
}
