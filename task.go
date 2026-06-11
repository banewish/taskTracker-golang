package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type TaskStatus string

const (
	StatusNotDone    TaskStatus = "not done"
	StatusInProgress TaskStatus = "in progress"
	StatusDone       TaskStatus = "done"
)

type Task struct {
	ID        int        `json:"id"`
	Title     string     `json:"title"`
	Status    TaskStatus `json:"status"`
	CreatedAt time.Time  `json:"createdAt,omitempty"`
	UpdatedAt time.Time  `json:"updatedAt,omitempty"`
}

func (task *Task) UnmarshalJSON(data []byte) error {
	type rawTask struct {
		ID        int         `json:"id"`
		Title     string      `json:"title"`
		Status    *TaskStatus `json:"status"`
		Done      *bool       `json:"done"`
		CreatedAt *time.Time  `json:"createdAt"`
		UpdatedAt *time.Time  `json:"updatedAt"`
	}

	var raw rawTask
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	task.ID = raw.ID
	task.Title = raw.Title
	if raw.CreatedAt != nil {
		task.CreatedAt = *raw.CreatedAt
	}
	if raw.UpdatedAt != nil {
		task.UpdatedAt = *raw.UpdatedAt
	}

	if raw.Status != nil {
		switch *raw.Status {
		case StatusNotDone, StatusInProgress, StatusDone:
			task.Status = *raw.Status
			return nil
		default:
			return fmt.Errorf("invalid task status %q", *raw.Status)
		}
	}

	if raw.Done != nil {
		if *raw.Done {
			task.Status = StatusDone
		} else {
			task.Status = StatusNotDone
		}
		return nil
	}

	return fmt.Errorf("missing task status")
}
