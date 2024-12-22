package handlers

import "github.com/bwmarrin/discordgo"

func findOption(name string, options []*discordgo.ApplicationCommandInteractionDataOption) *discordgo.ApplicationCommandInteractionDataOption {
	for _, v := range options {
		if v.Name == name {
			return v
		}
	}
	return nil
}
