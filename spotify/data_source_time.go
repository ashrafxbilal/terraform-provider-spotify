package spotify

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTime() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTimeRead,
		Schema: map[string]*schema.Schema{
			"current_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current time in RFC3339 format",
			},
			"hour": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Current hour (0-23)",
			},
			"minute": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Current minute (0-59)",
			},
			"day_of_week": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current day of the week (Monday, Tuesday, etc.)",
			},
			"is_weekend": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the current day is a weekend (Saturday or Sunday)",
			},
			"time_of_day": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time of day (morning, afternoon, evening, night)",
			},
			"suggested_moods": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Three suggested moods based on time of day and day of week",
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
			"suggested_genres": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Three suggested genres based on time of day and day of week",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"genre": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Selected genre for playlist generation. If not provided, defaults to first suggested genre",
			},
		},
	}
}

func dataSourceTimeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get current time in local timezone
	now := time.Now().Local()

	// Format day of week
	dayOfWeek := now.Weekday().String()

	// Determine if it's a weekend
	isWeekend := now.Weekday() == time.Saturday || now.Weekday() == time.Sunday

	// Determine time of day
	hour := now.Hour()
	timeOfDay := "night" // Default
	if hour >= 5 && hour < 12 {
		timeOfDay = "morning"
	} else if hour >= 12 && hour < 17 {
		timeOfDay = "afternoon"
	} else if hour >= 17 && hour < 22 {
		timeOfDay = "evening"
	}

	// Generate suggested moods based on time of day and day of week
	suggestedMoods := getSuggestedMoods(timeOfDay, isWeekend)
	defaultMood := suggestedMoods[0]

	// Generate suggested genres based on time of day and day of week
	suggestedGenres := getSuggestedGenres(timeOfDay, isWeekend)
	defaultGenre := suggestedGenres[0]

	// Check if user provided a custom mood
	var selectedMood string
	if v, ok := d.GetOk("mood"); ok {
		selectedMood = v.(string)
	} else {
		// Default to first suggested mood
		selectedMood = defaultMood
	}

	// Check if user provided a custom genre
	var selectedGenre string
	if v, ok := d.GetOk("genre"); ok {
		selectedGenre = v.(string)
	} else {
		// Default to first suggested genre
		selectedGenre = defaultGenre
	}

	// Set ID and values
	d.SetId(fmt.Sprintf("%d", now.Unix()))
	d.Set("current_time", now.Format(time.RFC3339))
	d.Set("hour", hour)
	d.Set("minute", now.Minute())
	d.Set("day_of_week", dayOfWeek)
	d.Set("is_weekend", isWeekend)
	d.Set("time_of_day", timeOfDay)
	d.Set("suggested_moods", suggestedMoods)
	d.Set("mood", selectedMood)
	d.Set("suggested_genres", suggestedGenres)
	d.Set("genre", selectedGenre)

	return diags
}

// getSuggestedMoods returns three suggested moods based on time of day and whether it's a weekend
func getSuggestedMoods(timeOfDay string, isWeekend bool) []string {
	switch timeOfDay {
	case "morning":
		if isWeekend {
			return []string{"relaxed", "peaceful", "refreshed"}
		}
		return []string{"focused", "motivated", "energized"}
	case "afternoon":
		if isWeekend {
			return []string{"energetic", "playful", "adventurous"}
		}
		return []string{"productive", "determined", "inspired"}
	case "evening":
		if isWeekend {
			return []string{"party", "excited", "social"}
		}
		return []string{"chill", "relaxed", "unwinding"}
	case "night":
		return []string{"chill", "dreamy", "reflective"}
	default:
		return []string{"balanced", "neutral", "content"}
	}
}

// getSuggestedGenres returns three suggested genres based on time of day and whether it's a weekend
func getSuggestedGenres(timeOfDay string, isWeekend bool) []string {
	switch timeOfDay {
	case "morning":
		if isWeekend {
			return []string{"acoustic", "folk", "indie folk"}
		}
		return []string{"pop", "upbeat", "motivational"}
	case "afternoon":
		if isWeekend {
			return []string{"dance", "pop", "hip-hop"}
		}
		return []string{"rock", "alternative", "indie rock"}
	case "evening":
		if isWeekend {
			return []string{"electronic", "dance", "house"}
		}
		return []string{"indie", "alternative", "chill"}
	case "night":
		return []string{"ambient", "lo-fi", "chill electronic"}
	default:
		return []string{"pop", "top hits", "variety"}
	}
}