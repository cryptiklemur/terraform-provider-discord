package discord

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDiscordInvite() *schema.Resource {
	return &schema.Resource{
		Create: resourceInviteCreate,
		Read:   resourceInviteRead,
		Delete: resourceInviteDelete,

		Schema: map[string]*schema.Schema{
			"channel_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: descriptions["discord_resource_invite_channel"],
			},
			"max_age": {
				Type:        schema.TypeInt,
				ForceNew:    true,
				Optional:    true,
				Default:     86400,
				Description: descriptions["discord_resource_invite_max_age"],
			},
			"max_uses": {
				Type:        schema.TypeInt,
				ForceNew:    true,
				Optional:    true,
				Default:     0,
				Description: descriptions["discord_resource_invite_max_uses"],
			},
			"temporary": {
				Type:        schema.TypeBool,
				ForceNew:    true,
				Optional:    true,
				Default:     false,
				Description: descriptions["discord_resource_invite_temporary"],
			},
			"unique": {
				Type:        schema.TypeBool,
				ForceNew:    true,
				Optional:    true,
				Default:     false,
				Description: descriptions["discord_resource_invite_unique"],
			},
		},
	}
}

func resourceInviteCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*discordgo.Session)

	channelId := d.Get("channel_id").(string)
	channel, err := client.Channel(channelId)
	if err != nil {
		return errors.New("Channel does not exist with that ID: " + channelId)
	}

	invite, err := client.ChannelInviteCreate(channel.ID, discordgo.Invite{
		MaxAge:    d.Get("max_age").(int),
		MaxUses:   d.Get("max_uses").(int),
		Temporary: d.Get("temporary").(bool),
		Unique:    d.Get("unique").(bool),
	})
	if err != nil {
		return errors.New("Failed to create a invite: " + err.Error())
	}

	d.SetId(invite.Code)

	return nil
}

func resourceInviteRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*discordgo.Session)

	_, err := client.Invite(d.Id())
	if err != nil {
		d.SetId("")

		return nil
	}

	return nil
}

func resourceInviteDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*discordgo.Session)

	_, err := client.Invite(d.Id())
	if err != nil {
		return nil
	}

	_, _ = client.InviteDelete(d.Id())

	return nil
}
