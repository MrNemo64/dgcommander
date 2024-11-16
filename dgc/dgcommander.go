package dgc

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/bwmarrin/discordgo"
)

type DGCommander struct {
	lock     sync.RWMutex
	commands map[string]map[string]map[discordgo.ApplicationCommandType]command // Map of name -> guild/global -> type -> command
	session  *discordgo.Session
	log      *slog.Logger
}

func New(log *slog.Logger, session *discordgo.Session) *DGCommander {
	return &DGCommander{
		lock:     sync.RWMutex{},
		commands: make(map[string]map[string]map[discordgo.ApplicationCommandType]command),
		session:  session,
		log:      log,
	}
}

type command interface {
	manage(log *slog.Logger, ss *discordgo.Session, i *discordgo.InteractionCreate) error
	autocomplete(log *slog.Logger, ss *discordgo.Session, i *discordgo.InteractionCreate) error
}

func (c *DGCommander) AddCommand(b CommandBuilder) (*discordgo.ApplicationCommand, error) {
	definition := b.discordDefineForCreation()
	guild := b.guild()
	created, err := c.session.ApplicationCommandCreate(c.session.State.User.ID, guild, definition)
	if err != nil {
		return nil, err
	}
	// commands := c.getOrCreateCommandsWithNameInGuild(definition.Name, guild)
	// commands[definition.Type] = b.create()
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

func (c *DGCommander) ManageInteraction(ss *discordgo.Session, i *discordgo.InteractionCreate) {
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
		c.lock.Lock()
		command, found := c.getCommandByNameInGuildAndType(data.Name, i.GuildID, data.CommandType)
		c.lock.Unlock()
		if !found {
			log.Info("Unknown application command received", "received-command", data)
			c.respondError(ss, i.Interaction, fmt.Errorf("Unknown command"))
			return
		}
		if err := command.manage(log, ss, i); err != nil {
			c.respondError(ss, i.Interaction, err)
		}
	} else if i.Type == discordgo.InteractionApplicationCommandAutocomplete {
		data := i.ApplicationCommandData()
		c.lock.Lock()
		command, found := c.getCommandByNameInGuildAndType(data.Name, i.GuildID, discordgo.ChatApplicationCommand)
		c.lock.Unlock()
		if !found {
			log.Info("Unknown application command autocompletion received", "received-command", data)
			c.respondError(ss, i.Interaction, fmt.Errorf("Unknown command"))
			return
		}
		if err := command.autocomplete(log, ss, i); err != nil {
			c.respondError(ss, i.Interaction, err)
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

func (c *DGCommander) respondError(ss *discordgo.Session, i *discordgo.Interaction, err error) {
	if respondErr := ss.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				{
					Type:        discordgo.EmbedTypeArticle,
					Color:       0xf00,
					Title:       "Error handling interaction",
					Description: err.Error(),
				},
			},
		},
	}); respondErr != nil {
		c.log.Warn("Error respondig to interaction with error", "interaction", i, "error-to-respond", err, "err", respondErr)
	}
}

func (c *DGCommander) KnownCommandNames() []string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	known := make([]string, len(c.commands))
	i := 0
	for name := range c.commands {
		known[i] = name
		i++
	}
	return known
}
