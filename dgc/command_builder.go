package dgc

import (
	"slices"

	"github.com/MrNemo64/dgcommander/dgc/util"
	"github.com/bwmarrin/discordgo"
)

type CommandBuilder interface {
	discordDefineForCreation() *discordgo.ApplicationCommand
	create() command
}

type specificCommandBuilder interface {
}

type genericCommandBuilder[B specificCommandBuilder, M any] struct {
	upper        B
	name         util.Localizable[B]
	guildId      string
	nsfw         bool
	integrations []discordgo.ApplicationIntegrationType // https://discord.com/developers/docs/resources/application#application-object-application-integration-types
	contexts     []discordgo.InteractionContextType     // https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-interaction-context-types
	middlewares  []M
}

func (b *genericCommandBuilder[B, M]) discordDefineForCreation() *discordgo.ApplicationCommand {
	var integrations *[]discordgo.ApplicationIntegrationType
	if len(b.integrations) > 0 {
		integrations = &b.integrations
	}
	var contexts *[]discordgo.InteractionContextType
	if len(b.contexts) > 0 {
		contexts = &b.contexts
	}
	var nameLocalizations *map[discordgo.Locale]string
	if b.name.Localizations != nil {
		nameLocalizations = &b.name.Localizations
	}
	return &discordgo.ApplicationCommand{
		Name:              b.name.Value,
		NameLocalizations: nameLocalizations,
		NSFW:              &b.nsfw,
		IntegrationTypes:  integrations,
		Contexts:          contexts,
	}
}

func (b *genericCommandBuilder[B, M]) Name() *util.Localizable[B] {
	return &b.name
}

func (b *genericCommandBuilder[B, M]) ForGuild(guildId string) B {
	b.guildId = guildId
	return b.upper
}

func (b *genericCommandBuilder[B, M]) Nsfw(nsfw bool) B {
	b.nsfw = nsfw
	return b.upper
}

func (b *genericCommandBuilder[B, M]) AddMiddleware(middleware M) B {
	b.middlewares = append(b.middlewares, middleware)
	return b.upper
}

func (b *genericCommandBuilder[B, M]) GuildInstallable(installable bool) B {
	if installable {
		if !slices.Contains(b.integrations, discordgo.ApplicationIntegrationGuildInstall) {
			b.integrations = append(b.integrations, discordgo.ApplicationIntegrationGuildInstall)
		}
	} else {
		b.integrations = removeElement(b.integrations, discordgo.ApplicationIntegrationGuildInstall)
	}
	return b.upper
}

func (b *genericCommandBuilder[B, M]) UserInstallable(installable bool) B {
	if installable {
		if !slices.Contains(b.integrations, discordgo.ApplicationIntegrationUserInstall) {
			b.integrations = append(b.integrations, discordgo.ApplicationIntegrationUserInstall)
		}
	} else {
		b.integrations = removeElement(b.integrations, discordgo.ApplicationIntegrationUserInstall)
	}
	return b.upper
}

func (b *genericCommandBuilder[B, M]) AllowInGuilds(allowed bool) B {
	if allowed {
		if !slices.Contains(b.contexts, discordgo.InteractionContextGuild) {
			b.contexts = append(b.contexts, discordgo.InteractionContextGuild)
		}
	} else {
		b.contexts = removeElement(b.contexts, discordgo.InteractionContextGuild)
	}
	return b.upper
}

func (b *genericCommandBuilder[B, M]) AllowInBotDM(allowed bool) B {
	if allowed {
		if !slices.Contains(b.contexts, discordgo.InteractionContextBotDM) {
			b.contexts = append(b.contexts, discordgo.InteractionContextBotDM)
		}
	} else {
		b.contexts = removeElement(b.contexts, discordgo.InteractionContextBotDM)
	}
	return b.upper
}

func (b *genericCommandBuilder[B, M]) AllowInPrivateChannel(allowed bool) B {
	if allowed {
		if !slices.Contains(b.contexts, discordgo.InteractionContextPrivateChannel) {
			b.contexts = append(b.contexts, discordgo.InteractionContextPrivateChannel)
		}
	} else {
		b.contexts = removeElement(b.contexts, discordgo.InteractionContextPrivateChannel)
	}
	return b.upper
}

func (b *genericCommandBuilder[B, M]) AllowEverywhere(allowed bool) B {
	b.AllowInBotDM(allowed)
	b.AllowInGuilds(allowed)
	b.AllowInPrivateChannel(allowed)
	return b.upper
}
