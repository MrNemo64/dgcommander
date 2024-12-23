package dgc

import (
	"github.com/bwmarrin/discordgo"
)

type slashCommandArgumentListBuilder[B specificCommandBuilder] struct {
	upper     B
	arguments []slashCommandArgumentBuilder
}

func (b *slashCommandArgumentListBuilder[B]) discordDefineForCreation() []*discordgo.ApplicationCommandOption {
	requiredArgs := make([]*discordgo.ApplicationCommandOption, 0)
	optionalArgs := make([]*discordgo.ApplicationCommandOption, 0)
	for _, arg := range b.arguments {
		def := arg.discordDefineForCreation()
		if def.Required {
			requiredArgs = append(requiredArgs, def)
		} else {
			optionalArgs = append(optionalArgs, def)
		}
	}
	return append(requiredArgs, optionalArgs...)
}

func (b *slashCommandArgumentListBuilder[B]) create() {
	panic("TODO")
}

func (b *slashCommandArgumentListBuilder[B]) AddArgument(arg slashCommandArgumentBuilder) B {
	b.arguments = append(b.arguments, arg)
	return b.upper
}

func (b *slashCommandArgumentListBuilder[B]) AddArguments(args ...slashCommandArgumentBuilder) B {
	b.arguments = append(b.arguments, args...)
	return b.upper
}
