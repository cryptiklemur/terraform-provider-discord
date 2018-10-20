package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
	"github.com/polds/imgbase64"
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
			"region": {
				Type: schema.TypeString,
				Optional: true,
				Description: descriptions["discord_resource_server_region"],
			},
			"verification_level": {
				Type: schema.TypeInt,
				Optional: true,
				Default: 0,
				Description: descriptions["verification_level"],
				ValidateFunc: func(val interface{}, key string) (warns []string, errors []error) {
					v := val.(int)
					if v > 3 || v < 0 {
						errors = append(errors, fmt.Errorf("%q must be between 0 and 3 inclusive, got: %d", key, v))
					}

					return
				},
			},
			"default_message_notifications": {
				Type: schema.TypeInt,
				Optional: true,
				Default: 0,
				Description: descriptions["discord_resource_server_default_message_notifications"],
				ValidateFunc: func(val interface{}, key string) (warns []string, errors []error) {
					v := val.(int)
					if v != 0 && v != 1 {
						errors = append(errors, fmt.Errorf("%q must be 0 or 1, got: %d", key, v))
					}

					return
				},
			},
			"afk_channel_id": {
				Type: schema.TypeString,
				Optional: true,
				Description: descriptions["discord_resource_server_afk_channel_id"],
			},
			"afk_timeout": {
				Type: schema.TypeInt,
				Optional: true,
				Default: 300,
				Description: descriptions["discord_resource_server_afk_timeout"],
				ValidateFunc: func(val interface{}, key string) (warns []string, errors []error) {
					v := val.(int)
					if v < 0 {
						errors = append(errors, fmt.Errorf("%q must be greater than 0, got: %d", key, v))
					}

					return
				},
			},
			"icon_url": {
				Type: schema.TypeString,
				Optional: true,
				Description: descriptions["discord_resource_server_icon_url"],
			},
			"icon_local": {
				Type: schema.TypeString,
				Optional: true,
				Description: descriptions["discord_resource_server_icon_url"],
			},
			"icon_data_uri": {
				Type: schema.TypeString,
				Optional: true,
				Description: descriptions["discord_resource_server_icon_data_uri"],
			},
			"icon_hash": {
				Type: schema.TypeString,
				Computed: true,
				Description: descriptions["discord_resource_server_icon_hash"],
			},
			"owner_id": {
				Type: schema.TypeString,
				Optional: true,
				Description: descriptions["discord_resource_server_owner_id"],
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

	level := discordgo.VerificationLevel(d.Get("verification_level").(int))
	params := discordgo.GuildParams{
		DefaultMessageNotifications: d.Get("default_message_notifications").(int),
		VerificationLevel: &level,
	}
	edit := false
	if v, ok := d.GetOkExists("region"); ok {
		params.Region = v.(string)
		edit = true
	}
	if v, ok := d.GetOkExists("afk_channel_id"); ok {
		params.AfkChannelID = v.(string)
		edit = true
	}
	if v, ok := d.GetOkExists("afk_timeout"); ok {
		params.AfkTimeout = v.(int)
		edit = true
	}
	if v, ok := d.GetOkExists("icon_url"); ok {
		img := imgbase64.FromRemote(v.(string))
		params.Icon = img
		edit = true
	}
	if v, ok := d.GetOkExists("icon_local"); ok {
		img, err := imgbase64.FromLocal(v.(string))
		if err != nil {
			client.GuildDelete(guild.ID)

			return errors.New("Failed to fetch icon from: " + v.(string))
		}
		params.Icon = img
		edit = true
	}
	if v, ok := d.GetOkExists("icon_data_uri"); ok {
		params.Icon = v.(string)
		edit = true
	}
	if v, ok := d.GetOkExists("owner_id"); ok {
		params.OwnerID = v.(string)
		edit = true
	}

	log.Println("[DISCORD] Setting icon to: " + params.Icon)
	if edit {
		client.GuildEdit(guild.ID, params)
	}

	d.SetId(guild.ID)
	d.Set("owner", guild.OwnerID)

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
	d.Set("region", guild.Region)
	d.Set("default_message_notifications", guild.DefaultMessageNotifications)
	d.Set("afk_channel_id", guild.AfkChannelID)
	d.Set("afk_timeout", guild.AfkTimeout)
	d.Set("icon_hash", guild.Icon)
	d.Set("owner_id", guild.OwnerID)

	return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*discordgo.Session)

	changed := false
	params := discordgo.GuildParams{}

	if d.HasChange("name") {
		params.Name = d.Get("name").(string)
		changed = true
	}
	if d.HasChange("region") {
		params.Region = d.Get("region").(string)
		changed = true
	}
	if d.HasChange("default_message_notifications") {
		params.DefaultMessageNotifications = d.Get("default_message_notifications").(int)
		changed = true
	}
	if d.HasChange("afk_channel_id") {
		params.AfkChannelID = d.Get("afk_channel_id").(string)
		changed = true
	}
	if d.HasChange("afk_timeout") {
		params.AfkTimeout = d.Get("afk_timeout").(int)
		changed = true
	}
	if d.HasChange("icon_url") {
		img := imgbase64.FromRemote(d.Get("icon_url").(string))
		params.Icon = img
		changed = true
	}
	if d.HasChange("icon_local") {
		img, err := imgbase64.FromLocal(d.Get("icon_local").(string))
		if err != nil {
			return err
		}
		params.Icon = img
		changed = true
	}
	if d.HasChange("icon_data_uri") {
		params.Icon = d.Get("icon_data_uri").(string)
		changed = true
	}
	if d.HasChange("owner_id") {
		params.OwnerID = d.Get("owner_id").(string)
		changed = true
	}

	if changed {
		_, err := client.GuildEdit(d.Id(), params)
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
	pData.Set("region", guild.Region)
	pData.Set("default_message_notifications", guild.DefaultMessageNotifications)
	pData.Set("afk_channel_id", guild.AfkChannelID)
	pData.Set("afk_timeout", guild.AfkTimeout)
	pData.Set("icon_hash", guild.Icon)
	pData.Set("owner_id", guild.OwnerID)
	results = append(results, pData)

	return results, nil
}