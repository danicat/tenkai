package task

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

var tasksFile = "tasks.json"

// SetTasksFile allows changing the tasks file path, primarily for testing.
func SetTasksFile(filepath string) {
	tasksFile = filepath
}

// Task represents a single todo item.
type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

// LoadTasks reads tasks from the tasksFile.
func LoadTasks() ([]Task, error) {
	data, err := os.ReadFile(tasksFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil // Return empty list if file doesn't exist
		}
		return nil, fmt.Errorf("reading tasks file: %w", err)
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("unmarshalling tasks: %w", err)
	}
	return tasks, nil
}

// SaveTasks writes tasks to the tasksFile.
func SaveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling tasks: %w", err)
	}

	if err := os.WriteFile(tasksFile, data, 0644); err != nil {
		return fmt.Errorf("writing tasks file: %w", err)
	}
	return nil
}

// NextID returns the next available ID for a new task.
func NextID(tasks []Task) int {
	if len(tasks) == 0 {
		return 1
	}
	// Sort tasks by ID to find the maximum
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].ID < tasks[j].ID
	})
	return tasks[len(tasks)-1].ID + 1
}

// AddTask adds a new task to the list and saves it.
func AddTask(description string) (*Task, error) {
	tasks, err := LoadTasks()
	if err != nil {
		return nil, err
	}

	id := NextID(tasks)
	newTask := &Task{
		ID:          id,
		Description: description,
		Done:        false,
	}
	tasks = append(tasks, *newTask)

	if err := SaveTasks(tasks);
	err != nil {
		return nil, err
	}

	return newTask, nil
}

// ListPendingTasks returns all pending tasks.
func ListPendingTasks() ([]Task, error) {
	tasks, err := LoadTasks()
	if err != nil {
		return nil, err
	}

	var pendingTasks []Task
	for _, t := range tasks {
		if !t.Done {
			pendingTasks = append(pendingTasks, t)
		}
	}
	return pendingTasks, nil
}

// CompleteTask marks a task as done and saves the changes.
func CompleteTask(id int) (*Task, error) {
	tasks, err := LoadTasks()
	if err != nil {
		return nil, err	}

	var completedTask *Task
	found := false
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Done = true
			completedTask = &tasks[i]
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("task with ID %d not found", id)
	}

	if err := SaveTasks(tasks);
	err != nil {
		return nil, err
	}

	return completedTask, nil
}