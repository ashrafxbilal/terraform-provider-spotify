package spotify

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zmb3/spotify/v2"
)

func dataSourceFeaturedPlaylists() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFeaturedPlaylistsRead,
		Schema: map[string]*schema.Schema{
			"country": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An ISO 3166-1 alpha-2 country code to get featured playlists for a specific country",
			},
			"locale": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The desired language, consisting of an ISO 639-1 language code and an ISO 3166-1 alpha-2 country code, joined by an underscore",
			},
			"timestamp": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A timestamp in ISO 8601 format to get featured playlists for a specific date and time",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     20,
				Description: "The maximum number of playlists to return",
			},
			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The message from Spotify for the featured playlists",
			},
			"playlists": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of featured playlists",
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

func dataSourceFeaturedPlaylistsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*ProviderClient).SpotifyClient

	// Get parameters from schema
	country := d.Get("country").(string)
	locale := d.Get("locale").(string)
	timestampStr := d.Get("timestamp").(string)
	limit := d.Get("limit").(int)

	// Build options for the API call
	opts := []spotify.RequestOption{}

	// Add country if specified
	if country != "" {
		opts = append(opts, spotify.Country(country))
	}

	// Add locale if specified
	if locale != "" {
		opts = append(opts, spotify.Locale(locale))
	}

	// Add timestamp if specified
	if timestampStr != "" {
		timestamp, err := time.Parse(time.RFC3339, timestampStr)
		if err != nil {
			return diag.FromErr(fmt.Errorf("invalid timestamp format: %s. Use ISO 8601 format (e.g., 2023-01-01T12:00:00Z)", err))
		}
		// Convert timestamp to string in the format expected by the API
		timestampFormatted := timestamp.Format(time.RFC3339)
		opts = append(opts, spotify.Timestamp(timestampFormatted))
	}

	// Add limit
	if limit > 0 {
		opts = append(opts, spotify.Limit(limit))
	}

	// Get featured playlists
	message, playlistPage, err := client.FeaturedPlaylists(ctx, opts...)
	if err != nil {
		// Log the error but don't fail - allow empty results
		fmt.Printf("Warning: error getting featured playlists: %s\n", err)
		// Set an empty message and continue with empty playlists
		message = "No featured playlists available"
		// Return empty results instead of error
		d.Set("message", message)
		d.Set("playlists", []map[string]interface{}{})
		// Set the ID with error indicator
		d.SetId(fmt.Sprintf("featured-playlists-error-%d", time.Now().Unix()))
		return diags
	}

	// Set the message
	d.Set("message", message)

	// Process playlists
	playlists := make([]map[string]interface{}, 0, len(playlistPage.Playlists))
	for _, playlist := range playlistPage.Playlists {
		playlistMap := make(map[string]interface{})
		playlistMap["id"] = string(playlist.ID)
		playlistMap["name"] = playlist.Name
		playlistMap["description"] = playlist.Description
		playlistMap["owner"] = playlist.Owner.DisplayName
		playlistMap["image_url"] = ""
		if len(playlist.Images) > 0 {
			playlistMap["image_url"] = playlist.Images[0].URL
		}
		playlistMap["tracks_total"] = fmt.Sprintf("%d", playlist.Tracks.Total)
		playlistMap["url"] = playlist.ExternalURLs["spotify"]

		playlists = append(playlists, playlistMap)
	}

	// Set the playlists
	d.Set("playlists", playlists)

	// Set the ID
	d.SetId(fmt.Sprintf("featured-playlists-%d", time.Now().Unix()))

	return diags
}