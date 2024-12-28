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

type stringSlashCommandArgumentBuilder struct {
	genericSlashCommandArgumentBuilder[*stringSlashCommandArgumentBuilder]
	minLength *int
	maxLength int
}

func NewStringArgument() *stringSlashCommandArgumentBuilder {
	b := &stringSlashCommandArgumentBuilder{}
	b.genericSlashCommandArgumentBuilder.kind = discordgo.ApplicationCommandOptionString
	b.genericSlashCommandArgumentBuilder.upper = b
	return b
}

func (b *stringSlashCommandArgumentBuilder) DiscordDefineForCreation() *discordgo.ApplicationCommandOption {
	o := b.genericSlashCommandArgumentBuilder.DiscordDefineForCreation()
	if b.minLength != nil {
		o.MinLength = b.minLength
	}
	o.MaxLength = b.maxLength
	return o
}

func (b *stringSlashCommandArgumentBuilder) createSpecific() SlashCommandArgument {
	return &StringSlashCommandArgument{inlinedSlashCommandArgument[string]{b.name}}
}

// 0-6000
func (b *stringSlashCommandArgumentBuilder) MinLength(min int) *stringSlashCommandArgumentBuilder {
	b.minLength = &min
	return b
}

// 1-6000
func (b *stringSlashCommandArgumentBuilder) MaxLength(max int) *stringSlashCommandArgumentBuilder {
	b.maxLength = max
	return b
}

// Integer

type integerSlashCommandArgumentBuilder struct {
	genericSlashCommandArgumentBuilder[*integerSlashCommandArgumentBuilder]
	minValue *int
	maxValue int
}

func NewIntegerArgument() *integerSlashCommandArgumentBuilder {
	b := &integerSlashCommandArgumentBuilder{}
	b.genericSlashCommandArgumentBuilder.kind = discordgo.ApplicationCommandOptionInteger
	b.genericSlashCommandArgumentBuilder.upper = b
	return b
}

func (b *integerSlashCommandArgumentBuilder) DiscordDefineForCreation() *discordgo.ApplicationCommandOption {
	o := b.genericSlashCommandArgumentBuilder.DiscordDefineForCreation()
	if b.minValue != nil {
		fv := float64(*b.minValue)
		o.MinValue = &fv
	}
	o.MaxValue = float64(b.maxValue)
	return o
}

func (b *integerSlashCommandArgumentBuilder) createSpecific() SlashCommandArgument {
	return &IntegerSlashCommandArgument{b.name}
}

func (b *integerSlashCommandArgumentBuilder) MinValue(min int) *integerSlashCommandArgumentBuilder {
	b.minValue = &min
	return b
}

func (b *integerSlashCommandArgumentBuilder) MaxValue(max int) *integerSlashCommandArgumentBuilder {
	b.maxValue = max
	return b
}

// Number

type numberSlashCommandArgumentBuilder struct {
	genericSlashCommandArgumentBuilder[*numberSlashCommandArgumentBuilder]
	minValue *float64
	maxValue float64
}

func NewNumberArgument() *numberSlashCommandArgumentBuilder {
	b := &numberSlashCommandArgumentBuilder{}
	b.genericSlashCommandArgumentBuilder.kind = discordgo.ApplicationCommandOptionNumber
	b.genericSlashCommandArgumentBuilder.upper = b
	return b
}

func (b *numberSlashCommandArgumentBuilder) DiscordDefineForCreation() *discordgo.ApplicationCommandOption {
	o := b.genericSlashCommandArgumentBuilder.DiscordDefineForCreation()
	if b.minValue != nil {
		o.MinValue = b.minValue
	}
	o.MaxValue = b.maxValue
	return o
}

func (b *numberSlashCommandArgumentBuilder) createSpecific() SlashCommandArgument {
	return &NumberSlashCommandArgument{inlinedSlashCommandArgument[float64]{b.name}}
}

func (b *numberSlashCommandArgumentBuilder) MinValue(min float64) *numberSlashCommandArgumentBuilder {
	b.minValue = &min
	return b
}

func (b *numberSlashCommandArgumentBuilder) MaxValue(max float64) *numberSlashCommandArgumentBuilder {
	b.maxValue = max
	return b
}

// User

type userSlashCommandArgumentBuilder struct {
	genericSlashCommandArgumentBuilder[*userSlashCommandArgumentBuilder]
}

func NewUserArgument() *userSlashCommandArgumentBuilder {
	b := &userSlashCommandArgumentBuilder{}
	b.genericSlashCommandArgumentBuilder.kind = discordgo.ApplicationCommandOptionUser
	b.genericSlashCommandArgumentBuilder.upper = b
	return b
}

func (b *userSlashCommandArgumentBuilder) createSpecific() SlashCommandArgument {
	a := &UserSlashCommandArgument{
		genericExtractingSlashCommandArgument: genericExtractingSlashCommandArgument[discordgo.User, *UserSlashCommandArgument]{
			name: b.name,
		},
	}
	a.genericExtractingSlashCommandArgument.specific = a
	return a
}

// Role

type roleSlashCommandArgumentBuilder struct {
	genericSlashCommandArgumentBuilder[*roleSlashCommandArgumentBuilder]
}

func NewRoleArgument() *roleSlashCommandArgumentBuilder {
	b := &roleSlashCommandArgumentBuilder{}
	b.genericSlashCommandArgumentBuilder.kind = discordgo.ApplicationCommandOptionRole
	b.genericSlashCommandArgumentBuilder.upper = b
	return b
}

func (b *roleSlashCommandArgumentBuilder) createSpecific() SlashCommandArgument {
	a := &RoleSlashCommandArgument{
		genericExtractingSlashCommandArgument: genericExtractingSlashCommandArgument[discordgo.Role, *RoleSlashCommandArgument]{
			name: b.name,
		},
	}
	a.genericExtractingSlashCommandArgument.specific = a
	return a
}

// Mentionable

type mentionableSlashCommandArgumentBuilder struct {
	genericSlashCommandArgumentBuilder[*mentionableSlashCommandArgumentBuilder]
}

func NewMentionableArgument() *mentionableSlashCommandArgumentBuilder {
	b := &mentionableSlashCommandArgumentBuilder{}
	b.genericSlashCommandArgumentBuilder.kind = discordgo.ApplicationCommandOptionMentionable
	b.genericSlashCommandArgumentBuilder.upper = b
	return b
}

func (b *mentionableSlashCommandArgumentBuilder) createSpecific() SlashCommandArgument {
	return &MentionableSlashCommandArgument{b.name}
}

// Attachment

type attachmentSlashCommandArgumentBuilder struct {
	genericSlashCommandArgumentBuilder[*attachmentSlashCommandArgumentBuilder]
}

func NewAttachmentArgument() *attachmentSlashCommandArgumentBuilder {
	b := &attachmentSlashCommandArgumentBuilder{}
	b.genericSlashCommandArgumentBuilder.kind = discordgo.ApplicationCommandOptionAttachment
	b.genericSlashCommandArgumentBuilder.upper = b
	return b
}

func (b *attachmentSlashCommandArgumentBuilder) createSpecific() SlashCommandArgument {
	a := &AttachmentSlashCommandArgument{
		genericExtractingSlashCommandArgument: genericExtractingSlashCommandArgument[discordgo.MessageAttachment, *AttachmentSlashCommandArgument]{
			name: b.name,
		},
	}
	a.genericExtractingSlashCommandArgument.specific = a
	return a
}

// Channel

type channelSlashCommandArgumentBuilder struct {
	genericSlashCommandArgumentBuilder[*channelSlashCommandArgumentBuilder]
	channelTypes []discordgo.ChannelType
}

func NewChannelArgument() *channelSlashCommandArgumentBuilder {
	b := &channelSlashCommandArgumentBuilder{}
	b.genericSlashCommandArgumentBuilder.kind = discordgo.ApplicationCommandOptionChannel
	b.genericSlashCommandArgumentBuilder.upper = b
	return b
}

func (b *channelSlashCommandArgumentBuilder) DiscordDefineForCreation() *discordgo.ApplicationCommandOption {
	o := b.genericSlashCommandArgumentBuilder.DiscordDefineForCreation()
	o.ChannelTypes = b.channelTypes
	return o
}

func (b *channelSlashCommandArgumentBuilder) createSpecific() SlashCommandArgument {
	a := &ChannelSlashCommandArgument{
		genericExtractingSlashCommandArgument: genericExtractingSlashCommandArgument[discordgo.Channel, *ChannelSlashCommandArgument]{
			name: b.name,
		},
	}
	a.genericExtractingSlashCommandArgument.specific = a
	return a
}

func (b *channelSlashCommandArgumentBuilder) AllowChannel(channel discordgo.ChannelType) *channelSlashCommandArgumentBuilder {
	b.channelTypes = append(b.channelTypes, channel)
	return b
}

func (b *channelSlashCommandArgumentBuilder) AllowChannels(channel ...discordgo.ChannelType) *channelSlashCommandArgumentBuilder {
	b.channelTypes = append(b.channelTypes, channel...)
	return b
}
