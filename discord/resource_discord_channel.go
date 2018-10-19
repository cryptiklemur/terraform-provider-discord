package discord

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDiscordChannel() *schema.Resource {
	return &schema.Resource{
		Create: resourceChannelCreate,
		Read:   resourceChannelRead,
		Update: resourceChannelUpdate,
		Delete: resourceChannelDelete,

		Schema: map[string]*schema.Schema{
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
				Type: schema.TypeString,
				Optional: true,
				Description: descriptions["discord_resource_channel_category"],
			},
			"type": {
				Type:     schema.TypeString,
				Default:  "text",
				Optional: true,
				ForceNew: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errors []error) {
					v := val.(string)

					if _, ok := getChannelType(v); !ok {
						errors = append(errors, fmt.Errorf("%q must be one of: text, voice, category, got: %d", key, v))
					}

					return
				},
				Description: descriptions["discord_resource_channel_type"],
			},
			"topic": {
				Type:     schema.TypeString,
				Optional: true,
				Description: descriptions["discord_resource_channel_topic"],
			},
			"nsfw": {
				Type:     schema.TypeBool,
				Optional: true,
				Description: descriptions["discord_resource_channel_nsfw"],
			},
			"position": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: descriptions["discord_resource_channel_position"],
			},
			"bitrate": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: descriptions["discord_resource_channel_bitrate"],
			},
			"userlimit": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: descriptions["discord_resource_channel_userlimit"],
			},
		},
	}
}

func getChannelType(name string) (string, bool) {
	switch name {
	case "text":
	case "voice":
		return name, true
	case "category":
		return "4", true
	}

	return name, false
}

func validateChannel(d *schema.ResourceData) (bool, error) {
	channelType := d.Get("type").(string)

	if channelType == "category" {
		if _, ok := d.GetOkExists("category"); ok {
			return false, errors.New("category cannot be a child of another category")
		}
		if _, ok := d.GetOkExists("nsfw"); ok {
			return false, errors.New("nsfw is not allowed on categories")
		}
	}

	if channelType == "voice" {
		if _, ok := d.GetOkExists("topic"); ok {
			return false, errors.New("topic is not allowed on voice channels")
		}
		if _, ok := d.GetOkExists("nsfw"); ok {
			return false, errors.New("nsfw is not allowed on voice channels")
		}
	}

	if channelType == "text" {
		if _, ok := d.GetOkExists("bitrate"); ok {
			return false, errors.New("bitrate is not allowed on text channels")
		}
		if _, ok := d.GetOkExists("user_limit"); ok {
			return false, errors.New("user_limit is not allowed on text channels")
		}
	}

	return true, nil
}

func resourceChannelCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*discordgo.Session)

	if ok, reason := validateChannel(d); !ok {
		return reason
	}

	serverId := d.Get("server_id").(string)
	server, err := client.Guild(serverId)
	if err != nil {
		return errors.New("Guild does not exist with that ID: " + serverId)
	}

	name := d.Get("name").(string)
	channelType, _ := getChannelType(d.Get("type").(string))
	channel, err := client.GuildChannelCreate(server.ID, name, channelType)
	if err != nil {
		return errors.New("Failed to create a channel: " + err.Error())
	}

	params := discordgo.ChannelEdit{}
	edit := false
	if v, ok := d.GetOkExists("topic"); ok {
		params.Topic = v.(string)
		edit = true
	}
	if v, ok := d.GetOkExists("nsfw"); ok {
		params.NSFW = v.(bool)
		edit = true
	}
	if v, ok := d.GetOkExists("position"); ok {
		params.Position = v.(int)
		edit = true
	}
	if v, ok := d.GetOkExists("bitrate"); ok {
		params.Bitrate = v.(int)
		edit = true
	}
	if v, ok := d.GetOkExists("user_limit"); ok {
		params.UserLimit = v.(int)
		edit = true
	}
	if v, ok := d.GetOkExists("category"); ok {
		params.ParentID = v.(string)
		edit = true
	}

	if edit {
		client.ChannelEditComplex(channel.ID, &params)
	}

	d.SetId(channel.ID)

	return nil
}

func resourceChannelRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*discordgo.Session)

	channel, err := client.Channel(d.Id())
	if err != nil {
		d.SetId("")

		return nil
	}

	d.Set("type", channel.Type)
	d.Set("name", channel.Name)
	d.Set("position", channel.Position)
	d.Set("category", channel.ParentID)
	if channel.Type == discordgo.ChannelTypeGuildVoice {
		d.Set("bitrate", channel.Bitrate)
	}
	if channel.Type == discordgo.ChannelTypeGuildText {
		d.Set("topic", channel.Topic)
		d.Set("nsfw", channel.NSFW)
	}

	return nil
}

func resourceChannelUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*discordgo.Session)
	if ok, reason := validateChannel(d); !ok {
		return reason
	}

	params := discordgo.ChannelEdit{}
	changed := false

	if d.HasChange("name") {
		params.Name = d.Get("name").(string)
		changed = true
	}
	if d.HasChange("topic") {
		params.Topic = d.Get("topic").(string)
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
		_, err := client.ChannelEditComplex(d.Id(), &params)
		if err != nil {
			return err
		}
	}

	return resourceChannelRead(d, m)
}

func resourceChannelDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*discordgo.Session)

	_, err := client.Channel(d.Id())
	if err != nil {
		return nil
	}

	_, _ = client.ChannelDelete(d.Id())

	return nil
}
