package main

import (
	"github.com/aequasi/discord-terraform/discord"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: discord.Provider})
}
