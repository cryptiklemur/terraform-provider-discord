package discord

import (
    "github.com/andersfylling/disgord"
    "github.com/bwmarrin/discordgo"
)

func getTextChannelType(channelType uint) (string, bool) {
    switch channelType {
    case 0:
        return "text", true
    case 2:
        return "voice", true
    case 4:
        return "category", true
    case 5:
        return "news", true
    case 6:
        return "store", true
    }

    return "text", false
}

func getDiscordChannelType(name string) (uint, bool) {
    switch name {
    case "text":
        return 0, true
    case "voice":
        return 2, true
    case "category":
        return 4, true
    case "news":
        return 5, true
    case "store":
        return 6, true
    }

    return 0, false
}

type Channel struct {
    ServerId string
    ChannelId string
    Channel *discordgo.Channel
}

func findChannelById(array []*disgord.Channel, id disgord.Snowflake) *disgord.Channel {
    for _, element := range array {
        if element.ID == id {
            return element
        }
    }

    return nil
}