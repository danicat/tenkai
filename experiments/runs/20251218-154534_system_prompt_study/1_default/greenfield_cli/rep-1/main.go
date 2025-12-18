package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Error response structure.
type ErrorResponse struct {
	Error string `json:"error"`
}

// printError prints an error message in JSON format and exits.
var exitFunc = os.Exit // Allows mocking os.Exit in tests

// printError prints an error message in JSON format and exits.
func printError(msg string, code int) {
	resp := ErrorResponse{Error: msg}
	jsonOutput, _ := json.Marshal(resp)
	fmt.Println(string(jsonOutput))
	exitFunc(code)
}

func main() {
	// Define subcommands
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addDescription := addCmd.String("", "", "Task description") // Positional argument workaround

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	doneCmd := flag.NewFlagSet("done", flag.ExitOnError)
	doneID := doneCmd.Int("", 0, "Task ID") // Positional argument workaround

	// Check which subcommand is called
	if len(os.Args) < 2 {
		printError("invalid arguments", 42)
	}

	switch os.Args[1] {
	case "add":
		// Parse add command's arguments
		if err := addCmd.Parse(os.Args[2:]); err != nil {
			printError(fmt.Sprintf("failed to parse add command arguments: %v", err), 1)
		}
		// Custom parsing for positional argument
		if addCmd.NArg() == 0 {
			printError("invalid arguments", 42)
		}
		description := strings.Join(addCmd.Args(), " ")
		if description == "" {
			printError("invalid arguments", 42)
		}

		tasks, err := loadTasks()
		if err != nil {
			printError(fmt.Sprintf("failed to load tasks: %v", err), 1)
		}

		nextID := 1
		if len(tasks) > 0 {
			nextID = tasks[len(tasks)-1].ID + 1
		}

		newTask := Task{
			ID:          nextID,
			Description: description,
			Done:        false,
		}
		tasks = append(tasks, newTask)

		if err := saveTasks(tasks); err != nil {
			printError(fmt.Sprintf("failed to save tasks: %v", err), 1)
		}

		type AddResponse struct {
			ID     int    `json:"id"`
			Status string `json:"status"`
		}
		resp := AddResponse{ID: newTask.ID, Status: "added"}
		jsonOutput, _ := json.Marshal(resp)
		fmt.Println(string(jsonOutput))

	case "list":
		listCmd.Parse(os.Args[2:])
		tasks, err := loadTasks()
		if err != nil {
			printError(fmt.Sprintf("failed to load tasks: %v", err), 1)
		}

		var pendingTasks []Task
		for _, task := range tasks {
			if !task.Done {
				pendingTasks = append(pendingTasks, task)
			}
		}

		jsonOutput, _ := json.Marshal(pendingTasks)
		fmt.Println(string(jsonOutput))

	case "done":
		// Parse done command's arguments
		if err := doneCmd.Parse(os.Args[2:]); err != nil {
			printError(fmt.Sprintf("failed to parse done command arguments: %v", err), 1)
		}
		// Custom parsing for positional argument
		if doneCmd.NArg() == 0 {
			printError("invalid arguments", 42)
		}
		taskIDStr := doneCmd.Args()[0]
		taskID, err := strconv.Atoi(taskIDStr)
		if err != nil || taskID <= 0 {
			printError("invalid arguments", 42)
		}

		tasks, err := loadTasks()
		if err != nil {
			printError(fmt.Sprintf("failed to load tasks: %v", err), 1)
		}

		found := false
		for i := range tasks {
			if tasks[i].ID == taskID {
				tasks[i].Done = true
				found = true
				break
			}
		}

		if !found {
			printError(fmt.Sprintf("task with ID %d not found", taskID), 1)
		}

		if err := saveTasks(tasks); err != nil {
			printError(fmt.Sprintf("failed to save tasks: %v", err), 1)
		}

		type DoneResponse struct {
			ID     int    `json:"id"`
			Status string `json:"status"`
		}
		resp := DoneResponse{ID: taskID, Status: "completed"}
		jsonOutput, _ := json.Marshal(resp)
		fmt.Println(string(jsonOutput))

	default:
		printError("invalid arguments", 42)
	}
}
