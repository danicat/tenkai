package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"todo-cli/internal/todo"
)

// JSONStorer implements the todo.Storer interface for JSON file persistence.
type JSONStorer struct {
	FilePath string
}

// NewJSONStorer creates a new JSONStorer.
func NewJSONStorer(filePath string) *JSONStorer {
	return &JSONStorer{FilePath: filePath}
}

// Load loads tasks from the JSON file.
func (s *JSONStorer) Load() ([]*todo.Task, error) {
	_, err := os.Stat(s.FilePath)
	if os.IsNotExist(err) {
		// If file doesn't exist, return empty list and no error
		return []*todo.Task{},
			nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to stat file %s: %w", s.FilePath, err)
	}

	data, err := ioutil.ReadFile(s.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", s.FilePath, err)
	}

	if len(data) == 0 {
		return []*todo.Task{},
			nil
	}

	var tasks []*todo.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tasks from %s: %w", s.FilePath, err)
	}

	return tasks, nil
}

// Save saves tasks to the JSON file.
func (s *JSONStorer) Save(tasks []*todo.Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tasks: %w", err)
	}

	if err := ioutil.WriteFile(s.FilePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write tasks to file %s: %w", s.FilePath, err)
	}

	return nil
}
