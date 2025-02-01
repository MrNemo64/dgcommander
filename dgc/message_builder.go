package dgc

import (
	"github.com/bwmarrin/discordgo"
)

type MessageCommandBuilder struct {
	genericCommandBuilder[*MessageCommandBuilder, MessageMiddleware]
	handler MessageCommandHandler
}

func NewMessageCommand() *MessageCommandBuilder {
	b := &MessageCommandBuilder{}
	b.genericCommandBuilder.upper = b
	b.genericCommandBuilder.name.Upper = b
	return b
}

func (b *MessageCommandBuilder) create() command {
	return &messageCommand{
		handler:     b.handler,
		middlewares: b.middlewares,
	}
}

func (b *MessageCommandBuilder) discordDefineForCreation() *discordgo.ApplicationCommand {
	c := b.genericCommandBuilder.discordDefineForCreation()
	c.Type = discordgo.MessageApplicationCommand
	return c
}

func (b *MessageCommandBuilder) Handler(handler MessageCommandHandler) *MessageCommandBuilder {
	b.handler = handler
	return b
}
