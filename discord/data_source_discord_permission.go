package discord

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var permissions map[string]int

func dataSourceDiscordPermission() *schema.Resource {
	permissions = map[string]int{
		"create_instant_invite":  0x00000001,
		"kick_members":           0x00000002,
		"ban_members":            0x00000004,
		"administrator":          0x00000008,
		"manage_channels":        0x000000010,
		"manage_guild":           0x000000020,
		"add_reactions":          0x000000040,
		"view_audit_log":         0x000000080,
		"priority_speaker":       0x000000100,
		"stream":                 0x000000200,
		"view_channel":           0x000000400,
		"send_messages":          0x000000800,
		"send_tts_messages":      0x000001000,
		"manage_messages":        0x00002000,
		"embed_links":            0x00004000,
		"attach_files":           0x00008000,
		"read_message_history":   0x00010000,
		"mention_everyone":       0x00020000,
		"use_external_emojis":    0x00040000,
		"view_guild_insights":    0x00080000,
		"connect":                0x00100000,
		"speak":                  0x00200000,
		"mute_members":           0x00400000,
		"deafen_members":         0x00800000,
		"move_members":           0x01000000,
		"use_vad":                0x02000000,
		"change_nickname":        0x04000000,
		"manage_nicknames":       0x08000000,
		"manage_roles":           0x10000000,
		"manage_webhooks":        0x20000000,
		"manage_emojis":          0x40000000,
		"request_to_speak":       0x100000000,
		"manage_threads":         0x400000000,
		"create_public_threads":  0x800000000,
		"create_private_threads": 0x1000000000,
		"use_external_stickers":  0x2000000000,
		"send_thread_messages":   0x4000000000,
	}

	schemaMap := make(map[string]*schema.Schema)
	schemaMap["allow_extends"] = &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
	}
	schemaMap["deny_extends"] = &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
	}
	schemaMap["allow_bits"] = &schema.Schema{
		Type:     schema.TypeInt,
		Computed: true,
	}
	schemaMap["deny_bits"] = &schema.Schema{
		Type:     schema.TypeInt,
		Computed: true,
	}
	for k := range permissions {
		schemaMap[k] = &schema.Schema{
			Optional: true,
			Type:     schema.TypeString,
			Default:  "unset",
			ValidateDiagFunc: func(v interface{}, path cty.Path) (diags diag.Diagnostics) {
				str := v.(string)
				allowed := [3]string{"allow", "unset", "deny"}

				if !contains(allowed, str) {
					return append(diags, diag.Errorf("%s is not an allowed value. Pick one of: allowed, unset, deny", str)...)
				}

				return diags
			},
		}
	}

	return &schema.Resource{
		ReadContext: dataSourceDiscordPermissionRead,
		Schema:      schemaMap,
	}
}

func dataSourceDiscordPermissionRead(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	var allowBits int
	var denyBits int
	for perm, bit := range permissions {
		v := d.Get(perm).(string)
		if v == "allow" {
			allowBits |= bit
		}
		if v == "deny" {
			denyBits |= bit
		}
	}

	d.SetId(strconv.Itoa(Hashcode(fmt.Sprintf("%d:%d", allowBits, denyBits))))
	d.Set("allow_bits", allowBits|(d.Get("allow_extends").(int)))
	d.Set("deny_bits", denyBits|(d.Get("deny_extends").(int)))

	return diags
}
