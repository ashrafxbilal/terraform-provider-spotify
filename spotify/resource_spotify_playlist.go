package spotify

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zmb3/spotify/v2"

	// "github.com/ashrafxbilal/terraform-provider-spotify/spotify/errors" // Uncomment when needed
	"github.com/ashrafxbilal/terraform-provider-spotify/spotify/logging"
	"github.com/ashrafxbilal/terraform-provider-spotify/spotify/utils"
)

func resourceSpotifyPlaylist() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSpotifyPlaylistCreate,
		ReadContext:   resourceSpotifyPlaylistRead,
		UpdateContext: resourceSpotifyPlaylistUpdate,
		DeleteContext: resourceSpotifyPlaylistDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the playlist",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the playlist",
			},
			"public": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the playlist is public",
			},
			"collaborative": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the playlist is collaborative",
			},
			"tracks": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The tracks in the playlist",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"snapshot_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Spotify snapshot ID for the playlist",
			},
			"spotify_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Spotify URL for the playlist",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp of the last update",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceSpotifyPlaylistCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*ProviderClient).SpotifyClient
	logger := logging.DefaultLogger.WithContext(ctx)

	// Log the operation
	logger.Info("Creating playlist", "name", d.Get("name").(string))

	// Get the current user's ID
	user, err := client.CurrentUser(ctx)
	if err != nil {
		return utils.HandleAPIError(ctx, err, "get current user for", "playlist creation", "")
	}

	// Create the playlist
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	public := d.Get("public").(bool)
	collaborative := d.Get("collaborative").(bool)

	// Validate required fields
	if name == "" {
		return utils.HandleValidationError(ctx, "Playlist name is required", map[string]string{
			"name": "empty",
		})
	}

	playlist, err := client.CreatePlaylistForUser(ctx, user.ID, name, description, public, collaborative)
	if err != nil {
		return utils.HandleAPIError(ctx, err, "create", "playlist", name)
	}

	// Add tracks if specified
	if v, ok := d.GetOk("tracks"); ok {
		tracks := expandTracks(v.([]interface{}))
		if len(tracks) > 0 {
			logger.Info("Adding tracks to playlist", "playlist_id", playlist.ID, "track_count", len(tracks))
			_, err := client.AddTracksToPlaylist(ctx, playlist.ID, tracks...)
			if err != nil {
				return utils.HandleAPIError(ctx, err, "add tracks to", "playlist", string(playlist.ID))
			}
		}
	}

	// Set the ID and other computed values
	playlistID := string(playlist.ID)
	d.SetId(playlistID)
	d.Set("snapshot_id", playlist.SnapshotID)
	d.Set("spotify_url", playlist.ExternalURLs["spotify"])
	d.Set("last_updated", time.Now().Format(time.RFC3339))

	// Log successful operation
	logger.Info("Successfully created playlist", "playlist_id", playlistID, "name", name)

	return resourceSpotifyPlaylistRead(ctx, d, m)
}

func resourceSpotifyPlaylistRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*ProviderClient).SpotifyClient
	logger := logging.DefaultLogger.WithContext(ctx)
	resourceID := d.Id()

	// Log the operation
	logger.Info("Reading playlist", "playlist_id", resourceID)

	// Get the playlist
	playlistID := spotify.ID(resourceID)
	playlist, err := client.GetPlaylist(ctx, playlistID)
	if err != nil {
		// Check if the playlist was not found
		if utils.IsSpotifyNotFoundError(err) {
			logger.Warn("Playlist not found, removing from state", "playlist_id", resourceID)
			d.SetId("")
			return diags
		}

		// Handle other API errors with standardized error handling
		return utils.HandleAPIError(ctx, err, "read", "playlist", resourceID)
	}

	// Set the values
	d.Set("name", playlist.Name)
	d.Set("description", playlist.Description)
	d.Set("public", playlist.IsPublic)
	d.Set("collaborative", playlist.Collaborative)
	d.Set("snapshot_id", playlist.SnapshotID)
	d.Set("spotify_url", playlist.ExternalURLs["spotify"])

	// Get the tracks
	tracks, err := getPlaylistTracks(ctx, client, playlistID)
	if err != nil {
		return utils.HandleAPIError(ctx, err, "read tracks for", "playlist", resourceID)
	}

	d.Set("tracks", tracks)

	// Log successful operation
	logger.Info("Successfully read playlist", "playlist_id", resourceID, "track_count", len(tracks))

	return diags
}

func resourceSpotifyPlaylistUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*ProviderClient).SpotifyClient
	playlistID := spotify.ID(d.Id())

	// Update basic details if changed
	if d.HasChanges("name", "description", "public", "collaborative") {
		name := d.Get("name").(string)
		description := d.Get("description").(string)
		public := d.Get("public").(bool)

		// Update playlist details
		err := client.ChangePlaylistName(ctx, playlistID, name)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating playlist name: %s", err))
		}

		err = client.ChangePlaylistDescription(ctx, playlistID, description)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating playlist description: %s", err))
		}

		err = client.ChangePlaylistAccess(ctx, playlistID, public)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating playlist access: %s", err))
		}

		// Collaborative status needs to be updated separately
		if d.HasChange("collaborative") {
			// Note: v2 API doesn't have SetPlaylistCollaborative, this would need to be handled differently
			// For now, we'll log a warning
			fmt.Printf("Warning: Setting collaborative status is not supported in the current Spotify API version")
		}
	}

	// Update tracks if changed
	if d.HasChange("tracks") {
		// Get current tracks
		currentTracks, err := getPlaylistTracks(ctx, client, playlistID)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error getting current playlist tracks: %s", err))
		}

		// Get desired tracks
		desiredTracksRaw := d.Get("tracks").([]interface{})
		desiredTracks := expandTracks(desiredTracksRaw)

		// Convert current tracks to Spotify IDs
		currentTrackIDs := make([]spotify.ID, len(currentTracks))
		for i, track := range currentTracks {
			currentTrackIDs[i] = spotify.ID(track)
		}

		// Replace all tracks
		err = client.ReplacePlaylistTracks(ctx, playlistID, desiredTracks...)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error replacing playlist tracks: %s", err))
		}
	}

	d.Set("last_updated", time.Now().Format(time.RFC3339))

	return resourceSpotifyPlaylistRead(ctx, d, m)
}

func resourceSpotifyPlaylistDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*ProviderClient).SpotifyClient

	// Spotify doesn't have a true delete operation for playlists
	// Instead, we need to unfollow the playlist
	playlistID := spotify.ID(d.Id())
	err := client.UnfollowPlaylist(ctx, playlistID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error unfollowing playlist: %s", err))
	}

	return diags
}

// Helper functions

func expandTracks(tracks []interface{}) []spotify.ID {
	result := make([]spotify.ID, len(tracks))
	for i, v := range tracks {
		result[i] = spotify.ID(v.(string))
	}
	return result
}

func getPlaylistTracks(ctx context.Context, client *spotify.Client, playlistID spotify.ID) ([]string, error) {
	var allTracks []string
	limit := 100
	offset := 0

	for {
		tracksPage, err := client.GetPlaylistItems(ctx, playlistID, spotify.Limit(limit), spotify.Offset(offset))
		if err != nil {
			return nil, err
		}

		for _, item := range tracksPage.Items {
			if item.Track.Track != nil {
				allTracks = append(allTracks, string(item.Track.Track.ID))
			}
		}

		if len(tracksPage.Items) < limit {
			break
		}

		offset += limit
	}

	return allTracks, nil
}