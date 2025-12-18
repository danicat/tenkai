package todo

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

// tempDir is a temporary directory for test files.
var tempDir string

func setupTest(t *testing.T) {
	// Create a temporary directory for tasks.json
	var err error
	tempDir, err = ioutil.TempDir("", "todo_test_")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	// Change the working directory to the temporary directory
	// so that tasksFile ("tasks.json") is created there.
	currentDir, _ := os.Getwd()
	if err := os.Chdir(tempDir);
		err != nil {
			t.Fatalf("failed to change to temp dir: %v", err)
		}
	
	// Reset global state for each test
	tasks = []Task{}
	nextID = 1

	// Ensure cleanup even if the test panics.
	t.Cleanup(func() {
		// Change back to original directory
		_ = os.Chdir(currentDir)
		// Clean up the temporary directory
		_ = os.RemoveAll(tempDir)
		// Restore original tasksFile path (if it was changed for testing)
		tasksFile = filepath.Join(currentDir, "tasks.json") // Assuming tasks.json should be in the original directory for production
	})
	
	// This is crucial: update tasksFile to point to the temporary directory
	tasksFile = filepath.Join(tempDir, "tasks.json")

}

func teardownTest(t *testing.T) {
	// No need to explicitly remove tasksFile or tempDir here,
	// as it's handled by t.Cleanup in setupTest.
}

func TestAddTask(t *testing.T) {
	setupTest(t)
	defer teardownTest(t)

	desc1 := "Buy groceries"
	task1, err := AddTask(desc1)
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}

	if task1.ID != 1 {
		t.Errorf("Expected ID 1, got %d", task1.ID)
	}
	if task1.Description != desc1 {
		t.Errorf("Expected description %q, got %q", desc1, task1.Description)
	}
	if task1.Done != false {
		t.Errorf("Expected Done false, got %t", task1.Done)
	}

	// Verify persistence by reloading
	tasks = []Task{} // Clear in-memory tasks to simulate new app start
	nextID = 1      // Reset nextID
	if err := loadTasks(); err != nil {
		t.Fatalf("loadTasks after AddTask failed: %v", err)
	}

	if len(tasks) != 1 {
		t.Fatalf("Expected 1 task after reload, got %d", len(tasks))
	}
	if tasks[0].ID != 1 || tasks[0].Description != desc1 || tasks[0].Done != false {
		t.Errorf("Reloaded task mismatch: %+v", tasks[0])
	}
	if nextID != 2 {
		t.Errorf("Expected nextID 2, got %d", nextID)
	}

	desc2 := "Pay bills"
	task2, err := AddTask(desc2)
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}

	if task2.ID != 2 {
		t.Errorf("Expected ID 2, got %d", task2.ID)
	}
	// Verify persistence again for second task
	tasks = []Task{}
	nextID = 1
	if err := loadTasks(); err != nil {
		t.Fatalf("loadTasks after second AddTask failed: %v", err)
	}
	if len(tasks) != 2 {
		t.Fatalf("Expected 2 tasks after reload, got %d", len(tasks))
	}
	if nextID != 3 {
		t.Errorf("Expected nextID 3, got %d", nextID)
	}

}

func TestListTasks(t *testing.T) {
	setupTest(t)
	defer teardownTest(t)

	// No tasks initially
	initialTasks, err := ListTasks(false)
	if err != nil {
		t.Fatalf("ListTasks failed: %v", err)
	}
	if len(initialTasks) != 0 {
		t.Errorf("Expected 0 tasks, got %d", len(initialTasks))
	}

	_, _ = AddTask("Task 1")
	_, _ = AddTask("Task 2")
	completedTask, _ := AddTask("Task 3")
	_, _ = CompleteTask(completedTask.ID)

	// List pending tasks
	pendingTasks, err := ListTasks(false)
	if err != nil {
		t.Fatalf("ListTasks (pending) failed: %v", err)
	}
	if len(pendingTasks) != 2 {
		t.Errorf("Expected 2 pending tasks, got %d", len(pendingTasks))
	}
	if pendingTasks[0].ID != 1 || pendingTasks[1].ID != 2 {
		t.Errorf("Expected tasks with IDs 1 and 2, got %+v", pendingTasks)
	}

	// List all tasks
	allTasks, err := ListTasks(true)
	if err != nil {
		t.Fatalf("ListTasks (all) failed: %v", err)
	}
	if len(allTasks) != 3 {
		t.Errorf("Expected 3 total tasks, got %d", len(allTasks))
	}
	if allTasks[2].ID != 3 || allTasks[2].Done != true {
		t.Errorf("Expected task 3 to be done, got %+v", allTasks[2])
	}
}

func TestCompleteTask(t *testing.T) {
	setupTest(t)
	defer teardownTest(t)

	task1, _ := AddTask("Task to complete")

	// Complete an existing task
	completed, err := CompleteTask(task1.ID)
	if err != nil {
		t.Fatalf("CompleteTask failed: %v", err)
	}
	if completed.ID != task1.ID || completed.Done != true {
		t.Errorf("Expected task %d to be done, got %+v", task1.ID, completed)
	}

	// Verify persistence
	tasks = []Task{}
	nextID = 1
	if err := loadTasks(); err != nil {
		t.Fatalf("loadTasks after CompleteTask failed: %v", err)
	}
	if len(tasks) != 1 || tasks[0].Done != true {
		t.Errorf("Persisted task not correctly marked as done: %+v", tasks[0])
	}

	// Try to complete a non-existent task
	_, err = CompleteTask(999)
	if err == nil || err.Error() != "task with ID 999 not found" {
		t.Errorf("Expected 'task not found' error, got %v", err)
	}

	// Try to complete an already completed task
	_, err = CompleteTask(task1.ID)
	if err == nil || err.Error() != fmt.Sprintf("task %d is already completed", task1.ID) {
		t.Errorf("Expected 'already completed' error, got %v", err)
	}

	// Ensure nextID is updated correctly after completion/load
	if nextID != 2 {
		t.Errorf("Expected nextID 2 after completion/load, got %d", nextID)
	}
}

func TestLoadSaveEmptyFile(t *testing.T) {
	setupTest(t)
	defer teardownTest(t)

	// Ensure tasks.json is created if it doesn't exist and empty tasks are loaded
	if err := loadTasks(); err != nil {
		t.Fatalf("loadTasks on empty file failed: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks, got %d", len(tasks))
	}
	if nextID != 1 {
		t.Errorf("Expected nextID 1 for empty tasks, got %d", nextID)
	}

	// Add a task and save
	_, err := AddTask("First task")
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}

	// Simulate clear and reload
	tasks = []Task{}
	nextID = 1
	if err := loadTasks(); err != nil {
		t.Fatalf("loadTasks after adding task failed: %v", err)
	}

	if len(tasks) != 1 || tasks[0].ID != 1 {
		t.Errorf("Expected 1 task with ID 1, got %+v", tasks)
	}
}

func TestConcurrentAccess(t *testing.T) {
	setupTest(t)
	defer teardownTest(t)

	// Test concurrent AddTask operations
	var wg sync.WaitGroup
	numGoroutines := 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			description := fmt.Sprintf("Concurrent Task %d", i)
			_, err := AddTask(description)
			if err != nil {
				t.Errorf("AddTask failed in goroutine %d: %v", i, err)
			}
		}(i)
	}

	wg.Wait()

	// Reload tasks to ensure all are persisted and IDs are correct.
	tasks = []Task{}
	nextID = 1
	if err := loadTasks(); err != nil {
		t.Fatalf("loadTasks after concurrent AddTask failed: %v", err)
	}

	if len(tasks) != numGoroutines {
		t.Errorf("Expected %d tasks after concurrent adds, got %d", numGoroutines, len(tasks))
	}

	idMap := make(map[int]bool)
	for _, task := range tasks {
		if idMap[task.ID] {
			t.Errorf("Duplicate task ID found: %d", task.ID)
		}
		idMap[task.ID] = true
		if task.ID < 1 || task.ID > numGoroutines {
			t.Errorf("Task ID %d out of expected range", task.ID)
		}
	}
}
