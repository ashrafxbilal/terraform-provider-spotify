package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

func dataSourceWeatherRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Detect location
	loc, err := detectLocation()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error detecting location: %s", err))
	}

	// Get weather data
	weather, err := getCurrentWeather(loc.Lat, loc.Lon)
	if err != nil {
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
	d.Set("temperature", weather.Current.Temperature)
	d.Set("lat", loc.Lat)
	d.Set("lon", loc.Lon)
	d.Set("city", loc.City)
	d.Set("suggested_moods", suggestedMoods)
	d.Set("mood", selectedMood)
	d.Set("is_sunny", weather.Current.Temperature > 20)

	return diags
}

// detectLocation gets the user's location based on IP
func detectLocation() (*GeoLocation, error) {
	resp, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return nil, fmt.Errorf("error getting location: %w", err)
	}
	defer resp.Body.Close()

	var loc GeoLocation
	if err := json.NewDecoder(resp.Body).Decode(&loc); err != nil {
		return nil, fmt.Errorf("error decoding location: %w", err)
	}
	return &loc, nil
}

// getCurrentWeather gets weather data for a location
func getCurrentWeather(lat, lon float64) (*WeatherResponse, error) {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f&current=temperature_2m&timezone=auto", lat, lon)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching weather: %w", err)
	}
	defer resp.Body.Close()

	var weather WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return nil, fmt.Errorf("error decoding weather response: %w", err)
	}
	return &weather, nil
}