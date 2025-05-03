# Terraform Provider for Spotify

This Terraform provider allows you to manage Spotify resources like playlists and tracks using infrastructure as code. It also provides data sources for retrieving information based on weather, time, and user profiles to create dynamic, mood-based playlists.

## Features

- Create and manage Spotify playlists
- Add and remove tracks from playlists
- Generate track recommendations based on:
  - Current weather conditions
  - Time of day
  - Custom moods and genres
- Retrieve user profile information

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
- [Go](https://golang.org/doc/install) >= 1.20 (for building the provider)
- A Spotify account and registered application

## Dependency Management

This project uses pinned dependency versions to ensure reproducible builds. For more information on our dependency management strategy, please see [DEPENDENCIES.md](DEPENDENCIES.md).

## Versioning

This project follows semantic versioning principles to ensure compatibility and clearly communicate changes. For detailed guidelines on our versioning strategy, please see [VERSIONING.md](VERSIONING.md).

## CI/CD Pipeline and Testing

This project implements a comprehensive CI/CD pipeline with GitHub Actions, including automated testing, security scanning, and scheduled playlist refreshes. For more information, please see [CI_CD.md](CI_CD.md).

### Running Acceptance Tests

To run the acceptance tests for this provider, you can use the included setup script:

```sh
./scripts/setup_test_env.sh
```

This script will:
1. Set up all required environment variables
2. Source variables from an existing `.env` file if present
3. Prompt for any missing credentials
4. Offer to run the auth proxy to obtain a refresh token if needed

After running the script, you can execute the tests with:

```sh
go test -v ./spotify -timeout 120m
```

## Docker Containers

This project provides Docker containers for both development and runtime use, making it easier to contribute to the project and use the provider without installing dependencies locally. The containers are automatically built and published to Docker Hub. For more information on using the Docker containers, please see [DOCKER.md](DOCKER.md).

## Building the Provider

1. Clone the repository
2. Build the provider using the Go `install` command:

```sh
go build -o terraform-provider-spotify
```

## Installing the Provider

To use the provider in your Terraform configuration, you'll need to set up a development override in your `~/.terraformrc` file:

```hcl
provider_installation {
  dev_overrides {
    "local/spotify" = "/path/to/your/terraform-provider-spotify"
  }
  direct {}
}
```

## Authentication

This provider requires a Spotify API client ID, client secret, and refresh token. You can obtain these by:

1. Creating a Spotify application in the [Spotify Developer Dashboard](https://developer.spotify.com/dashboard/applications)
2. Setting up a redirect URI (use `https://glitch.com/~spotify-oauth-redirect`)
3. Generating a refresh token using the included auth proxy (see below)

## Configuration

```hcl
terraform {
  required_providers {
    spotify = {
      source  = "local/spotify"
      version = "0.1.0"
    }
  }
}

provider "spotify" {
  client_id     = "your-client-id"     # Or use environment variable: TF_VAR_spotify_client_id
  client_secret = "your-client-secret" # Or use environment variable: TF_VAR_spotify_client_secret
  redirect_uri  = "your-redirect-uri"  # Or use environment variable: TF_VAR_spotify_redirect_uri
  refresh_token = "your-refresh-token" # Or use environment variable: TF_VAR_spotify_refresh_token
}
```

## Example Usage

### Creating a Weather-Based Playlist

```hcl
# Get weather data with suggested moods
data "spotify_weather" "current" {}

# Get tracks based on weather mood
data "spotify_tracks" "weather_based" {
  mood  = data.spotify_weather.current.mood
  limit = 10
}

# Create a weather-based playlist
resource "spotify_playlist" "weather_based" {
  name        = "${data.spotify_weather.current.mood} Weather Vibes"
  description = "Tracks based on the current weather: ${data.spotify_weather.current.temperature}°C"
  public      = true
  tracks      = data.spotify_tracks.weather_based.ids
}
```

### Creating a Time-Based Playlist

```hcl
# Get time-based data with suggested moods and genres
data "spotify_time" "now" {}

# Get tracks based on time of day
data "spotify_tracks" "time_based" {
  genre = data.spotify_time.now.genre
  mood  = data.spotify_time.now.mood
  limit = 10
}

# Create a time-based playlist
resource "spotify_playlist" "time_based" {
  name        = "${data.spotify_time.now.time_of_day} ${data.spotify_time.now.day_of_week} Mix"
  description = "Tracks for ${data.spotify_time.now.time_of_day} vibes on ${data.spotify_time.now.day_of_week}"
  public      = true
  tracks      = data.spotify_tracks.time_based.ids
}
```

## Data Sources

### spotify_user

Retrieves information about the authenticated Spotify user.

```hcl
data "spotify_user" "me" {}

output "user_id" {
  value = data.spotify_user.me.id
}
```

### spotify_weather

Retrieves current weather information and suggests moods based on temperature.

```hcl
data "spotify_weather" "current" {}

output "suggested_weather_moods" {
  value = data.spotify_weather.current.suggested_moods
}
```

### spotify_time

Provides time-based information and suggests moods and genres based on time of day.

```hcl
data "spotify_time" "now" {}

output "suggested_time_moods" {
  value = data.spotify_time.now.suggested_moods
}

output "suggested_time_genres" {
  value = data.spotify_time.now.suggested_genres
}
```

### spotify_tracks

Retrieves track recommendations based on genre, artist, and mood.

```hcl
data "spotify_tracks" "recommendations" {
  genre = "pop"
  mood  = "energetic"
  limit = 20
}

output "track_names" {
  value = data.spotify_tracks.recommendations.names
}
```

## Resources

### spotify_playlist

Manages a Spotify playlist.

```hcl
resource "spotify_playlist" "example" {
  name        = "My Terraform Playlist"
  description = "Created and managed by Terraform"
  public      = true
  tracks      = ["spotify:track:4iV5W9uYEdYUVa79Axb7Rh", "spotify:track:1301WleyT98MSxVHPZCA6M"]
}
```

### spotify_playlist_track

Manages a track within a playlist.

```hcl
resource "spotify_playlist_track" "example" {
  playlist_id = spotify_playlist.example.id
  track_id    = "4iV5W9uYEdYUVa79Axb7Rh"
  position    = 0
}
```

### spotify_playlist_cover

Manages a custom cover image for a playlist. You can provide an image URL or generate a cover with emojis based on mood or weather.

```hcl
# Example with image URL
resource "spotify_playlist_cover" "example_url" {
  playlist_id = spotify_playlist.example.id
  image_url   = "https://example.com/image.jpg"
}

# Example with emoji based on mood
resource "spotify_playlist_cover" "example_mood" {
  playlist_id = spotify_playlist.example.id
  mood        = "energetic"
  background_color = "#FF5733"
}

# Example with custom emoji
resource "spotify_playlist_cover" "example_emoji" {
  playlist_id = spotify_playlist.example.id
  emoji       = "🎸"
  background_color = "#1DB954"
}
```

## Getting Started

### Prerequisites

- [Go](https://golang.org/doc/install) 1.16 or higher
- [Terraform](https://www.terraform.io/downloads.html) 0.13 or higher
- A Spotify account with a registered application

### Setup Guide

1. **Clone the repository**

   ```sh
   git clone https://github.com/yourusername/terraform-provider-spotify.git
   cd terraform-provider-spotify
   ```

2. **Set up environment variables**

   Copy the example environment file and edit it with your credentials:

   ```sh
   cp .env.example .env
   # Edit .env with your Spotify credentials
   ```

3. **Get Spotify API tokens**

   Use the included auth proxy to get your refresh token:

   ```sh
   # Load environment variables
   source scripts/load_env.sh
   
   # Run the auth proxy
   cd spotify_auth_proxy
   go run spotify-auth.go
   ```

   Follow the prompts to authorize the application and get your tokens.
   Update your `.env` file with the obtained refresh token.

4. **Build the provider**

   ```sh
   # Return to the project root
   cd ..
   
   # Build the provider
   make build
   
   # Install the provider for local use
   make install
   ```

5. **Run an example**

   ```sh
   # Load environment variables again
   source scripts/load_env.sh
   
   # Run the example
   cd examples/basic_playlist
   terraform init
   terraform apply
   ```

## Environment Variables

Instead of hardcoding sensitive values in your Terraform configuration, use environment variables:

```sh
# These are set automatically when you source scripts/load_env.sh
export TF_VAR_spotify_client_id="your-client-id"
export TF_VAR_spotify_client_secret="your-client-secret"
export TF_VAR_spotify_redirect_uri="your-redirect-uri"
export TF_VAR_spotify_refresh_token="your-refresh-token"
```

Then in your Terraform configuration:

```hcl
provider "spotify" {
  client_id     = var.spotify_client_id
  client_secret = var.spotify_client_secret
  redirect_uri  = var.spotify_redirect_uri
  refresh_token = var.spotify_refresh_token
}

variable "spotify_client_id" {}
variable "spotify_client_secret" {}
variable "spotify_redirect_uri" {}
variable "spotify_refresh_token" {}
```

## Contributing

Contributions are welcome! Here's how you can contribute to this project:

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-new-feature`
3. Make your changes
4. Run tests: `go test ./...`
5. Commit your changes: `git commit -am 'Add some feature'`
6. Push to the branch: `git push origin feature/my-new-feature`
7. Submit a pull request

## License

MIT