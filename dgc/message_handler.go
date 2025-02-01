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

type messageMiddleHandler struct {
	ctx     *MessageExecutionContext
	handler MessageCommandHandler
}

func (mmh *messageMiddleHandler) handle() error {
	return mmh.handler(mmh.ctx)
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

	mmh := messageMiddleHandler{
		ctx:     &ctx,
		handler: c.handler,
	}

	mc := newMiddlewareChain(&ctx, c.middlewares, mmh.handle)
	if err := mc.startChain(); err != nil {
		return ctx.acknowledged, err
	}
	return ctx.acknowledged, nil
}
