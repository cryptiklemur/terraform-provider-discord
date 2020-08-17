package discord

import (
    "fmt"
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
            StateContext: resourceRoleImportState,
        },

        Schema: map[string]*schema.Schema{
            "role_id": {
                Type:        schema.TypeString,
                Computed:    true,
                Description: descriptions["discord_resource_role_id"],
            },
            "server_id": {
                Type:        schema.TypeString,
                Required:    true,
                Description: descriptions["discord_resource_role_server"],
            },
            "name": {
                Type:        schema.TypeString,
                Required:    true,
                Description: descriptions["discord_resource_role_name"],
            },
            "permissions": {
                Type:        schema.TypeInt,
                Optional:    true,
                Default:     0,
                Description: descriptions["discord_resource_role_permissions"],
            },
            "color": {
                Type:        schema.TypeInt,
                Optional:    true,
                Default:     0xEB4034,
                Description: descriptions["discord_resource_role_color"],
            },
            "hoist": {
                Type:        schema.TypeBool,
                Optional:    true,
                Default:     false,
                Description: descriptions["discord_resource_role_hoist"],
            },
            "mentionable": {
                Type:        schema.TypeBool,
                Optional:    true,
                Default:     false,
                Description: descriptions["discord_resource_role_mentionable"],
            },
            "position": {
                Type:        schema.TypeInt,
                Optional:    true,
                Default:     1,
                Description: descriptions["discord_resource_role_position"],
                ValidateFunc: func(val interface{}, key string) (warns []string, errors []error) {
                    v := val.(int)

                    if v < 0 {
                        errors = append(errors, fmt.Errorf("position must be greater than 0, got: %d", v))
                    }

                    return
                },
            },
            "managed": {
                Type:        schema.TypeBool,
                Computed:    true,
                Description: descriptions["discord_resource_role_managed"],
            },
        },
    }
}

func resourceRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*Context).Client

    serverId := d.Get("server_id").(string)
    server, err := client.Guild(serverId)
    if err != nil {
        return diag.Errorf("Server does not exist with that ID: %s", serverId)
    }

    role, err := client.GuildRoleCreate(server.ID)
    if err != nil {
        return diag.Errorf("Failed to create a channel: %s", err.Error())
    }

    role, err = client.GuildRoleEdit(
        server.ID,
        role.ID,
        d.Get("name").(string),
        d.Get("color").(int),
        d.Get("hoist").(bool),
        d.Get("permissions").(int),
        d.Get("mentionable").(bool),
    )

    roles, err := client.GuildRoles(server.ID)
    if err != nil {
        return diag.Errorf("Failed to fetch roles: %s", err.Error())
    }
    index, exists := findRoleIndex(roles, role)
    if !exists {
        return diag.Errorf("Role somehow does not exists",)
    }

    moveRole(roles, index, d.Get("position").(int))
    roles, err = client.GuildRoleReorder(server.ID, roles)
    if err != nil {
        return diag.Errorf("Failed to re-order roles: %s", err.Error())
    }

    d.SetId(fmt.Sprintf("%s:%s", server.ID, role.ID))
    d.Set("role_id", role.ID)
    d.Set("server_id", server.ID)

    return resourceRoleRead(ctx, d, m)
}

func resourceRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    client := m.(*Context).Client

    r, err := getRole(client, d.Id())
    if err != nil {
        return diag.Errorf("Failed to fetch role %s: %s", d.Id(), err.Error())
    }
    role := r.Role

    d.Set("name", role.Name)
    d.Set("position", role.Position)
    d.Set("color", role.Color)
    d.Set("hoist", role.Hoist)
    d.Set("mentionable", role.Mentionable)
    d.Set("permissions", role.Permissions)

    return diags
}

func resourceRoleUpdate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    client := m.(*Context).Client

    r, err := getRole(client, d.Id())
    if err != nil {
        return diag.Errorf("Failed to fetch role %s: %s", d.Id(), err.Error())
    }

    role, err := client.GuildRoleEdit(
        r.ServerId,
        r.RoleId,
        d.Get("name").(string),
        d.Get("color").(int),
        d.Get("hoist").(bool),
        d.Get("permissions").(int),
        d.Get("mentionable").(bool),
    )
    if err != nil {
        return diag.Errorf("Failed to edit role: %s", err.Error())
    }

    if d.HasChange("position") {
        roles, err := client.GuildRoles(r.ServerId)
        if err != nil {
            return diag.Errorf("Failed to fetch roles: %s", err.Error())
        }
        index, exists := findRoleIndex(roles, role)
        if !exists {
            return diag.Errorf("Failed to find role")
        }

        moveRole(roles, index, d.Get("position").(int))
        roles, err = client.GuildRoleReorder(r.ServerId, roles)
        if err != nil {
            return diag.Errorf("Failed to re-order roles: %s", err.Error())
        }
    }

    return diags
}

func resourceRoleDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    client := m.(*Context).Client

    r, err := getRole(client, d.Id())
    if err != nil {
        return diag.Errorf("Failed to fetch role %s: %s", d.Id(), err.Error())
    }

    err = client.GuildRoleDelete(r.ServerId, r.RoleId)
    if err != nil {
        return diag.Errorf("Failed to delete role: %s", err.Error())
    }

    return diags
}

func resourceRoleImportState(_ context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
    client := m.(*Context).Client
    results := make([]*schema.ResourceData, 1, 1)
    results[0] = d

    serverId, roleId, err := parseTwoIds(d.Id())
    if err != nil {
        return nil, err
    }

    d.Set("role_id", roleId)
    d.Set("server_id", serverId)

    r, err := getRole(client, d.Id())
    if err != nil {
        return nil, err
    }

    pData := resourceDiscordRole().Data(nil)
    pData.SetId(d.Id())
    pData.SetType("discord_channel")
    d.Set("name", r.Role.Name)
    d.Set("color", r.Role.Color)
    d.Set("hoist", r.Role.Hoist)
    d.Set("permissions", r.Role.Permissions)
    d.Set("mentionable", r.Role.Mentionable)
    d.Set("position", r.Role.Position)
    d.Set("managed", r.Role.Managed)
    results = append(results, pData)

    return results, nil
}
