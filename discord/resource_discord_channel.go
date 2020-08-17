package discord

import (
    "errors"
    "fmt"
    "github.com/bwmarrin/discordgo"
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
            StateContext: resourceChannelImportState,
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

                    if _, ok := getDiscordGoChannelType(v); !ok {
                        errors = append(errors, fmt.Errorf("%q must be one of: text, voice, category, got: %d", key, v))
                    }

                    return
                },
                Description: descriptions["discord_resource_channel_type"],
            },
            "topic": {
                Type:        schema.TypeString,
                Optional:    true,
                Default:     "​",
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

func getTextChannelType(channelType discordgo.ChannelType) (string, bool) {
    switch channelType {
    case discordgo.ChannelTypeGuildText:
        return "text", true
    case discordgo.ChannelTypeGuildVoice:
        return "voice", true
    case discordgo.ChannelTypeGuildCategory:
        return "category", true
    case discordgo.ChannelTypeGuildNews:
        return "news", true
    case discordgo.ChannelTypeGuildStore:
        return "store", true
    }

    return "text", false
}

func getDiscordGoChannelType(name string) (discordgo.ChannelType, bool) {
    switch name {
    case "text":
        return discordgo.ChannelTypeGuildText, true
    case "voice":
        return discordgo.ChannelTypeGuildVoice, true
    case "category":
        return discordgo.ChannelTypeGuildCategory, true
    case "news":
        return discordgo.ChannelTypeGuildNews, true
    case "store":
        return discordgo.ChannelTypeGuildStore, true
    }

    return -1, false
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

func resourceChannelCreate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    client := m.(*Context).Client

    if ok, reason := validateChannel(d); !ok {
        return diag.FromErr(reason)
    }

    serverId := d.Get("server_id").(string)
    server, err := client.Guild(serverId)
    if err != nil {
        return diag.Errorf("Server does not exist with that ID: %s", serverId)
    }

    channelType, _ := getDiscordGoChannelType(d.Get("type").(string))
    channel, err := client.GuildChannelCreate(
        server.ID,
        d.Get("name").(string),
        channelType,
    )
    if err != nil {
        return diag.Errorf("Failed to create a channel: %s", err.Error())
    }

    params := discordgo.ChannelEdit{}
    edit := false
    if v, ok := d.GetOk("topic"); ok {
        params.Topic = v.(string)
        edit = true
    }
    if v, ok := d.GetOk("nsfw"); ok {
        params.NSFW = v.(bool)
        edit = true
    }
    if v, ok := d.GetOk("position"); ok {
        params.Position = v.(int)
        edit = true
    }
    if v, ok := d.GetOk("bitrate"); ok {
        params.Bitrate = v.(int)
        edit = true
    }
    if v, ok := d.GetOk("user_limit"); ok {
        params.UserLimit = v.(int)
        edit = true
    }
    if v, ok := d.GetOk("category"); ok {
        _, parentId, err := parseTwoIds(v.(string))
        if err != nil {
            return diag.FromErr(err)
        }
        params.ParentID = parentId
        edit = true
    }

    if edit {
        channel, err = client.ChannelEditComplex(channel.ID, &params)
        if err != nil {
            return diag.FromErr(err)
        }
    }

    d.SetId(fmt.Sprintf("%s:%s", serverId, channel.ID))
    d.Set("server_id", serverId)
    d.Set("channel_id", channel.ID)

    return diags
}

func resourceChannelRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    client := m.(*Context).Client

    c, err := getChannel(client, d.Id())
    if err != nil {
        return diag.Errorf("Failed to fetch channel %s: %e", d.Id(), err.Error())
    }
    channel := c.Channel

    channelType, ok := getTextChannelType(channel.Type)
    if !ok {
        return diag.Errorf("Invalid channel type: %d", channel.Type)
    }
    d.Set("type", channelType)

    d.Set("name", channel.Name)
    d.Set("position", channel.Position)
    d.Set("category", fmt.Sprintf("%s:%s", c.ServerId, channel.ParentID))
    if channel.Type == discordgo.ChannelTypeGuildVoice {
        d.Set("bitrate", channel.Bitrate)
    }
    if channel.Type == discordgo.ChannelTypeGuildText {
        d.Set("topic", channel.Topic)
        d.Set("nsfw", channel.NSFW)
    }

    return diags
}

func resourceChannelUpdate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    client := m.(*Context).Client
    if ok, reason := validateChannel(d); !ok {
        return diag.FromErr(reason)
    }

    c, err := getChannel(client, d.Id())
    if err != nil {
        return diag.Errorf("Failed to fetch channel %s: %e", d.Id(), err.Error())
    }
    channel := c.Channel

    params := &discordgo.ChannelEdit{
        Name:                 channel.Name,
        Topic:                channel.Topic,
        NSFW:                 channel.NSFW,
        Position:             channel.Position,
        Bitrate:              channel.Bitrate,
        UserLimit:            channel.UserLimit,
        PermissionOverwrites: channel.PermissionOverwrites,
        ParentID:             channel.ParentID,
        RateLimitPerUser:     channel.RateLimitPerUser,
    }
    changed := false

    if d.HasChange("name") {
        params.Name = d.Get("name").(string)
        changed = true
    }
    if d.HasChange("topic") {
        params.Topic = d.Get("topic").(string)
        if params.Topic == "" {
            params.Topic = "​"
        }

        changed = true
    }
    if d.HasChange("nsfw") {
        params.NSFW = d.Get("nsfw").(bool)
        changed = true
    }
    if d.HasChange("bitrate") {
        params.Bitrate = d.Get("bitrate").(int)
        changed = true
    }
    if d.HasChange("user_limit") {
        params.UserLimit = d.Get("user_limit").(int)
        changed = true
    }
    if d.HasChange("position") {
        params.Position = d.Get("position").(int)
        changed = true
    }
    if d.HasChange("category") {
        params.ParentID = d.Get("category").(string)
        changed = true
    }

    if changed {
        _, err := client.ChannelEditComplex(c.ChannelId, params)
        if err != nil {
            return diag.FromErr(err)
        }
    }

    return diags
}

func resourceChannelDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    client := m.(*Context).Client

    c, err := getChannel(client, d.Id())
    if err != nil {
        return diag.Errorf("Failed to fetch channel %s: %e", d.Id(), err.Error())
    }

    client.ChannelDelete(c.ChannelId)

    return diags
}

func resourceChannelImportState(_ context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
    results := make([]*schema.ResourceData, 1, 1)
    results[0] = d

    client := m.(*Context).Client
    c, err := getChannel(client, d.Id())
    if err != nil {
        return results, err
    }

    pData := resourceDiscordChannel().Data(nil)
    pData.SetId(d.Id())
    pData.SetType("discord_channel")

    channelType, ok := getTextChannelType(c.Channel.Type)
    if !ok {
        return results, errors.New(fmt.Sprint("Invalid channel type: %d", c.Channel.Type))
    }
    d.Set("type", channelType)

    d.Set("server_id", c.Channel.GuildID)
    d.Set("type", c.Channel.Type)
    d.Set("name", c.Channel.Name)
    d.Set("position", c.Channel.Position)
    d.Set("category", c.Channel.ParentID)
    if c.Channel.Type == discordgo.ChannelTypeGuildVoice {
        d.Set("bitrate", c.Channel.Bitrate)
    }
    if c.Channel.Type == discordgo.ChannelTypeGuildText {
        d.Set("topic", c.Channel.Topic)
        d.Set("nsfw", c.Channel.NSFW)
    }
    results = append(results, pData)

    return results, nil
}
