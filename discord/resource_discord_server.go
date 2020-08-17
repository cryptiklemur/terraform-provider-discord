package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/polds/imgbase64"
	"golang.org/x/net/context"
	"log"
)

func resourceDiscordServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServerCreate,
		ReadContext:   resourceServerRead,
		UpdateContext: resourceServerUpdate,
		DeleteContext: resourceServerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceServerImportState,
		},

		Schema: map[string]*schema.Schema{
			"server_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: descriptions["discord_resource_server_id"],
			},
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
						errors = append(errors, fmt.Errorf("verification_level must be between 0 and 3 inclusive, got: %d", v))
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
						errors = append(errors, fmt.Errorf("default_message_notifications must be 0 or 1, got: %d", v))
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
						errors = append(errors, fmt.Errorf("afk_timeout must be greater than 0, got: %d", v))
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

func resourceServerCreate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*Context).Client

	name := d.Get("name").(string)
	server, err := client.GuildCreate(name)
	if err != nil {
		return diag.Errorf("Failed to create server: %s", err.Error())
	}

	if d.Get("empty").(bool) {
		log.Print("[DISCORD] Empty marked as true. Wiping server.")
		channels, err := client.GuildChannels(server.ID)
		if err != nil {
			return diag.Errorf("Failed to fetch channels for new server: %s", err.Error())
		}

		for _, channel := range channels {
			log.Println("[DISCORD] Deleting Channel: " + channel.ID)
			_, err := client.ChannelDelete(channel.ID)
			if err != nil {
				return diag.Errorf("Failed to delete channel for new server: %s", err.Error())
			}
		}
	}

	params := discordgo.GuildParams{
		Name:                        server.Name,
		Region:                      server.Region,
		VerificationLevel:           &server.VerificationLevel,
		AfkChannelID:                server.AfkChannelID,
		AfkTimeout:                  server.AfkTimeout,
		Icon:                        server.Icon,
		OwnerID:                     server.OwnerID,
		Splash:                      server.Splash,
	}
	edit := false
	if v, ok := d.GetOk("region"); ok {
		params.Region = v.(string)
		edit = true
	}
	if v, ok := d.GetOk("afk_channel_id"); ok {
		params.AfkChannelID = v.(string)
		edit = true
	}
	if v, ok := d.GetOk("afk_timeout"); ok {
		params.AfkTimeout = v.(int)
		edit = true
	}
	if v, ok := d.GetOk("verification_level"); ok {
		level := discordgo.VerificationLevel(v.(int))
		params.VerificationLevel = &level
		edit = true
	}
	if v, ok := d.GetOk("default_message_notifications"); ok {
		params.DefaultMessageNotifications = v.(int)
		edit = true
	}
	if v, ok := d.GetOk("icon_url"); ok {
		img := imgbase64.FromRemote(v.(string))
		params.Icon = img
		edit = true
	}
	if v, ok := d.GetOk("icon_local"); ok {
		params.Icon = v.(string)
		edit = true
	}
	if v, ok := d.GetOk("icon_data_uri"); ok {
		params.Icon = v.(string)
		edit = true
	}
	if v, ok := d.GetOk("owner_id"); ok {
		params.OwnerID = v.(string)
		edit = true
	}

	log.Println("[DISCORD] Setting icon to: " + params.Icon)
	if edit {
		_, err = client.GuildEdit(server.ID, params)
		if err != nil {
			return diag.Errorf("Failed to edit server: %s", err.Error())
		}
	}

	d.SetId(server.ID)
	d.Set("owner", server.OwnerID)

	return diags
}

func resourceServerRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*Context).Client

	server, err := client.Guild(d.Id())
	if err != nil {
		return diag.Errorf("Error fetching server: %s", err.Error())
	}

	d.Set("name", server.Name)
	d.Set("region", server.Region)
	d.Set("default_message_notifications", server.DefaultMessageNotifications)
	d.Set("afk_channel_id", server.AfkChannelID)
	d.Set("afk_timeout", server.AfkTimeout)
	d.Set("icon_hash", server.Icon)
	d.Set("verification_level", server.VerificationLevel)
	d.Set("default_message_notifications", server.DefaultMessageNotifications)

	// We don't want to set the owner to null, should only change this if its changing to something else
	if d.Get("owner_id").(string) != "" {
		d.Set("owner_id", server.OwnerID)
	}

	return diags
}

func resourceServerUpdate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*Context).Client

	server, err := client.Guild(d.Id())
	if err != nil {
		return diag.Errorf("Error fetching server: %s", err.Error())
	}

	changed := false
	params := discordgo.GuildParams{
		Name:                        server.Name,
		Region:                      server.Region,
		VerificationLevel:           &server.VerificationLevel,
		AfkChannelID:                server.AfkChannelID,
		AfkTimeout:                  server.AfkTimeout,
		Icon:                        server.Icon,
		OwnerID:                     server.OwnerID,
		Splash:                      server.Splash,
	}

	if d.HasChange("name") {
		params.Name = d.Get("name").(string)
		changed = true
	}
	if d.HasChange("region") {
		params.Region = d.Get("region").(string)
		changed = true
	}
	if d.HasChange("verification_level") {
		level := discordgo.VerificationLevel(d.Get("verification_level").(int))
		params.VerificationLevel = &level
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
		params.Icon = d.Get("icon_local").(string)
		log.Printf("[DISCORD] Setting icon to: %s", params.Icon)
		changed = true
	}

	if d.HasChange("icon_data_uri") {
		params.Icon = d.Get("icon_data_uri").(string)
		changed = true
	}

	ownerId, hasOwner := d.GetOk("owner_id")
	if d.HasChange("owner_id") {
		if hasOwner {
			params.OwnerID = ownerId.(string)
			changed = true
		}
	} else {
		if hasOwner {
			params.OwnerID = server.OwnerID
			changed = true
		}
	}

	if changed {
		_, err := client.GuildEdit(d.Id(), params)
		if err != nil {
			return diag.Errorf("Failed to edit server: %s", err.Error())
		}
	}

	return diags
}

func resourceServerDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*Context).Client

	_, err := client.GuildDelete(d.Id())
	if err != nil {
		return diag.Errorf("Failed to delete server: %s", err)
	}

	return diags
}

func resourceServerImportState(_ context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	results := make([]*schema.ResourceData, 1, 1)
	results[0] = d

	client := m.(*Context).Client
	server, err := client.Guild(d.Id())
	if err != nil {
		return results, nil
	}

	pData := resourceDiscordServer().Data(nil)
	pData.SetId(d.Id())
	pData.SetType("discord_server")
	pData.Set("name", server.Name)
	pData.Set("region", server.Region)
	pData.Set("default_message_notifications", server.DefaultMessageNotifications)
	pData.Set("afk_channel_id", server.AfkChannelID)
	pData.Set("afk_timeout", server.AfkTimeout)
	pData.Set("icon_hash", server.Icon)
	pData.Set("owner_id", server.OwnerID)
	results = append(results, pData)

	return results, nil
}