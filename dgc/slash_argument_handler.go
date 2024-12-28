package dgc

import (
	"github.com/bwmarrin/discordgo"
)

// SlashCommandArgument represents an argument in a slash command that can parse its value
// from provided argument parsing information.
type SlashCommandArgument interface {
	// Parse extracts the value of this argument based on the provided parsing information.
	//
	// It returns:
	//  - valueName: the name associated with the parsed value.
	//  - value: the parsed value of this argument.
	//
	// Errors:
	//  - ErrArgumentHasNoValue: if the value of this argument is not present.
	//  - ErrArgumentHasInvalidValue: if the value of this argument has an invalid type.
	Parse(info *ArgumentParsingInformation) (valueName string, value any, err error)
	Name() string
}

type inlinedSlashCommandArgument[T any] struct{ name string }

func (a *inlinedSlashCommandArgument[T]) Parse(info *ArgumentParsingInformation) (string, any, error) {
	op := info.FindOption(a.name)
	if op == nil {
		return "", nil, ErrArgumentHasNoValue.New(a.name)
	}
	value, ok := op.Value.(T)
	if !ok {
		if !info.Autocompleting {
			return "", nil, ErrArgumentHasInvalidValue.New(a.name, op.Value, nameOfT[T]())
		}
		// if we are autocompleting, a string is technically considered a valid value
		// since the user can type whatever for us to autocomplete
		values, ok := op.Value.(string)
		if !ok {
			return "", nil, ErrArgumentHasInvalidValue.New(a.name, op.Value, nameOfT[T]())
		}
		return a.name, values, nil
	}
	return a.name, value, nil
}

func (a *inlinedSlashCommandArgument[T]) Name() string { return a.name }

type BooleanSlashCommandArgument struct {
	inlinedSlashCommandArgument[bool]
}

func (b *BooleanSlashCommandArgument) Parse(info *ArgumentParsingInformation) (string, any, error) {
	return b.inlinedSlashCommandArgument.Parse(info)
}

type NumberSlashCommandArgument struct {
	inlinedSlashCommandArgument[float64]
}

func (b *NumberSlashCommandArgument) Parse(info *ArgumentParsingInformation) (string, any, error) {
	return b.inlinedSlashCommandArgument.Parse(info)
}

type StringSlashCommandArgument struct {
	inlinedSlashCommandArgument[string]
}

func (b *StringSlashCommandArgument) Parse(info *ArgumentParsingInformation) (string, any, error) {
	return b.inlinedSlashCommandArgument.Parse(info)
}

type IntegerSlashCommandArgument struct{ name string }

func (a *IntegerSlashCommandArgument) Parse(info *ArgumentParsingInformation) (string, any, error) {
	op := info.FindOption(a.name)
	if op == nil {
		return "", nil, ErrArgumentHasNoValue.New(a.name)
	}
	value, ok := op.Value.(float64)
	if !ok {
		if !info.Autocompleting {
			return "", nil, ErrArgumentHasInvalidValue.New(a.name, op.Value, "int64")
		}
		// if we are autocompleting, a string is technically considered a valid value
		// since the user can type whatever for us to autocomplete
		values, ok := op.Value.(string)
		if !ok {
			return "", nil, ErrArgumentHasInvalidValue.New(a.name, op.Value, "int64")
		}
		return a.name, values, nil
	}
	return a.name, int64(value), nil
}

func (a *IntegerSlashCommandArgument) Name() string { return a.name }

type specificExtractingSlashCommandArgument[T any] interface {
	extract(info *ArgumentParsingInformation, id string) (*T, bool)
}

type genericExtractingSlashCommandArgument[T any, S specificExtractingSlashCommandArgument[T]] struct {
	specific S
	name     string
}

func (a *genericExtractingSlashCommandArgument[T, S]) Parse(info *ArgumentParsingInformation) (string, any, error) {
	op := info.FindOption(a.name)
	if op == nil {
		return "", nil, ErrArgumentHasNoValue.New(a.name)
	}
	valueId, ok := op.Value.(string)
	if !ok {
		return "", nil, ErrArgumentHasInvalidValue.New(a.name, op.Value, nameOfT[T]())
	}
	value, found := a.specific.extract(info, valueId)
	if !found {
		return "", nil, ErrArgumentHasNoValue.New(a.name)
	}
	return a.name, value, nil
}

func (a *genericExtractingSlashCommandArgument[T, S]) Name() string { return a.name }

type ChannelSlashCommandArgument struct {
	genericExtractingSlashCommandArgument[discordgo.Channel, *ChannelSlashCommandArgument]
}

func (b *ChannelSlashCommandArgument) Parse(info *ArgumentParsingInformation) (string, any, error) {
	return b.genericExtractingSlashCommandArgument.Parse(info)
}

func (ChannelSlashCommandArgument) extract(info *ArgumentParsingInformation, id string) (*discordgo.Channel, bool) {
	v, f := info.Resolved.Channels[id]
	return v, f
}

type AttachmentSlashCommandArgument struct {
	genericExtractingSlashCommandArgument[discordgo.MessageAttachment, *AttachmentSlashCommandArgument]
}

func (b *AttachmentSlashCommandArgument) Parse(info *ArgumentParsingInformation) (string, any, error) {
	return b.genericExtractingSlashCommandArgument.Parse(info)
}

func (AttachmentSlashCommandArgument) extract(info *ArgumentParsingInformation, id string) (*discordgo.MessageAttachment, bool) {
	v, f := info.Resolved.Attachments[id]
	return v, f
}

type UserSlashCommandArgument struct {
	genericExtractingSlashCommandArgument[discordgo.User, *UserSlashCommandArgument]
}

func (b *UserSlashCommandArgument) Parse(info *ArgumentParsingInformation) (string, any, error) {
	return b.genericExtractingSlashCommandArgument.Parse(info)
}

func (UserSlashCommandArgument) extract(info *ArgumentParsingInformation, id string) (*discordgo.User, bool) {
	v, f := info.Resolved.Users[id]
	return v, f
}

type RoleSlashCommandArgument struct {
	genericExtractingSlashCommandArgument[discordgo.Role, *RoleSlashCommandArgument]
}

func (b *RoleSlashCommandArgument) Parse(info *ArgumentParsingInformation) (string, any, error) {
	return b.genericExtractingSlashCommandArgument.Parse(info)
}

func (RoleSlashCommandArgument) extract(info *ArgumentParsingInformation, id string) (*discordgo.Role, bool) {
	v, f := info.Resolved.Roles[id]
	return v, f
}

type MemberSlashCommandArgument struct {
	genericExtractingSlashCommandArgument[discordgo.Member, *MemberSlashCommandArgument]
}

func (b *MemberSlashCommandArgument) Parse(info *ArgumentParsingInformation) (string, any, error) {
	return b.genericExtractingSlashCommandArgument.Parse(info)
}

func (MemberSlashCommandArgument) extract(info *ArgumentParsingInformation, id string) (*discordgo.Member, bool) {
	v, f := info.Resolved.Members[id]
	if !f {
		return nil, false
	}
	u, f := info.Resolved.Users[id]
	if !f {
		return nil, false
	}
	v.User = u
	return v, f
}

type MentionableSlashCommandArgument struct {
	name string
}

func (a *MentionableSlashCommandArgument) Parse(info *ArgumentParsingInformation) (string, any, error) {
	op := info.FindOption(a.name)
	if op == nil {
		return "", nil, ErrArgumentHasNoValue.New(a.name)
	}
	valueId, ok := op.Value.(string)
	if !ok {
		return "", nil, ErrArgumentHasInvalidValue.New(a.name, op.Value, "mentionable")
	}
	role, found := info.Resolved.Roles[valueId]
	if found {
		return a.name, role, nil
	} else {
		user, found := info.Resolved.Users[valueId]
		if !found {
			return "", nil, ErrArgumentHasNoValue.New(a.name)
		}
		member, found := info.Resolved.Members[valueId]
		if !found {
			return a.name, user, nil
		}
		member.User = user
		return a.name, member, nil
	}
}

func (a *MentionableSlashCommandArgument) Name() string { return a.name }
