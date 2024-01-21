package main

import "github.com/bwmarrin/discordgo"

type InteractionHandleFunc func(s *discordgo.Session, i *discordgo.InteractionCreate)

type Commander interface {
	Command() *discordgo.ApplicationCommand
	HandleFuncs() map[string]InteractionHandleFunc
}

type CommandRegistrar interface {
	Register(c Commander) error
	Handle(name string, s *discordgo.Session, i *discordgo.InteractionCreate)
	Unregister(id string) error
	UnregisterAll() error
}

type CommandRegistry struct {
	session  *discordgo.Session
	commands []*discordgo.ApplicationCommand
	handlers map[string]InteractionHandleFunc
}

func NewCommandRegistry(session *discordgo.Session) *CommandRegistry {
	return &CommandRegistry{
		session:  session,
		commands: make([]*discordgo.ApplicationCommand, 0),
		handlers: make(map[string]InteractionHandleFunc),
	}
}

func (r *CommandRegistry) Register(c Commander) error {
	registered, err := r.session.ApplicationCommandCreate(r.session.State.User.ID, "", c.Command())
	if err != nil {
		return err
	}

	r.commands = append(r.commands, registered)

	for name, handler := range c.HandleFuncs() {
		r.handlers[name] = handler
	}

	return nil
}

func (r *CommandRegistry) Handle(name string, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if handler, ok := r.handlers[name]; ok {
		handler(s, i)
	}
}

func (r *CommandRegistry) Unregister(id string) error {
	return r.session.ApplicationCommandDelete(r.session.State.User.ID, "", id)
}

func (r *CommandRegistry) UnregisterAll() error {
	for _, command := range r.commands {
		err := r.Unregister(command.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

var _ CommandRegistrar = (*CommandRegistry)(nil)
