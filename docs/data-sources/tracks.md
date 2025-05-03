---
page_title: "spotify_tracks Data Source - terraform-provider-spotify"
subcategory: ""
description: |-
  Retrieves and recommends Spotify tracks based on various criteria.
---

# Data Source: spotify_tracks

Retrieves and recommends Spotify tracks based on various criteria such as genre, mood, artist, and more. This data source is useful for creating dynamic playlists with recommended tracks.

## Example Usage

```terraform
# Get tracks based on genre and mood
data "spotify_tracks" "chill_electronic" {
  genre = "electronic"
  mood  = "chill"
  limit = 20
}

# Get tracks based on seed tracks
data "spotify_tracks" "similar" {
  seed_tracks = ["spotify:track:4iV5W9uYEdYUVa79Axb7Rh", "spotify:track:1301WleyT98MSxVHPZCA6M"]
  limit = 15
}

# Get tracks based on artist
data "spotify_tracks" "artist_tracks" {
  seed_artists = ["spotify:artist:4Z8W4fKeB5YxbusRsdQVPb"]
  limit = 10
}

# Create a playlist with the recommended tracks
resource "spotify_playlist" "recommended" {
  name        = "Recommended Chill Electronic"
  description = "Tracks recommended based on electronic genre and chill mood"
  public      = true
  tracks      = data.spotify_tracks.chill_electronic.ids
}
```

## Dynamic Recommendations Example

```terraform
data "spotify_time" "now" {}
data "spotify_weather" "current" {}

data "spotify_tracks" "dynamic" {
  genre = data.spotify_time.now.genre
  mood  = data.spotify_weather.current.mood
  limit = 25
}

resource "spotify_playlist" "dynamic" {
  name        = "${data.spotify_weather.current.mood} ${data.spotify_time.now.time_of_day} Mix"
  description = "Auto-generated based on weather and time of day"
  public      = true
  tracks      = data.spotify_tracks.dynamic.ids
}
```

## Argument Reference

* `genre` - (Optional) A genre to use for recommendations.
* `mood` - (Optional) A mood to use for recommendations.
* `seed_tracks` - (Optional) A list of Spotify track URIs or IDs to use as seeds for recommendations.
* `seed_artists` - (Optional) A list of Spotify artist URIs or IDs to use as seeds for recommendations.
* `seed_genres` - (Optional) A list of genres to use as seeds for recommendations.
* `limit` - (Optional) The maximum number of tracks to return. Defaults to 20.
* `min_popularity` - (Optional) The minimum popularity of the tracks (0-100).
* `max_popularity` - (Optional) The maximum popularity of the tracks (0-100).
* `target_popularity` - (Optional) The target popularity of the tracks (0-100).
* `min_energy` - (Optional) The minimum energy of the tracks (0.0-1.0).
* `max_energy` - (Optional) The maximum energy of the tracks (0.0-1.0).
* `target_energy` - (Optional) The target energy of the tracks (0.0-1.0).
* `min_tempo` - (Optional) The minimum tempo of the tracks (BPM).
* `max_tempo` - (Optional) The maximum tempo of the tracks (BPM).
* `target_tempo` - (Optional) The target tempo of the tracks (BPM).

## Attribute Reference

* `id` - A unique identifier for this data source.
* `ids` - A list of Spotify track IDs.
* `uris` - A list of Spotify track URIs.
* `tracks` - A list of track objects with the following attributes:
  * `id` - The Spotify ID of the track.
  * `name` - The name of the track.
  * `artist` - The name of the primary artist.
  * `album` - The name of the album.
  * `popularity` - The popularity of the track (0-100).
  * `duration_ms` - The duration of the track in milliseconds.
  * `explicit` - Whether the track has explicit lyrics.
  * `uri` - The Spotify URI of the track.