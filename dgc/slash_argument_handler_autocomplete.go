package dgc

type SlashCommandAutocompleteArgument interface {
	SlashCommandArgument
	Autocomplete() any
}

type SlashCommandAutocompleteArgumentHandler[T any] func() any

type genericSlashCommandAutocompleteArgumentHandler[A SlashCommandArgument, T any] struct {
	arg     A
	handler SlashCommandAutocompleteArgumentHandler[T]
}

// String

type StringSlashCommandAutocompleteArgumentHandler struct {
	StringSlashCommandArgument
	handler SlashCommandAutocompleteArgumentHandler[string]
}

func (arg *StringSlashCommandAutocompleteArgumentHandler) Autocomplete() any {
	return arg.handler()
}

// Integer

type IntegerSlashCommandAutocompleteArgumentHandler struct {
	IntegerSlashCommandArgument
	handler SlashCommandAutocompleteArgumentHandler[int64]
}

func (arg *IntegerSlashCommandAutocompleteArgumentHandler) Autocomplete() any {
	return arg.handler()
}

// Number

type NumberSlashCommandAutocompleteArgumentHandler struct {
	NumberSlashCommandArgument
	handler SlashCommandAutocompleteArgumentHandler[float64]
}

func (arg *NumberSlashCommandAutocompleteArgumentHandler) Autocomplete() any {
	return arg.handler()
}
