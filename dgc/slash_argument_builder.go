package dgc

import (
	"github.com/bwmarrin/discordgo"
)

type SlashCommandArgumentBuilder interface {
	DiscordDefineForCreation() *discordgo.ApplicationCommandOption
	// If the returned requiredName is not nil, when the arguments of a given command are being parsed, a value associated
	// to the returned requiredName must be present, if not, the handler of the command won't be invoked, returning an error
	// ErrMissingRequiredArguments to the user who used the command
	Create() (requiredName *string, argument SlashCommandArgument)
}

type specificSlashCommandArgumentBuilder interface {
	createSpecific() SlashCommandArgument
}

type genericSlashCommandArgumentBuilder[B specificSlashCommandArgumentBuilder] struct {
	upper       B
	kind        discordgo.ApplicationCommandOptionType
	name        string
	description string
	required    bool
}

func (arg *genericSlashCommandArgumentBuilder[B]) DiscordDefineForCreation() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        arg.kind,
		Name:        arg.name,
		Description: arg.description,
		Required:    arg.required,
	}
}

func (arg *genericSlashCommandArgumentBuilder[B]) Create() (*string, SlashCommandArgument) {
	if arg.required {
		return &arg.name, arg.upper.createSpecific()
	}
	return nil, arg.upper.createSpecific()
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
	b.genericSlashCommandArgumentBuilder.kind = discordgo.ApplicationCommandOptionBoolean
	b.genericSlashCommandArgumentBuilder.upper = b
	return b
}

func (b *booleanArgumentBuilder) createSpecific() SlashCommandArgument {
	return &BooleanSlashCommandArgument{inlinedSlashCommandArgument[bool]{b.name}}
}

// String

type stringArgumentBuilder struct {
	genericSlashCommandArgumentBuilder[*stringArgumentBuilder]
	minLength *int
	maxLength int
}

func NewStringArgument() *stringArgumentBuilder {
	b := &stringArgumentBuilder{}
	b.genericSlashCommandArgumentBuilder.kind = discordgo.ApplicationCommandOptionString
	b.genericSlashCommandArgumentBuilder.upper = b
	return b
}

func (b *stringArgumentBuilder) DiscordDefineForCreation() *discordgo.ApplicationCommandOption {
	o := b.genericSlashCommandArgumentBuilder.DiscordDefineForCreation()
	if b.minLength != nil {
		o.MinLength = b.minLength
	}
	o.MaxLength = b.maxLength
	return o
}

func (b *stringArgumentBuilder) createSpecific() SlashCommandArgument {
	return &StringSlashCommandArgument{inlinedSlashCommandArgument[string]{b.name}}
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
	b.genericSlashCommandArgumentBuilder.kind = discordgo.ApplicationCommandOptionInteger
	b.genericSlashCommandArgumentBuilder.upper = b
	return b
}

func (b *integerArgumentBuilder) DiscordDefineForCreation() *discordgo.ApplicationCommandOption {
	o := b.genericSlashCommandArgumentBuilder.DiscordDefineForCreation()
	if b.minValue != nil {
		fv := float64(*b.minValue)
		o.MinValue = &fv
	}
	o.MaxValue = float64(b.maxValue)
	return o
}

func (b *integerArgumentBuilder) createSpecific() SlashCommandArgument {
	return &IntegerSlashCommandArgument{b.name}
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
	b.genericSlashCommandArgumentBuilder.kind = discordgo.ApplicationCommandOptionNumber
	b.genericSlashCommandArgumentBuilder.upper = b
	return b
}

func (b *numberArgumentBuilder) DiscordDefineForCreation() *discordgo.ApplicationCommandOption {
	o := b.genericSlashCommandArgumentBuilder.DiscordDefineForCreation()
	if b.minValue != nil {
		o.MinValue = b.minValue
	}
	o.MaxValue = b.maxValue
	return o
}

func (b *numberArgumentBuilder) createSpecific() SlashCommandArgument {
	return &NumberSlashCommandArgument{inlinedSlashCommandArgument[float64]{b.name}}
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
	b.genericSlashCommandArgumentBuilder.kind = discordgo.ApplicationCommandOptionUser
	b.genericSlashCommandArgumentBuilder.upper = b
	return b
}

func (b *userArgumentBuilder) createSpecific() SlashCommandArgument {
	return &UserSlashCommandArgument{extractingSlashCommandArgument[discordgo.User]{b.name}}
}

// Role

type roleArgumentBuilder struct {
	genericSlashCommandArgumentBuilder[*roleArgumentBuilder]
}

func NewRoleArgument() *roleArgumentBuilder {
	b := &roleArgumentBuilder{}
	b.genericSlashCommandArgumentBuilder.kind = discordgo.ApplicationCommandOptionRole
	b.genericSlashCommandArgumentBuilder.upper = b
	return b
}

func (b *roleArgumentBuilder) createSpecific() SlashCommandArgument {
	return &RoleSlashCommandArgument{extractingSlashCommandArgument[discordgo.Role]{b.name}}
}

// Mentionable

type mentionableArgumentBuilder struct {
	genericSlashCommandArgumentBuilder[*mentionableArgumentBuilder]
}

func NewMentionableArgument() *mentionableArgumentBuilder {
	b := &mentionableArgumentBuilder{}
	b.genericSlashCommandArgumentBuilder.kind = discordgo.ApplicationCommandOptionMentionable
	b.genericSlashCommandArgumentBuilder.upper = b
	return b
}

func (b *mentionableArgumentBuilder) createSpecific() SlashCommandArgument {
	return &MentionableSlashCommandArgument{b.name}
}

// Attachment

type attachmentArgumentBuilder struct {
	genericSlashCommandArgumentBuilder[*attachmentArgumentBuilder]
}

func NewAttachmentArgument() *attachmentArgumentBuilder {
	b := &attachmentArgumentBuilder{}
	b.genericSlashCommandArgumentBuilder.kind = discordgo.ApplicationCommandOptionAttachment
	b.genericSlashCommandArgumentBuilder.upper = b
	return b
}

func (b *attachmentArgumentBuilder) createSpecific() SlashCommandArgument {
	return &AttachmentSlashCommandArgument{extractingSlashCommandArgument[discordgo.MessageAttachment]{b.name}}
}

// Channel

type channelArgumentBuilder struct {
	genericSlashCommandArgumentBuilder[*channelArgumentBuilder]
	channelTypes []discordgo.ChannelType
}

func NewChannelArgument() *channelArgumentBuilder {
	b := &channelArgumentBuilder{}
	b.genericSlashCommandArgumentBuilder.kind = discordgo.ApplicationCommandOptionChannel
	b.genericSlashCommandArgumentBuilder.upper = b
	return b
}

func (b *channelArgumentBuilder) DiscordDefineForCreation() *discordgo.ApplicationCommandOption {
	o := b.genericSlashCommandArgumentBuilder.DiscordDefineForCreation()
	o.ChannelTypes = b.channelTypes
	return o
}

func (b *channelArgumentBuilder) createSpecific() SlashCommandArgument {
	return &ChannelSlashCommandArgument{extractingSlashCommandArgument[discordgo.Channel]{b.name}}
}

func (b *channelArgumentBuilder) AllowChannel(channel discordgo.ChannelType) *channelArgumentBuilder {
	b.channelTypes = append(b.channelTypes, channel)
	return b
}

func (b *channelArgumentBuilder) AllowChannels(channel ...discordgo.ChannelType) *channelArgumentBuilder {
	b.channelTypes = append(b.channelTypes, channel...)
	return b
}
