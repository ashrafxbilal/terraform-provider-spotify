###############################################################################
# Variables for Spotify Personalized Playlist Generator
#
# This file contains all variable definitions used in the main.tf configuration
###############################################################################

variable "spotify_client_id" {
  description = "Spotify API Client ID"
  type        = string
  sensitive   = true
}

variable "spotify_client_secret" {
  description = "Spotify API Client Secret"
  type        = string
  sensitive   = true
}

variable "spotify_redirect_uri" {
  description = "Spotify API Redirect URI"
  type        = string
}

variable "spotify_refresh_token" {
  description = "Spotify API Refresh Token"
  type        = string
  sensitive   = true
}