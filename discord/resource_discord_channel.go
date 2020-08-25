package discord

import (
    "errors"
    "github.com/andersfylling/disgord"
    "github.com/hashicorp/go-cty/cty"
    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "golang.org/x/net/context"
    "log"
    "strings"
)

func getChannelSchema(channelType string, s map[string]*schema.Schema) map[string]*schema.Schema {
    addedSchema := map[string]*schema.Schema{
        "server_id": {
            Type:     schema.TypeString,
            Required: true,
        },
        "type": {
            Type:     schema.TypeString,
            Required: true,
            ValidateDiagFunc: func(i interface{}, path cty.Path) (diags diag.Diagnostics) {
                if i.(string) != channelType {
                    diags = append(diags, diag.Errorf("type must be %s, %s passed", channelType, i.(string))...)
                }

                return diags
            },
            DefaultFunc: func() (interface{}, error) {
                return channelType, nil
            },
        },
        "name": {
            Type:     schema.TypeString,
            Required: true,
        },
        "position": {
            Type:     schema.TypeInt,
            Default:  1,
            Optional: true,
        },
    }

    if channelType != "category" {
        addedSchema["category"] = &schema.Schema{
            Type:     schema.TypeString,
            Optional: true,
        }
        addedSchema["sync_perms_with_category"] = &schema.Schema{
            Type:     schema.TypeBool,
            Optional: true,
            Default:  true,
        }
    }

    if s != nil {
        for k, v := range s {
            addedSchema[k] = v
        }
    }

    return addedSchema
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
    channelType := d.Get("type").(string)
    channelTypeInt, _ := getDiscordChannelType(channelType)

    var topic string
    var bitrate uint
    var userlimit uint
    var nsfw bool
    var parentId disgord.Snowflake

    if channelType == "text" {
        if v, ok := d.GetOk("topic"); ok {
            topic = v.(string)
        }
        if v, ok := d.GetOk("nsfw"); ok {
            nsfw = v.(bool)
        }
    } else if channelType == "voice" {
        if v, ok := d.GetOk("bitrate"); ok {
            bitrate = uint(v.(int))
        }
        if v, ok := d.GetOk("userlimit"); ok {
            userlimit = uint(v.(int))
        }
    }

    if channelType != "category" {
        if v, ok := d.GetOk("category"); ok {
            parentId = getId(v.(string))
        }
    }

    channel, err := client.CreateGuildChannel(ctx, serverId, d.Get("name").(string), &disgord.CreateGuildChannelParams{
        Type:     channelTypeInt,
        Topic:     topic,
        Bitrate:   bitrate,
        UserLimit: userlimit,
        ParentID:  parentId,
        NSFW:      nsfw,
        Position:  d.Get("position").(int),
    })

    if err != nil {
        return diag.Errorf("Failed to create channel: %s", err.Error())
    }

    d.SetId(channel.ID.String())
    d.Set("server_id", serverId)
    d.Set("channel_id", channel.ID.String())

    if channelType != "category" {
        if v, ok := d.GetOk("sync_perms_with_category"); ok && v.(bool) {
            if channel.ParentID.IsZero() {
                return append(diags, diag.Errorf("Can't sync permissions with category. Channel (%s) doesn't have a category", channel.ID.String())...)
            }
            parent, err := client.GetChannel(ctx, channel.ParentID)
            if err != nil {
                return append(diags, diag.Errorf("Can't sync permissions with category. Channel (%s) doesn't have a category", channel.ID.String())...)
            }

            if err = syncChannelPermissions(client, ctx, parent, channel); err != nil {
                return append(diags, diag.Errorf("Can't sync permissions with category: %s", channel.ID.String(), err.Error())...)
            }
        }
    }

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
    } else if channelType == "voice" {
        d.Set("bitrate", channel.Bitrate)
        d.Set("userlimit", channel.UserLimit)
    }

    if channelType != "category" {
        if !channel.ParentID.IsZero() {
            parent, err := client.GetChannel(ctx, channel.ParentID)
            if err != nil {
                return diag.Errorf("Failed to fetch category of channel %s: %s", channel.ID.String(), err.Error())
            }

            synced := arePermissionsSynced(channel, parent)
            log.Printf("Are permissions synced between %s and %s: %s", channel.Name, parent.Name, synced)
            d.Set("sync_perms_with_category", synced)
        } else {
            d.Set("sync_perms_with_category", false)
        }
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
    if d.HasChange("position") {
        builder.SetPosition(d.Get("position").(int))
    }

    if channelType == "text" {
        if d.HasChange("topic") {
            builder.SetTopic(d.Get("topic").(string))
        }
        if d.HasChange("nsfw") {
            builder.SetNsfw(d.Get("nsfw").(bool))
        }
    } else if channelType == "voice" {
        if d.HasChange("bitrate") {
            builder.SetBitrate(uint(d.Get("bitrate").(int)))
        }
        if d.HasChange("user_limit") {
            builder.SetUserLimit(uint(d.Get("user_limit").(int)))
        }
    }
    if channelType != "category" && d.HasChange("category") {
        if d.Get("category").(string) != "" {
            builder.SetParentID(getId(d.Get("category").(string)))
        } else {
            builder.RemoveParentID()
        }
    }

    channel, err := builder.Execute()
    if err != nil {
        return diag.Errorf("Failed to update channel %s: %s", d.Id(), err.Error())
    }

    if channelType != "category" {
        if v, ok := d.GetOk("sync_perms_with_category"); ok && v.(bool) {
            if channel.ParentID.IsZero() {
                return append(diags, diag.Errorf("Can't sync permissions with category. Channel (%s) doesn't have a category", channel.ID.String())...)
            }
            parent, err := client.GetChannel(ctx, channel.ParentID)
            if err != nil {
                return append(diags, diag.Errorf("Can't sync permissions with category. Channel (%s) doesn't have a category", channel.ID.String())...)
            }

            if err = syncChannelPermissions(client, ctx, parent, channel); err != nil {
                return append(diags, diag.Errorf("Can't sync permissions with category: %s", channel.ID.String(), err.Error())...)
            }
        }
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
