package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func setupTest(t *testing.T) func() {
	// Store original values to restore them later
	oldTasksFile := tasksFile

	// Set tasksFile to a unique name for this test run
	tasksFile = "test_tasks_" + t.Name() + ".json"

	// Reinitialize currentTasks for a clean slate for each test
	currentTasks = &tasks{}
	currentTasks.List = []Task{}

	// Clean up any existing test file
	os.Remove(tasksFile)

	return func() {
		// Clean up after the test
		os.Remove(tasksFile)

		// Restore original tasksFile name
		tasksFile = oldTasksFile

		// Ensure currentTasks is completely reset at the end too
		currentTasks = &tasks{}
		currentTasks.List = []Task{}
	}
}

func captureOutput(f func()) string {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer w.Close()

	f()

	out, _ := ioutil.ReadAll(r)
	os.Stdout = oldStdout
	return string(out)
}

func TestAddCommand(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	// Test adding a single task
	output := captureOutput(func() {
		exitCode := run([]string{"todo", "add", "Buy milk"})
		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", exitCode)
		}
	})

	expected := `{"id":1,"status":"added"}` + "\n"
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}

	// Verify task is in storage
	loadTasks()
	if len(currentTasks.List) != 1 || currentTasks.List[0].Description != "Buy milk" || currentTasks.List[0].ID != 1 || currentTasks.List[0].Done != false {
		t.Errorf("Task not added correctly to storage: %+v", currentTasks.List)
	}

	// Test adding another task
	output = captureOutput(func() {
		exitCode := run([]string{"todo", "add", "Go for a run"})
		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", exitCode)
		}
	})

	expected = `{"id":2,"status":"added"}` + "\n"
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}

	loadTasks()
	if len(currentTasks.List) != 2 || currentTasks.List[1].ID != 2 || currentTasks.List[1].Description != "Go for a run" || currentTasks.List[1].Done != false {
		t.Errorf("Second task not added correctly: %+v", currentTasks.List)
	}
}

func TestListCommand(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	// Add some tasks first
	_ = run([]string{"todo", "add", "Task One"})
	_ = run([]string{"todo", "add", "Task Two"})

	output := captureOutput(func() {
		exitCode := run([]string{"todo", "list"})
		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", exitCode)
		}
	})

	expectedTasks := []Task{
		{ID: 1, Description: "Task One", Done: false},
		{ID: 2, Description: "Task Two", Done: false},
	}
	expected, _ := json.Marshal(expectedTasks)
	expected = append(expected, '\n')

	if output != string(expected) {
		t.Errorf("Expected %q, got %q", string(expected), output)
	}

	// Mark one task as done and list again
	_ = run([]string{"todo", "done", "1"})
	output = captureOutput(func() {
		exitCode := run([]string{"todo", "list"})
		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", exitCode)
		}
	})

	expectedTasks = []Task{
		{ID: 2, Description: "Task Two", Done: false},
	}
	expected, _ = json.Marshal(expectedTasks)
	expected = append(expected, '\n')

	if output != string(expected) {
		t.Errorf("Expected %q, got %q", string(expected), output)
	}
}

func TestDoneCommand(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	// Add a task
	_ = run([]string{"todo", "add", "Finish project"})

	// Mark it as done
	output := captureOutput(func() {
		exitCode := run([]string{"todo", "done", "1"})
		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", exitCode)
		}
	})

	expected := `{"id":1,"status":"completed"}` + "\n"
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}

	// Verify in storage
	loadTasks()
	if len(currentTasks.List) != 1 || currentTasks.List[0].Done != true {
		t.Errorf("Task not marked as done correctly: %+v", currentTasks.List)
	}

	// Test marking a non-existent task as done
	var exitCode int
	output = captureOutput(func() {
		exitCode = run([]string{"todo", "done", "99"})
	})

	if exitCode != 42 {
		t.Errorf("Expected exit code 42 for non-existent task, got %d", exitCode)
	}
	expected = `{"error":"task with ID 99 not found"}` + "\n"
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}

func TestErrorHandling(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	// Test invalid command
	var exitCode int
	output := captureOutput(func() {
		exitCode = run([]string{"todo", "unknown"})
	})

	if exitCode != 42 {
		t.Errorf("Expected exit code 42 for unknown command, got %d", exitCode)
	}
	expected := `{"error":"unknown command: unknown"}` + "\n"
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}

	// Test add command with no description
	output = captureOutput(func() {
		exitCode = run([]string{"todo", "add"})
	})

	if exitCode != 42 {
		t.Errorf("Expected exit code 42 for add without description, got %d", exitCode)
	}
	expected = `{"error":"invalid arguments: missing task description"}` + "\n"
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}

	// Test done command with no ID
	output = captureOutput(func() {
		exitCode = run([]string{"todo", "done"})
	})

	if exitCode != 42 {
		t.Errorf("Expected exit code 42 for done without ID, got %d", exitCode)
	}
	expected = `{"error":"invalid arguments: missing task ID"}` + "\n"
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}

	// Test done command with invalid ID
	output = captureOutput(func() {
		exitCode = run([]string{"todo", "done", "abc"})
	})

	if exitCode != 42 {
		t.Errorf("Expected exit code 42 for done with invalid ID, got %d", exitCode)
	}
	expected = `{"error":"invalid task ID: abc"}` + "\n"
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}

func TestLoadSaveTasks(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	// Add some tasks manually to currentTasks and save
	currentTasks.List = []Task{
		{ID: 1, Description: "Task A", Done: false},
		{ID: 2, Description: "Task B", Done: true},
	}
	if err := saveTasks(); err != nil {
		t.Fatalf("Failed to save tasks: %v", err)
	}

	// Clear currentTasks and load from file
	currentTasks.Lock()
	currentTasks.List = []Task{}
	currentTasks.Unlock()

	if err := loadTasks(); err != nil {
		t.Fatalf("Failed to load tasks: %v", err)
	}

	if len(currentTasks.List) != 2 || currentTasks.List[0].ID != 1 || currentTasks.List[1].Description != "Task B" || currentTasks.List[1].Done != true {
		t.Errorf("Loaded tasks do not match saved tasks: %+v", currentTasks.List)
	}

	// Test loading from an empty file
	os.Remove(tasksFile)
	if err := ioutil.WriteFile(tasksFile, []byte{}, 0644); err != nil {
		t.Fatalf("Failed to create empty tasks file: %v", err)
	}
	currentTasks.Lock()
	currentTasks.List = []Task{{ID: 1, Description: "dummy", Done: false}}
	currentTasks.Unlock()
	if err := loadTasks(); err != nil {
		t.Fatalf("Failed to load from empty file: %v", err)
	}
	if len(currentTasks.List) != 0 {
		t.Errorf("Expected empty list for empty file, got %+v", currentTasks.List)
	}
}
