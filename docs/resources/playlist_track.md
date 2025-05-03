---
page_title: "spotify_playlist_track Resource - terraform-provider-spotify"
subcategory: ""
description: |-
  Manages tracks in a Spotify playlist.
---

# Resource: spotify_playlist_track

Manages tracks in a Spotify playlist. This resource allows you to add, remove, and reorder tracks in an existing Spotify playlist.

## Example Usage

```terraform
resource "spotify_playlist" "example" {
  name        = "My Terraform Playlist"
  description = "Created and managed by Terraform"
  public      = true
}

resource "spotify_playlist_track" "track1" {
  playlist_id = spotify_playlist.example.id
  track_id    = "spotify:track:4iV5W9uYEdYUVa79Axb7Rh"
  position    = 0
}

resource "spotify_playlist_track" "track2" {
  playlist_id = spotify_playlist.example.id
  track_id    = "spotify:track:1301WleyT98MSxVHPZCA6M"
  position    = 1
}
```

## Argument Reference

* `playlist_id` - (Required) The Spotify ID of the playlist.
* `track_id` - (Required) The Spotify track URI or ID to add to the playlist.
* `position` - (Optional) The position to insert the track in the playlist (0-based index). If not specified, the track will be added to the end of the playlist.

## Attribute Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - A composite ID in the format `{playlist_id}:{track_id}`.
* `added_at` - The timestamp when the track was added to the playlist.
* `added_by` - The Spotify user ID of the user who added the track.
* `is_local` - Whether the track is a local file.

## Import

Spotify playlist tracks can be imported using the composite ID, e.g.,

```
$ terraform import spotify_playlist_track.example 3cEYpjA9oz9GiPac4AsH4n:4iV5W9uYEdYUVa79Axb7Rh
```