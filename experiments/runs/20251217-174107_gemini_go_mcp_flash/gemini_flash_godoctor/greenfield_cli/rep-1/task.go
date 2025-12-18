package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

// Task represents a single todo item.
type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

// TodoManager manages the todo list's state and persistence.
type TodoManager struct {
	tasks    []Task
	nextID   int
	mu       sync.Mutex
	filename string
}

// NewTodoManager creates and initializes a new TodoManager.
func NewTodoManager(filename string) *TodoManager {
	t := &TodoManager{
		filename: filename,
		nextID:   1,
		tasks:    []Task{},
	}
	t.LoadTasks() // Load tasks on initialization
	return t
}

// LoadTasks loads tasks from the manager's file.
func (t *TodoManager) LoadTasks() {
	data, err := ioutil.ReadFile(t.filename)
	if err != nil {
		if os.IsNotExist(err) {
			t.tasks = []Task{}
			t.nextID = 1
			return
		}
		// If file exists but can't be read, initialize empty tasks.
		t.tasks = []Task{}
		t.nextID = 1
		return
	}

	if err := json.Unmarshal(data, &t.tasks); err != nil {
		// If unmarshalling fails, initialize empty tasks.
		t.tasks = []Task{}
		t.nextID = 1
		return
	}

	if len(t.tasks) > 0 {
		t.nextID = t.tasks[len(t.tasks)-1].ID + 1
	} else {
		t.nextID = 1
	}
}

// SaveTasks saves the current tasks to the manager's file.
func (t *TodoManager) SaveTasks() error {
	data, err := json.MarshalIndent(t.tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tasks: %w", err)
	}

	if err := ioutil.WriteFile(t.filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write tasks file: %w", err)
	}
	return nil
}

// AddTask adds a new task with the given description.
func (t *TodoManager) AddTask(description string) Task {
	t.mu.Lock()
	defer t.mu.Unlock()

	task := Task{
		ID:          t.nextID,
		Description: description,
		Done:        false,
	}
	t.tasks = append(t.tasks, task)
	t.nextID++
	t.SaveTasks() // Save after adding
	return task
}

// GetTasks returns all tasks.
func (t *TodoManager) GetTasks() []Task {
	t.mu.Lock()
	defer t.mu.Unlock()
	// Return a copy to prevent external modification of the internal slice
	return append([]Task{}, t.tasks...)
}

// CompleteTask marks a task as done by its ID.
func (t *TodoManager) CompleteTask(id int) (Task, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for i, task := range t.tasks {
		if task.ID == id {
			t.tasks[i].Done = true
			t.SaveTasks() // Save after completing
			return t.tasks[i], nil
		}
	}
	return Task{}, fmt.Errorf("task with ID %d not found", id)
}
