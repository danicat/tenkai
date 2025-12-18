package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"

	"todo-cli/internal/todo"
)

// outputJSON prints the given data as JSON to stdout.
func outputJSON(data interface{}) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1) // Exit with a generic error code if JSON encoding fails
	}
}

// handleExitError prints an error message as JSON and exits with the specified code.
func handleExitError(err error, exitCode int) {
	errorMsg := map[string]string{"error": err.Error()}
	outputJSON(errorMsg)
	os.Exit(exitCode)
}

func main() {
	// Subcommands
	// addCmd is not used; description is taken from os.Args[2]

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listAll := listCmd.Bool("all", false, "List all tasks, including completed ones")

	// doneCmd is not used; ID is taken from os.Args[2]

	if len(os.Args) < 2 {
		handleExitError(fmt.Errorf("invalid arguments"), 42)
	}

	switch os.Args[1] {
	case "add":
		// Custom parsing for positional argument for "add"
		if len(os.Args) < 3 {
			handleExitError(fmt.Errorf("invalid arguments: missing task description"), 42)
		}
		description := os.Args[2]
		if description == "" {
			handleExitError(fmt.Errorf("invalid arguments: task description cannot be empty"), 42)
		}
		task, err := todo.AddTask(description)
		if err != nil {
			handleExitError(fmt.Errorf("failed to add task: %w", err), 1)
		}
		outputJSON(map[string]interface{}{"id": task.ID, "status": "added"})

	case "list":
		listCmd.Parse(os.Args[2:])
		tasks, err := todo.ListTasks(*listAll)
		if err != nil {
			handleExitError(fmt.Errorf("failed to list tasks: %w", err), 1)
		}
		outputJSON(tasks)

	case "done":
		// Custom parsing for positional argument for "done"
		if len(os.Args) < 3 {
			handleExitError(fmt.Errorf("invalid arguments: missing task ID"), 42)
		}
		idStr := os.Args[2]
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			handleExitError(fmt.Errorf("invalid arguments: task ID must be a positive integer"), 42)
		}
	
		task, err := todo.CompleteTask(id)
		if err != nil {
			handleExitError(fmt.Errorf("failed to complete task: %w", err), 1)
		}
		outputJSON(map[string]interface{}{"id": task.ID, "status": "completed"})

	default:
		handleExitError(fmt.Errorf("invalid command: %s", os.Args[1]), 42)
	}
}
