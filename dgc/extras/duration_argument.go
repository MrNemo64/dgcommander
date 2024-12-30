package extras

import (
	"time"

	"github.com/MrNemo64/dgcommander/dgc"
	"github.com/MrNemo64/dgcommander/dgc/util"
	"github.com/bwmarrin/discordgo"
)

type DurationSlashCommandArgumentBuilder struct {
	name        util.Localizable[*DurationSlashCommandArgumentBuilder]
	description util.Localizable[*DurationSlashCommandArgumentBuilder]
	required    bool
}

func (b *DurationSlashCommandArgumentBuilder) Name() *util.Localizable[*DurationSlashCommandArgumentBuilder] {
	return &b.name
}

func (b *DurationSlashCommandArgumentBuilder) Description() *util.Localizable[*DurationSlashCommandArgumentBuilder] {
	return &b.description
}

func (b *DurationSlashCommandArgumentBuilder) Required(required bool) *DurationSlashCommandArgumentBuilder {
	b.required = required
	return b
}

func (b *DurationSlashCommandArgumentBuilder) Create() (*string, dgc.SlashCommandArgument) {
	if b.required {
		return &b.name.Value, &DurationSlashCommandArgument{b.name.Value}
	}
	return nil, &DurationSlashCommandArgument{b.name.Value}
}

func (b *DurationSlashCommandArgumentBuilder) DiscordDefineForCreation() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionString,
		Name:                     b.name.Value,
		NameLocalizations:        b.name.Localizations,
		Description:              b.description.Value,
		DescriptionLocalizations: b.description.Localizations,
		Required:                 b.required,
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
