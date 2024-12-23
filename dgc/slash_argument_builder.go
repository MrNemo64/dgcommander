package dgc

import (
	"github.com/bwmarrin/discordgo"
)

type slashCommandArgumentBuilder interface {
	discordDefineForCreation() *discordgo.ApplicationCommandOption
	create() (bool, string, SlashCommandArgument)
}

type specificSlashCommandArgumentBuilder interface {
	create() SlashCommandArgument
}

type genericSlashCommandArgumentBuilder[B specificSlashCommandArgumentBuilder] struct {
	upper       B
	kind        discordgo.ApplicationCommandOptionType
	name        string
	description string
	required    bool
}

func (arg *genericSlashCommandArgumentBuilder[B]) discordDefineForCreation() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        arg.kind,
		Name:        arg.name,
		Description: arg.description,
		Required:    arg.required,
	}
}

func (arg *genericSlashCommandArgumentBuilder[B]) create() (bool, string, SlashCommandArgument) {
	return arg.required, arg.name, arg.upper.create()
}

func (b *genericSlashCommandArgumentBuilder[B]) Name(name string) B {
	b.name = name
	return b.upper
}

func (b *genericSlashCommandArgumentBuilder[B]) Description(description string) B {
	b.description = description
	return b.upper
}

func (b *genericSlashCommandArgumentBuilder[B]) Required(required bool) B {
	b.required = required
	return b.upper
}

// Boolean

type booleanArgumentBuilder struct {
	genericSlashCommandArgumentBuilder[*booleanArgumentBuilder]
}

func NewBooleanArgument() *booleanArgumentBuilder {
	b := &booleanArgumentBuilder{}
	b.genericSlashCommandArgumentBuilder.upper = b
	return b
}

func (b *booleanArgumentBuilder) create() SlashCommandArgument {
	panic("TODO")
}

// String

type stringArgumentBuilder struct {
	genericSlashCommandArgumentBuilder[*stringArgumentBuilder]
	minLength *int
	maxLength int
}

func NewStringArgument() *stringArgumentBuilder {
	b := &stringArgumentBuilder{}
	b.genericSlashCommandArgumentBuilder.upper = b
	return b
}

func (b *stringArgumentBuilder) discordDefineForCreation() *discordgo.ApplicationCommandOption {
	o := b.genericSlashCommandArgumentBuilder.discordDefineForCreation()
	if b.minLength != nil {
		o.MinLength = b.minLength
	}
	o.MaxLength = b.maxLength
	return o
}

func (b *stringArgumentBuilder) create() SlashCommandArgument {
	panic("TODO")
}

// 0-6000
func (b *stringArgumentBuilder) MinLength(min int) *stringArgumentBuilder {
	b.minLength = &min
	return b
}

// 1-6000
func (b *stringArgumentBuilder) MaxLength(max int) *stringArgumentBuilder {
	b.maxLength = max
	return b
}

// Integer

type integerArgumentBuilder struct {
	genericSlashCommandArgumentBuilder[*integerArgumentBuilder]
	minValue *int
	maxValue int
}

func NewIntegerArgument() *integerArgumentBuilder {
	b := &integerArgumentBuilder{}
	b.genericSlashCommandArgumentBuilder.upper = b
	return b
}

func (b *integerArgumentBuilder) discordDefineForCreation() *discordgo.ApplicationCommandOption {
	o := b.genericSlashCommandArgumentBuilder.discordDefineForCreation()
	if b.minValue != nil {
		fv := float64(*b.minValue)
		o.MinValue = &fv
	}
	o.MaxValue = float64(b.maxValue)
	return o
}

func (b *integerArgumentBuilder) create() SlashCommandArgument {
	panic("TODO")
}

func (b *integerArgumentBuilder) MinValue(min int) *integerArgumentBuilder {
	b.minValue = &min
	return b
}

func (b *integerArgumentBuilder) MaxValue(max int) *integerArgumentBuilder {
	b.maxValue = max
	return b
}

// Number

type numberArgumentBuilder struct {
	genericSlashCommandArgumentBuilder[*numberArgumentBuilder]
	minValue *float64
	maxValue float64
}

func NewNumberArgument() *numberArgumentBuilder {
	b := &numberArgumentBuilder{}
	b.genericSlashCommandArgumentBuilder.upper = b
	return b
}

func (b *numberArgumentBuilder) discordDefineForCreation() *discordgo.ApplicationCommandOption {
	o := b.genericSlashCommandArgumentBuilder.discordDefineForCreation()
	if b.minValue != nil {
		o.MinValue = b.minValue
	}
	o.MaxValue = b.maxValue
	return o
}

func (b *numberArgumentBuilder) create() SlashCommandArgument {
	panic("TODO")
}

func (b *numberArgumentBuilder) MinValue(min float64) *numberArgumentBuilder {
	b.minValue = &min
	return b
}

func (b *numberArgumentBuilder) MaxValue(max float64) *numberArgumentBuilder {
	b.maxValue = max
	return b
}

// User

type userArgumentBuilder struct {
	genericSlashCommandArgumentBuilder[*userArgumentBuilder]
}

func NewUserArgument() *userArgumentBuilder {
	b := &userArgumentBuilder{}
	b.genericSlashCommandArgumentBuilder.upper = b
	return b
}

func (b *userArgumentBuilder) create() SlashCommandArgument {
	panic("TODO")
}

// Role

type roleArgumentBuilder struct {
	genericSlashCommandArgumentBuilder[*roleArgumentBuilder]
}

func NewRoleArgument() *roleArgumentBuilder {
	b := &roleArgumentBuilder{}
	b.genericSlashCommandArgumentBuilder.upper = b
	return b
}

func (b *roleArgumentBuilder) create() SlashCommandArgument {
	panic("TODO")
}

// Mentionable

type mentionableArgumentBuilder struct {
	genericSlashCommandArgumentBuilder[*mentionableArgumentBuilder]
}

func NewMentionableArgument() *mentionableArgumentBuilder {
	b := &mentionableArgumentBuilder{}
	b.genericSlashCommandArgumentBuilder.upper = b
	return b
}

func (b *mentionableArgumentBuilder) create() SlashCommandArgument {
	panic("TODO")
}

// Attachment

type attachmentArgumentBuilder struct {
	genericSlashCommandArgumentBuilder[*attachmentArgumentBuilder]
}

func NewAttachmentArgument() *attachmentArgumentBuilder {
	b := &attachmentArgumentBuilder{}
	b.genericSlashCommandArgumentBuilder.upper = b
	return b
}

func (b *attachmentArgumentBuilder) create() SlashCommandArgument {
	panic("TODO")
}

// Channel

type channelArgumentBuilder struct {
	genericSlashCommandArgumentBuilder[*channelArgumentBuilder]
	channelTypes []discordgo.ChannelType
}

func NewChannelArgument() *channelArgumentBuilder {
	b := &channelArgumentBuilder{}
	b.genericSlashCommandArgumentBuilder.upper = b
	return b
}

func (b *channelArgumentBuilder) discordDefineForCreation() *discordgo.ApplicationCommandOption {
	o := b.genericSlashCommandArgumentBuilder.discordDefineForCreation()
	o.ChannelTypes = b.channelTypes
	return o
}

func (b *channelArgumentBuilder) create() SlashCommandArgument {
	panic("TODO")
}

func (b *channelArgumentBuilder) AllowChannel(channel discordgo.ChannelType) *channelArgumentBuilder {
	b.channelTypes = append(b.channelTypes, channel)
	return b
}

func (b *channelArgumentBuilder) AllowChannels(channel ...discordgo.ChannelType) *channelArgumentBuilder {
	b.channelTypes = append(b.channelTypes, channel...)
	return b
}
