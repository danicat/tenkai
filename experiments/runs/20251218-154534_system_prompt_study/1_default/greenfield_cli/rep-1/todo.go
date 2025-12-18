package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const tasksFile = "tasks.json"

// Task represents a single todo item.
type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

// loadTasks reads tasks from the tasksFile.
func loadTasks() ([]Task, error) {
	_, err := os.Stat(tasksFile)
	if os.IsNotExist(err) {
		return []Task{}, nil // Return empty list if file doesn't exist
	} else if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(tasksFile)
	if err != nil {
		return nil, err
	}

	var tasks []Task
	if len(data) == 0 {
		return []Task{}, nil
	}

	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// saveTasks writes tasks to the tasksFile.
func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(tasksFile, data, 0644)
}
