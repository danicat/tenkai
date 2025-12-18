package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var tasksFile = "tasks.json"

// LoadTasks reads tasks from the tasks.json file.
func LoadTasks() ([]Task, error) {
	_, err := os.Stat(tasksFile)
	if os.IsNotExist(err) {
		return []Task{}, nil // Return empty list if file doesn't exist
	}

	data, err := ioutil.ReadFile(tasksFile)
	if err != nil {
		return nil, fmt.Errorf("error reading tasks file: %w", err)
	}

	if len(data) == 0 {
		return []Task{}, nil // Return empty list if file is empty
	}

	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling tasks: %w", err)
	}

	return tasks, nil
}

// SaveTasks writes tasks to the tasks.json file.
func SaveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling tasks: %w", err)
	}

	err = ioutil.WriteFile(tasksFile, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing tasks file: %w", err)
	}

	return nil
}
