package dgc

import "github.com/bwmarrin/discordgo"

type genericSlashCommandAutocompleteArgumentBuilder[T any, B specificSlashCommandArgumentBuilder] struct {
	genericSlashCommandArgumentBuilder[B]
	handler SlashCommandAutocompleteArgumentHandler[T]
}

func (b *genericSlashCommandAutocompleteArgumentBuilder[T, B]) discordDefineForCreation() *discordgo.ApplicationCommandOption {
	d := b.genericSlashCommandArgumentBuilder.DiscordDefineForCreation()
	d.Autocomplete = true
	return d
}

func (b *genericSlashCommandAutocompleteArgumentBuilder[T, B]) Handler(handler SlashCommandAutocompleteArgumentHandler[T]) B {
	b.handler = handler
	return b.upper
}

// String

type stringSlashCommandAutocompleteArgumentBuilder struct {
	genericSlashCommandAutocompleteArgumentBuilder[string, *stringSlashCommandAutocompleteArgumentBuilder]
}

func (b *stringSlashCommandAutocompleteArgumentBuilder) createSpecific() SlashCommandArgument {
	return &StringSlashCommandAutocompleteArgumentHandler{
		StringSlashCommandArgument: StringSlashCommandArgument{
			inlinedSlashCommandArgument[string]{b.name},
		},
		handler: b.handler,
	}
}

func NewStringAutocompleteArgument() *stringSlashCommandAutocompleteArgumentBuilder {
	b := &stringSlashCommandAutocompleteArgumentBuilder{}
	b.genericSlashCommandAutocompleteArgumentBuilder.genericSlashCommandArgumentBuilder.kind = discordgo.ApplicationCommandOptionString
	b.genericSlashCommandAutocompleteArgumentBuilder.genericSlashCommandArgumentBuilder.upper = b
	return b
}

// Integer

type integerSlashCommandAutocompleteArgumentBuilder struct {
	genericSlashCommandAutocompleteArgumentBuilder[int64, *integerSlashCommandAutocompleteArgumentBuilder]
}

func (b *integerSlashCommandAutocompleteArgumentBuilder) createSpecific() SlashCommandArgument {
	return &IntegerSlashCommandAutocompleteArgumentHandler{
		IntegerSlashCommandArgument: IntegerSlashCommandArgument{
			name: b.name,
		},
		handler: b.handler,
	}
}

func NewIntegerAutocompleteArgument() *integerSlashCommandAutocompleteArgumentBuilder {
	b := &integerSlashCommandAutocompleteArgumentBuilder{}
	b.genericSlashCommandAutocompleteArgumentBuilder.genericSlashCommandArgumentBuilder.kind = discordgo.ApplicationCommandOptionInteger
	b.genericSlashCommandAutocompleteArgumentBuilder.genericSlashCommandArgumentBuilder.upper = b
	return b
}

// Number

type numberSlashCommandAutocompleteArgumentBuilder struct {
	genericSlashCommandAutocompleteArgumentBuilder[float64, *numberSlashCommandAutocompleteArgumentBuilder]
}

func (b *numberSlashCommandAutocompleteArgumentBuilder) createSpecific() SlashCommandArgument {
	return &NumberSlashCommandAutocompleteArgumentHandler{
		NumberSlashCommandArgument: NumberSlashCommandArgument{
			inlinedSlashCommandArgument[float64]{b.name},
		},
		handler: b.handler,
	}
}

func NewNumberAutocompleteArgument() *numberSlashCommandAutocompleteArgumentBuilder {
	b := &numberSlashCommandAutocompleteArgumentBuilder{}
	b.genericSlashCommandAutocompleteArgumentBuilder.genericSlashCommandArgumentBuilder.kind = discordgo.ApplicationCommandOptionNumber
	b.genericSlashCommandAutocompleteArgumentBuilder.genericSlashCommandArgumentBuilder.upper = b
	return b
}
