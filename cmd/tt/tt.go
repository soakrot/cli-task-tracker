package main

import (
	"fmt"
	"github.com/soakrot/cli-task-tracker/internal/command"
	"github.com/soakrot/cli-task-tracker/internal/store"
	"os"
	"text/tabwriter"
)

const storePath string = "$XDG_DATA_HOME/task-tracker/tasks.json"

const (
	bold  = "\033[1m"
	reset = "\033[0m"
)

var s = store.Store{Tasks: make(map[int]*store.Task), NextID: 1}

var commands = map[string]func([]string) error{
	"add": func(args []string) error {
		add := command.AddCmd(args)
		id, err := s.AddTask(add.Title)
		fmt.Println(id)
		return err
	},
	"update": func(args []string) error {
		update := command.UpdateCmd(args)
		err := s.UpdateTask(update.Id, update.Title)
		return err
	},
	"delete": func(args []string) error {
		del := command.DeleteCmd(args)
		out, err := s.DeleteTask(del.Id)
		fmt.Println(out)
		return err
	},
	"delete-all": func(args []string) error {
		delAllArgs := command.DeleteAllCmd(args)
		err := s.DeleteAll(delAllArgs.Status)
		return err
	},
	"mark": func(args []string) error {
		mark := command.MarkCmd(args)
		err := s.MarkTask(mark.Id, mark.Status)
		return err
	},
	"list": func(args []string) error {
		list := command.ListCmd(args)
		err := s.ListTasks(list.Status)
		return err
	},
}

func main() {
	if len(os.Args) == 1 {
		printUsage()
		os.Exit(1)
	}

	err := store.LoadData(&s)
	if err != nil {
		fmt.Println(fmt.Errorf("Error while loading data: %w", err))
	}

	cmd, ok := commands[os.Args[1]]
	if !ok {
		printUsage()
		os.Exit(1)
	}

	err = cmd(os.Args[2:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func printUsage() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, bold+"Usage:"+reset)
	fmt.Fprintln(w, "  "+bold+"task-tacker "+reset+"<command> <title>")
	fmt.Fprintln(w, "  "+bold+"task-tacker "+reset+"update <task-id> <title>")
	fmt.Fprintln(w, "  "+bold+"task-tacker "+reset+"delete <task-id>")
	fmt.Fprintln(w, "  "+bold+"task-tacker "+reset+"list [<done | todo | in-progress>]")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, bold+"Commands:"+reset)
	fmt.Fprintln(w, "  "+bold+"add\t"+reset+"<title>\tAdd a new task and return its ID")
	fmt.Fprintln(w, "  "+bold+"update\t"+reset+"<id> <title>\tUpdate the title of task \033[1mid\033[0m")
	fmt.Fprintln(w, "  "+bold+"delete\t"+reset+"<id>\tDelete task \033[1mid\033[0m permanently")
	fmt.Fprintln(w, "  "+bold+"list\t"+reset+"<status>\tList all tasks, or filter by \033[1mstatus\033[0m")
	fmt.Fprintln(w, "  "+bold+"mark\t"+reset+"<id> <status>\tMark task as todo | in-progress | done")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, bold+"Options:"+reset)
	fmt.Fprintln(w, "  "+bold+"-h, --help:"+reset+"\tPrint help")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, bold+"Status values:\t"+reset+"done | todo | in-progress")
	w.Flush()
}
