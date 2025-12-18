package main

import (
	"os"
	"reflect"
	"testing"
)

// Helper to create a new TodoManager for each test
func newTestTodoManager(t *testing.T) *TodoManager {
	tempFile := "temp_tasks_" + t.Name() + ".json"
	os.Remove(tempFile) // Ensure a clean slate

	mgr := NewTodoManager(tempFile)
	return mgr
}

func TestTodoManager_LoadAndSaveTasks(t *testing.T) {
	mgr := newTestTodoManager(t)
	defer os.Remove(mgr.filename)

	// Directly manipulate tasks for testing Load/Save
	mgr.mu.Lock()
	mgr.tasks = []Task{{ID: 1, Description: "Test task 1", Done: false}}
	mgr.nextID = 2
	mgr.mu.Unlock()

	err := mgr.SaveTasks()
	if err != nil {
		t.Fatalf("SaveTasks failed: %v", err)
	}

	// Create a new manager to ensure loading from file
	newMgr := NewTodoManager(mgr.filename)

	if len(newMgr.tasks) != 1 {
		t.Fatalf("Expected 1 task after loading, got %d", len(newMgr.tasks))
	}

	expectedTask := Task{ID: 1, Description: "Test task 1", Done: false}
	if !reflect.DeepEqual(newMgr.tasks[0], expectedTask) {
		t.Errorf("Loaded task mismatch. Expected %+v, got %+v", expectedTask, newMgr.tasks[0])
	}

	if newMgr.nextID != 2 {
		t.Errorf("Expected nextID to be 2, got %d", newMgr.nextID)
	}
}

func TestTodoManager_AddTask(t *testing.T) {
	mgr := newTestTodoManager(t)
	defer os.Remove(mgr.filename)

	description := "Buy groceries"
	task := mgr.AddTask(description)

	if task.ID != 1 {
		t.Errorf("Expected task ID to be 1, got %d", task.ID)
	}
	if task.Description != description {
		t.Errorf("Expected description %q, got %q", description, task.Description)
	}
	if task.Done != false {
		t.Errorf("Expected task to be not done, got %t", task.Done)
	}

	loadedTasks := mgr.GetTasks()
	if len(loadedTasks) != 1 {
		t.Fatalf("Expected 1 task after adding, got %d", len(loadedTasks))
	}
	if loadedTasks[0].Description != description {
		t.Errorf("Loaded task description mismatch. Expected %q, got %q", description, loadedTasks[0].Description)
	}

	// Add another task to check ID increment
	description2 := "Pay bills"
	task2 := mgr.AddTask(description2)

	if task2.ID != 2 {
		t.Errorf("Expected second task ID to be 2, got %d", task2.ID)
	}
}

func TestTodoManager_GetTasks(t *testing.T) {
	mgr := newTestTodoManager(t)
	defer os.Remove(mgr.filename)

	mgr.AddTask("Task A")
	mgr.AddTask("Task B")

	loadedTasks := mgr.GetTasks()

	if len(loadedTasks) != 2 {
		t.Fatalf("Expected 2 tasks, got %d", len(loadedTasks))
	}

	expectedDescriptions := []string{"Task A", "Task B"}
	for i, task := range loadedTasks {
		if task.Description != expectedDescriptions[i] {
			t.Errorf("Task description mismatch at index %d. Expected %q, got %q", i, expectedDescriptions[i], task.Description)
		}
	}
}

func TestTodoManager_CompleteTask(t *testing.T) {
	mgr := newTestTodoManager(t)
	defer os.Remove(mgr.filename)

	task1 := mgr.AddTask("Task to complete 1")
	mgr.AddTask("Task to complete 2")

	completedTask, err := mgr.CompleteTask(task1.ID)
	if err != nil {
		t.Fatalf("CompleteTask failed: %v", err)
	}

	if completedTask.Done != true {
		t.Errorf("Expected task to be done, got %t", completedTask.Done)
	}

	loadedTasks := mgr.GetTasks()
	found := false
	for _, task := range loadedTasks {
		if task.ID == task1.ID {
			found = true
			if task.Done != true {
				t.Errorf("Task in loaded list not marked as done")
			}
			break
		}
	}
	if !found {
		t.Errorf("Completed task not found in loaded list")
	}

	// Test completing a non-existent task
	_, err = mgr.CompleteTask(999)
	if err == nil {
		t.Errorf("Expected error when completing non-existent task, got nil")
	}
	expectedErrorMsg := "task with ID 999 not found"
	if err.Error() != expectedErrorMsg {
		t.Errorf("Expected error message %q, got %q", expectedErrorMsg, err.Error())
	}
}
