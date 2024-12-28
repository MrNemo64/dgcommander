package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/MrNemo64/dgcommander/dgc"
	"github.com/MrNemo64/dgcommander/dgc/extras"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var repo = newRepo()

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

	commander := dgc.New(slog.Default(), ss)

	taskSelectArg := dgc.NewIntegerAutocompleteArgument().
		Name("task").
		Description("Task to perform the command on").
		Handler(autocompleteTasks).
		Required(true)

	cmd, err := commander.AddCommand(
		dgc.NewMultiSlashCommandBuilder().
			Name("tasks").
			Description("Main command").
			AllowEverywhere(true).
			AddSubCommand(dgc.NewSubCommand().
				Name("list").
				Description("List all tasks").
				AddArgument(dgc.NewUserArgument().
					Name("user").
					Description("User to run the command on. Defaults to the user running the command").
					Required(false)).
				Handler(showTasks),
			).
			AddSubCommand(dgc.NewSubCommand().
				Name("toggle").
				Description("Toggles a task, setting it as done or not").
				AddArgument(taskSelectArg).
				Handler(toggleTask),
			).
			AddSubCommand(dgc.NewSubCommand().
				Name("delete").
				Description("Deletes a task").
				AddArgument(taskSelectArg).
				Handler(deleteTask),
			).
			AddSubCommand(dgc.NewSubCommand().
				Name("create").
				Description("Creates a task").
				AddArguments(
					dgc.NewStringArgument().Name("name").Description("Task name").Required(true),
					dgc.NewStringArgument().Name("description").Description("Task description").Required(false),
					extras.NewDurationArgument().Name("duration").Description("Task duration").Required(false),
				).
				Handler(createTask),
			),
	)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := ss.ApplicationCommandDelete(ss.State.User.ID, cmd.GuildID, cmd.ID); err != nil {
			panic(err)
		}
	}()

	fmt.Println("Running")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("Clossing")

}

func autocompleteTasks(sender *discordgo.User, ctx *dgc.SlashAutocompleteContext) error {
	list := repo.getUserTasks(ctx.GetUserOr("user", sender).ID)
	field := ctx.GetStringOr("task", "")
	for _, task := range list.tasks {
		if strings.Contains(task.name, field) || strings.Contains(task.description, field) {
			ctx.AddChoice(task.name, task.id)
		}
	}
	return nil
}

func showTasks(sender *discordgo.User, ctx *dgc.SlashExecutionContext) error {
	user := ctx.GetUserOr("user", sender)
	list := repo.getUserTasks(user.ID)
	return ctx.RespondWithMessage(&discordgo.InteractionResponseData{
		Content: fmt.Sprintf("Tasks of <@%s>", user.ID),
		Embeds: mapf(list.tasks, func(task *task) *discordgo.MessageEmbed {
			return task.toEmbed()
		}),
	})
}

func toggleTask(sender *discordgo.User, ctx *dgc.SlashExecutionContext) error {
	list := repo.getUserTasks(sender.ID)
	selectedTaskId := ctx.GetRequiredInteger("task")
	task := list.findTask(selectedTaskId)
	if task == nil {
		return ctx.RespondWithMessage(&discordgo.InteractionResponseData{
			Content: fmt.Sprintf("The task with id %d does not exist", selectedTaskId),
		})
	}
	if task.completedAt == nil {
		now := time.Now()
		task.completedAt = &now
	} else {
		task.completedAt = nil
	}
	return ctx.RespondWithMessage(&discordgo.InteractionResponseData{
		Content: "Task updated",
		Embeds:  []*discordgo.MessageEmbed{task.toEmbed()},
	})
}

func deleteTask(sender *discordgo.User, ctx *dgc.SlashExecutionContext) error {
	list := repo.getUserTasks(sender.ID)
	selectedTaskId := ctx.GetRequiredInteger("task")
	msg := "The task with id %d does not exist"
	if list.deleteTask(selectedTaskId) {
		msg = "Deleted the task with id %d"
	}
	return ctx.RespondWithMessage(&discordgo.InteractionResponseData{
		Content: fmt.Sprintf(msg, selectedTaskId),
	})
}

func createTask(sender *discordgo.User, ctx *dgc.SlashExecutionContext) error {
	list := repo.getUserTasks(sender.ID)
	name := ctx.GetRequiredString("name")
	description := ctx.GetStringOr("description", "")
	var duration *time.Duration
	if durationv, found := dgc.GetArgument[time.Duration](ctx, "duration"); found {
		duration = &durationv
	}
	task := list.createTask(name, description, duration)
	return ctx.RespondWithMessage(&discordgo.InteractionResponseData{
		Content: fmt.Sprintf("Created task with id %d", task.id),
		Embeds:  []*discordgo.MessageEmbed{task.toEmbed()},
	})
}
