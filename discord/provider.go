package discord

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"secret": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"discord_server":             resourceDiscordServer(),
			"discord_category_channel":   resourceDiscordCategoryChannel(),
			"discord_text_channel":       resourceDiscordTextChannel(),
			"discord_voice_channel":      resourceDiscordVoiceChannel(),
			"discord_news_channel":       resourceDiscordNewsChannel(),
			"discord_channel_permission": resourceDiscordChannelPermission(),
			"discord_invite":             resourceDiscordInvite(),
			"discord_role":               resourceDiscordRole(),
			"discord_role_everyone":      resourceDiscordRoleEveryone(),
			"discord_member_roles":       resourceDiscordMemberRoles(),
			"discord_message":            resourceDiscordMessage(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"discord_permission":  dataSourceDiscordPermission(),
			"discord_color":       dataSourceDiscordColor(),
			"discord_local_image": dataSourceDiscordLocalImage(),
			"discord_role":        dataSourceDiscordRole(),
			"discord_server":      dataSourceDiscordServer(),
			"discord_member":      dataSourceDiscordMember(),
		},

		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	config := Config{
		Token: d.Get("token").(string),
	}

	client, err := config.Client()
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return client, diags
}
