package main

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
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

	cmds, err := ss.ApplicationCommands(ss.State.Application.ID, "")
	if err != nil {
		panic(err)
	}
	for _, cmd := range cmds {
		if err := ss.ApplicationCommandDelete(ss.State.Application.ID, "", cmd.ID); err != nil {
			panic(err)
		}
		fmt.Printf("Deleted command %s\n", cmd.Name)
	}
	if _, err := fmt.Println("Deleted commands"); err != nil {
		panic(err)
	}
}
