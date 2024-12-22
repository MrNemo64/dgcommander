package dgc2

import (
	"github.com/bwmarrin/discordgo"
)

type messageCommandBuilder struct {
	genericCommandBuilder[*messageCommandBuilder]
	handler MessageCommandHandler
}

func NewMessageCommand() *messageCommandBuilder {
	b := &messageCommandBuilder{}
	b.genericCommandBuilder.upper = b
	return b
}

func (b *messageCommandBuilder) create() command {
	return &messageCommand{
		handler: b.handler,
	}
}

func (b *messageCommandBuilder) discordDefineForCreation() *discordgo.ApplicationCommand {
	c := b.genericCommandBuilder.discordDefineForCreation()
	c.Type = discordgo.MessageApplicationCommand
	return c
}

func (b *messageCommandBuilder) Handler(handler MessageCommandHandler) *messageCommandBuilder {
	b.handler = handler
	return b
}
