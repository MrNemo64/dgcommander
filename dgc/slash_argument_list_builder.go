package dgc

import (
	"github.com/bwmarrin/discordgo"
)

type slashCommandArgumentListBuilder[B specificCommandBuilder] struct {
	upper     B
	arguments []SlashCommandArgumentBuilder
}

func (b *slashCommandArgumentListBuilder[B]) discordDefineForCreation() []*discordgo.ApplicationCommandOption {
	requiredArgs := make([]*discordgo.ApplicationCommandOption, 0)
	optionalArgs := make([]*discordgo.ApplicationCommandOption, 0)
	for _, arg := range b.arguments {
		def := arg.DiscordDefineForCreation()
		if def.Required {
			requiredArgs = append(requiredArgs, def)
		} else {
			optionalArgs = append(optionalArgs, def)
		}
	}
	return append(requiredArgs, optionalArgs...)
}

func (b *slashCommandArgumentListBuilder[B]) create() slashCommandArgumentListDefinition {
	var required []string
	var args []SlashCommandArgument

	for _, arg := range b.arguments {
		name, v := arg.Create()
		if name != nil {
			required = append(required, *name)
		}
		args = append(args, v)
	}

	return slashCommandArgumentListDefinition{
		required:  required,
		arguments: args,
	}
}

func (b *slashCommandArgumentListBuilder[B]) AddArgument(arg SlashCommandArgumentBuilder) B {
	b.arguments = append(b.arguments, arg)
	return b.upper
}

func (b *slashCommandArgumentListBuilder[B]) AddArguments(args ...SlashCommandArgumentBuilder) B {
	b.arguments = append(b.arguments, args...)
	return b.upper
}
