package dgc

import (
	"errors"
)

var (
	ErrUserCommandHasNoUser = errors.New("User command was used but no user was given")
)

type UserCommandHandler func(ctx *UserExecutionContext) error
type UserMiddleware = func(ctx *UserExecutionContext, next func()) error

type userCommand struct {
	middlewares []UserMiddleware
	handler     UserCommandHandler
}

func (c *userCommand) execute(info *RespondingContext) (bool, error) {
	data := info.I.ApplicationCommandData()
	targetUser := data.TargetID
	user, found := data.Resolved.Users[targetUser]
	if !found {
		return false, ErrUserCommandHasNoUser
	}
	member, found := data.Resolved.Members[targetUser]
	if found {
		member.GuildID = info.I.GuildID
		member.User = user
	}
	ctx := UserExecutionContext{
		RespondingContext: info,
		User:              user,
		Member:            member,
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
