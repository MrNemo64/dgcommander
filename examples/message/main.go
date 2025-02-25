package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/MrNemo64/dgcommander/dgc"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting messages example")
	if err := godotenv.Load("../../.env"); err != nil {
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

	commander := dgc.New(context.Background(), slog.Default(), ss, dgc.DefaultTimeProvider{})

	cmd, err := commander.AddCommand(
		dgc.NewMessageCommand().
			Name().Set("Resend message").
			Handler(handleResend).
			AllowEverywhere(true),
	)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := ss.ApplicationCommandDelete(ss.State.User.ID, cmd.GuildID, cmd.ID); err != nil {
			panic(err)
		}
		fmt.Println("Deleted command")
	}()

	fmt.Println("Running")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("Clossing")
}

func handleResend(ctx *dgc.MessageExecutionContext) error {
	defer ctx.Finish()
	fmt.Printf("Called by %s (%s) on message %s\n", ctx.Sender.Username, ctx.Sender.ID, ctx.Message.ID)
	return ctx.RespondWithMessage(&discordgo.InteractionResponseData{
		Content: ctx.Message.Content,
		Embeds:  ctx.Message.Embeds,
	})
}
