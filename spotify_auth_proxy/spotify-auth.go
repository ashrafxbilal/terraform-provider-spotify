// Package main provides a command-line utility for obtaining Spotify API tokens
// for use with the Terraform Spotify provider.
//
// Usage:
//   1. Set required environment variables:
//      - SPOTIFY_CLIENT_ID: Your Spotify application client ID
//      - SPOTIFY_CLIENT_SECRET: Your Spotify application client secret
//      - SPOTIFY_REDIRECT_URI: Your registered redirect URI
//      - SPOTIFY_SCOPES: (Optional) Space-separated list of required scopes
//   2. Run the program and follow the prompts
//   3. Use the obtained refresh token in your Terraform configuration
package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Configuration constants
var (
	// These values should be provided via environment variables
	clientID     = os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")
	redirectURI  = os.Getenv("SPOTIFY_REDIRECT_URI")
	scopes       = os.Getenv("SPOTIFY_SCOPES")
)

// validateEnvironment checks if all required environment variables are set
func validateEnvironment() error {
	missing := []string{}

	if clientID == "" {
		missing = append(missing, "SPOTIFY_CLIENT_ID")
	}
	if clientSecret == "" {
		missing = append(missing, "SPOTIFY_CLIENT_SECRET")
	}
	if redirectURI == "" {
		missing = append(missing, "SPOTIFY_REDIRECT_URI")
	}
	if scopes == "" {
		// Set default scopes if not provided
		scopes = "ugc-image-upload user-top-read user-read-recently-played user-read-private playlist-modify-public playlist-modify-private"
		fmt.Println("Warning: SPOTIFY_SCOPES not set, using default scopes")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	return nil
}

func main() {
	// Validate environment variables
	if err := validateEnvironment(); err != nil {
		log.Fatalf("Environment validation failed: %v\n\nPlease set the required environment variables:\n\nexport SPOTIFY_CLIENT_ID=your_client_id\nexport SPOTIFY_CLIENT_SECRET=your_client_secret\nexport SPOTIFY_REDIRECT_URI=your_redirect_uri\nexport SPOTIFY_SCOPES=your_scopes\n", err)
	}

	authURL := fmt.Sprintf(
		"https://accounts.spotify.com/authorize?client_id=%s&response_type=code&redirect_uri=%s&scope=%s",
		url.QueryEscape(clientID),
		url.QueryEscape(redirectURI),
		url.QueryEscape(scopes),
	)

	fmt.Println("1. Open the following URL in your browser and authorize the app:")
	fmt.Println(authURL)
	fmt.Println("\n2. After authorization, you'll be redirected to the Glitch page.")
	fmt.Println("   Copy the 'code' value from the URL in the address bar and paste it below.")
	fmt.Print("\nPaste code here: ")

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			log.Fatalf("Error reading input: %v", err)
		}
		log.Fatalf("No input provided")
	}
	code := scanner.Text()

	if code == "" {
		log.Fatalf("Authorization code cannot be empty")
	}

	tokens, err := exchangeCodeForToken(code)
	if err != nil {
		log.Fatalf("Failed to exchange code for token: %v", err)
	}

	fmt.Println("\nAccess Token:", tokens.AccessToken)
	fmt.Println("Refresh Token:", tokens.RefreshToken)
	fmt.Println("\nTo use these tokens with Terraform, set the following environment variables:")
	fmt.Println("export TF_VAR_spotify_client_id=\"" + clientID + "\"")
	fmt.Println("export TF_VAR_spotify_client_secret=\"" + clientSecret + "\"")
	fmt.Println("export TF_VAR_spotify_redirect_uri=\"" + redirectURI + "\"")
	fmt.Println("export TF_VAR_spotify_refresh_token=\"" + tokens.RefreshToken + "\"")
}

type SpotifyTokenResponse struct {
    AccessToken  string `json:"access_token"`
    TokenType    string `json:"token_type"`
    Scope        string `json:"scope"`
    ExpiresIn    int    `json:"expires_in"`
    RefreshToken string `json:"refresh_token"`
}

// ErrorResponse represents an error response from the Spotify API
type ErrorResponse struct {
	Error       string `json:"error"`
	Description string `json:"error_description"`
}

func exchangeCodeForToken(code string) (*SpotifyTokenResponse, error) {
    data := url.Values{}
    data.Set("grant_type", "authorization_code")
    data.Set("code", code)
    data.Set("redirect_uri", redirectURI)

    req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    authHeader := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
    req.Header.Set("Authorization", "Basic "+authHeader)
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()

    // Read the response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response body: %w", err)
    }

    // Check for error response
    if resp.StatusCode != http.StatusOK {
        var errorResp ErrorResponse
        if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Error != "" {
            return nil, fmt.Errorf("API error: %s - %s", errorResp.Error, errorResp.Description)
        }
        return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(body))
    }

    // Parse the successful response
    var tokenData SpotifyTokenResponse
    if err := json.Unmarshal(body, &tokenData); err != nil {
        return nil, fmt.Errorf("failed to parse response: %w", err)
    }

    // Validate the response
    if tokenData.AccessToken == "" {
        return nil, fmt.Errorf("received empty access token")
    }
    if tokenData.RefreshToken == "" {
        return nil, fmt.Errorf("received empty refresh token")
    }

    return &tokenData, nil
}
