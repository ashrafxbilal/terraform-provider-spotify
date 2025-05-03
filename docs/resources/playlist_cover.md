---
page_title: "spotify_playlist_cover Resource - terraform-provider-spotify"
subcategory: ""
description: |-
  Manages a Spotify playlist cover image.
---

# Resource: spotify_playlist_cover

Manages a Spotify playlist cover image. This resource allows you to create and update custom cover images for Spotify playlists, including dynamic images based on mood, weather, and other contextual factors.

## Example Usage

```terraform
resource "spotify_playlist" "example" {
  name        = "My Playlist with Custom Cover"
  description = "A playlist with a custom cover image"
  public      = true
}

resource "spotify_playlist_cover" "example_cover" {
  playlist_id     = spotify_playlist.example.id
  background_color = "#FF5500"
  emoji           = "🎵"
}
```

## Dynamic Cover Example

```terraform
data "spotify_weather" "current" {}
data "spotify_time" "now" {}

resource "spotify_playlist" "dynamic" {
  name        = "${data.spotify_weather.current.mood} ${data.spotify_time.now.time_of_day} Mix"
  description = "Auto-generated based on weather and time of day"
  public      = true
}

resource "spotify_playlist_cover" "dynamic_cover" {
  playlist_id     = spotify_playlist.dynamic.id
  mood            = data.spotify_weather.current.mood
  weather         = data.spotify_weather.current.condition
  force_update    = true
}
```

## Argument Reference

* `playlist_id` - (Required) The Spotify ID of the playlist.
* `image_url` - (Optional) A URL to an image to use as the playlist cover. The image must be a JPEG and less than 256KB in size.
* `emoji` - (Optional) An emoji to use for generating a cover image.
* `mood` - (Optional) A mood to use for generating a cover image (e.g., "energetic", "chill", "melancholy").
* `weather` - (Optional) A weather condition to use for generating a cover image (e.g., "sunny", "rainy", "cloudy").
* `background_color` - (Optional) A hex color code for the background of the generated cover image.
* `force_update` - (Optional) Force update the playlist cover image even if no changes are detected. Defaults to `false`.

## Attribute Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - The Spotify ID of the playlist.

## Notes

You must provide at least one of `image_url`, `emoji`, `mood`, or `weather` to generate a cover image. If multiple options are provided, they will be used in the following order of precedence:

1. `image_url`
2. `emoji`
3. `mood`
4. `weather`

The `force_update` parameter is useful for ensuring the cover image is updated when using dynamic values that may not change in Terraform's state but should trigger an update (such as weather conditions that change externally).