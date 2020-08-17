package discord

import (
    "github.com/bwmarrin/discordgo"
)

type Channel struct {
    ServerId string
    ChannelId string
    Channel *discordgo.Channel
}

func findChannelById(array []*discordgo.Channel, id string) *discordgo.Channel {
    for _, element := range array {
        if element.ID == id {
            return element
        }
    }

    return nil
}

func getChannel(client *discordgo.Session, combinedId string) (*Channel, error) {
    var c Channel

    serverId, channelId, err := parseTwoIds(combinedId)
    if err != nil {
        return nil, err
    }

    channels, err := client.GuildChannels(serverId)
    if err != nil {
        return nil, err
    }

    c.ServerId = serverId
    c.ChannelId = channelId
    c.Channel = findChannelById(channels, channelId)

    return &c, nil
}
