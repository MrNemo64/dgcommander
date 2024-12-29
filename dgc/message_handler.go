package dgc

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

var (
	ErrMessageCommandHasNoMessage = errors.New("Message command was used but no message was given")
)

type MessageCommandHandler func(ctx *MessageExecutionContext, sender *discordgo.User) error

type messageCommand struct {
	handler MessageCommandHandler
}

func (c *messageCommand) manage(info *InvokationInformation) (bool, error) {
	data := info.I.ApplicationCommandData()
	targetMessage := data.TargetID
	message, found := data.Resolved.Messages[targetMessage]
	if !found {
		return false, ErrMessageCommandHasNoMessage
	}
	ctx := MessageExecutionContext{
		respondingContext: respondingContext{
			executionContext: newExecutionContext(info),
		},
		Message: message,
	}
	err := c.handler(&ctx, info.Sender)
	return ctx.alreadyResponded, err
}
