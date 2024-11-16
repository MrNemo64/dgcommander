package handlers

import (
	"errors"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

type MessageHandler func(ctx *MessageExecutionContext, sender *discordgo.User) error

type MessageExecutionContext struct {
	ExecutionContext
	Message *discordgo.Message
}

type MessageCommand struct {
	Handler MessageHandler
}

func (b *MessageCommand) Autocomplete(log *slog.Logger, sender *discordgo.User, ss *discordgo.Session, i *discordgo.InteractionCreate) error {
	return errors.New("message commands cannot be autocompleted")
}

func (c *MessageCommand) Manage(log *slog.Logger, sender *discordgo.User, ss *discordgo.Session, i *discordgo.InteractionCreate) (bool, error) {
	data := i.ApplicationCommandData()
	targetMessage := data.TargetID
	message, found := data.Resolved.Messages[targetMessage]
	if !found {
		return false, errors.New("Message command was used but no message was given")
	}
	ctx := MessageExecutionContext{
		ExecutionContext: ExecutionContext{
			log:     log,
			Session: ss,
			I:       i,
		},
		Message: message,
	}
	err := c.Handler(&ctx, sender)
	return ctx.alreadyResponded, err
}
