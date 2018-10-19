package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
	"log"
)

func resourceDiscordServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,
		Importer: &schema.ResourceImporter{
			State: resourceServerImportState,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: descriptions["discord_resource_server_name"],
			},
			"empty": {
				Type: schema.TypeBool,
				Default: true,
				ForceNew: true,
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
		log.Println("DISCORD: Empty marked as true. Wiping server.")
		channels, err := client.GuildChannels(guild.ID)
		if err != nil {
			return errors.New("Failed to fetch channels for new guild")
		}

		for _, channel := range channels {
			log.Println("DISCORD: Deleting Channel: " + channel.ID)
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
		log.Fatal(err)
		d.SetId("")
		return nil
	}


	d.Set("name", guild.Name)

	return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*discordgo.Session)

	if d.HasChange("name") {
		_, err := client.GuildEdit(d.Id(), discordgo.GuildParams{Name: d.Get("name").(string)})
		if err != nil {
			return err
		}
	}


	return resourceServerRead(d, m)
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*discordgo.Session)

	_, err := client.Guild(d.Id())
	if err != nil {
		log.Fatal(err)
		return nil
	}

	client.GuildDelete(d.Id())

	return nil
}

func resourceServerImportState(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	results := make([]*schema.ResourceData, 1, 1)
	results[0] = d

	client := m.(*discordgo.Session)
	guild, err := client.Guild(d.Id())
	if err != nil {
		return results, nil
	}

	server := resourceDiscordServer()
	pData := server.Data(nil)
	pData.SetId(d.Id())
	pData.SetType("discord_server")
	pData.Set("name", guild.Name)
	results = append(results, pData)

	return results, nil
}