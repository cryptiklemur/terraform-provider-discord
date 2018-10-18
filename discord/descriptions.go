package discord

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"discord_provider_token":        "The token to use for authenticating with the Discord API",
		"discord_resource_server_name":  "What to name the new server",
		"discord_resource_server_empty": "Whether or not the server should be created with no channels",
	}
}
