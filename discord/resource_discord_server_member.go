package discord

import (
    "github.com/andersfylling/disgord"
    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "golang.org/x/net/context"
)

func resourceDiscordServerMember() *schema.Resource {
    return &schema.Resource{
        CreateContext: resourceServerMemberCreate,
        ReadContext:   resourceServerMemberRead,
        UpdateContext: resourceServerMemberUpdate,
        DeleteContext: resourceServerMemberDelete,
        Importer: &schema.ResourceImporter{
            StateContext: schema.ImportStatePassthroughContext,
        },

        Schema: map[string]*schema.Schema{
            "server_id": {
                Type:        schema.TypeString,
                Required:    true,
                Description: descriptions["discord_resource_server_member_server"],
            },
            "user_id": {
                Type:        schema.TypeString,
                Required:    true,
                Description: descriptions["discord_resource_server_member_user"],
            },
            "joined_at": {
                Type:        schema.TypeString,
                Computed:    true,
                Description: descriptions["discord_resource_server_member_joined_at"],
            },
            "premium_since": {
                Type:        schema.TypeString,
                Computed:    true,
                Description: descriptions["discord_resource_server_member_premium_since"],
            },
            "roles": {
                Type:        schema.TypeSet,
                Elem:        &schema.Schema{Type: schema.TypeString},
                Optional:    true,
                Description: descriptions["discord_resource_server_member_roles"],
                Set:         schema.HashString,
            },
            "in_server": {
                Type:     schema.TypeBool,
                Computed: true,
            },
        },
    }
}

func resourceServerMemberCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics

    client := m.(*Context).Client

    serverId := getId(d.Get("server_id").(string))
    userId := getId(d.Get("user_id").(string))

    _, err := client.GetMember(ctx, serverId, userId)
    d.SetId(userId.String())
    d.Set("in_server", err == nil)

    if err == nil {
        diags = append(diags, resourceServerMemberRead(ctx, d, m)...)
        diags = append(diags, resourceServerMemberUpdate(ctx, d, m)...)
    }

    return diags
}

func resourceServerMemberRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    client := m.(*Context).Client
    serverId := getId(d.Get("server_id").(string))
    userId := getId(d.Get("user_id").(string))

    member, err := client.GetMember(ctx, serverId, userId)
    d.Set("in_server", err == nil)
    if err != nil {
        d.Set("joined_at", nil)
        d.Set("premium_since", nil)
        d.Set("roles", nil)
        return diags
    }

    d.Set("joined_at", member.JoinedAt)
    d.Set("premium_since", member.PremiumSince)

    roles := make([]string, 0, len(member.Roles))
    for _, r := range member.Roles {
        roles = append(roles, r.String())
    }
    d.Set("roles", roles)

    return diags
}

func resourceServerMemberUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    client := m.(*Context).Client

    serverId := getId(d.Get("server_id").(string))
    userId := getId(d.Get("user_id").(string))
    builder := client.UpdateGuildMember(ctx, serverId, userId)

    _, err := client.GetMember(ctx, serverId, userId)
    d.Set("in_server", err == nil)
    if err != nil {
        d.Set("joined_at", nil)
        d.Set("premium_since", nil)
        d.Set("roles", nil)
        return diags
    }

    if _, v := d.GetChange("roles"); v != nil {
        items := v.(*schema.Set).List()
        roles := make([]disgord.Snowflake, 0, len(items))
        for _, r := range items {
            roles = append(roles, getId(r.(string)))
        }

        builder.SetRoles(roles)

        err = builder.Execute()
        if err != nil {
            return diag.Errorf("Failed to edit member: %s", err.Error())
        }
    }

    return diags
}

func resourceServerMemberDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    client := m.(*Context).Client
    serverId := getId(d.Get("server_id").(string))
    userId := getId(d.Get("user_id").(string))

    err := client.KickMember(ctx, serverId, userId, "Removed via Terraform")
    if err != nil {
        return diag.Errorf("Failed to remove member from the server: %s", err.Error())
    }

    return diags
}
