package handlers

import "github.com/bwmarrin/discordgo"

type CommandArguments struct{}

type ArgumentList struct {
}

func (al *ArgumentList) ParseOptions(options []*discordgo.ApplicationCommandInteractionDataOption) (CommandArguments, error) {
	return CommandArguments{}, nil
}
