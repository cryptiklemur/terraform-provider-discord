package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/net/context"
)

func resourceDiscordInvite() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInviteCreate,
		ReadContext:   resourceInviteRead,
		DeleteContext: resourceInviteDelete,

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

func resourceInviteCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*Context).Client

	_, channelId, err := parseTwoIds(d.Get("channel_id").(string))
	if err != nil {
		return diag.Errorf("Channel does not exist with that ID: %s", d.Get("channel_id").(string))
	}

	invite, err := client.ChannelInviteCreate(channelId, discordgo.Invite{
		MaxAge:    d.Get("max_age").(int),
		MaxUses:   d.Get("max_uses").(int),
		Temporary: d.Get("temporary").(bool),
		Unique:    d.Get("unique").(bool),
	})
	if err != nil {
		return diag.Errorf("Failed to create a invite: %s", err.Error())
	}

	d.SetId(invite.Code)

	return diags
}

func resourceInviteRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*Context).Client

	_, err := client.Invite(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceInviteDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*Context).Client

	_, err := client.InviteDelete(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
