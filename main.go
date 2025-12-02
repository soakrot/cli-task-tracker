package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

const storePath string = "$XDG_DATA_HOME/task-tracker/tasks.json"

type taskStatus string

const (
	todo       taskStatus = "todo"
	inProgress taskStatus = "in-progress"
	done       taskStatus = "done"
)

type Task struct {
	Id          uint
	CreatedAt   int64
	UpdatedAt   int64
	Description string
	Status      taskStatus
}

type Tasks map[uint]*Task

type Store struct {
	Tasks  Tasks `json:"tasks"`
	NextID uint  `json:"nextID"`
}

const (
	bold  = "\033[1m"
	reset = "\033[0m"
)

func printHelp() {
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

func (t *Task) setStatus(status string) error {
	if ok, err := isValidStatus(status); !ok {
		return err
	}
	t.Status = taskStatus(status)
	return nil
}

func isValidStatus(status string) (bool, error) {
	switch taskStatus(strings.Trim(status, " ")) {
	case todo, inProgress, done:
		return true, nil
	default:
		return false, errors.New("Invalid task status")
	}
}

func createTask(id uint, description, status string) (*Task, error) {
	task := Task{
		Id:          id,
		Description: description,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}
	if err := task.setStatus(status); err != nil {
		return nil, fmt.Errorf("Error when creating a task: %w", err)
	}
	return &task, nil
}

func (s *Store) AddTask(description string) (int, error) {
	if len(strings.Trim(description, " ")) == 0 {
		return -1, errors.New("Invalid or empty title")
	}

	newTask, err := createTask(s.NextID, description, "todo")
	if err != nil {
		return -1, fmt.Errorf("Error while adding a task: %w", err)
	}
	s.Tasks[newTask.Id] = newTask
	s.NextID++

	return int(newTask.Id), nil
}

func (s *Store) UpdateTask(id uint, c string) error {
	if len(s.Tasks) == 0 {
		return errors.New("There are no tasks, start working!")
	}

	if _, ok := s.Tasks[id]; !ok {
		return errors.New("Invalid task ID")
	}

	if len(strings.Trim(c, " ")) == 0 {
		return errors.New("Empty content")
	}

	s.Tasks[id].Description = c
	s.Tasks[id].UpdatedAt = time.Now().Unix()

	return nil
}

func printTask(writer *tabwriter.Writer, task *Task) {
	fmt.Fprintf(writer,
		"%d\t%q\t%s\t%s\t%s\n",
		task.Id,
		task.Description,
		task.Status,
		time.Unix(task.CreatedAt, 0).Format(time.RFC1123),
		time.Unix(task.UpdatedAt, 0).Format(time.RFC1123),
	)
}

func printTasks(tasks Tasks, status string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "id\tdescription\tstatus\tcreated-at\tupdated-at")
	for _, task := range tasks {
		if task.Status == taskStatus(status) {
			printTask(w, task)
		}
		if len(status) == 0 {
			printTask(w, task)
		}
	}
	w.Flush()
}

func (s *Store) ListTasks(status string) error {
	if len(s.Tasks) == 0 {
		return errors.New("There are no tasks, start working!")
	}

	if len(strings.Trim(string(status), " ")) == 0 {
		printTasks(s.Tasks, "")
		return nil
	}

	if ok, err := isValidStatus(status); !ok {
		return fmt.Errorf("list: %w: '%s'\n", err, status)
	}

	printTasks(s.Tasks, status)
	return nil
}

func isFileExists(path string) bool {
	_, err := os.Lstat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func loadData(s *Store) error {
	dataPath := os.ExpandEnv(filepath.Join(os.TempDir(), "tasks.json"))
	// dataPath := os.ExpandEnv(storePath)

	if !isFileExists(dataPath) {
		if err := os.MkdirAll(filepath.Dir(dataPath), 0o766); err != nil {
			return errors.New("error occured while creating a directory")
		}

		f, _ := os.Create(dataPath)
		fmt.Println("created ", dataPath)

		b, err := json.Marshal(s)
		if err != nil {
			return err
		}
		w := bufio.NewWriter(f)
		if _, err := w.Write(b); err != nil {
			return err
		}
		w.Flush()

		defer f.Close()
		return nil
	}

	data, err := os.ReadFile(dataPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, s); err != nil {
		return err
	}

	return nil
}

func writeData(s *Store) error {
	dataPath := filepath.Join(os.TempDir(), "tasks.json")

	data, err := json.Marshal(s)
	if err != nil {
		return err
	}

	errF := os.WriteFile(dataPath, data, 0644)
	if errF != nil {
		panic(errF)
	}

	return nil
}

func main() {
	if len(os.Args) == 1 {
		printHelp()
		os.Exit(1)
	}

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	store := Store{Tasks: make(map[uint]*Task), NextID: 1}
	err := loadData(&store)
	if err != nil {
		fmt.Println(fmt.Errorf("Error while loading data: %w", err))
	}

	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		input := addCmd.Arg(0)
		id, err := store.AddTask(input)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		writeData(&store)
		fmt.Println(id)
	case "update":
		updateCmd.Parse(os.Args[2:])
		content := updateCmd.Arg(1)
		id, err := strconv.ParseUint(updateCmd.Arg(0), 10, 0)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := store.UpdateTask(uint(id), content); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		writeData(&store)
	case "delete":
		// TODO: implement task deletion
		writeData(&store)
	case "mark":
		// TODO: implement task marking
		writeData(&store)
	case "list":
		listCmd.Parse(os.Args[2:])
		status := listCmd.Arg(0)
		if err := store.ListTasks(status); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		os.Exit(1)
	}
}
