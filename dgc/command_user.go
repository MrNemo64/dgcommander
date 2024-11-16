package dgc

import (
	"context"
	"errors"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

type UserCommandHandler func(context context.Context, sender *discordgo.User) error

type userCommand struct {
	handler UserCommandHandler
}

func (c *userCommand) manage(log *slog.Logger, ss *discordgo.Session, i *discordgo.InteractionCreate) error {
	return nil
}

func (*userCommand) autocomplete(_ *slog.Logger, _ *discordgo.Session, _ *discordgo.InteractionCreate) error {
	return errors.New("Message commands cannot be autocompleted")
}
