package handlers

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

type SlashHandler func(sender *discordgo.User, ctx *SlashExecutionContext) error

type SlashCommand interface {
	Manage(log *slog.Logger, sender *discordgo.User, ss *discordgo.Session, i *discordgo.InteractionCreate) (interactionAcknowledged bool, err error)
	execute(log *slog.Logger, sender *discordgo.User, ss *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) (interactionAcknowledged bool, err error)
}

type SlashExecutionContext struct {
	ExecutionContext
	args *CommandArguments
}

type SlashSimpleCommand struct {
	args    ArgumentList
	handler SlashHandler
}

func NewSlashSimpleCommand(args ArgumentList, handler SlashHandler) *SlashSimpleCommand {
	return &SlashSimpleCommand{
		args:    args,
		handler: handler,
	}
}

func (c *SlashSimpleCommand) Manage(log *slog.Logger, sender *discordgo.User, ss *discordgo.Session, i *discordgo.InteractionCreate) (bool, error) {
	return c.execute(log, sender, ss, i, i.ApplicationCommandData().Options)
}

func (b *SlashSimpleCommand) Autocomplete(log *slog.Logger, sender *discordgo.User, ss *discordgo.Session, i *discordgo.InteractionCreate) error {
	return errors.New("return not implemented")
}

func (c *SlashSimpleCommand) execute(log *slog.Logger, sender *discordgo.User, ss *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) (bool, error) {
	args, err := c.args.ParseOptions(options)
	if err != nil {
		return false, err
	}
	ctx := SlashExecutionContext{
		ExecutionContext: ExecutionContext{
			log:     log,
			Session: ss,
			I:       i,
		},
		args: &args,
	}
	err = c.handler(sender, &ctx)
	return ctx.alreadyResponded, err
}

type SlashGroupCommand struct {
	subCommands map[string]SlashCommand
}

func NewSlashGroupCommand(subCommands map[string]SlashCommand) *SlashGroupCommand {
	return &SlashGroupCommand{subCommands: subCommands}
}

func (c *SlashGroupCommand) Manage(log *slog.Logger, sender *discordgo.User, ss *discordgo.Session, i *discordgo.InteractionCreate) (bool, error) {
	return c.execute(log, sender, ss, i, i.ApplicationCommandData().Options)
}

func (b *SlashGroupCommand) Autocomplete(log *slog.Logger, sender *discordgo.User, ss *discordgo.Session, i *discordgo.InteractionCreate) error {
	return errors.New("message commands cannot be autocompleted")
}

func (c *SlashGroupCommand) execute(log *slog.Logger, sender *discordgo.User, ss *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) (bool, error) {
	if len(options) != 1 {
		return false, errors.New("a command group or sub command is expected to only have one option present")
	}
	option := options[0]
	if option.Type != discordgo.ApplicationCommandOptionSubCommandGroup && option.Type != discordgo.ApplicationCommandOptionSubCommand {
		return false, fmt.Errorf("the option %+v has not a valid type for a command group or sub command group", option)
	}
	command, found := c.subCommands[option.Name]
	if !found {
		return false, errors.New("unknown command")
	}
	return command.execute(log, sender, ss, i, option.Options)
}
