package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type commandResult struct {
	tasks   []Task
	changed bool
	message string
}

func main() {
	message, err := run(os.Args)
	if err != nil {
		if errors.Is(err, errHelp) {
			printUsage()
			return
		}

		fmt.Println(err)
		if errors.Is(err, errUsage) {
			printUsage()
		}
		os.Exit(1)
	}

	if message != "" {
		fmt.Println(message)
	}
}

func run(args []string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("please provide a command: %w", errUsage)
	}

	command := strings.ToLower(args[1])
	if command == "help" || command == "-h" || command == "--help" {
		return "", errHelp
	}

	tasks, err := loadTasks()
	if err != nil {
		return "", fmt.Errorf("failed to load tasks: %w", err)
	}

	result, err := executeCommand(command, args, tasks)
	if err != nil {
		return "", err
	}

	if result.changed {
		if err := saveTasks(result.tasks); err != nil {
			return "", fmt.Errorf("failed to save tasks: %w", err)
		}
	}

	return result.message, nil
}

func executeCommand(command string, args []string, tasks []Task) (commandResult, error) {
	switch command {
	case "add":
		return handleAddCommand(args, tasks)
	case "list":
		listTasks(tasks)
		return commandResult{tasks: tasks}, nil
	case "list-done":
		listTasksByStatus(tasks, StatusDone)
		return commandResult{tasks: tasks}, nil
	case "list-not-done":
		listNotDone(tasks)
		return commandResult{tasks: tasks}, nil
	case "list-in-progress":
		listTasksByStatus(tasks, StatusInProgress)
		return commandResult{tasks: tasks}, nil
	case "status":
		return handleStatusCommand(args, tasks)
	case "done":
		return handleStatusAliasCommand(args, tasks, StatusDone, "Task marked done.")
	case "in-progress":
		return handleStatusAliasCommand(args, tasks, StatusInProgress, "Task marked in progress.")
	case "not-done":
		return handleStatusAliasCommand(args, tasks, StatusNotDone, "Task marked not done.")
	case "help", "-h", "--help":
		return commandResult{}, errHelp
	case "delete":
		return handleDeleteCommand(args, tasks)
	case "rename":
		return handleRenameCommand(args, tasks)
	default:
		return commandResult{}, fmt.Errorf("unknown command %q: %w", command, errUsage)
	}
}

var errUsage = errors.New("usage")

var errHelp = errors.New("help")

func handleAddCommand(args []string, tasks []Task) (commandResult, error) {
	if len(args) < 3 {
		return commandResult{}, fmt.Errorf("please provide a task title: %w", errUsage)
	}

	title := strings.TrimSpace(strings.Join(args[2:], " "))
	if title == "" {
		return commandResult{}, fmt.Errorf("task title cannot be empty")
	}

	updatedTasks := addTask(tasks, title)
	return commandResult{tasks: updatedTasks, changed: true, message: "Task added."}, nil
}

func handleStatusCommand(args []string, tasks []Task) (commandResult, error) {
	if len(args) < 4 {
		return commandResult{}, fmt.Errorf("please provide a task ID and status: %w", errUsage)
	}

	id, err := parseIDArg(args)
	if err != nil {
		return commandResult{}, err
	}

	status, err := parseTaskStatus(args[3])
	if err != nil {
		return commandResult{}, err
	}

	updatedTasks, found := setTaskStatus(tasks, id, status)

	if !found {
		return commandResult{}, fmt.Errorf("task with ID %d was not found", id)
	}

	return commandResult{tasks: updatedTasks, changed: true, message: fmt.Sprintf("Task marked %s.", status)}, nil
}

func handleStatusAliasCommand(args []string, tasks []Task, status TaskStatus, message string) (commandResult, error) {
	id, err := parseIDArg(args)
	if err != nil {
		return commandResult{}, err
	}

	updatedTasks, found := setTaskStatus(tasks, id, status)
	if !found {
		return commandResult{}, fmt.Errorf("task with ID %d was not found", id)
	}

	return commandResult{tasks: updatedTasks, changed: true, message: message}, nil
}

func handleDeleteCommand(args []string, tasks []Task) (commandResult, error) {
	id, err := parseIDArg(args)
	if err != nil {
		return commandResult{}, err
	}

	updatedTasks, found := deleteTask(tasks, id)
	if !found {
		return commandResult{}, fmt.Errorf("task with ID %d was not found", id)
	}

	return commandResult{tasks: updatedTasks, changed: true, message: "Task deleted."}, nil
}

func handleRenameCommand(args []string, tasks []Task) (commandResult, error) {
	id, err := parseIDArg(args)
	if err != nil {
		return commandResult{}, err
	}

	if len(args) < 4 {
		return commandResult{}, fmt.Errorf("please provide a new title for the task: %w", errUsage)
	}

	newTitle := strings.TrimSpace(strings.Join(args[3:], " "))
	if newTitle == "" {
		return commandResult{}, fmt.Errorf("task title cannot be empty")
	}

	updatedTasks, found := renameTask(tasks, id, newTitle)
	if !found {
		return commandResult{}, fmt.Errorf("task with ID %d was not found", id)
	}

	return commandResult{tasks: updatedTasks, changed: true, message: "Task renamed."}, nil
}

func parseIDArg(args []string) (int, error) {
	if len(args) < 3 {
		return 0, fmt.Errorf("please provide a task ID: %w", errUsage)
	}

	id, err := strconv.Atoi(args[2])
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("task ID must be a positive integer")
	}

	return id, nil
}

func parseTaskStatus(value string) (TaskStatus, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "done":
		return StatusDone, nil
	case "in-progress", "in progress":
		return StatusInProgress, nil
	case "not-done", "not done":
		return StatusNotDone, nil
	default:
		return "", fmt.Errorf("invalid task status %q; use done, in-progress, or not-done", value)
	}
}

func printUsage() {
	fmt.Println("Task Tracker usage:")
	fmt.Println("  go run . add <title>")
	fmt.Println("  go run . list")
	fmt.Println("  go run . list-done")
	fmt.Println("  go run . list-not-done")
	fmt.Println("  go run . list-in-progress")
	fmt.Println("  go run . status <id> <done|in-progress|not-done>")
	fmt.Println("  go run . done <id>")
	fmt.Println("  go run . in-progress <id>")
	fmt.Println("  go run . not-done <id>")
	fmt.Println("  go run . delete <id>")
	fmt.Println("  go run . rename <id> <new title>")
	fmt.Println("  go run . help")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run . status 2 in-progress")
	fmt.Println("  go run . status 4 done")
}
