package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	tasksFile = "tasks.json"
)

// osExit is a package-level variable that can be swapped for testing
var osExit = os.Exit

// Task represents a single todo item
type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

// loadTasks loads tasks from tasks.json
func loadTasks() ([]Task, error) {
	data, err := ioutil.ReadFile(tasksFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil // Return empty list if file doesn't exist
		}
		return nil, fmt.Errorf("failed to read tasks file: %w", err)
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tasks: %w", err)
	}
	return tasks, nil
}

// saveTasks saves tasks to tasks.json
func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tasks: %w", err)
	}

	if err := ioutil.WriteFile(tasksFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write tasks file: %w", err)
	}
	return nil
}

// generateErrorOutput generates an error message in JSON format
func generateErrorOutput(msg string) string {
	output := map[string]string{"error": msg}
	jsonOutput, _ := json.Marshal(output)
	return string(jsonOutput)
}

// runCommand executes the appropriate command and returns output and exit code
func runCommand(args []string) (string, int) {
	if len(args) < 2 {
		return generateErrorOutput("invalid arguments"), 42
	}

	command := args[1]

	switch command {
	case "add":
		return handleAdd(args)
	case "list":
		return handleList()
	case "done":
		return handleDone(args)
	default:
		return generateErrorOutput("unknown command"), 42
	}
}

func main() {
	output, code := runCommand(os.Args)
	fmt.Println(output)
	if code != 0 {
		osExit(code)
	}
}

func handleAdd(args []string) (string, int) {
	if len(args) < 3 || strings.TrimSpace(args[2]) == "" {
		return generateErrorOutput("invalid arguments"), 42
	}
	description := args[2]

	tasks, err := loadTasks()
	if err != nil {
		return generateErrorOutput(err.Error()), 1
	}

	nextID := 1
	if len(tasks) > 0 {
		for _, task := range tasks {
			if task.ID >= nextID {
				nextID = task.ID + 1
			}
		}
	}

	newTask := Task{
		ID:          nextID,
		Description: description,
		Done:        false,
	}
	tasks = append(tasks, newTask)

	if err := saveTasks(tasks); err != nil {
		return generateErrorOutput(err.Error()), 1
	}

	result := map[string]interface{}{"id": newTask.ID, "status": "added"}
	jsonOutput, _ := json.Marshal(result)
	return string(jsonOutput), 0
}

func handleList() (string, int) {
	tasks, err := loadTasks()
	if err != nil {
		return generateErrorOutput(err.Error()), 1
	}

	var pendingTasks []Task
	for _, task := range tasks {
		if !task.Done {
			pendingTasks = append(pendingTasks, task)
		}
	}

	jsonOutput, _ := json.Marshal(pendingTasks)
	return string(jsonOutput), 0
}

func handleDone(args []string) (string, int) {
	if len(args) < 3 || strings.TrimSpace(args[2]) == "" {
		return generateErrorOutput("invalid arguments"), 42
	}

	idStr := args[2]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return generateErrorOutput("invalid task ID"), 42
	}

	tasks, err := loadTasks()
	if err != nil {
		return generateErrorOutput(err.Error()), 1
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
		return generateErrorOutput(fmt.Sprintf("task with ID %d not found", id)), 42
	}

	if err := saveTasks(tasks); err != nil {
		return generateErrorOutput(err.Error()), 1
	}

	result := map[string]interface{}{"id": id, "status": "completed"}
	jsonOutput, _ := json.Marshal(result)
	return string(jsonOutput), 0
}
