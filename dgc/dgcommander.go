package dgc

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type InvokationInformation struct {
	DGC          *DGCommander
	Session      *discordgo.Session
	log          *slog.Logger
	Sender       *discordgo.User
	I            *discordgo.InteractionCreate
	timeProvider TimeProvider
	ReceivedTime time.Time
}

type command interface {
	execute(*RespondingContext) (interactionAcknowledged bool, err error)
}

type autocomplete interface {
	autocomplete(*InvokationInformation) (interactionAcknowledged bool, err error)
}

type TimeProvider interface {
	Now() time.Time
}

type DefaultTimeProvider struct{}

func (DefaultTimeProvider) Now() time.Time { return time.Now() }

type CommandMiddleware = func(ctx *RespondingContext, next func()) error

type DGCommander struct {
	// TODO add some kind of middleware to commands
	lock         sync.RWMutex
	commands     map[string]map[string]map[discordgo.ApplicationCommandType]command // Map of name -> guild/global -> type -> command
	session      *discordgo.Session
	log          *slog.Logger
	timeProvider TimeProvider
	ctx          context.Context

	commandsMiddleware []CommandMiddleware
}

func New(ctx context.Context, log *slog.Logger, session *discordgo.Session, timeProvider TimeProvider) *DGCommander {
	dgc := &DGCommander{
		lock:         sync.RWMutex{},
		commands:     make(map[string]map[string]map[discordgo.ApplicationCommandType]command),
		session:      session,
		log:          log,
		timeProvider: timeProvider,
		ctx:          ctx,
	}
	session.AddHandler(dgc.manageInteraction)
	return dgc
}

func (c *DGCommander) AddCommand(b CommandBuilder) (*discordgo.ApplicationCommand, error) {
	definition := b.discordDefineForCreation()
	guild := definition.GuildID
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

func (c *DGCommander) AddMiddleware(middleware CommandMiddleware) *DGCommander {
	c.commandsMiddleware = append(c.commandsMiddleware, middleware)
	return c
}

type dgcMiddleHandler struct {
	ctx *RespondingContext
	c   *DGCommander
	ss  *discordgo.Session
	i   *discordgo.InteractionCreate
}

func (dgcmh *dgcMiddleHandler) handle() error {
	dgcmh.c.lock.RLock()
	data := dgcmh.i.ApplicationCommandData()
	command, found := dgcmh.c.getCommandByNameInGuildAndType(data.Name, dgcmh.i.GuildID, data.CommandType)
	dgcmh.c.lock.RUnlock()
	if !found {
		dgcmh.ctx.log.Info("Unknown application command received", "received-command", data)
		dgcmh.c.respondError(dgcmh.ss, dgcmh.i.Interaction, false, fmt.Errorf("Unknown command"))
		return nil
	}
	interactionAcknowledged, err := command.execute(dgcmh.ctx)
	if err != nil {
		dgcmh.c.respondError(dgcmh.ss, dgcmh.i.Interaction, interactionAcknowledged, err)
	}
	return nil
}

func (c *DGCommander) manageInteraction(ss *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionPing || i.Type == discordgo.InteractionMessageComponent || i.Type == discordgo.InteractionModalSubmit {
		return
	}
	info := InvokationInformation{
		DGC:          c,
		Session:      ss,
		log:          nil,
		Sender:       i.Interaction.User,
		I:            i,
		timeProvider: c.timeProvider,
		ReceivedTime: c.timeProvider.Now(),
	}
	if info.Sender == nil {
		info.Sender = i.Interaction.Member.User
	}
	info.log = c.log.With("sender", info.Sender.ID) // TODO add information about the interaction to be able to identify it
	if i.Type == discordgo.InteractionApplicationCommand {
		ctx := RespondingContext{
			executionContext: newExecutionContext(c.ctx, &info),
		}

		dgcmh := dgcMiddleHandler{
			ctx: &ctx,
			c:   c,
			ss:  ss,
			i:   i,
		}

		chain := newMiddlewareChain(&ctx, c.commandsMiddleware, dgcmh.handle)
		if err := chain.startChain(); err != nil {
			c.respondError(ss, i.Interaction, ctx.acknowledged, err)
			return
		}
	} else if i.Type == discordgo.InteractionApplicationCommandAutocomplete {
		data := i.ApplicationCommandData()
		c.lock.RLock()
		command, found := c.getCommandByNameInGuildAndType(data.Name, i.GuildID, discordgo.ChatApplicationCommand)
		c.lock.RUnlock()
		if !found {
			info.log.Info("Unknown application command autocompletion received", "received-command", data)
			c.respondError(ss, i.Interaction, false, fmt.Errorf("Unknown command"))
			return
		}
		ac, ok := command.(autocomplete)
		if !ok {
			info.log.Info("Received a request to autocomplete a non autocompletable command", "received-command", data)
			c.respondError(ss, i.Interaction, false, fmt.Errorf("Cannot autocomplete"))
		}
		interactionAcknowledged, err := ac.autocomplete(&info)
		if err != nil {
			c.respondError(ss, i.Interaction, interactionAcknowledged, err)
		}
	} else {
		info.log.Warn("Unknown interaction type")
	}
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
