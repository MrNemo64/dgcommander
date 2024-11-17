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
	fmt.Println("Running choices example")
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

	builder := dgc.NewSimpleSlash().
		Name("string-choices").
		Description("example command").
		AddArguments(
			dgc.NewStringArgument().
				Name("first-arg").
				Description("the first arg").
				Required(true),
			dgc.NewBooleanArgument().
				Name("bool-arg").
				Description("the seccond arg").
				Required(false),
		).
		Handler(func(sender *discordgo.User, ctx *handlers.SlashExecutionContext) error {
			fmt.Println("Called")
			return nil
		})

	cmd, err := commander.AddCommand(builder)
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
