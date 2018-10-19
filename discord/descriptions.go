package discord

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		// Provider
		"discord_provider_token": "The token to use for authenticating with the Discord API",

		// Server
		"discord_resource_server_name":  "What to name the new server",
		"discord_resource_server_empty": "Whether or not the server should be created with no channels",
		"discord_resource_server_region": "Guild Voice Region",
		"discord_resource_server_verification_level": "Verification Level: 0 - None, 1 - Low, 2 - Medium, 3 - High, 4 - Very High",
		"discord_resource_server_default_message_notifications": "Message Notification Level: 0 - All Messages, 1 - Only Mentions",
		"discord_resource_server_afk_channel_id": "ID For AFK Channel",
		"discord_resource_server_afk_timeout": "AFK Timeout in Seconds",
		"discord_resource_server_icon_url": "URL to image for the server icon",
		"discord_resource_server_icon_local": "Local path to image for the server icon",
		"discord_resource_server_icon_data_uri": "Base 64 encoded Data URI for server icon",
		"discord_resource_server_icon_hash": "Server Icon Hash",
		"discord_resource_server_owner_id": "ID of user to own the server",

		// Channel
		"discord_resource_channel_guild": "What guild this channel belongs to",
		"discord_resource_channel_name":  "What to name the new channel",
		"discord_resource_channel_type":  "What type of channel is this (text, voice, category)",

		// Invite
		"discord_resource_invite_channel":   "What channel this invite belongs to",
		"discord_resource_invite_max_age":   "Duration of invite in seconds before expiry, or 0 for never",
		"discord_resource_invite_max_uses":  "Max number of uses or 0 for unlimited",
		"discord_resource_invite_temporary": "Whether this invite only grants temporary membership",
		"discord_resource_invite_unique":    "if true, don't try to reuse a similar invite (useful for creating many unique one time use invites)",
	}
}
