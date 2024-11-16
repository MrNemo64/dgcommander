package dgc

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type DGCommander struct {
	lock                  sync.RWMutex
	commands              map[string]map[string]map[discordgo.ApplicationCommandType]command // Map of name -> guild/global -> type -> command
	session               *discordgo.Session
	log                   *slog.Logger
	responseDuration      time.Duration // https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-callback
	tokenValidityDuration time.Duration
}

func New(log *slog.Logger, session *discordgo.Session) *DGCommander {
	dgc := &DGCommander{
		lock:                  sync.RWMutex{},
		commands:              make(map[string]map[string]map[discordgo.ApplicationCommandType]command),
		session:               session,
		log:                   log,
		responseDuration:      3 * time.Second,
		tokenValidityDuration: 15 * time.Minute,
	}
	session.AddHandler(dgc.manageInteraction)
	return dgc
}

type command interface {
	Manage(log *slog.Logger, sender *discordgo.User, ss *discordgo.Session, i *discordgo.InteractionCreate) (interactionAcknowledged bool, err error)
	Autocomplete(log *slog.Logger, sender *discordgo.User, ss *discordgo.Session, i *discordgo.InteractionCreate) error
}

func (c *DGCommander) AddCommand(b CommandBuilder) (*discordgo.ApplicationCommand, error) {
	definition := b.discordDefineForCreation()
	guild := b.guild()
	created, err := c.session.ApplicationCommandCreate(c.session.State.User.ID, guild, definition)
	if err != nil {
		return nil, err
	}
	c.lock.Lock()
	commands := c.getOrCreateCommandsWithNameInGuild(definition.Name, guild)
	commands[definition.Type] = b.create()
	c.lock.Unlock()
	return created, nil
}

func (c *DGCommander) getOrCreateCommandsWithName(name string) map[string]map[discordgo.ApplicationCommandType]command {
	if commands, found := c.commands[name]; found {
		return commands
	}
	commands := make(map[string]map[discordgo.ApplicationCommandType]command)
	c.commands[name] = commands
	return commands
}

func (c *DGCommander) getOrCreateCommandsWithNameInGuild(name, guild string) map[discordgo.ApplicationCommandType]command {
	commands := c.getOrCreateCommandsWithName(name)
	if commandsInGuild, found := commands[guild]; found {
		return commandsInGuild
	}
	commandsInGuild := make(map[discordgo.ApplicationCommandType]command)
	commands[guild] = commandsInGuild
	return commandsInGuild
}

func (c *DGCommander) manageInteraction(ss *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionPing || i.Type == discordgo.InteractionMessageComponent || i.Type == discordgo.InteractionModalSubmit {
		return
	}
	sender := i.Interaction.User
	if sender == nil {
		sender = i.Interaction.Member.User
	}
	log := c.log.With("sender", sender.ID, "interaction", i.Interaction) // TODO
	if i.Type == discordgo.InteractionApplicationCommand {
		data := i.ApplicationCommandData()
		c.lock.RLock()
		command, found := c.getCommandByNameInGuildAndType(data.Name, i.GuildID, data.CommandType)
		c.lock.RUnlock()
		if !found {
			log.Info("Unknown application command received", "received-command", data)
			c.respondError(ss, i.Interaction, false, fmt.Errorf("Unknown command"))
			return
		}
		interactionAcknowledged, err := command.Manage(log, sender, ss, i)
		if err != nil {
			c.respondError(ss, i.Interaction, interactionAcknowledged, err)
		}
	} else if i.Type == discordgo.InteractionApplicationCommandAutocomplete {
		data := i.ApplicationCommandData()
		c.lock.RLock()
		command, found := c.getCommandByNameInGuildAndType(data.Name, i.GuildID, discordgo.ChatApplicationCommand)
		c.lock.RUnlock()
		if !found {
			log.Info("Unknown application command autocompletion received", "received-command", data)
			c.respondError(ss, i.Interaction, false, fmt.Errorf("Unknown command"))
			return
		}
		if err := command.Autocomplete(log, sender, ss, i); err != nil {
			c.respondError(ss, i.Interaction, false, err)
		}
	} else {
		log.Warn("Unknown interaction type")
	}
}

func (c *DGCommander) getCommandByNameInGuildAndType(name, guild string, kind discordgo.ApplicationCommandType) (command, bool) {
	commandsWithName, found := c.commands[name]
	if !found {
		return nil, false
	}
	commandsInGuild, found := commandsWithName[guild]
	if !found {
		commandsInGuild, found = commandsWithName[""]
		if !found {
			return nil, false
		}
	}
	command, found := commandsInGuild[kind]
	return command, found
}

func (c *DGCommander) respondError(ss *discordgo.Session, i *discordgo.Interaction, interactionAcknowledged bool, err error) {
	embeds := []*discordgo.MessageEmbed{
		{
			Type:        discordgo.EmbedTypeArticle,
			Color:       0xff0000,
			Title:       "Error handling interaction",
			Description: err.Error(),
		},
	}
	if interactionAcknowledged {
		if _, respondErr := ss.FollowupMessageCreate(i, false, &discordgo.WebhookParams{
			Embeds: embeds,
		}); respondErr != nil {
			c.log.Warn("Error respondig to interaction with error", "interaction", i, "error-to-respond", err, "err", respondErr)
		}
	} else {
		if respondErr := ss.InteractionRespond(i, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:  discordgo.MessageFlagsEphemeral,
				Embeds: embeds,
			},
		}); respondErr != nil {
			c.log.Warn("Error respondig to interaction with error", "interaction", i, "error-to-respond", err, "err", respondErr)
		}
	}
}
