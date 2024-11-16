package handlers

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

type SlashHandler func(sender *discordgo.User, ctx *SlashExecutionContext) error

type SlashExecutionContext struct {
	Session *discordgo.Session
	I       *discordgo.InteractionCreate
	log     *slog.Logger
	args    *CommandArguments
}

type SlashSimpleCommand struct {
	args    ArgumentList
	handler SlashHandler
}

func (c *SlashSimpleCommand) manage(log *slog.Logger, sender *discordgo.User, ss *discordgo.Session, i *discordgo.InteractionCreate) error {
	return c.execute(log, sender, ss, i, i.ApplicationCommandData().Options)
}

func (c *SlashSimpleCommand) execute(log *slog.Logger, sender *discordgo.User, ss *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) error {
	args, err := c.args.ParseOptions(options)
	if err != nil {
		return err
	}
	ctx := SlashExecutionContext{
		Session: ss,
		I:       i,
		log:     log,
		args:    &args,
	}
	return c.handler(sender, &ctx)
}
