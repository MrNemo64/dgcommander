package dgc2

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

type SlashExecutionContext struct {
	executionContext
}

type SlashCommandHandler func(sender *discordgo.User, ctx *SlashExecutionContext) error

type slashCommand struct{}

func (c *slashCommand) manage(log *slog.Logger, sender *discordgo.User, ss *discordgo.Session, i *discordgo.InteractionCreate) (bool, error) {
	panic("TODO")
}

// TODO
