package discord

import (
    "github.com/bwmarrin/discordgo"
)

type Role struct {
    ServerId string
    RoleId string
    Role *discordgo.Role
}

func insertRole(array []*discordgo.Role, value *discordgo.Role, index int) []*discordgo.Role {
    return append(array[:index], append([]*discordgo.Role{value}, array[index:]...)...)
}

func removeRole(array []*discordgo.Role, index int) []*discordgo.Role {
    return append(array[:index], array[index+1:]...)
}

func moveRole(array []*discordgo.Role, srcIndex int, dstIndex int) []*discordgo.Role {
    value := array[srcIndex]
    return insertRole(removeRole(array, srcIndex), value, dstIndex)
}

func findRoleIndex(array []*discordgo.Role, value *discordgo.Role) (int, bool) {
    for index, element := range array {
        if element.ID == value.ID {
            return index, true
        }
    }

    return -1, false
}

func findRoleById(array []*discordgo.Role, id string) *discordgo.Role {
    for _, element := range array {
        if element.ID == id {
            return element
        }
    }

    return nil
}

func getRole(client *discordgo.Session, combinedId string) (*Role, error) {
    var c Role

    serverId, roleId, err := parseTwoIds(combinedId)
    if err != nil {
        return nil, err
    }

    roles, err := client.GuildRoles(serverId)
    if err != nil {
        return nil, err
    }

    c.ServerId = serverId
    c.RoleId = roleId
    c.Role = findRoleById(roles, roleId)

    return &c, nil
}
