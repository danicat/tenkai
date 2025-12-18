package todo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

var tasksFile = "tasks.json"

// Task represents a single todo item.
type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

// Global variable to hold tasks and a mutex for concurrent access.
var (
	tasks     []Task
	nextID    int
	tasksLock sync.Mutex
)

func init() {
	// Initialize nextID to 1. It will be updated after loading tasks.
	nextID = 1
}

// loadTasks loads tasks from the tasksFile.
// This function should only be called when tasksLock is already held.
func loadTasks() error {
	data, err := ioutil.ReadFile(tasksFile)
	if err != nil {
		if os.IsNotExist(err) {
			tasks = []Task{}
			nextID = 1
			return nil // File doesn't exist, start with empty tasks
		}
		return fmt.Errorf("reading tasks file: %w", err)
	}

	if len(data) == 0 {
		tasks = []Task{}
		nextID = 1
		return nil
	}

	if err := json.Unmarshal(data, &tasks); err != nil {
		return fmt.Errorf("unmarshaling tasks: %w", err)
	}

	// Determine the next available ID
	if len(tasks) > 0 {
		for _, task := range tasks {
			if task.ID >= nextID {
				nextID = task.ID + 1
			}
		}
	} else {
		nextID = 1
	}

	return nil
}

// saveTasks saves tasks to the tasksFile.
// This function should only be called when tasksLock is already held.
func saveTasks() error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling tasks: %w", err)
	}

	if err := ioutil.WriteFile(tasksFile, data, 0644); err != nil {
		return fmt.Errorf("writing tasks file: %w", err)
	}
	return nil
}

// AddTask adds a new task with the given description.
func AddTask(description string) (Task, error) {
	tasksLock.Lock()
	defer tasksLock.Unlock()

	if err := loadTasks(); err != nil {
		return Task{}, err
	}

	newTask := Task{
		ID:          nextID,
		Description: description,
		Done:        false,
	}
	tasks = append(tasks, newTask)
	nextID++

	if err := saveTasks(); err != nil {
		return Task{}, err
	}
	return newTask, nil
}

// ListTasks returns all tasks, filtering by 'done' status if specified.
func ListTasks(includeDone bool) ([]Task, error) {
	tasksLock.Lock()
	defer tasksLock.Unlock()

	if err := loadTasks(); err != nil {
		return nil, err
	}

	if includeDone {
		return tasks, nil
	}

	var pendingTasks []Task
	for _, task := range tasks {
		if !task.Done {
			pendingTasks = append(pendingTasks, task)
		}
	}
	return pendingTasks, nil
}

// CompleteTask marks a task as done by its ID.
func CompleteTask(id int) (Task, error) {
	tasksLock.Lock()
	defer tasksLock.Unlock()

	if err := loadTasks(); err != nil {
		return Task{}, err
	}

	for i := range tasks {
		if tasks[i].ID == id {
			if tasks[i].Done {
				return tasks[i], fmt.Errorf("task %d is already completed", id)
			}
			tasks[i].Done = true
			if err := saveTasks(); err != nil {
				return tasks[i], err
			}
			return tasks[i], nil
		}
	}
	return Task{}, fmt.Errorf("task with ID %d not found", id)
}
