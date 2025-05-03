---
page_title: "spotify_weather Data Source - terraform-provider-spotify"
subcategory: ""
description: |-
  Retrieves current weather information and suggests moods based on weather conditions.
---

# Data Source: spotify_weather

Retrieves current weather information and suggests moods based on weather conditions. This data source is useful for creating dynamic playlists that change based on the current weather.

## Example Usage

```terraform
data "spotify_weather" "current" {}

output "weather_condition" {
  value = data.spotify_weather.current.condition
}

output "weather_mood" {
  value = data.spotify_weather.current.mood
}

output "temperature" {
  value = "${data.spotify_weather.current.temperature}°C"
}

resource "spotify_playlist" "weather_based" {
  name        = "${data.spotify_weather.current.mood} Weather Mix"
  description = "Music for ${data.spotify_weather.current.condition} weather at ${data.spotify_weather.current.temperature}°C"
  public      = true
}
```

## Custom Location Example

```terraform
data "spotify_weather" "paris" {
  location = {
    city = "Paris"
    country = "FR"
  }
}

output "paris_weather" {
  value = "${data.spotify_weather.paris.condition} at ${data.spotify_weather.paris.temperature}°C"
}
```

## Argument Reference

* `location` - (Optional) The location to get weather for. If not specified, the provider will attempt to determine the current location.
  * `city` - (Optional) The city name.
  * `country` - (Optional) The country code.
  * `lat` - (Optional) The latitude coordinate.
  * `lon` - (Optional) The longitude coordinate.
* `mood` - (Optional) Override the automatically determined mood with a custom mood.

## Attribute Reference

* `id` - A unique identifier for this data source.
* `condition` - The current weather condition (e.g., "clear", "cloudy", "rainy", "snowy").
* `temperature` - The current temperature in Celsius.
* `humidity` - The current humidity percentage.
* `wind_speed` - The current wind speed in meters per second.
* `mood` - A suggested mood based on the weather conditions.
* `suggested_moods` - A list of suggested moods appropriate for the current weather.
* `suggested_genres` - A list of suggested genres appropriate for the current weather.
* `location` - The location used for the weather data.
  * `city` - The city name.
  * `country` - The country code.
  * `lat` - The latitude coordinate.
  * `lon` - The longitude coordinate.