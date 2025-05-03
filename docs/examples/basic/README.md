# Basic Spotify Provider Example

This example demonstrates how to use the Spotify Terraform Provider to create a simple playlist with a custom cover image.

## Features

This example:

1. Creates a playlist with predefined tracks
2. Adds a custom cover image with an emoji
3. Outputs information about the created playlist and user

## Prerequisites

- [Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
- A Spotify account with a registered application

## Setup

1. Create a `terraform.tfvars` file with your credentials:

```hcl
spotify_client_id     = "your-client-id"
spotify_client_secret = "your-client-secret"
spotify_redirect_uri  = "your-redirect-uri"
spotify_refresh_token = "your-refresh-token"
```

2. Initialize Terraform:

```sh
terraform init
```

3. Apply the configuration:

```sh
terraform apply
```

## What This Example Creates

### Data Sources

- `spotify_user` - Gets information about the current user

### Resources

- `spotify_playlist.basic_playlist` - A simple playlist with predefined tracks
- `spotify_playlist_cover.basic_cover` - A custom cover for the playlist with an emoji

## Outputs

- `playlist_info` - Information about the created playlist
- `user_info` - Information about the current user

## Notes

- The playlist will be created in your Spotify account and will be visible in your Spotify client
- You can modify the track list in the `main.tf` file to include your preferred tracks
- The emoji and background color for the cover image can be customized in the `main.tf` file