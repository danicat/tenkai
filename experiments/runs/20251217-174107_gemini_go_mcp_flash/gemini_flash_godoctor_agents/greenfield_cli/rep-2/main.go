package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

const tasksFile = "tasks.json"

// Task represents a single todo item
type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

func loadOrCreateTasks(filename string) ([]Task, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil // Return empty list if file doesn't exist
		}
		return nil, err // Other read error
	}

	if len(data) == 0 {
		return []Task{}, nil // Return empty list if file is empty
	}

	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func saveTasks(filename string, tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

func formatError(err string) string {
	output := map[string]string{"error": err}
	jsonOutput, _ := json.Marshal(output)
	return string(jsonOutput)
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

// runApp contains the core logic of the CLI, returning the output string and exit code.
func runApp(args []string) (string, int) {
	tasks, err := loadOrCreateTasks(tasksFile)
	if err != nil {
		return formatError(fmt.Sprintf("failed to load tasks: %v", err)), 1
	}

	// Subcommands
	addCmd := flag.NewFlagSet("add", flag.ContinueOnError) // Changed to ContinueOnError
	addDescription := addCmd.String("description", "", "Description of the task")

	listCmd := flag.NewFlagSet("list", flag.ContinueOnError) // Changed to ContinueOnError

	doneCmd := flag.NewFlagSet("done", flag.ContinueOnError) // Changed to ContinueOnError
	doneID := doneCmd.Int("id", 0, "ID of the task to complete")

	// Prevent flag.Parse from printing usage on error, handle it manually
	addCmd.SetOutput(ioutil.Discard)
	listCmd.SetOutput(ioutil.Discard)
	doneCmd.SetOutput(ioutil.Discard)

	if len(args) < 2 {
		return formatError("no command provided"), 42
	}

	switch args[1] {
	case "add":
		err := addCmd.Parse(args[2:])
		if err != nil || *addDescription == "" {
			return formatError("invalid arguments: description is required"), 42
		}
		task := Task{
			ID:          getNextID(tasks),
			Description: *addDescription,
			Done:        false,
		}
		tasks = append(tasks, task)
		if err := saveTasks(tasksFile, tasks); err != nil {
			return formatError(fmt.Sprintf("failed to save tasks: %v", err)), 1
		}
		output := map[string]interface{}{"id": task.ID, "status": "added"}
		jsonOutput, _ := json.Marshal(output)
		return string(jsonOutput), 0

	case "list":
		err := listCmd.Parse(args[2:])
		if err != nil {
			// This case might happen if list command gets unexpected flags
			return formatError(fmt.Sprintf("invalid arguments for list: %v", err)), 42
		}
		pendingTasks := make([]Task, 0)
		for _, task := range tasks {
			if !task.Done {
				pendingTasks = append(pendingTasks, task)
			}
		}
		jsonOutput, err := json.Marshal(pendingTasks)
		if err != nil {
			return formatError(fmt.Sprintf("failed to marshal tasks: %v", err)), 1
		}
		return string(jsonOutput), 0

	case "done":
		err := doneCmd.Parse(args[2:])
		if err != nil || *doneID == 0 {
			return formatError("invalid arguments: task ID is required"), 42
		}
		found := false
		for i := range tasks {
			if tasks[i].ID == *doneID {
				tasks[i].Done = true
				found = true
				break
			}
		}
		if !found {
			return formatError(fmt.Sprintf("task with ID %d not found", *doneID)), 1
		}
		if err := saveTasks(tasksFile, tasks); err != nil {
			return formatError(fmt.Sprintf("failed to save tasks: %v", err)), 1
		}
		output := map[string]interface{}{"id": *doneID, "status": "completed"}
		jsonOutput, _ := json.Marshal(output)
		return string(jsonOutput), 0

	default:
		return formatError(fmt.Sprintf("unknown command: %s", args[1])), 42
	}
}

func main() {
	output, code := runApp(os.Args)
	fmt.Println(output)
	os.Exit(code)
}
