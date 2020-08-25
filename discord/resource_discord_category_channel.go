package discord

import (
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDiscordCategoryChannel() *schema.Resource {
    return &schema.Resource{
        CreateContext: resourceChannelCreate,
        ReadContext:   resourceChannelRead,
        UpdateContext: resourceChannelUpdate,
        DeleteContext: resourceChannelDelete,
        Importer: &schema.ResourceImporter{
            StateContext: schema.ImportStatePassthroughContext,
        },
        Schema: getChannelSchema("category", nil),
    }
}
