package main

import (
	"testing"
)

func TestAddTask(t *testing.T) {
	tasks := []Task{
		{
			ID:     1,
			Title:  "existing task",
			Status: StatusNotDone,
		},
	}

	updated := addTask(tasks, "new task")
	if updated == nil {
		t.Fatalf("expected updated tasks, got nil")
	}
	if len(updated) != 2 {
		t.Fatalf("Expected 2 tasks, got %d", len(updated))
	}
}

func TestListTasks(t *testing.T) {
	tasks := []Task{
		{
			ID:     1,
			Title:  "task 1",
			Status: StatusNotDone,
		},
		{
			ID:     2,
			Title:  "task 2",
			Status: StatusInProgress,
		},
	}

	// call the function to ensure it doesn't panic
	listTasks(tasks)
}

func TestDeleteTask(t *testing.T) {
	tasks := []Task{
		{
			ID:     1,
			Title:  "task 1",
			Status: StatusNotDone,
		},
		{
			ID:     2,
			Title:  "task 2",
			Status: StatusInProgress,
		},
	}

	// call the function to ensure it doesn't panic
	deleteTask(tasks, 1)
}

func TestRenameTask(t *testing.T) {
	tests := []struct {
		name      string
		tasks     []Task
		id        int
		newTitle  string
		wantFound bool
		wantTitle string
	}{
		{
			name:      "renames existing task",
			tasks:     []Task{{ID: 1, Title: "old title", Status: StatusNotDone}},
			id:        1,
			newTitle:  "new title",
			wantFound: true,
			wantTitle: "new title",
		},
		{
			name:      "returns false when task does not exist",
			tasks:     []Task{{ID: 1, Title: "old title", Status: StatusNotDone}},
			id:        2,
			newTitle:  "new title",
			wantFound: false,
			wantTitle: "old title",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			updated, found := renameTask(test.tasks, test.id, test.newTitle)

			if found != test.wantFound {
				t.Fatalf("expected found=%v, got %v", test.wantFound, found)
			}

			if updated[0].Title != test.wantTitle {
				t.Fatalf("expected title %q, got %q", test.wantTitle, updated[0].Title)
			}
		})
	}
}
