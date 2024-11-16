package dgc

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

type slashCommand struct {
}

func (c *slashCommand) manage(_ *slog.Logger, _ *discordgo.Session, _ *discordgo.InteractionCreate) error {
	return nil
}
