package main

import (
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
		AddArguments(
			dgc.NewStringArgument().
				Name("first-arg").
				Description("the first arg").
				Required(true),
			dgc.NewBooleanArgument().
				Name("bool-arg"),
		)

	if _, err := commander.AddCommand(builder); err != nil {
		panic(err)
	}

	fmt.Println("Running")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
}
