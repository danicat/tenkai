package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

// mockExit is a helper to prevent os.Exit from terminating the test
var exitCode int
var originalExit func(code int) = exit

func mockExit(code int) {
	exitCode = code
	panic("exit was called") // Panic to stop execution in the test function
}

// setUpTest prepares the test environment
func setUpTest() {
	// Restore exit after each test
	exit = mockExit
	// Ensure tasks.json does not exist before tests
	os.Remove(tasksFile)
	exitCode = 0 // Reset exit code
}

// tearDownTest cleans up the test environment
func tearDownTest() {
	exit = originalExit  // Restore original exit
	os.Remove(tasksFile) // Clean up tasks.json
}

func TestAddCommand(t *testing.T) {
	setUpTest()
	defer tearDownTest()

	// Redirect stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	os.Args = []string{"todo", "add", "Buy milk"}
	defer func() {
		if r := recover(); r != nil {
			// exit was called, check the exit code if needed
		}
		w.Close()
		os.Stdout = oldStdout
	}()
	main()

	w.Close()
	out, _ := ioutil.ReadAll(r)

	expected := SuccessResponse{ID: 1, Status: "added"}
	var actual SuccessResponse
	err := json.Unmarshal(out, &actual)
	if err != nil {
		t.Fatalf("Failed to unmarshal output: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}

	tasks, err := loadTasks()
	if err != nil {
		t.Fatalf("Failed to load tasks after add: %v", err)
	}
	if len(tasks) != 1 || tasks[0].Description != "Buy milk" || tasks[0].ID != 1 || tasks[0].Done != false {
		t.Errorf("Task not added correctly: %v", tasks)
	}
}

func TestAddCommand_InvalidArgs(t *testing.T) {
	setUpTest()
	defer tearDownTest()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	os.Args = []string{"todo", "add"}

	defer func() {
		if r := recover(); r != nil {
			if exitCode != 42 {
				t.Errorf("Expected exit code 42, got %d", exitCode)
			}
		} else {
			t.Errorf("Expected exit to be called")
		}
		w.Close()
		os.Stdout = oldStdout
	}()
	main()

	w.Close()
	out, _ := ioutil.ReadAll(r)

	expected := ErrorResponse{Error: "invalid arguments"}
	var actual ErrorResponse
	err := json.Unmarshal(out, &actual)
	if err != nil {
		t.Fatalf("Failed to unmarshal output: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestListCommand(t *testing.T) {
	setUpTest()
	defer tearDownTest()

	// Add a task directly to the file for testing list
	tasks := []Task{{
		ID:          1,
		Description: "Task 1",
		Done:        false,
	}, {
		ID:          2,
		Description: "Task 2",
		Done:        true,
	}}
	saveTasks(tasks)

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	os.Args = []string{"todo", "list"}
	defer func() {
		if r := recover(); r != nil {
			// exit was called, check the exit code if needed
		}
		w.Close()
		os.Stdout = oldStdout
	}()
	main()

	w.Close()
	out, _ := ioutil.ReadAll(r)

	expected := []Task{{
		ID:          1,
		Description: "Task 1",
		Done:        false,
	}}
	var actual []Task
	err := json.Unmarshal(out, &actual)
	if err != nil {
		t.Fatalf("Failed to unmarshal output: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestDoneCommand(t *testing.T) {
	setUpTest()
	defer tearDownTest()

	// Add a task directly to the file for testing done
	tasks := []Task{{
		ID:          1,
		Description: "Task 1",
		Done:        false,
	}}
	saveTasks(tasks)

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	os.Args = []string{"todo", "done", "1"}
	defer func() {
		if r := recover(); r != nil {
			// exit was called, check the exit code if needed
		}
		w.Close()
		os.Stdout = oldStdout
	}()
	main()

	w.Close()
	out, _ := ioutil.ReadAll(r)

	expected := SuccessResponse{ID: 1, Status: "completed"}
	var actual SuccessResponse
	err := json.Unmarshal(out, &actual)
	if err != nil {
		t.Fatalf("Failed to unmarshal output: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}

	// Verify task is marked as done in the file
	updatedTasks, err := loadTasks()
	if err != nil {
		t.Fatalf("Failed to load tasks after done: %v", err)
	}
	if len(updatedTasks) != 1 || updatedTasks[0].ID != 1 || updatedTasks[0].Done != true {
		t.Errorf("Task not marked as done correctly: %v", updatedTasks)
	}
}

func TestDoneCommand_InvalidID(t *testing.T) {
	setUpTest()
	defer tearDownTest()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	os.Args = []string{"todo", "done", "abc"}

	defer func() {
		if r := recover(); r != nil {
			if exitCode != 42 {
				t.Errorf("Expected exit code 42, got %d", exitCode)
			}
		} else {
			t.Errorf("Expected exit to be called")
		}
		w.Close()
		os.Stdout = oldStdout
	}()
	main()

	w.Close()
	out, _ := ioutil.ReadAll(r)

	expected := ErrorResponse{Error: "invalid task ID"}
	var actual ErrorResponse
	err := json.Unmarshal(out, &actual)
	if err != nil {
		t.Fatalf("Failed to unmarshal output: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestDoneCommand_NotFound(t *testing.T) {
	setUpTest()
	defer tearDownTest()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	os.Args = []string{"todo", "done", "999"}

	defer func() {
		if r := recover(); r != nil {
			if exitCode != 1 {
				t.Errorf("Expected exit code 1, got %d", exitCode)
			}
		} else {
			t.Errorf("Expected exit to be called")
		}
		w.Close()
		os.Stdout = oldStdout
	}()
	main()

	w.Close()
	out, _ := ioutil.ReadAll(r)

	expected := ErrorResponse{Error: "task not found"}
	var actual ErrorResponse
	err := json.Unmarshal(out, &actual)
	if err != nil {
		t.Fatalf("Failed to unmarshal output: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestUnknownCommand(t *testing.T) {
	setUpTest()
	defer tearDownTest()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	os.Args = []string{"todo", "unknown"}

	defer func() {
		if r := recover(); r != nil {
			if exitCode != 1 {
				t.Errorf("Expected exit code 1, got %d", exitCode)
			}
		} else {
			t.Errorf("Expected exit to be called")
		}
		w.Close()
		os.Stdout = oldStdout
	}()
	main()

	w.Close()
	out, _ := ioutil.ReadAll(r)

	expected := ErrorResponse{Error: "unknown command"}
	var actual ErrorResponse
	err := json.Unmarshal(out, &actual)
	if err != nil {
		t.Fatalf("Failed to unmarshal output: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestNoCommand(t *testing.T) {
	setUpTest()
	defer tearDownTest()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	os.Args = []string{"todo"}

	defer func() {
		if r := recover(); r != nil {
			if exitCode != 42 {
				t.Errorf("Expected exit code 42, got %d", exitCode)
			}
		} else {
			t.Errorf("Expected exit to be called")
		}
		w.Close()
		os.Stdout = oldStdout
	}()
	main()

	w.Close()
	out, _ := ioutil.ReadAll(r)

	expected := ErrorResponse{Error: "invalid arguments"}
	var actual ErrorResponse
	err := json.Unmarshal(out, &actual)
	if err != nil {
		t.Fatalf("Failed to unmarshal output: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}
