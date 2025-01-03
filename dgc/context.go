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
	ctx          context.Context
	cancelCtx    context.CancelFunc
	// TODO add info about
	// - context like if running on a guild/dm
	// - channel where we are being run
	// - locale about the sender
}

func newExecutionContext(info *InvokationInformation) *executionContext {
	c, f := context.WithCancel(context.Background()) // todo receive context from background
	ctx := &executionContext{
		InvokationInformation: info,
		ctx:                   c,
		cancelCtx:             f,
	}
	ctx.timer = time.AfterFunc(ctx.EndTime().Sub(info.timeProvider.Now()), func() { ctx.cancelCtx() })
	return ctx
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

func (c *executionContext) Ctx() context.Context {
	return c.ctx
}

type respondingContext struct {
	*executionContext
	alreadyResponded bool
}

func (ctx *respondingContext) respond(resp *discordgo.InteractionResponse) error {
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

// User

type UserExecutionContext struct {
	respondingContext
	User   *discordgo.User
	Member *discordgo.Member // Nil if not running on a guild
}

// Slash

type SlashExecutionContext struct {
	respondingContext
	slashCommandArgumentList
}

type SlashAutocompleteContext struct {
	*executionContext
	slashCommandArgumentList
	choices []*discordgo.ApplicationCommandOptionChoice
}

func (ctx *SlashAutocompleteContext) AddChoice(name string, value any) *SlashAutocompleteContext {
	ctx.choices = append(ctx.choices, &discordgo.ApplicationCommandOptionChoice{Name: name, Value: value})
	return ctx
}

func (ctx *SlashAutocompleteContext) makeChoices() []*discordgo.ApplicationCommandOptionChoice {
	return ctx.choices
}
