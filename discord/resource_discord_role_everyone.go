package discord

import (
    "github.com/andersfylling/disgord"
    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "golang.org/x/net/context"
)

func resourceDiscordRoleEveryone() *schema.Resource {
    return &schema.Resource{
        CreateContext: resourceRoleEveryoneRead,
        ReadContext:   resourceRoleEveryoneRead,
        UpdateContext: resourceRoleEveryoneUpdate,
        DeleteContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
            return []diag.Diagnostic{{
                Severity: diag.Warning,
                Summary:  "Deleting the everyone role is not allowed",
            }}
        },
        Importer: &schema.ResourceImporter{
            StateContext: resourceRoleEveryoneImport,
        },

        Schema: map[string]*schema.Schema{
            "server_id": {
                Type:     schema.TypeString,
                Required: true,
                ForceNew: true,
            },
            "permissions": {
                Type:     schema.TypeInt,
                Optional: true,
                Default:  0,
                ForceNew: false,
            },
        },
    }
}

func resourceRoleEveryoneImport(ctx context.Context, data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
    data.SetId(data.Id())
    data.Set("server_id", getId(data.Id()).String())

    return schema.ImportStatePassthroughContext(ctx, data, i)
}

func resourceRoleEveryoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    client := m.(*Context).Client

    serverId := getId(d.Get("server_id").(string))
    d.SetId(serverId.String())

    server, err := client.GetGuild(ctx, serverId)
    if err != nil {
        return diag.Errorf("Failed to fetch server %s: %s", serverId.String(), err.Error())
    }

    role, err := server.Role(serverId)
    if err != nil {
        return diag.Errorf("Failed to fetch role %s: %s", d.Id(), err.Error())
    }

    d.Set("permissions", role.Permissions)

    return diags
}

func resourceRoleEveryoneUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    client := m.(*Context).Client

    serverId := getId(d.Get("server_id").(string))
    d.SetId(serverId.String())
    builder := client.UpdateGuildRole(ctx, serverId, serverId)

    builder.SetPermissions(disgord.PermissionBit(d.Get("permissions").(int)))

    role, err := builder.Execute()
    if err != nil {
        return diag.Errorf("Failed to update role %s: %s", d.Id(), err.Error())
    }

    d.Set("permissions", role.Permissions)

    return diags
}
