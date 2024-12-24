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
		dgc.NewMultiSlashCommandBuilder().
			Name("calculate").
			Description("Collection of simple operations").
			AllowEverywhere(true).
			AddSubCommand(dgc.NewSubCommand().
				Name("sum").
				Description("Calculates the sum of 2 numbers `a+b`").
				AddArguments(
					dgc.NewNumberArgument().Name("a").Description("First value of the sum").Required(true),
					dgc.NewNumberArgument().Name("b").Description("Seccond value of the sum").Required(true),
				).
				Handler(handleSum),
			),
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

func handleSum(sender *discordgo.User, ctx *dgc.SlashExecutionContext) error {
	a := ctx.GetRequiredNumber("a")
	b := ctx.GetRequiredNumber("b")
	return ctx.RespondWithMessage(&discordgo.InteractionResponseData{
		Content: fmt.Sprintf("The result of `%.2f + %.2f` is `%.2f`", a, b, a+b),
	})
}
