package main

import (
	"context"
	"fmt"
	"log/slog"
	"math"
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
		dgc.NewMultiSlashCommandBuilder().
			Name().Set("calculate").
			Description().Set("Collection of simple operations").
			AllowEverywhere(true).
			AddSubCommand(dgc.NewSubCommand().
				Name().Set("sum").
				Description().Set("Calculates the sum of 2 numbers `a+b`").
				AddArguments(
					dgc.NewNumberArgument().Name().Set("a").Description().Set("First value of the sum").Required(true),
					dgc.NewNumberArgument().Name().Set("b").Description().Set("Seccond value of the sum").Required(true),
				).
				Handler(handleSum),
			).
			AddSubCommandGroup(dgc.NewSubCommandGroup().
				Name().Set("trigonometry").
				Description().Set("Trigonometry related functions").
				AddSubCommand(dgc.NewSubCommand().
					Name().Set("sin").
					Description().Set("Calculates the sin of the given angle").
					AddArguments(
						dgc.NewNumberArgument().Name().Set("angle").Description().Set("The angle to calculate the sin").Required(true),
						dgc.NewNumberChoicesArgument().Name().Set("degree").Description().Set("Degree of type to calculate").Required(false).
							AddChoice("degrees", math.Pi/180.01).
							AddChoice("radians", 1),
					).
					Handler(handleSin),
				),
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

func handleSum(ctx *dgc.SlashExecutionContext) error {
	defer ctx.Finish()
	a := ctx.GetRequiredNumber("a")
	b := ctx.GetRequiredNumber("b")
	return ctx.RespondWithMessage(&discordgo.InteractionResponseData{
		Content: fmt.Sprintf("The result of `%.2f + %.2f` is `%.2f`", a, b, a+b),
	})
}

func handleSin(ctx *dgc.SlashExecutionContext) error {
	defer ctx.Finish()
	angle := ctx.GetRequiredNumber("angle")
	degree := ctx.GetNumberOr("degree", 1) // default is radians
	angle *= degree
	return ctx.RespondWithMessage(&discordgo.InteractionResponseData{
		Content: fmt.Sprintf("The `sin(%.2f)` is `%.2f`", angle, math.Sin(angle)),
	})
}
