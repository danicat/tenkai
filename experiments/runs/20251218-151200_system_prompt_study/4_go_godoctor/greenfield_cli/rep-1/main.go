package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

const tasksFile = "tasks.json"

// Task represents a single todo item.
type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

// Global variable to hold tasks in memory
var tasks []Task

var exit = os.Exit

// printError prints an error message in JSON format to stdout and exits with the specified code.
func printError(msg string, code int) {
	fmt.Printf(`{"error": "%s"}
`, msg)
	exit(code)
}

// loadTasks loads tasks from the tasks.json file.
func loadTasks() error {
	data, err := ioutil.ReadFile(tasksFile)
	if err != nil {
		if os.IsNotExist(err) {
			tasks = []Task{} // Initialize an empty slice if file doesn't exist
			return nil
		}
		return fmt.Errorf("error reading tasks file: %w", err)
	}

	if len(data) == 0 {
		tasks = []Task{} // Handle empty file case
		return nil
	}

	if err := json.Unmarshal(data, &tasks); err != nil {
		return fmt.Errorf("error unmarshaling tasks: %w", err)
	}
	return nil
}

// saveTasks saves the current tasks to the tasks.json file.
func saveTasks() error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling tasks: %w", err)
	}
	if err := ioutil.WriteFile(tasksFile, data, 0644); err != nil {
		return fmt.Errorf("error writing tasks file: %w", err)
	}
	return nil
}

// getNextID returns the next available ID for a new task.
func getNextID() int {
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
	if err := loadTasks(); err != nil {
		printError(err.Error(), 1)
	}

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addDescription := addCmd.String("d", "", "Description of the task")


	doneCmd := flag.NewFlagSet("done", flag.ExitOnError)
	doneID := doneCmd.Int("id", 0, "ID of the task to mark as done")

	if len(os.Args) < 2 {
		printError("expected 'add', 'list', or 'done' subcommand", 42)
	}

	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		if *addDescription == "" {
			printError("task description cannot be empty", 42)
		}
		addTask(*addDescription)
	case "list":
		listTasks()
	case "done":
		doneCmd.Parse(os.Args[2:])
		if *doneID == 0 {
			printError("task ID cannot be empty", 42)
		}
		completeTask(*doneID)
	default:
		printError("unknown subcommand", 42)
	}
}

func addTask(description string) {
	id := getNextID()
	newTask := Task{
		ID:          id,
		Description: description,
		Done:        false,
	}
	tasks = append(tasks, newTask)

	if err := saveTasks(); err != nil {
		printError(err.Error(), 1)
	}

	fmt.Printf(`{"id": %d, "status": "added"}
`, id)
}

func listTasks() {
	pendingTasks := []Task{}
	for _, task := range tasks {
		if !task.Done {
			pendingTasks = append(pendingTasks, task)
		}
	}

	data, err := json.MarshalIndent(pendingTasks, "", "  ")
	if err != nil {
		printError(err.Error(), 1)
	}

	fmt.Println(string(data))
}

func completeTask(id int) {
	found := false
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Done = true
			found = true
			break
		}
	}

	if !found {
		printError(fmt.Sprintf("task with ID %d not found", id), 42)
	}

	if err := saveTasks(); err != nil {
		printError(err.Error(), 1)
	}

	fmt.Printf(`{"id": %d, "status": "completed"}
`, id)
}
