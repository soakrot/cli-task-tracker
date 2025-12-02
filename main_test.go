package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestAddTask(t *testing.T) {
	store := Store{Tasks: make(map[uint]*Task), NextID: 1}

	tests := []struct {
		id          uint
		description string
		status      string
		want        uint
	}{
		{1, "task 1", "todo", 1},
		{2, "task 2", "todo", 2},
		{3, "task 3", "todo", 3},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s, %s", tt.description, tt.status)
		t.Run(testname, func(t *testing.T) {
			res, _ := store.AddTask(tt.description)
			if res == -1 {
				t.Errorf("got %d, want %d", res, tt.want)
			}
		})
	}

}

func TestUpdateTask(t *testing.T) {
	store := Store{Tasks: make(map[uint]*Task), NextID: 1}
	store.AddTask("Task 1")
	store.AddTask("Task 2")
	store.AddTask("Task 3")

	tests := []struct {
		id   uint
		c    string
		want string
	}{
		{1, "something else", "something else"},
		{2, "something else 2", "something else 2"},
		{3, "ummm", "ummm"},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.c)
		t.Run(testname, func(t *testing.T) {
			err := store.UpdateTask(tt.id, tt.c)
			task := store.Tasks[tt.id].Description
			if err != nil {
				t.Errorf("error: %s", err)
			}

			if task != tt.want {
				t.Errorf("got %s, want %s", task, tt.want)
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	store := Store{Tasks: make(map[uint]*Task), NextID: 1}
	store.AddTask("Task 1")
	store.AddTask("Task 2")
	store.AddTask("Task 3")

	tests := []struct {
		id   uint
		want string
	}{
		{1, "task 1"},
		{2, "task 2"},
		{3, "task 3"},
		{4, ""},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.want)
		t.Run(testname, func(t *testing.T) {
			out, _ := store.DeleteTask(tt.id)

			if strings.ToLower(out) != tt.want {
				t.Errorf("got %q, want %q", out, tt.want)
			}
		})
	}
}
