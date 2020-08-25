package discord

import (
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDiscordVoiceChannel() *schema.Resource {
    return &schema.Resource{
        CreateContext: resourceChannelCreate,
        ReadContext:   resourceChannelRead,
        UpdateContext: resourceChannelUpdate,
        DeleteContext: resourceChannelDelete,
        Importer: &schema.ResourceImporter{
            StateContext: schema.ImportStatePassthroughContext,
        },
        Schema: getChannelSchema("voice", map[string]*schema.Schema{
            "bitrate": {
                Type:     schema.TypeInt,
                Optional: true,
                Default:  64000,
            },
            "user_limit": {
                Type:     schema.TypeInt,
                Optional: true,
            },
        }),
    }
}
