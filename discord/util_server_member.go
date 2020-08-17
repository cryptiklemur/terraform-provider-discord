package discord

import (
    "github.com/bwmarrin/discordgo"
)

type ServerMember struct {
    ServerId string
    UserId string
    Member *discordgo.Member
}

func findMemberById(array []*discordgo.Member, id string) *discordgo.Member {
    for _, element := range array {
        if element.User.ID == id {
            return element
        }
    }

    return nil
}

func getServerMember(client *discordgo.Session, combinedId string) (*ServerMember, error) {
    var c ServerMember

    serverId, memberId, err := parseTwoIds(combinedId)
    if err != nil {
        return nil, err
    }

    c.ServerId = serverId
    c.UserId = memberId
    c.Member, err = client.GuildMember(serverId, memberId)
    if err != nil {
        return &c, err
    }

    return &c, nil
}
