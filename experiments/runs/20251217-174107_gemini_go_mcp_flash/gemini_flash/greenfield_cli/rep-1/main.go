package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

var tasksFile = "tasks.json"

// Task represents a single todo item
type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

// readTasks reads tasks from the tasks.json file
func readTasks(filePath string) ([]Task, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil // Return empty slice if file doesn't exist
		}
		return nil, fmt.Errorf("error reading tasks file %s: %w", filePath, err)
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("error unmarshalling tasks from %s: %w", filePath, err)
	}
	return tasks, nil
}

// saveTasks writes tasks to the tasks.json file
func saveTasks(filePath string, tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling tasks: %w", err)
	}
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("error writing tasks file %s: %w", filePath, err)
	}
	return nil
}

func printErrorAndExit(err string, code int, exitFunc func(int)) {
	output := map[string]string{"error": err}
	jsonOutput, _ := json.Marshal(output)
	fmt.Println(string(jsonOutput))
	exitFunc(code)
}

func main() {
	runCLI(os.Args, os.Exit)
}

func runCLI(args []string, exitFunc func(int)) {
	if len(args) < 2 {
		printErrorAndExit("invalid arguments", 42, exitFunc)
	}

	command := args[1]

	switch command {
	case "add":
		if len(args) < 3 {
			printErrorAndExit("invalid arguments", 42, exitFunc)
		}
		description := args[2]
		addTask(description, tasksFile, exitFunc)
	case "list":
		listTasks(tasksFile, exitFunc)
	case "done":
		if len(args) < 3 {
			printErrorAndExit("invalid arguments", 42, exitFunc)
		}
		idStr := args[2]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			printErrorAndExit("invalid arguments: task ID must be an integer", 42, exitFunc)
		}
		completeTask(id, tasksFile, exitFunc)
	default:
		printErrorAndExit(fmt.Sprintf("unknown command: %s", command), 42, exitFunc)
	}
}

func addTask(description string, filePath string, exitFunc func(int)) {
	tasks, err := readTasks(filePath)
	if err != nil {
		printErrorAndExit(err.Error(), 1, exitFunc)
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

	if err := saveTasks(filePath, tasks); err != nil {
		printErrorAndExit(err.Error(), 1, exitFunc)
	}

	output := map[string]interface{}{"id": newTask.ID, "status": "added"}
	jsonOutput, _ := json.Marshal(output)
	fmt.Println(string(jsonOutput))
}

func listTasks(filePath string, exitFunc func(int)) {
	tasks, err := readTasks(filePath)
	if err != nil {
		printErrorAndExit(err.Error(), 1, exitFunc)
	}

	var pendingTasks []Task
	for _, task := range tasks {
		if !task.Done {
			pendingTasks = append(pendingTasks, task)
		}
	}

	jsonOutput, _ := json.Marshal(pendingTasks)
	fmt.Println(string(jsonOutput))
}

func completeTask(id int, filePath string, exitFunc func(int)) {
	tasks, err := readTasks(filePath)
	if err != nil {
		printErrorAndExit(err.Error(), 1, exitFunc)
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
		printErrorAndExit(fmt.Sprintf("task with ID %d not found", id), 42, exitFunc)
	}

	if err := saveTasks(filePath, tasks); err != nil {
		printErrorAndExit(err.Error(), 1, exitFunc)
	}

	output := map[string]interface{}{"id": id, "status": "completed"}
	jsonOutput, _ := json.Marshal(output)
	fmt.Println(string(jsonOutput))
}
