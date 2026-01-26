package store

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
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
	Id          int
	CreatedAt   int64
	UpdatedAt   int64
	Description string
	Status      taskStatus
}

type Tasks map[int]*Task

type Store struct {
	Tasks  Tasks `json:"tasks"`
	NextID int   `json:"nextID"`
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
		return false, errors.New("Invalid status")
	}
}

func createTask(id int, description, status string) (*Task, error) {
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
	WriteData(s)

	return int(newTask.Id), nil
}

func (s *Store) UpdateTask(id int, c string) error {
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
	WriteData(s)

	return nil
}

func (s *Store) DeleteTask(id int) (string, error) {
	if len(s.Tasks) == 0 {
		return "", errors.New("There are no tasks, start working!")
	}

	out, ok := s.Tasks[id]
	if !ok {
		return "", errors.New("Invalid task ID")
	}

	if len(s.Tasks) == 1 {
		s.Tasks = map[int]*Task{}
		s.NextID = 1
		WriteData(s)
		return out.Description, nil
	}

	delete(s.Tasks, id)
	WriteData(s)

	return out.Description, nil
}

func (s *Store) DeleteAll(status string) error {
	skipStatusCheck := false
	if len(status) == 0 {
		skipStatusCheck = true
	}

	if ok, err := isValidStatus(status); !ok && !skipStatusCheck {
		return err
	}

	s.Tasks = map[int]*Task{}
	s.NextID = 1
	WriteData(s)

	return nil
}

func (s *Store) MarkTask(id int, status string) error {
	if len(s.Tasks) == 0 {
		return errors.New("there are no tasks, start working!")
	}

	if _, ok := s.Tasks[id]; !ok {
		return errors.New("Invalid task ID")
	}

	if ok, err := isValidStatus(status); !ok {
		return err
	}

	s.Tasks[id].setStatus(status)
	WriteData(s)

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

func LoadData(s *Store) error {
	dataPath := os.ExpandEnv(storePath)

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

func WriteData(s *Store) error {
	dataPath := os.ExpandEnv(storePath)

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
