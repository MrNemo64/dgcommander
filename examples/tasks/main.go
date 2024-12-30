package main

import (
	"context"
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

	commander := dgc.New(context.Background(), slog.Default(), ss, dgc.DefaultTimeProvider{})

	taskSelectArg := dgc.NewIntegerAutocompleteArgument().
		Name().Set("task").
		Description().Set("Task to perform the command on").
		Handler(autocompleteTasks).
		Required(true)

	cmd, err := commander.AddCommand(
		dgc.NewMultiSlashCommandBuilder().
			Name().Set("tasks").
			Description().Set("Main command").
			AllowEverywhere(true).
			AddSubCommand(dgc.NewSubCommand().
				Name().Set("list").
				Description().Set("List all tasks").
				AddArgument(dgc.NewUserArgument().
					Name().Set("user").
					Description().Set("User to run the command on. Defaults to the user running the command").
					Required(false)).
				Handler(showTasks),
			).
			AddSubCommand(dgc.NewSubCommand().
				Name().Set("toggle").
				Description().Set("Toggles a task, setting it as done or not").
				AddArgument(taskSelectArg).
				Handler(toggleTask),
			).
			AddSubCommand(dgc.NewSubCommand().
				Name().Set("delete").
				Description().Set("Deletes a task").
				AddArgument(taskSelectArg).
				Handler(deleteTask),
			).
			AddSubCommand(dgc.NewSubCommand().
				Name().Set("create").
				Description().Set("Creates a task").
				AddArguments(
					dgc.NewStringArgument().Name().Set("name").Description().Set("Task name").Required(true),
					dgc.NewStringArgument().Name().Set("description").Description().Set("Task description").Required(false),
					extras.NewDurationArgument().Name().Set("duration").Description().Set("Task duration").Required(false),
				).
				Handler(createTask),
			).
			AddSubCommand(dgc.NewSubCommand().
				Name().Set("edit").
				Description().Set("Edit a task").
				AddArguments(
					taskSelectArg,
					dgc.NewStringArgument().Name().Set("name").Description().Set("Task name").Required(false),
					dgc.NewStringArgument().Name().Set("description").Description().Set("Task description").Required(false),
					extras.NewDurationArgument().Name().Set("duration").Description().Set("Task duration").Required(false),
				).
				Handler(editTask),
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

func autocompleteTasks(ctx *dgc.SlashAutocompleteContext) error {
	defer ctx.Finish()
	list := repo.getUserTasks(ctx.GetUserOr("user", ctx.Sender).ID)
	field := ctx.GetStringOr("task", "")
	for _, task := range list.tasks {
		if strings.Contains(task.name, field) || strings.Contains(task.description, field) {
			ctx.AddChoice(task.name, task.id)
		}
	}
	return nil
}

func showTasks(ctx *dgc.SlashExecutionContext) error {
	defer ctx.Finish()
	user := ctx.GetUserOr("user", ctx.Sender)
	list := repo.getUserTasks(user.ID)
	return ctx.RespondWithMessage(&discordgo.InteractionResponseData{
		Content: fmt.Sprintf("Tasks of <@%s>", user.ID),
		Embeds: mapf(list.tasks, func(task *task) *discordgo.MessageEmbed {
			return task.toEmbed()
		}),
	})
}

func toggleTask(ctx *dgc.SlashExecutionContext) error {
	defer ctx.Finish()
	list := repo.getUserTasks(ctx.Sender.ID)
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

func deleteTask(ctx *dgc.SlashExecutionContext) error {
	defer ctx.Finish()
	list := repo.getUserTasks(ctx.Sender.ID)
	selectedTaskId := ctx.GetRequiredInteger("task")
	msg := "The task with id %d does not exist"
	if list.deleteTask(selectedTaskId) {
		msg = "Deleted the task with id %d"
	}
	return ctx.RespondWithMessage(&discordgo.InteractionResponseData{
		Content: fmt.Sprintf(msg, selectedTaskId),
	})
}

func createTask(ctx *dgc.SlashExecutionContext) error {
	defer ctx.Finish()
	list := repo.getUserTasks(ctx.Sender.ID)
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

func editTask(ctx *dgc.SlashExecutionContext) error {
	defer ctx.Finish()
	list := repo.getUserTasks(ctx.Sender.ID)
	selectedTaskId := ctx.GetRequiredInteger("task")
	task := list.findTask(selectedTaskId)
	if task == nil {
		return ctx.RespondWithMessage(&discordgo.InteractionResponseData{
			Content: fmt.Sprintf("The task with id %d does not exist", selectedTaskId),
		})
	}
	task.name = ctx.GetStringOr("name", task.name)
	task.description = ctx.GetStringOr("description", task.description)
	if durationv, found := dgc.GetArgument[time.Duration](ctx, "duration"); found {
		if durationv > 0 {
			task.duration = &durationv
		} else {
			task.duration = nil
		}
	}
	return ctx.RespondWithMessage(&discordgo.InteractionResponseData{
		Content: "Task updated",
		Embeds:  []*discordgo.MessageEmbed{task.toEmbed()},
	})
}
