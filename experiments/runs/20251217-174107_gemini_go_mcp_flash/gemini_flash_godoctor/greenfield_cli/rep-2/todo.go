package main

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
	"sync"
)

const tasksFile = "tasks.json"

// Task represents a single todo item.
type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

// TodoList manages a list of tasks.
type TodoList struct {
	tasks []Task
	mu    sync.Mutex // Mutex to protect access to tasks
}

// NewTodoList creates a new TodoList instance.
func NewTodoList() *TodoList {
	return &TodoList{
		tasks: []Task{},
	}
}

// LoadTasks loads tasks from the tasksFile.
func (tl *TodoList) LoadTasks() error {
	tl.mu.Lock()
	defer tl.mu.Unlock()

	data, err := os.ReadFile(tasksFile)
	if err != nil {
		if os.IsNotExist(err) {
			tl.tasks = []Task{}
			return nil
		}
		return err
	}

	return json.Unmarshal(data, &tl.tasks)
}

// SaveTasks saves tasks to the tasksFile.
func (tl *TodoList) SaveTasks() error {
	tl.mu.Lock()
	defer tl.mu.Unlock()

	data, err := json.MarshalIndent(tl.tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(tasksFile, data, 0644)
}

// AddTask adds a new task to the list.
func (tl *TodoList) AddTask(description string) (Task, error) {
	tl.mu.Lock()
	defer tl.mu.Unlock()

	id := 1
	if len(tl.tasks) > 0 {
		// Find the maximum ID and increment it for the new task
		sort.Slice(tl.tasks, func(i, j int) bool {
			return tl.tasks[i].ID < tl.tasks[j].ID
		})
		id = tl.tasks[len(tl.tasks)-1].ID + 1
	}

	task := Task{
		ID:          id,
		Description: description,
		Done:        false,
	}
	tl.tasks = append(tl.tasks, task)

	return task, nil
}

// ListPendingTasks returns all pending tasks.
func (tl *TodoList) ListPendingTasks() []Task {
	tl.mu.Lock()
	defer tl.mu.Unlock()

	var pending []Task
	for _, task := range tl.tasks {
		if !task.Done {
			pending = append(pending, task)
		}
	}
	return pending
}

// CompleteTask marks a task as done by its ID.
func (tl *TodoList) CompleteTask(id int) (Task, error) {
	tl.mu.Lock()
	defer tl.mu.Unlock()

	for i, task := range tl.tasks {
		if task.ID == id {
			tl.tasks[i].Done = true
			return tl.tasks[i], nil
		}
	}
	return Task{}, errors.New("task not found")
}
