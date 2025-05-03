---
page_title: "spotify_playlist Resource - terraform-provider-spotify"
subcategory: ""
description: |-
  Manages a Spotify playlist.
---

# Resource: spotify_playlist

Manages a Spotify playlist. This resource allows you to create, update, and delete playlists in your Spotify account, as well as manage the tracks within them.

## Example Usage

```terraform
resource "spotify_playlist" "example" {
  name        = "My Terraform Playlist"
  description = "Created and managed by Terraform"
  public      = true
  tracks      = ["spotify:track:4iV5W9uYEdYUVa79Axb7Rh", "spotify:track:1301WleyT98MSxVHPZCA6M"]
}
```

## Dynamic Playlist Example

```terraform
data "spotify_time" "now" {}
data "spotify_weather" "current" {}
data "spotify_tracks" "recommended" {
  genre = data.spotify_time.now.genre
  mood  = data.spotify_weather.current.mood
  limit = 20
}

resource "spotify_playlist" "dynamic" {
  name        = "${data.spotify_weather.current.mood} ${data.spotify_time.now.time_of_day} Mix"
  description = "Auto-generated based on weather and time of day"
  public      = true
  tracks      = data.spotify_tracks.recommended.ids
}
```

## Argument Reference

* `name` - (Required) The name of the playlist.
* `description` - (Optional) The description of the playlist.
* `public` - (Optional) Whether the playlist is public. Defaults to `true`.
* `tracks` - (Optional) A list of Spotify track URIs or IDs to add to the playlist.
* `collaborative` - (Optional) Whether the playlist is collaborative. Defaults to `false`.

## Attribute Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - The Spotify ID of the playlist.
* `snapshot_id` - The current snapshot ID of the playlist.
* `owner` - The Spotify user ID of the playlist owner.
* `followers` - The number of followers the playlist has.
* `images` - A list of image URLs associated with the playlist.

## Import

Spotify playlists can be imported using the playlist ID, e.g.,

```
$ terraform import spotify_playlist.example 3cEYpjA9oz9GiPac4AsH4n
```