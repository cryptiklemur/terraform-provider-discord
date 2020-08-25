package discord

import (
    "context"
    "github.com/andersfylling/disgord"
    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDiscordRole() *schema.Resource {
    return &schema.Resource{
        ReadContext: dataSourceDiscordRoleRead,
        Schema: map[string]*schema.Schema{
            "server_id": {
                Type:     schema.TypeString,
                Required: true,
            },
            "role_id": {
                ExactlyOneOf: []string{"role_id", "name"},
                Type:         schema.TypeString,
                Optional:     true,
            },
            "name": {
                ExactlyOneOf: []string{"role_id", "name"},
                Type:         schema.TypeString,
                Optional:     true,
            },
            "position": {
                Type: schema.TypeString,
                Computed: true,
            },
            "color": {
                Type: schema.TypeInt,
                Computed: true,
            },
            "permissions": {
                Type: schema.TypeInt,
                Computed: true,
            },
            "hoist": {
                Type: schema.TypeBool,
                Computed: true,
            },
            "mentionable": {
                Type: schema.TypeBool,
                Computed: true,
            },
            "managed": {
                Type: schema.TypeBool,
                Computed: true,
            },
        },
    }
}

func dataSourceDiscordRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    var err error
    var role *disgord.Role
    client := m.(*Context).Client

    serverId := getId(d.Get("server_id").(string))
    server, err := client.GetGuild(ctx, serverId)
    if err != nil {
        return diag.Errorf("Failed to fetch server %s: %s", serverId.String(), err.Error())
    }

    if v, ok := d.GetOk("role_id"); ok {
        role, err = server.Role(getId(v.(string)))
        if err != nil {
            return diag.Errorf("Failed to fetch role %s: %s", v.(string), err.Error())
        }
    }

    if v, ok := d.GetOk("name"); ok {
        roles, err := server.RoleByName(v.(string))
        if err != nil {
            return diag.Errorf("Failed to fetch role %s: %s", v.(string), err.Error())
        }

        if len(roles) <= 0 {
            return diag.Errorf("Failed to fetch role %s", v.(string))
        }

        role = roles[0]
    }

    d.SetId(role.ID.String())
    d.Set("role_id", role.ID.String())
    d.Set("name", role.Name)
    d.Set("position", len(server.Roles)-role.Position)
    d.Set("color", role.Color)
    d.Set("hoist", role.Hoist)
    d.Set("mentionable", role.Mentionable)
    d.Set("permissions", role.Permissions)
    d.Set("managed", role.Managed)

    return diags
}
