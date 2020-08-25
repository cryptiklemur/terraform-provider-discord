package discord

import (
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDiscordTextChannel() *schema.Resource {
    return &schema.Resource{
        CreateContext: resourceChannelCreate,
        ReadContext:   resourceChannelRead,
        UpdateContext: resourceChannelUpdate,
        DeleteContext: resourceChannelDelete,
        Importer: &schema.ResourceImporter{
            StateContext: schema.ImportStatePassthroughContext,
        },
        Schema: getChannelSchema("text", map[string]*schema.Schema{
            "topic": {
                Type:     schema.TypeString,
                Optional: true,
            },
            "nsfw": {
                Type:     schema.TypeBool,
                Optional: true,
                Default:  false,
            },
        }),
    }
}
