package testing

import (
	"context"

	"github.com/zmb3/spotify/v2"
)

// MockSpotifyClient provides a minimal mock implementation of the Spotify client for testing
type MockSpotifyClient struct {
	// Mock responses
	User     *spotify.PrivateUser
	Playlist *spotify.FullPlaylist

	// Error responses
	UserError     error
	PlaylistError error

	// Function mocks for more complex behaviors
	FeaturedPlaylistsFunc func(ctx context.Context, opts ...spotify.RequestOption) (string, *spotify.SimplePlaylistPage, error)
	GetTracksFunc         func(ctx context.Context, ids ...spotify.ID) ([]*spotify.FullTrack, error)
}

// CurrentUser mocks the CurrentUser method
func (m *MockSpotifyClient) CurrentUser(ctx context.Context) (*spotify.PrivateUser, error) {
	if m.UserError != nil {
		return nil, m.UserError
	}

	if m.User == nil {
		// Return a default mock user if none is set
		return &spotify.PrivateUser{
			User: spotify.User{
				ID:          "mock-user-id",
				DisplayName: "Mock User",
			},
			Email: "mock@example.com",
		}, nil
	}

	return m.User, nil
}

// GetPlaylist mocks the GetPlaylist method
func (m *MockSpotifyClient) GetPlaylist(ctx context.Context, id spotify.ID) (*spotify.FullPlaylist, error) {
	if m.PlaylistError != nil {
		return nil, m.PlaylistError
	}

	if m.Playlist == nil {
		// Return a default mock playlist if none is set
		return &spotify.FullPlaylist{
			SimplePlaylist: spotify.SimplePlaylist{
				ID:   id,
				Name: "Mock Playlist",
			},
		}, nil
	}

	return m.Playlist, nil
}

// FeaturedPlaylists mocks the FeaturedPlaylists method
func (m *MockSpotifyClient) FeaturedPlaylists(ctx context.Context, opts ...spotify.RequestOption) (string, *spotify.SimplePlaylistPage, error) {
	if m.FeaturedPlaylistsFunc != nil {
		return m.FeaturedPlaylistsFunc(ctx, opts...)
	}

	// Default implementation
	return "Featured Playlists", &spotify.SimplePlaylistPage{
		Playlists: []spotify.SimplePlaylist{
			{
				ID:   "default-playlist-id",
				Name: "Default Featured Playlist",
			},
		},
	}, nil
}

// GetTracks mocks the GetTracks method
func (m *MockSpotifyClient) GetTracks(ctx context.Context, ids ...spotify.ID) ([]*spotify.FullTrack, error) {
	if m.GetTracksFunc != nil {
		return m.GetTracksFunc(ctx, ids...)
	}

	// Default implementation
	tracks := make([]*spotify.FullTrack, len(ids))
	for i, id := range ids {
		tracks[i] = &spotify.FullTrack{
			SimpleTrack: spotify.SimpleTrack{
				ID:   id,
				Name: "Track " + string(id),
			},
		}
	}

	return tracks, nil
}