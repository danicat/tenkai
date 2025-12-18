package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/user/todo-cli/internal/task"
)

var testTasksFile = "test_main_tasks.json"

func TestMain(m *testing.M) {
	task.SetTasksFileForTest(testTasksFile)
	code := m.Run()
	os.Remove(testTasksFile)
	os.Exit(code)
}

func executeCommand(args []string) (string, int) {
	var buf bytes.Buffer
	exitCode := run(args, &buf)
	return strings.TrimSpace(buf.String()), exitCode
}

func TestAddCommand(t *testing.T) {
	os.Remove(testTasksFile) // Clean up before running tests

	// Test adding a task
	out, code := executeCommand([]string{"add", "Buy milk"})
	if code != 0 {
		t.Fatalf("Expected exit code 0, got %d. Output: %s", code, out)
	}

	expected := `{"id":1,"status":"added"}`
	if out != expected {
		t.Errorf("Expected %s, got %s", expected, out)
	}

	// Add another task
	out, code = executeCommand([]string{"add", "Go for a run"})
	if code != 0 {
		t.Fatalf("Expected exit code 0, got %d. Output: %s", code, out)
	}
	expected = `{"id":2,"status":"added"}`
	if out != expected {
		t.Errorf("Expected %s, got %s", expected, out)
	}

	// Test missing description
	out, code = executeCommand([]string{"add"})
	if code != exitCodeInvalidArgs {
		t.Errorf("Expected exit code %d, got %d", exitCodeInvalidArgs, code)
	}
	expected = fmt.Sprintf(`{"error":"%s"}`, errorInvalidArgs)
	if out != expected {
		t.Errorf("Expected %s, got %s", expected, out)
	}
}

func TestListCommand(t *testing.T) {
	os.Remove(testTasksFile)
	executeCommand([]string{"add", "Task A"})
	executeCommand([]string{"add", "Task B"})
	executeCommand([]string{"done", "1"})

	out, code := executeCommand([]string{"list"})
	if code != 0 {
		t.Fatalf("Expected exit code 0, got %d. Output: %s", code, out)
	}

	var outputData output
	if err := json.Unmarshal([]byte(out), &outputData); err != nil {
		t.Fatalf("Failed to unmarshal JSON output: %v, output: %s", err, out)
	}

	if len(outputData.Tasks) != 1 || outputData.Tasks[0].ID != 2 || outputData.Tasks[0].Description != "Task B" || outputData.Tasks[0].Done != false {
		t.Errorf("Expected 1 pending task (ID 2, 'Task B'), got: %+v", outputData.Tasks)
	}
}

func TestDoneCommand(t *testing.T) {
	os.Remove(testTasksFile)
	executeCommand([]string{"add", "Task to complete"})

	// Complete an existing task
	out, code := executeCommand([]string{"done", "1"})
	if code != 0 {
		t.Fatalf("Expected exit code 0, got %d. Output: %s", code, out)
	}
	expected := `{"id":1,"status":"completed"}`
	if out != expected {
		t.Errorf("Expected %s, got %s", expected, out)
	}

	// Try to complete already completed task
	out, code = executeCommand([]string{"done", "1"})
	if code != 1 {
		t.Errorf("Expected exit code 1 for already completed task, got %d", code)
	}
	expected = `{"error":"task 1 is already completed"}`
	if out != expected {
		t.Errorf("Expected %s, got %s", expected, out)
	}

	// Try to complete non-existent task
	out, code = executeCommand([]string{"done", "99"})
	if code != 1 {
		t.Errorf("Expected exit code 1 for non-existent task, got %d", code)
	}
	expected = `{"error":"task with ID 99 not found"}`
	if out != expected {
		t.Errorf("Expected %s, got %s", expected, out)
	}

	// Test invalid ID argument
	out, code = executeCommand([]string{"done", "abc"})
	if code != exitCodeInvalidArgs {
		t.Errorf("Expected exit code %d, got %d", exitCodeInvalidArgs, code)
	}
	expected = fmt.Sprintf(`{"error":"%s"}`, errorInvalidArgs)
	if out != expected {
		t.Errorf("Expected %s, got %s", expected, out)
	}
}

func TestUnknownCommand(t *testing.T) {
	out, code := executeCommand([]string{"unknown"})
	if code != exitCodeInvalidArgs {
		t.Errorf("Expected exit code %d, got %d", exitCodeInvalidArgs, code)
	}
	expected := fmt.Sprintf(`{"error":"unknown command: unknown"}`)
	if out != expected {
		t.Errorf("Expected %s, got %s", expected, out)
	}
}

func TestNoArguments(t *testing.T) {
	out, code := executeCommand([]string{})
	if code != exitCodeInvalidArgs {
		t.Errorf("Expected exit code %d, got %d", exitCodeInvalidArgs, code)
	}
	expected := fmt.Sprintf(`{"error":"%s"}`, errorInvalidArgs)
	if out != expected {
		t.Errorf("Expected %s, got %s", expected, out)
	}
}
