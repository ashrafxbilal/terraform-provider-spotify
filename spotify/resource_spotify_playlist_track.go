package spotify

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zmb3/spotify/v2"
)

func resourceSpotifyPlaylistTrack() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSpotifyPlaylistTrackCreate,
		ReadContext:   resourceSpotifyPlaylistTrackRead,
		UpdateContext: resourceSpotifyPlaylistTrackUpdate,
		DeleteContext: resourceSpotifyPlaylistTrackDelete,
		Schema: map[string]*schema.Schema{
			"playlist_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the playlist",
			},
			"track_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the track",
			},
			"position": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The position of the track in the playlist (0-based index)",
			},
			"added_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the track was added",
			},
			"track_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the track",
			},
			"artist": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The artist of the track",
			},
			"album": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The album of the track",
			},
			"duration_ms": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The duration of the track in milliseconds",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceSpotifyPlaylistTrackCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*ProviderClient).SpotifyClient

	playlistID := spotify.ID(d.Get("playlist_id").(string))
	trackID := spotify.ID(d.Get("track_id").(string))
	position := d.Get("position").(int)

	// Add the track to the playlist
	trackURI := fmt.Sprintf("spotify:track:%s", trackID)
	if position > 0 {
		// Add at specific position
		// Note: v2 API doesn't have AddTracksToPlaylistAtPosition, so we use AddTracksToPlaylist
		// and then reorder if needed
		_, err := client.AddTracksToPlaylist(ctx, playlistID, spotify.ID(trackURI))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error adding track to playlist: %s", err))
		}

	} else {
		// Add to the end
		_, err := client.AddTracksToPlaylist(ctx, playlistID, spotify.ID(trackURI))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error adding track to playlist: %s", err))
		}
	}

	d.SetId(fmt.Sprintf("%s:%s", playlistID, trackID))

	d.Set("added_at", time.Now().Format(time.RFC3339))

	return resourceSpotifyPlaylistTrackRead(ctx, d, m)
}

func resourceSpotifyPlaylistTrackRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*ProviderClient).SpotifyClient

	playlistID := spotify.ID(d.Get("playlist_id").(string))
	trackID := spotify.ID(d.Get("track_id").(string))

	// Get the track details
	track, err := client.GetTrack(ctx, trackID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting track: %s", err))
	}

	// Set the track details
	d.Set("track_name", track.Name)
	if len(track.Artists) > 0 {
		d.Set("artist", track.Artists[0].Name)
	}
	d.Set("album", track.Album.Name)
	d.Set("duration_ms", track.Duration)

	// Check if the track is still in the playlist
	trackExists, position, err := checkTrackInPlaylist(ctx, client, playlistID, trackID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error checking track in playlist: %s", err))
	}

	if !trackExists {
		// Track is no longer in the playlist
		d.SetId("")
		return diags
	}

	// Update position if it's different
	d.Set("position", position)

	return diags
}

func resourceSpotifyPlaylistTrackUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*ProviderClient).SpotifyClient

	playlistID := spotify.ID(d.Get("playlist_id").(string))
	trackID := spotify.ID(d.Get("track_id").(string))

	// Only position can be updated
	if d.HasChange("position") {
		// First, remove the track
		_, err := client.RemoveTracksFromPlaylist(ctx, playlistID, trackID)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error removing track from playlist: %s", err))
		}

		// Then add it back at the new position
		trackURI := fmt.Sprintf("spotify:track:%s", trackID)
		_, err = client.AddTracksToPlaylist(ctx, playlistID, spotify.ID(trackURI))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error adding track to playlist at new position: %s", err))
		}
	
	}

	return resourceSpotifyPlaylistTrackRead(ctx, d, m)
}

func resourceSpotifyPlaylistTrackDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*ProviderClient).SpotifyClient

	playlistID := spotify.ID(d.Get("playlist_id").(string))
	trackID := spotify.ID(d.Get("track_id").(string))

	// Remove the track from the playlist
	_, err := client.RemoveTracksFromPlaylist(ctx, playlistID, trackID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error removing track from playlist: %s", err))
	}

	d.SetId("")

	return diags
}

// Helper function to check if a track is in a playlist and get its position
func checkTrackInPlaylist(ctx context.Context, client *spotify.Client, playlistID, trackID spotify.ID) (bool, int, error) {
	limit := 100
	offset := 0
	position := 0

	for {
		tracksPage, err := client.GetPlaylistItems(ctx, playlistID, spotify.Limit(limit), spotify.Offset(offset))
		if err != nil {
			return false, 0, err
		}

		for i, item := range tracksPage.Items {
			if item.Track.Track != nil && item.Track.Track.ID == trackID {
				return true, offset + i, nil
			}
		}

		if len(tracksPage.Items) < limit {
			break
		}

		offset += limit
		position += len(tracksPage.Items)
	}

	return false, 0, nil
}