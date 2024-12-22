package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type ArgumentInstance interface {
	parse(*discordgo.ApplicationCommandInteractionDataOption) (name string, value any, err error)
}

type ArgumentList struct {
	required map[string]ArgumentInstance
	optional map[string]ArgumentInstance
}

func NewArgumentList(required, optional map[string]ArgumentInstance) *ArgumentList {
	return &ArgumentList{
		required: required,
		optional: optional,
	}
}

func (al *ArgumentList) ParseOptions(options []*discordgo.ApplicationCommandInteractionDataOption) (CommandArguments, error) {

	return CommandArguments{}, nil
}

type CommandArguments struct {
	values map[string]any
}

func (args *CommandArguments) GetRequiredArgument(name string) any {
	arg, found := args.values[name]
	if !found {
		panic(fmt.Errorf("missing required argument %s, maybe you didn't mark it as required in the command definition?", name))
	}
	return arg
}

func (args *CommandArguments) GetArgument(name string) (value any, found bool) {
	v, f := args.values[name]
	return v, f
}

func (args *CommandArguments) GetArgumentOrDefault(name string, def any) any {
	if value, found := args.GetArgument(name); found {
		return value
	}
	return def
}

func (args *CommandArguments) GetRequiredBool(name string) bool {
	arg, found := args.values[name]
	if !found {
		panic(fmt.Errorf("missing required boolean argument %s, maybe you didn't mark it as required in the command definition?", name))
	}
	if value, ok := arg.(bool); ok {
		return value
	}
	panic(fmt.Errorf("required boolean argument %s is of type %t, maybe you didn't use the correct type in the command definition?", name, arg))
}

func (args *CommandArguments) GetBool(name string) (value bool, found bool) {
	arg, found := args.values[name]
	if !found {
		return false, false
	}
	if value, ok := arg.(bool); ok {
		return value, true
	}
	panic(fmt.Errorf("boolean argument %s is of type %t (%v), maybe you didn't use the correct type in the command definition?", name, arg, arg))
}

func (args *CommandArguments) GetBoolOrDefault(name string, def bool) bool {
	if value, found := args.GetBool(name); found {
		return value
	}
	return def
}

func (args *CommandArguments) GetRequiredInteger(name string) int64 {
	arg, found := args.values[name]
	if !found {
		panic(fmt.Errorf("missing required integer argument %s, maybe you didn't mark it as required in the command definition?", name))
	}
	if value, ok := arg.(int64); ok {
		return value
	}
	panic(fmt.Errorf("required integer argument %s is of type %t (%v), maybe you didn't use the correct type in the command definition?", name, arg, arg))
}

func (args *CommandArguments) GetInteger(name string) (value int64, found bool) {
	arg, found := args.values[name]
	if !found {
		return 0, false
	}
	if value, ok := arg.(int64); ok {
		return value, true
	}
	panic(fmt.Errorf("integer argument %s is of type %t (%v), maybe you didn't use the correct type in the command definition?", name, arg, arg))
}

func (args *CommandArguments) GetIntegerOrDefault(name string, def int64) int64 {
	if value, found := args.GetInteger(name); found {
		return value
	}
	return def
}
