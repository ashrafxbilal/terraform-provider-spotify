###############################################################################
# Terraform Configuration for Spotify Personalized Playlist Generator
#
# This configuration creates a diverse personalized playlist based on:
# - User's listening history (short, medium, and long-term)
# - Current weather conditions
# - Time of day and day of week
###############################################################################

###############################################################################
# Provider Configuration
###############################################################################

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

###############################################################################
# Data Sources - Contextual Information
###############################################################################

data "spotify_user" "me" {}

data "spotify_weather" "current" {
  ## Uncomment and set your preferred mood from the suggested ones or your own custom mood
  ## Available moods now include: energetic, chill, cozy, melancholy, upbeat, focus, workout, romantic
  # mood = "focus"  # Example of using the new 'focus' mood
}

data "spotify_time" "now" {
  ## Uncomment and set your preferred mood from the suggested ones or your own custom mood
  # mood = "your_custom_mood"
  
  ## Uncomment and set your preferred genre from the suggested ones or your own custom genre
  # genre = "your_custom_genre"
}

###############################################################################
# Data Sources - User Preferences
###############################################################################

data "spotify_user_preferences" "short_term_taste" {
  time_range = "short_term"
}

data "spotify_user_preferences" "medium_term_taste" {
  time_range = "medium_term"
}

data "spotify_user_preferences" "long_term_taste" {
  time_range = "long_term" 
}

###############################################################################
# Data Sources - Track Recommendations
###############################################################################

# Small selection from user's absolute top genre
data "spotify_tracks" "top_genre_limited" {
  genre = length(data.spotify_user_preferences.medium_term_taste.suggested_seed_genres) > 0 ? "${data.spotify_user_preferences.medium_term_taste.suggested_seed_genres[0]}${data.spotify_time.now.hour}" : "pop${data.spotify_time.now.hour}"
  mood  = data.spotify_weather.current.mood
  limit = 5
}

# Tracks based on user's short-term 2nd top genre (for diversity)
data "spotify_tracks" "short_term_genre" {
  genre = length(data.spotify_user_preferences.short_term_taste.suggested_seed_genres) > 1 ? data.spotify_user_preferences.short_term_taste.suggested_seed_genres[1] : "pop"
  mood  = data.spotify_weather.current.mood
  limit = 10 # Lower limit to avoid API errors
}

# Tracks based on user's medium-term 3rd top genre (for diversity)
data "spotify_tracks" "medium_term_genre" {
  genre = length(data.spotify_user_preferences.medium_term_taste.suggested_seed_genres) > 2 ? data.spotify_user_preferences.medium_term_taste.suggested_seed_genres[2] : "rock"
  mood  = data.spotify_weather.current.mood
  limit = 10
}

# Tracks based on user's long-term 2nd top genre (for diversity)
data "spotify_tracks" "long_term_genre" {
  genre = length(data.spotify_user_preferences.long_term_taste.suggested_seed_genres) > 1 ? data.spotify_user_preferences.long_term_taste.suggested_seed_genres[1] : "electronic"
  mood  = data.spotify_weather.current.mood
  limit = 10
}

# Tracks based on weather mood using user's 4th top genre
data "spotify_tracks" "weather_based" {
  genre = length(data.spotify_user_preferences.medium_term_taste.suggested_seed_genres) > 3 ? "${data.spotify_user_preferences.medium_term_taste.suggested_seed_genres[3]}${data.spotify_weather.current.temperature}" : "pop${data.spotify_weather.current.temperature}"
  mood  = data.spotify_weather.current.mood
  limit = 15
}

# Tracks based on time of day using user's 3rd top genre
data "spotify_tracks" "time_based" {
  genre = length(data.spotify_user_preferences.long_term_taste.suggested_seed_genres) > 2 ? "${data.spotify_user_preferences.long_term_taste.suggested_seed_genres[2]}${data.spotify_time.now.time_of_day}${data.spotify_time.now.hour}" : "rock${data.spotify_time.now.time_of_day}${data.spotify_time.now.hour}"
  mood  = data.spotify_time.now.mood
  limit = 15 
}

# Tracks based on day of week using user's 4th top genre
data "spotify_tracks" "day_based" {
  genre = length(data.spotify_user_preferences.short_term_taste.suggested_seed_genres) > 3 ? "${data.spotify_user_preferences.short_term_taste.suggested_seed_genres[3]}${data.spotify_time.now.day_of_week}${data.spotify_time.now.hour}" : "electronic${data.spotify_time.now.day_of_week}${data.spotify_time.now.hour}"
  mood  = data.spotify_time.now.mood
  limit = 15 
}

###############################################################################
# Resources
###############################################################################


# Get tracks from Spotify's Global 50 playlist
data "spotify_tracks" "global_50_tracks" {
  # Always run this data source to ensure we get popular global tracks
  
  # Use pop genre for global hits with day of week to vary results
  genre = "pop${data.spotify_time.now.day_of_week}"
  
  # Use energetic mood for popular tracks
  mood  = "energetic"
  
  # Get more tracks from global hits
  limit = 15
  
  # Target popular tracks
  popularity = 90
}

# Dedicated source for English rock tracks
data "spotify_tracks" "english_rock_tracks" {
  # Use rock genre which is predominantly English, add time of day for variation
  genre = "rock${data.spotify_time.now.time_of_day}"
  
  # Match current mood for consistency
  mood  = data.spotify_weather.current.mood
  
  # Ensure we get enough English tracks
  limit = 8
  
  # Use medium-term history for stable preferences
  time_range = "medium_term"
}

# Dedicated source for English pop tracks
data "spotify_tracks" "english_pop_tracks" {
  # Use pop genre which is predominantly English, add temperature for variation
  genre = "pop${data.spotify_weather.current.temperature}"
  
  # Match current mood for consistency
  mood  = data.spotify_weather.current.mood
  
  # Ensure we get enough English tracks
  limit = 8
  
  # Use short-term history for recent preferences
  time_range = "short_term"
}

# Get tracks from new releases based on user's top genres (past month)
data "spotify_tracks" "new_releases_tracks" {
  # Always run this data source to ensure we get recent music
  
  # Use user's top genre for new releases with day and time for variation
  genre = length(data.spotify_user_preferences.short_term_taste.suggested_seed_genres) > 0 ? "${data.spotify_user_preferences.short_term_taste.suggested_seed_genres[0]}${data.spotify_time.now.day_of_week}${data.spotify_time.now.time_of_day}" : "pop${data.spotify_time.now.day_of_week}"
  
  # Use current mood
  mood  = data.spotify_weather.current.mood
  
  # Increase limit to get more recent tracks
  limit = 15
  
  # Add time_range parameter to get recent tracks (past month)
  time_range = "short_term"
}

resource "spotify_playlist" "combined_playlist" {
  name        = "${data.spotify_time.now.day_of_week} Vibes @ ${data.spotify_weather.current.temperature}°C ✨"
  description = "it's giving ${data.spotify_time.now.mood} mood core for ${data.spotify_time.now.day_of_week} ${data.spotify_time.now.time_of_day} szn."
  public      = true
  
  # Ensure the playlist is created with all tracks before applying the cover
  lifecycle {
    create_before_destroy = true
  }
  
  # Combine all track sources, prioritizing English tracks, new releases and global hits
  tracks = distinct(concat(
    data.spotify_tracks.english_rock_tracks.ids,                             # Dedicated English rock tracks
    data.spotify_tracks.english_pop_tracks.ids,                              # Dedicated English pop tracks
    data.spotify_tracks.new_releases_tracks.ids,                              # Latest releases from user's top genre
    data.spotify_tracks.global_50_tracks.ids,                                # Global 50 popular tracks
    data.spotify_tracks.top_genre_limited.ids,                                # Small selection from top genre
    data.spotify_tracks.short_term_genre.ids,                                 # 2nd favorite genre
    data.spotify_tracks.medium_term_genre.ids,                                # 3rd favorite genre
    data.spotify_tracks.long_term_genre.ids,                                  # 2nd long-term genre
    data.spotify_tracks.weather_based.ids,                                    # 4th favorite genre
    data.spotify_tracks.time_based.ids,                                       # 3rd long-term genre
    data.spotify_tracks.day_based.ids                                         # 4th short-term genre
  ))
}

