package dgc

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

var (
	ErrUnknownSlashCommand    = makeError("unknwon command: %s")
	ErrInvalidAmountOfOptions = makeError("a command group or sub command is expected to only have one option present %d are present")
	ErrInvalidSubCommandType  = makeError("the option %+v has not a valid type for a command group or sub command group")

	ErrNoFocusedArgument                 = makeError("there is no argument focused")
	ErrNoAutocompletingArgumentForOption = makeError("no autocomplete argument found for the option %+v")
)

type SlashCommandHandler func(sender *discordgo.User, ctx *SlashExecutionContext) error

type slashCommand interface {
	command
	autocomplete
}

type genericSlashCommand interface {
	slashCommand
	doManage(log *slog.Logger, sender *discordgo.User, ss *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) (interactionAcknowledged bool, err error)
	doAutocomplete(log *slog.Logger, sender *discordgo.User, ss *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) (interactionAcknowledged bool, err error)
}

type simpleSlashCommand struct {
	handler SlashCommandHandler
	args    slashCommandArgumentListDefinition
}

func (c *simpleSlashCommand) manage(log *slog.Logger, sender *discordgo.User, ss *discordgo.Session, i *discordgo.InteractionCreate) (bool, error) {
	return c.doManage(log, sender, ss, i, i.ApplicationCommandData().Options)
}

func (c *simpleSlashCommand) autocomplete(log *slog.Logger, sender *discordgo.User, ss *discordgo.Session, i *discordgo.InteractionCreate) (bool, error) {
	return c.doAutocomplete(log, sender, ss, i, i.ApplicationCommandData().Options)
}

func (c *simpleSlashCommand) doManage(log *slog.Logger, sender *discordgo.User, ss *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) (bool, error) {
	args, err := c.args.parse(i.ApplicationCommandData().Resolved, options)
	if err != nil {
		return false, err
	}
	ctx := SlashExecutionContext{
		respondingContext: respondingContext{
			executionContext: executionContext{
				log:     log,
				Session: ss,
				I:       i,
			},
		},
		slashCommandArgumentList: args,
	}
	err = c.handler(sender, &ctx)
	return ctx.alreadyResponded, err
}

func (c *simpleSlashCommand) doAutocomplete(log *slog.Logger, sender *discordgo.User, ss *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) (bool, error) {
	arg, err := c.findFocusedArgument(options)
	if err != nil {
		return false, err
	}
	ctx := SlashAutocompleteContext{
		executionContext: executionContext{
			log:     log,
			Session: ss,
			I:       i,
		},
	}
	if err := arg.Autocomplete(&ctx); err != nil {
		return false, err
	}
	choices := ctx.makeChoices()
	err = ss.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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

type multiSlashCommand struct {
	subCommands map[string]genericSlashCommand
}

func (c *multiSlashCommand) manage(log *slog.Logger, sender *discordgo.User, ss *discordgo.Session, i *discordgo.InteractionCreate) (bool, error) {
	return c.doManage(log, sender, ss, i, i.ApplicationCommandData().Options)
}

func (c *multiSlashCommand) autocomplete(log *slog.Logger, sender *discordgo.User, ss *discordgo.Session, i *discordgo.InteractionCreate) (interactionAcknowledged bool, err error) {
	return c.doAutocomplete(log, sender, ss, i, i.ApplicationCommandData().Options)
}

func (c *multiSlashCommand) doManage(log *slog.Logger, sender *discordgo.User, ss *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) (bool, error) {
	command, option, err := c.findSubCommand(options)
	if err != nil {
		return false, err
	}
	return command.doManage(log, sender, ss, i, option.Options)
}

func (c *multiSlashCommand) doAutocomplete(log *slog.Logger, sender *discordgo.User, ss *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) (bool, error) {
	command, option, err := c.findSubCommand(options)
	if err != nil {
		return false, err
	}
	return command.doAutocomplete(log, sender, ss, i, option.Options)
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
