package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"


	"todo-cli/internal/storage"
	"todo-cli/internal/todo"
)

const tasksFilePath = "tasks.json"
const invalidArgsExitCode = 42

func main() {
	// Initialize storer and manager

storer := storage.NewJSONStorer(tasksFilePath)
manager, err := todo.NewManager(storer)
if err != nil {
	printErrorAndExit(fmt.Sprintf("failed to initialize manager: %v", err), 1)
}

	// Define subcommands

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addTaskDesc := addCmd.String("t", "", "Description of the task to add")

	doneCmd := flag.NewFlagSet("done", flag.ExitOnError)
	doneTaskID := doneCmd.Int("id", 0, "ID of the task to complete")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	// Parse main command
	if len(os.Args) < 2 {
		printErrorAndExit("invalid arguments", invalidArgsExitCode)
	}

	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		if *addTaskDesc == "" {
			printErrorAndExit("invalid arguments: task description cannot be empty", invalidArgsExitCode)
		}
		task, err := manager.Add(*addTaskDesc)
		if err != nil {
			printErrorAndExit(fmt.Sprintf("failed to add task: %v", err), 1)
		}
		printJSONOutput(map[string]interface{}{"id": task.ID, "status": "added"})

	case "list":
		listCmd.Parse(os.Args[2:])
		tasks := manager.List()
		// Filter for pending tasks only for 'list' command
	
pendingTasks := []*todo.Task{}
		for _, task := range tasks {
			if !task.Done {
				pendingTasks = append(pendingTasks, task)
			}
		}
		printJSONOutput(pendingTasks)

	case "done":
		doneCmd.Parse(os.Args[2:])
		if *doneTaskID == 0 {
			printErrorAndExit("invalid arguments: task ID cannot be 0", invalidArgsExitCode)
		}
		task, err := manager.Complete(*doneTaskID)
		if err != nil {
			printErrorAndExit(fmt.Sprintf("failed to complete task: %v", err), 1)
		}
		printJSONOutput(map[string]interface{}{"id": task.ID, "status": "completed"})

	default:
		printErrorAndExit("invalid arguments: unknown command", invalidArgsExitCode)
	}
}

func printJSONOutput(data interface{}) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		printErrorAndExit(fmt.Sprintf("failed to encode JSON output: %v", err), 1)
	}
}

func printErrorAndExit(msg string, code int) {
	output := map[string]string{"error": msg}
	jsonOutput, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "{\"error\": \"failed to marshal error: %v\"}\n", err)
	} else {
		fmt.Fprintln(os.Stderr, string(jsonOutput))
	}
	os.Exit(code)
}
