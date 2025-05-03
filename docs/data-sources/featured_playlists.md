---
page_title: "spotify_featured_playlists Data Source - terraform-provider-spotify"
subcategory: ""
description: |-
  Retrieves Spotify's featured playlists.
---

# Data Source: spotify_featured_playlists

Retrieves Spotify's featured playlists, which are curated by Spotify and often change based on time of day, day of week, and region.

## Example Usage

```terraform
data "spotify_featured_playlists" "featured" {
  limit = 10
}

output "featured_playlist_names" {
  value = [for playlist in data.spotify_featured_playlists.featured.playlists : playlist.name]
}

output "featured_message" {
  value = data.spotify_featured_playlists.featured.message
}
```

## Country-Specific Example

```terraform
data "spotify_featured_playlists" "uk_featured" {
  country = "GB"
  limit   = 5
}

output "uk_featured_playlists" {
  value = data.spotify_featured_playlists.uk_featured.playlists
}
```

## Argument Reference

* `country` - (Optional) An ISO 3166-1 alpha-2 country code to get featured playlists for a specific country.
* `locale` - (Optional) The desired language, consisting of an ISO 639-1 language code and an ISO 3166-1 alpha-2 country code, joined by an underscore (e.g., "es_MX").
* `timestamp` - (Optional) A timestamp in ISO 8601 format to get featured playlists for a specific time.
* `limit` - (Optional) The maximum number of playlists to return. Default: 20, Maximum: 50.

## Attribute Reference

* `id` - A unique identifier for this data source.
* `message` - The message that accompanies the featured playlists.
* `playlists` - A list of featured playlist objects with the following attributes:
  * `id` - The Spotify ID of the playlist.
  * `name` - The name of the playlist.
  * `description` - The description of the playlist.
  * `owner` - The Spotify user ID of the playlist owner.
  * `public` - Whether the playlist is public.
  * `collaborative` - Whether the playlist is collaborative.
  * `tracks_total` - The total number of tracks in the playlist.
  * `image_url` - The URL of the playlist's cover image.
  * `url` - The Spotify URL of the playlist.
* `playlist_ids` - A list of Spotify playlist IDs.