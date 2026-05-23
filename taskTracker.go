package main

import (
	"bufio"
	"fmt"
	"os"
)

type Task struct { // defining a struct to represent a task
	ID    int
	Title string
	Done  bool
}

func main() {
	tasks := []Task{}
	fmt.Println("Enter title: ")
	var title string // creating scanner for user input
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		title = scanner.Text() // reading user input and storing it in title variable
	}

	tasks = addTask(tasks, title)
	listTasks(tasks)
}

func addTask(tasks []Task, title string) []Task {

	newTask := Task{
		ID:    len(tasks) + 1,
		Title: title,
		Done:  false,
	}
	return append(tasks, newTask) // adding the new task to the list of tasks and returning the updated list
}

func listTasks(tasks []Task) { // function to list all tasks
	for _, task := range tasks {
		fmt.Println(task.ID, task.Title, task.Done)
	}
}
