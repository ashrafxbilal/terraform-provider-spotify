#!/bin/bash

# Script to set up environment variables for Terraform Spotify Provider acceptance tests

echo "Setting up environment variables for Terraform Spotify Provider acceptance tests..."

# Check if .env file exists and source it if it does
if [ -f ".env" ]; then
  echo "Found .env file, sourcing variables..."
  source .env
fi

# Set TF_ACC to enable acceptance testing
export TF_ACC=1

# Check if required Spotify variables are set
if [ -z "$SPOTIFY_CLIENT_ID" ]; then
  read -p "Enter your Spotify Client ID: " SPOTIFY_CLIENT_ID
  export SPOTIFY_CLIENT_ID
fi

if [ -z "$SPOTIFY_CLIENT_SECRET" ]; then
  read -p "Enter your Spotify Client Secret: " SPOTIFY_CLIENT_SECRET
  export SPOTIFY_CLIENT_SECRET
fi

if [ -z "$SPOTIFY_REDIRECT_URI" ]; then
  read -p "Enter your Spotify Redirect URI (default: http://localhost:8080/callback): " SPOTIFY_REDIRECT_URI
  SPOTIFY_REDIRECT_URI=${SPOTIFY_REDIRECT_URI:-http://localhost:8080/callback}
  export SPOTIFY_REDIRECT_URI
fi

if [ -z "$SPOTIFY_REFRESH_TOKEN" ]; then
  echo "No refresh token found. You need to obtain a refresh token."
  echo "Would you like to run the auth proxy to get a refresh token? (y/n)"
  read run_auth
  
  if [[ "$run_auth" == "y" || "$run_auth" == "Y" ]]; then
    echo "Running auth proxy..."
    echo "After authorization, copy the refresh token and enter it below."
    bash ./scripts/run_auth.sh &
    auth_pid=$!
    
    # Wait for user to get the token
    read -p "Enter your Spotify Refresh Token: " SPOTIFY_REFRESH_TOKEN
    export SPOTIFY_REFRESH_TOKEN
    
    # Kill the auth proxy
    kill $auth_pid 2>/dev/null
  else
    read -p "Enter your Spotify Refresh Token: " SPOTIFY_REFRESH_TOKEN
    export SPOTIFY_REFRESH_TOKEN
  fi
fi

# Print the configured environment variables
echo ""
echo "Environment variables set:"
echo "TF_ACC=$TF_ACC"
echo "SPOTIFY_CLIENT_ID=$SPOTIFY_CLIENT_ID"
echo "SPOTIFY_CLIENT_SECRET=****" # Don't print the actual secret
echo "SPOTIFY_REDIRECT_URI=$SPOTIFY_REDIRECT_URI"
echo "SPOTIFY_REFRESH_TOKEN=****" # Don't print the actual token
if [ ! -z "$WEATHER_API_KEY" ]; then
  echo "WEATHER_API_KEY=****" # Don't print the actual key
fi

echo ""
echo "You can now run the acceptance tests with:"
echo "go test -v ./spotify -timeout 120m"
echo ""
echo "To run a specific test, use:"
echo "go test -v ./spotify -timeout 120m -run=TestAccSpotifyPlaylist"