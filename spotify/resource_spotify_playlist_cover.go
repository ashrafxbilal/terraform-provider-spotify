package spotify

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ashrafxbilal/terraform-provider-spotify/spotify/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
)

func resourceSpotifyPlaylistCover() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSpotifyPlaylistCoverCreate,
		ReadContext:   resourceSpotifyPlaylistCoverRead,
		UpdateContext: resourceSpotifyPlaylistCoverUpdate,
		DeleteContext: resourceSpotifyPlaylistCoverDelete,
		Schema: map[string]*schema.Schema{
			"playlist_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the playlist",
			},
			"image_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "URL of the image to use as playlist cover",
			},
			"emoji": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Emoji to use for generating a cover image",
			},
			"mood": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Mood to determine emoji for the cover image",
			},
			"weather": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Weather condition to determine emoji for the cover image",
			},
			"background_color": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "#1DB954", // Spotify green
				Description: "Background color for the generated cover image (hex code)",
			},
			"force_update": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Force update the playlist cover image even if no changes are detected",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp of the last update",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceSpotifyPlaylistCoverCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*ProviderClient).SpotifyClient
	playlistID := spotify.ID(d.Get("playlist_id").(string))

	// Generate a unique ID for this resource
	d.SetId(fmt.Sprintf("%s-cover-%d", playlistID, time.Now().Unix()))

	// Set the cover image
	diags := setPlaylistCoverImage(ctx, d, client)
	if diags.HasError() {
		return diags
	}

	// Use the safe version to handle potential errors
	if !utils.SetResourceDataSafe(d, "last_updated", time.Now().Format(time.RFC3339), ctx) {
		return diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error setting last_updated",
			Detail:   "Could not set last_updated timestamp",
		}}
	}

	return resourceSpotifyPlaylistCoverRead(ctx, d, m)
}

func resourceSpotifyPlaylistCoverRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	// Nothing to read as Spotify doesn't provide an API to get the current cover image data
	// We can only see the image URL in the playlist object, but not compare with our source
	return diags
}

func resourceSpotifyPlaylistCoverUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*ProviderClient).SpotifyClient

	// Check if any of the fields that affect the image have changed or if force_update is true
	forceUpdate := d.Get("force_update").(bool)
	imageChanged := d.HasChange("image_url") || d.HasChange("emoji") || d.HasChange("mood") ||
		d.HasChange("weather") || d.HasChange("background_color")

	// Only update if there are changes or force_update is true
	if imageChanged || forceUpdate {
		// Set the cover image
		diags := setPlaylistCoverImage(ctx, d, client)
		if diags.HasError() {
			return diags
		}

		// Use the safe version to handle potential errors
		if !utils.SetResourceDataSafe(d, "last_updated", time.Now().Format(time.RFC3339), ctx) {
			return diag.Diagnostics{diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error setting last_updated",
				Detail:   "Could not set last_updated timestamp",
			}}
		}
	}

	return resourceSpotifyPlaylistCoverRead(ctx, d, m)
}

func resourceSpotifyPlaylistCoverDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	// Spotify doesn't provide an API to delete a playlist cover image
	// The cover will remain until replaced
	d.SetId("")
	return diags
}

func setPlaylistCoverImage(ctx context.Context, d *schema.ResourceData, client *spotify.Client) diag.Diagnostics {
	var diags diag.Diagnostics
	playlistID := spotify.ID(d.Get("playlist_id").(string))

	// Get image data as base64
	imageData, err := getImageData(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error preparing image data: %s", err))
	}

	// The Spotify API endpoint for setting a playlist cover image
	// is not directly exposed in the zmb3/spotify library, so we need to make a direct API call
	// URL: https://api.spotify.com/v1/playlists/{playlist_id}/images

	// Create the request URL
	url := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/images", playlistID)

	// Create a new HTTP request
	// Important: Spotify expects the raw base64 string without any prefixes
	req, err := http.NewRequestWithContext(ctx, "PUT", url, strings.NewReader(imageData))
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating request: %s", err))
	}

	// Set the content type header - this is important for Spotify's API
	req.Header.Set("Content-Type", "image/jpeg")

	// Get the access token
	token := getAccessToken(ctx)
	if token == "" {
		return diag.FromErr(fmt.Errorf("failed to get access token"))
	}

	// Add the Authorization header to the request
	req.Header.Set("Authorization", "Bearer "+token)

	// Create a new HTTP client with reasonable timeouts
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Make the request
	resp, err := httpClient.Do(req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error uploading image: %s", err))
	}
	// Use our utility function to safely close the response body
	defer utils.HandleResponseBodyClose(ctx, resp)

	// Check the response status
	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf("error from Spotify API: %s (status code: %d)", string(body), resp.StatusCode))
	}

	return diags
}

func getImageData(d *schema.ResourceData) (string, error) {
	// Check if image_url is provided
	if imageURL, ok := d.GetOk("image_url"); ok {
		// Download the image from the URL
		resp, err := http.Get(imageURL.(string))
		if err != nil {
			return "", fmt.Errorf("error downloading image: %w", err)
		}
		// Use a separate function to safely close the response body
		ctx := context.Background() // Use background context since we don't have one passed in
		defer utils.HandleResponseBodyClose(ctx, resp)

		// Read the image data
		imageBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("error reading image data: %w", err)
		}

		// Encode the image as base64
		return base64.StdEncoding.EncodeToString(imageBytes), nil
	}

	// If emoji is provided, generate an image with the emoji
	if emoji, ok := d.GetOk("emoji"); ok {
		return generatePatternCoverImage(emoji.(string), d.Get("background_color").(string))
	}

	// If mood is provided, select an appropriate emoji
	if mood, ok := d.GetOk("mood"); ok {
		emoji := getMoodEmoji(mood.(string))
		return generatePatternCoverImage(emoji, d.Get("background_color").(string))
	}

	// If weather is provided, select an appropriate emoji
	if weather, ok := d.GetOk("weather"); ok {
		emoji := getWeatherEmoji(weather.(string))
		return generatePatternCoverImage(emoji, d.Get("background_color").(string))
	}

	// Default emoji if nothing else is provided
	return generatePatternCoverImage("🎵", d.Get("background_color").(string))
}

func generatePatternCoverImage(emoji string, backgroundColor string) (string, error) {
	// We don't need to seed crypto/rand as it's automatically seeded with secure entropy
	// The hash is still useful for deterministic patterns if needed
	//seedValue := hash(emoji + backgroundColor)

	// Create a 300x300 RGBA image
	img := image.NewRGBA(image.Rect(0, 0, 300, 300))

	// Get the primary color based on the emoji/mood/weather
	primaryColor := getThemeColor(emoji)

	// Parse the background color from hex string (this will be the secondary color)
	secondaryColor, err := parseHexColor(backgroundColor)
	if err != nil {
		// Default to Spotify green if color parsing fails
		secondaryColor = color.RGBA{29, 185, 84, 255} // #1DB954 (Spotify green)
	}

	// Get the pattern type based on the emoji
	patternType := getPatternType(emoji)

	// Generate a visually appealing cover based on the pattern type
	switch patternType {
	case "gradient":
		// Create a gradient background
		drawGradient(img, primaryColor, secondaryColor)

	case "waves":
		// Create a wavy pattern (good for chill, relaxed moods)
		drawWavePattern(img, primaryColor, secondaryColor)

	case "rays":
		// Create a sunburst/ray pattern (good for energetic, sunny)
		drawRayPattern(img, primaryColor, secondaryColor)

	case "circles":
		// Create concentric circles (good for focus, target-oriented)
		drawCirclePattern(img, primaryColor, secondaryColor)

	case "dots":
		// Create a dotted pattern (good for playful, upbeat)
		drawDotPattern(img, primaryColor, secondaryColor)

	default:
		// Default to gradient if no specific pattern
		drawGradient(img, primaryColor, secondaryColor)
	}

	// Encode the image as JPEG
	var buf bytes.Buffer
	opts := jpeg.Options{
		Quality: 90,
	}
	if err := jpeg.Encode(&buf, img, &opts); err != nil {
		return "", fmt.Errorf("error encoding JPEG: %w", err)
	}

	// Convert JPEG to base64
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// drawGradient creates a gradient from top-left to bottom-right
func drawGradient(img *image.RGBA, startColor, endColor color.RGBA) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Calculate the gradient position (0.0 to 1.0)
			pos := float64(x+y) / float64(width+height)

			// Interpolate between the two colors
			r := uint8(float64(startColor.R) + pos*float64(endColor.R-startColor.R))
			g := uint8(float64(startColor.G) + pos*float64(endColor.G-startColor.G))
			b := uint8(float64(startColor.B) + pos*float64(endColor.B-startColor.B))

			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}
}

// drawWavePattern creates a wavy pattern
func drawWavePattern(img *image.RGBA, primaryColor, secondaryColor color.RGBA) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Fill with secondary color first
	draw.Draw(img, bounds, &image.Uniform{secondaryColor}, image.Point{}, draw.Src)

	// Draw waves with primary color
	for y := 0; y < height; y++ {
		// Calculate wave amplitude based on y position
		amplitude := 30.0 * math.Sin(float64(y)/20.0)

		for x := 0; x < width; x++ {
			// Create wave effect
			if float64(x) < (float64(width)/2)+amplitude {
				img.Set(x, y, primaryColor)
			}
		}
	}
}

// drawRayPattern creates a sunburst/ray pattern
func drawRayPattern(img *image.RGBA, primaryColor, secondaryColor color.RGBA) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	centerX, centerY := width/2, height/2

	// Fill with secondary color first
	draw.Draw(img, bounds, &image.Uniform{secondaryColor}, image.Point{}, draw.Src)

	// Draw rays from center
	numRays := 12
	for i := 0; i < numRays; i++ {
		angle := float64(i) * (2 * math.Pi / float64(numRays))

		for r := 0; r < width; r++ {
			// Calculate position
			x := centerX + int(float64(r)*math.Cos(angle))
			y := centerY + int(float64(r)*math.Sin(angle))

			// Check if in bounds
			if x >= 0 && x < width && y >= 0 && y < height {
				// Draw ray
				img.Set(x, y, primaryColor)

				// Make ray thicker
				for w := 1; w < 20; w++ {
					offsetAngle := angle + float64(w)*0.01
					offsetX := centerX + int(float64(r)*math.Cos(offsetAngle))
					offsetY := centerY + int(float64(r)*math.Sin(offsetAngle))

					if offsetX >= 0 && offsetX < width && offsetY >= 0 && offsetY < height {
						img.Set(offsetX, offsetY, primaryColor)
					}
				}
			}
		}
	}
}

// drawCirclePattern creates concentric circles
func drawCirclePattern(img *image.RGBA, primaryColor, secondaryColor color.RGBA) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	centerX, centerY := width/2, height/2

	// Fill with secondary color first
	draw.Draw(img, bounds, &image.Uniform{secondaryColor}, image.Point{}, draw.Src)

	// Draw concentric circles
	maxRadius := int(math.Sqrt(float64(width*width+height*height)) / 2)
	for r := maxRadius; r > 0; r -= 20 {
		// Alternate colors
		circleColor := primaryColor
		if (r/20)%2 == 0 {
			circleColor = secondaryColor
		}

		// Draw circle
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				// Calculate distance from center
				dx, dy := x-centerX, y-centerY
				distance := int(math.Sqrt(float64(dx*dx + dy*dy)))

				// Draw circle with some thickness
				if distance <= r && distance > r-10 {
					img.Set(x, y, circleColor)
				}
			}
		}
	}
}

// drawDotPattern creates a dotted pattern
func drawDotPattern(img *image.RGBA, primaryColor, secondaryColor color.RGBA) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Fill with secondary color first
	draw.Draw(img, bounds, &image.Uniform{secondaryColor}, image.Point{}, draw.Src)

	// Draw dots in a grid pattern
	dotSpacing := 30
	dotRadius := 10

	for y := dotSpacing / 2; y < height; y += dotSpacing {
		for x := dotSpacing / 2; x < width; x += dotSpacing {
			// Add some randomness to dot positions using crypto/rand
			offsetX, offsetY := secureRandomOffset(2)

			// Draw dot
			for dy := -dotRadius; dy <= dotRadius; dy++ {
				for dx := -dotRadius; dx <= dotRadius; dx++ {
					// Check if point is inside the circle
					if dx*dx+dy*dy <= dotRadius*dotRadius {
						px, py := x+dx+offsetX, y+dy+offsetY

						// Check bounds
						if px >= 0 && px < width && py >= 0 && py < height {
							img.Set(px, py, primaryColor)
						}
					}
				}
			}
		}
	}
}

// secureRandomInt generates a cryptographically secure random integer between 0 and max-1
func secureRandomInt(max int) (int, error) {
	if max <= 0 {
		return 0, fmt.Errorf("max must be greater than 0")
	}

	// Create a buffer to hold 8 bytes (uint64)
	b := make([]byte, 8)

	// Read random bytes
	_, err := rand.Read(b)
	if err != nil {
		return 0, fmt.Errorf("error generating random number: %w", err)
	}

	// Convert bytes to uint64
	n := binary.BigEndian.Uint64(b)

	// Return n mod max to get a number in the range [0, max-1]
	return int(n % uint64(max)), nil
}

// secureRandomOffset generates a random offset in the range [-range, +range]
func secureRandomOffset(offsetRange int) (int, int) {
	// Generate random X offset
	xRand, err := secureRandomInt(offsetRange*2 + 1)
	if err != nil {
		// If there's an error, return 0 (no offset)
		return 0, 0
	}

	// Generate random Y offset
	yRand, err := secureRandomInt(offsetRange*2 + 1)
	if err != nil {
		// If there's an error, return 0 (no offset)
		return xRand - offsetRange, 0
	}

	// Convert from [0, 2*range] to [-range, +range]
	return xRand - offsetRange, yRand - offsetRange
}

// parseHexColor converts a hex color string (e.g., "#FF5733") to color.RGBA
func parseHexColor(hexColor string) (color.RGBA, error) {
	// Remove the # prefix if present
	hexColor = strings.TrimPrefix(hexColor, "#")

	// Parse the hex color
	if len(hexColor) != 6 {
		return color.RGBA{}, fmt.Errorf("invalid hex color format: %s", hexColor)
	}

	// Parse the RGB components
	r, err := strconv.ParseUint(hexColor[0:2], 16, 8)
	if err != nil {
		return color.RGBA{}, fmt.Errorf("error parsing red component: %w", err)
	}

	g, err := strconv.ParseUint(hexColor[2:4], 16, 8)
	if err != nil {
		return color.RGBA{}, fmt.Errorf("error parsing green component: %w", err)
	}

	b, err := strconv.ParseUint(hexColor[4:6], 16, 8)
	if err != nil {
		return color.RGBA{}, fmt.Errorf("error parsing blue component: %w", err)
	}

	return color.RGBA{uint8(r), uint8(g), uint8(b), 255}, nil
}

// getPatternType returns the pattern type to use for the given emoji/mood/weather
func getPatternType(emoji string) string {
	// Map emojis to appropriate pattern types
	patternTypes := map[string]string{
		// Mood emojis
		"⚡":  "rays",     // energetic -> rays pattern
		"😌":  "waves",    // chill -> waves pattern
		"🧸":  "gradient", // cozy -> gradient pattern
		"😢":  "waves",    // melancholy -> waves pattern
		"🥳":  "dots",     // upbeat -> dots pattern
		"🧠":  "circles",  // focus -> circles pattern
		"💪":  "rays",     // workout -> rays pattern
		"❤️": "gradient", // romantic -> gradient pattern
		"😊":  "dots",     // happy -> dots pattern
		"😔":  "waves",    // sad -> waves pattern
		"😡":  "rays",     // angry -> rays pattern
		"🤩":  "dots",     // excited -> dots pattern

		// Weather emojis
		"☀️": "rays",     // sunny -> rays pattern
		"☁️": "gradient", // cloudy -> gradient pattern
		"🌧️": "waves",    // rainy -> waves pattern
		"❄️": "dots",     // snowy -> dots pattern
		"⛈️": "rays",     // stormy -> rays pattern
		"🌫️": "gradient", // foggy -> gradient pattern
		"🌬️": "waves",    // windy -> waves pattern
		"🔥":  "rays",     // hot -> rays pattern
		"🧊":  "circles",  // cold -> circles pattern
		"🌈":  "gradient", // clear -> gradient pattern
		"🎵":  "circles",  // default music note -> circles pattern
	}

	// Get the pattern type for the emoji
	patternType, ok := patternTypes[emoji]
	if !ok {
		// Default to gradient if emoji not found
		patternType = "gradient"
	}

	return patternType
}

// getThemeColor returns a color associated with the emoji/mood/weather
func getThemeColor(emoji string) color.RGBA {
	// Map emojis to colors
	themeColors := map[string]color.RGBA{
		"⚡":  color.RGBA{255, 215, 0, 255},   // Gold for energetic
		"😌":  color.RGBA{135, 206, 235, 255}, // Sky Blue for chill
		"🧸":  color.RGBA{139, 69, 19, 255},   // Brown for cozy
		"😢":  color.RGBA{70, 130, 180, 255},  // Steel Blue for melancholy
		"🥳":  color.RGBA{255, 105, 180, 255}, // Hot Pink for upbeat
		"🧠":  color.RGBA{128, 0, 128, 255},   // Purple for focus
		"💪":  color.RGBA{220, 20, 60, 255},   // Crimson for workout
		"❤️": color.RGBA{255, 0, 0, 255},     // Red for romantic
		"😊":  color.RGBA{255, 215, 0, 255},   // Gold for happy
		"😔":  color.RGBA{70, 130, 180, 255},  // Steel Blue for sad
		"😡":  color.RGBA{178, 34, 34, 255},   // Firebrick for angry
		"🤩":  color.RGBA{255, 140, 0, 255},   // Dark Orange for excited
		"☀️": color.RGBA{255, 215, 0, 255},   // Gold for sunny
		"☁️": color.RGBA{211, 211, 211, 255}, // Light Gray for cloudy
		"🌧️": color.RGBA{70, 130, 180, 255},  // Steel Blue for rainy
		"❄️": color.RGBA{255, 250, 250, 255}, // Snow for snowy
		"⛈️": color.RGBA{47, 79, 79, 255},    // Dark Slate Gray for stormy
		"🌫️": color.RGBA{220, 220, 220, 255}, // Gainsboro for foggy
		"🌬️": color.RGBA{176, 196, 222, 255}, // Light Steel Blue for windy
		"🔥":  color.RGBA{255, 69, 0, 255},    // Orange Red for hot
		"🧊":  color.RGBA{173, 216, 230, 255}, // Light Blue for cold
		"🌈":  color.RGBA{147, 112, 219, 255}, // Medium Purple for clear
	}

	// Return the color for the emoji, or a default color if not found
	if color, ok := themeColors[emoji]; ok {
		return color
	}

	// Default color (Spotify green)
	return color.RGBA{29, 185, 84, 255}
}

func getMoodEmoji(mood string) string {
	// Map moods to appropriate emojis
	moodEmojis := map[string]string{
		"energetic":  "⚡",
		"chill":      "😌",
		"cozy":       "🧸",
		"melancholy": "😢",
		"upbeat":     "🥳",
		"focus":      "🧠",
		"workout":    "💪",
		"romantic":   "❤️",
		"happy":      "😊",
		"sad":        "😔",
		"angry":      "😡",
		"relaxed":    "😌",
		"excited":    "🤩",
	}

	if emoji, ok := moodEmojis[strings.ToLower(mood)]; ok {
		return emoji
	}

	// Default emoji if mood is not recognized
	return "🎵"
}

func getWeatherEmoji(weather string) string {
	// Map weather conditions to appropriate emojis
	weatherEmojis := map[string]string{
		"sunny":        "☀️",
		"cloudy":       "☁️",
		"rainy":        "🌧️",
		"snowy":        "❄️",
		"stormy":       "⛈️",
		"foggy":        "🌫️",
		"windy":        "🌬️",
		"hot":          "🔥",
		"cold":         "🧊",
		"clear":        "🌈",
		"thunderstorm": "⚡",
	}

	if emoji, ok := weatherEmojis[strings.ToLower(weather)]; ok {
		return emoji
	}

	// Default emoji if weather condition is not recognized
	return "🌤️"
}

// getAccessToken retrieves the access token from the provider client
// This is a workaround since we can't access the token directly from the Spotify client
func getAccessToken(ctx context.Context) string {
	// First, try to get the access token from the environment variable
	token := os.Getenv("SPOTIFY_ACCESS_TOKEN")
	if token != "" {
		return token
	}

	// If the access token is not available, try to refresh it using the refresh token
	fmt.Println("Access token not found, attempting to refresh")
	refreshToken := os.Getenv("SPOTIFY_REFRESH_TOKEN")
	if refreshToken == "" {
		fmt.Println("Error: SPOTIFY_REFRESH_TOKEN environment variable not set")
		return ""
	}

	// Get client credentials from environment variables
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	redirectURI := os.Getenv("SPOTIFY_REDIRECT_URI")

	if clientID == "" || clientSecret == "" || redirectURI == "" {
		fmt.Println("Error: Missing Spotify client credentials in environment variables")
		return ""
	}

	// Create OAuth2 config
	oauthConfig := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURI,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.spotify.com/authorize",
			TokenURL: "https://accounts.spotify.com/api/token",
		},
		Scopes: []string{
			"playlist-modify-public",
			"playlist-modify-private",
			"ugc-image-upload",
		},
	}

	// Use refresh token to obtain a new access token
	tkn := &oauth2.Token{RefreshToken: refreshToken}
	tokenSource := oauthConfig.TokenSource(ctx, tkn)
	newToken, err := tokenSource.Token()
	if err != nil {
		fmt.Printf("Error refreshing access token: %s\n", err)
		return ""
	}

	// Save the new access token to the environment variable for future use
	os.Setenv("SPOTIFY_ACCESS_TOKEN", newToken.AccessToken)

	// Return the new access token
	return newToken.AccessToken
}

// Simple hash function for seeding the random number generator
func hash(s string) int {
	h := 0
	for i := 0; i < len(s); i++ {
		h = 31*h + int(s[i])
	}
	return h
}
