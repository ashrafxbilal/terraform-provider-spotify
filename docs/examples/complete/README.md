# Complete Spotify Provider Example

This example demonstrates how to use all resources and data sources provided by the Spotify Terraform Provider to create dynamic, personalized playlists.

## Features

This example creates three different playlists:

1. **Time-based playlist** - Changes based on the time of day
2. **Weather-based playlist** - Changes based on current weather conditions
3. **Personalized playlist** - Based on the user's listening history

Each playlist is created with a custom cover image that reflects its theme.

## Prerequisites

- [Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
- A Spotify account with a registered application
- An OpenWeatherMap API key (optional, for weather-based features)

## Setup

1. Create a `terraform.tfvars` file with your credentials:

```hcl
spotify_client_id     = "your-client-id"
spotify_client_secret = "your-client-secret"
spotify_redirect_uri  = "your-redirect-uri"
spotify_refresh_token = "your-refresh-token"
weather_api_key       = "your-weather-api-key"  # Optional
```

2. Initialize Terraform:

```sh
terraform init
```

3. Apply the configuration:

```sh
terraform apply
```

## What This Example Creates

### Data Sources

- `spotify_user` - Gets information about the current user
- `spotify_time` - Gets time-based information and suggestions
- `spotify_weather` - Gets weather-based information and suggestions
- `spotify_user_preference` - Gets the user's listening preferences (short and long term)
- `spotify_featured_playlists` - Gets Spotify's featured playlists
- `spotify_new_releases` - Gets new album releases
- `spotify_tracks` - Gets recommended tracks based on various criteria

### Resources

- `spotify_playlist.time_playlist` - A playlist based on the time of day
- `spotify_playlist.weather_playlist` - A playlist based on current weather
- `spotify_playlist.personalized_playlist` - A playlist based on user preferences
- `spotify_playlist_cover.weather_cover` - A custom cover for the weather playlist
- `spotify_playlist_cover.time_cover` - A custom cover for the time playlist
- `spotify_playlist_cover.personalized_cover` - A custom cover for the personalized playlist

## Outputs

- `user_info` - Information about the current user
- `created_playlists` - Information about the created playlists
- `current_context` - Information about the current time and weather

## Notes

- The playlists will be created in your Spotify account and will be visible in your Spotify client
- The weather-based playlist requires an OpenWeatherMap API key
- The time and weather playlists include a `force_update` parameter to ensure they update even when Terraform doesn't detect changes