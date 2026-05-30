package main

import "encoding/json"

type TaskStatus string

const (
	StatusNotDone    TaskStatus = "not done"
	StatusInProgress TaskStatus = "in progress"
	StatusDone       TaskStatus = "done"
)

type Task struct {
	ID     int        `json:"id"`
	Title  string     `json:"title"`
	Status TaskStatus `json:"status"`
}

func (task *Task) UnmarshalJSON(data []byte) error {
	type rawTask struct {
		ID     int         `json:"id"`
		Title  string      `json:"title"`
		Status TaskStatus  `json:"status"`
		Done   interface{} `json:"done"`
	}

	var raw rawTask
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	task.ID = raw.ID
	task.Title = raw.Title

	switch raw.Status {
	case StatusNotDone, StatusInProgress, StatusDone:
		task.Status = raw.Status
		return nil
	}

	if done, ok := raw.Done.(bool); ok {
		if done {
			task.Status = StatusDone
		} else {
			task.Status = StatusNotDone
		}
		return nil
	}

	task.Status = StatusNotDone
	return nil
}
