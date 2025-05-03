###############################################################################
# Playlist Cover Image Configuration
###############################################################################

# Create a dynamic cover image for the playlist based on mood and weather
# This is automatically applied when the playlist is created
resource "spotify_playlist_cover" "dynamic_cover" {
  playlist_id = spotify_playlist.combined_playlist.id
  
  # Use the current weather mood for the emoji
  mood = data.spotify_weather.current.mood
  
  # Use a background color that matches the mood
  background_color = lookup({
    "energetic" = "#FF5733",  # Bright orange-red for energy
    "chill"     = "#33A1FF",  # Calm blue for chill vibes
    "cozy"      = "#A15E49",  # Warm brown for cozy feelings
    "melancholy" = "#6E6E6E", # Gray for melancholy
    "upbeat"    = "#FFCC33",  # Sunny yellow for upbeat
    "focus"     = "#9933FF",  # Purple for focus
    "workout"   = "#33FF57",  # Green for workout
    "romantic"  = "#FF33A1",  # Pink for romantic
  }, data.spotify_weather.current.mood, "#1DB954") # Default to Spotify green

  force_update    = true
  
  # This resource is automatically created with the playlist
  # No separate terraform apply needed
  lifecycle {
    create_before_destroy = true
  }
}

# Output the playlist URL with cover image
output "playlist_with_cover" {
  value       = spotify_playlist.combined_playlist.spotify_url
  description = "URL to the playlist with the custom cover image"
}