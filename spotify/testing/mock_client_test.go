package testing

import (
	"context"
	"errors"
	"testing"

	"github.com/zmb3/spotify/v2"
)

func TestMockSpotifyClient_CurrentUser(t *testing.T) {
	ctx := context.Background()
	
	// Test with default user
	mockClient := &MockSpotifyClient{}
	user, err := mockClient.CurrentUser(ctx)
	
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if user.ID != "mock-user-id" {
		t.Errorf("Expected ID 'mock-user-id', got %v", user.ID)
	}
	
	// Test with custom user
	customUser := &spotify.PrivateUser{
		User: spotify.User{
			ID: "custom-id",
		},
	}
	
	mockClient = &MockSpotifyClient{User: customUser}
	user, err = mockClient.CurrentUser(ctx)
	
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if user.ID != "custom-id" {
		t.Errorf("Expected ID 'custom-id', got %v", user.ID)
	}
	
	// Test with error
	testErr := errors.New("test error")
	mockClient = &MockSpotifyClient{UserError: testErr}
	_, err = mockClient.CurrentUser(ctx)
	
	if err != testErr {
		t.Errorf("Expected error %v, got %v", testErr, err)
	}
}

func TestMockSpotifyClient_GetPlaylist(t *testing.T) {
	ctx := context.Background()
	playlistID := spotify.ID("test-id")
	
	// Test with default playlist
	mockClient := &MockSpotifyClient{}
	playlist, err := mockClient.GetPlaylist(ctx, playlistID)
	
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if playlist.ID != playlistID {
		t.Errorf("Expected ID %v, got %v", playlistID, playlist.ID)
	}
	
	if playlist.Name != "Mock Playlist" {
		t.Errorf("Expected name 'Mock Playlist', got %v", playlist.Name)
	}
	
	// Test with custom playlist
	customPlaylist := &spotify.FullPlaylist{
		SimplePlaylist: spotify.SimplePlaylist{
			ID:   "custom-id",
			Name: "Custom Playlist",
		},
	}
	
	mockClient = &MockSpotifyClient{Playlist: customPlaylist}
	playlist, err = mockClient.GetPlaylist(ctx, playlistID)
	
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if playlist.ID != "custom-id" {
		t.Errorf("Expected ID 'custom-id', got %v", playlist.ID)
	}
	
	if playlist.Name != "Custom Playlist" {
		t.Errorf("Expected name 'Custom Playlist', got %v", playlist.Name)
	}
	
	// Test with error
	testErr := errors.New("test error")
	mockClient = &MockSpotifyClient{PlaylistError: testErr}
	_, err = mockClient.GetPlaylist(ctx, playlistID)
	
	if err != testErr {
		t.Errorf("Expected error %v, got %v", testErr, err)
	}
}