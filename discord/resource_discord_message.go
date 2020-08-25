package discord

import (
    "github.com/andersfylling/disgord"
    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "golang.org/x/net/context"
    "log"
    "os"
    "strings"
    "time"
)

func resourceDiscordMessage() *schema.Resource {
    return &schema.Resource{
        CreateContext: resourceMessageCreate,
        ReadContext:   resourceMessageRead,
        UpdateContext: resourceMessageUpdate,
        DeleteContext: resourceMessageDelete,
        Importer: &schema.ResourceImporter{
            StateContext: schema.ImportStatePassthroughContext,
        },

        Schema: map[string]*schema.Schema{
            "channel_id": {
                Type:     schema.TypeString,
                Required: true,
                ForceNew: true,
            },
            "server_id": {
                Type:     schema.TypeString,
                Computed: true,
            },
            "author": {
                Type:     schema.TypeString,
                Computed: true,
            },
            "content": {
                AtLeastOneOf: []string{"content", "embed"},
                Type:         schema.TypeString,
                Optional:     true,
                DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
                    return old == strings.TrimSuffix(new, "\r\n")
                },
            },
            "timestamp": {
                Type:     schema.TypeString,
                Computed: true,
            },
            "edited_timestamp": {
                Type:     schema.TypeString,
                Computed: true,
                Optional: true,
            },
            "tts": {
                Type:     schema.TypeBool,
                Optional: true,
                Default:  false,
            },
            "embed": {
                AtLeastOneOf: []string{"content", "embed"},
                Type:         schema.TypeList,
                Optional:     true,
                MaxItems:     1,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "title": {
                            Type:     schema.TypeString,
                            Optional: true,
                        },
                        "description": {
                            Type:     schema.TypeString,
                            Optional: true,
                        },
                        "url": {
                            Type:     schema.TypeString,
                            Optional: true,
                        },
                        "timestamp": {
                            Type:     schema.TypeString,
                            Optional: true,
                        },
                        "color": {
                            Type:     schema.TypeInt,
                            Optional: true,
                        },
                        "footer": {
                            Type:     schema.TypeList,
                            Optional: true,
                            MaxItems: 1,
                            Elem: &schema.Resource{
                                Schema: map[string]*schema.Schema{
                                    "text": {
                                        Type:     schema.TypeString,
                                        Required: true,
                                    },
                                    "icon_url": {
                                        Type:     schema.TypeString,
                                        Optional: true,
                                    },
                                },
                            },
                        },
                        "image": {
                            Type:     schema.TypeList,
                            Optional: true,
                            MaxItems: 1,
                            Elem: &schema.Resource{
                                Schema: map[string]*schema.Schema{
                                    "url": {
                                        Type:     schema.TypeString,
                                        Required: true,
                                    },
                                    "proxy_url": {
                                        Type:     schema.TypeString,
                                        Computed: true,
                                    },
                                    "height": {
                                        Type:     schema.TypeInt,
                                        Optional: true,
                                    },
                                    "width": {
                                        Type:     schema.TypeInt,
                                        Optional: true,
                                    },
                                },
                            },
                        },
                        "thumbnail": {
                            Type:     schema.TypeList,
                            Optional: true,
                            MaxItems: 1,
                            Elem: &schema.Resource{
                                Schema: map[string]*schema.Schema{
                                    "url": {
                                        Type:     schema.TypeString,
                                        Required: true,
                                    },
                                    "proxy_url": {
                                        Type:     schema.TypeString,
                                        Computed: true,
                                    },
                                    "height": {
                                        Type:     schema.TypeInt,
                                        Optional: true,
                                    },
                                    "width": {
                                        Type:     schema.TypeInt,
                                        Optional: true,
                                    },
                                },
                            },
                        },
                        "video": {
                            Type:     schema.TypeList,
                            Optional: true,
                            MaxItems: 1,
                            Elem: &schema.Resource{
                                Schema: map[string]*schema.Schema{
                                    "url": {
                                        Type:     schema.TypeString,
                                        Required: true,
                                    },
                                    "height": {
                                        Type:     schema.TypeInt,
                                        Optional: true,
                                    },
                                    "width": {
                                        Type:     schema.TypeInt,
                                        Optional: true,
                                    },
                                },
                            },
                        },
                        "provider": {
                            Type:     schema.TypeList,
                            Optional: true,
                            MaxItems: 1,
                            Elem: &schema.Resource{
                                Schema: map[string]*schema.Schema{
                                    "name": {
                                        Type:     schema.TypeString,
                                        Optional: true,
                                    },
                                    "url": {
                                        Type:     schema.TypeString,
                                        Optional: true,
                                    },
                                },
                            },
                        },
                        "author": {
                            Type:     schema.TypeList,
                            Optional: true,
                            MaxItems: 1,
                            Elem: &schema.Resource{
                                Schema: map[string]*schema.Schema{
                                    "name": {
                                        Type:     schema.TypeString,
                                        Optional: true,
                                    },
                                    "url": {
                                        Type:     schema.TypeString,
                                        Optional: true,
                                    },
                                    "icon_url": {
                                        Type:     schema.TypeString,
                                        Optional: true,
                                    },
                                    "proxy_icon_url": {
                                        Type:     schema.TypeString,
                                        Computed: true,
                                    },
                                },
                            },
                        },
                        "fields": {
                            Type:     schema.TypeList,
                            Optional: true,
                            Elem: &schema.Resource{
                                Schema: map[string]*schema.Schema{
                                    "name": {
                                        Type:     schema.TypeString,
                                        Required: true,
                                    },
                                    "value": {
                                        Type:     schema.TypeString,
                                        Optional: true,
                                    },
                                    "inline": {
                                        Type:     schema.TypeBool,
                                        Optional: true,
                                    },
                                },
                            },
                        },
                    },
                },
            },
            "pinned": {
                Type:     schema.TypeBool,
                Optional: true,
                Default:  false,
            },
            "type": {
                Type:     schema.TypeInt,
                Computed: true,
            },
        },
    }
}

func resourceMessageCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    client := m.(*Context).Client

    channelId := getId(d.Get("channel_id").(string))
    params := &disgord.CreateMessageParams{
        Content: d.Get("content").(string),
        Tts:     d.Get("tts").(bool),
    }

    if v, ok := d.GetOk("embed"); ok {
        embed, err := buildEmbed(v.([]interface{}))
        if err != nil {
            return diag.Errorf("Failed to create message in %s: %s", channelId.String(), err.Error())
        }

        params.Embed = embed
    }

    message, err := client.CreateMessage(ctx, channelId, params)
    if err != nil {
        return diag.Errorf("Failed to create message in %s: %s", channelId.String(), err.Error())
    }

    d.SetId(message.ID.String())
    d.Set("type", int(message.Type))
    d.Set("timestamp", message.Timestamp.Format(time.RFC3339))
    d.Set("author", message.Author.ID.String())
    if len(message.Embeds) > 0 {
        d.Set("embed", unbuildEmbed(message.Embeds[0]))
    } else {
        d.Set("embed", nil)
    }
    if !message.GuildID.IsZero() {
        d.Set("server_id", message.GuildID.String())
    }

    if d.Get("pinned").(bool) {
        err = client.PinMessage(ctx, message)
        if err != nil {
            diags = append(diags, diag.Errorf("Failed to pin message %s in %s: %s", message.ID.String(), channelId.String(), err.Error())...)
        }
    }

    return diags
}

func resourceMessageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    client := m.(*Context).Client

    channelId := getId(d.Get("channel_id").(string))
    messageId := getId(d.Id())
    message, err := client.GetMessage(ctx, channelId, messageId)
    if err != nil {
        return diag.Errorf("Failed to fetch message %s in %s: %s", messageId.String(), channelId.String(), err.Error())
    }

    if !message.GuildID.IsZero() {
        d.Set("server_id", message.GuildID.String())
    }
    d.Set("type", int(message.Type))
    d.Set("tts", message.Tts)
    d.Set("timestamp", message.Timestamp.Format(time.RFC3339))
    d.Set("author", message.Author.ID.String())
    d.Set("content", message.Content)
    d.Set("pinned", message.Pinned)

    if len(message.Embeds) > 0 {
        d.Set("embed", unbuildEmbed(message.Embeds[0]))
    }
    d.Set("edited_timestamp", message.EditedTimestamp.Format(time.RFC3339))

    return diags
}

func resourceMessageUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    f, err := os.OpenFile("text.log",
        os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Println(err)
    }
    defer f.Close()

    logger := log.New(f, "", log.LstdFlags)

    var diags diag.Diagnostics
    client := m.(*Context).Client

    channelId := getId(d.Get("channel_id").(string))
    messageId := getId(d.Id())
    builder := client.UpdateMessage(ctx, channelId, messageId)

    if d.HasChange("content") {
        message, _ := client.GetMessage(ctx, channelId, messageId)
        logger.Printf("old: %#v\n\n", message.Content)
        logger.Printf("new: %#v\n\n", d.Get("content").(string))
        builder.SetContent(d.Get("content").(string))
    }
    if d.HasChange("embed") {
        var embed *disgord.Embed
        _, n := d.GetChange("embed")
        if len(n.([]interface{})) > 0 {
            e, err := buildEmbed(n.([]interface{}))
            if err != nil {
                return diag.Errorf("Failed to edit message %s in %s: %s", messageId.String(), channelId.String(), err.Error())
            }

            embed = e
        }

        builder.SetEmbed(embed)
    }

    message, err := builder.Execute()
    if err != nil {
        return diag.Errorf("Failed to update message %s in %s: %s", channelId.String(), messageId.String(), err.Error())
    }

    if len(message.Embeds) > 0 {
        d.Set("embed", unbuildEmbed(message.Embeds[0]))
    } else {
        d.Set("embed", nil)
    }

    d.Set("edited_timestamp", message.EditedTimestamp.Format(time.RFC3339))

    return diags
}

func resourceMessageDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    client := m.(*Context).Client

    channelId := getId(d.Get("channel_id").(string))
    messageId := getId(d.Id())
    err := client.DeleteMessage(ctx, channelId, messageId)
    if err != nil {
        return diag.Errorf("Failed to delete message %s in %s: %s", messageId.String(), channelId.String(), err.Error())
    }

    return diags
}
