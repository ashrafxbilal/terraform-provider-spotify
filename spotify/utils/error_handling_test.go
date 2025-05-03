package utils

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestHandleAPIError(t *testing.T) {  
	// Test with a simple error
	ctx := context.Background()
	simpleErr := errors.New("simple error")
	diags := HandleAPIError(ctx, simpleErr, "test operation", "playlist", "123")
	
	if len(diags) == 0 {
		t.Error("Expected non-empty diagnostics, got empty")
	}
	
	// Check that the error message contains "API Error"
	if !strings.Contains(diags[0].Summary, "API Error") {
		t.Errorf("Expected diagnostics to contain 'API Error', got '%s'", diags[0].Summary)
	}
	
	// Check that the detail contains our operation
	if !strings.Contains(diags[0].Detail, "test operation") {
		t.Errorf("Expected diagnostics detail to contain 'test operation', got '%s'", diags[0].Detail)
	}
}

func TestHandleValidationError(t *testing.T) {
	// Test with a validation error message
	ctx := context.Background()
	fields := map[string]string{"field": "value"}
	diags := HandleValidationError(ctx, "test validation error", fields)
	
	if len(diags) == 0 {
		t.Error("Expected non-empty diagnostics, got empty")
	}
	
	// Check that the error message contains "Validation Error"
	if !strings.Contains(diags[0].Summary, "Validation Error") {
		t.Errorf("Expected diagnostics to contain 'Validation Error', got '%s'", diags[0].Summary)
	}
	
	// Check that the detail contains our error message
	if !strings.Contains(diags[0].Detail, "test validation error") {
		t.Errorf("Expected diagnostics detail to contain 'test validation error', got '%s'", diags[0].Detail)
	}
}

func TestIsSpotifyNotFoundError(t *testing.T) {
	// Test with a not found error
	notFoundErr := errors.New("not found")
	if !IsSpotifyNotFoundError(notFoundErr) {
		t.Errorf("Expected 'not found' to be recognized as a not found error")
	}
	
	// Test with a Not Found error (different case)
	notFoundErr2 := errors.New("Not Found")
	if !IsSpotifyNotFoundError(notFoundErr2) {
		t.Errorf("Expected 'Not Found' to be recognized as a not found error")
	}
	
	// Test with a non-not found error
	otherErr := errors.New("other error")
	if IsSpotifyNotFoundError(otherErr) {
		t.Errorf("Expected 'other error' to not be recognized as a not found error")
	}
}

// Helper function to check if a string contains another string
func containsString(s, substr string) bool {
	return strings.Contains(s, substr)
}