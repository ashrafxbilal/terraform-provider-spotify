package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zmb3/spotify/v2"

	"github.com/ashrafxbilal/terraform-provider-spotify/spotify/errors"
	"github.com/ashrafxbilal/terraform-provider-spotify/spotify/logging"
)

// HandleAPIError standardizes API error handling across resources and data sources
func HandleAPIError(ctx context.Context, err error, operation string, resourceType string, resourceID string) diag.Diagnostics {
	logger := logging.DefaultLogger.WithContext(ctx)

	// Log the error with context
	logger.Error("API error occurred",
		"operation", operation,
		"resource_type", resourceType,
		"resource_id", resourceID,
		"error", err.Error(),
	)

	// Create a standardized error with context
	spotifyErr := errors.NewAPIError(
		fmt.Sprintf("Failed to %s %s", operation, resourceType),
		err,
		map[string]string{
			"resource_id": resourceID,
			"operation":   operation,
		},
	)

	// Convert to diagnostics and return
	return spotifyErr.ToDiagnostics()
}

// HandleNotFoundError standardizes not found error handling
func HandleNotFoundError(ctx context.Context, resourceType string, resourceID string) diag.Diagnostics {
	logger := logging.DefaultLogger.WithContext(ctx)

	// Log the error with context
	logger.Warn("Resource not found",
		"resource_type", resourceType,
		"resource_id", resourceID,
	)

	// Create a standardized error
	spotifyErr := errors.NewNotFoundError(resourceType, resourceID)

	// Convert to diagnostics and return
	return spotifyErr.ToDiagnostics()
}

// HandleAuthError standardizes authentication error handling
func HandleAuthError(ctx context.Context, err error, operation string) diag.Diagnostics {
	logger := logging.DefaultLogger.WithContext(ctx)

	// Log the error with context
	logger.Error("Authentication error",
		"operation", operation,
		"error", err.Error(),
	)

	// Create a standardized error
	spotifyErr := errors.NewAuthError(
		"Failed to authenticate with Spotify API",
		err,
		map[string]string{
			"operation": operation,
		},
	)

	// Convert to diagnostics and return
	return spotifyErr.ToDiagnostics()
}

// HandleValidationError standardizes validation error handling
func HandleValidationError(ctx context.Context, message string, fields map[string]string) diag.Diagnostics {
	logger := logging.DefaultLogger.WithContext(ctx)

	// Create field string for logging
	fieldPairs := make([]string, 0, len(fields))
	for k, v := range fields {
		fieldPairs = append(fieldPairs, fmt.Sprintf("%s: %s", k, v))
	}

	// Log the validation error
	logger.Warn("Validation error",
		"message", message,
		"fields", strings.Join(fieldPairs, ", "),
	)

	// Create a standardized error
	spotifyErr := errors.NewValidationError(message, fields)

	// Convert to diagnostics and return
	return spotifyErr.ToDiagnostics()
}

// IsSpotifyNotFoundError checks if an error from the Spotify API is a not found error
func IsSpotifyNotFoundError(err error) bool {
	// This is a simplified example - in a real implementation,
	// you would need to check the specific error types or status codes
	// returned by the Spotify API client
	return err != nil && strings.Contains(strings.ToLower(err.Error()), "not found")
}

// Example of how to use these utilities in a resource read function
func ExampleResourceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger := logging.DefaultLogger.WithContext(ctx)
	_ = m.(spotify.Client) // Client is unused in this example
	resourceID := d.Id()

	// Log the operation
	logger.Info("Reading resource", "resource_id", resourceID)

	// Validate required fields
	if resourceID == "" {
		return HandleValidationError(ctx, "Resource ID is required", map[string]string{
			"id": "empty",
		})
	}

	// Example API call
	// playlist, err := client.GetPlaylist(ctx, spotify.ID(resourceID))
	// Simulating an error for demonstration
	var err error = fmt.Errorf("example error")

	if err != nil {
		// Check for specific error types
		if IsSpotifyNotFoundError(err) {
			// Resource not found, remove from state
			d.SetId("")
			return nil
		}

		// Handle other API errors
		return HandleAPIError(ctx, err, "read", "playlist", resourceID)
	}

	// Log successful operation
	logger.Info("Successfully read resource", "resource_id", resourceID)

	// Return empty diagnostics (success)
	return diag.Diagnostics{}
}