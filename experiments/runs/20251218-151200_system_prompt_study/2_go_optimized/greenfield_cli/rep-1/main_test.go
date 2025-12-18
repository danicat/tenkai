package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"testing"
)

// Mocks for os.Exit and os.Stdout to capture output and prevent premature exit.
var (
	exitCode       int
	mu             sync.Mutex // guards exitCode and outBuffer
	oldStdout      *os.File
	oldStderr      *os.File
	capturedStdout bytes.Buffer
)

// mockExit is a test helper to capture the exit code.
func mockExit(code int) {
	mu.Lock()
	exitCode = code
	mu.Unlock()
	panic("os.Exit was called with code: " + strconv.Itoa(code))
}

// captureOutput redirects stdout and stderr to a buffer.
func captureOutput() {
	mu.Lock()
	oldStdout = os.Stdout
	oldStderr = os.Stderr
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w
	capturedStdout.Reset()
	go func() {
		defer r.Close()
		_, _ = io.Copy(&capturedStdout, r)
	}()
	mu.Unlock()
}

// stopCaptureOutput restores stdout and returns the captured output.
func stopCaptureOutput() string {
	mu.Lock()
	defer mu.Unlock()
	// Close the writer to signal the goroutine to finish copying
	_ = os.Stdout.Close()
	_ = os.Stderr.Close()

	// Restore original stdout/stderr
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	return capturedStdout.String()
}

// TestMain sets up and tears down test environment.
func TestMain(m *testing.M) {
	// Save original os.Exit
	origOsExit := osExit

	// Replace os.Exit with mockExit during tests
	osExit = mockExit

	// Run tests
	code := m.Run()

	// Restore original os.Exit
	osExit = origOsExit

	os.Exit(code)
}

// setupTest creates a temporary tasks.json for testing and returns its path.
func setupTest(t *testing.T) string {
	tmpDir := t.TempDir()
	return filepath.Join(tmpDir, "tasks.json")
}

func TestAddTask(t *testing.T) {
	tasksFilePath := setupTest(t)

	description := "Buy groceries"

	defer func() {
		if r := recover(); r != nil {
			if s, ok := r.(string); ok && (s == "os.Exit was called with code: 0" || s == "os.Exit was called with code: 1") {
				// Expected panic from mockExit
			} else {
				t.Fatalf("unexpected panic: %v", r)
			}
		}
	}()

	captureOutput()
	addTask(description)
	output := stopCaptureOutput()

	expectedJSON := `{"id":1,"status":"added"}` + "\n"
	if output != expectedJSON {
		t.Errorf("addTask(%q) got %q, want %q", description, output, expectedJSON)
	}

	tasks, err := loadTasks(tasksFilePath)
	if err != nil {
		t.Fatalf("failed to load tasks after add: %v", err)
	}

	if len(tasks) != 1 || tasks[0].Description != description || tasks[0].ID != 1 || tasks[0].Done != false {
		t.Errorf("task not added correctly: %+v", tasks)
	}
}

func TestListTasks(t *testing.T) {
	tasksFilePath := setupTest(t)

	// Add some tasks first
	_ = saveTasks(tasksFilePath, []Task{
		{ID: 1, Description: "Task 1", Done: false},
		{ID: 2, Description: "Task 2", Done: true},
		{ID: 3, Description: "Task 3", Done: false},
	})

	defer func() {
		if r := recover(); r != nil {
			if s, ok := r.(string); ok && s == "os.Exit was called with code: 0" {
				// Expected panic from mockExit
			} else {
				t.Fatalf("unexpected panic: %v", r)
			}
		}
	}()

	captureOutput()
	listTasks()
	output := stopCaptureOutput()

	expectedTasks := []Task{
		{ID: 1, Description: "Task 1", Done: false},
		{ID: 3, Description: "Task 3", Done: false},
	}
	expectedJSONBytes, _ := json.Marshal(expectedTasks)
	expectedJSON := string(expectedJSONBytes) + "\n"

	if output != expectedJSON {
		t.Errorf("listTasks() got %q, want %q", output, expectedJSON)
	}
}

func TestCompleteTask(t *testing.T) {
	tasksFilePath := setupTest(t)

	// Add some tasks first
	err := saveTasks(tasksFilePath, []Task{
		{ID: 1, Description: "Task 1", Done: false},
		{ID: 2, Description: "Task 2", Done: false},
	})
	if err != nil {
		t.Fatalf("failed to save initial tasks: %v", err)
	}

	defer func() {
		if r := recover(); r != nil {
			if s, ok := r.(string); ok && s == "os.Exit was called with code: 0" {
				// Expected panic from mockExit
			} else {
				t.Fatalf("unexpected panic: %v", r)
			}
		}
	}()

	captureOutput()
	completeTask(2)
	output := stopCaptureOutput()

	expectedJSON := `{"id":2,"status":"completed"}` + "\n"
	if output != expectedJSON {
		t.Errorf("completeTask(2) got %q, want %q", output, expectedJSON)
	}

	tasks, err := loadTasks(tasksFilePath)
	if err != nil {
		t.Fatalf("failed to load tasks after complete: %v", err)
	}

	if len(tasks) != 2 || tasks[0].Done != false || tasks[1].Done != true {
		t.Errorf("task not completed correctly: %+v", tasks)
	}
}

func TestCompleteTask_NotFound(t *testing.T) {
	tasksFilePath := setupTest(t)

	// Add some tasks first
	_ = saveTasks(tasksFilePath, []Task{
		{ID: 1, Description: "Task 1", Done: false},
	})

	defer func() {
		if r := recover(); r != nil {
			if s, ok := r.(string); ok && s == "os.Exit was called with code: 1" {
				// Expected panic from mockExit due to error exit
			} else {
				t.Fatalf("unexpected panic: %v", r)
			}
		}
	}()

	captureOutput()
	completeTask(99)
	output := stopCaptureOutput()

	expectedJSON := `{"error":"task with ID 99 not found"}` + "\n"
	if output != expectedJSON {
		t.Errorf("completeTask(99) got %q, want %q", output, expectedJSON)
	}
}

func TestMain_InvalidCommand(t *testing.T) {
	tasksFilePath := setupTest(t)
	oldTasksFile := tasksFile
	defer func() { tasksFile = oldTasksFile }()
	tasksFile = tasksFilePath

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"todo", "invalidcommand"}

	defer func() {
		if r := recover(); r != nil {
			if s, ok := r.(string); ok && s == "os.Exit was called with code: 42" {
				// Expected panic from mockExit with code 42
			} else {
				t.Fatalf("unexpected panic: %v", r)
			}
		} else {
			t.Fatalf("main did not panic with os.Exit")
		}
	}()

	captureOutput()
	main()
	output := stopCaptureOutput()

	expectedJSON := `{"error":"unknown command"}` + "\n"
	if output != expectedJSON {
		t.Errorf("main with invalid command got %q, want %q", output, expectedJSON)
	}
}

func TestMain_AddMissingDescription(t *testing.T) {
	tasksFilePath := setupTest(t)
	oldTasksFile := tasksFile
	defer func() { tasksFile = oldTasksFile }()
	tasksFile = tasksFilePath

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"todo", "add"}

	defer func() {
		if r := recover(); r != nil {
			if s, ok := r.(string); ok && s == "os.Exit was called with code: 42" {
				// Expected panic from mockExit with code 42
			} else {
				t.Fatalf("unexpected panic: %v", r)
			}
		} else {
			t.Fatalf("main did not panic with os.Exit")
		}
	}()

	captureOutput()
	main()
	output := stopCaptureOutput()

	expectedJSON := `{"error":"missing task description"}` + "\n"
	if output != expectedJSON {
		t.Errorf("main with missing description got %q, want %q", output, expectedJSON)
	}
}

func TestMain_DoneMissingID(t *testing.T) {
	tasksFilePath := setupTest(t)
	oldTasksFile := tasksFile
	defer func() { tasksFile = oldTasksFile }()
	tasksFile = tasksFilePath

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"todo", "done"}

	defer func() {
		if r := recover(); r != nil {
			if s, ok := r.(string); ok && s == "os.Exit was called with code: 42" {
				// Expected panic from mockExit with code 42
			} else {
				t.Fatalf("unexpected panic: %v", r)
			}
		} else {
			t.Fatalf("main did not panic with os.Exit")
		}
	}()

	captureOutput()
	main()
	output := stopCaptureOutput()

	expectedJSON := `{"error":"missing task ID"}` + "\n"
	if output != expectedJSON {
		t.Errorf("main with missing ID got %q, want %q", output, expectedJSON)
	}
}

func TestMain_DoneInvalidID(t *testing.T) {
	tasksFilePath := setupTest(t)
	oldTasksFile := tasksFile
	defer func() { tasksFile = oldTasksFile }()
	tasksFile = tasksFilePath

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"todo", "done", "abc"}

	defer func() {
		if r := recover(); r != nil {
			if s, ok := r.(string); ok && s == "os.Exit was called with code: 42" {
				// Expected panic from mockExit with code 42
			} else {
				t.Fatalf("unexpected panic: %v", r)
			}
		} else {
			t.Fatalf("main did not panic with os.Exit")
		}
	}()

	captureOutput()
	main()
	output := stopCaptureOutput()

	expectedJSON := `{"error":"invalid task ID"}` + "\n"
	if output != expectedJSON {
		t.Errorf("main with invalid ID got %q, want %q", output, expectedJSON)
	}
}

func TestLoadSaveTasks(t *testing.T) {
	tasksFilePath := setupTest(t)

	tasks := []Task{
		{ID: 1, Description: "Test Task 1", Done: false},
		{ID: 2, Description: "Test Task 2", Done: true},
	}

	err := saveTasks(tasksFilePath, tasks)
	if err != nil {
		t.Fatalf("failed to save tasks: %v", err)
	}

	loadedTasks, err := loadTasks(tasksFilePath)
	if err != nil {
		t.Fatalf("failed to load tasks: %v", err)
	}

	if len(loadedTasks) != len(tasks) {
		t.Fatalf("loaded tasks count mismatch: got %d, want %d", len(loadedTasks), len(tasks))
	}

	for i, task := range tasks {
		if loadedTasks[i].ID != task.ID || loadedTasks[i].Description != task.Description || loadedTasks[i].Done != task.Done {
			t.Errorf("loaded task mismatch at index %d: got %+v, want %+v", i, loadedTasks[i], task)
		}
	}
}

func TestLoadTasks_EmptyFile(t *testing.T) {
	tasksFilePath := setupTest(t)

	// Create an empty tasks.json file
	err := os.WriteFile(tasksFilePath, []byte{}, 0644)
	if err != nil {
		t.Fatalf("failed to create empty tasks file: %v", err)
	}

	// Test loading from an empty file
	loadedTasks, err := loadTasks(tasksFilePath)
	if err != nil {
		t.Fatalf("failed to load tasks from empty file: %v", err)
	}

	if len(loadedTasks) != 0 {
		t.Errorf("expected 0 tasks from empty file, got %d", len(loadedTasks))
	}
}

func TestLoadTasks_NonExistentFile(t *testing.T) {
	tasksFilePath := setupTest(t)

	// Ensure tasks.json does not exist
	_ = os.Remove(tasksFilePath) // Use _ to ignore error, as it might not exist

	loadedTasks, err := loadTasks(tasksFilePath)
	if err != nil {
		t.Fatalf("failed to load tasks from non-existent file: %v", err)
	}

	if len(loadedTasks) != 0 {
		t.Errorf("expected 0 tasks from non-existent file, got %d", len(loadedTasks))
	}
}

