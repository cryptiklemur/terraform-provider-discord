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
			"type": {
				Type:     schema.TypeString,
				Default:  "text",
				Optional: true,
				ForceNew: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errors []error) {
					v := val.(string)
					if v != "text" && v != "voice" && v != "category" {
						errors = append(errors, fmt.Errorf("%q must be one of: text, voice, category, got: %d", key, v))
					}

					return
				},
				Description: descriptions["discord_resource_channel_type"],
			},
		},
	}
}

func resourceChannelCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*discordgo.Session)

	serverId := d.Get("server_id").(string)
	server, err := client.Guild(serverId)
	if err != nil {
		return errors.New("Guild does not exist with that ID: " + serverId)
	}

	name := d.Get("name").(string)
	channelType := d.Get("type").(string)
	channel, err := client.GuildChannelCreate(server.ID, name, channelType)
	if err != nil {
		return errors.New("Failed to create a channel: " + err.Error())
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

	return nil
}

func resourceChannelUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)
	client := m.(*discordgo.Session)

	if d.HasChange("name") {
		_, err := client.ChannelEdit(d.Id(), d.Get("name").(string))
		if err != nil {
			return err
		}

		d.SetPartial("name")
	}

	d.Partial(false)

	return nil
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
