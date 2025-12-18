package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

// captureOutput redirects stdout to a buffer and returns its content.
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// cleanUp removes the tasks.json file after each test.
func cleanUp() {
	os.Remove(tasksFile)
}

func TestAddCommand(t *testing.T) {
	defer cleanUp()

	output, code := runApp([]string{"cmd", "add", "-description", "Buy groceries"})

	if code != 0 {
		t.Fatalf("Expected exit code 0, got %d. Output: %s", code, output)
	}

	expected := map[string]interface{}{"id": float64(1), "status": "added"}
	var actual map[string]interface{}
	err := json.Unmarshal([]byte(output), &actual)
	if err != nil {
		t.Fatalf("Failed to unmarshal output: %v, output: %s", err, output)
	}

	if actual["id"] != expected["id"] || actual["status"] != expected["status"] {
		t.Errorf("Expected %v, got %v", expected, actual)
	}

	// Verify task is saved
	data, err := ioutil.ReadFile(tasksFile)
	if err != nil {
		t.Fatalf("Failed to read tasks file: %v", err)
	}
	var savedTasks []Task
	err = json.Unmarshal(data, &savedTasks)
	if err != nil {
		t.Fatalf("Failed to unmarshal saved tasks: %v", err)
	}

	if len(savedTasks) != 1 || savedTasks[0].Description != "Buy groceries" || savedTasks[0].ID != 1 || savedTasks[0].Done != false {
		t.Errorf("Saved tasks mismatch: %v", savedTasks)
	}
}

func TestListCommand(t *testing.T) {
	defer cleanUp()

	// Add a task first
	_, _ = runApp([]string{"cmd", "add", "-description", "Walk the dog"})

	// Add another task
	_, _ = runApp([]string{"cmd", "add", "-description", "Pay bills"})

	// Mark one task as done
	_, _ = runApp([]string{"cmd", "done", "-id", "1"})

	// List pending tasks
	output, code := runApp([]string{"cmd", "list"})

	if code != 0 {
		t.Fatalf("Expected exit code 0, got %d. Output: %s", code, output)
	}

	expectedTasks := []Task{
		{ID: 2, Description: "Pay bills", Done: false},
	}
	var actualTasks []Task
	err := json.Unmarshal([]byte(output), &actualTasks)
	if err != nil {
		t.Fatalf("Failed to unmarshal output: %v, output: %s", err, output)
	}

	// Compare by marshaling and unmarshaling to handle potential order differences if any
	expectedJSON, _ := json.Marshal(expectedTasks)
	actualJSON, _ := json.Marshal(actualTasks)

	if !bytes.Equal(expectedJSON, actualJSON) {
		t.Errorf("Expected %s, got %s", string(expectedJSON), string(actualJSON))
	}


	// Test listing when no tasks are present
	cleanUp() // Clear tasks.json and in-memory tasks
	output, code = runApp([]string{"cmd", "list"})
	if code != 0 {
		t.Fatalf("Expected exit code 0, got %d. Output: %s", code, output)
	}
	if output != "[]\n" && output != "[]" {
		t.Errorf("Expected empty array for list, got: %s", output)
	}
}

func TestDoneCommand(t *testing.T) {
	defer cleanUp()

	// Add a task
	_, _ = runApp([]string{"cmd", "add", "-description", "Read a book"})

	// Mark task as done
	output, code := runApp([]string{"cmd", "done", "-id", "1"})

	if code != 0 {
		t.Fatalf("Expected exit code 0, got %d. Output: %s", code, output)
	}

	expected := map[string]interface{}{"id": float64(1), "status": "completed"}
	var actual map[string]interface{}
	err := json.Unmarshal([]byte(output), &actual)
	if err != nil {
		t.Fatalf("Failed to unmarshal output: %v, output: %s", err, output)
	}

	if actual["id"] != expected["id"] || actual["status"] != expected["status"] {
		t.Errorf("Expected %v, got %v", expected, actual)
	}

	// Verify task status in file
	data, err := ioutil.ReadFile(tasksFile)
	if err != nil {
		t.Fatalf("Failed to read tasks file: %v", err)
	}
	var savedTasks []Task
	err = json.Unmarshal(data, &savedTasks)
	if err != nil {
		t.Fatalf("Failed to unmarshal saved tasks: %v", err)
	}

	if len(savedTasks) != 1 || savedTasks[0].ID != 1 || savedTasks[0].Done != true {
		t.Errorf("Saved task not marked as done: %v", savedTasks)
	}
}

func TestInvalidAddArguments(t *testing.T) {
	defer cleanUp()

	output, code := runApp([]string{"cmd", "add"})

	expectedError := formatError("invalid arguments: description is required")

	if output != expectedError + "\n" && output != expectedError {
		t.Errorf("Expected error %s, got %s", expectedError, output)
	}

	if code != 42 {
		t.Errorf("Expected exit code 42, got %d", code)
	}
}

func TestUnknownCommand(t *testing.T) {
	defer cleanUp()

	output, code := runApp([]string{"cmd", "foobar"})

	expectedError := formatError("unknown command: foobar")

	if output != expectedError + "\n" && output != expectedError {
		t.Errorf("Expected error %s, got %s", expectedError, output)
	}

	if code != 42 {
		t.Errorf("Expected exit code 42, got %d", code)
	}
}

func TestInvalidDoneArguments(t *testing.T) {
	defer cleanUp()

	output, code := runApp([]string{"cmd", "done"})

	expectedError := formatError("invalid arguments: task ID is required")

	if output != expectedError + "\n" && output != expectedError {
		t.Errorf("Expected error %s, got %s", expectedError, output)
	}

	if code != 42 {
		t.Errorf("Expected exit code 42, got %d", code)
	}
}

func TestTaskNotFound(t *testing.T) {
	defer cleanUp()

	output, code := runApp([]string{"cmd", "done", "-id", "999"})

	expectedError := formatError("task with ID 999 not found")

	if output != expectedError + "\n" && output != expectedError {
		t.Errorf("Expected error %s, got %s", expectedError, output)
	}

	if code != 1 {
		t.Errorf("Expected exit code 1, got %d", code)
	}
}

func TestNoCommandProvided(t *testing.T) {
	defer cleanUp()

	output, code := runApp([]string{"cmd"})

	expectedError := formatError("no command provided")

	if output != expectedError + "\n" && output != expectedError {
		t.Errorf("Expected error %s, got %s", expectedError, output)
	}

	if code != 42 {
		t.Errorf("Expected exit code 42, got %d", code)
	}
}