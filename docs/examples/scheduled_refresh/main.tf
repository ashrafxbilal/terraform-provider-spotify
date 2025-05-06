terraform {
  required_providers {
    spotify = {
      source  = "local/spotify"
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

# Get current user information
data "spotify_user" "me" {}

# Get current time data with suggested moods and genres
data "spotify_time" "now" {}

# Get current weather data
data "spotify_weather" "current" {}

# Get tracks based on time of day
data "spotify_tracks" "time_based" {
  genre = data.spotify_time.now.genre
  mood  = data.spotify_time.now.mood
  limit = 15
}

# Get tracks based on weather
data "spotify_tracks" "weather_based" {
  mood  = data.spotify_weather.current.mood
  limit = 15
}

# Create or update a time-based playlist
resource "spotify_playlist" "daily_time_mix" {
  name        = "${data.spotify_time.now.time_of_day} Mix - ${formatdate("YYYY-MM-DD", timestamp())}"
  description = "Tracks for ${data.spotify_time.now.time_of_day} vibes. Auto-updated daily."
  public      = true
  tracks      = data.spotify_tracks.time_based.ids
}

# Create or update a weather-based playlist
resource "spotify_playlist" "daily_weather_mix" {
  name        = "${data.spotify_weather.current.mood} Weather Mix - ${formatdate("YYYY-MM-DD", timestamp())}"
  description = "Tracks based on the current weather: ${data.spotify_weather.current.temperature}°C. Auto-updated daily."
  public      = true
  tracks      = data.spotify_tracks.weather_based.ids
}

# Add custom cover image to time-based playlist
resource "spotify_playlist_cover" "time_cover" {
  playlist_id = spotify_playlist.daily_time_mix.id
  mood        = data.spotify_time.now.mood
  force_update = true
}

# Add custom cover image to weather-based playlist
resource "spotify_playlist_cover" "weather_cover" {
  playlist_id = spotify_playlist.daily_weather_mix.id
  mood        = data.spotify_weather.current.mood
  force_update = true
}