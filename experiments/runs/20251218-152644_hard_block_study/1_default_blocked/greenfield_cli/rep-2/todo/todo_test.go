package todo

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Setup: Create a temporary tasks.json for testing
	tmpFile, err := ioutil.TempFile("", "tasks_test_*.json")
	if err != nil {
		panic(err)
	}
	tmpFileName := tmpFile.Name()
	tmpFile.Close()

	// Override the tasksFile variable for tests
	oldTasksFile := tasksFile
	tasksFile = tmpFileName

	code := m.Run()

	// Teardown: Remove the temporary tasks.json
	os.Remove(tmpFileName)

	// Restore the original tasksFile variable
	tasksFile = oldTasksFile
	os.Exit(code)
}

func TestAddListTasks(t *testing.T) {
	// Clear tasks before test
	err := ioutil.WriteFile(tasksFile, []byte("[]"), 0644)
	if err != nil {
		t.Fatalf("Failed to clear tasks file: %v", err)
	}
	nextID = 1 // Reset nextID for consistent test results

	addResult, err := Add("Buy milk")
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}
	if addResult.ID != 1 || addResult.Status != "added" {
		t.Errorf("Unexpected add result: %+v", addResult)
	}

	addResult2, err := Add("Walk dog")
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}
	if addResult2.ID != 2 || addResult2.Status != "added" {
		t.Errorf("Unexpected add result: %+v", addResult2)
	}

	tasks, err := List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}

	if tasks[0].Description != "Buy milk" || tasks[0].Done != false {
		t.Errorf("Unexpected task 1: %+v", tasks[0])
	}
	if tasks[1].Description != "Walk dog" || tasks[1].Done != false {
		t.Errorf("Unexpected task 2: %+v", tasks[1])
	}
}

func TestCompleteTask(t *testing.T) {
	// Clear tasks before test
	err := ioutil.WriteFile(tasksFile, []byte("[]"), 0644)
	if err != nil {
		t.Fatalf("Failed to clear tasks file: %v", err)
	}
	nextID = 1 // Reset nextID for consistent test results

	Add("Buy groceries") // ID should be 1
	Add("Pay bills")     // ID should be 2

	completeResult, err := Complete(1)
	if err != nil {
		t.Fatalf("Complete failed: %v", err)
	}
	if completeResult.ID != 1 || completeResult.Status != "completed" {
		t.Errorf("Unexpected complete result: %+v", completeResult)
	}

	tasks, err := List()
	if err != nil {
		t.Fatalf("List failed after complete: %v", err)
	}

	if len(tasks) != 1 || tasks[0].ID != 2 {
		t.Errorf("Expected 1 pending task with ID 2, got: %+v", tasks)
	}

	// Test completing a non-existent task
	err = nil
	_, err = Complete(999)
	if err == nil || err.Error() != "task with ID 999 not found" {
		t.Errorf("Expected 'task with ID 999 not found' error, got: %v", err)
	}

	// Test completing an already completed task (should still return completed for idempotent)
	completeResult2, err := Complete(1)
	if err != nil {
		t.Fatalf("Complete of already done task failed: %v", err)
	}
	if completeResult2.ID != 1 || completeResult2.Status != "completed" {
		t.Errorf("Unexpected complete result for already done task: %+v", completeResult2)
	}
}

func TestAddErrorHandling(t *testing.T) {
	// Clear tasks before test
	err := ioutil.WriteFile(tasksFile, []byte("[]"), 0644)
	if err != nil {
		t.Fatalf("Failed to clear tasks file: %v", err)
	}
	nextID = 1 // Reset nextID for consistent test results

	_, err = Add("")
	if err == nil || err.Error() != "invalid arguments" {
		t.Errorf("Expected 'invalid arguments' error for empty description, got: %v", err)
	}
}

func TestCompleteErrorHandling(t *testing.T) {
	// Clear tasks before test
	err := ioutil.WriteFile(tasksFile, []byte("[]"), 0644)
	if err != nil {
		t.Fatalf("Failed to clear tasks file: %v", err)
	}
	nextID = 1 // Reset nextID for consistent test results

	_, err = Complete(0)
	if err == nil || err.Error() != "invalid arguments" {
		t.Errorf("Expected 'invalid arguments' error for ID 0, got: %v", err)
	}

	_, err = Complete(-1)
	if err == nil || err.Error() != "invalid arguments" {
		t.Errorf("Expected 'invalid arguments' error for negative ID, got: %v", err)
	}
}

func TestNextIDPersistence(t *testing.T) {
	// Clear tasks before test
	err := ioutil.WriteFile(tasksFile, []byte("[]"), 0644)
	if err != nil {
		t.Fatalf("Failed to clear tasks file: %v", err)
	}
	nextID = 1 // Reset nextID for consistent test results

	// Add a task, then simulate a restart by re-reading tasks.
	Add("Task One") // ID 1
	Add("Task Two") // ID 2

	// Simulate a fresh start by reloading tasks
	_, err = readTasks()
	if err != nil {
		t.Fatalf("readTasks failed: %v", err)
	}

	// nextID should be updated from the loaded tasks
	if nextID != 3 {
		t.Errorf("Expected nextID to be 3 after loading existing tasks, got %d", nextID)
	}

	addResult, err := Add("Task Three")
	if err != nil {
		t.Fatalf("Add after restart failed: %v", err)
	}

	if addResult.ID != 3 {
		t.Errorf("Expected next ID to be 3 after restart, got %d", addResult.ID)
	}

	// Verify persistence of the new task
	finalTasks, err := readTasks()
	if err != nil {
		t.Fatalf("readTasks failed after final add: %v", err)
	}
	if len(finalTasks) != 3 || finalTasks[2].Description != "Task Three" {
		t.Errorf("Expected 3 tasks with 'Task Three' as the last, got %+v", finalTasks)
	}
}
