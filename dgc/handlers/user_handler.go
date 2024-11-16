package handlers

import "github.com/bwmarrin/discordgo"

type UserHandler func(ctx *UserExecutionContext, sender *discordgo.User) error

type UserExecutionContext struct {
	ExecutionContext
}
