package dgc

import (
	"errors"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

type messageCommand struct{}

func (c *messageCommand) manage(_ *slog.Logger, _ *discordgo.Session, _ *discordgo.InteractionCreate) error {
	return nil
}

func (*messageCommand) autocomplete(_ *slog.Logger, _ *discordgo.Session, _ *discordgo.InteractionCreate) error {
	return errors.New("Message commands cannot be autocompleted")
}
