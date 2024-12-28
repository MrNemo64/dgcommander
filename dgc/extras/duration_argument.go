package extras

import (
	"time"

	"github.com/MrNemo64/dgcommander/dgc"
	"github.com/bwmarrin/discordgo"
)

type DurationSlashCommandArgumentBuilder struct {
	name        string
	description string
	required    bool
}

func (b *DurationSlashCommandArgumentBuilder) Name(name string) *DurationSlashCommandArgumentBuilder {
	b.name = name
	return b
}

func (b *DurationSlashCommandArgumentBuilder) Description(description string) *DurationSlashCommandArgumentBuilder {
	b.description = description
	return b
}

func (b *DurationSlashCommandArgumentBuilder) Required(required bool) *DurationSlashCommandArgumentBuilder {
	b.required = required
	return b
}

func (b *DurationSlashCommandArgumentBuilder) Create() (*string, dgc.SlashCommandArgument) {
	if b.required {
		return &b.name, &DurationSlashCommandArgument{b.name}
	}
	return nil, &DurationSlashCommandArgument{b.name}
}

func (b *DurationSlashCommandArgumentBuilder) DiscordDefineForCreation() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        b.name,
		Description: b.description,
		Required:    b.required,
	}
}

func NewDurationArgument() *DurationSlashCommandArgumentBuilder {
	return &DurationSlashCommandArgumentBuilder{}
}

type DurationSlashCommandArgument struct {
	name string
}

func (a *DurationSlashCommandArgument) Name() string { return a.name }

func (a *DurationSlashCommandArgument) Parse(info *dgc.ArgumentParsingInformation) (string, any, error) {
	op := info.FindOption(a.name)
	if op == nil {
		return "", nil, dgc.ErrArgumentHasNoValue.New(a.name)
	}
	value, ok := op.Value.(string)
	if !ok {
		return "", nil, dgc.ErrArgumentHasInvalidValue.New(a.name, op.Value, "duration")
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return "", nil, err
	}
	return a.name, duration, nil
}
