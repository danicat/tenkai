package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"
)

// mockExit replaces main.exit for testing purposes.
func mockExit(code int) {
	panic(strconv.Itoa(code))
}

// captureOutput captures stdout and stderr.
func captureOutput(f func()) (string, string) {
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	rPStdout, wPStdout, _ := os.Pipe()
	rPStderr, wPStderr, _ := os.Pipe()

	os.Stdout = wPStdout
	os.Stderr = wPStderr

	outChan := make(chan string)
	errChan := make(chan string)

	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, rPStdout)
		outChan <- buf.String()
	}()

	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, rPStderr)
		errChan <- buf.String()
	}()

	f()

	wPStdout.Close()
	wPStderr.Close()

	out := <-outChan
	err := <-errChan

	os.Stdout = oldStdout
	os.Stderr = oldStderr

	return out, err
}

// setupTest cleans up tasks.json and resets the tasks slice before each test.
func setupTest() {
	os.Remove(tasksFile) // Clean up any existing tasks.json
	tasks = []Task{}    // Reset in-memory tasks
}

func TestAddCommand(t *testing.T) {
	setupTest()
	defer func() { exit = os.Exit }() // Ensure os.Exit is restored
	exit = mockExit

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"todo", "add", "-d", "Buy groceries"}

	out, _ := captureOutput(main)

	expected := `{"id": 1, "status": "added"}`
	if !strings.Contains(out, expected) {
		t.Errorf("Add command output mismatch.\nExpected: %s\nGot: %s", expected, out)
	}

	// Verify task file content
	data, err := ioutil.ReadFile(tasksFile)
	if err != nil {
		t.Fatalf("Failed to read tasks file: %v", err)
	}

	var loadedTasks []Task
	if err := json.Unmarshal(data, &loadedTasks); err != nil {
		t.Fatalf("Failed to unmarshal tasks: %v", err)
	}

	if len(loadedTasks) != 1 || loadedTasks[0].ID != 1 || loadedTasks[0].Description != "Buy groceries" || loadedTasks[0].Done != false {
		t.Errorf("Task not correctly saved. Got: %+v", loadedTasks)
	}
}

func TestAddCommand_MissingDescription(t *testing.T) {
	setupTest()
	defer func() { exit = os.Exit }() // Ensure os.Exit is restored
	exit = mockExit

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"todo", "add"}

	exitCode := 0
	out, _ := captureOutput(func() {
		defer func() {
			if r := recover(); r != nil {
				exitCode, _ = strconv.Atoi(r.(string))
			}
		}()
		main()
	})

	expectedError := `{"error": "task description cannot be empty"}`
	if !strings.Contains(out, expectedError) {
		t.Errorf("Expected error for missing description not found.\nExpected: %s\nGot: %s", expectedError, out)
	}

	if exitCode != 42 {
		t.Errorf("Expected exit code 42 for missing description, got %d", exitCode)
	}
}

func TestListCommand(t *testing.T) {
	setupTest()
	defer func() { exit = os.Exit }() // Ensure os.Exit is restored
	exit = mockExit

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Add some tasks first
	os.Args = []string{"todo", "add", "-d", "Task 1"}
	captureOutput(main)
	os.Args = []string{"todo", "add", "-d", "Task 2"}
	captureOutput(main)

	// List tasks
	os.Args = []string{"todo", "list"}
	out, _ := captureOutput(main)

	expected := `[
  {
    "id": 1,
    "description": "Task 1",
    "done": false
  },
  {
    "id": 2,
    "description": "Task 2",
    "done": false
  }
]`

	if !strings.Contains(out, expected) {
		t.Errorf("List command output mismatch.\nExpected substring: %s\nGot: %s", expected, out)
	}
}

func TestCompleteCommand(t *testing.T) {
	setupTest()
	defer func() { exit = os.Exit }() // Ensure os.Exit is restored
	exit = mockExit

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Add a task
	os.Args = []string{"todo", "add", "-d", "Task to complete"}
	captureOutput(main)

	// Complete the task
	os.Args = []string{"todo", "done", "-id", "1"}
	out, _ := captureOutput(main)

	expected := `{"id": 1, "status": "completed"}`
	if !strings.Contains(out, expected) {
		t.Errorf("Complete command output mismatch.\nExpected: %s\nGot: %s", expected, out)
	}

	// Verify task file content
	data, err := ioutil.ReadFile(tasksFile)
	if err != nil {
		t.Fatalf("Failed to read tasks file: %v", err)
	}

	var loadedTasks []Task
	if err := json.Unmarshal(data, &loadedTasks); err != nil {
		t.Fatalf("Failed to unmarshal tasks: %v", err)
	}

	if len(loadedTasks) != 1 || loadedTasks[0].ID != 1 || loadedTasks[0].Description != "Task to complete" || loadedTasks[0].Done != true {
		t.Errorf("Task not correctly completed. Got: %+v", loadedTasks)
	}
}

func TestCompleteCommand_TaskNotFound(t *testing.T) {
	setupTest()
	defer func() { exit = os.Exit }() // Ensure os.Exit is restored
	exit = mockExit

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"todo", "done", "-id", "99"}

	exitCode := 0
	out, _ := captureOutput(func() {
		defer func() {
			if r := recover(); r != nil {
				exitCode, _ = strconv.Atoi(r.(string))
			}
		}()
		main()
	})

	expectedError := `{"error": "task with ID 99 not found"}`
	if !strings.Contains(out, expectedError) {
		t.Errorf("Expected error for task not found not present.\nExpected: %s\nGot: %s", expectedError, out)
	}

	if exitCode != 42 {
		t.Errorf("Expected exit code 42 for task not found, got %d", exitCode)
	}
}

func TestCompleteCommand_MissingID(t *testing.T) {
	setupTest()
	defer func() { exit = os.Exit }() // Ensure os.Exit is restored
	exit = mockExit

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"todo", "done"}

	exitCode := 0
	out, _ := captureOutput(func() {
		defer func() {
			if r := recover(); r != nil {
				exitCode, _ = strconv.Atoi(r.(string))
			}
		}()
		main()
	})

	expectedError := `{"error": "task ID cannot be empty"}`
	if !strings.Contains(out, expectedError) {
		t.Errorf("Expected error for missing ID not found.\nExpected: %s\nGot: %s", expectedError, out)
	}

	if exitCode != 42 {
		t.Errorf("Expected exit code 42 for missing ID, got %d", exitCode)
	}
}

func TestUnknownCommand(t *testing.T) {
	setupTest()
	defer func() { exit = os.Exit }() // Ensure os.Exit is restored
	exit = mockExit

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"todo", "unknown"}

	exitCode := 0
	out, _ := captureOutput(func() {
		defer func() {
			if r := recover(); r != nil {
				exitCode, _ = strconv.Atoi(r.(string))
			}
		}()
		main()
	})

	expectedError := `{"error": "unknown subcommand"}`
	if !strings.Contains(out, expectedError) {
		t.Errorf("Expected error for unknown subcommand not found.\nExpected: %s\nGot: %s", expectedError, out)
	}

	if exitCode != 42 {
		t.Errorf("Expected exit code 42 for unknown subcommand, got %d", exitCode)
	}
}

func TestNoCommand(t *testing.T) {
	setupTest()
	defer func() { exit = os.Exit }() // Ensure os.Exit is restored
	exit = mockExit

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"todo"}

	exitCode := 0
	out, _ := captureOutput(func() {
		defer func() {
			if r := recover(); r != nil {
				exitCode, _ = strconv.Atoi(r.(string))
			}
		}()
		main()
	})

	expectedError := `{"error": "expected 'add', 'list', or 'done' subcommand"}`
	if !strings.Contains(out, expectedError) {
		t.Errorf("Expected error for no subcommand not found.\nExpected: %s\nGot: %s", expectedError, out)
	}

	if exitCode != 42 {
		t.Errorf("Expected exit code 42 for no subcommand, got %d", exitCode)
	}
}
