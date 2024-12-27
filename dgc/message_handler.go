package dgc

import (
	"errors"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

var (
	ErrMessageCommandHasNoMessage = errors.New("Message command was used but no message was given")
)

type MessageCommandHandler func(ctx *MessageExecutionContext, sender *discordgo.User) error

type messageCommand struct {
	handler MessageCommandHandler
}

func (c *messageCommand) manage(log *slog.Logger, sender *discordgo.User, ss *discordgo.Session, i *discordgo.InteractionCreate) (bool, error) {
	data := i.ApplicationCommandData()
	targetMessage := data.TargetID
	message, found := data.Resolved.Messages[targetMessage]
	if !found {
		return false, ErrMessageCommandHasNoMessage
	}
	ctx := MessageExecutionContext{
		respondingContext: respondingContext{
			executionContext: executionContext{
				log:     log,
				Session: ss,
				I:       i,
			},
		},
		Message: message,
	}
	err := c.handler(&ctx, sender)
	return ctx.alreadyResponded, err
}
