package discord

import (
    "errors"
    "fmt"
    "github.com/andersfylling/disgord"
    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "golang.org/x/net/context"
    "strings"
)

func resourceDiscordChannel() *schema.Resource {
    return &schema.Resource{
        CreateContext: resourceChannelCreate,
        ReadContext:   resourceChannelRead,
        UpdateContext: resourceChannelUpdate,
        DeleteContext: resourceChannelDelete,
        Importer: &schema.ResourceImporter{
            StateContext: schema.ImportStatePassthroughContext,
        },

        Schema: map[string]*schema.Schema{
            "channel_id": {
                Type:        schema.TypeString,
                Computed:    true,
                Description: descriptions["discord_resource_channel_id"],
            },
            "name": {
                Type:        schema.TypeString,
                Required:    true,
                Description: descriptions["discord_resource_channel_name"],
            },
            "server_id": {
                Type:        schema.TypeString,
                Required:    true,
                Description: descriptions["discord_resource_channel_server"],
            },
            "category": {
                Type:        schema.TypeString,
                Optional:    true,
                Description: descriptions["discord_resource_channel_category"],
            },
            "type": {
                Type:     schema.TypeString,
                Default:  "text",
                Optional: true,
                ForceNew: true,
                ValidateFunc: func(val interface{}, key string) (warns []string, errors []error) {
                    v := val.(string)

                    if _, ok := getDiscordChannelType(v); !ok {
                        errors = append(errors, fmt.Errorf("%q must be one of: text, voice, category, got: %d", key, v))
                    }

                    return
                },
                Description: descriptions["discord_resource_channel_type"],
            },
            "topic": {
                Type:        schema.TypeString,
                Optional:    true,
                Default:     "",
                Description: descriptions["discord_resource_channel_topic"],
            },
            "nsfw": {
                Type:        schema.TypeBool,
                Optional:    true,
                Description: descriptions["discord_resource_channel_nsfw"],
            },
            "position": {
                Type:        schema.TypeInt,
                Default:     1,
                Optional:    true,
                Description: descriptions["discord_resource_channel_position"],
            },
            "bitrate": {
                Type:        schema.TypeInt,
                Optional:    true,
                Description: descriptions["discord_resource_channel_bitrate"],
            },
            "user_limit": {
                Type:        schema.TypeInt,
                Optional:    true,
                Description: descriptions["discord_resource_channel_userlimit"],
            },
        },
    }
}

func validateChannel(d *schema.ResourceData) (bool, error) {
    channelType := d.Get("type").(string)

    if channelType == "category" {
        if _, ok := d.GetOk("category"); ok {
            return false, errors.New("category cannot be a child of another category")
        }
        if _, ok := d.GetOk("nsfw"); ok {
            return false, errors.New("nsfw is not allowed on categories")
        }
    }

    if channelType == "voice" {
        if _, ok := d.GetOk("topic"); ok {
            return false, errors.New("topic is not allowed on voice channels")
        }
        if _, ok := d.GetOk("nsfw"); ok {
            return false, errors.New("nsfw is not allowed on voice channels")
        }
    }

    if channelType == "text" {
        if _, ok := d.GetOk("bitrate"); ok {
            return false, errors.New("bitrate is not allowed on text channels")
        }
        if _, ok := d.GetOk("user_limit"); ok {
            if d.Get("user_limit").(int) > 0 {
                return false, errors.New("user_limit is not allowed on text channels")
            }
        }
        name := d.Get("name").(string)
        if strings.ToLower(name) != name {
            return false, errors.New("name must be lowercase")
        }
    }

    return true, nil
}

func resourceChannelCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    client := m.(*Context).Client

    if ok, reason := validateChannel(d); !ok {
        return diag.FromErr(reason)
    }

    serverId := getMajorId(d.Get("server_id"))
    channelType, _ := getDiscordChannelType(d.Get("type").(string))

    channel, err := client.CreateGuildChannel(ctx, serverId, d.Get("name").(string), &disgord.CreateGuildChannelParams{
        Type:      channelType,
        Topic:     d.Get("topic").(string),
        Bitrate:   d.Get("bitrate").(uint),
        UserLimit: d.Get("user_limit").(uint),
        ParentID:  getId(d.Get("category").(string)),
        NSFW:      d.Get("nsfw").(bool),
        Position:  d.Get("position").(int),
    })

    if err != nil {
        return diag.Errorf("Failed to create channel: %s", err.Error())
    }

    d.SetId(channel.ID.String())
    d.Set("server_id", serverId)
    d.Set("channel_id", channel.ID.String())

    return diags
}

func resourceChannelRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    client := m.(*Context).Client

    channel, err := client.GetChannel(ctx, getId(d.Id()))
    if err != nil {
        return diag.Errorf("Failed to fetch channel %s: %e", d.Id(), err.Error())
    }

    channelType, ok := getTextChannelType(channel.Type)
    if !ok {
        return diag.Errorf("Invalid channel type: %d", channel.Type)
    }

    d.Set("type", channelType)
    d.Set("name", channel.Name)
    d.Set("position", channel.Position)

    if channelType == "text" {
        d.Set("topic", channel.Topic)
        d.Set("nsfw", channel.NSFW)
    }

    if channelType == "voice" {
        d.Set("bitrate", channel.Bitrate)
    }

    if !channel.ParentID.IsZero() {
        d.Set("category", channel.ParentID.String())
    } else {
        d.Set("category", nil)
    }

    return diags
}

func resourceChannelUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    client := m.(*Context).Client
    if ok, reason := validateChannel(d); !ok {
        return diag.FromErr(reason)
    }

    channelType := d.Get("type").(string)
    builder := client.UpdateChannel(ctx, getId(d.Id()))

    if d.HasChange("name") {
        builder.SetName(d.Get("name").(string))
    }
    if d.HasChange("topic") && channelType == "text" {
        builder.SetTopic(d.Get("topic").(string))
    }
    if d.HasChange("nsfw") && channelType == "text" {
        builder.SetNsfw(d.Get("nsfw").(bool))
    }
    if _, v := d.GetChange("bitrate"); v.(int) > 0 && channelType == "voice" {
        builder.SetBitrate(uint(d.Get("bitrate").(int)))
    }
    if d.HasChange("user_limit") && channelType == "voice" {
        builder.SetUserLimit(uint(d.Get("user_limit").(int)))
    }
    if d.HasChange("position") {
        builder.SetPosition(d.Get("position").(int))
    }
    if d.HasChange("category") {
        if d.Get("category").(string) != "" {
            builder.SetParentID(getId(d.Get("category").(string)))
        } else {
            builder.RemoveParentID()
        }
    }

    _, err := builder.Execute()
    if err != nil {
        return diag.Errorf("Failed to update channel %s: %s", d.Id(), err.Error())
    }

    return diags
}

func resourceChannelDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    client := m.(*Context).Client

    _, err := client.DeleteChannel(ctx, getId(d.Id()))
    if err != nil {
        return diag.Errorf("Failed to delete channel %s: %e", d.Id(), err.Error())
    }

    return diags
}
