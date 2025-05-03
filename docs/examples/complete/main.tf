# Complete Spotify Provider Example
# This example demonstrates using all resources and data sources

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
  weather_api_key = var.weather_api_key
}

# Get current user information
data "spotify_user" "me" {}

# Get current time data
data "spotify_time" "now" {}

# Get current weather data
data "spotify_weather" "current" {}

# Get user preferences
data "spotify_user_preference" "short_term" {
  time_range = "short_term"
  limit      = 10
}

data "spotify_user_preference" "long_term" {
  time_range = "long_term"
  limit      = 10
}

# Get featured playlists
data "spotify_featured_playlists" "featured" {
  limit = 5
}

# Get new releases
data "spotify_new_releases" "latest" {
  limit = 5
}

# Get tracks based on time of day
data "spotify_tracks" "time_based" {
  genre = data.spotify_time.now.genre
  mood  = data.spotify_time.now.mood
  limit = 10
}

# Get tracks based on weather
data "spotify_tracks" "weather_based" {
  mood  = data.spotify_weather.current.mood
  limit = 10
}

# Get tracks based on user preferences
data "spotify_tracks" "preference_based" {
  seed_artists = slice(data.spotify_user_preference.long_term.artist_ids, 0, 2)
  seed_tracks  = slice(data.spotify_user_preference.short_term.track_ids, 0, 2)
  limit        = 10
}

# Create a time-based playlist
resource "spotify_playlist" "time_playlist" {
  name        = "${data.spotify_time.now.time_of_day} Mix - ${formatdate("YYYY-MM-DD", timestamp())}"
  description = "Tracks for ${data.spotify_time.now.time_of_day} vibes. Created by Terraform."
  public      = true
  tracks      = data.spotify_tracks.time_based.ids
}

# Create a weather-based playlist
resource "spotify_playlist" "weather_playlist" {
  name        = "${data.spotify_weather.current.mood} Weather Mix - ${formatdate("YYYY-MM-DD", timestamp())}"
  description = "Music for ${data.spotify_weather.current.condition} weather at ${data.spotify_weather.current.temperature}°C. Created by Terraform."
  public      = true
  tracks      = data.spotify_tracks.weather_based.ids
}

# Create a personalized playlist based on user preferences
resource "spotify_playlist" "personalized_playlist" {
  name        = "${data.spotify_user.me.display_name}'s Personal Mix"
  description = "Based on your listening history. Created by Terraform."
  public      = false
  tracks      = data.spotify_tracks.preference_based.ids
}

# Add a custom cover image to the weather playlist
resource "spotify_playlist_cover" "weather_cover" {
  playlist_id     = spotify_playlist.weather_playlist.id
  mood            = data.spotify_weather.current.mood
  weather         = data.spotify_weather.current.condition
  force_update    = true
}

# Add a custom cover image to the time playlist
resource "spotify_playlist_cover" "time_cover" {
  playlist_id     = spotify_playlist.time_playlist.id
  mood            = data.spotify_time.now.mood
  background_color = "#3498DB"
  force_update    = true
}

# Add a custom cover image to the personalized playlist
resource "spotify_playlist_cover" "personalized_cover" {
  playlist_id     = spotify_playlist.personalized_playlist.id
  emoji           = "🎵"
  background_color = "#9B59B6"
}

# Output user information
output "user_info" {
  value = {
    name      = data.spotify_user.me.display_name
    id        = data.spotify_user.me.id
    followers = data.spotify_user.me.followers
  }
}

# Output created playlists
output "created_playlists" {
  value = {
    time_playlist = {
      name = spotify_playlist.time_playlist.name
      id   = spotify_playlist.time_playlist.id
    }
    weather_playlist = {
      name = spotify_playlist.weather_playlist.name
      id   = spotify_playlist.weather_playlist.id
    }
    personalized_playlist = {
      name = spotify_playlist.personalized_playlist.name
      id   = spotify_playlist.personalized_playlist.id
    }
  }
}

# Output current context
output "current_context" {
  value = {
    time_of_day = data.spotify_time.now.time_of_day
    weather     = data.spotify_weather.current.condition
    temperature = "${data.spotify_weather.current.temperature}°C"
  }
}