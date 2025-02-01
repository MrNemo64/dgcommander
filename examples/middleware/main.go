package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MrNemo64/dgcommander/dgc"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting middleware example")
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

	commander.AddMiddleware(func(ctx *dgc.RespondingContext, next func()) error {
		t := time.Now()
		fmt.Println("Pre next middleware")
		next()
		fmt.Printf("Post next middleware, took %dms\n", time.Since(t).Milliseconds())
		return nil
	})

	cmd, err := commander.AddCommand(
		dgc.NewMessageCommand().
			Name().Set("Message command").
			AddMiddleware(func(ctx *dgc.MessageExecutionContext, next func()) error {
				fmt.Println("Pre next message middleware")
				next()
				fmt.Println("Post next message middleware")
				return nil
			}).
			Handler(handler).
			AllowEverywhere(true),
	)
	if err != nil {
		panic(err)
	}
	cmd, err = commander.AddCommand(
		dgc.NewUserCommand().
			Name().Set("User command").
			AddMiddleware(func(ctx *dgc.UserExecutionContext, next func()) error {
				fmt.Println("Pre next user middleware")
				next()
				fmt.Println("Post next user middleware")
				return nil
			}).
			Handler(handler).
			AllowEverywhere(true),
	)
	if err != nil {
		panic(err)
	}
	cmd, err = commander.AddCommand(
		dgc.NewMultiSlashCommandBuilder().
			Name().Set("slash-multi-command").
			Description().Set("Description").
			AddMiddleware(func(ctx *dgc.RespondingContext, next func()) error {
				fmt.Println("Pre next multi middleware")
				next()
				fmt.Println("Post next multi middleware")
				return nil
			}).
			AllowEverywhere(true).
			AddSubCommand(dgc.NewSubCommand().
				Name().Set("slash-multi-single").
				Description().Set("Description").
				AddMiddleware(func(ctx *dgc.SlashExecutionContext, next func()) error {
					fmt.Println("Pre next multi single middleware 1")
					next()
					fmt.Println("Post next multi single middleware 1")
					return nil
				}).
				AddMiddleware(func(ctx *dgc.SlashExecutionContext, next func()) error {
					fmt.Println("Pre next multi single middleware 2")
					next()
					fmt.Println("Post next multi single middleware 2")
					return nil
				}).
				Handler(handler),
			),
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

func handler[T interface {
	Finish()
	RespondWithMessage(*discordgo.InteractionResponseData) error
}](ctx T) error {
	defer ctx.Finish()
	fmt.Println("Called a command")
	return ctx.RespondWithMessage(&discordgo.InteractionResponseData{
		Content: "Called a command",
	})
}
