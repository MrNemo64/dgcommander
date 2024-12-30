package dgc

import (
	"github.com/bwmarrin/discordgo"
)

type UserCommandBuilder struct {
	genericCommandBuilder[*UserCommandBuilder]
	handler UserCommandHandler
}

func NewUserCommand() *UserCommandBuilder {
	b := &UserCommandBuilder{}
	b.genericCommandBuilder.upper = b
	b.genericCommandBuilder.name.Upper = b
	return b
}

func (b *UserCommandBuilder) create() command {
	return &userCommand{
		handler: b.handler,
	}
}

func (b *UserCommandBuilder) discordDefineForCreation() *discordgo.ApplicationCommand {
	c := b.genericCommandBuilder.discordDefineForCreation()
	c.Type = discordgo.UserApplicationCommand
	return c
}

func (b *UserCommandBuilder) Handler(handler UserCommandHandler) *UserCommandBuilder {
	b.handler = handler
	return b
}
