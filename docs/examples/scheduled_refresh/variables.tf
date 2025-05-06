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
}

variable "spotify_refresh_token" {
  description = "Spotify API refresh token"
  type        = string
  sensitive   = true
}