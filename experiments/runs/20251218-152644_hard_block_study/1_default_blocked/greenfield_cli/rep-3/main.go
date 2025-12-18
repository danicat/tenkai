package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var exit = os.Exit

const tasksFile = "tasks.json"

// Task represents a todo item
type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

// ErrorResponse represents a JSON error output
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse represents a JSON success output for add/done
type SuccessResponse struct {
	ID     int    `json:"id,omitempty"`
	Status string `json:"status"`
}

func printErrorAndExit(message string, exitCode int) {
	resp := ErrorResponse{Error: message}
	jsonOutput, _ := json.Marshal(resp)
	fmt.Println(string(jsonOutput))
	exit(exitCode)
}

func loadTasks() ([]Task, error) {
	if _, err := os.Stat(tasksFile); os.IsNotExist(err) {
		return []Task{}, nil
	}

	data, err := ioutil.ReadFile(tasksFile)
	if err != nil {
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

func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tasks: %w", err)
	}
	err = ioutil.WriteFile(tasksFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write tasks file: %w", err)
	}
	return nil
}

func getNextID(tasks []Task) int {
	if len(tasks) == 0 {
		return 1
	}
	maxID := 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	return maxID + 1
}

func main() {
	if len(os.Args) < 2 {
		printErrorAndExit("invalid arguments", 42)
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 3 {
			printErrorAndExit("invalid arguments", 42)
		}
		description := strings.Join(os.Args[2:], " ")
		if description == "" {
			printErrorAndExit("invalid arguments", 42)
		}

		tasks, err := loadTasks()
		if err != nil {
			printErrorAndExit(err.Error(), 1)
		}

		id := getNextID(tasks)
		newTask := Task{ID: id, Description: description, Done: false}
		tasks = append(tasks, newTask)

		err = saveTasks(tasks)
		if err != nil {
			printErrorAndExit(err.Error(), 1)
		}

		resp := SuccessResponse{ID: id, Status: "added"}
		jsonOutput, _ := json.Marshal(resp)
		fmt.Println(string(jsonOutput))

	case "list":
		tasks, err := loadTasks()
		if err != nil {
			printErrorAndExit(err.Error(), 1)
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
		if len(os.Args) < 3 {
			printErrorAndExit("invalid arguments", 42)
		}

		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			printErrorAndExit("invalid task ID", 42)
		}

		tasks, err := loadTasks()
		if err != nil {
			printErrorAndExit(err.Error(), 1)
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
			printErrorAndExit("task not found", 1)
		}

		err = saveTasks(tasks)
		if err != nil {
			printErrorAndExit(err.Error(), 1)
		}

		resp := SuccessResponse{ID: id, Status: "completed"}
		jsonOutput, _ := json.Marshal(resp)
		fmt.Println(string(jsonOutput))

	default:
		printErrorAndExit("unknown command", 1)
	}
}
