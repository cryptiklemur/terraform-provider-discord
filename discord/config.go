package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
)

type Config struct {
	Token string
	ClientID string
	Secret string
}

type Context struct {
	Client *discordgo.Session
	Config *Config
}

var ready = true

func (c *Config) Client() (*Context, error) {
	var client *discordgo.Session

	client, err := discordgo.New("Bot " + c.Token)
	if err != nil {
		fmt.Println("Error connecting to now.", err)

		return nil, err
	}

	return &Context{
		Client: client,
		Config: c,
	}, nil
}

func readyHandler(s *discordgo.Session, e *discordgo.Ready) {
	ready = true

	for _, server := range e.Guilds {
		log.Println("Bot is in server: " + server.ID + " - " + server.Name)
	}
}
