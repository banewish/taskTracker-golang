package main

import (
	"encoding/json"
	"errors"
	"os"
)

const tasksFile = "tasks.json"

func loadTasks() ([]Task, error) {
	data, err := os.ReadFile(tasksFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []Task{}, nil
		}
		return nil, err
	}

	if len(data) == 0 {
		return []Task{}, nil
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(tasksFile, data, 0o644)
}
