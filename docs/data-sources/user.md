---
page_title: "spotify_user Data Source - terraform-provider-spotify"
subcategory: ""
description: |-
  Retrieves information about the current Spotify user.
---

# Data Source: spotify_user

Retrieves information about the current Spotify user, including profile details, subscription level, and follower count.

## Example Usage

```terraform
data "spotify_user" "me" {}

output "user_name" {
  value = data.spotify_user.me.display_name
}

output "user_profile_image" {
  value = length(data.spotify_user.me.images) > 0 ? data.spotify_user.me.images[0] : "No profile image"
}

resource "spotify_playlist" "personal" {
  name        = "${data.spotify_user.me.display_name}'s Terraform Playlist"
  description = "A personalized playlist for ${data.spotify_user.me.display_name}"
  public      = true
}
```

## Argument Reference

This data source has no arguments.

## Attribute Reference

* `id` - The Spotify user ID.
* `display_name` - The user's display name.
* `email` - The user's email address.
* `product` - The user's Spotify subscription level (e.g., "premium", "free").
* `followers` - The total number of followers.
* `images` - A list of the user's profile images.