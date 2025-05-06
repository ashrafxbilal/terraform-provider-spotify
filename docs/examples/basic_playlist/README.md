# Spotify Personalized Playlist Generator

This example demonstrates how to use the Terraform Spotify Provider to create a personalized playlist based on:

- Your listening history (short, medium, and long-term)
- Current weather conditions
- Time of day and day of week
- Global popular tracks
- Latest releases from your favorite artists/genres

## Prerequisites

- [Terraform](https://www.terraform.io/downloads.html) 0.13+
- A Spotify account with a registered application
- Spotify API credentials (Client ID, Client Secret, Redirect URI)
- A Spotify refresh token (obtained using the included auth proxy)

## Setup

1. **Set up your Spotify API credentials**

   First, you need to register an application in the [Spotify Developer Dashboard](https://developer.spotify.com/dashboard/applications).
   
   When creating your app, add `https://glitch.com/~spotify-oauth-redirect` as a redirect URI.

2. **Get a refresh token**

   Use the included auth proxy to get a refresh token:

   ```bash
   # Set environment variables
   export SPOTIFY_CLIENT_ID="your-client-id"
   export SPOTIFY_CLIENT_SECRET="your-client-secret"
   export SPOTIFY_REDIRECT_URI="https://glitch.com/~spotify-oauth-redirect"
   
   # Run the auth proxy
   cd ../../spotify_auth_proxy
   go run spotify-auth.go
   ```

   Follow the prompts to authorize the application and get a refresh token.

3. **Configure your variables**

   Copy the example variables file and fill in your values:

   ```bash
   cp terraform.tfvars.example terraform.tfvars
   ```

   Edit `terraform.tfvars` with your Spotify API credentials and refresh token.

## Usage

1. **Initialize Terraform**

   ```bash
   terraform init
   ```

2. **Plan the changes**

   ```bash
   terraform plan
   ```

3. **Apply the changes to create your playlist**

   ```bash
   terraform apply
   ```

4. **Check your Spotify account**

   A new playlist should appear in your Spotify account with a name based on the current day of the week, temperature, and a Gen-Z style description.

## Customization

You can customize the playlist generation by modifying the following:

- **Weather-based mood**: Uncomment and set your preferred mood in the `spotify_weather` data source
- **Time-based mood**: Uncomment and set your preferred mood in the `spotify_time` data source
- **Genre preferences**: Modify the genre selections in the various `spotify_tracks` data sources

## Structure

This example follows Terraform best practices with separate files for:

- `main.tf`: Main configuration and resources
- `variables.tf`: Input variable declarations
- `outputs.tf`: Output declarations
- `terraform.tfvars`: Variable values (not committed to version control)

## Notes

- The playlist is regenerated each time you run `terraform apply`
- The playlist name and description are dynamically generated based on current conditions
- Track selection is influenced by your Spotify listening history