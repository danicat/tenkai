package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
)

var tasksFilePath = "tasks.json"

// Task represents a single todo item
type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

var (
	tasks     []Task
	nextID    int
	tasksLock sync.Mutex // To protect tasks slice and nextID
	osExit    = os.Exit  // To allow overriding os.Exit in tests
)

func init() {
	tasksLock.Lock()
	defer tasksLock.Unlock()

	loadedTasks, err := loadTasks()
	if err != nil {
		if !os.IsNotExist(err) {
			// If file exists but has other error, log it. For a simple CLI, we just start fresh.
			// In a real app, proper error handling or logging would be needed.
			fmt.Fprintf(os.Stderr, `{"error": "failed to load tasks: %v"}
`, err)
			osExit(1)
		}
		// tasks.json does not exist, initialize empty
		tasks = []Task{}
		nextID = 1
		return
	}
	tasks = loadedTasks
	if len(tasks) > 0 {
		nextID = tasks[len(tasks)-1].ID + 1
	} else {
		nextID = 1
	}
}

func loadTasks() ([]Task, error) {
	data, err := ioutil.ReadFile(tasksFilePath)
	if err != nil {
		return nil, err
	}
	var loadedTasks []Task
	err = json.Unmarshal(data, &loadedTasks)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling tasks: %w", err)
	}
	return loadedTasks, nil
}

func saveTasks(t []Task) error {
	tasksLock.Lock()
	defer tasksLock.Unlock()

	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling tasks: %w", err)
	}
	err = ioutil.WriteFile(tasksFilePath, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing tasks file: %w", err)
	}
	return nil
}

func handleDoneCommand(args []string) {
	if len(args) < 1 || args[0] == "" {
		printErrorAndExit(42, "missing task ID")
	}

	taskID, err := parseTaskID(args[0])
	if err != nil {
		printErrorAndExit(42, fmt.Sprintf("invalid task ID: %v", err))
	}

	tasksLock.Lock()
	defer tasksLock.Unlock()

	found := false
	for i := range tasks {
		if tasks[i].ID == taskID {
			tasks[i].Done = true
			found = true
			break
		}
	}

	if !found {
		printErrorAndExit(42, fmt.Sprintf("task with ID %d not found", taskID))
	}

	err = saveTasks(tasks)
	if err != nil {
		printErrorAndExit(1, fmt.Sprintf("failed to save tasks: %v", err))
	}

	output := map[string]interface{}{"id": taskID, "status": "completed"}
	jsonOutput, err := json.Marshal(output)
	if err != nil {
		printErrorAndExit(1, fmt.Sprintf("failed to marshal JSON output: %v", err))
	}
	fmt.Println(string(jsonOutput))
}

// Helper to parse task ID
func parseTaskID(s string) (int, error) {
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("%q is not a valid integer", s)
	}
	if id <= 0 {
		return 0, fmt.Errorf("task ID must be a positive integer")
	}
	return id, nil
}

func handleListCommand() {
	tasksLock.Lock()
	defer tasksLock.Unlock()

	var pendingTasks []Task
	for _, task := range tasks {
		if !task.Done {
			pendingTasks = append(pendingTasks, task)
		}
	}

	jsonOutput, err := json.Marshal(pendingTasks)
	if err != nil {
		printErrorAndExit(1, fmt.Sprintf("failed to marshal JSON output: %v", err))
	}
	fmt.Println(string(jsonOutput))
}

func handleAddCommand(args []string) {
	if len(args) < 1 || args[0] == "" {
		printErrorAndExit(42, "missing task description")
	}

	description := args[0]

	tasksLock.Lock()
	defer tasksLock.Unlock()

	newTask := Task{
		ID:          nextID,
		Description: description,
		Done:        false,
	}
	tasks = append(tasks, newTask)
	nextID++

	err := saveTasks(tasks)
	if err != nil {
		printErrorAndExit(1, fmt.Sprintf("failed to save tasks: %v", err))
	}

	output := map[string]interface{}{"id": newTask.ID, "status": "added"}
	jsonOutput, err := json.Marshal(output)
	if err != nil {
		printErrorAndExit(1, fmt.Sprintf("failed to marshal JSON output: %v", err))
	}
	fmt.Println(string(jsonOutput))
}

func printErrorAndExit(code int, msg string) {
	output := map[string]string{"error": msg}
	jsonOutput, err := json.Marshal(output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling error message: %v\n", err)
		osExit(1)
	}
	fmt.Println(string(jsonOutput))
	osExit(code)
}

func main() {
	if len(os.Args) < 2 {
		printErrorAndExit(42, "missing command")
	}

	switch os.Args[1] {
	case "add":
		handleAddCommand(os.Args[2:])
	case "list":
		handleListCommand()
	case "done":
		handleDoneCommand(os.Args[2:])
	default:
		printErrorAndExit(42, fmt.Sprintf("unknown command: %s", os.Args[1]))
	}
}
