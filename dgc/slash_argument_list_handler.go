package dgc

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

var (
	ErrMissingRequiredArguments = makeError("required arguments without values: %+v")                                                                        // arguments marked as required are not present after parsing all the arguments
	ErrMissingRequiredArgument  = makeError("missing required %s argument %s, maybe you didn't mark it as required in the command definition?")              // the handler requested as required a value that is not present
	ErrArgumentIsWrongType      = makeError("required %s argument %s is of type %T (%+v), maybe you didn't use the correct type in the command definition?") // the handler requested a value as one type but it is another tpye

	ErrArgumentHasInvalidValue = ArgumentHasInvalidValueError{DGCError: makeError("the argument %s has a value of type %T (%#v) but the expected type was %s")}
	ErrArgumentHasNoValue      = ArgumentHasNoValueError{DGCError: makeError("the %s has no value")}
)

type slashCommandArgumentListDefinition struct {
	arguments []SlashCommandArgument
	required  []string
}

type ArgumentParsingInformation struct {
	Options  []*discordgo.ApplicationCommandInteractionDataOption
	Resolved *discordgo.ApplicationCommandInteractionDataResolved
}

// Searches for an option in the slice of options with the same name, returning the first one found, nil if not found
func (api *ArgumentParsingInformation) FindOption(name string) *discordgo.ApplicationCommandInteractionDataOption {
	for _, o := range api.Options {
		if o.Name == name {
			return o
		}
	}
	return nil
}

func (d *slashCommandArgumentListDefinition) parse(resolved *discordgo.ApplicationCommandInteractionDataResolved, options []*discordgo.ApplicationCommandInteractionDataOption) (slashCommandArgumentList, error) {
	var allErrors []error
	list := slashCommandArgumentList{values: make(map[string]any)}
	parseInfo := ArgumentParsingInformation{Options: options, Resolved: resolved}
	for _, argument := range d.arguments {
		name, value, err := argument.Parse(&parseInfo)
		if err != nil {
			if !errors.Is(err, ErrArgumentHasNoValue) {
				allErrors = append(allErrors, err)
			}
			continue
		}
		list.values[name] = value
	}
	missingRequired := missingKeys(d.required, list.values)
	if len(missingRequired) > 0 {
		err := ErrMissingRequiredArguments.withArgs(missingRequired)
		if len(allErrors) > 0 {
			allErrors = append([]error{err}, allErrors...)
		} else {
			allErrors = []error{err}
		}
	}
	if len(allErrors) > 0 {
		return list, errors.Join(allErrors...)
	}
	return list, nil
}

func (d *slashCommandArgumentListDefinition) autocompleteArgumentFor(option *discordgo.ApplicationCommandInteractionDataOption) SlashCommandAutocompleteArgument {
	for _, arg := range d.arguments {
		auto, ok := arg.(SlashCommandAutocompleteArgument)
		if !ok {
			continue
		}
		if auto.IsForOption(option) {
			return auto
		}
	}
	return nil
}

type slashCommandArgumentList struct {
	values map[string]any
}
