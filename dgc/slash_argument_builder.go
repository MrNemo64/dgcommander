package dgc

import (
	"github.com/MrNemo64/dgcommander/dgc/util"
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
	name        util.Localizable[B]
	description util.Localizable[B]
	required    bool
}

func (arg *genericSlashCommandArgumentBuilder[B]) DiscordDefineForCreation() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:                     arg.kind,
		Name:                     arg.name.Value,
		NameLocalizations:        arg.name.Localizations,
		Description:              arg.description.Value,
		DescriptionLocalizations: arg.description.Localizations,
		Required:                 arg.required,
	}
}

func newGenericSlashCommandArgumentBuilder[B specificSlashCommandArgumentBuilder](upper B, kind discordgo.ApplicationCommandOptionType) genericSlashCommandArgumentBuilder[B] {
	b := genericSlashCommandArgumentBuilder[B]{}
	b.kind = kind
	b.upper = upper
	b.name.Upper = upper
	b.description.Upper = upper
	return b
}

func (arg *genericSlashCommandArgumentBuilder[B]) Create() (*string, SlashCommandArgument) {
	if arg.required {
		return &arg.name.Value, arg.upper.createSpecific()
	}
	return nil, arg.upper.createSpecific()
}

func (b *genericSlashCommandArgumentBuilder[B]) Name() *util.Localizable[B] {
	return &b.name
}

func (b *genericSlashCommandArgumentBuilder[B]) Description() *util.Localizable[B] {
	return &b.description
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
	b.genericSlashCommandArgumentBuilder = newGenericSlashCommandArgumentBuilder(b, discordgo.ApplicationCommandOptionBoolean)
	return b
}

func (b *booleanArgumentBuilder) createSpecific() SlashCommandArgument {
	return &BooleanSlashCommandArgument{inlinedSlashCommandArgument[bool]{b.name.Value}}
}

// String

type stringSlashCommandArgumentBuilder struct {
	genericSlashCommandArgumentBuilder[*stringSlashCommandArgumentBuilder]
	minLength *int
	maxLength int
}

func NewStringArgument() *stringSlashCommandArgumentBuilder {
	b := &stringSlashCommandArgumentBuilder{}
	b.genericSlashCommandArgumentBuilder = newGenericSlashCommandArgumentBuilder(b, discordgo.ApplicationCommandOptionString)
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
	return &StringSlashCommandArgument{inlinedSlashCommandArgument[string]{b.name.Value}}
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
	b.genericSlashCommandArgumentBuilder = newGenericSlashCommandArgumentBuilder(b, discordgo.ApplicationCommandOptionInteger)
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
	return &IntegerSlashCommandArgument{b.name.Value}
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
	b.genericSlashCommandArgumentBuilder = newGenericSlashCommandArgumentBuilder(b, discordgo.ApplicationCommandOptionNumber)
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
	return &NumberSlashCommandArgument{inlinedSlashCommandArgument[float64]{b.name.Value}}
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
	b.genericSlashCommandArgumentBuilder = newGenericSlashCommandArgumentBuilder(b, discordgo.ApplicationCommandOptionUser)
	return b
}

func (b *userSlashCommandArgumentBuilder) createSpecific() SlashCommandArgument {
	a := &UserSlashCommandArgument{
		genericExtractingSlashCommandArgument: genericExtractingSlashCommandArgument[discordgo.User, *UserSlashCommandArgument]{
			name: b.name.Value,
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
	b.genericSlashCommandArgumentBuilder = newGenericSlashCommandArgumentBuilder(b, discordgo.ApplicationCommandOptionRole)
	return b
}

func (b *roleSlashCommandArgumentBuilder) createSpecific() SlashCommandArgument {
	a := &RoleSlashCommandArgument{
		genericExtractingSlashCommandArgument: genericExtractingSlashCommandArgument[discordgo.Role, *RoleSlashCommandArgument]{
			name: b.name.Value,
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
	b.genericSlashCommandArgumentBuilder = newGenericSlashCommandArgumentBuilder(b, discordgo.ApplicationCommandOptionMentionable)
	return b
}

func (b *mentionableSlashCommandArgumentBuilder) createSpecific() SlashCommandArgument {
	return &MentionableSlashCommandArgument{b.name.Value}
}

// Attachment

type attachmentSlashCommandArgumentBuilder struct {
	genericSlashCommandArgumentBuilder[*attachmentSlashCommandArgumentBuilder]
}

func NewAttachmentArgument() *attachmentSlashCommandArgumentBuilder {
	b := &attachmentSlashCommandArgumentBuilder{}
	b.genericSlashCommandArgumentBuilder = newGenericSlashCommandArgumentBuilder(b, discordgo.ApplicationCommandOptionAttachment)
	return b
}

func (b *attachmentSlashCommandArgumentBuilder) createSpecific() SlashCommandArgument {
	a := &AttachmentSlashCommandArgument{
		genericExtractingSlashCommandArgument: genericExtractingSlashCommandArgument[discordgo.MessageAttachment, *AttachmentSlashCommandArgument]{
			name: b.name.Value,
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
	b.genericSlashCommandArgumentBuilder = newGenericSlashCommandArgumentBuilder(b, discordgo.ApplicationCommandOptionChannel)
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
			name: b.name.Value,
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
