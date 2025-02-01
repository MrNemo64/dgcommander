package dgc

import (
	"context"
	"errors"
	"time"

	"github.com/bwmarrin/discordgo"
)

// General

var (
	ErrAlreadyResponded = errors.New("already responded, send a follow up")
)

type executionContext struct {
	*InvokationInformation
	acknowledged bool
	timer        *time.Timer
	Ctx          context.Context
	cancelCtx    context.CancelFunc
	// TODO add info about
	// - context like if running on a guild/dm
	// - channel where we are being run
	// - locale about the sender
}

func newExecutionContext(ctx context.Context, info *InvokationInformation) *executionContext {
	c, f := context.WithCancel(ctx)
	ectx := &executionContext{
		InvokationInformation: info,
		Ctx:                   c,
		cancelCtx:             f,
	}
	ectx.timer = time.AfterFunc(ectx.EndTime().Sub(info.timeProvider.Now()), func() { ectx.cancelCtx() })
	return ectx
}

func (ctx *executionContext) Finish() {
	if !ctx.timer.Stop() {
		ctx.cancelCtx()
	}
}

func (ctx *executionContext) EndTime() time.Time {
	// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-callback
	if ctx.acknowledged {
		return ctx.ReceivedTime.Add(15 * time.Minute)
	}
	return ctx.ReceivedTime.Add(3 * time.Second)
}

type RespondingContext struct {
	*executionContext
	alreadyResponded bool
}

func (ctx *RespondingContext) respond(resp *discordgo.InteractionResponse) error {
	if ctx.alreadyResponded {
		return ErrAlreadyResponded
	}
	ctx.acknowledged = true
	ctx.timer.Reset(ctx.EndTime().Sub(ctx.timeProvider.Now()))
	if err := ctx.Session.InteractionRespond(ctx.I.Interaction, resp); err != nil {
		return err
	}
	ctx.alreadyResponded = true
	return nil
}

func (ctx *RespondingContext) RespondWithMessage(message *discordgo.InteractionResponseData) error {
	return ctx.respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: message,
	})
}

func (ctx *RespondingContext) RespondWithModal(modal *discordgo.InteractionResponseData) error {
	return ctx.respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: modal,
	})
}

func (ctx *RespondingContext) RespondLatter() error {
	return ctx.respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: nil,
	})
}

func (ctx *RespondingContext) AddFollowup(wait bool, data *discordgo.WebhookParams) (*discordgo.Message, error) {
	return ctx.Session.FollowupMessageCreate(ctx.I.Interaction, wait, data)
}

// Message

type MessageExecutionContext struct {
	*RespondingContext
	Message *discordgo.Message
}

// User

type UserExecutionContext struct {
	*RespondingContext
	User   *discordgo.User
	Member *discordgo.Member // Nil if not running on a guild
}

// Slash

type SlashExecutionContext struct {
	*RespondingContext
	slashCommandArgumentList
}

type SlashAutocompleteContext struct {
	*executionContext
	slashCommandArgumentList
	choices []*discordgo.ApplicationCommandOptionChoice
}

func (ctx *SlashAutocompleteContext) AddChoice(name string, value any) *SlashAutocompleteContext {
	ctx.choices = append(ctx.choices, &discordgo.ApplicationCommandOptionChoice{Name: name, Value: value, NameLocalizations: nil})
	return ctx
}

func (ctx *SlashAutocompleteContext) AddLocalizedChoice(name string, value any, localizations map[discordgo.Locale]string) *SlashAutocompleteContext {
	ctx.choices = append(ctx.choices, &discordgo.ApplicationCommandOptionChoice{Name: name, Value: value, NameLocalizations: localizations})
	return ctx
}

func (ctx *SlashAutocompleteContext) makeChoices() []*discordgo.ApplicationCommandOptionChoice {
	return ctx.choices
}
