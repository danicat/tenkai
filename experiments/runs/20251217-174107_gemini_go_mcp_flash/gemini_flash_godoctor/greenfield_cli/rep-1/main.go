package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
)

const tasksFile = "tasks.json"

func main() {

	// Initialize TodoManager

	t := NewTodoManager(tasksFile)



	// Subcommands

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)

	doneCmd := flag.NewFlagSet("done", flag.ExitOnError)

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)



	if len(os.Args) < 2 {

		printError("invalid arguments", 42)

	}



	switch os.Args[1] {

	case "add":

		addCmd.Parse(os.Args[2:])

		if len(addCmd.Args()) < 1 || addCmd.Args()[0] == "" {

			printError("task description cannot be empty", 42)

		}

		description := addCmd.Args()[0]

		task := t.AddTask(description)

		printJSON(map[string]interface{}{"id": task.ID, "status": "added"})



	case "list":

		listCmd.Parse(os.Args[2:])

		tasks := t.GetTasks()

		// Filter for pending tasks only

		pendingTasks := []Task{}

		for _, task := range tasks {

			if !task.Done {

				pendingTasks = append(pendingTasks, task)

			}

		}

		printJSON(pendingTasks)



	case "done":

		doneCmd.Parse(os.Args[2:])

		if len(doneCmd.Args()) < 1 || doneCmd.Args()[0] == "" {

			printError("task ID cannot be empty", 42)

		}

		id, err := strconv.Atoi(doneCmd.Args()[0])

		if err != nil {

			printError("invalid task ID", 42)

		}

		task, err := t.CompleteTask(id)

		if err != nil {

			printError(err.Error(), 1)

		}

		printJSON(map[string]interface{}{"id": task.ID, "status": "completed"})



	default:

		printError("unknown command", 42)

	}

}

func printJSON(data interface{}) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ") // Pretty print for readability
	if err := encoder.Encode(data); err != nil {
		// Fallback to plain error if JSON encoding fails
		fmt.Fprintf(os.Stderr, "{\"error\": \"failed to encode JSON: %v\"}\n", err)
		os.Exit(1)
	}
}

func printError(message string, exitCode int) {
	errorOutput := map[string]string{"error": message}
	encoder := json.NewEncoder(os.Stdout)
	if err := encoder.Encode(errorOutput); err != nil {
		// Fallback to plain error if JSON encoding fails
		fmt.Fprintf(os.Stderr, "{\"error\": \"failed to encode JSON error: %v\"}\n", err)
	}
	os.Exit(exitCode)
}