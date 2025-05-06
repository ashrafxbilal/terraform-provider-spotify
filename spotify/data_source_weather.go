package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/ashrafxbilal/terraform-provider-spotify/spotify/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// GeoLocation represents location data
type GeoLocation struct {
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
	City string  `json:"city"`
}

// WeatherResponse represents weather data
type WeatherResponse struct {
	Current struct {
		Temperature float64 `json:"temperature_2m"`
	} `json:"current"`
}

func dataSourceWeather() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWeatherRead,
		Schema: map[string]*schema.Schema{
			"temperature": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Current temperature in Celsius",
			},
			"lat": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Latitude of the detected location",
			},
			"lon": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Longitude of the detected location",
			},
			"city": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "City name of the detected location",
			},
			"suggested_moods": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Three suggested moods based on the current weather",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"mood": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Selected mood for playlist generation. If not provided, defaults to first suggested mood",
			},
			"is_sunny": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the current weather is considered sunny (temperature > 20°C)",
			},
		},
	}
}

// setResourceDataWithErrorCheck sets a value in the ResourceData and checks for errors
func setResourceDataWithErrorCheck(d *schema.ResourceData, key string, value interface{}, ctx context.Context) diag.Diagnostics {
	logger := logging.DefaultLogger.WithContext(ctx)
	if err := d.Set(key, value); err != nil {
		logger.Error(fmt.Sprintf("Error setting %s", key), "error", err.Error())
		return diag.FromErr(fmt.Errorf("error setting %s: %s", key, err))
	}
	return nil
}

func dataSourceWeatherRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	logger := logging.DefaultLogger.WithContext(ctx)

	// Detect location
	loc, err := detectLocation()
	if err != nil {
		logger.Error("Failed to detect location", "error", err.Error())
		return diag.FromErr(fmt.Errorf("error detecting location: %s", err))
	}

	// Get weather data
	weather, err := getCurrentWeather(loc.Lat, loc.Lon)
	if err != nil {
		logger.Error("Failed to get weather data", "error", err.Error())
		return diag.FromErr(fmt.Errorf("error getting weather data: %s", err))
	}

	// Generate three suggested moods based on temperature
	var suggestedMoods []string
	var defaultMood string

	// Primary mood based on temperature
	if weather.Current.Temperature > 25 {
		defaultMood = "energetic"
		suggestedMoods = []string{"energetic", "upbeat", "lively"}
	} else if weather.Current.Temperature < 10 {
		defaultMood = "cozy"
		suggestedMoods = []string{"cozy", "mellow", "relaxed"}
	} else {
		defaultMood = "chill"
		suggestedMoods = []string{"chill", "balanced", "focused"}
	}

	// Check if user provided a custom mood
	var selectedMood string
	if v, ok := d.GetOk("mood"); ok {
		selectedMood = v.(string)
	} else {
		// Default to first suggested mood
		selectedMood = defaultMood
	}

	// Set values
	d.SetId(fmt.Sprintf("%f-%f-%d", loc.Lat, loc.Lon, time.Now().Unix()))

	// Set resource data with error checking
	if diagErr := setResourceDataWithErrorCheck(d, "temperature", weather.Current.Temperature, ctx); diagErr != nil {
		return diagErr
	}

	if diagErr := setResourceDataWithErrorCheck(d, "lat", loc.Lat, ctx); diagErr != nil {
		return diagErr
	}

	if diagErr := setResourceDataWithErrorCheck(d, "lon", loc.Lon, ctx); diagErr != nil {
		return diagErr
	}

	if diagErr := setResourceDataWithErrorCheck(d, "city", loc.City, ctx); diagErr != nil {
		return diagErr
	}

	if diagErr := setResourceDataWithErrorCheck(d, "suggested_moods", suggestedMoods, ctx); diagErr != nil {
		return diagErr
	}

	if diagErr := setResourceDataWithErrorCheck(d, "mood", selectedMood, ctx); diagErr != nil {
		return diagErr
	}

	if diagErr := setResourceDataWithErrorCheck(d, "is_sunny", weather.Current.Temperature > 20, ctx); diagErr != nil {
		return diagErr
	}

	return diags
}

// LocationRequestBuilder implements the Builder pattern for location API requests
type LocationRequestBuilder struct {
	baseURL string
	timeout time.Duration
}

// NewLocationRequestBuilder creates a new builder with default values
func NewLocationRequestBuilder() *LocationRequestBuilder {
	return &LocationRequestBuilder{
		baseURL: "http://ip-api.com/json/",
		timeout: 10 * time.Second,
	}
}

// WithTimeout sets a custom timeout for the HTTP request
func (b *LocationRequestBuilder) WithTimeout(timeout time.Duration) *LocationRequestBuilder {
	b.timeout = timeout
	return b
}

// Execute builds the request, executes it, and returns the location response
func (b *LocationRequestBuilder) Execute() (*GeoLocation, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: b.timeout,
	}

	// Execute the request
	resp, err := client.Get(b.baseURL)
	if err != nil {
		return nil, fmt.Errorf("error getting location: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// Decode the response
	var loc GeoLocation
	if err := json.NewDecoder(resp.Body).Decode(&loc); err != nil {
		return nil, fmt.Errorf("error decoding location: %w", err)
	}

	return &loc, nil
}

// detectLocation gets the user's location based on IP using the Builder pattern
func detectLocation() (*GeoLocation, error) {
	// Create a new builder and execute the request
	return NewLocationRequestBuilder().Execute()
}

// WeatherRequestBuilder implements the Builder pattern for weather API requests
type WeatherRequestBuilder struct {
	baseURL    string
	latitude   float64
	longitude  float64
	parameters map[string]string
	timeout    time.Duration
}

// NewWeatherRequestBuilder creates a new builder with default values
func NewWeatherRequestBuilder() *WeatherRequestBuilder {
	return &WeatherRequestBuilder{
		baseURL:    "https://api.open-meteo.com/v1/forecast",
		parameters: make(map[string]string),
		timeout:    10 * time.Second,
	}
}

// WithCoordinates sets the latitude and longitude
func (b *WeatherRequestBuilder) WithCoordinates(lat, lon float64) *WeatherRequestBuilder {
	b.latitude = lat
	b.longitude = lon
	return b
}

// WithParameter adds a custom parameter to the request
func (b *WeatherRequestBuilder) WithParameter(key, value string) *WeatherRequestBuilder {
	b.parameters[key] = value
	return b
}

// WithTimeout sets a custom timeout for the HTTP request
func (b *WeatherRequestBuilder) WithTimeout(timeout time.Duration) *WeatherRequestBuilder {
	b.timeout = timeout
	return b
}

// Validate checks if the request parameters are valid
func (b *WeatherRequestBuilder) Validate() error {
	// Validate coordinates
	if b.latitude < -90 || b.latitude > 90 {
		return fmt.Errorf("invalid latitude: %.4f (must be between -90 and 90)", b.latitude)
	}
	if b.longitude < -180 || b.longitude > 180 {
		return fmt.Errorf("invalid longitude: %.4f (must be between -180 and 180)", b.longitude)
	}
	return nil
}

// Build constructs the URL and returns it
func (b *WeatherRequestBuilder) Build() (string, error) {
	// Validate parameters
	if err := b.Validate(); err != nil {
		return "", err
	}

	// Parse base URL
	parsedURL, err := url.Parse(b.baseURL)
	if err != nil {
		return "", fmt.Errorf("error parsing base URL: %w", err)
	}

	// Add parameters
	params := url.Values{}
	params.Add("latitude", fmt.Sprintf("%.4f", b.latitude))
	params.Add("longitude", fmt.Sprintf("%.4f", b.longitude))

	// Add default parameters if not overridden
	if _, exists := b.parameters["current"]; !exists {
		params.Add("current", "temperature_2m")
	}

	if _, exists := b.parameters["timezone"]; !exists {
		params.Add("timezone", "auto")
	}

	// Add any custom parameters
	for key, value := range b.parameters {
		params.Add(key, value)
	}

	parsedURL.RawQuery = params.Encode()
	return parsedURL.String(), nil
}

// Execute builds the request, executes it, and returns the weather response
func (b *WeatherRequestBuilder) Execute() (*WeatherResponse, error) {
	// Build the URL
	requestURL, err := b.Build()
	if err != nil {
		return nil, err
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: b.timeout,
	}

	// Execute the request
	resp, err := client.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("error fetching weather: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// Decode the response
	var weather WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return nil, fmt.Errorf("error decoding weather response: %w", err)
	}

	return &weather, nil
}

// getCurrentWeather gets weather data for a location using the Builder pattern
func getCurrentWeather(lat, lon float64) (*WeatherResponse, error) {
	// Create a new builder and configure it
	weatherRequest := NewWeatherRequestBuilder().
		WithCoordinates(lat, lon).
		WithParameter("current", "temperature_2m").
		WithParameter("timezone", "auto")

	// Execute the request
	return weatherRequest.Execute()
}
