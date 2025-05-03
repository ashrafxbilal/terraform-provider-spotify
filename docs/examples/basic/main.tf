# Basic Spotify Provider Example
# This example demonstrates creating a simple playlist

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

# Get current user information
data "spotify_user" "me" {}

# Create a simple playlist
resource "spotify_playlist" "basic_playlist" {
  name        = "My Terraform Playlist"
  description = "Created with Terraform by ${data.spotify_user.me.display_name}"
  public      = true
  
  # Add some tracks to the playlist
  # These are example track IDs - replace with your own
  tracks = [
    "spotify:track:4iV5W9uYEdYUVa79Axb7Rh", # Bohemian Rhapsody
    "spotify:track:1301WleyT98MSxVHPZCA6M", # Hotel California
    "spotify:track:3z8h0TU7ReDPLIbEnYhWZb"  # Bohemian Rhapsody
  ]
}

# Add a custom cover image to the playlist
resource "spotify_playlist_cover" "basic_cover" {
  playlist_id     = spotify_playlist.basic_playlist.id
  emoji           = "🎵"
  background_color = "#1DB954" # Spotify green
}

# Output the playlist information
output "playlist_info" {
  value = {
    name        = spotify_playlist.basic_playlist.name
    id          = spotify_playlist.basic_playlist.id
    description = spotify_playlist.basic_playlist.description
    tracks      = length(spotify_playlist.basic_playlist.tracks)
  }
}

# Output user information
output "user_info" {
  value = {
    name      = data.spotify_user.me.display_name
    id        = data.spotify_user.me.id
  }
}