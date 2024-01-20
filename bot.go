package main

import "github.com/bwmarrin/discordgo"

type Bot struct {
	session        *discordgo.Session
	commandHandler map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

func NewBot(token string) (*Bot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		session:        session,
		commandHandler: make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)),
	}, nil
}

func (b *Bot) setup() error {
	b.session.Identify.Intents = discordgo.IntentGuildMembers | discordgo.IntentGuildMessages |
		discordgo.IntentGuilds | discordgo.IntentDirectMessages

	b.session.AddHandler(b.ready)
	b.session.AddHandler(b.handleApplicationCommand)

	// TODO: Register command
	return nil
}

func (b *Bot) ready(s *discordgo.Session, _ *discordgo.Ready) {
	_ = s.UpdateGameStatus(0, "초기화하는 중...")
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

	if handler, ok := b.commandHandler[name]; ok {
		handler(s, i)
	}
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

// TODO: Add command and handler
