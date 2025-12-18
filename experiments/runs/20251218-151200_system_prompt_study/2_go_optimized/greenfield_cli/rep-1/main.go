package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

var (
	tasksFile = "tasks.json"
	osExit    = os.Exit // Allows mocking os.Exit in tests
)

// Task represents a todo item
type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

// handleError prints an error message in JSON format and exits with the given code.
func handleError(err error, exitCode int) {
	outputError(err.Error(), exitCode)
}

// outputError prints an error message in JSON format and exits with the given code.
func outputError(msg string, exitCode int) {
	jsonOutput, _ := json.Marshal(map[string]string{"error": msg})
	fmt.Println(string(jsonOutput))
	osExit(exitCode)
}

// outputJSON prints the given data in JSON format and exits with code 0.
func outputJSON(data interface{}) {
	jsonOutput, err := json.Marshal(data)
	if err != nil {
		handleError(fmt.Errorf("failed to marshal JSON: %w", err), 1)
	}
	fmt.Println(string(jsonOutput))
	osExit(0)
}

func main() {
	if len(os.Args) < 2 {
		outputError("no command provided", 42)
	}

	command := os.Args[1]

	switch command {
	case "add":
		// Expects `todo add "description"`
		if len(os.Args) < 3 {
			outputError("missing task description", 42)
		}
		description := os.Args[2]
		addTask(description)
	case "list":
		listTasks()
	case "done":
		// Expects `todo done <ID>`
		if len(os.Args) < 3 {
			outputError("missing task ID", 42)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			outputError("invalid task ID", 42)
		}
		completeTask(id)
	default:
		outputError("unknown command", 42)
	}
}

// loadTasks loads tasks from tasks.json
func loadTasks(filePath string) ([]Task, error) {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return []Task{}, nil // Return empty list if file doesn't exist
	}
	if err != nil {
		return nil, fmt.Errorf("error checking tasks file: %w", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading tasks file: %w", err)
	}

	var tasks []Task
	if len(data) == 0 {
		return []Task{}, nil // Return empty list if file is empty
	}

	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling tasks: %w", err)
	}

	return tasks, nil
}

// saveTasks saves tasks to tasks.json
func saveTasks(filePath string, tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling tasks: %w", err)
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing tasks file: %w", err)
	}

	return nil
}

// addTask adds a new task to the list
func addTask(description string) {
	tasks, err := loadTasks(tasksFile)
	if err != nil {
		handleError(err, 1)
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

	err = saveTasks(tasksFile, tasks)
	if err != nil {
		handleError(err, 1)
	}

	outputJSON(map[string]interface{}{"id": newTask.ID, "status": "added"})
}

// listTasks lists all pending tasks
func listTasks() {
	tasks, err := loadTasks(tasksFile)
	if err != nil {
		handleError(err, 1)
	}

	var pendingTasks []Task
	for _, task := range tasks {
		if !task.Done {
			pendingTasks = append(pendingTasks, task)
		}
	}

	outputJSON(pendingTasks)
}

// completeTask marks a task as done by its ID
func completeTask(id int) {
	tasks, err := loadTasks(tasksFile)
	if err != nil {
		handleError(err, 1)
	}

	found := false
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Done = true
			found = true
			break
		}
	}

	if !found {
		outputError(fmt.Sprintf("task with ID %d not found", id), 1)
	}

	err = saveTasks(tasksFile, tasks)
	if err != nil {
		handleError(err, 1)
	}

	outputJSON(map[string]interface{}{"id": id, "status": "completed"})
}
