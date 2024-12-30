package dgc

import (
	"errors"
)

var (
	ErrMessageCommandHasNoMessage = errors.New("Message command was used but no message was given")
)

type MessageCommandHandler func(ctx *MessageExecutionContext) error

type messageCommand struct {
	handler MessageCommandHandler
}

func (c *messageCommand) execute(info *InvokationInformation) (bool, error) {
	data := info.I.ApplicationCommandData()
	targetMessage := data.TargetID
	message, found := data.Resolved.Messages[targetMessage]
	if !found {
		return false, ErrMessageCommandHasNoMessage
	}
	ctx := MessageExecutionContext{
		respondingContext: respondingContext{
			executionContext: newExecutionContext(info.DGC.ctx, info),
		},
		Message: message,
	}
	err := c.handler(&ctx)
	return ctx.alreadyResponded, err
}
