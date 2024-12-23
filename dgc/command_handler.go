package dgc

import (
	"errors"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

var (
	ErrAlreadyResponded = errors.New("already responded, send a follow up")
)

type executionContext struct {
	Session          *discordgo.Session
	I                *discordgo.InteractionCreate
	log              *slog.Logger
	alreadyResponded bool
}

func (ctx *executionContext) respond(resp *discordgo.InteractionResponse) error {
	if ctx.alreadyResponded {
		return ErrAlreadyResponded
	}
	if err := ctx.Session.InteractionRespond(ctx.I.Interaction, resp); err != nil {
		return err
	}
	ctx.alreadyResponded = true
	return nil
}

func (ctx *executionContext) RespondWithMessage(message *discordgo.InteractionResponseData) error {
	return ctx.respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: message,
	})
}

func (ctx *executionContext) RespondWithModal(modal *discordgo.InteractionResponseData) error {
	return ctx.respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: modal,
	})
}

func (ctx *executionContext) RespondLatter() error {
	return ctx.respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: nil,
	})
}

func (ctx *executionContext) AddFollowup(wait bool, data *discordgo.WebhookParams) (*discordgo.Message, error) {
	return ctx.Session.FollowupMessageCreate(ctx.I.Interaction, wait, data)
}
