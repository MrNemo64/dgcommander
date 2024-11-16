package dgc

import (
	"slices"

	"github.com/bwmarrin/discordgo"
)

type CommandBuilder interface {
	discordDefineForCreation() *discordgo.ApplicationCommand
	guild() string
	create() command
}

type commandBuilder[H any, B any] struct {
	upper        B
	handler      H
	name         string
	guildId      string
	nsfw         bool
	integrations []discordgo.ApplicationIntegrationType // https://discord.com/developers/docs/resources/application#application-object-application-integration-types
	contexts     []discordgo.InteractionContextType     // https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-interaction-context-types
}

func (b *commandBuilder[H, B]) guild() string { return b.guildId }

func (b *commandBuilder[H, B]) discordDefineForCreation() *discordgo.ApplicationCommand {
	var integrations *[]discordgo.ApplicationIntegrationType
	if len(b.integrations) > 0 {
		integrations = &b.integrations
	}
	var contexts *[]discordgo.InteractionContextType
	if len(b.contexts) > 0 {
		contexts = &b.contexts
	}
	return &discordgo.ApplicationCommand{
		Name:             b.name,
		NSFW:             &b.nsfw,
		IntegrationTypes: integrations,
		Contexts:         contexts,
	}
}

func (b *commandBuilder[H, B]) Handler(handler H) B {
	b.handler = handler
	return b.upper
}

func (b *commandBuilder[H, B]) Name(name string) B {
	b.name = name
	return b.upper
}

func (b *commandBuilder[H, B]) ForGuild(guildId string) B {
	b.guildId = guildId
	return b.upper
}

func (b *commandBuilder[H, B]) Nsfw(nsfw bool) B {
	b.nsfw = nsfw
	return b.upper
}

func (b *commandBuilder[H, B]) GuildInstallable(installable bool) B {
	if installable {
		if !slices.Contains(b.integrations, discordgo.ApplicationIntegrationGuildInstall) {
			b.integrations = append(b.integrations, discordgo.ApplicationIntegrationGuildInstall)
		}
	} else {
		b.integrations = removeElement(b.integrations, discordgo.ApplicationIntegrationGuildInstall)
	}
	return b.upper
}

func (b *commandBuilder[H, B]) UserInstallable(installable bool) B {
	if installable {
		if !slices.Contains(b.integrations, discordgo.ApplicationIntegrationUserInstall) {
			b.integrations = append(b.integrations, discordgo.ApplicationIntegrationUserInstall)
		}
	} else {
		b.integrations = removeElement(b.integrations, discordgo.ApplicationIntegrationUserInstall)
	}
	return b.upper
}

func (b *commandBuilder[H, B]) AllowInGuilds(allowed bool) B {
	if allowed {
		if !slices.Contains(b.contexts, discordgo.InteractionContextGuild) {
			b.contexts = append(b.contexts, discordgo.InteractionContextGuild)
		}
	} else {
		b.contexts = removeElement(b.contexts, discordgo.InteractionContextGuild)
	}
	return b.upper
}

func (b *commandBuilder[H, B]) AllowInBotDM(allowed bool) B {
	if allowed {
		if !slices.Contains(b.contexts, discordgo.InteractionContextBotDM) {
			b.contexts = append(b.contexts, discordgo.InteractionContextBotDM)
		}
	} else {
		b.contexts = removeElement(b.contexts, discordgo.InteractionContextBotDM)
	}
	return b.upper
}

func (b *commandBuilder[H, B]) AllowInPrivateChannel(allowed bool) B {
	if allowed {
		if !slices.Contains(b.contexts, discordgo.InteractionContextPrivateChannel) {
			b.contexts = append(b.contexts, discordgo.InteractionContextPrivateChannel)
		}
	} else {
		b.contexts = removeElement(b.contexts, discordgo.InteractionContextPrivateChannel)
	}
	return b.upper
}

func (b *commandBuilder[H, B]) AllowEverywhere(allowed bool) B {
	b.AllowInBotDM(allowed)
	b.AllowInGuilds(allowed)
	b.AllowInPrivateChannel(allowed)
	return b.upper
}
