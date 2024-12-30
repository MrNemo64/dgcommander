package dgc

import (
	"errors"
)

var (
	ErrUserCommandHasNoUser = errors.New("User command was used but no user was given")
)

type UserCommandHandler func(ctx *UserExecutionContext) error

type userCommand struct {
	handler UserCommandHandler
}

func (c *userCommand) execute(info *InvokationInformation) (bool, error) {
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
		respondingContext: respondingContext{
			executionContext: newExecutionContext(info.DGC.ctx, info),
		},
		User:   user,
		Member: member,
	}
	err := c.handler(&ctx)
	return ctx.alreadyResponded, err
}
