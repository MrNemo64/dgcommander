package dgc

import "github.com/bwmarrin/discordgo"

type SlashCommandAutocompleteArgument interface {
	SlashCommandArgument
	IsForOption(option *discordgo.ApplicationCommandInteractionDataOption) bool
	Autocomplete(*discordgo.User, *SlashAutocompleteContext) error
}

type SlashCommandAutocompleteArgumentHandler func(*discordgo.User, *SlashAutocompleteContext) error

type genericSlashCommandAutocompleteArgumentHandler[A SlashCommandArgument] struct {
	arg     A
	handler SlashCommandAutocompleteArgumentHandler
}

func (arg *genericSlashCommandAutocompleteArgumentHandler[A]) IsForOption(option *discordgo.ApplicationCommandInteractionDataOption) bool {
	return arg.arg.Name() == option.Name
}

func (arg *genericSlashCommandAutocompleteArgumentHandler[A]) Autocomplete(sender *discordgo.User, ctx *SlashAutocompleteContext) error {
	return arg.handler(sender, ctx)
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
