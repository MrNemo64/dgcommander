package dgc

import "github.com/bwmarrin/discordgo"

type ArgumentChoice[T any] struct {
	Name  string
	Value T
}

type choicesArgumentBuilder[T any, B any] struct {
	upper   B
	choices []ArgumentChoice[T]
}

func (b *choicesArgumentBuilder[T, B]) WithChoice(choice ArgumentChoice[T]) B {
	b.choices = append(b.choices, choice)
	return b.upper
}

func (b *choicesArgumentBuilder[T, B]) WithChoices(choice ...ArgumentChoice[T]) B {
	b.choices = append(b.choices, choice...)
	return b.upper
}

func (b *choicesArgumentBuilder[T, B]) discordDefineForCreation() []*discordgo.ApplicationCommandOptionChoice {
	c := make([]*discordgo.ApplicationCommandOptionChoice, len(b.choices))
	for i, choice := range b.choices {
		c[i] = &discordgo.ApplicationCommandOptionChoice{
			Name:  choice.Name,
			Value: choice.Value,
		}
	}
	return c
}

type choicesStringArgumentBuilder struct {
	baseCommandArgumentBuilder[*choicesStringArgumentBuilder]
	choicesArgumentBuilder[string, *choicesStringArgumentBuilder]
}

func NewStringChoicesArgument() *choicesStringArgumentBuilder {
	b := &choicesStringArgumentBuilder{baseCommandArgumentBuilder: baseCommandArgumentBuilder[*choicesStringArgumentBuilder]{kind: discordgo.ApplicationCommandOptionString}}
	b.baseCommandArgumentBuilder.upper = b
	b.choicesArgumentBuilder.upper = b
	return b
}

func (b *choicesStringArgumentBuilder) discordDefineForCreation() *discordgo.ApplicationCommandOption {
	c := b.baseCommandArgumentBuilder.discordDefineForCreation()
	c.Choices = b.choicesArgumentBuilder.discordDefineForCreation()
	return c
}

type choicesIntegerArgumentBuilder struct {
	baseCommandArgumentBuilder[*choicesIntegerArgumentBuilder]
	choicesArgumentBuilder[int64, *choicesIntegerArgumentBuilder]
}

func NewIntegerChoicesArgument() *choicesIntegerArgumentBuilder {
	b := &choicesIntegerArgumentBuilder{baseCommandArgumentBuilder: baseCommandArgumentBuilder[*choicesIntegerArgumentBuilder]{kind: discordgo.ApplicationCommandOptionInteger}}
	b.baseCommandArgumentBuilder.upper = b
	b.choicesArgumentBuilder.upper = b
	return b
}

func (b *choicesIntegerArgumentBuilder) discordDefineForCreation() *discordgo.ApplicationCommandOption {
	c := b.baseCommandArgumentBuilder.discordDefineForCreation()
	c.Choices = b.choicesArgumentBuilder.discordDefineForCreation()
	return c
}

type choicesNumberArgumentBuilder struct {
	baseCommandArgumentBuilder[*choicesNumberArgumentBuilder]
	choicesArgumentBuilder[float64, *choicesNumberArgumentBuilder]
}

func NewNumberChoicesArgument() *choicesNumberArgumentBuilder {
	b := &choicesNumberArgumentBuilder{baseCommandArgumentBuilder: baseCommandArgumentBuilder[*choicesNumberArgumentBuilder]{kind: discordgo.ApplicationCommandOptionNumber}}
	b.baseCommandArgumentBuilder.upper = b
	b.choicesArgumentBuilder.upper = b
	return b
}

func (b *choicesNumberArgumentBuilder) discordDefineForCreation() *discordgo.ApplicationCommandOption {
	c := b.baseCommandArgumentBuilder.discordDefineForCreation()
	c.Choices = b.choicesArgumentBuilder.discordDefineForCreation()
	return c
}
