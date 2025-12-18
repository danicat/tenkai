package task

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

// TasksManager handles operations related to tasks.
type TasksManager struct {
	mu    sync.Mutex
	tasks []Task
	nextID int
	filePath string
}

// NewTasksManager creates a new TasksManager and loads existing tasks.
func NewTasksManager(filePath string) (*TasksManager, error) {
	tm := &TasksManager{
		filePath: filePath,
		nextID: 1,
	}
	err := tm.loadTasks()
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load tasks: %w", err)
	}
	if len(tm.tasks) > 0 {
		tm.nextID = tm.tasks[len(tm.tasks)-1].ID + 1
	}
	return tm, nil
}

// loadTasks loads tasks from the JSON file.
func (tm *TasksManager) loadTasks() error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	data, err := ioutil.ReadFile(tm.filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &tm.tasks)
}

// saveTasks saves tasks to the JSON file.
func (tm *TasksManager) saveTasks() error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	data, err := json.MarshalIndent(tm.tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tasks: %w", err)
	}
	return ioutil.WriteFile(tm.filePath, data, 0644)
}

// Add adds a new task to the list.
func (tm *TasksManager) Add(description string) (*Task, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	newTask := Task{
		ID:          tm.nextID,
		Description: description,
		Done:        false,
	}
	tm.tasks = append(tm.tasks, newTask)
	tm.nextID++

	if err := tm.saveTasks(); err != nil {
		return nil, fmt.Errorf("failed to save tasks after adding: %w", err)
	}
	return &newTask, nil
}

// List returns all tasks.
func (tm *TasksManager) List() []Task {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	return tm.tasks
}

// Complete marks a task as done by its ID.
func (tm *TasksManager) Complete(id int) (*Task, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	for i := range tm.tasks {
		if tm.tasks[i].ID == id {
			tm.tasks[i].Done = true
			if err := tm.saveTasks(); err != nil {
				return nil, fmt.Errorf("failed to save tasks after completing: %w", err)
			}
			return &tm.tasks[i], nil
		}
	}
	return nil, fmt.Errorf("task with ID %d not found", id)
}
