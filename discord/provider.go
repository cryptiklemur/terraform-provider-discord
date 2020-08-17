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
                Type:        schema.TypeString,
                Required:    true,
                Description: descriptions["discord_provider_token"],
            },
            "client_id": {
                Type:        schema.TypeString,
                Required:    true,
                Description: descriptions["discord_provider_client_id"],
            },
            "secret": {
                Type:        schema.TypeString,
                Required:    true,
                Description: descriptions["discord_provider_secret"],
            },
        },

        ResourcesMap: map[string]*schema.Resource{
            "discord_server":       resourceDiscordServer(),
            "discord_channel":      resourceDiscordChannel(),
            "discord_invite":       resourceDiscordInvite(),
            "discord_role":         resourceDiscordRole(),
            "discord_server_member": resourceDiscordServerMember(),
        },

        DataSourcesMap: map[string]*schema.Resource{
            "discord_local_image": dataSourceDiscordLocalImage(),
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
