package version

import (
	"testing"
	"strings"
)

func TestGetVersion(t *testing.T) {
	// Store original values
	originalVersion := Version
	originalVersionPrerelease := VersionPrerelease
	
	// Restore original values after test
	defer func() {
		Version = originalVersion
		VersionPrerelease = originalVersionPrerelease
	}()

	// Test version without prerelease
	Version = "1.0.0"
	VersionPrerelease = ""
	result := GetVersion()
	if result != "1.0.0" {
		t.Errorf("GetVersion() = %v, want %v", result, "1.0.0")
	}

	// Test version with prerelease
	Version = "1.0.0"
	VersionPrerelease = "beta1"
	result = GetVersion()
	if result != "1.0.0-beta1" {
		t.Errorf("GetVersion() = %v, want %v", result, "1.0.0-beta1")
	}
}

func TestGetVersionInfo(t *testing.T) {
	// Store original values
	originalVersion := Version
	originalVersionPrerelease := VersionPrerelease
	originalGitCommit := GitCommit
	originalBuildDate := BuildDate
	
	// Restore original values after test
	defer func() {
		Version = originalVersion
		VersionPrerelease = originalVersionPrerelease
		GitCommit = originalGitCommit
		BuildDate = originalBuildDate
	}()

	// Set test values
	Version = "1.0.0"
	VersionPrerelease = "beta1"
	GitCommit = "abc123"
	BuildDate = "2023-01-01"

	// Call the function
	result := GetVersionInfo()

	// Check that the result contains all expected values
	expectedValues := []string{"1.0.0-beta1", "abc123", "2023-01-01"}
	for _, val := range expectedValues {
		if !strings.Contains(result, val) {
			t.Errorf("GetVersionInfo() = %v, should contain %v", result, val)
		}
	}
}