package spotify

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The client ID for Spotify API authentication",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The client secret for Spotify API authentication",
			},
			"redirect_uri": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The redirect URI for Spotify API authentication",
			},
			"refresh_token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The refresh token for Spotify API",
			},
			"weather_api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "API key for OpenWeatherMap",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"spotify_playlist":       resourceSpotifyPlaylist(),
			"spotify_playlist_track": resourceSpotifyPlaylistTrack(),
			"spotify_playlist_cover": resourceSpotifyPlaylistCover(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"spotify_tracks":             dataSourceSpotifyTracks(),
			"spotify_weather":            dataSourceWeather(),
			"spotify_time":               dataSourceTime(),
			"spotify_user":               dataSourceUser(),
			"spotify_user_preferences":   dataSourceUserPreferences(),
			"spotify_featured_playlists": dataSourceFeaturedPlaylists(),
			"spotify_new_releases":       dataSourceNewReleases(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

// ProviderClient holds the Spotify client and other API clients
type ProviderClient struct {
	SpotifyClient *spotify.Client
	WeatherAPIKey string
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	clientID := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)
	redirectURI := d.Get("redirect_uri").(string)
	refreshToken := d.Get("refresh_token").(string)
	weatherAPIKey := d.Get("weather_api_key").(string)

	oauthConfig := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURI,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.spotify.com/authorize",
			TokenURL: "https://accounts.spotify.com/api/token",
		},
		Scopes: []string{
			"user-read-private",
			"playlist-modify-public",
			"playlist-modify-private",
			"user-top-read",
			"user-read-recently-played",
			"ugc-image-upload",
		},
	}

	// Use refresh token to obtain a new access token
	token := &oauth2.Token{RefreshToken: refreshToken}
	tokenSource := oauthConfig.TokenSource(ctx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("error refreshing access token: %s", err))
	}

	// Print token information for debugging
	fmt.Printf("Access Token: %s...\n", newToken.AccessToken[:10])
	fmt.Printf("Token Type: %s\n", newToken.TokenType)
	fmt.Printf("Token Expiry: %s\n", newToken.Expiry.String())

	// Set credentials as environment variables for other resources to use
	// This is a workaround for resources that need to make direct API calls
	if err := os.Setenv("SPOTIFY_ACCESS_TOKEN", newToken.AccessToken); err != nil {
		return nil, diag.FromErr(fmt.Errorf("error setting SPOTIFY_ACCESS_TOKEN environment variable: %s", err))
	}
	if err := os.Setenv("SPOTIFY_REFRESH_TOKEN", refreshToken); err != nil {
		return nil, diag.FromErr(fmt.Errorf("error setting SPOTIFY_REFRESH_TOKEN environment variable: %s", err))
	}
	if err := os.Setenv("SPOTIFY_CLIENT_ID", clientID); err != nil {
		return nil, diag.FromErr(fmt.Errorf("error setting SPOTIFY_CLIENT_ID environment variable: %s", err))
	}
	if err := os.Setenv("SPOTIFY_CLIENT_SECRET", clientSecret); err != nil {
		return nil, diag.FromErr(fmt.Errorf("error setting SPOTIFY_CLIENT_SECRET environment variable: %s", err))
	}
	if err := os.Setenv("SPOTIFY_REDIRECT_URI", redirectURI); err != nil {
		return nil, diag.FromErr(fmt.Errorf("error setting SPOTIFY_REDIRECT_URI environment variable: %s", err))
	}

	httpClient := oauthConfig.Client(ctx, newToken)
	spotifyClient := spotify.New(httpClient)

	// Test with explicit API call format
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/recommendations?seed_genres=pop&limit=1", nil)
	if err != nil {
		fmt.Printf("Error creating request: %s\n", err)
	} else {
		req.Header.Set("Authorization", "Bearer "+newToken.AccessToken)
		resp, httpErr := httpClient.Do(req)
		if httpErr != nil {
			fmt.Printf("Error making direct API call: %s\n", httpErr)
		} else {
			fmt.Printf("Direct API call status: %d\n", resp.StatusCode)
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Error reading response body: %s\n", err)
			} else {
				fmt.Printf("Response body: %s\n", string(body))
			}
			if err := resp.Body.Close(); err != nil {
				fmt.Printf("Error closing response body: %s\n", err)
			}
		}
	}

	// Test API connection
	user, err := spotifyClient.CurrentUser(ctx)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("error authenticating with Spotify: %s", err))
	}
	fmt.Printf("Logged in as: %s\n", user.ID)

	// Test recommendations API with a simple query
	fmt.Println("Testing recommendations API...")
	seeds := spotify.Seeds{Genres: []string{"pop"}}
	attrs := spotify.NewTrackAttributes()
	recs, err := spotifyClient.GetRecommendations(ctx, seeds, attrs, spotify.Limit(1))
	if err != nil {
		fmt.Printf("Warning: Recommendations API test failed: %s\n", err)
	} else {
		fmt.Printf("Recommendations API test successful, got %d tracks\n", len(recs.Tracks))
	}

	return &ProviderClient{
		SpotifyClient: spotifyClient,
		WeatherAPIKey: weatherAPIKey,
	}, diags
}
