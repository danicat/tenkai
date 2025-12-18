package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

// Task represents a single todo item.
type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

// Global variable to store tasks (will be loaded from/saved to file)
var tasks []Task

var tasksFile = "tasks.json"

func init() {
	if tf := os.Getenv("TASKS_FILE"); tf != "" {
		tasksFile = tf
	}
}

func main() {
	// Subcommands
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	doneCmd := flag.NewFlagSet("done", flag.ExitOnError)

	// Add command flags
	addDescription := addCmd.String("text", "", "Description of the task")

	// Done command flags
	doneID := doneCmd.Int("id", 0, "ID of the task to mark as done")

	if len(os.Args) < 2 {
		fmt.Println(`{"error": "expected 'add', 'list', or 'done' subcommands"}`)
		os.Exit(42)
	}

	loadTasks() // Load tasks at the beginning

	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		if *addDescription == "" {
			fmt.Println(`{"error": "invalid arguments"}`)
			os.Exit(42)
		}
		addTask(*addDescription)
	case "list":
		listCmd.Parse(os.Args[2:])
		listTasks()
	case "done":
		doneCmd.Parse(os.Args[2:])
		if *doneID == 0 {
			fmt.Println(`{"error": "invalid arguments"}`)
			os.Exit(42)
		}
		completeTask(*doneID)
	default:
		fmt.Println(`{"error": "expected 'add', 'list', or 'done' subcommands"}`)
		os.Exit(42)
	}

	saveTasks() // Save tasks before exiting
}

func loadTasks() {
	data, err := os.ReadFile(tasksFile)
	if err != nil {
		if os.IsNotExist(err) {
			tasks = []Task{}
			return
		}
		printErrorAndExit(fmt.Sprintf("failed to read tasks file: %v", err))
	}

	err = json.Unmarshal(data, &tasks)
	if err != nil {
		printErrorAndExit(fmt.Sprintf("failed to unmarshal tasks: %v", err))
	}
}

func saveTasks() {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		printErrorAndExit(fmt.Sprintf("failed to marshal tasks: %v", err))
	}

	err = os.WriteFile(tasksFile, data, 0644)
	if err != nil {
		printErrorAndExit(fmt.Sprintf("failed to write tasks file: %v", err))
	}
}

func addTask(description string) {
	newID := 1
	if len(tasks) > 0 {
		newID = tasks[len(tasks)-1].ID + 1
	}
	task := Task{
		ID:          newID,
		Description: description,
		Done:        false,
	}
	tasks = append(tasks, task)
	fmt.Printf(`{"id": %d, "status": "added"}` + "\n", task.ID)
}

func listTasks() {

pendingTasks := []Task{}
	for _, task := range tasks {
		if !task.Done {
			pendingTasks = append(pendingTasks, task)
		}
	}
	data, err := json.Marshal(pendingTasks)
	if err != nil {
		printErrorAndExit(fmt.Sprintf("failed to marshal tasks for listing: %v", err))
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
		fmt.Println(`{"error": "task not found"}`)
		os.Exit(1)
	}
	fmt.Printf(`{"id": %d, "status": "completed"}` + "\n", id)
}

func printErrorAndExit(msg string) {
	fmt.Printf(`{"error": "%s"}` + "\n", msg)
	os.Exit(1)
}
