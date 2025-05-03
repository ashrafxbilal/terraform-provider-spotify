---
page_title: "spotify_user_preference Data Source - terraform-provider-spotify"
subcategory: ""
description: |-
  Retrieves user preferences and listening habits from Spotify.
---

# Data Source: spotify_user_preference

Retrieves user preferences and listening habits from Spotify, including top artists, tracks, and genres. This data source is useful for creating personalized playlists based on a user's listening history.

## Example Usage

```terraform
data "spotify_user_preference" "my_preferences" {
  time_range = "medium_term"
}

output "top_artists" {
  value = data.spotify_user_preference.my_preferences.top_artists
}

output "top_genres" {
  value = data.spotify_user_preference.my_preferences.top_genres
}

resource "spotify_playlist" "personalized" {
  name        = "My Personalized Mix"
  description = "Based on my top genres: ${join(", ", slice(data.spotify_user_preference.my_preferences.top_genres, 0, 3))}"
  public      = true
}
```

## Multiple Time Ranges Example

```terraform
data "spotify_user_preference" "short_term" {
  time_range = "short_term"
}

data "spotify_user_preference" "medium_term" {
  time_range = "medium_term"
}

data "spotify_user_preference" "long_term" {
  time_range = "long_term"
}

output "recent_favorites" {
  value = data.spotify_user_preference.short_term.top_tracks
}

output "consistent_favorites" {
  value = data.spotify_user_preference.long_term.top_tracks
}
```

## Argument Reference

* `time_range` - (Optional) The time range to use for top artists and tracks. Valid values: "short_term" (approximately last 4 weeks), "medium_term" (approximately last 6 months), or "long_term" (calculated from several years of data and including all new data as it becomes available). Default: "medium_term".
* `limit` - (Optional) The maximum number of items to return for top artists and tracks. Default: 20, Maximum: 50.

## Attribute Reference

* `id` - A unique identifier for this data source.
* `top_artists` - A list of the user's top artists with the following attributes:
  * `id` - The Spotify ID of the artist.
  * `name` - The name of the artist.
  * `genres` - A list of genres associated with the artist.
  * `popularity` - The popularity of the artist (0-100).
  * `image_url` - The URL of the artist's image.
  * `url` - The Spotify URL of the artist.
* `top_tracks` - A list of the user's top tracks with the following attributes:
  * `id` - The Spotify ID of the track.
  * `name` - The name of the track.
  * `artist` - The name of the primary artist.
  * `album` - The name of the album.
  * `popularity` - The popularity of the track (0-100).
  * `uri` - The Spotify URI of the track.
* `top_genres` - A list of the user's top genres based on their top artists.
* `artist_ids` - A list of Spotify artist IDs for the user's top artists.
* `track_ids` - A list of Spotify track IDs for the user's top tracks.