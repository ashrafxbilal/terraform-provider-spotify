#!/bin/bash
# Setup script for Terraform Spotify Provider

set -e

echo "Setting up Terraform Spotify Provider..."

# Check if .env file exists
if [ ! -f ".env" ]; then
  echo "Creating .env file from template..."
  cp .env.example .env
  echo "Please edit .env with your Spotify credentials before continuing."
  echo "Press Enter when you're ready to continue..."
  read
fi

# Load environment variables
echo "Loading environment variables..."
source scripts/load_env.sh

# Build and install the provider
echo "Building and installing the provider..."
make build
make install

echo ""
echo "Setup complete! You can now run the examples:"
echo "cd examples/basic_playlist"
echo "terraform init"
echo "terraform apply"
echo ""
echo "Or get your Spotify tokens if you haven't already:"
echo "cd spotify_auth_proxy"
echo "go run spotify-auth.go"