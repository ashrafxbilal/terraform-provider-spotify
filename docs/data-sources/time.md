---
page_title: "spotify_time Data Source - terraform-provider-spotify"
subcategory: ""
description: |-
  Provides time-based information and suggests moods and genres based on time of day.
---

# Data Source: spotify_time

Provides time-based information and suggests moods and genres based on time of day. This data source is useful for creating dynamic playlists that change based on the time of day or day of the week.

## Example Usage

```terraform
data "spotify_time" "now" {}

output "time_of_day" {
  value = data.spotify_time.now.time_of_day
}

output "suggested_moods" {
  value = data.spotify_time.now.suggested_moods
}

output "suggested_genres" {
  value = data.spotify_time.now.suggested_genres
}

resource "spotify_playlist" "time_based" {
  name        = "${data.spotify_time.now.time_of_day} Vibes - ${formatdate("YYYY-MM-DD", timestamp())}"
  description = "Music for ${data.spotify_time.now.time_of_day} listening"
  public      = true
}
```

## Argument Reference

This data source has no required arguments.

* `custom_time` - (Optional) A custom time to use instead of the current time. Format: "HH:MM".
* `custom_day` - (Optional) A custom day to use instead of the current day. Valid values: "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday".

## Attribute Reference

* `id` - A unique identifier for this data source.
* `hour` - The current hour (0-23).
* `minute` - The current minute (0-59).
* `day_of_week` - The current day of the week ("monday", "tuesday", etc.).
* `is_weekend` - Whether the current day is a weekend day.
* `time_of_day` - The time of day category ("morning", "afternoon", "evening", "night").
* `mood` - A suggested mood based on the time of day.
* `genre` - A suggested genre based on the time of day and day of week.
* `suggested_moods` - A list of suggested moods appropriate for the current time of day.
* `suggested_genres` - A list of suggested genres appropriate for the current time of day.