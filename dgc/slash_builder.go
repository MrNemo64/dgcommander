package dgc

import (
	"github.com/bwmarrin/discordgo"
)

type genericSlashCommandBuilder[B specificCommandBuilder] struct {
	genericCommandBuilder[B]
	description string
}

func (b *genericSlashCommandBuilder[B]) discordDefineForCreation() *discordgo.ApplicationCommand {
	c := b.genericCommandBuilder.discordDefineForCreation()
	c.Type = discordgo.ChatApplicationCommand
	c.Description = b.description
	return c
}

func (b *genericSlashCommandBuilder[B]) Description(description string) B {
	b.description = description
	return b.upper
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
	b.slashCommandArgumentListBuilder.upper = b
	return b
}

func (b *SimpleSlashCommandBuilder) discordDefineForCreation() *discordgo.ApplicationCommand {
	c := b.genericCommandBuilder.discordDefineForCreation()
	c.Options = b.slashCommandArgumentListBuilder.discordDefineForCreation()
	return c
}

func (b *SimpleSlashCommandBuilder) create() command {
	panic("TODO")
}

func (b *SimpleSlashCommandBuilder) Handler(handler SlashCommandHandler) *SimpleSlashCommandBuilder {
	b.handler = handler
	return b
}

// Multi

type subcommandLikeBuilder interface {
	discordDefineForCreation() *discordgo.ApplicationCommandOption
}

type MultiSlashCommandBuilder struct {
	genericSlashCommandBuilder[*MultiSlashCommandBuilder]
	subCommands []subcommandLikeBuilder
}

func NewMultiSlashCommandBuilder() *MultiSlashCommandBuilder {
	b := &MultiSlashCommandBuilder{}
	b.genericSlashCommandBuilder.genericCommandBuilder.upper = b
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
	panic("TODO")
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
	name        string
	description string
}

func NewSubCommand() *SubSlashCommandBuilder {
	b := &SubSlashCommandBuilder{}
	b.slashCommandArgumentListBuilder.upper = b
	return b
}

func (b *SubSlashCommandBuilder) discordDefineForCreation() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        b.name,
		Description: b.description,
		Options:     b.slashCommandArgumentListBuilder.discordDefineForCreation(),
	}
}

func (b *SubSlashCommandBuilder) Name(name string) *SubSlashCommandBuilder {
	b.name = name
	return b
}

func (b *SubSlashCommandBuilder) Description(description string) *SubSlashCommandBuilder {
	b.description = description
	return b
}

func (b *SubSlashCommandBuilder) Handler(handler SlashCommandHandler) *SubSlashCommandBuilder {
	b.handler = handler
	return b
}

type SubSlashCommandGroupBuilder struct {
	commands    []*SubSlashCommandBuilder
	name        string
	description string
}

func NewSubCommandGroup() *SubSlashCommandGroupBuilder {
	return &SubSlashCommandGroupBuilder{}
}

func (b *SubSlashCommandGroupBuilder) discordDefineForCreation() *discordgo.ApplicationCommandOption {
	c := &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
		Name:        b.name,
		Description: b.description,
		Options:     make([]*discordgo.ApplicationCommandOption, len(b.commands)),
	}
	for i, subCommand := range b.commands {
		c.Options[i] = subCommand.discordDefineForCreation()
	}
	return c
}

func (b *SubSlashCommandGroupBuilder) Name(name string) *SubSlashCommandGroupBuilder {
	b.name = name
	return b
}

func (b *SubSlashCommandGroupBuilder) Description(description string) *SubSlashCommandGroupBuilder {
	b.description = description
	return b
}

func (b *SubSlashCommandGroupBuilder) AddSubCommand(command *SubSlashCommandBuilder) *SubSlashCommandGroupBuilder {
	b.commands = append(b.commands, command)
	return b
}

func (b *SubSlashCommandGroupBuilder) AddSubCommands(commands ...*SubSlashCommandBuilder) *SubSlashCommandGroupBuilder {
	b.commands = append(b.commands, commands...)
	return b
}
