package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"reflect"
	"testing"
)

// Helper to set up a temporary tasks file for tests
func setupTestTasksFile(t *testing.T, content []Task) {
	data, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal test tasks: %v", err)
	}

	err = os.WriteFile(tasksFile, data, 0644)
	if err != nil {
		t.Fatalf("failed to write test tasks file: %v", err)
	}
}

// Helper to clean up the temporary tasks file
func cleanupTestTasksFile(t *testing.T) {
	err := os.Remove(tasksFile)
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("failed to remove test tasks file: %v", err)
	}
}

func TestLoadAndSaveTasks(t *testing.T) {
	cleanupTestTasksFile(t)
	defer cleanupTestTasksFile(t)

	t.Run("empty file", func(t *testing.T) {
		tasks, err := loadTasks()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(tasks) != 0 {
			t.Errorf("expected empty tasks, got %v", tasks)
		}
	})

	t.Run("non-empty file", func(t *testing.T) {
		expectedTasks := []Task{
			{ID: 1, Description: "Test task 1", Done: false},
			{ID: 2, Description: "Test task 2", Done: true},
		}
		setupTestTasksFile(t, expectedTasks)

		tasks, err := loadTasks()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(tasks, expectedTasks) {
			t.Errorf("expected %v, got %v", expectedTasks, tasks)
		}

		// Test saving an empty list
		saveTasks([]Task{})
		loadedTasks, err := loadTasks()
		if err != nil {
			t.Fatalf("unexpected error after saving empty list: %v", err)
		}
		if len(loadedTasks) != 0 {
			t.Errorf("expected empty tasks after saving empty list, got %v", loadedTasks)
		}
	})
}

// captureOutput is a helper to run the main function and capture its output and exit code.
func captureOutput(f func()) (string, int) {
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	oldExitFunc := exit // Save the original exit function

	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w

	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	writers := io.MultiWriter(outBuf, errBuf)

	exitCode := 0
	exit = func(code int) {
		exitCode = code
		panic("exit was called") // Panic to stop execution without actually exiting
	}

	defer func() {
		os.Stdout = oldStdout
		os.Stderr = oldStderr
		exit = oldExitFunc // Restore the original exit function
	}()

	func() {
		defer func() {
			if r := recover(); r != nil {
				if r != "exit was called" {
					panic(r)
				}
			}
		}()
		f()
	}()

	w.Close()
	io.Copy(writers, r) // Copy all output from the pipe to our buffers
	r.Close()

	return outBuf.String(), exitCode
}

func TestMainCommands(t *testing.T) {
	cleanupTestTasksFile(t)
	defer cleanupTestTasksFile(t)

	t.Run("add command", func(t *testing.T) {
		// Test adding a task
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		os.Args = []string{"todo", "add", "Buy groceries"}

		output, exitCode := captureOutput(main)
		if exitCode != ExitCodeOK {
			t.Errorf("expected exit code %d, got %d", ExitCodeOK, exitCode)
		}

		expectedOutput := "{\"id\":1,\"status\":\"added\"}\n"
		if output != expectedOutput {
			t.Errorf("expected output %q, got %q", expectedOutput, output)
		}

		tasks, err := loadTasks()
		if err != nil {
			t.Fatalf("failed to load tasks after add: %v", err)
		}
		if len(tasks) != 1 || tasks[0].Description != "Buy groceries" || tasks[0].ID != 1 || tasks[0].Done != false {
			t.Errorf("task not added correctly: %v", tasks)
		}

		// Test adding another task
		os.Args = []string{"todo", "add", "Clean room"}
		output, exitCode = captureOutput(main)
		if exitCode != ExitCodeOK {
			t.Errorf("expected exit code %d, got %d", ExitCodeOK, exitCode)
		}
		expectedOutput = "{\"id\":2,\"status\":\"added\"}\n"
		if output != expectedOutput {
			t.Errorf("expected output %q, got %q", expectedOutput, output)
		}

		tasks, err = loadTasks()
		if err != nil {
			t.Fatalf("failed to load tasks after second add: %v", err)
		}
		if len(tasks) != 2 || tasks[1].Description != "Clean room" || tasks[1].ID != 2 || tasks[1].Done != false {
			t.Errorf("second task not added correctly: %v", tasks)
		}
	})

	t.Run("list command", func(t *testing.T) {
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		os.Args = []string{"todo", "list"}

		output, exitCode := captureOutput(main)
		if exitCode != ExitCodeOK {
			t.Errorf("expected exit code %d, got %d", ExitCodeOK, exitCode)
		}

		expectedTasks := []Task{
			{ID: 1, Description: "Buy groceries", Done: false},
			{ID: 2, Description: "Clean room", Done: false},
		}
		expectedOutputBytes, _ := json.Marshal(expectedTasks)
		expectedOutput := string(expectedOutputBytes) + "\n"

		if output != expectedOutput {
			t.Errorf("expected output %q, got %q", expectedOutput, output)
		}
	})

	t.Run("done command", func(t *testing.T) {
		// Test completing a task
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		os.Args = []string{"todo", "done", "1"}

		output, exitCode := captureOutput(main)
		if exitCode != ExitCodeOK {
			t.Errorf("expected exit code %d, got %d", ExitCodeOK, exitCode)
		}

		expectedOutput := "{\"id\":1,\"status\":\"completed\"}\n"
		if output != expectedOutput {
			t.Errorf("expected output %q, got %q", expectedOutput, output)
		}

		tasks, err := loadTasks()
		if err != nil {
			t.Fatalf("failed to load tasks after done: %v", err)
		}
		if len(tasks) != 2 || tasks[0].Done != true || tasks[1].Done != false {
			t.Errorf("task not marked as done correctly: %v", tasks)
		}

		// Test listing after completing a task
		os.Args = []string{"todo", "list"}
		output, exitCode = captureOutput(main)
		if exitCode != ExitCodeOK {
			t.Errorf("expected exit code %d, got %d", ExitCodeOK, exitCode)
		}

		expectedTasks := []Task{
			{ID: 2, Description: "Clean room", Done: false},
		}
		expectedOutputBytes, _ := json.Marshal(expectedTasks)
		expectedOutput = string(expectedOutputBytes) + "\n"

		if output != expectedOutput {
			t.Errorf("expected output %q, got %q", expectedOutput, output)
		}

		// Test done on non-existent ID
		os.Args = []string{"todo", "done", "99"}
		output, exitCode = captureOutput(main)
		if exitCode != ExitCodeError {
			t.Errorf("expected exit code %d, got %d", ExitCodeError, exitCode)
		}
		expectedErrorOutput := "{\"error\":\"task with ID 99 not found\"}\n"
		if output != expectedErrorOutput {
			t.Errorf("expected error output %q, got %q", expectedErrorOutput, output)
		}

		// Test done with invalid ID format
		os.Args = []string{"todo", "done", "abc"}
		output, exitCode = captureOutput(main)
		if exitCode != ExitCodeInvalidArguments {
			t.Errorf("expected exit code %d, got %d", ExitCodeInvalidArguments, exitCode)
		}
		expectedErrorOutput = "{\"error\":\"invalid task ID: strconv.Atoi: parsing \\\"abc\\\": invalid syntax\"}\n"
		if output != expectedErrorOutput {
			t.Errorf("expected error output %q, got %q", expectedErrorOutput, output)
		}
	})

	t.Run("invalid arguments", func(t *testing.T) {
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		// No command
		os.Args = []string{"todo"}
		output, exitCode := captureOutput(main)
		if exitCode != ExitCodeInvalidArguments {
			t.Errorf("expected exit code %d, got %d", ExitCodeInvalidArguments, exitCode)
		}
		expectedOutput := "{\"error\":\"invalid arguments\"}\n"
		if output != expectedOutput {
			t.Errorf("expected output %q, got %q", expectedOutput, output)
		}

		// Unknown command
		os.Args = []string{"todo", "unknown"}
		output, exitCode = captureOutput(main)
		if exitCode != ExitCodeInvalidArguments {
			t.Errorf("expected exit code %d, got %d", ExitCodeInvalidArguments, exitCode)
		}
		expectedOutput = "{\"error\":\"unknown command: unknown\"}\n"
		if output != expectedOutput {
			t.Errorf("expected output %q, got %q", expectedOutput, output)
		}

		// Add with no description
		os.Args = []string{"todo", "add"}
		output, exitCode = captureOutput(main)
		if exitCode != ExitCodeInvalidArguments {
			t.Errorf("expected exit code %d, got %d", ExitCodeInvalidArguments, exitCode)
		}
		expectedOutput = "{\"error\":\"missing task description\"}\n"
		if output != expectedOutput {
			t.Errorf("expected output %q, got %q", expectedOutput, output)
		}
	})
}
