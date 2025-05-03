# Variables for the Spotify Provider

variable "spotify_client_id" {
  description = "Spotify API client ID"
  type        = string
  sensitive   = true
}

variable "spotify_client_secret" {
  description = "Spotify API client secret"
  type        = string
  sensitive   = true
}

variable "spotify_redirect_uri" {
  description = "Spotify API redirect URI"
  type        = string
  default     = "http://localhost:8080/callback"
}

variable "spotify_refresh_token" {
  description = "Spotify API refresh token"
  type        = string
  sensitive   = true
}

variable "weather_api_key" {
  description = "OpenWeatherMap API key for weather-based playlists"
  type        = string
  sensitive   = true
  default     = ""
}