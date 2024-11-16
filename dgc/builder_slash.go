package dgc

import (
	"github.com/MrNemo64/dgcommander/dgc/handlers"
	"github.com/bwmarrin/discordgo"
)

type slashBuilder[B any] struct {
	commandBuilder[handlers.SlashHandler, B]
	description string
}

func (b *slashBuilder[B]) Description(description string) B {
	b.description = description
	return b.upper
}

func (b *slashBuilder[B]) discordDefineForCreation() *discordgo.ApplicationCommand {
	c := b.commandBuilder.discordDefineForCreation()
	c.Type = discordgo.ChatApplicationCommand
	c.Description = b.description
	return c
}

type slashSimpleBuilder struct {
	slashBuilder[*slashSimpleBuilder]
	argumentListBuilder[*slashSimpleBuilder]
}

func NewSimpleSlash() *slashSimpleBuilder {
	b := &slashSimpleBuilder{
		slashBuilder: slashBuilder[*slashSimpleBuilder]{
			commandBuilder: commandBuilder[handlers.SlashHandler, *slashSimpleBuilder]{},
		},
		argumentListBuilder: argumentListBuilder[*slashSimpleBuilder]{},
	}
	b.slashBuilder.commandBuilder.upper = b
	b.argumentListBuilder.upper = b
	return b
}

func (b *slashSimpleBuilder) create() command {
	panic("not implemented")
}

func (b *slashSimpleBuilder) discordDefineForCreation() *discordgo.ApplicationCommand {
	c := b.slashBuilder.discordDefineForCreation()
	c.Options = b.argumentListBuilder.discordDefineForCreation()
	return c
}

type slashComplexBuilder struct {
	slashBuilder[*slashComplexBuilder]
	subCommands []subCommandLike
}

func NewComplexSlash() *slashComplexBuilder {
	b := &slashComplexBuilder{
		slashBuilder: slashBuilder[*slashComplexBuilder]{
			commandBuilder: commandBuilder[handlers.SlashHandler, *slashComplexBuilder]{},
		},
	}
	b.slashBuilder.commandBuilder.upper = b
	return b
}

func (b *slashComplexBuilder) AddSubCommandGroup(group *slashSubCommandGroupBuilder) *slashComplexBuilder {
	b.subCommands = append(b.subCommands, group)
	return b
}

func (b *slashComplexBuilder) AddSubCommand(command *slashSubCommandBuilder) *slashComplexBuilder {
	b.subCommands = append(b.subCommands, command)
	return b
}

func (b *slashComplexBuilder) create() command {
	panic("not implemented")
}

func (b *slashComplexBuilder) discordDefineForCreation() *discordgo.ApplicationCommand {
	c := b.slashBuilder.discordDefineForCreation()
	c.Options = make([]*discordgo.ApplicationCommandOption, len(b.subCommands))
	for i, subCommand := range b.subCommands {
		c.Options[i] = subCommand.discordDefineForCreation()
	}
	return c
}

type subCommandLike interface {
	isSubCommandLike() bool
	discordDefineForCreation() *discordgo.ApplicationCommandOption
}

type slashSubCommandBuilder struct {
	argumentListBuilder[*slashSubCommandBuilder]
	name        string
	description string
}

func NewSubCommand() *slashSubCommandBuilder {
	b := &slashSubCommandBuilder{argumentListBuilder: argumentListBuilder[*slashSubCommandBuilder]{}}
	b.argumentListBuilder.upper = b
	return b
}

func (*slashSubCommandBuilder) isSubCommandLike() bool { return true }

func (b *slashSubCommandBuilder) Name(name string) *slashSubCommandBuilder {
	b.name = name
	return b
}

func (b *slashSubCommandBuilder) Description(description string) *slashSubCommandBuilder {
	b.description = description
	return b
}

func (b *slashSubCommandBuilder) discordDefineForCreation() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        b.name,
		Description: b.description,
		Options:     b.argumentListBuilder.discordDefineForCreation(),
	}
}

type slashSubCommandGroupBuilder struct {
	subCommands []*slashSubCommandBuilder
	name        string
	description string
}

func NewSubCommandGroup() *slashSubCommandGroupBuilder {
	return &slashSubCommandGroupBuilder{}
}

func (*slashSubCommandGroupBuilder) isSubCommandLike() bool { return true }

func (b *slashSubCommandGroupBuilder) Name(name string) *slashSubCommandGroupBuilder {
	b.name = name
	return b
}

func (b *slashSubCommandGroupBuilder) Description(description string) *slashSubCommandGroupBuilder {
	b.description = description
	return b
}

func (b *slashSubCommandGroupBuilder) AddSubCommand(command *slashSubCommandBuilder) *slashSubCommandGroupBuilder {
	b.subCommands = append(b.subCommands, command)
	return b
}

func (b *slashSubCommandGroupBuilder) AddSubCommands(commands ...*slashSubCommandBuilder) *slashSubCommandGroupBuilder {
	b.subCommands = append(b.subCommands, commands...)
	return b
}

func (b *slashSubCommandGroupBuilder) discordDefineForCreation() *discordgo.ApplicationCommandOption {
	c := &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
		Name:        b.name,
		Description: b.description,
		Options:     make([]*discordgo.ApplicationCommandOption, len(b.subCommands)),
	}
	for i, subCommand := range b.subCommands {
		c.Options[i] = subCommand.discordDefineForCreation()
	}
	return c
}
