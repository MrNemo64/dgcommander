package dgc

import (
	"github.com/MrNemo64/dgcommander/dgc/handlers"
	"github.com/bwmarrin/discordgo"
)

type commandArgument interface {
	isCommandArgument() bool
	discordDefineForCreation() *discordgo.ApplicationCommandOption
	create() (bool, string, handlers.ArgumentInstance)
}

type argumentListBuilder[B any] struct {
	upper     B
	arguments []commandArgument
}

func (h *argumentListBuilder[B]) AddArgument(arg commandArgument) B {
	h.arguments = append(h.arguments, arg)
	return h.upper
}

func (h *argumentListBuilder[B]) AddArguments(args ...commandArgument) B {
	h.arguments = append(h.arguments, args...)
	return h.upper
}

func (h *argumentListBuilder[B]) discordDefineForCreation() []*discordgo.ApplicationCommandOption {
	requiredArgs := make([]*discordgo.ApplicationCommandOption, 0)
	optionalArgs := make([]*discordgo.ApplicationCommandOption, 0)
	for _, arg := range h.arguments {
		def := arg.discordDefineForCreation()
		if def.Required {
			requiredArgs = append(requiredArgs, def)
		} else {
			optionalArgs = append(optionalArgs, def)
		}
	}
	return append(requiredArgs, optionalArgs...)
}

func (h *argumentListBuilder[B]) create() handlers.ArgumentList {
	required := make(map[string]handlers.ArgumentInstance)
	optional := make(map[string]handlers.ArgumentInstance)
	for _, argument := range h.arguments {
		r, name, arg := argument.create()
		if r {
			required[name] = arg
		} else {
			optional[name] = arg
		}
	}
	return *handlers.NewArgumentList(required, optional)
}

type specificCommandArgumentBuilder interface {
	create() handlers.ArgumentInstance
}

type baseCommandArgumentBuilder[B specificCommandArgumentBuilder] struct {
	upper       B
	kind        discordgo.ApplicationCommandOptionType
	name        string
	description string
	required    bool
}

func (*baseCommandArgumentBuilder[T]) isCommandArgument() bool { return true }

func (arg *baseCommandArgumentBuilder[T]) discordDefineForCreation() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        arg.kind,
		Name:        arg.name,
		Description: arg.description,
		Required:    arg.required,
	}
}

func (arg *baseCommandArgumentBuilder[T]) create() (bool, string, handlers.ArgumentInstance) {
	return arg.required, arg.name, arg.upper.create()
}

func (b *baseCommandArgumentBuilder[T]) Name(name string) T {
	b.name = name
	return b.upper
}

func (b *baseCommandArgumentBuilder[T]) Description(description string) T {
	b.description = description
	return b.upper
}

func (b *baseCommandArgumentBuilder[T]) Required(required bool) T {
	b.required = required
	return b.upper
}

type booleanArgumentBuilder struct {
	baseCommandArgumentBuilder[*booleanArgumentBuilder]
}

func NewBooleanArgument() *booleanArgumentBuilder {
	b := &booleanArgumentBuilder{baseCommandArgumentBuilder: baseCommandArgumentBuilder[*booleanArgumentBuilder]{kind: discordgo.ApplicationCommandOptionBoolean}}
	b.baseCommandArgumentBuilder.upper = b
	return b
}

func (b *booleanArgumentBuilder) create() handlers.ArgumentInstance {

}

type userArgumentBuilder struct {
	baseCommandArgumentBuilder[*userArgumentBuilder]
}

func NewUserArgument() *userArgumentBuilder {
	b := &userArgumentBuilder{baseCommandArgumentBuilder: baseCommandArgumentBuilder[*userArgumentBuilder]{kind: discordgo.ApplicationCommandOptionUser}}
	b.baseCommandArgumentBuilder.upper = b
	return b
}

type roleArgumentBuilder struct {
	baseCommandArgumentBuilder[*roleArgumentBuilder]
}

func NewRoleArgument() *roleArgumentBuilder {
	b := &roleArgumentBuilder{baseCommandArgumentBuilder: baseCommandArgumentBuilder[*roleArgumentBuilder]{kind: discordgo.ApplicationCommandOptionRole}}
	b.baseCommandArgumentBuilder.upper = b
	return b
}

type mentionableArgumentBuilder struct {
	baseCommandArgumentBuilder[*mentionableArgumentBuilder]
}

func NewMentionableArgument() *mentionableArgumentBuilder {
	b := &mentionableArgumentBuilder{baseCommandArgumentBuilder: baseCommandArgumentBuilder[*mentionableArgumentBuilder]{kind: discordgo.ApplicationCommandOptionMentionable}}
	b.baseCommandArgumentBuilder.upper = b
	return b
}

type attachmentArgumentBuilder struct {
	baseCommandArgumentBuilder[*attachmentArgumentBuilder]
}

func NewAttachmentArgument() *attachmentArgumentBuilder {
	b := &attachmentArgumentBuilder{baseCommandArgumentBuilder: baseCommandArgumentBuilder[*attachmentArgumentBuilder]{kind: discordgo.ApplicationCommandOptionAttachment}}
	b.baseCommandArgumentBuilder.upper = b
	return b
}

type channelArgumentBuilder struct {
	baseCommandArgumentBuilder[*channelArgumentBuilder]
	channelTypes []discordgo.ChannelType
}

func NewChannelArgument() *channelArgumentBuilder {
	b := &channelArgumentBuilder{baseCommandArgumentBuilder: baseCommandArgumentBuilder[*channelArgumentBuilder]{kind: discordgo.ApplicationCommandOptionChannel}}
	b.baseCommandArgumentBuilder.upper = b
	return b
}

func (b *channelArgumentBuilder) AllowChannel(channel discordgo.ChannelType) *channelArgumentBuilder {
	b.channelTypes = append(b.channelTypes, channel)
	return b
}

type stringArgumentBuilder struct {
	baseCommandArgumentBuilder[*stringArgumentBuilder]
	minLength *int
	maxLength *int
}

func NewStringArgument() *stringArgumentBuilder {
	b := &stringArgumentBuilder{baseCommandArgumentBuilder: baseCommandArgumentBuilder[*stringArgumentBuilder]{kind: discordgo.ApplicationCommandOptionString}}
	b.baseCommandArgumentBuilder.upper = b
	return b
}

func (b *stringArgumentBuilder) MinLength(min int) *stringArgumentBuilder {
	b.minLength = &min
	return b
}

func (b *stringArgumentBuilder) MaxLength(max int) *stringArgumentBuilder {
	b.maxLength = &max
	return b
}

type integerArgumentBuilder struct {
	baseCommandArgumentBuilder[*integerArgumentBuilder]
	min *int64
	max *int64
}

func NewIntegerArgument() *integerArgumentBuilder {
	b := &integerArgumentBuilder{baseCommandArgumentBuilder: baseCommandArgumentBuilder[*integerArgumentBuilder]{kind: discordgo.ApplicationCommandOptionInteger}}
	b.baseCommandArgumentBuilder.upper = b
	return b
}

func (b *integerArgumentBuilder) Min(min int64) *integerArgumentBuilder {
	b.min = &min
	return b
}

func (b *integerArgumentBuilder) Max(max int64) *integerArgumentBuilder {
	b.max = &max
	return b
}

type numberArgumentBuilder struct {
	baseCommandArgumentBuilder[*numberArgumentBuilder]
	min *float64
	max *float64
}

func NewNumberArgument() *numberArgumentBuilder {
	b := &numberArgumentBuilder{baseCommandArgumentBuilder: baseCommandArgumentBuilder[*numberArgumentBuilder]{kind: discordgo.ApplicationCommandOptionNumber}}
	b.baseCommandArgumentBuilder.upper = b
	return b
}

func (b *numberArgumentBuilder) Min(min float64) *numberArgumentBuilder {
	b.min = &min
	return b
}

func (b *numberArgumentBuilder) Max(max float64) *numberArgumentBuilder {
	b.max = &max
	return b
}
