package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"

	"experiment/internal/task"
)

const tasksFilePath = "tasks.json"
const invalidArgsExitCode = 42

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	tm, err := task.NewTasksManager(tasksFilePath)
	if err != nil {
		return exitWithError(fmt.Errorf("failed to initialize task manager: %w", err))
	}

	if len(args) < 2 {
		return exitWithInvalidArgs("no command provided")
	}

	command := args[1]
	switch command {
	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		addDesc := addCmd.String("desc", "", "Description of the task")
		addCmd.Parse(args[2:])

		if *addDesc == "" {
			return exitWithInvalidArgs("description cannot be empty")
		}
		newTask, err := tm.Add(*addDesc)
		if err != nil {
			return exitWithError(fmt.Errorf("failed to add task: %w", err))
		}
		printJSON(map[string]interface{}{"id": newTask.ID, "status": "added"})
	case "list":
		tasks := tm.List()
		var output []map[string]interface{}
		for _, t := range tasks {
			output = append(output, map[string]interface{}{
				"id":          t.ID,
				"description": t.Description,
				"done":        t.Done,
			})
		}
		printJSON(output)
	case "done":
		doneCmd := flag.NewFlagSet("done", flag.ExitOnError)
		doneIDStr := doneCmd.String("id", "", "ID of the task to complete")
		doneCmd.Parse(args[2:])

		if *doneIDStr == "" {
			return exitWithInvalidArgs("task ID cannot be empty")
		}
		taskID, err := strconv.Atoi(*doneIDStr)
		if err != nil {
			return exitWithInvalidArgs("invalid task ID")
		}
		completedTask, err := tm.Complete(taskID)
		if err != nil {
			return exitWithError(fmt.Errorf("failed to complete task: %w", err))
		}
		printJSON(map[string]interface{}{"id": completedTask.ID, "status": "completed"})
	default:
		return exitWithInvalidArgs(fmt.Sprintf("unknown command: %s", command))
	}
	return nil
}

func printJSON(data interface{}) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "") // No indentation for strict JSON output
	if err := encoder.Encode(data); err != nil {
		// Fallback to plain error if JSON encoding fails for some reason
		fmt.Fprintf(os.Stderr, "{\"error\": \"failed to encode JSON output: %v\"}\n", err)
		os.Exit(1)
	}
}

func exitWithError(err error) error {
	printJSON(map[string]string{"error": err.Error()})
	os.Exit(1)
	return err // This line is unreachable but satisfies return type
}

func exitWithInvalidArgs(msg string) error {
	printJSON(map[string]string{"error": msg})
	os.Exit(invalidArgsExitCode)
	return fmt.Errorf("%s", msg) // This line is unreachable but satisfies return type
}
