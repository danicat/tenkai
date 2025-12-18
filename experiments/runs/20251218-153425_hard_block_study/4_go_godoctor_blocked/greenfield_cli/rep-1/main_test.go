package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

// save and restore original os.Args
var originalArgs = os.Args

func setArgs(args []string) {
	os.Args = append([]string{"todo"}, args...)
}

func restoreArgs() {
	os.Args = originalArgs
}

// captureOutput captures stdout, stderr, and the exit code of a function that calls osExit.
func captureOutput(f func()) (string, string, int) {
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	var stdoutBuf, stderrBuf bytes.Buffer
	outReader, outWriter, _ := os.Pipe()
	errReader, errWriter, _ := os.Pipe()

	os.Stdout = outWriter
	os.Stderr = errWriter

	var exitCode int
	outputDone := make(chan struct{})    // Signifies that output has been fully read
	executionDone := make(chan struct{}) // Signifies that the function f() has finished executing

	// Goroutine to read from the pipes
	go func() {
		defer close(outputDone)
		outWriter.Close() // Close the write-end of the pipe in the main goroutine
		errWriter.Close() // Close the write-end of the pipe in the main goroutine

		_, _ = stdoutBuf.ReadFrom(outReader)
		_, _ = stderrBuf.ReadFrom(errReader)
		outReader.Close()
		errReader.Close()
	}()

	// Override osExit temporarily to capture the exit code
	originalOsExit := osExit
	osExit = func(code int) {
		exitCode = code
		panic("os.Exit was called") // Panic to stop the execution flow in the current goroutine
	}

	// Execute the function that might call osExit
	func() {
		defer func() {
			if r := recover(); r != nil && r.(string) != "os.Exit was called" {
				panic(r) // Re-panic if it's not our expected exit panic
			}
			close(executionDone)
		}()
		f()
	}()

	// Wait for the function f() to finish executing
	<-executionDone

	// Restore original os.Exit
	osExit = originalOsExit

	// Restore original stdout/stderr
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	// Wait for the output capture goroutine to finish reading
	<-outputDone

	return strings.TrimSpace(stdoutBuf.String()), strings.TrimSpace(stderrBuf.String()), exitCode
}

func TestMain(m *testing.M) {
	// TestMain does not need to override os.Exit here anymore.
	// The captureOutput function handles it for individual tests.

	exitVal := m.Run()
	os.Exit(exitVal)
}

func setupTestEnv() func() {
	tasksLock.Lock()
	defer tasksLock.Unlock()

	// Create a temporary file for tasks.json
	tmpFile, err := ioutil.TempFile("", "tasks_test_*.json")
	if err != nil {
		panic(err)
	}
	_ = tmpFile.Close()

	oldTasksFilePath := tasksFilePath
	tasksFilePath = tmpFile.Name()

	// Clear existing tasks and reset nextID for the test
	tasks = []Task{}
	nextID = 1

	// Ensure an empty tasks.json exists for tests that expect it
	_ = saveTasks(tasks)

	return func() {
		// Cleanup: remove the temporary file and restore original tasksFilePath
		_ = os.Remove(tasksFilePath)
		tasksFilePath = oldTasksFilePath
		tasks = []Task{}
		nextID = 1
	}
}

func Test_loadAndSaveTasks(t *testing.T) {
	cleanup := setupTestEnv()
	defer cleanup()

	expectedTasks := []Task{
		{ID: 1, Description: "Test Task 1", Done: false},
		{ID: 2, Description: "Test Task 2", Done: true},
	}

	err := saveTasks(expectedTasks)
	if err != nil {
		t.Fatalf("Failed to save tasks: %v", err)
	}

	loadedTasks, err := loadTasks()
	if err != nil {
		t.Fatalf("Failed to load tasks: %v", err)
	}

	if len(loadedTasks) != len(expectedTasks) {
		t.Fatalf("Loaded tasks count mismatch. Expected %d, Got %d", len(expectedTasks), len(loadedTasks))
	}

	for i, task := range loadedTasks {
		if task != expectedTasks[i] {
			t.Errorf("Task mismatch at index %d. Expected %+v, Got %+v", i, expectedTasks[i], task)
		}
	}
}

func Test_parseTaskID(t *testing.T) {
	tests := []struct {
		input     string
		expected  int
		expectErr bool
	}{
		{"1", 1, false},
		{"0", 0, true},
		{"-5", 0, true},
		{"abc", 0, true},
		{"+10", 10, false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			actual, err := parseTaskID(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("parseTaskID(%q) error status mismatch. Expected error: %v, Got: %v", tt.input, tt.expectErr, err != nil)
			}
			if actual != tt.expected {
				t.Errorf("parseTaskID(%q) = %d, expected %d", tt.input, actual, tt.expected)
			}
		})
	}
}

func Test_handleAddCommand(t *testing.T) {
	cleanup := setupTestEnv()
	defer cleanup()

	// Test adding a task successfully
	stdout, _, exitCode := captureOutput(func() {
		setArgs([]string{"add", "Buy groceries"})
		handleAddCommand(os.Args[2:])
	})
	restoreArgs()

	if exitCode != 0 {
		t.Errorf("handleAddCommand exited with code %d, expected 0", exitCode)
	}

	expectedJSON := `{"id":1,"status":"added"}`
	if stdout != expectedJSON {
		t.Errorf("handleAddCommand output mismatch. Expected %s, Got %s", expectedJSON, stdout)
	}

	loadedTasks, _ := loadTasks()
	if len(loadedTasks) != 1 || loadedTasks[0].Description != "Buy groceries" || loadedTasks[0].ID != 1 {
		t.Errorf("Task not added correctly: %+v", loadedTasks)
	}

	// Test adding another task
	stdout, _, exitCode = captureOutput(func() {
		setArgs([]string{"add", "Pay bills"})
		handleAddCommand(os.Args[2:])
	})
	restoreArgs()

	if exitCode != 0 {
		t.Errorf("handleAddCommand exited with code %d, expected 0", exitCode)
	}

	expectedJSON = `{"id":2,"status":"added"}`
	if stdout != expectedJSON {
		t.Errorf("handleAddCommand output mismatch. Expected %s, Got %s", expectedJSON, stdout)
	}

	loadedTasks, _ = loadTasks()
	if len(loadedTasks) != 2 || loadedTasks[1].Description != "Pay bills" || loadedTasks[1].ID != 2 {
		t.Errorf("Second task not added correctly: %+v", loadedTasks)
	}

	// Test missing description
	stdout, _, exitCode = captureOutput(func() {
		setArgs([]string{"add"})
		handleAddCommand(os.Args[2:])
	})
	restoreArgs()

	if exitCode != 42 {
		t.Errorf("handleAddCommand with missing description exited with code %d, expected 42", exitCode)
	}
	expectedError := `{"error":"missing task description"}`
	if stdout != expectedError {
		t.Errorf("handleAddCommand error output mismatch. Expected %s, Got %s", expectedError, stdout)
	}
}

func Test_handleListCommand(t *testing.T) {
	cleanup := setupTestEnv()
	defer cleanup()

	// Add some tasks for listing
	saveTasks([]Task{
		{ID: 1, Description: "Task 1", Done: false},
		{ID: 2, Description: "Task 2", Done: true},
		{ID: 3, Description: "Task 3", Done: false},
	})

	// Manually set package-level tasks and nextID for the test to ensure consistency
	tasksLock.Lock()
	tasks = []Task{
		{ID: 1, Description: "Task 1", Done: false},
		{ID: 2, Description: "Task 2", Done: true},
		{ID: 3, Description: "Task 3", Done: false},
	}
	nextID = 4
	tasksLock.Unlock()

	// Test listing pending tasks
	stdout, _, exitCode := captureOutput(func() {
		setArgs([]string{"list"})
		handleListCommand()
	})
	restoreArgs()

	if exitCode != 0 {
		t.Errorf("handleListCommand exited with code %d, expected 0", exitCode)
	}

	expectedJSON := `[{"id":1,"description":"Task 1","done":false},{"id":3,"description":"Task 3","done":false}]`
	if stdout != expectedJSON {
		t.Errorf("handleListCommand output mismatch. Expected %s, Got %s", expectedJSON, stdout)
	}

	// Test listing with no tasks
	_ = os.Remove(tasksFilePath) // Use tasksFilePath directly
	tasksLock.Lock()
	tasks = []Task{}
	nextID = 1
	tasksLock.Unlock()

	stdout, _, exitCode = captureOutput(func() {
		setArgs([]string{"list"})
		handleListCommand()
	})
	restoreArgs()

	if exitCode != 0 {
		t.Errorf("handleListCommand (empty) exited with code %d, expected 0", exitCode)
	}
	expectedEmptyJSON := `[]`
	if stdout != expectedEmptyJSON {
		t.Errorf("handleListCommand (empty) output mismatch. Expected %s, Got %s", expectedEmptyJSON, stdout)
	}
}

func Test_handleDoneCommand(t *testing.T) {
	cleanup := setupTestEnv()
	defer cleanup()

	// Add some tasks initially
	saveTasks([]Task{
		{ID: 1, Description: "Task 1", Done: false},
		{ID: 2, Description: "Task 2", Done: false},
		{ID: 3, Description: "Task 3", Done: false},
	})

	tasksLock.Lock()
	tasks = []Task{
		{ID: 1, Description: "Task 1", Done: false},
		{ID: 2, Description: "Task 2", Done: false},
		{ID: 3, Description: "Task 3", Done: false},
	}
	nextID = 4
	tasksLock.Unlock()

	// Test marking a task as done successfully
	stdout, _, exitCode := captureOutput(func() {
		setArgs([]string{"done", "2"})
		handleDoneCommand(os.Args[2:])
	})
	restoreArgs()

	if exitCode != 0 {
		t.Errorf("handleDoneCommand exited with code %d, expected 0", exitCode)
	}
	expectedJSON := `{"id":2,"status":"completed"}`
	if stdout != expectedJSON {
		t.Errorf("handleDoneCommand output mismatch. Expected %s, Got %s", expectedJSON, stdout)
	}

	loadedTasks, _ := loadTasks()
	if !loadedTasks[1].Done {
		t.Errorf("Task 2 was not marked as done: %+v", loadedTasks[1])
	}

	// Test marking a non-existent task
	stdout, _, exitCode = captureOutput(func() {
		setArgs([]string{"done", "99"})
		handleDoneCommand(os.Args[2:])
	})
	restoreArgs()

	if exitCode != 42 {
		t.Errorf("handleDoneCommand (non-existent) exited with code %d, expected 42", exitCode)
	}
	expectedError := `{"error":"task with ID 99 not found"}`
	if stdout != expectedError {
		t.Errorf("handleDoneCommand (non-existent) error output mismatch. Expected %s, Got %s", expectedError, stdout)
	}

	// Test invalid ID format
	stdout, _, exitCode = captureOutput(func() {
		setArgs([]string{"done", "abc"})
		handleDoneCommand(os.Args[2:])
	})
	restoreArgs()

	if exitCode != 42 {
		t.Errorf("handleDoneCommand (invalid ID) exited with code %d, expected 42", exitCode)
	}
	expectedError = `{"error":"invalid task ID: \"abc\" is not a valid integer"}`
	if stdout != expectedError {
		t.Errorf("handleDoneCommand (invalid ID) error output mismatch. Expected %s, Got %s", expectedError, stdout)
	}

	// Test missing ID
	stdout, _, exitCode = captureOutput(func() {
		setArgs([]string{"done"})
		handleDoneCommand(os.Args[2:])
	})
	restoreArgs()

	if exitCode != 42 {
		t.Errorf("handleDoneCommand (missing ID) exited with code %d, expected 42", exitCode)
	}
	expectedError = `{"error":"missing task ID"}`
	if stdout != expectedError {
		t.Errorf("handleDoneCommand (missing ID) error output mismatch. Expected %s, Got %s", expectedError, stdout)
	}
}

func Test_MainFunctionDispatch(t *testing.T) {
	cleanup := setupTestEnv()
	defer cleanup()

	// Test unknown command
	stdout, _, exitCode := captureOutput(func() {
		setArgs([]string{"unknown"})
		main()
	})
	restoreArgs()

	if exitCode != 42 {
		t.Errorf("main (unknown command) exited with code %d, expected 42", exitCode)
	}
	expectedError := `{"error":"unknown command: unknown"}`
	if stdout != expectedError {
		t.Errorf("main (unknown command) output mismatch. Expected %s, Got %s", expectedError, stdout)
	}

	// Test missing command
	stdout, _, exitCode = captureOutput(func() {
		setArgs([]string{})
		main()
	})
	restoreArgs()

	if exitCode != 42 {
		t.Errorf("main (missing command) exited with code %d, expected 42", exitCode)
	}
	expectedError = `{"error":"missing command"}`
	if stdout != expectedError {
		t.Errorf("main (missing command) output mismatch. Expected %s, Got %s", expectedError, stdout)
	}
}
