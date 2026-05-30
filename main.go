package main

import (
	"fmt"

	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := strings.ToLower(os.Args[1])

	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("Failed to load tasks:", err)
		os.Exit(1)
	}

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a task title.")
			printUsage()
			os.Exit(1)
		}

		title := strings.TrimSpace(strings.Join(os.Args[2:], " "))
		if title == "" {
			fmt.Println("Task title cannot be empty.")
			os.Exit(1)
		}

		tasks = addTask(tasks, title)
		if err := saveTasks(tasks); err != nil {
			fmt.Println("Failed to save tasks:", err)
			os.Exit(1)
		}

		fmt.Println("Task added.")

	case "list":
		listTasks(tasks)

	case "list-done":
		listTasksByStatus(tasks, StatusDone)

	case "list-not-done":
		listNotDone(tasks)

	case "list-in-progress":
		listTasksByStatus(tasks, StatusInProgress)

	case "done":
		id, ok := parseIDArg(os.Args)
		if !ok {
			os.Exit(1)
		}

		var found bool
		tasks, found = markTaskDone(tasks, id)
		if !found {
			fmt.Printf("Task with ID %d was not found.\n", id)
			os.Exit(1)
		}

		if err := saveTasks(tasks); err != nil {
			fmt.Println("Failed to save tasks:", err)
			os.Exit(1)
		}

		fmt.Println("Task marked done.")

	case "in-progress":
		id, ok := parseIDArg(os.Args)
		if !ok {
			os.Exit(1)
		}

		var found bool
		tasks, found = markTaskInProgress(tasks, id)
		if !found {
			fmt.Printf("Task with ID %d was not found.\n", id)
			os.Exit(1)
		}

		if err := saveTasks(tasks); err != nil {
			fmt.Println("Failed to save tasks:", err)
			os.Exit(1)
		}

		fmt.Println("Task marked in progress.")

	case "not-done":
		id, ok := parseIDArg(os.Args)
		if !ok {
			os.Exit(1)
		}

		var found bool
		tasks, found = markTaskNotDone(tasks, id)
		if !found {
			fmt.Printf("Task with ID %d was not found.\n", id)
			os.Exit(1)
		}

		if err := saveTasks(tasks); err != nil {
			fmt.Println("Failed to save tasks:", err)
			os.Exit(1)
		}

		fmt.Println("Task marked not done.")

	case "delete":
		id, ok := parseIDArg(os.Args)
		if !ok {
			os.Exit(1)
		}

		var found bool
		tasks, found = deleteTask(tasks, id)
		if !found {
			fmt.Printf("Task with ID %d was not found.\n", id)
			os.Exit(1)
		}

		if err := saveTasks(tasks); err != nil {
			fmt.Println("Failed to save tasks:", err)
			os.Exit(1)
		}

		fmt.Println("Task deleted.")

	case "rename":
		id, ok := parseIDArg(os.Args)
		if !ok {
			os.Exit(1)
		}

		if len(os.Args) < 4 {
			fmt.Println("Please provide a new title for the task.")
			printUsage()
			os.Exit(1)
		}

		newTitle := strings.TrimSpace(strings.Join(os.Args[3:], " "))
		if newTitle == "" {
			fmt.Println("Task title cannot be empty.")
			os.Exit(1)
		}

		var found bool
		tasks, found = renameTask(tasks, id, newTitle)
		if !found {
			fmt.Printf("Task with ID %d was not found.\n", id)
			os.Exit(1)
		}

		if err := saveTasks(tasks); err != nil {
			fmt.Println("Failed to save tasks:", err)
			os.Exit(1)
		}

		fmt.Println("Task renamed.")

	default:
		fmt.Println("Unknown command:", command)
		printUsage()
		os.Exit(1)
	}

}

func parseIDArg(args []string) (int, bool) {
	if len(args) < 3 {
		fmt.Println("Please provide a task ID.")
		printUsage()
		return 0, false
	}

	id, err := strconv.Atoi(args[2])
	if err != nil || id <= 0 {
		fmt.Println("Task ID must be a positive integer.")
		return 0, false
	}

	return id, true
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  go run . add <title>")
	fmt.Println("  go run . list")
	fmt.Println("  go run . list-done")
	fmt.Println("  go run . list-not-done")
	fmt.Println("  go run . list-in-progress")
	fmt.Println("  go run . done <id>")
	fmt.Println("  go run . in-progress <id>")
	fmt.Println("  go run . not-done <id>")
	fmt.Println("  go run . delete <id>")
	fmt.Println("  go run . rename <id> <new title>")
}
