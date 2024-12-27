package dgc

import (
	"errors"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

// General

var (
	ErrAlreadyResponded = errors.New("already responded, send a follow up")
)

type executionContext struct {
	Session *discordgo.Session
	I       *discordgo.InteractionCreate
	log     *slog.Logger

	// TODO add things like time for the interaction to end
}

type respondingContext struct {
	executionContext
	alreadyResponded bool
}

func (ctx *respondingContext) respond(resp *discordgo.InteractionResponse) error {
	if ctx.alreadyResponded {
		return ErrAlreadyResponded
	}
	if err := ctx.Session.InteractionRespond(ctx.I.Interaction, resp); err != nil {
		return err
	}
	ctx.alreadyResponded = true
	return nil
}

func (ctx *respondingContext) RespondWithMessage(message *discordgo.InteractionResponseData) error {
	return ctx.respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: message,
	})
}

func (ctx *respondingContext) RespondWithModal(modal *discordgo.InteractionResponseData) error {
	return ctx.respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: modal,
	})
}

func (ctx *respondingContext) RespondLatter() error {
	return ctx.respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: nil,
	})
}

func (ctx *respondingContext) AddFollowup(wait bool, data *discordgo.WebhookParams) (*discordgo.Message, error) {
	return ctx.Session.FollowupMessageCreate(ctx.I.Interaction, wait, data)
}

// Message

type MessageExecutionContext struct {
	respondingContext
	Message *discordgo.Message
}

// Slash

type SlashExecutionContext struct {
	respondingContext
	slashCommandArgumentList
}

type SlashAutocompleteContext struct {
	executionContext
	choices []*discordgo.ApplicationCommandOptionChoice
}

func (ctx *SlashAutocompleteContext) AddChoice(name string, value any) *SlashAutocompleteContext {
	ctx.choices = append(ctx.choices, &discordgo.ApplicationCommandOptionChoice{Name: name, Value: value})
	return ctx
}

func (ctx *SlashAutocompleteContext) makeChoices() []*discordgo.ApplicationCommandOptionChoice {
	return ctx.choices
}
