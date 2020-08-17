package discord

import (
    "fmt"
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
            StateContext: resourceServerMemberImportState,
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
                Set:           schema.HashString,
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
    d.SetId(fmt.Sprintf("%s:%s", d.Get("server_id").(string), d.Get("user_id").(string)))

    client := m.(*Context).Client
    _, err := getServerMember(client, d.Id())
    d.Set("in_server", err == nil)

    if err == nil {
        diags = append(diags, resourceServerMemberRead(ctx, d, m)...)
        diags = append(diags, resourceServerMemberUpdate(ctx, d, m)...)
    }

    return diags
}

func resourceServerMemberRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    client := m.(*Context).Client

    u, err := getServerMember(client, d.Id())
    d.Set("in_server", err == nil)
    if err != nil {
        d.Set("joined_at", nil)
        d.Set("premium_since", nil)
        d.Set("roles", nil)
        return diags
    }
    member := u.Member

    d.Set("joined_at", member.JoinedAt)
    d.Set("premium_since", member.PremiumSince)
    d.Set("roles", member.Roles)

    return diags
}

func resourceServerMemberUpdate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    client := m.(*Context).Client

    u, err := getServerMember(client, d.Id())
    d.Set("in_server", err == nil)
    if err != nil {
        d.Set("joined_at", nil)
        d.Set("premium_since", nil)
        d.Set("roles", nil)
        return diags
    }

    if _, v := d.GetChange("roles"); v != nil {
        items := v.(*schema.Set).List()
        roles := make([]string, 0, len(items))
        for _, r := range items {
            _, roleId, err := parseTwoIds(r.(string))
            if err != nil {
                return diag.Errorf("Failed to edit member. Couldn't parse role ids: %s", err.Error())
            }
            roles = append(roles, roleId)
        }

        err := client.GuildMemberEdit(u.ServerId, u.UserId, roles)
        if err != nil {
            return diag.Errorf("Failed to edit member: %s", err.Error())
        }
    }

    return diags
}

func resourceServerMemberDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    client := m.(*Context).Client

    u, err := getServerMember(client, d.Id())
    if err != nil {
        return diag.Errorf("Failed to fetch server member %s: %s", d.Id(), err.Error())
    }

    err = client.GuildMemberDelete(u.ServerId, u.UserId)
    if err != nil {
        return diag.Errorf("Failed to remove member from the server: %s", err.Error())
    }

    return diags
}

func resourceServerMemberImportState(_ context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
    results := make([]*schema.ResourceData, 1, 1)
    results[0] = d

    client := m.(*Context).Client
    u, err := getServerMember(client, d.Id())
    if err != nil {
        return nil, err
    }

    d.Set("server_id", u.ServerId)
    d.Set("user_id", u.UserId)

    member := resourceDiscordServerMember()

    pData := member.Data(nil)
    pData.SetId(d.Id())
    pData.SetType("discord_server_member")
    d.Set("joined_at", u.Member.JoinedAt)
    d.Set("premium_since", u.Member.PremiumSince)
    d.Set("roles", u.Member.Roles)
    results = append(results, pData)

    return results, nil
}
