package task_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"experiment/internal/task"
)

// Helper function to create a temporary tasks.json file
func createTempTasksFile(t *testing.T, content string) string {
	tmpFile, err := ioutil.TempFile("", "tasks_test_*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer tmpFile.Close()

	if content != "" {
		_, err = tmpFile.WriteString(content)
		if err != nil {
			t.Fatalf("failed to write to temp file: %v", err)
		}
	}
	return tmpFile.Name()
}

// TestNewTasksManager tests the initialization of TasksManager
func TestNewTasksManager(t *testing.T) {
	t.Parallel()

	// Test with a non-existent file
	nonExistentFile := filepath.Join(os.TempDir(), "non_existent_tasks.json")
	tm, err := task.NewTasksManager(nonExistentFile)
	if err != nil {
		t.Fatalf("unexpected error creating manager with non-existent file: %v", err)
	}
	if tm == nil {
		t.Fatal("TasksManager should not be nil")
	}
	if len(tm.List()) != 0 {
		t.Errorf("expected no tasks, got %d", len(tm.List()))
	}
	os.Remove(nonExistentFile) // Clean up

	// Test with an existing file
	existingContent := `[{"id": 1, "description": "Test Task 1", "done": false}]`
	existingFile := createTempTasksFile(t, existingContent)
	defer os.Remove(existingFile)

	tm, err = task.NewTasksManager(existingFile)
	if err != nil {
		t.Fatalf("unexpected error creating manager with existing file: %v", err)
	}
	if len(tm.List()) != 1 {
		t.Errorf("expected 1 task, got %d", len(tm.List()))
	}
	if tm.List()[0].ID != 1 || tm.List()[0].Description != "Test Task 1" {
		t.Errorf("task mismatch, got: %+v", tm.List()[0])
	}

	// Test with an existing file with multiple tasks to ensure nextID is correct
	existingContentMulti := `[{"id": 10, "description": "T1", "done": false}, {"id": 20, "description": "T2", "done": true}]`
	existingFileMulti := createTempTasksFile(t, existingContentMulti)
	defer os.Remove(existingFileMulti)

	tm, err = task.NewTasksManager(existingFileMulti)
	if err != nil {
		t.Fatalf("unexpected error creating manager with existing multi-task file: %v", err)
	}
	// Add a new task to check next ID
	newTask, err := tm.Add("New task after existing")
	if err != nil {
		t.Fatalf("failed to add task: %v", err)
	}
	if newTask.ID != 21 { // Should be maxID + 1
		t.Errorf("expected nextID to be 21, got %d", newTask.ID)
	}
}

// TestAdd tests adding a task.
func TestAdd(t *testing.T) {
	t.Parallel()
	filePath := createTempTasksFile(t, "")
	defer os.Remove(filePath)

	tm, err := task.NewTasksManager(filePath)
	if err != nil {
		t.Fatalf("failed to create TasksManager: %v", err)
	}

	description := "Buy groceries"
	newTask, err := tm.Add(description)
	if err != nil {
		t.Fatalf("failed to add task: %v", err)
	}

	if newTask.Description != description {
		t.Errorf("expected description %q, got %q", description, newTask.Description)
	}
	if newTask.ID != 1 {
		t.Errorf("expected ID 1, got %d", newTask.ID)
	}
	if newTask.Done != false {
		t.Errorf("expected Done to be false, got %t", newTask.Done)
	}

	// Verify task is in the list
	tasks := tm.List()
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task in list, got %d", len(tasks))
	}
	if tasks[0].Description != description {
		t.Errorf("expected listed task description %q, got %q", description, tasks[0].Description)
	}

	// Verify task is persisted
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read tasks file: %v", err)
	}
	var loadedTasks []task.Task
	if err := json.Unmarshal(data, &loadedTasks); err != nil {
		t.Fatalf("failed to unmarshal tasks from file: %v", err)
	}
	if len(loadedTasks) != 1 {
		t.Fatalf("expected 1 task in persisted file, got %d", len(loadedTasks))
	}
	if loadedTasks[0].Description != description {
		t.Errorf("expected persisted task description %q, got %q", description, loadedTasks[0].Description)
	}

	// Add another task to check ID increment
	secondTask, err := tm.Add("Walk the dog")
	if err != nil {
		t.Fatalf("failed to add second task: %v", err)
	}
	if secondTask.ID != 2 {
		t.Errorf("expected second task ID 2, got %d", secondTask.ID)
	}
}

// TestList tests listing tasks.
func TestList(t *testing.T) {
	t.Parallel()
	filePath := createTempTasksFile(t, "")
	defer os.Remove(filePath)

	tm, err := task.NewTasksManager(filePath)
	if err != nil {
		t.Fatalf("failed to create TasksManager: %v", err)
	}

	tm.Add("Task 1")
	tm.Add("Task 2")

	tasks := tm.List()
	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}
	if tasks[0].Description != "Task 1" || tasks[1].Description != "Task 2" {
		t.Errorf("task descriptions mismatch: %+v", tasks)
	}
}

// TestComplete tests completing a task.
func TestComplete(t *testing.T) {
	t.Parallel()
	filePath := createTempTasksFile(t, "")
	defer os.Remove(filePath)

	tm, err := task.NewTasksManager(filePath)
	if err != nil {
		t.Fatalf("failed to create TasksManager: %v", err)
	}

	addedTask, err := tm.Add("Finish report")
	if err != nil {
		t.Fatalf("failed to add task: %v", err)
	}

	completedTask, err := tm.Complete(addedTask.ID)
	if err != nil {
		t.Fatalf("failed to complete task: %v", err)
	}

	if completedTask.ID != addedTask.ID {
		t.Errorf("expected completed task ID %d, got %d", addedTask.ID, completedTask.ID)
	}
	if completedTask.Done != true {
		t.Errorf("expected task to be done, got %t", completedTask.Done)
	}

	// Verify task is updated in the list
	tasks := tm.List()
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
	if tasks[0].ID != addedTask.ID || tasks[0].Done != true {
		t.Errorf("listed task not updated: %+v", tasks[0])
	}

	// Verify task is persisted as done
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read tasks file: %v", err)
	}
	var loadedTasks []task.Task
	if err := json.Unmarshal(data, &loadedTasks); err != nil {
		t.Fatalf("failed to unmarshal tasks from file: %v", err)
	}
	if len(loadedTasks) != 1 {
		t.Fatalf("expected 1 task in persisted file, got %d", len(loadedTasks))
	}
	if loadedTasks[0].ID != addedTask.ID || loadedTasks[0].Done != true {
		t.Errorf("persisted task not updated: %+v", loadedTasks[0])
	}

	// Test completing a non-existent task
	_, err = tm.Complete(999)
	if err == nil {
		t.Error("expected error when completing non-existent task, got none")
	}
	expectedErrMsg := "task with ID 999 not found"
	if err.Error() != expectedErrMsg {
		t.Errorf("expected error message %q, got %q", expectedErrMsg, err.Error())
	}
}
