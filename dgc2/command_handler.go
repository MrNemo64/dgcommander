package dgc2

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

type executionContext struct {
	Session          *discordgo.Session
	I                *discordgo.InteractionCreate
	log              *slog.Logger
	alreadyResponded bool
}
