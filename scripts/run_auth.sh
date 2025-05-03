#!/bin/bash

# Get the directory where the script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Load environment variables
echo "Loading environment variables..."
set -a
source "$PROJECT_ROOT/.env"
set +a
echo "Environment variables loaded successfully."

# Change to the auth proxy directory
cd "$PROJECT_ROOT/spotify_auth_proxy"
echo "Running Spotify auth script..."

# Run the auth script
go run spotify-auth.go

# Keep the terminal open
echo ""
echo "Press Enter to close this window..."
read