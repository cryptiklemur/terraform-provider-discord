package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"log"
	"time"
)

type Config struct {
	Token string
}

var ready = true

func (c *Config) Client() (interface{}, error) {
	var client *discordgo.Session

	client, err := discordgo.New("Bot " + c.Token)
	if err != nil {
		fmt.Println("Error connecting to discord.", err)

		return nil, err
	}

	err = client.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)

		return nil, err
	}

	client.AddHandlerOnce(readyHandler)

	i := 0
	for true {
		i++
		if ready {
			break
		}

		if i > 120 {
			return nil, errors.New("Bot failed to connect")
		}
		time.Sleep(time.Second)
	}

	return client, nil
}

func readyHandler(s *discordgo.Session, e *discordgo.Ready) {
	ready = true

	for _, guild := range e.Guilds {
		log.Println("Bot is in guild: " + guild.ID + " - " + guild.Name)
	}
}