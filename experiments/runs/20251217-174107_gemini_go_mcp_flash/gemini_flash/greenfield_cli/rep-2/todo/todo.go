package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	sync"
)

const tasksFile = "tasks.json"

var (
	mu sync.Mutex // Mutex to protect concurrent access to tasks.json
)

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

// AddTask adds a new task to the list.
func AddTask(description string) (Task, error) {
	mu.Lock()
	defer mu.Unlock()

	tasks, err := loadTasks()
	if err != nil {
		return Task{}, err
	}

	id := 1
	if len(tasks) > 0 {
		id = tasks[len(tasks)-1].ID + 1
	}

	newTask := Task{
		ID:          id,
		Description: description,
		Done:        false,
	}

	tasks = append(tasks, newTask)
	err = saveTasks(tasks)
	if err != nil {
		return Task{}, err
	}

	return newTask, nil
}

// ListTasks returns all tasks.
func ListTasks() ([]Task, error) {
	mu.Lock()
	defer mu.Unlock()
	return loadTasks()
}

// CompleteTask marks a task as done.
func CompleteTask(id int) (Task, error) {
	mu.Lock()
	defer mu.Unlock()

	tasks, err := loadTasks()
	if err != nil {
		return Task{}, err
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Done = true
			err = saveTasks(tasks)
			if err != nil {
				return Task{}, err
			}
			return tasks[i], nil
		}
	}

	return Task{}, errors.New("task not found")
}

func loadTasks() ([]Task, error) {
	data, err := ioutil.ReadFile(tasksFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil // No tasks file yet, return empty list
		}
		return nil, fmt.Errorf("failed to read tasks file: %w", err)
	}

	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tasks: %w", err)
	}

	return tasks, nil
}

func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tasks: %w", err)
	}

	err = ioutil.WriteFile(tasksFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write tasks file: %w", err)
	}

	return nil
}
