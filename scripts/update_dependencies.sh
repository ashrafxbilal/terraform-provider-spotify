#!/bin/bash

# Script to check and update dependencies for the Terraform Spotify Provider
# This script helps maintain up-to-date dependencies while ensuring stability

# Get the directory where the script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Change to the project root directory
cd "$PROJECT_ROOT"

# Print header
echo "===================================================="
echo "Dependency Update Tool for Terraform Spotify Provider"
echo "===================================================="
echo 

# Check for outdated dependencies
echo "Checking for outdated dependencies..."
echo "---------------------------------------------------"
go list -u -m all | grep -v 'indirect' | grep -E '\[.*\]'
echo 

# Ask if user wants to update dependencies
read -p "Do you want to update direct dependencies? (y/n): " answer
if [[ "$answer" != "y" && "$answer" != "Y" ]]; then
    echo "Update canceled."
    exit 0
fi

# Update dependencies
echo "\nUpdating dependencies..."
echo "---------------------------------------------------"

# Update direct dependencies only
go get -u github.com/hashicorp/terraform-plugin-sdk/v2
go get -u github.com/zmb3/spotify/v2
go get -u golang.org/x/oauth2

# Tidy up the go.mod file
go mod tidy

echo "\nDependencies updated. Please review changes in go.mod and go.sum."
echo "Run tests to ensure everything still works correctly."
echo "---------------------------------------------------"
echo "Changes in go.mod:"
git diff go.mod

# Remind about security checks
echo "\nRemember to run security checks on the updated dependencies:"
echo "go list -json -m all | nancy sleuth"
echo "or use another security scanning tool."

echo "\nDone!"