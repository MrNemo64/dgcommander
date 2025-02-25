package dgc

import "github.com/bwmarrin/discordgo"

type genericSlashCommandAutocompleteArgumentBuilder[B specificSlashCommandArgumentBuilder] struct {
	genericSlashCommandArgumentBuilder[B]
	handler SlashCommandAutocompleteArgumentHandler
}

func (b *genericSlashCommandAutocompleteArgumentBuilder[B]) DiscordDefineForCreation() *discordgo.ApplicationCommandOption {
	d := b.genericSlashCommandArgumentBuilder.DiscordDefineForCreation()
	d.Autocomplete = true
	return d
}

func (b *genericSlashCommandAutocompleteArgumentBuilder[B]) Handler(handler SlashCommandAutocompleteArgumentHandler) B {
	b.handler = handler
	return b.upper
}

// String

type stringSlashCommandAutocompleteArgumentBuilder struct {
	genericSlashCommandAutocompleteArgumentBuilder[*stringSlashCommandAutocompleteArgumentBuilder]
}

func (b *stringSlashCommandAutocompleteArgumentBuilder) createSpecific() SlashCommandArgument {
	arg := &StringSlashCommandAutocompleteArgumentHandler{
		StringSlashCommandArgument: StringSlashCommandArgument{
			inlinedSlashCommandArgument[string]{b.name.Value},
		},
		genericSlashCommandAutocompleteArgumentHandler: genericSlashCommandAutocompleteArgumentHandler[*StringSlashCommandAutocompleteArgumentHandler]{
			handler: b.handler,
		},
	}
	arg.genericSlashCommandAutocompleteArgumentHandler.arg = arg
	return arg
}

func NewStringAutocompleteArgument() *stringSlashCommandAutocompleteArgumentBuilder {
	b := &stringSlashCommandAutocompleteArgumentBuilder{}
	b.genericSlashCommandAutocompleteArgumentBuilder.genericSlashCommandArgumentBuilder = newGenericSlashCommandArgumentBuilder(b, discordgo.ApplicationCommandOptionString)
	return b
}

// Integer

type integerSlashCommandAutocompleteArgumentBuilder struct {
	genericSlashCommandAutocompleteArgumentBuilder[*integerSlashCommandAutocompleteArgumentBuilder]
}

func (b *integerSlashCommandAutocompleteArgumentBuilder) createSpecific() SlashCommandArgument {
	arg := &IntegerSlashCommandAutocompleteArgumentHandler{
		IntegerSlashCommandArgument: IntegerSlashCommandArgument{
			name: b.name.Value,
		},
		genericSlashCommandAutocompleteArgumentHandler: genericSlashCommandAutocompleteArgumentHandler[*IntegerSlashCommandAutocompleteArgumentHandler]{
			handler: b.handler,
		},
	}
	arg.genericSlashCommandAutocompleteArgumentHandler.arg = arg
	return arg
}

func NewIntegerAutocompleteArgument() *integerSlashCommandAutocompleteArgumentBuilder {
	b := &integerSlashCommandAutocompleteArgumentBuilder{}
	b.genericSlashCommandAutocompleteArgumentBuilder.genericSlashCommandArgumentBuilder = newGenericSlashCommandArgumentBuilder(b, discordgo.ApplicationCommandOptionInteger)
	return b
}

// Number

type numberSlashCommandAutocompleteArgumentBuilder struct {
	genericSlashCommandAutocompleteArgumentBuilder[*numberSlashCommandAutocompleteArgumentBuilder]
}

func (b *numberSlashCommandAutocompleteArgumentBuilder) createSpecific() SlashCommandArgument {
	arg := &NumberSlashCommandAutocompleteArgumentHandler{
		NumberSlashCommandArgument: NumberSlashCommandArgument{
			inlinedSlashCommandArgument[float64]{b.name.Value},
		},
		genericSlashCommandAutocompleteArgumentHandler: genericSlashCommandAutocompleteArgumentHandler[*NumberSlashCommandAutocompleteArgumentHandler]{
			handler: b.handler,
		},
	}
	arg.genericSlashCommandAutocompleteArgumentHandler.arg = arg
	return arg
}

func NewNumberAutocompleteArgument() *numberSlashCommandAutocompleteArgumentBuilder {
	b := &numberSlashCommandAutocompleteArgumentBuilder{}
	b.genericSlashCommandAutocompleteArgumentBuilder.genericSlashCommandArgumentBuilder = newGenericSlashCommandArgumentBuilder(b, discordgo.ApplicationCommandOptionNumber)
	return b
}
