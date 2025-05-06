# CI/CD and Testing Documentation

## Overview

This document outlines the CI/CD pipeline, testing strategy, and version management approach for the Terraform Spotify Provider.

## CI/CD Pipeline

The CI/CD pipeline is implemented using GitHub Actions and consists of the following workflows:

### Continuous Integration (CI)

Location: `.github/workflows/ci.yml`

Triggered on:
- Push to main branch
- Pull requests to main branch

Steps:
1. **Build and Test**
   - Set up Go environment
   - Verify dependencies
   - Build the provider
   - Run unit tests
   - Run linter

2. **Security Scan**
   - Run Gosec security scanner
   - Scan dependencies with Nancy

### Continuous Deployment (CD)

Location: `.github/workflows/release.yml`

Triggered on:
- Push of a tag matching the pattern `v*` (e.g., v0.1.0)

Steps:
1. **Release**
   - Set up Go environment
   - Import GPG key for signing
   - Run GoReleaser to build and publish the provider

### Scheduled Playlist Refresh

Location: `.github/workflows/scheduled-refresh.yml`

Triggered on:
- Daily schedule (8:00 AM UTC)
- Manual trigger

Steps:
1. **Refresh Playlists**
   - Set up Terraform
   - Configure credentials
   - Apply Terraform configuration to update playlists

## Testing Strategy

The provider implements a comprehensive testing strategy:

### Unit Tests

Run with: `make test`

These tests verify individual components and functions without making actual API calls to Spotify.

### Acceptance Tests

Run with: `make testacc`

These tests make actual API calls to Spotify and verify the provider's functionality end-to-end. They require valid Spotify API credentials set as environment variables.

Required environment variables for acceptance tests:
- `SPOTIFY_CLIENT_ID`
- `SPOTIFY_CLIENT_SECRET`
- `SPOTIFY_REDIRECT_URI`
- `SPOTIFY_REFRESH_TOKEN`

### Test Files

- `spotify/provider_test.go`: Tests for provider configuration
- `spotify/resource_spotify_playlist_test.go`: Tests for playlist resource

## Version Management

Version management is implemented using:

1. **Version Package**
   - Location: `version/version.go`
   - Defines version constants and functions
   - Version information is injected during build via ldflags

2. **GoReleaser**
   - Location: `.goreleaser.yml`
   - Automates the release process
   - Creates binaries for multiple platforms
   - Generates checksums and signs releases

3. **Makefile Targets**
   - `make release`: Tags and pushes a new release
   - Version is automatically extracted from `version/version.go`

## Scheduled Triggers for Playlist Refresh

The provider supports scheduled playlist refreshes through:

1. **GitHub Actions Workflow**
   - Daily automatic refresh
   - Manual trigger option

2. **Example Configuration**
   - Location: `examples/scheduled_refresh/`
   - Creates time-based and weather-based playlists
   - Updates playlist covers

3. **Makefile Target**
   - `make scheduled-refresh`: Manually runs the scheduled refresh

## Dependency Management

Dependencies are managed through:

1. **Pinned Versions**
   - All dependencies have specific versions in `go.mod`

2. **Update Script**
   - Location: `scripts/update_dependencies.sh`
   - Checks for outdated dependencies
   - Updates dependencies in a controlled manner

3. **Makefile Target**
   - `make deps`: Runs the dependency update script

## Getting Started with CI/CD

### Setting Up GitHub Secrets

To use the CI/CD pipeline, set up the following GitHub secrets:

1. For releases:
   - `GPG_PRIVATE_KEY`: Your GPG private key for signing releases
   - `PASSPHRASE`: Passphrase for your GPG key

2. For scheduled playlist refresh:
   - `SPOTIFY_CLIENT_ID`: Your Spotify API client ID
   - `SPOTIFY_CLIENT_SECRET`: Your Spotify API client secret
   - `SPOTIFY_REDIRECT_URI`: Your Spotify API redirect URI
   - `SPOTIFY_REFRESH_TOKEN`: Your Spotify API refresh token
   - `WEATHER_API_KEY`: Your OpenWeatherMap API key
   - `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`: For state storage (if using S3 backend)

### Creating a Release

To create a new release:

1. Update the version in `version/version.go`
2. Run `make release`

This will tag the current commit and push it to GitHub, triggering the release workflow.