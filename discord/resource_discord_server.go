package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDiscordServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: descriptions["discord_resource_server_name"],
			},
			"empty": {
				Type: schema.TypeBool,
				Default: true,
				Optional: true,
				Description: descriptions["discord_resource_server_empty"],
			},
		},
	}
}

func resourceServerCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*discordgo.Session)

	name := d.Get("name").(string)
	guild, err := client.GuildCreate(name)
	if err != nil {
		return err
	}

	if d.Get("empty").(bool) {
		for _, channel := range guild.Channels {
			client.ChannelDelete(channel.ID)
		}
	}

	d.SetId(guild.ID)

	return nil
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*discordgo.Session)

	guild, err := client.Guild(d.Id())
	if err != nil {
		d.SetId("")
		return nil
	}

	d.Set("name", guild.Name)

	return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)
	client := m.(*discordgo.Session)

	if d.HasChange("name") {
		_, err := client.GuildEdit(d.Id(), discordgo.GuildParams{Name: d.Get("name").(string)})
		if err != nil {
			return err
		}

		d.SetPartial("name")
	}

	d.Partial(false)

	return nil
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*discordgo.Session)

	_, err := client.Guild(d.Id())
	if err != nil {
		return nil
	}

	client.GuildDelete(d.Id())

	return nil
}
