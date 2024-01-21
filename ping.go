package main

import "github.com/bwmarrin/discordgo"

type PingCommand struct{}

func (c *PingCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Ping!",
	}
}

func (c *PingCommand) HandleFuncs() map[string]InteractionHandleFunc {
	return map[string]InteractionHandleFunc{
		"ping": c.ping,
	}
}

func (c *PingCommand) ping(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong!",
		},
	})
}
