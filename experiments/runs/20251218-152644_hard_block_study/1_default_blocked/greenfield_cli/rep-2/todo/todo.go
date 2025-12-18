package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

var tasksFile = "tasks.json"

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

type AddResult struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

type CompleteResult struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

var (
	mu     sync.Mutex // Mutex to protect access to tasksFile
	nextID = 1
)

func readTasks() ([]Task, error) {
	mu.Lock()
	defer mu.Unlock()

	data, err := ioutil.ReadFile(tasksFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil
		}
		return nil, fmt.Errorf("error reading tasks file: %w", err)
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("error unmarshaling tasks: %w", err)
	}

	if len(tasks) > 0 {
		for _, task := range tasks {
			if task.ID >= nextID {
				nextID = task.ID + 1
			}
		}
	}

	return tasks, nil
}

func writeTasks(tasks []Task) error {
	mu.Lock()
	defer mu.Unlock()

	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling tasks: %w", err)
	}

	if err := ioutil.WriteFile(tasksFile, data, 0644); err != nil {
		return fmt.Errorf("error writing tasks file: %w", err)
	}

	return nil
}

// Add a new task
func Add(description string) (*AddResult, error) {
	if description == "" {
		return nil, errors.New("invalid arguments")
	}

	tasks, err := readTasks()
	if err != nil {
		return nil, err
	}

	newTask := Task{
		ID:          nextID,
		Description: description,
		Done:        false,
	}
	nextID++
	tasks = append(tasks, newTask)

	if err := writeTasks(tasks); err != nil {
		return nil, err
	}

	return &AddResult{ID: newTask.ID, Status: "added"}, nil
}

// List all pending tasks
func List() ([]Task, error) {
	tasks, err := readTasks()
	if err != nil {
		return nil, err
	}

	var pendingTasks []Task
	for _, task := range tasks {
		if !task.Done {
			pendingTasks = append(pendingTasks, task)
		}
	}

	return pendingTasks, nil
}

// Complete a task by ID
func Complete(id int) (*CompleteResult, error) {
	if id <= 0 {
		return nil, errors.New("invalid arguments")
	}

	tasks, err := readTasks()
	if err != nil {
		return nil, err
	}

	found := false
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Done = true
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("task with ID %d not found", id)
	}

	if err := writeTasks(tasks); err != nil {
		return nil, err
	}

	return &CompleteResult{ID: id, Status: "completed"}, nil
}
