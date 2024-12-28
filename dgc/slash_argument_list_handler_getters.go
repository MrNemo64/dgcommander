package dgc

import "github.com/bwmarrin/discordgo"

// Generic

func GetRequiredArgument[T any](ctx *SlashExecutionContext, name string) T {
	arg, found := ctx.slashCommandArgumentList.values[name]
	if !found {
		panic(ErrMissingRequiredArgument.withArgs(nameOfT[T](), name))
	}
	if value, ok := arg.(T); ok {
		return value
	}
	panic(ErrArgumentIsWrongType.withArgs(nameOfT[T](), name, arg, arg))
}

func GetArgument[T any](ctx *SlashExecutionContext, name string) (value T, found bool) {
	arg, found := ctx.slashCommandArgumentList.values[name]
	if !found {
		var zero T
		return zero, false
	}
	if value, ok := arg.(T); ok {
		return value, true
	}
	panic(ErrArgumentIsWrongType.withArgs(nameOfT[T](), name, arg, arg))
}

func GetArgumentOr[T any](ctx *SlashExecutionContext, name string, def T) T {
	if value, found := GetArgument[T](ctx, name); found {
		return value
	}
	return def
}

// Boolean

func (args *slashCommandArgumentList) GetRequiredBool(name string) bool {
	arg, found := args.values[name]
	if !found {
		panic(ErrMissingRequiredArgument.withArgs("boolean", name))
	}
	if value, ok := arg.(bool); ok {
		return value
	}
	panic(ErrArgumentIsWrongType.withArgs("boolean", name, arg, arg))
}

func (args *slashCommandArgumentList) GetBool(name string) (value bool, found bool) {
	arg, found := args.values[name]
	if !found {
		return false, false
	}
	if value, ok := arg.(bool); ok {
		return value, true
	}
	panic(ErrArgumentIsWrongType.withArgs("boolean", name, arg, arg))
}

func (args *slashCommandArgumentList) GetBoolOr(name string, def bool) bool {
	if value, found := args.GetBool(name); found {
		return value
	}
	return def
}

// String

func (args *slashCommandArgumentList) GetRequiredString(name string) string {
	arg, found := args.values[name]
	if !found {
		panic(ErrMissingRequiredArgument.withArgs("string", name))
	}
	if value, ok := arg.(string); ok {
		return value
	}
	panic(ErrArgumentIsWrongType.withArgs("string", name, arg, arg))
}

func (args *slashCommandArgumentList) GetString(name string) (value string, found bool) {
	arg, found := args.values[name]
	if !found {
		return "", false
	}
	if value, ok := arg.(string); ok {
		return value, true
	}
	panic(ErrArgumentIsWrongType.withArgs("string", name, arg, arg))
}

func (args *slashCommandArgumentList) GetStringOr(name string, def string) string {
	if value, found := args.GetString(name); found {
		return value
	}
	return def
}

// Integer

func (args *slashCommandArgumentList) GetRequiredInteger(name string) int64 {
	arg, found := args.values[name]
	if !found {
		panic(ErrMissingRequiredArgument.withArgs("integer", name))
	}
	if value, ok := arg.(int64); ok {
		return value
	}
	panic(ErrArgumentIsWrongType.withArgs("integer", name, arg, arg))
}

func (args *slashCommandArgumentList) GetInteger(name string) (value int64, found bool) {
	arg, found := args.values[name]
	if !found {
		return 0, false
	}
	if value, ok := arg.(int64); ok {
		return value, true
	}
	panic(ErrArgumentIsWrongType.withArgs("integer", name, arg, arg))
}

func (args *slashCommandArgumentList) GetIntegerOr(name string, def int64) int64 {
	if value, found := args.GetInteger(name); found {
		return value
	}
	return def
}

// Number

func (args *slashCommandArgumentList) GetRequiredNumber(name string) float64 {
	arg, found := args.values[name]
	if !found {
		panic(ErrMissingRequiredArgument.withArgs("number", name))
	}
	if value, ok := arg.(float64); ok {
		return value
	}
	panic(ErrArgumentIsWrongType.withArgs("number", name, arg, arg))
}

func (args *slashCommandArgumentList) GetNumber(name string) (value float64, found bool) {
	arg, found := args.values[name]
	if !found {
		return 0, false
	}
	if value, ok := arg.(float64); ok {
		return value, true
	}
	panic(ErrArgumentIsWrongType.withArgs("number", name, arg, arg))
}

func (args *slashCommandArgumentList) GetNumberOr(name string, def float64) float64 {
	if value, found := args.GetNumber(name); found {
		return value
	}
	return def
}

// User

func (args *slashCommandArgumentList) GetRequiredUser(name string) *discordgo.User {
	arg, found := args.values[name]
	if !found {
		panic(ErrMissingRequiredArgument.withArgs("user", name))
	}
	if value, ok := arg.(*discordgo.User); ok {
		return value
	}
	panic(ErrArgumentIsWrongType.withArgs("user", name, arg, arg))
}

func (args *slashCommandArgumentList) GetUser(name string) (value *discordgo.User, found bool) {
	arg, found := args.values[name]
	if !found {
		return nil, false
	}
	if value, ok := arg.(*discordgo.User); ok {
		return value, true
	}
	panic(ErrArgumentIsWrongType.withArgs("user", name, arg, arg))
}

func (args *slashCommandArgumentList) GetUserOr(name string, def *discordgo.User) *discordgo.User {
	if value, found := args.GetUser(name); found {
		return value
	}
	return def
}

// Member

func (args *slashCommandArgumentList) GetRequiredMember(name string) *discordgo.Member {
	arg, found := args.values[name]
	if !found {
		panic(ErrMissingRequiredArgument.withArgs("member", name))
	}
	if value, ok := arg.(*discordgo.Member); ok {
		return value
	}
	panic(ErrArgumentIsWrongType.withArgs("member", name, arg, arg))
}

func (args *slashCommandArgumentList) GetMember(name string) (value *discordgo.Member, found bool) {
	arg, found := args.values[name]
	if !found {
		return nil, false
	}
	if value, ok := arg.(*discordgo.Member); ok {
		return value, true
	}
	panic(ErrArgumentIsWrongType.withArgs("member", name, arg, arg))
}

func (args *slashCommandArgumentList) GetMemberOr(name string, def *discordgo.Member) *discordgo.Member {
	if value, found := args.GetMember(name); found {
		return value
	}
	return def
}

// Role

func (args *slashCommandArgumentList) GetRequiredRole(name string) *discordgo.Role {
	arg, found := args.values[name]
	if !found {
		panic(ErrMissingRequiredArgument.withArgs("role", name))
	}
	if value, ok := arg.(*discordgo.Role); ok {
		return value
	}
	panic(ErrArgumentIsWrongType.withArgs("role", name, arg, arg))
}

func (args *slashCommandArgumentList) GetRole(name string) (value *discordgo.Role, found bool) {
	arg, found := args.values[name]
	if !found {
		return nil, false
	}
	if value, ok := arg.(*discordgo.Role); ok {
		return value, true
	}
	panic(ErrArgumentIsWrongType.withArgs("role", name, arg, arg))
}

func (args *slashCommandArgumentList) GetRoleOr(name string, def *discordgo.Role) *discordgo.Role {
	if value, found := args.GetRole(name); found {
		return value
	}
	return def
}

// Mentionable
// TODO

// Attachment

func (args *slashCommandArgumentList) GetRequiredAttachment(name string) *discordgo.MessageAttachment {
	arg, found := args.values[name]
	if !found {
		panic(ErrMissingRequiredArgument.withArgs("attachment", name))
	}
	if value, ok := arg.(*discordgo.MessageAttachment); ok {
		return value
	}
	panic(ErrArgumentIsWrongType.withArgs("attachment", name, arg, arg))
}

func (args *slashCommandArgumentList) GetAttachment(name string) (value *discordgo.MessageAttachment, found bool) {
	arg, found := args.values[name]
	if !found {
		return nil, false
	}
	if value, ok := arg.(*discordgo.MessageAttachment); ok {
		return value, true
	}
	panic(ErrArgumentIsWrongType.withArgs("attachment", name, arg, arg))
}

func (args *slashCommandArgumentList) GetAttachmentOr(name string, def *discordgo.MessageAttachment) *discordgo.MessageAttachment {
	if value, found := args.GetAttachment(name); found {
		return value
	}
	return def
}

// Channel

func (args *slashCommandArgumentList) GetRequiredChannel(name string) *discordgo.Channel {
	arg, found := args.values[name]
	if !found {
		panic(ErrMissingRequiredArgument.withArgs("channel", name))
	}
	if value, ok := arg.(*discordgo.Channel); ok {
		return value
	}
	panic(ErrArgumentIsWrongType.withArgs("channel", name, arg, arg))
}

func (args *slashCommandArgumentList) GetChannel(name string) (value *discordgo.Channel, found bool) {
	arg, found := args.values[name]
	if !found {
		return nil, false
	}
	if value, ok := arg.(*discordgo.Channel); ok {
		return value, true
	}
	panic(ErrArgumentIsWrongType.withArgs("channel", name, arg, arg))
}

func (args *slashCommandArgumentList) GetChannelOr(name string, def *discordgo.Channel) *discordgo.Channel {
	if value, found := args.GetChannel(name); found {
		return value
	}
	return def
}
