package errors

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// ErrorType defines the type of error for categorization and consistent handling
type ErrorType string

// Error types for consistent categorization
const (
	// ErrorTypeAPI represents errors from the Spotify API
	ErrorTypeAPI ErrorType = "API"

	// ErrorTypeAuth represents authentication errors
	ErrorTypeAuth ErrorType = "Authentication"

	// ErrorTypeValidation represents validation errors
	ErrorTypeValidation ErrorType = "Validation"

	// ErrorTypeNotFound represents resource not found errors
	ErrorTypeNotFound ErrorType = "NotFound"

	// ErrorTypeInternal represents internal provider errors
	ErrorTypeInternal ErrorType = "Internal"
)

// SpotifyError represents a standardized error structure for the Spotify provider
type SpotifyError struct {
	Type    ErrorType
	Message string
	Err     error
	Context map[string]string
}

// Error implements the error interface
func (e *SpotifyError) Error() string {
	base := fmt.Sprintf("%s Error: %s", e.Type, e.Message)
	if e.Err != nil {
		base = fmt.Sprintf("%s: %s", base, e.Err.Error())
	}

	// Add context information if available
	if len(e.Context) > 0 {
		contextStrings := make([]string, 0, len(e.Context))
		for k, v := range e.Context {
			contextStrings = append(contextStrings, fmt.Sprintf("%s: %s", k, v))
		}
		base = fmt.Sprintf("%s [%s]", base, strings.Join(contextStrings, ", "))
	}

	return base
}

// NewAPIError creates a new API error
func NewAPIError(message string, err error, context map[string]string) *SpotifyError {
	return &SpotifyError{
		Type:    ErrorTypeAPI,
		Message: message,
		Err:     err,
		Context: context,
	}
}

// NewAuthError creates a new authentication error
func NewAuthError(message string, err error, context map[string]string) *SpotifyError {
	return &SpotifyError{
		Type:    ErrorTypeAuth,
		Message: message,
		Err:     err,
		Context: context,
	}
}

// NewValidationError creates a new validation error
func NewValidationError(message string, context map[string]string) *SpotifyError {
	return &SpotifyError{
		Type:    ErrorTypeValidation,
		Message: message,
		Context: context,
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(resourceType string, identifier string) *SpotifyError {
	return &SpotifyError{
		Type:    ErrorTypeNotFound,
		Message: fmt.Sprintf("%s with ID '%s' not found", resourceType, identifier),
		Context: map[string]string{
			"resource_type": resourceType,
			"identifier":    identifier,
		},
	}
}

// NewInternalError creates a new internal error
func NewInternalError(message string, err error) *SpotifyError {
	return &SpotifyError{
		Type:    ErrorTypeInternal,
		Message: message,
		Err:     err,
	}
}

// ToDiagnostics converts a SpotifyError to Terraform diagnostics
func (e *SpotifyError) ToDiagnostics() diag.Diagnostics {
	var severity diag.Severity

	// Determine severity based on error type
	switch e.Type {
	case ErrorTypeValidation, ErrorTypeNotFound:
		severity = diag.Error
	case ErrorTypeAPI, ErrorTypeAuth, ErrorTypeInternal:
		severity = diag.Error
	default:
		severity = diag.Error
	}

	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: severity,
			Summary:  fmt.Sprintf("%s Error", e.Type),
			Detail:   e.Error(),
		},
	}
}

// ToDiag is a helper function to convert any error to diagnostics
func ToDiag(err error) diag.Diagnostics {
	if err == nil {
		return nil
	}

	// If it's already a SpotifyError, use its method
	if spotifyErr, ok := err.(*SpotifyError); ok {
		return spotifyErr.ToDiagnostics()
	}

	// Otherwise, create a generic internal error
	return NewInternalError("An unexpected error occurred", err).ToDiagnostics()
}

// IsNotFoundError checks if the error is a not found error
func IsNotFoundError(err error) bool {
	if spotifyErr, ok := err.(*SpotifyError); ok {
		return spotifyErr.Type == ErrorTypeNotFound
	}
	return false
}

// IsAuthError checks if the error is an authentication error
func IsAuthError(err error) bool {
	if spotifyErr, ok := err.(*SpotifyError); ok {
		return spotifyErr.Type == ErrorTypeAuth
	}
	return false
}

// WithContext adds context to an existing error
func WithContext(err error, key string, value string) error {
	if spotifyErr, ok := err.(*SpotifyError); ok {
		if spotifyErr.Context == nil {
			spotifyErr.Context = make(map[string]string)
		}
		spotifyErr.Context[key] = value
		return spotifyErr
	}

	// If it's not a SpotifyError, create a new internal error with context
	return &SpotifyError{
		Type:    ErrorTypeInternal,
		Message: "An unexpected error occurred",
		Err:     err,
		Context: map[string]string{key: value},
	}
}