package task

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

const tasksFile = "tasks.json"

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

var (
	tasks    []Task
	nextID   int
	dataLock sync.Mutex // Mutex to protect access to tasks and nextID
	loadOnce sync.Once  // Ensures loadTasks is called only once
)

func loadTasks() error {

	_, err := os.Stat(tasksFile)
	if os.IsNotExist(err) {
		tasks = []Task{}
		nextID = 1
		return nil
	} else if err != nil {
		return fmt.Errorf("checking tasks file: %w", err)
	}

	file, err := ioutil.ReadFile(tasksFile)
	if err != nil {
		return fmt.Errorf("reading tasks file: %w", err)
	}

	if len(file) == 0 {
		tasks = []Task{}
		nextID = 1
		return nil
	}

	if err := json.Unmarshal(file, &tasks); err != nil {
		return fmt.Errorf("unmarshalling tasks: %w", err)
	}

	nextID = 1
	if len(tasks) > 0 {
		for _, t := range tasks {
			if t.ID >= nextID {
				nextID = t.ID + 1
			}
		}
	}
	return nil
}

func saveTasks() error {

	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling tasks: %w", err)
	}

	if err := ioutil.WriteFile(tasksFile, data, 0644); err != nil {
		return fmt.Errorf("writing tasks file: %w", err)
	}
	return nil
}

func AddTask(description string) (*Task, error) {
	loadOnce.Do(func() {
		// loadTasks does not acquire dataLock internally
		if err := loadTasks(); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading tasks: %v\n", err)
		}
	})

	dataLock.Lock()
	defer dataLock.Unlock()

	if description == "" {
		return nil, fmt.Errorf("description cannot be empty")
	}

	task := Task{
		ID:          nextID,
		Description: description,
		Done:        false,
	}
	tasks = append(tasks, task)
	nextID++

	if err := saveTasks(); err != nil {
		return nil, fmt.Errorf("saving tasks: %w", err)
	}
	return &task, nil
}

func ListTasks() []Task {
	loadOnce.Do(func() {
		if err := loadTasks(); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading tasks: %v\n", err)
		}
	})

	dataLock.Lock()
	defer dataLock.Unlock()
	return tasks
}

func CompleteTask(id int) (*Task, error) {
	loadOnce.Do(func() {
		if err := loadTasks(); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading tasks: %v\n", err)
		}
	})

	dataLock.Lock()
	defer dataLock.Unlock()

	for i, t := range tasks {
		if t.ID == id {
			tasks[i].Done = true
			if err := saveTasks(); err != nil {
				return nil, fmt.Errorf("saving tasks: %w", err)
			}
			return &tasks[i], nil
		}
	}
	return nil, fmt.Errorf("task with ID %d not found", id)
}
