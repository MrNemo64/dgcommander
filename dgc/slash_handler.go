package dgc

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

var (
	ErrUnknownSlashCommand    = makeError("unknwon command: %s")
	ErrInvalidAmountOfOptions = makeError("a command group or sub command is expected to only have one option present %d are present")
	ErrInvalidSubCommandType  = makeError("the option %+v has not a valid type for a command group or sub command group")
)

type SlashExecutionContext struct {
	executionContext
	slashCommandArgumentList
}

type SlashCommandHandler func(sender *discordgo.User, ctx *SlashExecutionContext) error

type genericSlashCommand interface {
	manage(log *slog.Logger, sender *discordgo.User, ss *discordgo.Session, i *discordgo.InteractionCreate) (bool, error)
	execute(log *slog.Logger, sender *discordgo.User, ss *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) (interactionAcknowledged bool, err error)
}

type simpleSlashCommand struct {
	handler SlashCommandHandler
	args    slashCommandArgumentListDefinition
}

func (c *simpleSlashCommand) manage(log *slog.Logger, sender *discordgo.User, ss *discordgo.Session, i *discordgo.InteractionCreate) (bool, error) {
	return c.execute(log, sender, ss, i, i.ApplicationCommandData().Options)
}

func (c *simpleSlashCommand) execute(log *slog.Logger, sender *discordgo.User, ss *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) (bool, error) {
	args, err := c.args.parse(i.ApplicationCommandData().Resolved, options)
	if err != nil {
		return false, err
	}
	ctx := SlashExecutionContext{
		executionContext: executionContext{
			log:     log,
			Session: ss,
			I:       i,
		},
		slashCommandArgumentList: args,
	}
	err = c.handler(sender, &ctx)
	return ctx.alreadyResponded, err
}

type multiSlashCommand struct {
	subCommands map[string]genericSlashCommand
}

func (c *multiSlashCommand) manage(log *slog.Logger, sender *discordgo.User, ss *discordgo.Session, i *discordgo.InteractionCreate) (bool, error) {
	return c.execute(log, sender, ss, i, i.ApplicationCommandData().Options)
}

func (c *multiSlashCommand) execute(log *slog.Logger, sender *discordgo.User, ss *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) (bool, error) {
	if len(options) != 1 {
		return false, ErrInvalidAmountOfOptions.withArgs(len(options))
	}
	option := options[0]
	if option.Type != discordgo.ApplicationCommandOptionSubCommandGroup && option.Type != discordgo.ApplicationCommandOptionSubCommand {
		return false, ErrInvalidSubCommandType.withArgs(option)
	}
	command, found := c.subCommands[option.Name]
	if !found {
		return false, ErrUnknownSlashCommand.withArgs(option.Name)
	}
	return command.execute(log, sender, ss, i, option.Options)
}
