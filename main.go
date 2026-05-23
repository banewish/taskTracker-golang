package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("Failed to load tasks:", err)
		os.Exit(1)
	}

	fmt.Println("Enter title:")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		fmt.Println("No title provided")
		return
	}

	title := scanner.Text()
	tasks = addTask(tasks, title)

	if err := saveTasks(tasks); err != nil {
		fmt.Println("Failed to save tasks:", err)
		os.Exit(1)
	}

	listTasks(tasks)
}

func listTasks(tasks []Task) {
	for _, task := range tasks {
		fmt.Println(task.ID, task.Title, task.Done)
	}
}
