package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	session  *discordgo.Session
	registry CommandRegistrar

	logging bool
}

func NewBot(token string, logging bool) (*Bot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		session:  session,
		registry: NewCommandRegistry(session),
		logging:  logging,
	}, nil
}

func (b *Bot) setup() error {
	b.session.Identify.Intents = discordgo.IntentGuildMembers | discordgo.IntentGuildMessages |
		discordgo.IntentGuilds | discordgo.IntentDirectMessages

	b.session.AddHandler(b.ready)
	b.session.AddHandler(b.guildCreate)
	b.session.AddHandler(b.guildDelete)
	b.session.AddHandler(b.handleApplicationCommand)

	return nil
}

func (b *Bot) ready(s *discordgo.Session, _ *discordgo.Ready) {
	_ = s.UpdateGameStatus(0, "초기화하는 중...")
}

func (b *Bot) guildCreate(s *discordgo.Session, g *discordgo.GuildCreate) {
	if g.Unavailable {
		return
	}

	err := s.UpdateGameStatus(0, fmt.Sprintf("%d개의 서버에서 대기", len(s.State.Guilds)))
	if err != nil {
		if b.logging {
			Error("Failed to update game status", err)
		}
	}
}

func (b *Bot) guildDelete(s *discordgo.Session, g *discordgo.GuildDelete) {
	if g.Unavailable {
		return
	}

	err := s.UpdateGameStatus(0, fmt.Sprintf("%d개의 서버에서 대기", len(s.State.Guilds)))
	if err != nil {
		if b.logging {
			Error("Failed to update game status", err)
		}
	}
}

func (b *Bot) handleApplicationCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var name string

	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		name = i.ApplicationCommandData().Name
	case discordgo.InteractionMessageComponent:
		name = i.MessageComponentData().CustomID
	default:
		return
	}

	b.registry.Handle(name, s, i)
}

func (b *Bot) Run() error {
	err := b.setup()
	if err != nil {
		return err
	}

	return b.session.Open()
}

func (b *Bot) Close() error {
	return b.session.Close()
}
