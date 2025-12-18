package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
)

const tasksFile = "tasks.json"

// Task represents a single todo item.
type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

// Global variable to hold tasks and a mutex for safe access.
var (
	tasks     []Task
	tasksLock sync.Mutex
	nextID    int
)

func init() {
	tasks = make([]Task, 0)
	nextID = 1
}

// loadTasks reads tasks from the tasksFile.
func loadTasks() error {
	tasksLock.Lock()
	defer tasksLock.Unlock()

	data, err := ioutil.ReadFile(tasksFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No tasks file yet, return empty list
		}
		return fmt.Errorf("error reading tasks file: %w", err)
	}

	if len(data) == 0 {
		tasks = []Task{}
		nextID = 1
		return nil
	}

	if err := json.Unmarshal(data, &tasks); err != nil {
		return fmt.Errorf("error unmarshaling tasks: %w", err)
	}

	// Determine nextID based on existing tasks
	nextID = 1
	for _, task := range tasks {
		if task.ID >= nextID {
			nextID = task.ID + 1
		}
	}
	return nil
}

// saveTasks writes current tasks to the tasksFile.
func saveTasks() error {
	tasksLock.Lock()
	defer tasksLock.Unlock()

	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling tasks: %w", err)
	}

	if err := ioutil.WriteFile(tasksFile, data, 0644); err != nil {
		return fmt.Errorf("error writing tasks file: %w", err)
	}
	return nil
}

// cliError represents an error that should result in a specific exit code.
type cliError struct {
	errorMsg string
	exitCode int
}

func (e *cliError) Error() string {
	return e.errorMsg
}

// newCliError creates a new cliError.
func newCliError(msg string, code int) *cliError {
	return &cliError{errorMsg: msg, exitCode: code}
}

// addTask adds a new task.
func addTask(description string) error {
	if description == "" {
		return newCliError("description cannot be empty", 42)
	}

	tasksLock.Lock()
	defer tasksLock.Unlock()

	task := Task{
		ID:          nextID,
		Description: description,
		Done:        false,
	}
	tasks = append(tasks, task)
	nextID++

	if err := saveTasks(); err != nil {
		return newCliError(err.Error(), 1)
	}

	output := map[string]interface{}{
		"id":     task.ID,
		"status": "added",
	}
	jsonOutput, _ := json.Marshal(output)
	fmt.Println(string(jsonOutput))
	return nil
}

// listTasks lists all pending tasks.
func listTasks() error {
	tasksLock.Lock()
	defer tasksLock.Unlock()

	pendingTasks := []Task{}
	for _, task := range tasks {
		if !task.Done {
			pendingTasks = append(pendingTasks, task)
		}
	}

	jsonOutput, _ := json.Marshal(pendingTasks)
	fmt.Println(string(jsonOutput))
	return nil
}

// completeTask marks a task as done.
func completeTask(id int) error {
	tasksLock.Lock()
	defer tasksLock.Unlock()

	found := false
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Done = true
			found = true
			break
		}
	}

	if !found {
		return newCliError(fmt.Sprintf("task with ID %d not found", id), 42)
	}

	if err := saveTasks(); err != nil {
		return newCliError(err.Error(), 1)
	}

	output := map[string]interface{}{
		"id":     id,
		"status": "completed",
	}
	jsonOutput, _ := json.Marshal(output)
	fmt.Println(string(jsonOutput))
	return nil
}

// runCLI is the main entry point logic, returning an error for testing.
func runCLI(fset *flag.FlagSet, args []string) error {
	// fset.Parse consumes arguments for its flags, then fset.Args() returns non-flag arguments.
	// The 'args' parameter here should already be os.Args[1:] like in main().

	if err := fset.Parse(args); err != nil {
		// This could be due to an invalid flag or a parsing error before command processing
		// flag.ExitOnError handles os.Exit, so we capture the error if it occurs.
		return newCliError(fmt.Sprintf("flag parsing error: %v", err), 42) // Use a generic 42 for flag errors.
	}

	if err := loadTasks(); err != nil {
		return err
	}

	parsedArgs := fset.Args()
	if len(parsedArgs) == 0 {
		return newCliError("no command provided", 42)
	}

	command := parsedArgs[0]

	var err error
	switch command {
	case "add":
		if len(parsedArgs) < 2 {
			err = newCliError("missing task description for 'add' command", 42)
		} else {
			err = addTask(parsedArgs[1])
		}
	case "list":
		err = listTasks()
	case "done":
		if len(parsedArgs) < 2 {
			err = newCliError("missing task ID for 'done' command", 42)
		} else {
			id, parseErr := strconv.Atoi(parsedArgs[1])
			if parseErr != nil {
				err = newCliError(fmt.Sprintf("invalid task ID: %v", parseErr), 42)
			} else {
				err = completeTask(id)
			}
		}
	default:
		err = newCliError(fmt.Sprintf("unknown command: %s", command), 42)
	}

	return err
}

func main() {
	// Use default flag.CommandLine for the actual executable
	err := runCLI(flag.CommandLine, os.Args[1:])
	if err != nil {
		var cerr *cliError
		if errors.As(err, &cerr) {
			fmt.Printf("{\"error\": %q}\n", cerr.Error())
			os.Exit(cerr.exitCode)
		} else {
			fmt.Printf("{\"error\": %q}\n", err.Error())
			os.Exit(1)
		}
	}
}
