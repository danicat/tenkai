package task

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

var tasksFile = "tasks.json"

// SetTasksFileForTest allows changing the tasks file path for testing.
func SetTasksFileForTest(filename string) {
	tasksFile = filename
}

// Task represents a single todo item.
type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

// Store manages loading and saving tasks to a file.
type Store struct {
	mu    sync.Mutex
	tasks []Task
	nextID int
}

// NewStore creates a new TaskStore and loads tasks from the file.
func NewStore() (*Store, error) {
	s := &Store{
		tasks: make([]Task, 0),
		nextID: 1,
	}
	if err := s.loadTasks(); err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, start with an empty list
			return s, nil
		}
		return nil, fmt.Errorf("failed to load tasks: %w", err)
	}
	s.updateNextID()
	return s, nil
}

// AddTask adds a new task to the store.
func (s *Store) AddTask(description string) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task := Task{
		ID:          s.nextID,
		Description: description,
		Done:        false,
	}
	s.tasks = append(s.tasks, task)
	s.nextID++

	if err := s.saveTasks(); err != nil {
		return Task{}, fmt.Errorf("failed to save tasks after adding: %w", err)
	}
	return task, nil
}

// ListTasks returns all tasks, optionally filtered by done status.
func (s *Store) ListTasks(includeDone bool) []Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	if includeDone {
		return s.tasks
	}

	var pendingTasks []Task
	for _, t := range s.tasks {
		if !t.Done {
			pendingTasks = append(pendingTasks, t)
		}
	}
	return pendingTasks
}

// CompleteTask marks a task as done.
func (s *Store) CompleteTask(id int) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.tasks {
		if s.tasks[i].ID == id {
			if s.tasks[i].Done {
				return s.tasks[i], fmt.Errorf("task %d is already completed", id)
			}
			s.tasks[i].Done = true
			if err := s.saveTasks(); err != nil {
				return Task{}, fmt.Errorf("failed to save tasks after completing: %w", err)
			}
			return s.tasks[i], nil
		}
	}
	return Task{}, fmt.Errorf("task with ID %d not found", id)
}

// loadTasks loads tasks from tasks.json.
func (s *Store) loadTasks() error {
	data, err := os.ReadFile(tasksFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &s.tasks)
}

// saveTasks saves tasks to tasks.json.
func (s *Store) saveTasks() error {
	data, err := json.MarshalIndent(s.tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tasks: %w", err)
	}
	return os.WriteFile(tasksFile, data, 0644)
}

// updateNextID sets nextID to be one greater than the maximum ID in the current tasks.
func (s *Store) updateNextID() {
	maxID := 0
	for _, t := range s.tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	s.nextID = maxID + 1
}

