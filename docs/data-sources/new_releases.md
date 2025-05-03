---
page_title: "spotify_new_releases Data Source - terraform-provider-spotify"
subcategory: ""
description: |-
  Retrieves new album releases from Spotify.
---

# Data Source: spotify_new_releases

Retrieves new album releases from Spotify. This data source is useful for creating playlists with the latest music or staying up-to-date with new releases.

## Example Usage

```terraform
data "spotify_new_releases" "latest" {
  limit = 10
}

output "new_album_names" {
  value = [for album in data.spotify_new_releases.latest.albums : album.name]
}

output "new_album_artists" {
  value = [for album in data.spotify_new_releases.latest.albums : album.artists]
}
```

## Country-Specific Example

```terraform
data "spotify_new_releases" "us_releases" {
  country = "US"
  limit   = 5
}

output "us_new_releases" {
  value = data.spotify_new_releases.us_releases.albums
}
```

## Argument Reference

* `country` - (Optional) An ISO 3166-1 alpha-2 country code to get new releases for a specific country.
* `limit` - (Optional) The maximum number of new releases to return. Default: 20, Maximum: 50.

## Attribute Reference

* `id` - A unique identifier for this data source.
* `albums` - A list of new release album objects with the following attributes:
  * `id` - The Spotify ID of the album.
  * `name` - The name of the album.
  * `type` - The type of the album (e.g., "album", "single").
  * `artists` - A comma-separated list of the artists' names.
  * `release_date` - The release date of the album.
  * `release_date_precision` - The precision of the release date ("year", "month", or "day").
  * `image_url` - The URL of the album's cover image.
  * `total_tracks` - The total number of tracks in the album.
  * `url` - The Spotify URL of the album.
* `album_ids` - A list of Spotify album IDs.