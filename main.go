package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/ashrafxbilal/terraform-provider-spotify/spotify"
	"github.com/ashrafxbilal/terraform-provider-spotify/version"
	"log"
)

func main() {
	// Log version information
	log.Printf(version.GetVersionInfo())

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: spotify.Provider,
	})
}