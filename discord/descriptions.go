package discord

import "fmt"

var descriptions map[string]string

func init() {
    oauthUrl := "https://discord-token.vercel.app/grant"

    descriptions = map[string]string{
        // Provider
        "discord_provider_token":     "The token to use for authenticating with the Discord API",
        "discord_provider_client_id": "The client to use for authenticating with the Discord API",
        "discord_provider_secret":    "The secret to use for authenticating with the Discord API",

        // Server
        "discord_resource_server_id":                            "ID of the server",
        "discord_resource_server_name":                          "What to name the new server",
        "discord_resource_server_region":                        "Guild Voice Region",
        "discord_resource_server_verification_level":            "Verification Level: 0 - None, 1 - Low, 2 - Medium, 3 - High, 4 - Very High",
        "discord_resource_server_explicit_content_filter":       "Explicit Media Content Filter Level: 0 - None, 1 - Non-Roled Members, 2 - All Members",
        "discord_resource_server_default_message_notifications": "Message Notification Level: 0 - All Messages, 1 - Only Mentions",
        "discord_resource_server_afk_channel_id":                "ID For AFK Channel",
        "discord_resource_server_afk_timeout":                   "AFK Timeout in Seconds",
        "discord_resource_server_icon_url":                      "URL to image for the server icon",
        "discord_resource_server_icon_data_uri":                 "Base 64 encoded Data URI for server icon",
        "discord_resource_server_icon_hash":                     "Server icon hash",
        "discord_resource_server_splash_url":                    "URL to image for the server splash",
        "discord_resource_server_splash_data_uri":               "Base 64 encoded Data URI for server splash",
        "discord_resource_server_splash_hash":                   "Server splash hash",
        "discord_resource_server_owner_id":                      "ID of user to own the server",
        "discord_resource_server_system_channel_id":             "ID of channel to send system messages to",

        // Channel
        "discord_resource_channel_server":    "What server this channel belongs to",
        "discord_resource_channel_name":      "What to name the new channel",
        "discord_resource_channel_type":      "What type of channel is this (text, voice, category)",
        "discord_resource_channel_category":  "What category this belongs to, if any",
        "discord_resource_channel_topic":     "Channel Topic",
        "discord_resource_channel_nsfw":      "If the channel is NSFW",
        "discord_resource_channel_position":  "The position of the channel in the left-hand listing",
        "discord_resource_channel_bitrate":   "The bitrate (in bits) of the voice channel; 8000 to 96000 (128000 for VIP servers)",
        "discord_resource_channel_userlimit": "The maximum amount of users that can be in this voice channel",

        // Invite
        "discord_resource_invite_channel":   "What channel this invite belongs to",
        "discord_resource_invite_max_age":   "Duration of invite in seconds before expiry, or 0 for never",
        "discord_resource_invite_max_uses":  "Max number of uses or 0 for unlimited",
        "discord_resource_invite_temporary": "Whether this invite only grants temporary membership",
        "discord_resource_invite_unique":    "if true, don't try to reuse a similar invite (useful for creating many unique one time use invites)",

        // Role
        "discord_resource_role_server":      "ID of the server",
        "discord_resource_role_name":        "Name of the role",
        "discord_resource_role_permissions": "Bitwise value of the enabled/disabled permissions",
        "discord_resource_role_color":       "RGB color value",
        "discord_resource_role_hoist":       "Whether the role should be displayed separately in the sidebar",
        "discord_resource_role_mentionable": "Whether the role should be mentionable",
        "discord_resource_role_members":     "Array of member ids that should have this role",
        "discord_resource_role_position":    "Where this role is positioned. If not set, becomes user managed",
        "discord_resource_role_managed":     "Whether or not the role is managed",

        // Guild Member
        "discord_resource_server_member_server":        "ID of the server",
        "discord_resource_server_member_user":          "ID of the user",
        "discord_resource_server_member_access_token":  fmt.Sprintf("User's access token. Fetch from %s", oauthUrl),
        "discord_resource_server_member_nick":          "Nickname of the user",
        "discord_resource_server_member_mute":          "whether the user is muted in voice channels",
        "discord_resource_server_member_deaf":          "whether the user is deafened in voice channels",
        "discord_resource_server_member_joined_at":     "When the user joined",
        "discord_resource_server_member_premium_since": "When the user started boosting the server",
        "discord_resource_server_member_roles":         "Array of role ids",
    }
}
