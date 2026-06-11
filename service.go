package main

import (
	"fmt"
	"time"
)

func nextTaskID(tasks []Task) int {
	maxID := 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	return maxID + 1
}

func taskTimestamp() time.Time {
	return time.Now().UTC().Truncate(time.Second)
}

func addTask(tasks []Task, title string) []Task {
	now := taskTimestamp()
	newTask := Task{
		ID:        nextTaskID(tasks),
		Title:     title,
		Status:    StatusNotDone,
		CreatedAt: now,
		UpdatedAt: now,
	}
	tasks = append(tasks, newTask)
	return tasks
}

func listTasks(tasks []Task) {
	printTaskHeader()
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}

	for _, task := range tasks {
		printTaskRow(task)
	}
}

func printTaskHeader() {
	fmt.Printf("%-4s %-12s %-28s %-16s %-16s\n", "ID", "STATUS", "TITLE", "CREATED", "UPDATED")
}

func formatTaskTime(value time.Time) string {
	if value.IsZero() {
		return "-"
	}

	return value.Format("2006-01-02 15:04")
}

func printTaskRow(task Task) {
	fmt.Printf("%-4d %-12s %-28s %-16s %-16s\n", task.ID, task.Status, task.Title, formatTaskTime(task.CreatedAt), formatTaskTime(task.UpdatedAt))
}

func setTaskStatus(tasks []Task, id int, status TaskStatus) ([]Task, bool) {
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Status = status
			tasks[i].UpdatedAt = taskTimestamp()
			return tasks, true
		}
	}
	return tasks, false
}

func markTaskDone(tasks []Task, id int) ([]Task, bool) {
	return setTaskStatus(tasks, id, StatusDone)
}

func markTaskInProgress(tasks []Task, id int) ([]Task, bool) {
	return setTaskStatus(tasks, id, StatusInProgress)
}

func markTaskNotDone(tasks []Task, id int) ([]Task, bool) {
	return setTaskStatus(tasks, id, StatusNotDone)
}

func deleteTask(tasks []Task, id int) ([]Task, bool) {
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return tasks, true
		}
	}
	return tasks, false
}

func renameTask(tasks []Task, id int, newTitle string) ([]Task, bool) {
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Title = newTitle
			tasks[i].UpdatedAt = taskTimestamp()
			return tasks, true
		}
	}
	return tasks, false
}

func listTasksByStatus(tasks []Task, status TaskStatus) {
	printTaskHeader()
	printed := false
	for _, task := range tasks {
		if task.Status == status {
			printTaskRow(task)
			printed = true
		}
	}

	if !printed {
		fmt.Println("No tasks found.")
	}
}

func listNotDone(tasks []Task) {
	listTasksByStatus(tasks, StatusNotDone)
}
