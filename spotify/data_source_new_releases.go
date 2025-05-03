package spotify

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zmb3/spotify/v2"
)

func dataSourceNewReleases() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewReleasesRead,
		Schema: map[string]*schema.Schema{
			"country": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An ISO 3166-1 alpha-2 country code to get new releases for a specific country",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     20,
				Description: "The maximum number of new releases to return",
			},
			"albums": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of new release albums",
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
			"album_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of album IDs",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceNewReleasesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*ProviderClient).SpotifyClient

	// Get parameters from schema
	country := d.Get("country").(string)
	limit := d.Get("limit").(int)

	// Build options for the API call
	opts := []spotify.RequestOption{}

	// Add country if specified
	if country != "" {
		opts = append(opts, spotify.Country(country))
	}

	// Add limit
	if limit > 0 {
		opts = append(opts, spotify.Limit(limit))
	}

	// Get new releases
	albumPage, err := client.NewReleases(ctx, opts...)
	if err != nil {
		// Log the error but don't fail - allow empty results
		fmt.Printf("Warning: error getting new releases: %s\n", err)
		// Return empty results instead of error
		d.Set("albums", []map[string]interface{}{})
		d.Set("album_ids", []string{})
		// Set the ID with error indicator
		d.SetId(fmt.Sprintf("new-releases-error-%d", time.Now().Unix()))
		return diags
	}

	// Process albums
	albums := make([]map[string]interface{}, 0, len(albumPage.Albums))
	albumIDs := make([]string, 0, len(albumPage.Albums))

	for _, album := range albumPage.Albums {
		albumMap := make(map[string]interface{})
		albumMap["id"] = string(album.ID)
		albumMap["name"] = album.Name
		albumMap["type"] = album.AlbumType
		
		// Get artist names
		artistNames := ""
		for i, artist := range album.Artists {
			if i > 0 {
				artistNames += ", "
			}
			artistNames += artist.Name
		}
		albumMap["artists"] = artistNames
		
		// Get release date
		albumMap["release_date"] = album.ReleaseDate
		albumMap["release_date_precision"] = album.ReleaseDatePrecision
		
		// Get image URL
		albumMap["image_url"] = ""
		if len(album.Images) > 0 {
			albumMap["image_url"] = album.Images[0].URL
		}
		
		// Get total tracks
		albumMap["total_tracks"] = fmt.Sprintf("%d", album.TotalTracks)
		
		// Get external URL
		albumMap["url"] = album.ExternalURLs["spotify"]

		albums = append(albums, albumMap)
		albumIDs = append(albumIDs, string(album.ID))
	}

	// Set the albums and album IDs
	d.Set("albums", albums)
	d.Set("album_ids", albumIDs)

	// Set the ID
	d.SetId(fmt.Sprintf("new-releases-%d", time.Now().Unix()))

	return diags
}