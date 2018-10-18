package discord

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		// Provider
		"discord_provider_token":        "The token to use for authenticating with the Discord API",

		// Server
		"discord_resource_server_name":  "What to name the new server",
		"discord_resource_server_empty": "Whether or not the server should be created with no channels",

		// Channel
		"discord_resource_channel_guild": "What guild this channel belongs to",
		"discord_resource_channel_name":  "What to name the new channel",
		"discord_resource_channel_type":  "What type of channel is this (text, voice, category)",
	}
}
