---
page_title: "Spotify Provider"
subcategory: ""
description: |-
  The Spotify provider allows Terraform to manage Spotify resources like playlists and tracks.
---

# Spotify Provider

The Spotify provider allows you to manage Spotify resources like playlists and tracks using infrastructure as code. It also provides data sources for retrieving information based on weather, time, and user profiles to create dynamic, mood-based playlists.

## Example Usage

```terraform
terraform {
  required_providers {
    spotify = {
      source  = "ashrafxbilal/spotify"
      version = "0.1.0"
    }
  }
}

provider "spotify" {
  client_id     = var.spotify_client_id
  client_secret = var.spotify_client_secret
  redirect_uri  = var.spotify_redirect_uri
  refresh_token = var.spotify_refresh_token
}
```

## Authentication

The Spotify provider requires several credentials to authenticate with the Spotify API:

- **client_id** - Your Spotify application client ID
- **client_secret** - Your Spotify application client secret
- **redirect_uri** - The redirect URI configured for your Spotify application
- **refresh_token** - A refresh token obtained through the OAuth flow

These can be provided in the provider configuration block or as environment variables:

```sh
export SPOTIFY_CLIENT_ID="your-client-id"
export SPOTIFY_CLIENT_SECRET="your-client-secret"
export SPOTIFY_REDIRECT_URI="your-redirect-uri"
export SPOTIFY_REFRESH_TOKEN="your-refresh-token"
```

## Getting Started

To obtain the necessary credentials:

1. Create a Spotify application in the [Spotify Developer Dashboard](https://developer.spotify.com/dashboard/applications)
2. Note your Client ID and Client Secret
3. Add a Redirect URI (e.g., `http://localhost:8080/callback`)
4. Use the included auth proxy to obtain a refresh token:

```sh
cd spotify_auth_proxy
go run spotify-auth.go
```

## Schema

### Required

- **client_id** (String) - Your Spotify application client ID
- **client_secret** (String) - Your Spotify application client secret
- **redirect_uri** (String) - The redirect URI configured for your Spotify application
- **refresh_token** (String) - A refresh token obtained through the OAuth flow

### Optional

- **weather_api_key** (String) - OpenWeatherMap API key for weather-based playlists