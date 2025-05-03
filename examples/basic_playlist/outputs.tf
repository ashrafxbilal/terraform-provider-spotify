###############################################################################
# Outputs for Spotify Personalized Playlist Generator
#
# This file contains all output declarations used in the main.tf configuration
###############################################################################

###############################################################################
# Outputs - Contextual Suggestions
###############################################################################

# Suggested moods from weather
output "suggested_weather_moods" {
  value       = data.spotify_weather.current.suggested_moods
  description = "Choose one of these moods or set your own in the weather data source"
}

# Suggested moods from time
output "suggested_time_moods" {
  value       = data.spotify_time.now.suggested_moods
  description = "Choose one of these moods or set your own in the time data source"
}

# Suggested genres from time
output "suggested_time_genres" {
  value       = data.spotify_time.now.suggested_genres
  description = "Choose one of these genres or set your own in the time data source"
}

# Global hits information
output "global_hits" {
  value       = data.spotify_tracks.global_50_tracks.names
  description = "List of popular global tracks included in the playlist"
}

# New releases information
output "new_releases" {
  value       = data.spotify_tracks.new_releases_tracks.names
  description = "List of new releases from your favorite artists/genres"
}

# Track artists
output "track_artists" {
  value       = distinct(concat(
    data.spotify_tracks.global_50_tracks.artists,
    data.spotify_tracks.new_releases_tracks.artists
  ))
  description = "Artists featured in the global hits and new releases"
}

###############################################################################
# Outputs - User Preferences
###############################################################################

output "user_id" {
  value = data.spotify_user.me.id
}

# top genres from different time ranges
output "short_term_top_genres" {
  value       = data.spotify_user_preferences.short_term_taste.top_genres
  description = "Your top genres from the last 4 weeks"
}

output "medium_term_top_genres" {
  value       = data.spotify_user_preferences.medium_term_taste.top_genres
  description = "Your top genres from the last 6 months"
}

output "long_term_top_genres" {
  value       = data.spotify_user_preferences.long_term_taste.top_genres
  description = "Your top genres from several years"
}

# top artists from different time ranges
output "short_term_top_artists" {
  value       = data.spotify_user_preferences.short_term_taste.top_artists
  description = "Your top artists from the last 4 weeks"
}

output "medium_term_top_artists" {
  value       = data.spotify_user_preferences.medium_term_taste.top_artists
  description = "Your top artists from the last 6 months"
}

output "long_term_top_artists" {
  value       = data.spotify_user_preferences.long_term_taste.top_artists
  description = "Your top artists from several years"
}

###############################################################################
# Outputs - Current Context
###############################################################################

output "weather_mood" {
  value = data.spotify_weather.current.mood
}

output "time_mood" {
  value = data.spotify_time.now.mood
}