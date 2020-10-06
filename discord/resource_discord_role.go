package discord

import (
    "fmt"
    "github.com/andersfylling/disgord"
    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "golang.org/x/net/context"
)

func resourceDiscordRole() *schema.Resource {
    return &schema.Resource{
        CreateContext: resourceRoleCreate,
        ReadContext:   resourceRoleRead,
        UpdateContext: resourceRoleUpdate,
        DeleteContext: resourceRoleDelete,
        Importer: &schema.ResourceImporter{
            StateContext: resourceRoleImport,
        },

        Schema: map[string]*schema.Schema{
            "server_id": {
                Type:     schema.TypeString,
                Required: true,
                ForceNew: true,
            },
            "name": {
                Type:     schema.TypeString,
                Required: true,
                ForceNew: false,
            },
            "permissions": {
                Type:     schema.TypeInt,
                Optional: true,
                Default:  0,
                ForceNew: false,
            },
            "color": {
                Type:     schema.TypeInt,
                Optional: true,
                ForceNew: false,
            },
            "hoist": {
                Type:     schema.TypeBool,
                Optional: true,
                Default:  false,
                ForceNew: false,
            },
            "mentionable": {
                Type:     schema.TypeBool,
                Optional: true,
                Default:  false,
                ForceNew: false,
            },
            "position": {
                Type:     schema.TypeInt,
                Optional: true,
                Default:  1,
                ForceNew: false,

                ValidateFunc: func(val interface{}, key string) (warns []string, errors []error) {
                    v := val.(int)

                    if v < 0 {
                        errors = append(errors, fmt.Errorf("position must be greater than 0, got: %d", v))
                    }

                    return
                },
            },
            "managed": {
                Type:     schema.TypeBool,
                Computed: true,
            },
        },
    }
}

func resourceRoleImport(ctx context.Context, data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
    serverId, roleId, err := getBothIds(data.Id())
    if err != nil {
        return nil, err
    }

    data.SetId(roleId.String())
    data.Set("server_id", serverId.String())

    return schema.ImportStatePassthroughContext(ctx, data, i)
}

func resourceRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    client := m.(*Context).Client

    serverId := getId(d.Get("server_id").(string))
    server, err := client.GetGuild(ctx, serverId)
    if err != nil {
        return diag.Errorf("Server does not exist with that ID: %s", serverId)
    }

    role, err := client.CreateGuildRole(ctx, serverId, &disgord.CreateGuildRoleParams{
        Name:        d.Get("name").(string),
        Permissions: uint64(d.Get("permissions").(int)),
        Color:       uint(d.Get("color").(int)),
        Hoist:       d.Get("hoist").(bool),
        Mentionable: d.Get("mentionable").(bool),
    })
    if err != nil {
        return diag.Errorf("Failed to create role for %s: %s", serverId.String(), err.Error())
    }

    if newPosition, ok := d.GetOk("position"); ok {
        var oldRole *disgord.Role
        for _, r := range server.Roles {
            if r.Position == newPosition.(int) {
                oldRole = r
                break
            }
        }
        if oldRole == nil {
            _, err = client.UpdateGuildRolePositions(ctx, serverId, []disgord.UpdateGuildRolePositionsParams{
                {ID: role.ID, Position: newPosition.(int)},
            })
        } else {
            _, err = client.UpdateGuildRolePositions(ctx, serverId, []disgord.UpdateGuildRolePositionsParams{
                {ID: oldRole.ID, Position: role.Position},
                {ID: role.ID, Position: newPosition.(int)},
            })
        }
        if err != nil {
            diags = append(diags, diag.Errorf("Failed to re-order roles: %s", err.Error())...)
        } else {
            d.Set("position", newPosition)
        }
    }

    d.SetId(role.ID.String())
    d.Set("server_id", server.ID.String())
    d.Set("managed", role.Managed)

    return diags
}

func resourceRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    client := m.(*Context).Client

    serverId := getId(d.Get("server_id").(string))
    server, err := client.GetGuild(ctx, serverId)
    if err != nil {
        return diag.Errorf("Failed to fetch server %s: %s", serverId.String(), err.Error())
    }

    role, err := server.Role(getId(d.Id()))
    if err != nil {
        return diag.Errorf("Failed to fetch role %s: %s", d.Id(), err.Error())
    }

    d.Set("name", role.Name)
    d.Set("position", role.Position)
    d.Set("color", role.Color)
    d.Set("hoist", role.Hoist)
    d.Set("mentionable", role.Mentionable)
    d.Set("permissions", role.Permissions)
    d.Set("managed", role.Managed)

    return diags
}

func resourceRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    client := m.(*Context).Client

    serverId := getId(d.Get("server_id").(string))
    server, err := client.GetGuild(ctx, serverId)
    if err != nil {
        return diag.Errorf("Failed to fetch server %s: %s", serverId.String(), err.Error())
    }

    roleId := getId(d.Id())
    role, err := server.Role(roleId)
    if err != nil {
        return diag.Errorf("Failed to fetch role %s: %s", d.Id(), err.Error())
    }

    if d.HasChange("position") {
        _, newPosition := d.GetChange("position")
        var oldRole *disgord.Role
        for _, r := range server.Roles {
            if r.Position == newPosition.(int) {
                oldRole = r
                break
            }
        }
        if oldRole == nil {
            _, err = client.UpdateGuildRolePositions(ctx, serverId, []disgord.UpdateGuildRolePositionsParams{
                {ID: roleId, Position: newPosition.(int)},
            })
        } else {
            _, err = client.UpdateGuildRolePositions(ctx, serverId, []disgord.UpdateGuildRolePositionsParams{
                {ID: oldRole.ID, Position: role.Position},
                {ID: roleId, Position: newPosition.(int)},
            })
        }
        if err != nil {
            diags = append(diags, diag.Errorf("Failed to re-order roles: %s", err.Error())...)
        } else {
            d.Set("position", newPosition)
        }
    }

    builder := client.UpdateGuildRole(ctx, serverId, roleId)

    builder.SetName(d.Get("name").(string))
    if _, v := d.GetChange("color"); v.(int) > 0 {
        builder.SetColor(uint(v.(int)))
    }
    builder.SetHoist(d.Get("hoist").(bool))
    builder.SetMentionable(d.Get("mentionable").(bool))
    builder.SetPermissions(disgord.PermissionBit(d.Get("permissions").(int)))

    role, err = builder.Execute()
    if err != nil {
        return diag.Errorf("Failed to update role %s: %s", d.Id(), err.Error())
    }

    d.Set("name", role.Name)
    d.Set("position", role.Position)
    d.Set("color", role.Color)
    d.Set("hoist", role.Hoist)
    d.Set("mentionable", role.Mentionable)
    d.Set("permissions", role.Permissions)
    d.Set("managed", role.Managed)

    return diags
}

func resourceRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    client := m.(*Context).Client

    serverId := getId(d.Get("server_id").(string))
    roleId := getId(d.Id())

    err := client.DeleteGuildRole(ctx, serverId, roleId)
    if err != nil {
        return diag.Errorf("Failed to delete role: %s", err.Error())
    }

    return diags
}
