package handlers

import (
	"errors"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

type ExecutionContext struct {
	Session          *discordgo.Session
	I                *discordgo.InteractionCreate
	log              *slog.Logger
	alreadyResponded bool
}

func (ctx *ExecutionContext) respond(resp *discordgo.InteractionResponse) error {
	if ctx.alreadyResponded {
		return errors.New("already responded, send a follow up")
	}
	if err := ctx.Session.InteractionRespond(ctx.I.Interaction, resp); err != nil {
		return err
	}
	ctx.alreadyResponded = true
	return nil
}

func (ctx *ExecutionContext) RespondWithMessage(message *discordgo.InteractionResponseData) error {
	return ctx.respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: message,
	})
}

func (ctx *ExecutionContext) RespondWithModal(modal *discordgo.InteractionResponseData) error {
	return ctx.respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: modal,
	})
}

func (ctx *ExecutionContext) RespondLatter() error {
	return ctx.respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: nil,
	})
}

func (ctx *ExecutionContext) AddFollowup(wait bool, data *discordgo.WebhookParams) (*discordgo.Message, error) {
	return ctx.Session.FollowupMessageCreate(ctx.I.Interaction, wait, data)
}
