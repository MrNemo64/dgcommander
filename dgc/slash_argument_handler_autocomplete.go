package dgc

import "github.com/bwmarrin/discordgo"

type SlashCommandAutocompleteArgument interface {
	SlashCommandArgument
	IsForOption(option *discordgo.ApplicationCommandInteractionDataOption) bool
	Autocomplete(*SlashAutocompleteContext) error
}

type SlashCommandAutocompleteArgumentHandler func(*SlashAutocompleteContext) error

type genericSlashCommandAutocompleteArgumentHandler[A SlashCommandArgument] struct {
	arg     A
	handler SlashCommandAutocompleteArgumentHandler
}

func (arg *genericSlashCommandAutocompleteArgumentHandler[A]) IsForOption(option *discordgo.ApplicationCommandInteractionDataOption) bool {
	return arg.arg.Name() == option.Name
}

func (arg *genericSlashCommandAutocompleteArgumentHandler[A]) Autocomplete(ctx *SlashAutocompleteContext) error {
	return arg.handler(ctx)
}

// String

type StringSlashCommandAutocompleteArgumentHandler struct {
	StringSlashCommandArgument
	genericSlashCommandAutocompleteArgumentHandler[*StringSlashCommandAutocompleteArgumentHandler]
}

// Integer

type IntegerSlashCommandAutocompleteArgumentHandler struct {
	IntegerSlashCommandArgument
	genericSlashCommandAutocompleteArgumentHandler[*IntegerSlashCommandAutocompleteArgumentHandler]
}

// Number

type NumberSlashCommandAutocompleteArgumentHandler struct {
	NumberSlashCommandArgument
	genericSlashCommandAutocompleteArgumentHandler[*NumberSlashCommandAutocompleteArgumentHandler]
}
