package errors

import (
	"strings"
	"testing"
)

func TestNewAPIError(t *testing.T) {
	// Test with basic error message
	err := NewAPIError("test error", nil, nil)
	
	if !strings.Contains(err.Error(), "test error") {
		t.Errorf("Expected error message to contain 'test error', got '%s'", err.Error())
	}
	
	// Test with details
	details := map[string]string{"code": "404", "reason": "not found"}
	errWithDetails := NewAPIError("error with details", nil, details)
	
	if !strings.Contains(errWithDetails.Error(), "error with details") {
		t.Errorf("Expected error message to contain 'error with details', got '%s'", errWithDetails.Error())
	}
	
	// Check if context is included in error message
	if !strings.Contains(errWithDetails.Error(), "code: 404") {
		t.Errorf("Expected error to contain context 'code: 404', got '%s'", errWithDetails.Error())
	}
}

func TestNewValidationError(t *testing.T) {
	// Test with basic error message
	details := map[string]string{"field": "name"}
	err := NewValidationError("validation error", details)
	
	if !strings.Contains(err.Error(), "validation error") {
		t.Errorf("Expected error message to contain 'validation error', got '%s'", err.Error())
	}
	
	// Check if context is included in error message
	if !strings.Contains(err.Error(), "field: name") {
		t.Errorf("Expected error to contain context 'field: name', got '%s'", err.Error())
	}
}