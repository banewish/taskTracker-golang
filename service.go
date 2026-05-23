package main

func addTask(tasks []Task, title string) []Task {
	newTask := Task{
		ID:    len(tasks) + 1,
		Title: title,
		Done:  false,
	}
	return append(tasks, newTask)
}
