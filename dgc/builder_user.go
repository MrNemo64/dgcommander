package dgc

import "github.com/bwmarrin/discordgo"

type userBuilder struct {
	commandBuilder[*userBuilder]
}

func NewUser() *userBuilder {
	b := &userBuilder{}
	b.commandBuilder.upper = b
	return b
}

func (b *userBuilder) create() command {
	panic("not implemented")
}

func (b *userBuilder) discordDefineForCreation() *discordgo.ApplicationCommand {
	c := b.commandBuilder.discordDefineForCreation()
	c.Type = discordgo.UserApplicationCommand
	return c
}
