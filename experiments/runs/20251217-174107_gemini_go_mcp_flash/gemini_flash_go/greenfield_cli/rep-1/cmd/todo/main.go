package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"todo/internal/task"
)

const (
	errorInvalidArgs = 42
	errorGeneral     = 1
)

// outputJSON prints the given data as JSON to stdout.
func outputJSON(data interface{}) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(errorGeneral)
	}
}

// printErrorAndExit prints an error message as JSON to stdout and exits with the given code.
func printErrorAndExit(err error, code int) {
	outErr := map[string]string{"error": err.Error()}
	outputJSON(outErr)
	os.Exit(code)
}

func main() {
	if len(os.Args) < 2 {
		printErrorAndExit(fmt.Errorf("invalid arguments"), errorInvalidArgs)
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 3 {
			printErrorAndExit(fmt.Errorf("missing task description"), errorInvalidArgs)
		}
		description := strings.Join(os.Args[2:], " ")
		task, err := task.AddTask(description)
		if err != nil {
			printErrorAndExit(err, errorGeneral)
		}
		outputJSON(map[string]interface{}{"id": task.ID, "status": "added"})

	case "list":
		tasks := task.ListTasks()
		// Filter for undone tasks as per requirements (list all pending tasks)
	
pendingTasks := []task.Task{}
		for _, t := range tasks {
			if !t.Done {
				pendingTasks = append(pendingTasks, t)
			}
		}
		outputJSON(pendingTasks)

	case "done":
		if len(os.Args) < 3 {
			printErrorAndExit(fmt.Errorf("missing task ID"), errorInvalidArgs)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			printErrorAndExit(fmt.Errorf("invalid task ID"), errorInvalidArgs)
		}
		task, err := task.CompleteTask(id)
		if err != nil {
			printErrorAndExit(err, errorGeneral)
		}
		outputJSON(map[string]interface{}{"id": task.ID, "status": "completed"})

	default:
		printErrorAndExit(fmt.Errorf("unknown command: %s", command), errorInvalidArgs)
	}
}
