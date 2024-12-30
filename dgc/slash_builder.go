package dgc

import (
	"github.com/MrNemo64/dgcommander/dgc/util"
	"github.com/bwmarrin/discordgo"
)

type genericSlashCommandBuilder[B specificCommandBuilder] struct {
	genericCommandBuilder[B]
	description util.Localizable[B]
}

func (b *genericSlashCommandBuilder[B]) discordDefineForCreation() *discordgo.ApplicationCommand {
	c := b.genericCommandBuilder.discordDefineForCreation()
	c.Type = discordgo.ChatApplicationCommand
	c.Description = b.description.Value
	var descriptionLocalizations *map[discordgo.Locale]string
	if b.description.Localizations != nil {
		descriptionLocalizations = &b.description.Localizations
	}
	c.DescriptionLocalizations = descriptionLocalizations
	return c
}

func (b *genericSlashCommandBuilder[B]) Description() *util.Localizable[B] {
	return &b.description
}

// Simple

type SimpleSlashCommandBuilder struct {
	genericSlashCommandBuilder[*SimpleSlashCommandBuilder]
	slashCommandArgumentListBuilder[*SimpleSlashCommandBuilder]
	handler SlashCommandHandler
}

func NewSimpleSlashCommandBuilder() *SimpleSlashCommandBuilder {
	b := &SimpleSlashCommandBuilder{}
	b.genericSlashCommandBuilder.genericCommandBuilder.upper = b
	b.genericSlashCommandBuilder.genericCommandBuilder.name.Upper = b
	b.genericSlashCommandBuilder.description.Upper = b
	b.slashCommandArgumentListBuilder.upper = b
	return b
}

func (b *SimpleSlashCommandBuilder) discordDefineForCreation() *discordgo.ApplicationCommand {
	c := b.genericCommandBuilder.discordDefineForCreation()
	c.Options = b.slashCommandArgumentListBuilder.discordDefineForCreation()
	return c
}

func (b *SimpleSlashCommandBuilder) create() command {
	return &simpleSlashCommand{
		handler: b.handler,
		args:    b.slashCommandArgumentListBuilder.create(),
	}
}

func (b *SimpleSlashCommandBuilder) Handler(handler SlashCommandHandler) *SimpleSlashCommandBuilder {
	b.handler = handler
	return b
}

// Multi

type subcommandLikeBuilder interface {
	discordDefineForCreation() *discordgo.ApplicationCommandOption
	create() (string, genericSlashCommand)
}

type MultiSlashCommandBuilder struct {
	genericSlashCommandBuilder[*MultiSlashCommandBuilder]
	subCommands []subcommandLikeBuilder
}

func NewMultiSlashCommandBuilder() *MultiSlashCommandBuilder {
	b := &MultiSlashCommandBuilder{}
	b.genericSlashCommandBuilder.genericCommandBuilder.upper = b
	b.genericSlashCommandBuilder.genericCommandBuilder.name.Upper = b
	b.genericSlashCommandBuilder.description.Upper = b
	return b
}

func (b *MultiSlashCommandBuilder) discordDefineForCreation() *discordgo.ApplicationCommand {
	c := b.genericSlashCommandBuilder.discordDefineForCreation()
	c.Options = make([]*discordgo.ApplicationCommandOption, len(b.subCommands))
	for i, subCommand := range b.subCommands {
		c.Options[i] = subCommand.discordDefineForCreation()
	}
	return c
}

func (b *MultiSlashCommandBuilder) create() command {
	sub := make(map[string]genericSlashCommand)
	for _, v := range b.subCommands {
		name, command := v.create()
		sub[name] = command
	}
	return &multiSlashCommand{
		subCommands: sub,
	}
}

func (b *MultiSlashCommandBuilder) AddSubCommandGroup(group *SubSlashCommandGroupBuilder) *MultiSlashCommandBuilder {
	b.subCommands = append(b.subCommands, group)
	return b
}

func (b *MultiSlashCommandBuilder) AddSubCommand(command *SubSlashCommandBuilder) *MultiSlashCommandBuilder {
	b.subCommands = append(b.subCommands, command)
	return b
}

type SubSlashCommandBuilder struct {
	slashCommandArgumentListBuilder[*SubSlashCommandBuilder]
	handler     SlashCommandHandler
	name        util.Localizable[*SubSlashCommandBuilder]
	description util.Localizable[*SubSlashCommandBuilder]
}

func NewSubCommand() *SubSlashCommandBuilder {
	b := &SubSlashCommandBuilder{}
	b.slashCommandArgumentListBuilder.upper = b
	b.name.Upper = b
	b.description.Upper = b
	return b
}

func (b *SubSlashCommandBuilder) discordDefineForCreation() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionSubCommand,
		Name:                     b.name.Value,
		NameLocalizations:        b.name.Localizations,
		Description:              b.description.Value,
		DescriptionLocalizations: b.description.Localizations,
		Options:                  b.slashCommandArgumentListBuilder.discordDefineForCreation(),
	}
}

func (b *SubSlashCommandBuilder) create() (string, genericSlashCommand) {
	return b.name.Value, &simpleSlashCommand{
		handler: b.handler,
		args:    b.slashCommandArgumentListBuilder.create(),
	}
}

func (b *SubSlashCommandBuilder) Name() *util.Localizable[*SubSlashCommandBuilder] {
	return &b.name
}

func (b *SubSlashCommandBuilder) Description() *util.Localizable[*SubSlashCommandBuilder] {
	return &b.description
}

func (b *SubSlashCommandBuilder) Handler(handler SlashCommandHandler) *SubSlashCommandBuilder {
	b.handler = handler
	return b
}

type SubSlashCommandGroupBuilder struct {
	commands    []*SubSlashCommandBuilder
	name        util.Localizable[*SubSlashCommandGroupBuilder]
	description util.Localizable[*SubSlashCommandGroupBuilder]
}

func NewSubCommandGroup() *SubSlashCommandGroupBuilder {
	b := &SubSlashCommandGroupBuilder{}
	b.name.Upper = b
	b.description.Upper = b
	return b
}

func (b *SubSlashCommandGroupBuilder) discordDefineForCreation() *discordgo.ApplicationCommandOption {
	c := &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionSubCommandGroup,
		Name:                     b.name.Value,
		NameLocalizations:        b.name.Localizations,
		Description:              b.description.Value,
		DescriptionLocalizations: b.description.Localizations,
		Options:                  make([]*discordgo.ApplicationCommandOption, len(b.commands)),
	}
	for i, subCommand := range b.commands {
		c.Options[i] = subCommand.discordDefineForCreation()
	}
	return c
}

func (b *SubSlashCommandGroupBuilder) create() (string, genericSlashCommand) {
	sub := make(map[string]genericSlashCommand)
	for _, v := range b.commands {
		name, command := v.create()
		sub[name] = command
	}
	return b.name.Value, &multiSlashCommand{
		subCommands: sub,
	}
}

func (b *SubSlashCommandGroupBuilder) Name() *util.Localizable[*SubSlashCommandGroupBuilder] {
	return &b.name
}

func (b *SubSlashCommandGroupBuilder) Description() *util.Localizable[*SubSlashCommandGroupBuilder] {
	return &b.description
}

func (b *SubSlashCommandGroupBuilder) AddSubCommand(command *SubSlashCommandBuilder) *SubSlashCommandGroupBuilder {
	b.commands = append(b.commands, command)
	return b
}

func (b *SubSlashCommandGroupBuilder) AddSubCommands(commands ...*SubSlashCommandBuilder) *SubSlashCommandGroupBuilder {
	b.commands = append(b.commands, commands...)
	return b
}
