package main

import (
	"fmt"
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
