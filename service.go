package main

import "fmt"

func renumberTaskIDs(tasks []Task) []Task {
	for i := range tasks {
		tasks[i].ID = i + 1
	}
	return tasks
}

func addTask(tasks []Task, title string) []Task {
	newTask := Task{
		ID:    len(tasks) + 1,
		Title: title,
		Done:  false,
	}
	tasks = append(tasks, newTask)
	return renumberTaskIDs(tasks)
}

func listTasks(tasks []Task) {
	for _, task := range tasks {
		fmt.Println(task.ID, task.Title, task.Done)
	}
}

func markTaskDone(tasks []Task, id int) ([]Task, bool) {
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Done = true
			return tasks, true
		}
	}
	return tasks, false
}

func deleteTask(tasks []Task, id int) ([]Task, bool) {
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return renumberTaskIDs(tasks), true
		}
	}
	return tasks, false
}

func renameTask(tasks []Task, id int, newTitle string) ([]Task, bool) {
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Title = newTitle
			return tasks, true
		}
	}
	return tasks, false
}
