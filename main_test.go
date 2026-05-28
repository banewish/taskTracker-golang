package main

import (
	"testing"
)

func TestAddTask(t *testing.T) {
	tasks := []struct{
		ID int
		Title string
		Done bool
	}{
		{
			ID: 		1, 
			Title: 		"existing task", 
			Done: 		false,
		}}

	if updated := addTask(tasks, "new task")
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
			tasks:     []Task{{ID: 1, Title: "old title", Done: false}},
			id:        1,
			newTitle:  "new title",
			wantFound: true,
			wantTitle: "new title",
		},
		{
			name:      "returns false when task does not exist",
			tasks:     []Task{{ID: 1, Title: "old title", Done: false}},
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
