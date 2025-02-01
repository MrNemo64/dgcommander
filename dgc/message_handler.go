package dgc

import (
	"errors"
)

var (
	ErrMessageCommandHasNoMessage = errors.New("Message command was used but no message was given")
)

type MessageCommandHandler func(ctx *MessageExecutionContext) error
type MessageMiddleware = func(info *MessageExecutionContext, next func()) error

type messageCommand struct {
	middlewares []MessageMiddleware
	handler     MessageCommandHandler
}

func (c *messageCommand) execute(info *RespondingContext) (bool, error) {
	data := info.I.ApplicationCommandData()
	targetMessage := data.TargetID
	message, found := data.Resolved.Messages[targetMessage]
	if !found {
		return false, ErrMessageCommandHasNoMessage
	}
	ctx := MessageExecutionContext{
		RespondingContext: info,
		Message:           message,
	}
	mc := newMiddlewareChain(&ctx, c.middlewares)
	if err := mc.startChain(); err != nil {
		return ctx.acknowledged, err
	}
	if mc.allMiddlewaresCalled {
		err := c.handler(&ctx)
		return ctx.acknowledged, err
	}
	return ctx.acknowledged, nil
}
