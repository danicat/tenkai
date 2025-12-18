package todo

import (
	"fmt"
	"sort"
	"sync"
)

// Task represents a single todo item.
type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

// Storer defines the interface for persisting and loading tasks.
type Storer interface {
	Load() ([]*Task, error)
	Save(tasks []*Task) error
}

// Manager handles todo list operations.
type Manager struct {
	mu    sync.Mutex
	tasks []*Task
	nextID int
	storer Storer
}

// NewManager creates a new Manager with the given Storer.
func NewManager(storer Storer) (*Manager, error) {
	m := &Manager{
		storer: storer,
		nextID: 1, // Start IDs from 1
	}

	// Load tasks on initialization
	if err := m.Load(); err != nil {
		return nil, fmt.Errorf("failed to load tasks: %w", err)
	}

	return m, nil
}

// Add adds a new task to the list.
func (m *Manager) Add(description string) (*Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	task := &Task{
		ID:          m.nextID,
		Description: description,
		Done:        false,
	}
	m.tasks = append(m.tasks, task)
	m.nextID++

	if err := m.Save(); err != nil {
		return nil, fmt.Errorf("failed to save tasks after adding: %w", err)
	}

	return task, nil
}

// List returns all tasks, sorted by ID.
func (m *Manager) List() []*Task {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Return a copy to prevent external modification
	tasks := make([]*Task, len(m.tasks))
	copy(tasks, m.tasks)

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].ID < tasks[j].ID
	})

	return tasks
}

// Complete marks a task as done by its ID.
func (m *Manager) Complete(id int) (*Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, task := range m.tasks {
		if task.ID == id {
			if task.Done {
				return task, nil // Already completed
			}
			task.Done = true
			if err := m.Save(); err != nil {
				return nil, fmt.Errorf("failed to save tasks after completing: %w", err)
			}
			return task, nil
		}
	}
	return nil, fmt.Errorf("task with ID %d not found", id)
}

// Load loads tasks from the storer into the manager.
func (m *Manager) Load() error {
	loadedTasks, err := m.storer.Load()
	if err != nil {
		return err
	}
	m.tasks = loadedTasks
	// Determine nextID based on loaded tasks
	maxID := 0
	for _, task := range m.tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	m.nextID = maxID + 1
	return nil
}

// Save saves the current tasks from the manager to the storer.
func (m *Manager) Save() error {
	return m.storer.Save(m.tasks)
}
