package dgc

import (
	"github.com/MrNemo64/dgcommander/dgc/handlers"
	"github.com/bwmarrin/discordgo"
)

type messageBuilder struct {
	commandBuilder[handlers.MessageHandler, *messageBuilder]
}

func NewMessage() *messageBuilder {
	b := &messageBuilder{}
	b.commandBuilder.upper = b
	return b
}

func (b *messageBuilder) create() command {
	return &handlers.MessageCommand{
		Handler: b.handler,
	}
}

func (b *messageBuilder) discordDefineForCreation() *discordgo.ApplicationCommand {
	c := b.commandBuilder.discordDefineForCreation()
	c.Type = discordgo.MessageApplicationCommand
	return c
}
