package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var cliBinaryPath string

func TestMain(m *testing.M) {
	// Create a temporary directory for the binary
	tempDir, err := ioutil.TempDir("", "todo_cli_test_binary")
	if err != nil {
		fmt.Printf("Failed to create temp dir for binary: %v\n", err)
		os.Exit(1)
	}
	cliBinaryPath = filepath.Join(tempDir, "todo_cli")

	// Build the CLI binary
	cmd := exec.Command("go", "build", "-o", cliBinaryPath, ".")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to build CLI binary: %v\n", err)
		os.RemoveAll(tempDir)
		os.Exit(1)
	}

	// Run tests
	exitCode := m.Run()

	// Clean up
	os.RemoveAll(tempDir)
	os.Exit(exitCode)
}

// Helper to run the CLI command with specific arguments and return stdout/stderr
func runCLI(t *testing.T, args ...string) (string, string, int) {
	cmd := exec.Command(cliBinaryPath, args...)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			t.Fatalf("failed to run command: %v", err)
		}
	}

	return stdout.String(), stderr.String(), exitCode
}

// Helper to create a temporary tasks.json file for tests
func createTempTasksFile(t *testing.T, content string) string {
	tempDir, err := ioutil.TempDir("", "todo_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	tempFile := filepath.Join(tempDir, tasksFile)

	err = ioutil.WriteFile(tempFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to write temp tasks file: %v", err)
	}

	os.Setenv("TASKS_FILE", tempFile) // Set environment variable to point to temp file

	return tempDir
}

// TestAddTask tests the 'add' command
func TestAddTask(t *testing.T) {
	tempDir := createTempTasksFile(t, "[]")
	defer os.RemoveAll(tempDir)

	stdout, stderr, exitCode := runCLI(t, "add", "-text", "Buy groceries")
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d. Stderr: %s", exitCode, stderr)
	}

	expected := `{"id": 1, "status": "added"}
`
	if stdout != expected {
		t.Errorf("Expected stdout %q, got %q", expected, stdout)
	}

	// Verify task was added to the file
	var loadedTasks []Task
	data, err := ioutil.ReadFile(filepath.Join(tempDir, tasksFile))
	if err != nil {
		t.Fatalf("failed to read tasks file after add: %v", err)
	}
	err = json.Unmarshal(data, &loadedTasks)
	if err != nil {
		t.Fatalf("failed to unmarshal tasks after add: %v", err)
	}

	if len(loadedTasks) != 1 || loadedTasks[0].Description != "Buy groceries" || loadedTasks[0].ID != 1 || loadedTasks[0].Done != false {
		t.Errorf("Task not correctly saved: %+v", loadedTasks)
	}
}

// TestListTasks tests the 'list' command
func TestListTasks(t *testing.T) {
	initialTasks := `[{"id": 1, "description": "Buy milk", "done": false}, {"id": 2, "description": "Pay bills", "done": true}]`
	tempDir := createTempTasksFile(t, initialTasks)
	defer os.RemoveAll(tempDir)

	stdout, stderr, exitCode := runCLI(t, "list")
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d. Stderr: %s", exitCode, stderr)
	}

	expectedJSON := `[{"id":1,"description":"Buy milk","done":false}]`
	if strings.TrimSpace(stdout) != expectedJSON {
		t.Errorf("Expected stdout %q, got %q", expectedJSON, strings.TrimSpace(stdout))
	}
}

// TestCompleteTask tests the 'done' command
func TestCompleteTask(t *testing.T) {
	initialTasks := `[{"id": 1, "description": "Buy milk", "done": false}, {"id": 2, "description": "Pay bills", "done": false}]`
	tempDir := createTempTasksFile(t, initialTasks)
	defer os.RemoveAll(tempDir)

	stdout, stderr, exitCode := runCLI(t, "done", "-id", "1")
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d. Stderr: %s", exitCode, stderr)
	}

	expected := `{"id": 1, "status": "completed"}
`
	if stdout != expected {
		t.Errorf("Expected stdout %q, got %q", expected, stdout)
	}

	// Verify task was marked as done in the file
	var loadedTasks []Task
	data, err := ioutil.ReadFile(filepath.Join(tempDir, tasksFile))
	if err != nil {
		t.Fatalf("failed to read tasks file after done: %v", err)
	}
	err = json.Unmarshal(data, &loadedTasks)
	if err != nil {
		t.Fatalf("failed to unmarshal tasks after done: %v", err)
	}

	if len(loadedTasks) != 2 || loadedTasks[0].Done != true || loadedTasks[1].Done != false {
		t.Errorf("Task not correctly marked as done: %+v", loadedTasks)
	}
}

// TestInvalidArguments tests error handling for invalid arguments (exit code 42)
func TestInvalidArguments(t *testing.T) {
	tempDir := createTempTasksFile(t, "[]")
	defer os.RemoveAll(tempDir)

	// Test 'add' without text
	stdout, stderr, exitCode := runCLI(t, "add")
	if exitCode != 42 {
		t.Errorf("Expected exit code 42 for 'add' without text, got %d. Stderr: %s", exitCode, stderr)
	}
	expectedErr := `{"error": "invalid arguments"}
`
	if stdout != expectedErr {
		t.Errorf("Expected stdout %q, got %q for 'add' without text", expectedErr, stdout)
	}

	// Test 'done' without ID
	stdout, stderr, exitCode = runCLI(t, "done")
	if exitCode != 42 {
		t.Errorf("Expected exit code 42 for 'done' without ID, got %d. Stderr: %s", exitCode, stderr)
	}
	if stdout != expectedErr {
		t.Errorf("Expected stdout %q, got %q for 'done' without ID", expectedErr, stdout)
	}

	// Test unknown command
	stdout, stderr, exitCode = runCLI(t, "unknown")
	if exitCode != 42 {
		t.Errorf("Expected exit code 42 for unknown command, got %d. Stderr: %s", exitCode, stderr)
	}
	expectedUnknownCmdErr := `{"error": "expected 'add', 'list', or 'done' subcommands"}
`
	if stdout != expectedUnknownCmdErr {
		t.Errorf("Expected stdout %q, got %q for unknown command", expectedUnknownCmdErr, stdout)
	}
}

// TestTaskNotFound tests error handling when a task ID is not found (exit code 1)
func TestTaskNotFound(t *testing.T) {
	initialTasks := `[{"id": 1, "description": "Buy milk", "done": false}]`
	tempDir := createTempTasksFile(t, initialTasks)
	defer os.RemoveAll(tempDir)

	stdout, stderr, exitCode := runCLI(t, "done", "-id", "99")
	if exitCode != 1 {
		t.Errorf("Expected exit code 1 for task not found, got %d. Stderr: %s", exitCode, stderr)
	}

	expectedErr := `{"error": "task not found"}
`
	if stdout != expectedErr {
		t.Errorf("Expected stdout %q, got %q", expectedErr, stdout)
	}
}

// TestPersistence tests if tasks are loaded and saved correctly across runs
func TestPersistence(t *testing.T) {
	tempDir := createTempTasksFile(t, "[]")
	defer os.RemoveAll(tempDir)

	// Add a task in the first run
	stdout, _, exitCode := runCLI(t, "add", "-text", "Task One")
	if exitCode != 0 {
		t.Fatalf("First run 'add' failed: %s", stdout)
	}

	// Add another task in a separate run
	stdout, _, exitCode = runCLI(t, "add", "-text", "Task Two")
	if exitCode != 0 {
		t.Fatalf("Second run 'add' failed: %s", stdout)
	}

	// List tasks in a third run to check if both are present
	stdout, _, exitCode = runCLI(t, "list")
	if exitCode != 0 {
		t.Fatalf("Third run 'list' failed: %s", stdout)
	}

	expectedJSON := `[{"id":1,"description":"Task One","done":false},{"id":2,"description":"Task Two","done":false}]`
	if strings.TrimSpace(stdout) != expectedJSON {
		t.Errorf("Expected stdout %q, got %q after multiple runs", expectedJSON, strings.TrimSpace(stdout))
	}
}
