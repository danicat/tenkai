package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const tasksFile = "tasks.json"

// Task represents a todo item
type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

// Exit codes
const (
	ExitCodeOK               = 0
	ExitCodeError            = 1
	ExitCodeInvalidArguments = 42
)

// exit is a variable that can be overridden in tests to prevent os.Exit from terminating the test runner.
var exit = os.Exit

// handleError prints a JSON error message and exits with the specified code.
func handleError(err error, code int) {
	outputJSON(map[string]string{"error": err.Error()}, code)
}

// outputJSON prints a JSON object to stdout and exits with the specified code.
func outputJSON(data interface{}, code int) {
	jsonOutput, err := json.Marshal(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "{\"error\": \"failed to marshal JSON output: %s\"}\n", err)
		exit(ExitCodeError)
	}
	fmt.Println(string(jsonOutput))
	exit(code)
}

func main() {
	if len(os.Args) < 2 {
		handleError(fmt.Errorf("invalid arguments"), ExitCodeInvalidArguments)
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 3 {
			handleError(fmt.Errorf("missing task description"), ExitCodeInvalidArguments)
		}
		description := strings.Join(os.Args[2:], " ")

		tasks, err := loadTasks()
		if err != nil {
			handleError(err, ExitCodeError)
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

		err = saveTasks(tasks)
		if err != nil {
			handleError(err, ExitCodeError)
		}

		outputJSON(map[string]interface{}{"id": newTask.ID, "status": "added"}, ExitCodeOK)
	case "list":
		tasks, err := loadTasks()
		if err != nil {
			handleError(err, ExitCodeError)
		}

		var pendingTasks []Task
		for _, task := range tasks {
			if !task.Done {
				pendingTasks = append(pendingTasks, task)
			}
		}

		outputJSON(pendingTasks, ExitCodeOK)
	case "done":
		if len(os.Args) < 3 {
			handleError(fmt.Errorf("missing task ID"), ExitCodeInvalidArguments)
		}
		taskID, err := strconv.Atoi(os.Args[2])
		if err != nil {
			handleError(fmt.Errorf("invalid task ID: %w", err), ExitCodeInvalidArguments)
		}

		tasks, err := loadTasks()
		if err != nil {
			handleError(err, ExitCodeError)
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
			handleError(fmt.Errorf("task with ID %d not found", taskID), ExitCodeError)
		}

		err = saveTasks(tasks)
		if err != nil {
			handleError(err, ExitCodeError)
		}

		outputJSON(map[string]interface{}{"id": taskID, "status": "completed"}, ExitCodeOK)
	default:
		handleError(fmt.Errorf("unknown command: %s", command), ExitCodeInvalidArguments)
	}
}

// loadTasks reads tasks from tasks.json
func loadTasks() ([]Task, error) {
	data, err := os.ReadFile(tasksFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil // Return empty list if file doesn't exist
		}
		return nil, fmt.Errorf("failed to read tasks file: %w", err)
	}

	var tasks []Task
	if len(data) == 0 {
		return []Task{}, nil
	}

	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tasks: %w", err)
	}
	return tasks, nil
}

// saveTasks writes tasks to tasks.json
func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tasks: %w", err)
	}

	err = os.WriteFile(tasksFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write tasks file: %w", err)
	}
	return nil
}
