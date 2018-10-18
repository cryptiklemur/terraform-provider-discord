package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Config struct {
	Token string
}

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

	defer client.Close()

	return client, nil
}
