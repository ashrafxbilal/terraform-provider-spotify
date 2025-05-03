package spotify

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Spotify user ID",
			},
			"display_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user's display name",
			},
			"email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user's email address",
			},
			"product": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user's Spotify subscription level",
			},
			"followers": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of followers",
			},
			"images": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The user's profile images",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*ProviderClient).SpotifyClient

	// Get the current user's profile
	user, err := client.CurrentUser(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting current user: %s", err))
	}

	// Set the values
	d.SetId(user.ID)
	d.Set("display_name", user.DisplayName)
	d.Set("email", user.Email)
	d.Set("product", user.Product)
	d.Set("followers", user.Followers.Count)

	// Extract image URLs
	images := make([]string, len(user.Images))
	for i, img := range user.Images {
		images[i] = img.URL
	}
	d.Set("images", images)

	return diags
}