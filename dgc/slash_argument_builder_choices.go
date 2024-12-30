package dgc

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type ArgumentChoice[T any] struct {
	Name          string
	Value         T
	Localizations map[discordgo.Locale]string
}

type genericSlashCommandChoicesArgumentBuilder[T any, B specificSlashCommandArgumentBuilder] struct {
	upper   B
	choices []ArgumentChoice[T]
}

func (b *genericSlashCommandChoicesArgumentBuilder[T, B]) WithChoice(choice ArgumentChoice[T]) B {
	b.choices = append(b.choices, choice)
	return b.upper
}

func (b *genericSlashCommandChoicesArgumentBuilder[T, B]) WithChoices(choice ...ArgumentChoice[T]) B {
	b.choices = append(b.choices, choice...)
	return b.upper
}

func (b *genericSlashCommandChoicesArgumentBuilder[T, B]) AddChoice(name string, value T) B {
	return b.AddLocalizedChoice(name, value, nil)
}

func (b *genericSlashCommandChoicesArgumentBuilder[T, B]) AddLocalizedChoice(name string, value T, localization map[discordgo.Locale]string) B {
	b.choices = append(b.choices, ArgumentChoice[T]{name, value, nil})
	return b.upper
}

func (b *genericSlashCommandChoicesArgumentBuilder[T, B]) AddChoices(choices ...any) B {
	if len(choices)%2 != 0 {
		panic(fmt.Errorf("Called choicesArgumentBuilder.AddChoices but the last choice has no value: %v", choices))
	}
	for i := 0; i < len(choices); i += 2 {
		name, ok := choices[i].(string)
		if !ok {
			panic(fmt.Errorf("Called choicesArgumentBuilder.AddChoices but %v at index %d is not a string", choices[i], i))
		}
		value, ok := choices[i+1].(T)
		if !ok {
			panic(fmt.Errorf("Called choicesArgumentBuilder.AddChoices but %v at index %d is not a %s", choices[i], i, nameOfT[T]()))
		}
		b.choices = append(b.choices, ArgumentChoice[T]{name, value, nil})
	}
	return b.upper
}

func (b *genericSlashCommandChoicesArgumentBuilder[T, B]) discordDefineForCreation() []*discordgo.ApplicationCommandOptionChoice {
	c := make([]*discordgo.ApplicationCommandOptionChoice, len(b.choices))
	for i, choice := range b.choices {
		c[i] = &discordgo.ApplicationCommandOptionChoice{
			Name:              choice.Name,
			Value:             choice.Value,
			NameLocalizations: choice.Localizations,
		}
	}
	return c
}

// String

type stringSlashCommandChoicesArgumentBuilder struct {
	genericSlashCommandArgumentBuilder[*stringSlashCommandChoicesArgumentBuilder]
	genericSlashCommandChoicesArgumentBuilder[string, *stringSlashCommandChoicesArgumentBuilder]
}

func (b *stringSlashCommandChoicesArgumentBuilder) DiscordDefineForCreation() *discordgo.ApplicationCommandOption {
	o := b.genericSlashCommandArgumentBuilder.DiscordDefineForCreation()
	o.Choices = b.genericSlashCommandChoicesArgumentBuilder.discordDefineForCreation()
	return o
}

func (b *stringSlashCommandChoicesArgumentBuilder) createSpecific() SlashCommandArgument {
	return &StringSlashCommandArgument{inlinedSlashCommandArgument[string]{b.name.Value}}
}

func NewStringChoicesArgument() *stringSlashCommandChoicesArgumentBuilder {
	b := &stringSlashCommandChoicesArgumentBuilder{}
	b.genericSlashCommandArgumentBuilder = newGenericSlashCommandArgumentBuilder(b, discordgo.ApplicationCommandOptionString)
	b.genericSlashCommandChoicesArgumentBuilder.upper = b
	return b
}

// Integer

type integerSlashCommandChoicesArgumentBuilder struct {
	genericSlashCommandArgumentBuilder[*integerSlashCommandChoicesArgumentBuilder]
	genericSlashCommandChoicesArgumentBuilder[int64, *integerSlashCommandChoicesArgumentBuilder]
}

func (b *integerSlashCommandChoicesArgumentBuilder) DiscordDefineForCreation() *discordgo.ApplicationCommandOption {
	o := b.genericSlashCommandArgumentBuilder.DiscordDefineForCreation()
	o.Choices = b.genericSlashCommandChoicesArgumentBuilder.discordDefineForCreation()
	return o
}

func (b *integerSlashCommandChoicesArgumentBuilder) createSpecific() SlashCommandArgument {
	return &IntegerSlashCommandArgument{name: b.name.Value}
}

func NewIntegerChoicesArgument() *integerSlashCommandChoicesArgumentBuilder {
	b := &integerSlashCommandChoicesArgumentBuilder{}
	b.genericSlashCommandArgumentBuilder.kind = discordgo.ApplicationCommandOptionInteger
	b.genericSlashCommandArgumentBuilder.upper = b
	b.genericSlashCommandChoicesArgumentBuilder.upper = b
	return b
}

// Number

type numberSlashCommandChoicesArgumentBuilder struct {
	genericSlashCommandArgumentBuilder[*numberSlashCommandChoicesArgumentBuilder]
	genericSlashCommandChoicesArgumentBuilder[float64, *numberSlashCommandChoicesArgumentBuilder]
}

func (b *numberSlashCommandChoicesArgumentBuilder) DiscordDefineForCreation() *discordgo.ApplicationCommandOption {
	o := b.genericSlashCommandArgumentBuilder.DiscordDefineForCreation()
	o.Choices = b.genericSlashCommandChoicesArgumentBuilder.discordDefineForCreation()
	return o
}

func (b *numberSlashCommandChoicesArgumentBuilder) createSpecific() SlashCommandArgument {
	return &NumberSlashCommandArgument{inlinedSlashCommandArgument[float64]{b.name.Value}}
}

func NewNumberChoicesArgument() *numberSlashCommandChoicesArgumentBuilder {
	b := &numberSlashCommandChoicesArgumentBuilder{}
	b.genericSlashCommandArgumentBuilder = newGenericSlashCommandArgumentBuilder(b, discordgo.ApplicationCommandOptionNumber)
	b.genericSlashCommandChoicesArgumentBuilder.upper = b
	return b
}
