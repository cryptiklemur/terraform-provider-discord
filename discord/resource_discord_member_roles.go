package discord

import (
    "encoding/json"
    "fmt"
    "github.com/andersfylling/disgord"
    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "golang.org/x/net/context"
    "strconv"
)

type RoleSchema struct {
    RoleId  disgord.Snowflake `json:"role_id"`
    HasRole bool              `json:"has_role"`
}

func convertToRoleSchema(v interface{}) (*RoleSchema, error) {
    var roleSchema *RoleSchema
    j, _ := json.MarshalIndent(v, "", "    ")
    err := json.Unmarshal(j, &roleSchema)

    return roleSchema, err
}

func resourceDiscordMemberRoles() *schema.Resource {
    return &schema.Resource{
        CreateContext: resourceMemberRolesCreate,
        ReadContext:   resourceMemberRolesRead,
        UpdateContext: resourceMemberRolesUpdate,
        DeleteContext: resourceMemberRolesDelete,
        Importer: &schema.ResourceImporter{
            StateContext: schema.ImportStatePassthroughContext,
        },

        Schema: map[string]*schema.Schema{
            "user_id": {
                Type:     schema.TypeString,
                Required: true,
            },
            "server_id": {
                Type:     schema.TypeString,
                Required: true,
            },
            "role": {
                Type: schema.TypeSet,
                Required: true,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "role_id": {
                            Type:     schema.TypeString,
                            Required: true,
                        },
                        "has_role": {
                            Type:     schema.TypeBool,
                            Optional: true,
                            Default:  true,
                        },
                    },
                },
            },
        },
    }
}

func resourceMemberRolesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics

    client := m.(*Context).Client

    serverId := getId(d.Get("server_id").(string))
    userId := getId(d.Get("user_id").(string))

    _, err := client.GetMember(ctx, serverId, userId)
    if err != nil {
        return diag.Errorf("Could not get member %s in %s: %s", userId.String(), serverId.String(), err.Error())
    }

    d.SetId(strconv.Itoa(Hashcode(fmt.Sprintf("%s:%s", serverId.String(), userId.String()))))

    diags = append(diags, resourceMemberRolesRead(ctx, d, m)...)
    diags = append(diags, resourceMemberRolesUpdate(ctx, d, m)...)

    return diags
}

func resourceMemberRolesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    client := m.(*Context).Client
    serverId := getId(d.Get("server_id").(string))
    userId := getId(d.Get("user_id").(string))

    member, err := client.GetMember(ctx, serverId, userId)
    if err != nil {
        return diag.Errorf("Could not get member %s in %s: %s", userId.String(), serverId.String(), err.Error())
    }

    items := d.Get("role").(*schema.Set).List()
    roles := make([]*RoleSchema, 0, len(items))

    for _, r := range items {
        v, _ := convertToRoleSchema(r)
        if hasRole(member, v.RoleId) {
            roles = append(roles, &RoleSchema{RoleId: v.RoleId, HasRole: true})
        } else {
            roles = append(roles, &RoleSchema{RoleId: v.RoleId, HasRole: false})
        }
    }

    return diags
}

func resourceMemberRolesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    client := m.(*Context).Client

    serverId := getId(d.Get("server_id").(string))
    userId := getId(d.Get("user_id").(string))

    member, err := client.GetMember(ctx, serverId, userId)
    if err != nil {
        return diag.Errorf("Could not get member %s in %s: %s", userId.String(), serverId.String(), err.Error())
    }

    builder := client.UpdateGuildMember(ctx, serverId, userId)

    old, new := d.GetChange("role")
    oldItems := old.(*schema.Set).List()
    items := new.(*schema.Set).List()

    roles := member.Roles

    for _, r := range items {
        v, _ := convertToRoleSchema(r)
        hasRole := hasRole(member, v.RoleId)
        // if its supposed to have the role, and it does, continue
        if hasRole && v.HasRole {
            continue
        }
        // If its supposed to have the role, and it doesnt, add it
        if v.HasRole && !hasRole {
            roles = append(roles, v.RoleId)
        }
        // If its not supposed to have the role, and it does, remove it
        if !v.HasRole && hasRole {
            roles = removeRoleById(roles, v.RoleId)
        }
    }

    // If the change removed the role, and the user has it, remove it
    for _, r := range oldItems {
        v, _ := convertToRoleSchema(r)
        if wasRemoved(items, v) && v.HasRole {
            roles = removeRoleById(roles, v.RoleId)
        }
    }

    builder.SetRoles(roles)

    err = builder.Execute()
    if err != nil {
        return diag.Errorf("Failed to edit member %s: %s", userId.String(), err.Error())
    }

    return diags
}

func wasRemoved(items []interface{}, v *RoleSchema) bool {
    for _, i := range items {
        item, _ := convertToRoleSchema(i)
        if item.RoleId == v.RoleId {
            return false
        }
    }

    return true
}

func resourceMemberRolesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    client := m.(*Context).Client
    serverId := getId(d.Get("server_id").(string))
    userId := getId(d.Get("user_id").(string))

    member, err := client.GetMember(ctx, serverId, userId)
    if err != nil {
        return diag.Errorf("Could not get member %s in %s: %s", userId.String(), serverId.String(), err.Error())
    }

    builder := client.UpdateGuildMember(ctx, serverId, userId)

    items := d.Get("role").(*schema.Set).List()
    roles := member.Roles

    for _, r := range items {
        v, _ := convertToRoleSchema(r)
        hasRole := hasRole(member, v.RoleId)
        // if its supposed to have the role, and it does, remove it
        if hasRole && v.HasRole {
            roles = removeRoleById(roles, v.RoleId)
        }
    }

    builder.SetRoles(roles)

    return diags
}
