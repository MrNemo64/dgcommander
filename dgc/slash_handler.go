package dgc

import (
	"github.com/bwmarrin/discordgo"
)

var (
	ErrUnknownSlashCommand    = makeError("unknwon command: %s")
	ErrInvalidAmountOfOptions = makeError("a command group or sub command is expected to only have one option present %d are present")
	ErrInvalidSubCommandType  = makeError("the option %+v has not a valid type for a command group or sub command group")

	ErrNoFocusedArgument                 = makeError("there is no argument focused")
	ErrNoAutocompletingArgumentForOption = makeError("no autocomplete argument found for the option %+v")
)

type SlashCommandHandler func(ctx *SlashExecutionContext) error

type slashCommand interface {
	command
	autocomplete
}

type genericSlashCommand interface {
	slashCommand
	doExecute(info *RespondingContext, options []*discordgo.ApplicationCommandInteractionDataOption) (interactionAcknowledged bool, err error)
	doAutocomplete(info *InvokationInformation, options []*discordgo.ApplicationCommandInteractionDataOption) (interactionAcknowledged bool, err error)
}

type SlashSimpleMiddleware = func(ctx *SlashExecutionContext, next func()) error

type simpleSlashCommand struct {
	handler     SlashCommandHandler
	args        slashCommandArgumentListDefinition
	middlewares []SlashSimpleMiddleware
}

func (c *simpleSlashCommand) execute(info *RespondingContext) (bool, error) {
	return c.doExecute(info, info.I.ApplicationCommandData().Options)
}

func (c *simpleSlashCommand) autocomplete(info *InvokationInformation) (bool, error) {
	return c.doAutocomplete(info, info.I.ApplicationCommandData().Options)
}

func (c *simpleSlashCommand) doExecute(info *RespondingContext, options []*discordgo.ApplicationCommandInteractionDataOption) (bool, error) {
	args, err := c.args.parse(info.I.ApplicationCommandData().Resolved, options, false)
	if err != nil {
		return false, err
	}
	ctx := SlashExecutionContext{
		RespondingContext:        info,
		slashCommandArgumentList: args,
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

func (c *simpleSlashCommand) doAutocomplete(info *InvokationInformation, options []*discordgo.ApplicationCommandInteractionDataOption) (bool, error) {
	arg, err := c.findFocusedArgument(options)
	if err != nil {
		return false, err
	}
	args, err := c.args.parse(info.I.ApplicationCommandData().Resolved, options, true)
	if err != nil {
		return false, err
	}
	ctx := SlashAutocompleteContext{
		executionContext:         newExecutionContext(info.DGC.ctx, info),
		slashCommandArgumentList: args,
	}
	if err := arg.Autocomplete(&ctx); err != nil {
		return false, err
	}
	choices := ctx.makeChoices()
	err = info.Session.InteractionRespond(info.I.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
	return err != nil, err
}

func (c *simpleSlashCommand) findFocusedArgument(options []*discordgo.ApplicationCommandInteractionDataOption) (SlashCommandAutocompleteArgument, error) {
	var option *discordgo.ApplicationCommandInteractionDataOption
	for _, option = range options {
		if option.Focused {
			break
		}
	}
	if option == nil {
		return nil, ErrNoFocusedArgument
	}
	auto := c.args.autocompleteArgumentFor(option)
	if auto == nil {
		return nil, ErrNoAutocompletingArgumentForOption.withArgs(option)
	}
	return auto, nil
}

type SlashMultiMiddleware = func(ctx *RespondingContext, next func()) error

type multiSlashCommand struct {
	subCommands map[string]genericSlashCommand
	middlewares []SlashMultiMiddleware
}

func (c *multiSlashCommand) execute(info *RespondingContext) (bool, error) {
	return c.doExecute(info, info.I.ApplicationCommandData().Options)
}

func (c *multiSlashCommand) autocomplete(info *InvokationInformation) (interactionAcknowledged bool, err error) {
	return c.doAutocomplete(info, info.I.ApplicationCommandData().Options)
}

func (c *multiSlashCommand) doExecute(ctx *RespondingContext, options []*discordgo.ApplicationCommandInteractionDataOption) (bool, error) {
	command, option, err := c.findSubCommand(options)
	if err != nil {
		return false, err
	}
	mc := newMiddlewareChain(ctx, c.middlewares)
	if err := mc.startChain(); err != nil {
		return ctx.acknowledged, err
	}
	if !mc.allMiddlewaresCalled {
		return ctx.acknowledged, nil
	}
	return command.doExecute(ctx, option.Options)
}

func (c *multiSlashCommand) doAutocomplete(info *InvokationInformation, options []*discordgo.ApplicationCommandInteractionDataOption) (bool, error) {
	command, option, err := c.findSubCommand(options)
	if err != nil {
		return false, err
	}
	return command.doAutocomplete(info, option.Options)
}

func (c *multiSlashCommand) findSubCommand(options []*discordgo.ApplicationCommandInteractionDataOption) (genericSlashCommand, *discordgo.ApplicationCommandInteractionDataOption, error) {
	if len(options) != 1 {
		return nil, nil, ErrInvalidAmountOfOptions.withArgs(len(options))
	}
	option := options[0]
	if option.Type != discordgo.ApplicationCommandOptionSubCommandGroup && option.Type != discordgo.ApplicationCommandOptionSubCommand {
		return nil, nil, ErrInvalidSubCommandType.withArgs(option)
	}
	command, found := c.subCommands[option.Name]
	if !found {
		return nil, nil, ErrUnknownSlashCommand.withArgs(option.Name)
	}
	return command, option, nil
}
