package dgc2

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

type command interface {
	manage(log *slog.Logger, sender *discordgo.User, ss *discordgo.Session, i *discordgo.InteractionCreate) (interactionAcknowledged bool, err error)
}
