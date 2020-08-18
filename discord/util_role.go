package discord

import (
    "context"
    "github.com/andersfylling/disgord"
    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type Role struct {
    ServerId disgord.Snowflake
    RoleId disgord.Snowflake
    Role *disgord.Role
}

func insertRole(array []*disgord.Role, value *disgord.Role, index int) []*disgord.Role {
    return append(array[:index], append([]*disgord.Role{value}, array[index:]...)...)
}

func removeRole(array []*disgord.Role, index int) []*disgord.Role {
    return append(array[:index], array[index+1:]...)
}

func moveRole(array []*disgord.Role, srcIndex int, dstIndex int) []*disgord.Role {
    value := array[srcIndex]
    return insertRole(removeRole(array, srcIndex), value, dstIndex)
}

func findRoleIndex(array []*disgord.Role, value *disgord.Role) (int, bool) {
    for index, element := range array {
        if element.ID == value.ID {
            return index, true
        }
    }

    return -1, false
}

func findRoleById(array []*disgord.Role, id disgord.Snowflake) *disgord.Role {
    for _, element := range array {
        if element.ID == id {
            return element
        }
    }

    return nil
}

func reorderRoles(ctx context.Context, m interface{}, serverId disgord.Snowflake, role *disgord.Role, position int) (bool, diag.Diagnostics) {
    client := m.(*Context).Client

    roles, err := client.GetGuildRoles(ctx, serverId)
    if err != nil {
        return false, diag.Errorf("Failed to fetch roles: %s", err.Error())
    }
    index, exists := findRoleIndex(roles, role)
    if !exists {
        return false, diag.Errorf("Role somehow does not exists",)
    }

    moveRole(roles, index, position)

    params := make([]disgord.UpdateGuildRolePositionsParams, 0, len(roles))
    for index, r := range roles {
        params = append(params, disgord.UpdateGuildRolePositionsParams{ID: r.ID, Position: index})
    }

    roles, err = client.UpdateGuildRolePositions(ctx, serverId, params)
    if err != nil {
        return false, diag.Errorf("Failed to re-order roles: %s", err.Error())
    }

    return true, nil
}

func getRole(ctx context.Context, client *disgord.Client, serverId disgord.Snowflake, roleId disgord.Snowflake) (*disgord.Role, error) {
    roles, err := client.GetGuildRoles(ctx, serverId)
    if err != nil {
        return nil, err
    }

    role := findRoleById(roles, roleId)

    return role, nil
}
