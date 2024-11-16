package dgc

import "github.com/bwmarrin/discordgo"

type AutocompleteHanlder func()

type autocompletedArgumentBuilder[T any, B any] struct {
	upper   B
	handler AutocompleteHanlder
}

func (b *autocompletedArgumentBuilder[T, B]) Handler(handler AutocompleteHanlder) B {
	b.handler = handler
	return b.upper
}

type autocompletedStringArgumentBuilder struct {
	baseCommandArgumentBuilder[*autocompletedStringArgumentBuilder]
	autocompletedArgumentBuilder[string, *autocompletedStringArgumentBuilder]
}

func NewStringAutocompletedArgument() *autocompletedStringArgumentBuilder {
	b := &autocompletedStringArgumentBuilder{baseCommandArgumentBuilder: baseCommandArgumentBuilder[*autocompletedStringArgumentBuilder]{kind: discordgo.ApplicationCommandOptionString}}
	b.baseCommandArgumentBuilder.upper = b
	b.autocompletedArgumentBuilder.upper = b
	return b
}

func (b *autocompletedStringArgumentBuilder) discordDefineForCreation() *discordgo.ApplicationCommandOption {
	c := b.baseCommandArgumentBuilder.discordDefineForCreation()
	c.Autocomplete = true
	return c
}

type autocompletedIntegerArgumentBuilder struct {
	baseCommandArgumentBuilder[*autocompletedIntegerArgumentBuilder]
	autocompletedArgumentBuilder[float64, *autocompletedIntegerArgumentBuilder]
}

func NewIntegerAutocompletedArgument() *autocompletedIntegerArgumentBuilder {
	b := &autocompletedIntegerArgumentBuilder{baseCommandArgumentBuilder: baseCommandArgumentBuilder[*autocompletedIntegerArgumentBuilder]{kind: discordgo.ApplicationCommandOptionInteger}}
	b.baseCommandArgumentBuilder.upper = b
	b.autocompletedArgumentBuilder.upper = b
	return b
}

func (b *autocompletedIntegerArgumentBuilder) discordDefineForCreation() *discordgo.ApplicationCommandOption {
	c := b.baseCommandArgumentBuilder.discordDefineForCreation()
	c.Autocomplete = true
	return c
}

type autocompletedNumberArgumentBuilder struct {
	baseCommandArgumentBuilder[*autocompletedNumberArgumentBuilder]
	autocompletedArgumentBuilder[float64, *autocompletedNumberArgumentBuilder]
}

func NewNumberAutocompletedArgument() *autocompletedNumberArgumentBuilder {
	b := &autocompletedNumberArgumentBuilder{baseCommandArgumentBuilder: baseCommandArgumentBuilder[*autocompletedNumberArgumentBuilder]{kind: discordgo.ApplicationCommandOptionNumber}}
	b.baseCommandArgumentBuilder.upper = b
	b.autocompletedArgumentBuilder.upper = b
	return b
}

func (b *autocompletedNumberArgumentBuilder) discordDefineForCreation() *discordgo.ApplicationCommandOption {
	c := b.baseCommandArgumentBuilder.discordDefineForCreation()
	c.Autocomplete = true
	return c
}
