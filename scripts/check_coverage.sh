#!/bin/bash

# Script to check test coverage for the Terraform Spotify Provider

THRESHOLD=70

echo "Running tests with coverage..."
go test -coverprofile=coverage.out ./...

echo "\nCoverage by function:"
go tool cover -func=coverage.out

# Extract total coverage percentage
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print substr($3, 1, length($3)-1)}')
echo "\nTotal coverage: $COVERAGE%"

# Generate HTML report
echo "Generating HTML coverage report..."
go tool cover -html=coverage.out -o coverage.html
echo "HTML coverage report generated: coverage.html"

# Check if coverage meets threshold
if (( $(echo "$COVERAGE < $THRESHOLD" | bc -l) )); then
  echo "\n❌ Test coverage is below the $THRESHOLD% threshold (Current: $COVERAGE%)"
  exit 1
else
  echo "\n✅ Test coverage meets or exceeds the $THRESHOLD% threshold (Current: $COVERAGE%)"
  exit 0
fi