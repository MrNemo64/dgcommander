package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/MrNemo64/dgcommander/dgc"
	"github.com/MrNemo64/dgcommander/dgc/handlers"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting messages example")
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	ss, err := discordgo.New("Bot " + os.Getenv("EXAMPLES_TOKEN"))
	if err != nil {
		panic(err)
	}
	if err := ss.Open(); err != nil {
		panic(err)
	}
	defer ss.Close()

	commander := dgc.New(slog.Default(), ss)

	cmd, err := commander.AddCommand(
		dgc.NewMessage().
			Name("Resend message").
			Handler(handleResend).
			AllowEverywhere(true),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Running")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("Clossing")

	if err := ss.ApplicationCommandDelete(ss.State.User.ID, cmd.GuildID, cmd.ID); err != nil {
		panic(err)
	}
}

func handleResend(ctx *handlers.MessageExecutionContext, sender *discordgo.User) error {
	fmt.Printf("Called by %s (%s) on message %s\n", sender.Username, sender.ID, ctx.Message.ID)
	return ctx.RespondWithMessage(&discordgo.InteractionResponseData{
		Content: ctx.Message.Content,
		Embeds:  ctx.Message.Embeds,
	})
}
