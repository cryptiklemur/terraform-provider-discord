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
            StateContext: schema.ImportStatePassthroughContext,
        },

        Schema: map[string]*schema.Schema{
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

    serverId := getId(d.Get("server_id").(string))
    server, err := client.GetGuild(ctx, serverId)
    if err != nil {
        return diag.Errorf("Server does not exist with that ID: %s", serverId)
    }

    role, err := client.CreateGuildRole(ctx, serverId, &disgord.CreateGuildRoleParams{
        Name:        d.Get("name").(string),
        Permissions: d.Get("permissions").(uint64),
        Color:       d.Get("color").(uint),
        Hoist:       d.Get("hoist").(bool),
        Mentionable: d.Get("mentionable").(bool),
    })

    roles, err := client.GetGuildRoles(ctx, serverId)
    if err != nil {
        return diag.Errorf("Failed to fetch roles: %s", err.Error())
    }
    index, exists := findRoleIndex(roles, role)
    if !exists {
        return diag.Errorf("Role somehow does not exists")
    }

    moveRole(roles, index, d.Get("position").(int))

    params := make([]disgord.UpdateGuildRolePositionsParams, 0, len(roles))
    for index, r := range roles {
        params = append(params, disgord.UpdateGuildRolePositionsParams{ID: r.ID, Position: index})
    }

    if ok, err := reorderRoles(ctx, m, serverId, role, d.Get("position").(int)); !ok {
        return err
    }

    d.SetId(role.ID.String())
    d.Set("server_id", server.ID.String())

    return resourceRoleRead(ctx, d, m)
}

func resourceRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    client := m.(*Context).Client

    serverId := getId(d.Get("server_id").(string))
    role, err := getRole(ctx, client, serverId, getId(d.Id()))
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
    roleId := getId(d.Id())

    builder := client.UpdateGuildRole(ctx, serverId, roleId)

    builder.SetName(d.Get("name").(string))
    if _, v := d.GetChange("color"); v.(int) > 0 {
        builder.SetColor(uint(v.(int)))
    }
    builder.SetHoist(d.Get("hoist").(bool))
    builder.SetMentionable(d.Get("mentionable").(bool))
    if _, v := d.GetChange("permission"); v != nil {
        builder.SetPermissions(uint64(v.(int)))
    }

    role, err := builder.Execute()
    if err != nil {
        return diag.Errorf("Failed to update role %s: %s", d.Id(), err.Error())
    }

    if d.HasChange("position") {
        if ok, err := reorderRoles(ctx, m, serverId, role, d.Get("position").(int)); !ok {
            return err
        }
    }

    return diags
}

func resourceRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    client := m.(*Context).Client

    serverId, roleId, err := getBothIds(d.Id())
    if err != nil {
        return diag.Errorf("Failed to fetch role %s: %s", d.Id(), err.Error())
    }

    err = client.DeleteGuildRole(ctx, serverId, roleId)
    if err != nil {
        return diag.Errorf("Failed to delete role: %s", err.Error())
    }

    return diags
}
