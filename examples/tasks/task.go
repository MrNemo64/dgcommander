package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type task struct {
	id          int64
	name        string
	description string
	duration    *time.Duration
	completedAt *time.Time
}

func (t task) toEmbed() *discordgo.MessageEmbed {
	color := 0x0000ff
	fields := []*discordgo.MessageEmbedField{
		{
			Name:  "Id",
			Value: strconv.FormatInt(t.id, 10),
		},
	}
	if t.duration != nil {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "Duration",
			Value: t.duration.String(),
		})
	}
	if t.completedAt != nil {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "Completed At",
			Value: fmt.Sprintf("<t:%d:f>", t.completedAt.Unix()),
		})
		color = 0x00ff00
	}

	return &discordgo.MessageEmbed{
		Title:       t.name,
		Description: t.description,
		Fields:      fields,
		Color:       color,
	}
}

type taskList struct {
	tasks  []*task
	nextId int64
}

func (l *taskList) findTask(id int64) *task {
	for _, task := range l.tasks {
		if task.id == id {
			return task
		}
	}
	return nil
}

func (l *taskList) deleteTask(id int64) bool {
	tasks := make([]*task, 0)
	deleted := false
	for _, task := range l.tasks {
		if task.id != id {
			tasks = append(tasks, task)
		} else {
			deleted = true
		}
	}
	l.tasks = tasks
	return deleted
}

func (l *taskList) createTask(name, description string, duration *time.Duration) *task {
	task := &task{
		id:          l.nextId,
		name:        name,
		description: description,
		duration:    duration,
	}
	l.tasks = append(l.tasks, task)
	l.nextId++
	return task
}

type taskRepo struct {
	tasks map[string]*taskList
	m     sync.Mutex
}

func (r *taskRepo) getUserTasks(id string) *taskList {
	r.m.Lock()
	defer r.m.Unlock()
	list, found := r.tasks[id]
	if found {
		return list
	}
	list = &taskList{
		nextId: 4,
		tasks: []*task{
			{
				id:          1,
				name:        "Run command",
				description: "Run the main command",
			},
			{
				id:          2,
				name:        "Leave a star",
				description: "Star the repository",
			},
			{
				id:          3,
				name:        "Try it out",
				description: "Try out dgcommander",
			},
		},
	}
	r.tasks[id] = list
	return list
}

func newRepo() *taskRepo {
	return &taskRepo{
		tasks: make(map[string]*taskList),
		m:     sync.Mutex{},
	}
}
