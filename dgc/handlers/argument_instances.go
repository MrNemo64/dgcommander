package handlers

import "github.com/bwmarrin/discordgo"

type specificArgumentInstance interface {
	isOfType(discordgo.ApplicationCommandOptionType) bool
}

type genericArgumentInstance[T specificArgumentInstance] struct {
	specific T
	name     string
}

func (i *genericArgumentInstance[T]) parse(op *discordgo.ApplicationCommandInteractionDataOption) (name string, value any, err error) {
	if !i.specific.isOfType(op.Type) {
		return
	}
}
