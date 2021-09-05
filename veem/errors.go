package veem

import (
	"fmt"
	"time"
)

// APIError represents a Veem API error.
type APIError struct {
	// The type of the error
	ErrorType string `json:"error"`
	// ErrorDescription will usually be blank, it is present
	// during authentication errors.
	ErrorDescription string `json:"error_description"`
	// The error code, if present.
	Code int `json:"code"`
	// An error message, if present.
	Message string `json:"message"`
	// The time the error was produced.
	Timestamp time.Time `json:"timestamp"`
}

// Error implements the error interface.
func (a *APIError) Error() string {
	if a.ErrorDescription != "" {
		return fmt.Sprintf("(%s) %s", a.ErrorType, a.ErrorDescription)
	}
	return fmt.Sprintf("(%s) %s", a.ErrorType, a.Message)
}
