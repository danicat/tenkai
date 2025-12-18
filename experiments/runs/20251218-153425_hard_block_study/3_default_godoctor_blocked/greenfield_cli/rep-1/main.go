package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

var TasksFile = "tasks.json"

// TodoTask represents a single todo item
type TodoTask struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

// Exit functions for consistent error output
func ExitWithError(output io.Writer, err string, code int) int {
	out := map[string]string{"error": err}
	jsonOutput, _ := json.Marshal(out)
	fmt.Fprintln(output, string(jsonOutput))
	return code
}

func ExitWithJSON(output io.Writer, data interface{}) int {
	jsonOutput, err := json.Marshal(data)
	if err != nil {
		return ExitWithError(output, "failed to marshal JSON output", 1)
	}
	fmt.Fprintln(output, string(jsonOutput))
	return 0
}

func LoadTasks() ([]TodoTask, error) {
	data, err := ioutil.ReadFile(TasksFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []TodoTask{}, nil // No tasks file yet, return empty slice
		}
		return nil, fmt.Errorf("failed to read tasks file: %w", err)
	}

	var tasks []TodoTask
	if len(data) == 0 {
		return []TodoTask{}, nil // Empty file, return empty slice
	}
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tasks: %w", err)
	}
	return tasks, nil
}

func SaveTasks(tasks []TodoTask) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tasks: %w", err)
	}
	err = ioutil.WriteFile(TasksFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write tasks file: %w", err)
	}
	return nil
}

func runCLI(args []string, output io.Writer) int {
	if len(args) < 2 {
		return ExitWithError(output, "invalid arguments", 42)
	}

	switch args[1] {
	case "add":
		addCmd := flag.NewFlagSet("add", flag.ContinueOnError)
		var description string
		addCmd.StringVar(&description, "description", "", "Task description")

		// Temporarily redirect stderr to capture flag parsing errors
		oldStderr := os.Stderr
		r, w, _ := os.Pipe()
		os.Stderr = w
		err := addCmd.Parse(args[2:])
		w.Close()
		r.Close() // Close the read end as well
		// Restore stderr
		os.Stderr = oldStderr

		if err != nil { // Check for flag parsing errors
			return ExitWithError(output, "invalid arguments", 42)
		}

		if description == "" && len(addCmd.Args()) > 0 {
			description = addCmd.Args()[0]
		}
		if description == "" {
			return ExitWithError(output, "invalid arguments", 42)
		}

		tasks, err := LoadTasks()
		if err != nil {
			return ExitWithError(output, err.Error(), 1)
		}

		nextID := 1
		if len(tasks) > 0 {
			nextID = tasks[len(tasks)-1].ID + 1
		}

		newTask := TodoTask{
			ID:          nextID,
			Description: description,
			Done:        false,
		}
		tasks = append(tasks, newTask)

		err = SaveTasks(tasks)
		if err != nil {
			return ExitWithError(output, err.Error(), 1)
		}

		return ExitWithJSON(output, map[string]interface{}{"id": newTask.ID, "status": "added"})

	case "list":
		listCmd := flag.NewFlagSet("list", flag.ContinueOnError)
		// Temporarily redirect stderr to capture flag parsing errors
		oldStderr := os.Stderr
		r, w, _ := os.Pipe()
		os.Stderr = w
		err := listCmd.Parse(args[2:])
		w.Close()
		r.Close() // Close the read end as well
		// Restore stderr
		os.Stderr = oldStderr

		if err != nil { // Check for flag parsing errors
			return ExitWithError(output, "invalid arguments", 42)
		}

		tasks, err := LoadTasks()
		if err != nil {
			return ExitWithError(output, err.Error(), 1)
		}

		var pendingTasks []TodoTask
		for _, task := range tasks {
			if !task.Done {
				pendingTasks = append(pendingTasks, task)
			}
		}
		return ExitWithJSON(output, pendingTasks)

	case "done":
		doneCmd := flag.NewFlagSet("done", flag.ContinueOnError)
		// Temporarily redirect stderr to capture flag parsing errors
		oldStderr := os.Stderr
		r, w, _ := os.Pipe()
		os.Stderr = w
		err := doneCmd.Parse(args[2:])
		w.Close()
		r.Close() // Close the read end as well
		// Restore stderr
		os.Stderr = oldStderr

		if err != nil { // Check for flag parsing errors
			return ExitWithError(output, "invalid arguments", 42)
		}

		if len(doneCmd.Args()) == 0 {
			return ExitWithError(output, "invalid arguments", 42)
		}

		taskID, err := strconv.Atoi(doneCmd.Args()[0])
		if err != nil {
			return ExitWithError(output, "invalid task ID", 42)
		}

		tasks, err := LoadTasks()
		if err != nil {
			return ExitWithError(output, err.Error(), 1)
		}

		found := false
		for i, task := range tasks {
			if task.ID == taskID {
				tasks[i].Done = true
				found = true
				break
			}
		}

		if !found {
			return ExitWithError(output, fmt.Sprintf("task with ID %d not found", taskID), 42)
		}

		err = SaveTasks(tasks)
		if err != nil {
			return ExitWithError(output, err.Error(), 1)
		}

		return ExitWithJSON(output, map[string]interface{}{"id": taskID, "status": "completed"})

	default:
		return ExitWithError(output, "unknown command", 42)
	}
}

func main() {
	os.Exit(runCLI(os.Args, os.Stdout))
}
