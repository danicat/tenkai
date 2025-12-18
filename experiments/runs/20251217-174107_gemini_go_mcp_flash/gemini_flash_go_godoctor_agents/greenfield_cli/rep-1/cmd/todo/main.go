package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"

	"todo-cli/internal/task"
)

// Error messages and exit codes
const (
	exitCodeInvalidArgs = 42
	exitCodeGeneralError  = 1
)

func printJSON(data interface{}) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		fatalError(fmt.Errorf("failed to encode JSON: %w", err))
	}
}

func fatalError(err error) {
	printJSON(map[string]string{"error": err.Error()})
	os.Exit(exitCodeGeneralError)
}

func invalidArgsError(msg string) {
	printJSON(map[string]string{"error": msg})
	os.Exit(exitCodeInvalidArgs)
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <command> [arguments]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nCommands:\n")
		fmt.Fprintf(os.Stderr, "  add <description>\tAdd a new task\n")
		fmt.Fprintf(os.Stderr, "  list\t\t\tList all pending tasks\n")
		fmt.Fprintf(os.Stderr, "  done <id>\t\tMark a task as done\n")
		fmt.Fprintf(os.Stderr, "\n")
	}

	flag.Parse()

	if flag.NArg() < 1 {
		invalidArgsError("no command provided")
	}

	command := flag.Arg(0)

	switch command {
	case "add":
		if flag.NArg() < 2 {
			invalidArgsError("missing task description")
		}
		description := flag.Arg(1)
		newTask, err := task.AddTask(description)
		if err != nil {
			fatalError(fmt.Errorf("failed to add task: %w", err))
		}
		printJSON(map[string]interface{}{"id": newTask.ID, "status": "added"})

	case "list":
		tasks, err := task.ListPendingTasks()
		if err != nil {
			fatalError(fmt.Errorf("failed to list tasks: %w", err))
		}
		// Convert tasks to the required JSON format for listing
		type TaskListOutput struct {
			ID          int    `json:"id"`
			Description string `json:"description"`
			Done        bool   `json:"done"`
		}
		var output []TaskListOutput
		for _, t := range tasks {
			output = append(output, TaskListOutput{ID: t.ID, Description: t.Description, Done: t.Done})
		}
		printJSON(output)

	case "done":
		if flag.NArg() < 2 {
			invalidArgsError("missing task ID")
		}
		idStr := flag.Arg(1)
		id, err := strconv.Atoi(idStr)
		if err != nil {
			invalidArgsError(fmt.Sprintf("invalid task ID: %s", idStr))
		}

		completedTask, err := task.CompleteTask(id)
		if err != nil {
			fatalError(fmt.Errorf("failed to complete task: %w", err))
		}
		printJSON(map[string]interface{}{"id": completedTask.ID, "status": "completed"})

	default:
		invalidArgsError(fmt.Sprintf("unknown command: %s", command))
	}
}
