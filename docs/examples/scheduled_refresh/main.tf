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
  genre = length(data.spotify_user_preferences.medium_term_taste.suggested_seed_genres) > 0 ? "${data.spotify_user_preferences.medium_term_taste.suggested_seed_genres[0]}${data.spotify_time.now.hour}${formatdate("YYYYMMDD", timestamp())}-${local.run_id}" : "pop${data.spotify_time.now.hour}${formatdate("YYYYMMDD", timestamp())}-${local.run_id}"
  mood  = data.spotify_weather.current.mood
  limit = 10 # Increased limit for more variety
}

# Tracks based on user's short-term 2nd top genre (for diversity)
data "spotify_tracks" "short_term_genre" {
  genre = length(data.spotify_user_preferences.short_term_taste.suggested_seed_genres) > 1 ? "${data.spotify_user_preferences.short_term_taste.suggested_seed_genres[1]}${formatdate("YYYYMMDD", timestamp())}-${local.run_id}" : "pop${formatdate("YYYYMMDD", timestamp())}-${local.run_id}"
  mood  = data.spotify_weather.current.mood
  limit = 15 # Increased limit for more variety
  
  # Add artist parameter to get featured tracks from favorite artists
  artist = length(data.spotify_user_preferences.short_term_taste.suggested_seed_artists) > 0 ? data.spotify_user_preferences.short_term_taste.suggested_seed_artists[0] : ""
}

# Tracks based on user's medium-term 3rd top genre (for diversity)
data "spotify_tracks" "medium_term_genre" {
  genre = length(data.spotify_user_preferences.medium_term_taste.suggested_seed_genres) > 2 ? "${data.spotify_user_preferences.medium_term_taste.suggested_seed_genres[2]}${formatdate("YYYYMMDD", timestamp())}-${local.run_id}" : "rock${formatdate("YYYYMMDD", timestamp())}-${local.run_id}"
  mood  = data.spotify_weather.current.mood
  limit = 15 # Increased limit for more variety
  
  # Add artist parameter to get featured tracks from another favorite artist
  artist = length(data.spotify_user_preferences.medium_term_taste.suggested_seed_artists) > 1 ? data.spotify_user_preferences.medium_term_taste.suggested_seed_artists[1] : ""
  
  # Set popularity lower to find less common tracks
  popularity = 60
}

# Tracks based on user's long-term 2nd top genre (for diversity)
data "spotify_tracks" "long_term_genre" {
  genre = length(data.spotify_user_preferences.long_term_taste.suggested_seed_genres) > 1 ? "${data.spotify_user_preferences.long_term_taste.suggested_seed_genres[1]}${formatdate("YYYYMMDD", timestamp())}-${local.run_id}" : "electronic${formatdate("YYYYMMDD", timestamp())}-${local.run_id}"
  mood  = data.spotify_weather.current.mood
  limit = 15 # Increased limit for more variety
}

# Tracks based on weather mood using user's 4th top genre
data "spotify_tracks" "weather_based" {
  genre = length(data.spotify_user_preferences.medium_term_taste.suggested_seed_genres) > 3 ? "${data.spotify_user_preferences.medium_term_taste.suggested_seed_genres[3]}${data.spotify_weather.current.temperature}${formatdate("YYYYMMDD", timestamp())}-${local.run_id}" : "pop${data.spotify_weather.current.temperature}${formatdate("YYYYMMDD", timestamp())}-${local.run_id}"
  mood  = data.spotify_weather.current.mood
  limit = 25 # Increased limit for more variety
}

# Tracks based on time of day using user's 3rd top genre
data "spotify_tracks" "time_based" {
  genre = length(data.spotify_user_preferences.long_term_taste.suggested_seed_genres) > 2 ? "${data.spotify_user_preferences.long_term_taste.suggested_seed_genres[2]}${data.spotify_time.now.time_of_day}${data.spotify_time.now.hour}${formatdate("YYYYMMDD", timestamp())}-${local.run_id}" : "rock${data.spotify_time.now.time_of_day}${data.spotify_time.now.hour}${formatdate("YYYYMMDD", timestamp())}-${local.run_id}"
  mood  = data.spotify_time.now.mood
  limit = 25 # Increased limit for more variety
}

# Tracks based on day of week using user's 4th top genre
data "spotify_tracks" "day_based" {
  genre = length(data.spotify_user_preferences.short_term_taste.suggested_seed_genres) > 3 ? "${data.spotify_user_preferences.short_term_taste.suggested_seed_genres[3]}${data.spotify_time.now.day_of_week}${data.spotify_time.now.hour}${formatdate("YYYYMMDD", timestamp())}-${local.run_id}" : "electronic${data.spotify_time.now.day_of_week}${data.spotify_time.now.hour}${formatdate("YYYYMMDD", timestamp())}-${local.run_id}"
  mood  = data.spotify_time.now.mood
  limit = 25 # Increased limit for more variety
}

###############################################################################
# Resources
###############################################################################


# Get tracks from Spotify's Global 50 playlist
data "spotify_tracks" "global_50_tracks" {
  # Always run this data source to ensure we get popular global tracks
  
  # Use pop genre for global hits with day of week, date, and unique ID to vary results
  genre = "pop${data.spotify_time.now.day_of_week}${formatdate("YYYYMMDD", timestamp())}-${local.run_id}"
  
  # Use energetic mood for popular tracks
  mood  = "energetic"
  
  # Get more tracks from global hits
  limit = 25 # Increased limit for more variety
  
  # Target popular tracks
  popularity = 90
}

# Dedicated source for English rock tracks
data "spotify_tracks" "english_rock_tracks" {
  # Use rock genre which is predominantly English, add time of day, date, and unique ID for variation
  genre = "rock${data.spotify_time.now.time_of_day}${formatdate("YYYYMMDD", timestamp())}-${local.run_id}"
  
  # Match current mood for consistency
  mood  = data.spotify_weather.current.mood
  
  # Ensure we get enough English tracks
  limit = 15 # Increased limit for more variety
  
  # Use medium-term history for stable preferences
  time_range = "medium_term"
}

# Dedicated source for English pop tracks
data "spotify_tracks" "english_pop_tracks" {
  # Use pop genre which is predominantly English, add temperature, date, and unique ID for variation
  genre = "pop${data.spotify_weather.current.temperature}${formatdate("YYYYMMDD", timestamp())}-${local.run_id}"
  
  # Match current mood for consistency
  mood  = data.spotify_weather.current.mood
  
  # Ensure we get enough English tracks
  limit = 15 # Increased limit for more variety
  
  # Use short-term history for recent preferences
  time_range = "short_term"
}

# Get tracks from new releases based on user's top genres (past month)
data "spotify_tracks" "new_releases_tracks" {
  # Always run this data source to ensure we get recent music
  
  # Use user's top genre for new releases with day, time, date, and unique ID for variation
  genre = length(data.spotify_user_preferences.short_term_taste.suggested_seed_genres) > 0 ? "${data.spotify_user_preferences.short_term_taste.suggested_seed_genres[0]}${data.spotify_time.now.day_of_week}${data.spotify_time.now.time_of_day}${formatdate("YYYYMMDD", timestamp())}-${local.run_id}" : "pop${data.spotify_time.now.day_of_week}${formatdate("YYYYMMDD", timestamp())}-${local.run_id}"
  
  # Use current mood
  mood  = data.spotify_weather.current.mood
  
  # Increase limit to get more recent tracks
  limit = 25 # Increased limit for more variety
  
  # Add time_range parameter to get recent tracks (past month)
  time_range = "short_term"
  
  # Target newer tracks that user might not have heard yet
  popularity = 70
}

# Generate a unique ID string for this run using timestamp
locals {
  # Create a unique ID using current timestamp with nanoseconds and a hash of the current time
  run_id = "${formatdate("YYYYMMDDhhmmss", timestamp())}-${sha1(timestamp())}"
}

# Add a dedicated source for discovering similar but unheard tracks
data "spotify_tracks" "discovery_tracks" {
  # Use a mix of user's top genres to find similar music
  genre = length(data.spotify_user_preferences.medium_term_taste.suggested_seed_genres) > 0 ? "${data.spotify_user_preferences.medium_term_taste.suggested_seed_genres[0]}_similar_${formatdate("YYYYMMDDHHmmss", timestamp())}-${local.run_id}" : "discover_${formatdate("YYYYMMDDHHmmss", timestamp())}-${local.run_id}"
  
  # Use opposite mood to find contrasting tracks
  mood  = data.spotify_weather.current.mood == "energetic" ? "chill" : data.spotify_weather.current.mood == "chill" ? "energetic" : data.spotify_weather.current.mood == "focus" ? "upbeat" : data.spotify_weather.current.mood == "upbeat" ? "focus" : data.spotify_weather.current.mood == "melancholy" ? "workout" : data.spotify_weather.current.mood == "workout" ? "melancholy" : data.spotify_weather.current.mood == "cozy" ? "energetic" : data.spotify_weather.current.mood == "romantic" ? "workout" : "chill"
  
  # Get a good number of discovery tracks
  limit = 30
  
  # Target tracks with lower popularity - more obscure, less likely to have been heard
  # This helps find tracks similar to user's taste but very likely not heard before
  popularity = 30
}

# Add a dedicated source for featured collaborations with favorite artists
data "spotify_tracks" "featured_collaborations" {
  # Use a different genre to find diverse collaborations
  genre = length(data.spotify_user_preferences.long_term_taste.suggested_seed_genres) > 2 ? "${data.spotify_user_preferences.long_term_taste.suggested_seed_genres[2]}_collab_${formatdate("YYYYMMDDHHmmss", timestamp())}-${local.run_id}" : "collab_${formatdate("YYYYMMDDHHmmss", timestamp())}-${local.run_id}"
  
  # Use top artists as seed but look for tracks where they're featured rather than primary
  artist = length(data.spotify_user_preferences.long_term_taste.suggested_seed_artists) > 0 ? data.spotify_user_preferences.long_term_taste.suggested_seed_artists[0] : ""
  
  # Set a different mood to get variety
  mood  = data.spotify_weather.current.mood == "energetic" ? "chill" : "energetic"
  
  # Get a good number of collaboration tracks
  limit = 20
  
  # Use a different time range to find more variety
  time_range = "long_term"
  
  # Lower popularity to find more obscure collaborations
  popularity = 40
}

# Add a source for underground/niche tracks in user's favorite genres
data "spotify_tracks" "underground_tracks" {
  # Use user's top genre but specifically target underground tracks
  genre = length(data.spotify_user_preferences.short_term_taste.suggested_seed_genres) > 0 ? "${data.spotify_user_preferences.short_term_taste.suggested_seed_genres[0]}_underground_${formatdate("YYYYMMDDHHmmss", timestamp())}-${local.run_id}" : "indie_underground_${formatdate("YYYYMMDDHHmmss", timestamp())}-${local.run_id}"
  
  # Use current mood for consistency
  mood  = data.spotify_weather.current.mood
  
  # Get a good number of underground tracks
  limit = 25
  
  # Very low popularity to find truly underground tracks
  popularity = 20
}

# Add a source for completely different genres from what the user typically listens to
data "spotify_tracks" "genre_exploration" {
  # Use genres that are likely different from user's typical taste
  # This helps introduce completely new sounds
  genre = length(data.spotify_user_preferences.medium_term_taste.suggested_seed_genres) > 0 ? "${can(regex(" ", data.spotify_user_preferences.medium_term_taste.suggested_seed_genres[0])) ? 
      "experimental_${element(split(" ", data.spotify_user_preferences.medium_term_taste.suggested_seed_genres[0]), 0)}" :
      data.spotify_user_preferences.medium_term_taste.suggested_seed_genres[0] == "pop" ? "ambient_experimental" :
      data.spotify_user_preferences.medium_term_taste.suggested_seed_genres[0] == "rock" ? "electronic_idm" :
      data.spotify_user_preferences.medium_term_taste.suggested_seed_genres[0] == "hip-hop" ? "classical_contemporary" :
      data.spotify_user_preferences.medium_term_taste.suggested_seed_genres[0] == "r&b" ? "folk_acoustic" :
      data.spotify_user_preferences.medium_term_taste.suggested_seed_genres[0] == "electronic" ? "jazz_fusion" :
      data.spotify_user_preferences.medium_term_taste.suggested_seed_genres[0] == "indie" ? "world_afrobeat" : "jazz_experimental"}_${formatdate("YYYYMMDDHHmmss", timestamp())}-${local.run_id}" : "world_fusion_${formatdate("YYYYMMDDHHmmss", timestamp())}-${local.run_id}"
  
  # Use a mood that complements exploration
  mood  = "focus"
  
  # Get a good number of exploration tracks
  limit = 20
  
  # Medium popularity to ensure quality while still being discoverable
  popularity = 50
}

resource "spotify_playlist" "combined_playlist" {
  name        = "${data.spotify_time.now.day_of_week} Surprise Mix for ${data.spotify_weather.current.temperature}°C"
  description = "Refreshed ${formatdate("YYYY-MM-DD HH:mm", timestamp())} by Terraform, it's giving ${data.spotify_time.now.mood} mood core for ${data.spotify_time.now.day_of_week} ${data.spotify_time.now.time_of_day} szn."
  public      = true
  
  # Ensure the playlist is created with all tracks before applying the cover
  lifecycle {
    create_before_destroy = true
    # Force recreation of the playlist on each apply to ensure fresh tracks
    ignore_changes = []
  }
  
  # Combine all track sources with a focus on discovery and featured tracks
  # The distinct function ensures no duplicate tracks
  tracks = distinct(concat(
    # Highest priority for surprising and underground tracks
    data.spotify_tracks.underground_tracks.ids,                            # True underground tracks from favorite genres
    data.spotify_tracks.genre_exploration.ids,                             # Completely different genres for exploration
    data.spotify_tracks.discovery_tracks.ids,                              # Contrasting mood discovery tracks
    data.spotify_tracks.featured_collaborations.ids,                       # Obscure collaborations with favorite artists
    
    # Secondary priority for other discovery sources
    data.spotify_tracks.medium_term_genre.ids,                             # Less common tracks from favorite artists
    data.spotify_tracks.new_releases_tracks.ids,                           # Latest releases from top genre
    
    # Mix in other diverse sources (limited quantity to prioritize discovery)
    # Use conditional expressions to safely handle empty or small lists
    length(data.spotify_tracks.short_term_genre.ids) > 0 ? slice(data.spotify_tracks.short_term_genre.ids, 0, min(5, length(data.spotify_tracks.short_term_genre.ids))) : [],
    length(data.spotify_tracks.english_rock_tracks.ids) > 0 ? slice(data.spotify_tracks.english_rock_tracks.ids, 0, min(5, length(data.spotify_tracks.english_rock_tracks.ids))) : [],
    length(data.spotify_tracks.english_pop_tracks.ids) > 0 ? slice(data.spotify_tracks.english_pop_tracks.ids, 0, min(5, length(data.spotify_tracks.english_pop_tracks.ids))) : [],
    length(data.spotify_tracks.global_50_tracks.ids) > 0 ? slice(data.spotify_tracks.global_50_tracks.ids, 0, min(3, length(data.spotify_tracks.global_50_tracks.ids))) : [],
    length(data.spotify_tracks.top_genre_limited.ids) > 0 ? slice(data.spotify_tracks.top_genre_limited.ids, 0, min(3, length(data.spotify_tracks.top_genre_limited.ids))) : [],
    length(data.spotify_tracks.long_term_genre.ids) > 0 ? slice(data.spotify_tracks.long_term_genre.ids, 0, min(5, length(data.spotify_tracks.long_term_genre.ids))) : [],
    length(data.spotify_tracks.weather_based.ids) > 0 ? slice(data.spotify_tracks.weather_based.ids, 0, min(5, length(data.spotify_tracks.weather_based.ids))) : [],
    length(data.spotify_tracks.time_based.ids) > 0 ? slice(data.spotify_tracks.time_based.ids, 0, min(5, length(data.spotify_tracks.time_based.ids))) : [],
    length(data.spotify_tracks.day_based.ids) > 0 ? slice(data.spotify_tracks.day_based.ids, 0, min(5, length(data.spotify_tracks.day_based.ids))) : []
  ))
}

###############################################################################
# Playlist Cover Image
###############################################################################

# Create a dynamic cover image for the surprise playlist
resource "spotify_playlist_cover" "dynamic_cover" {
  playlist_id = spotify_playlist.combined_playlist.id
  
  # Use a custom emoji that represents exploration and surprise
  emoji = "🔮"
  
  # Add weather condition to enhance the cover image
  weather = data.spotify_weather.current.is_sunny ? "sunny" : "cloudy"
  
  # Use bold, unexpected colors that change with each refresh
  # These colors are more experimental and vibrant than standard mood colors
  background_color = lookup({
    "energetic" = "#FF00FF",  # Electric magenta for energy
    "chill"     = "#00FFCC",  # Bright turquoise for chill vibes
    "cozy"      = "#FF6600",  # Burnt orange for cozy feelings
    "melancholy" = "#6600CC", # Deep purple for melancholy
    "upbeat"    = "#FFFF00",  # Electric yellow for upbeat
    "focus"     = "#9900FF",  # Vivid violet for focus
    "workout"   = "#00FF66",  # Neon green for workout
    "romantic"  = "#FF0066",  # Hot pink for romantic
  }, data.spotify_weather.current.mood, "#1DB954") # Default to Spotify green

  # Force update ensures the cover is refreshed even when mood hasn't changed in Terraform state
  force_update = true
  
  # This resource is automatically created with the playlist
  lifecycle {
    create_before_destroy = true
  }
}
