package dgc

import "github.com/bwmarrin/discordgo"

type messageBuilder struct {
	commandBuilder[*messageBuilder]
}

func NewMessage() *messageBuilder {
	b := &messageBuilder{}
	b.commandBuilder.upper = b
	return b
}

func (b *messageBuilder) create() command {
	panic("not implemented")
}

func (b *messageBuilder) discordDefineForCreation() *discordgo.ApplicationCommand {
	c := b.commandBuilder.discordDefineForCreation()
	c.Type = discordgo.MessageApplicationCommand
	return c
}
