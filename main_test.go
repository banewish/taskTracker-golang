package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestAddTask(t *testing.T) {
	tasks := []Task{
		{
			ID:        1,
			Title:     "existing task",
			Status:    StatusNotDone,
			CreatedAt: taskTimestamp(),
			UpdatedAt: taskTimestamp(),
		},
	}

	updated := addTask(tasks, "new task")
	if updated == nil {
		t.Fatalf("expected updated tasks, got nil")
	}
	if len(updated) != 2 {
		t.Fatalf("Expected 2 tasks, got %d", len(updated))
	}

	if updated[1].ID != 2 {
		t.Fatalf("expected new task ID to be 2, got %d", updated[1].ID)
	}

	if updated[1].CreatedAt.IsZero() || updated[1].UpdatedAt.IsZero() {
		t.Fatal("expected new task to have timestamps")
	}

	if !updated[1].CreatedAt.Equal(updated[1].UpdatedAt) {
		t.Fatal("expected new task created and updated timestamps to match")
	}
}

func TestAddTaskAfterDeleteKeepsIDsUnique(t *testing.T) {
	tasks := []Task{
		{ID: 1, Title: "task 1", Status: StatusNotDone},
		{ID: 2, Title: "task 2", Status: StatusInProgress},
	}

	updated, found := deleteTask(tasks, 1)
	if !found {
		t.Fatal("expected task to be found")
	}

	updated = addTask(updated, "task 3")

	if len(updated) != 2 {
		t.Fatalf("expected 2 tasks after delete and add, got %d", len(updated))
	}

	if updated[0].ID != 2 {
		t.Fatalf("expected remaining task ID to stay 2, got %d", updated[0].ID)
	}

	if updated[1].ID != 3 {
		t.Fatalf("expected new task ID to be 3, got %d", updated[1].ID)
	}
}

func TestListTasks(t *testing.T) {
	tasks := []Task{{
		ID:        1,
		Title:     "task 1",
		Status:    StatusNotDone,
		CreatedAt: taskTimestamp(),
		UpdatedAt: taskTimestamp(),
	}}

	output := captureOutput(t, func() {
		listTasks(tasks)
	})

	if !strings.Contains(output, "task 1") {
		t.Fatalf("expected task title in output, got %q", output)
	}

	if !strings.Contains(output, string(StatusNotDone)) {
		t.Fatalf("expected task status in output, got %q", output)
	}
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

	updated, found := deleteTask(tasks, 1)
	if !found {
		t.Fatal("expected task to be found")
	}

	if len(updated) != 1 {
		t.Fatalf("expected 1 task after delete, got %d", len(updated))
	}

	if updated[0].ID != 2 {
		t.Fatalf("expected remaining task ID to stay 2, got %d", updated[0].ID)
	}

	if updated[0].Title != "task 2" {
		t.Fatalf("expected remaining task title to be %q, got %q", "task 2", updated[0].Title)
	}
}

func TestSetTaskStatus(t *testing.T) {
	tasks := []Task{{ID: 1, Title: "task 1", Status: StatusNotDone}}

	updated, found := setTaskStatus(tasks, 1, StatusInProgress)
	if !found {
		t.Fatal("expected task to be found")
	}

	if updated[0].Status != StatusInProgress {
		t.Fatalf("expected status %q, got %q", StatusInProgress, updated[0].Status)
	}

	if updated[0].UpdatedAt.IsZero() {
		t.Fatal("expected updated timestamp to be set")
	}
}

func TestParseTaskStatus(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    TaskStatus
		wantErr bool
	}{
		{name: "done", input: "done", want: StatusDone},
		{name: "in progress", input: "in-progress", want: StatusInProgress},
		{name: "not done", input: "not-done", want: StatusNotDone},
		{name: "invalid", input: "broken", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parseTaskStatus(test.input)
			if test.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}

			if err != nil {
				t.Fatal(err)
			}

			if got != test.want {
				t.Fatalf("expected %q, got %q", test.want, got)
			}
		})
	}
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

func TestLoadTasksRejectsInvalidStatus(t *testing.T) {
	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(oldwd); err != nil {
			t.Fatal(err)
		}
	}()

	data := []byte(`[{"id":1,"title":"bad task","status":"broken"}]`)
	if err := os.WriteFile(tasksFile, data, 0o644); err != nil {
		t.Fatal(err)
	}

	_, err = loadTasks()
	if err == nil {
		t.Fatal("expected loadTasks to fail for invalid status")
	}

	if !strings.Contains(err.Error(), "invalid task status") {
		t.Fatalf("expected invalid status error, got %v", err)
	}
}

func TestLoadTasksSupportsLegacyDoneField(t *testing.T) {
	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(oldwd); err != nil {
			t.Fatal(err)
		}
	}()

	data := []byte(`[{"id":1,"title":"legacy task","done":true}]`)
	if err := os.WriteFile(tasksFile, data, 0o644); err != nil {
		t.Fatal(err)
	}

	tasks, err := loadTasks()
	if err != nil {
		t.Fatal(err)
	}

	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}

	if tasks[0].Status != StatusDone {
		t.Fatalf("expected legacy done task to load as done, got %q", tasks[0].Status)
	}
}

func captureOutput(t *testing.T, fn func()) string {
	t.Helper()

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	os.Stdout = w
	fn()

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	os.Stdout = oldStdout

	var buffer bytes.Buffer
	if _, err := io.Copy(&buffer, r); err != nil {
		t.Fatal(err)
	}

	return buffer.String()
}
