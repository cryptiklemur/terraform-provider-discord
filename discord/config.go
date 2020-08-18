package discord

import (
    "github.com/andersfylling/disgord"
)

type Config struct {
    Token    string
    ClientID string
    Secret   string
}

type Context struct {
    Client *disgord.Client
    Config *Config
}

func (c *Config) Client() (*Context, error) {
    client := disgord.New(disgord.Config{
        BotToken: c.Token,
    })

    return &Context{Client: client, Config: c}, nil
}
