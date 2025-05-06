package spotify

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
	

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zmb3/spotify/v2"
)

// min returns the smaller of x or y.
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func dataSourceSpotifyTracks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSpotifyTracksRead,
		Schema: map[string]*schema.Schema{
			"genre": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Genre to search for tracks",
			},
			"artist": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Artist to search for tracks",
			},
			"mood": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Mood to search for tracks (energetic, chill, cozy, etc.)",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     20,
				Description: "Maximum number of tracks to return",
			},
			"time_range": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Time range for tracks (short_term, medium_term, long_term)",
			},
			"popularity": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Minimum popularity score (0-100) for tracks",
			},
			"ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of track IDs",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"names": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of track names",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"artists": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of track artists",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceSpotifyTracksRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*ProviderClient).SpotifyClient

	rand.Seed(time.Now().UnixNano())

	genre := d.Get("genre").(string)
	artist := d.Get("artist").(string)
	mood := d.Get("mood").(string)
	limit := d.Get("limit").(int)
	timeRange := d.Get("time_range").(string)
	popularity := d.Get("popularity").(int)

	// mapping mood to audio features
	var minEnergy, maxEnergy, minValence, maxValence, minAcousticness, maxAcousticness float32
	var minTempo, maxTempo, minDanceability, maxDanceability, minInstrumentalness, maxInstrumentalness float32
	
	switch mood {
	case "energetic":
		minEnergy = 0.7
		maxEnergy = 1.0
		minTempo = 120
		maxTempo = 180
		minValence = 0.6
		maxValence = 1.0
		minDanceability = 0.6
		maxDanceability = 1.0
	case "chill":
		minEnergy = 0.3
		maxEnergy = 0.6
		minTempo = 70
		maxTempo = 110
		minValence = 0.4
		maxValence = 0.7
		minDanceability = 0.3
		maxDanceability = 0.6
	case "cozy":
		minEnergy = 0.1
		maxEnergy = 0.4
		minTempo = 60
		maxTempo = 90
		minValence = 0.3
		maxValence = 0.6
		minAcousticness = 0.5
		maxAcousticness = 1.0
	case "melancholy":
		minEnergy = 0.2
		maxEnergy = 0.5
		minTempo = 60
		maxTempo = 90
		minValence = 0.0
		maxValence = 0.4
		minAcousticness = 0.4
		maxAcousticness = 0.8
	case "upbeat":
		minEnergy = 0.6
		maxEnergy = 0.9
		minTempo = 100
		maxTempo = 140
		minValence = 0.7
		maxValence = 1.0
		minDanceability = 0.5
		maxDanceability = 0.9
	case "focus":
		minEnergy = 0.3
		maxEnergy = 0.7
		minTempo = 80
		maxTempo = 120
		minValence = 0.3
		maxValence = 0.7
		minInstrumentalness = 0.5
		maxInstrumentalness = 1.0
	case "workout":
		minEnergy = 0.8
		maxEnergy = 1.0
		minTempo = 130
		maxTempo = 200
		minValence = 0.5
		maxValence = 1.0
		minDanceability = 0.6
		maxDanceability = 1.0
	case "romantic":
		minEnergy = 0.3
		maxEnergy = 0.6
		minTempo = 70
		maxTempo = 110
		minValence = 0.5
		maxValence = 0.8
		minAcousticness = 0.3
		maxAcousticness = 0.7
	default:
		// a balanced mix
		minEnergy = 0.4
		maxEnergy = 0.7
		minTempo = 90
		maxTempo = 130
		minValence = 0.4
		maxValence = 0.7
		minDanceability = 0.4
		maxDanceability = 0.7
	}

	// building recommendations options
	seeds := spotify.Seeds{}
	if genre != "" {
		seeds.Genres = []string{genre}
	}
	if artist != "" {
		// Search for the artist first
		searchResult, err := client.Search(ctx, artist, spotify.SearchTypeArtist, spotify.Limit(1))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error searching for artist: %s", err))
		}

		if searchResult.Artists != nil && len(searchResult.Artists.Artists) > 0 {
			seeds.Artists = []spotify.ID{searchResult.Artists.Artists[0].ID}
		}
	}
	
	// If time_range is specified, use it to get tracks from user's top artists
	if timeRange != "" {
		// Convert time_range to spotify.Range
		var spotifyRange spotify.Range
		switch timeRange {
		case "short_term":
			spotifyRange = spotify.ShortTermRange
		case "medium_term":
			spotifyRange = spotify.MediumTermRange
		case "long_term":
			spotifyRange = spotify.LongTermRange
		default:
			spotifyRange = spotify.ShortTermRange // Default to short_term for recent tracks
		}
		
		// Get user's top artists for the specified time range
		topArtists, err := client.CurrentUsersTopArtists(ctx, spotify.Limit(3), spotify.Timerange(spotifyRange))
		if err == nil && len(topArtists.Artists) > 0 {
			// Use the top artist as a seed if we don't already have an artist seed
			if len(seeds.Artists) == 0 && len(topArtists.Artists) > 0 {
				seeds.Artists = []spotify.ID{topArtists.Artists[0].ID}
			}
		}
	}

	// If no seeds are provided, use default genres based on mood
	if len(seeds.Artists) == 0 && len(seeds.Tracks) == 0 && len(seeds.Genres) == 0 {
		// Try to get available genre seeds from Spotify
		availableGenres, err := client.GetAvailableGenreSeeds(ctx)
		if err != nil {
			// If we can't get available genres, use hardcoded defaults
			defaultGenres := []string{"pop"}
			
			// Select appropriate genres based on mood
			switch mood {
			case "energetic":
				defaultGenres = []string{"pop", "dance", "electronic"}
			case "chill":
				defaultGenres = []string{"acoustic", "ambient", "chill"}
			case "cozy":
				defaultGenres = []string{"jazz", "folk", "indie"}
			default:
				defaultGenres = []string{"pop"}
			}
			
			// Use the first genre as seed
			seeds.Genres = []string{defaultGenres[0]}
		} else {
			// We have available genres, select appropriate ones based on mood
			var genreOptions []string
			switch mood {
			case "energetic":
				genreOptions = []string{"pop", "dance", "electronic", "edm", "party"}
			case "chill":
				genreOptions = []string{"acoustic", "ambient", "chill", "study"}
			case "cozy":
				genreOptions = []string{"jazz", "folk", "indie", "indie-pop"}
			default:
				genreOptions = []string{"pop", "rock"}
			}
			
			// Find the first available genre that matches our options
			for _, option := range genreOptions {
				for _, available := range availableGenres {
					if option == available {
						seeds.Genres = []string{option}
						break
					}
				}
				if len(seeds.Genres) > 0 {
					break
				}
			}
			
			// If no matching genre was found, use a safe default
			if len(seeds.Genres) == 0 && len(availableGenres) > 0 {
				seeds.Genres = []string{availableGenres[0]}
			}
		}
	}

	// audio features based on mood
	attrs := spotify.NewTrackAttributes()
	// Set energy parameters
	attrs = attrs.TargetEnergy(float64(minEnergy))
	attrs = attrs.MaxEnergy(float64(maxEnergy))
	
	// Set tempo parameters
	attrs = attrs.TargetTempo(float64(minTempo))
	if maxTempo > 0 {
		attrs = attrs.MaxTempo(float64(maxTempo))
	}
	
	// Set valence (positivity) parameters
	attrs = attrs.TargetValence(float64(minValence))
	attrs = attrs.MaxValence(float64(maxValence))
	
	// Set danceability parameters if defined for this mood
	if maxDanceability > 0 {
		attrs = attrs.TargetDanceability(float64(minDanceability))
		attrs = attrs.MaxDanceability(float64(maxDanceability))
	}
	
	// Set acousticness parameters if defined for this mood
	if maxAcousticness > 0 {
		attrs = attrs.TargetAcousticness(float64(minAcousticness))
		attrs = attrs.MaxAcousticness(float64(maxAcousticness))
	}
	
	// Set instrumentalness parameters if defined for this mood
	if maxInstrumentalness > 0 {
		attrs = attrs.TargetInstrumentalness(float64(minInstrumentalness))
		attrs = attrs.MaxInstrumentalness(float64(maxInstrumentalness))
	}
	
	// Set popularity if specified
	if popularity > 0 {
		attrs = attrs.MinPopularity(popularity)
	}
	
	// For Global 50 or new releases, we'll handle this differently
	// The Spotify API doesn't directly support filtering by release date in recommendations
	// Instead, we'll use the popularity parameter and time_range to target recent popular tracks

	// Get recommendations
	recommendations, err := client.GetRecommendations(ctx, seeds, attrs, spotify.Limit(limit))
	if err != nil {
		// Provide more detailed error information
		detailedErr := fmt.Errorf("error getting recommendations: %s\nSeeds used: %+v\nMood: %s", err, seeds, mood)
		
		// Log the error details for debugging
		fmt.Printf("Spotify API Error: %s\nFalling back to Search API\n", detailedErr)
		
		// Build a search query based on genre/mood
		searchQuery := ""
		if genre != "" {
			// If we have a genre, use it as the primary search term
			// but use simpler search queries that are more likely to succeed
			if mood != "" {
				switch mood {
				case "energetic":
					// Simplified query that's more likely to return results
					searchQuery = fmt.Sprintf("%s dance party", genre)
				case "chill":
					searchQuery = fmt.Sprintf("%s chill relax", genre)
				case "cozy":
					searchQuery = fmt.Sprintf("%s acoustic mellow", genre)
				case "melancholy":
					searchQuery = fmt.Sprintf("%s sad emotional", genre)
				case "upbeat":
					searchQuery = fmt.Sprintf("%s happy upbeat", genre)
				case "focus":
					searchQuery = fmt.Sprintf("%s focus instrumental", genre)
				case "workout":
					searchQuery = fmt.Sprintf("%s workout energetic", genre)
				case "romantic":
					searchQuery = fmt.Sprintf("%s love romantic", genre)
				default:
					searchQuery = fmt.Sprintf("%s", genre)
				}
			} else {
				searchQuery = fmt.Sprintf("%s", genre)
			}
		} else if mood != "" {
			// Map mood to appropriate search terms with more specific combinations
			switch mood {
			case "energetic":
				searchQuery = "tag:party tag:upbeat tag:dance energy:>0.7 popularity:>70"
			case "chill":
				searchQuery = "tag:chill tag:relaxing tag:ambient energy:<0.6 popularity:>60"
			case "cozy":
				searchQuery = "tag:mellow tag:acoustic tag:indie acoustic:>0.6 popularity:>50"
			case "melancholy":
				searchQuery = "tag:sad tag:emotional tag:melancholy valence:<0.4 popularity:>50"
			case "upbeat":
				searchQuery = "tag:happy tag:upbeat tag:summer valence:>0.7 popularity:>70"
			case "focus":
				searchQuery = "tag:focus tag:concentration tag:study instrumentalness:>0.5 popularity:>50"
			case "workout":
				searchQuery = "tag:workout tag:gym tag:fitness energy:>0.8 tempo:>130 popularity:>70"
			case "romantic":
				searchQuery = "tag:love tag:romantic tag:date valence:>0.5 popularity:>60"
			default:
				searchQuery = "tag:popular year:2020-2023"
			}
		} else {
			searchQuery = "tag:popular year:2020-2023 popularity:>75"
		}
		
		fmt.Printf("Using search query: %s\n", searchQuery)
		
		// Use the Search API instead
		searchResult, err := client.Search(ctx, searchQuery, spotify.SearchTypeTrack, spotify.Limit(limit))
		if err != nil {
			return diag.FromErr(fmt.Errorf("fallback search also failed: %s", err))
		}
		
		// Process search results instead of recommendations
		if searchResult.Tracks != nil && len(searchResult.Tracks.Tracks) > 0 {
			// Sort tracks by popularity (highest first) to prioritize more popular versions of similar tracks
			tracks := searchResult.Tracks.Tracks
			sort.Slice(tracks, func(i, j int) bool {
				return tracks[i].Popularity > tracks[j].Popularity
			})
			
			// Create maps to track seen track names, artists, and word prefixes to avoid duplicates
			seenNames := make(map[string]bool)
			seenArtists := make(map[string]int)
			seenPrefixes := make(map[string]int)
			
			// Filter tracks to avoid duplicates and ensure popular tracks
			var filteredTracks []spotify.FullTrack
			for _, track := range tracks {
				// Skip tracks with low popularity
				if track.Popularity < 50 {
					continue
				}
				
				// Use exact track name for comparison to avoid duplicates like "Bella's Finals"
				trackName := strings.ToLower(track.Name)
				
				// Get primary artist
				primaryArtist := ""
				if len(track.Artists) > 0 {
					primaryArtist = track.Artists[0].Name
				}
				
				// Limit tracks per artist to 2 to ensure variety
				artistCount := seenArtists[primaryArtist]
				if artistCount >= 2 {
					continue
				}
				
				// Check for common prefixes to avoid too many similar tracks
				words := strings.Fields(trackName)
				if len(words) > 0 {
					firstWord := words[0]
					// Apply more strict filtering for common prefixes to ensure diversity
					// Check if we've seen this prefix too many times
					if len(firstWord) > 3 { // Only consider meaningful prefixes (longer than 3 chars)
						// For any prefix, limit to just 1 song if it appears frequently in results
						if seenPrefixes[firstWord] >= 1 {
							continue
						}
					} else if len(firstWord) > 3 { // For other meaningful prefixes
						// Limit tracks with the same first word to 2 for better diversity
						if seenPrefixes[firstWord] >= 2 {
							continue
						}
					}
					seenPrefixes[firstWord]++
				}
				
				// If we haven't seen this exact track name before, add the track
				if trackName != "" && !seenNames[trackName] {
					seenNames[trackName] = true
					seenArtists[primaryArtist]++
					filteredTracks = append(filteredTracks, track)
				}
			}
			
			// Use the filtered tracks list
			trackIDs := make([]string, len(filteredTracks))
			trackNames := make([]string, len(filteredTracks))
			trackArtists := make([]string, len(filteredTracks))
			
			for i, track := range filteredTracks {
				trackIDs[i] = string(track.ID)
				trackNames[i] = track.Name
				if len(track.Artists) > 0 {
					trackArtists[i] = track.Artists[0].Name
				}
			}
			
			d.SetId(fmt.Sprintf("%d-%s-%s-%s-search", time.Now().Unix(), genre, artist, mood))
			d.Set("ids", trackIDs)
			d.Set("names", trackNames)
			d.Set("artists", trackArtists)
			
			return diags
		}
		
		return diag.FromErr(fmt.Errorf("no tracks found with search query: %s", searchQuery))
	}

	// Check if we got any tracks back
	if len(recommendations.Tracks) == 0 {
		// No tracks were returned, provide a helpful message
		detailedErr := fmt.Errorf("no tracks returned for the given criteria. Seeds used: %+v, Mood: %s", seeds, mood)
		fmt.Printf("Spotify API Warning: %s\n", detailedErr)
		
		// Continue execution but log the warning
	}

	// For recommendations, we can't sort by popularity as SimpleTrack doesn't have that field
	// But we can still filter by name prefix to avoid similar tracks
	tracks := recommendations.Tracks
	
	// Create maps to track seen track names, artists, and word prefixes to avoid duplicates
	seenNames := make(map[string]bool)
	seenArtists := make(map[string]int)
	seenPrefixes := make(map[string]int)
	
	// Filter tracks to avoid duplicates and ensure variety
	var filteredTracks []spotify.SimpleTrack
	for _, track := range tracks {
		// Use exact track name for comparison to avoid duplicates like "Bella's Finals"
		trackName := strings.ToLower(track.Name)
		
		// Get primary artist
		primaryArtist := ""
		if len(track.Artists) > 0 {
			primaryArtist = track.Artists[0].Name
		}
		
		// Limit tracks per artist to 2 to ensure variety
		artistCount := seenArtists[primaryArtist]
		if artistCount >= 2 {
			continue
		}
		
		// Check for common prefixes to avoid too many similar tracks
		words := strings.Fields(trackName)
		if len(words) > 0 {
			firstWord := words[0]
			// Apply more strict filtering for common prefixes to ensure diversity
			// Check if we've seen this prefix too many times
			if len(firstWord) > 3 { // Only consider meaningful prefixes (longer than 3 chars)
				// For any prefix, limit to just 1 song if it appears frequently in results
				if seenPrefixes[firstWord] >= 1 {
					continue
				}
			} else if len(firstWord) > 3 { // For other meaningful prefixes
				// Limit tracks with the same first word to 2 for better diversity
				if seenPrefixes[firstWord] >= 2 {
					continue
				}
			}
			seenPrefixes[firstWord]++
		}
		
		// If we haven't seen this exact track name before, add the track
		if trackName != "" && !seenNames[trackName] {
			seenNames[trackName] = true
			seenArtists[primaryArtist]++
			filteredTracks = append(filteredTracks, track)
		}
	}

	trackIDs := make([]string, len(filteredTracks))
	trackNames := make([]string, len(filteredTracks))
	trackArtists := make([]string, len(filteredTracks))

	for i, track := range filteredTracks {
		trackIDs[i] = string(track.ID)
		trackNames[i] = track.Name
		if len(track.Artists) > 0 {
			trackArtists[i] = track.Artists[0].Name
		}
	}

	d.SetId(fmt.Sprintf("%d-%s-%s-%s", time.Now().Unix(), genre, artist, mood))
	d.Set("ids", trackIDs)
	d.Set("names", trackNames)
	d.Set("artists", trackArtists)

	return diags
}