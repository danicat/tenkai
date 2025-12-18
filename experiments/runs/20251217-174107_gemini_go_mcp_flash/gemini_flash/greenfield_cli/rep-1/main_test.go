package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// mockExit is a mock function for os.Exit to prevent actual program termination during tests.
var mockExitCode int
var mockExitCalled bool

func mockExit(code int) {
	mockExitCode = code
	mockExitCalled = true
	panic("exit called") // Panic to stop execution in test
}

// testTasksFile is used to manage the tasks.json path during tests.
var testTasksFile string

func TestMain(m *testing.M) {
	// Create a temporary directory for test files
	testDir, err := ioutil.TempDir("", "todo_test_")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(testDir) // Clean up temporary directory after all tests

	// Set tasksFile to a path within the temporary directory for tests
	testTasksFile = filepath.Join(testDir, "tasks.json")

	// Run tests
	code := m.Run()

	os.Exit(code)
}

func TestAddTask(t *testing.T) {
	// Clear tasks.json before test
	if err := os.Remove(testTasksFile);
		err != nil && !os.IsNotExist(err) {
		t.Fatalf("failed to remove %s: %v", testTasksFile, err)
	}

	description := "Test Task 1"
	addTask(description, testTasksFile, mockExit)

	tasks, err := readTasks(testTasksFile)
	if err != nil {
		t.Fatalf("failed to read tasks: %v", err)
	}

	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}

	if tasks[0].Description != description || tasks[0].ID != 1 || tasks[0].Done != false {
		t.Errorf("unexpected task content: %+v", tasks[0])
	}

	description2 := "Test Task 2"
	addTask(description2, testTasksFile, mockExit)

	tasks, err = readTasks(testTasksFile)
	if err != nil {
		t.Fatalf("failed to read tasks: %v", err)
	}

	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}

	if tasks[1].Description != description2 || tasks[1].ID != 2 || tasks[1].Done != false {
		t.Errorf("unexpected task content: %+v", tasks[1])
	}
}

func TestListTasks(t *testing.T) {
	// Clear tasks.json before test
	if err := os.Remove(testTasksFile);
		err != nil && !os.IsNotExist(err) {
		t.Fatalf("failed to remove %s: %v", testTasksFile, err)
	}

	addTask("Task for listing 1", testTasksFile, mockExit)
	addTask("Task for listing 2", testTasksFile, mockExit)

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	listTasks(testTasksFile, mockExit)

	w.Close()
	os.Stdout = oldStdout
	out, _ := ioutil.ReadAll(r)

	var listedTasks []Task
	if err := json.Unmarshal(out, &listedTasks); err != nil {
		t.Fatalf("failed to unmarshal list output: %v", err)
	}

	if len(listedTasks) != 2 {
		t.Fatalf("expected 2 listed tasks, got %d", len(listedTasks))
	}

	if listedTasks[0].Description != "Task for listing 1" || listedTasks[1].Description != "Task for listing 2" {
		t.Errorf("unexpected listed tasks: %+v", listedTasks)
	}
}

func TestCompleteTask(t *testing.T) {
	// Clear tasks.json before test
	if err := os.Remove(testTasksFile);
		err != nil && !os.IsNotExist(err) {
		t.Fatalf("failed to remove %s: %v", testTasksFile, err)
	}

	addTask("Task to complete", testTasksFile, mockExit)
	addTask("Another task", testTasksFile, mockExit)

	completeTask(1, testTasksFile, mockExit)

	tasks, err := readTasks(testTasksFile)
	if err != nil {
		t.Fatalf("failed to read tasks: %v", err)
	}

	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks after completion, got %d", len(tasks))
	}

	if !tasks[0].Done {
		t.Errorf("task 1 should be done, but is not: %+v", tasks[0])
	}

	if tasks[1].Done {
		t.Errorf("task 2 should not be done, but is: %+v", tasks[1])
	}

	// Test completing a non-existent task
	mockExitCalled = false
	mockExitCode = 0

	defer func() {
		r := recover()
		if r == nil || r.(string) != "exit called" {
			t.Errorf("expected mockExit to be called")
		}
		if !mockExitCalled || mockExitCode != 42 {
			t.Errorf("expected mockExit to be called with code 42, got %d", mockExitCode)
		}
		// Restore mockExitCalled and mockExitCode for subsequent tests
		mockExitCalled = false
		mockExitCode = 0
	}()

	completeTask(99, testTasksFile, mockExit)
}

func TestRunCLI(t *testing.T) {
	// Save and restore os.Args and os.Stdout/Stderr
	oldArgs := os.Args
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	defer func() {
		os.Args = oldArgs
		os.Stdout = oldStdout
		os.Stderr = oldStderr
	}()

	// Temporarily set tasksFile to the testTasksFile for runCLI tests
	originalTasksFile := tasksFile
	tasksFile = testTasksFile
	defer func() { tasksFile = originalTasksFile }()

	t.Run("add command", func(t *testing.T) {
		// Clear tasks.json before test
		if err := os.Remove(testTasksFile);
			err != nil && !os.IsNotExist(err) {
			t.Fatalf("failed to remove %s: %v", testTasksFile, err)
		}

		oldStdoutSubtest := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		os.Args = []string{"todo", "add", "Test Task"}

		mockExitCalled = false
		defer func() {
			r := recover()
			if r != nil && r.(string) == "exit called" {
				// Expected if runCLI calls mockExit
			} else if r != nil {
				t.Fatalf("unexpected panic: %v", r)
			}
			w.Close()
			os.Stdout = oldStdoutSubtest
		}()
		runCLI(os.Args, mockExit)

		out, _ := ioutil.ReadAll(r)
		expected := `{"id":1,"status":"added"}`
		if strings.TrimSpace(string(out)) != expected {
			t.Errorf("expected %s, got %s", expected, strings.TrimSpace(string(out)))
		}
	})

	t.Run("list command", func(t *testing.T) {
		// Add a task for listing
		addTask("Task for listing", testTasksFile, mockExit)

		oldStdoutSubtest := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		os.Args = []string{"todo", "list"}

		mockExitCalled = false
		defer func() {
			r := recover()
			if r != nil && r.(string) == "exit called" {
				// Expected if runCLI calls mockExit
			} else if r != nil {
				t.Fatalf("unexpected panic: %v", r)
			}
			w.Close()
			os.Stdout = oldStdoutSubtest
		}()
		runCLI(os.Args, mockExit)

		out, _ := ioutil.ReadAll(r)
		expected := `[{"id":1,"description":"Task for listing","done":false}]`
		if strings.TrimSpace(string(out)) != expected {
			t.Errorf("expected %s, got %s", expected, strings.TrimSpace(string(out)))
		}
	})

	t.Run("done command", func(t *testing.T) {
		// Add a task to complete
		addTask("Task to complete", testTasksFile, mockExit)

		oldStdoutSubtest := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		os.Args = []string{"todo", "done", "1"}

		mockExitCalled = false
		defer func() {
			r := recover()
			if r != nil && r.(string) == "exit called" {
				// Expected if runCLI calls mockExit
			} else if r != nil {
				t.Fatalf("unexpected panic: %v", r)
			}
			w.Close()
			os.Stdout = oldStdoutSubtest
		}()
		runCLI(os.Args, mockExit)

		out, _ := ioutil.ReadAll(r)
		expected := `{"id":1,"status":"completed"}`
		if strings.TrimSpace(string(out)) != expected {
			t.Errorf("expected %s, got %s", expected, strings.TrimSpace(string(out)))
		}
	})

	t.Run("invalid command", func(t *testing.T) {
		oldStdoutSubtest := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		os.Args = []string{"todo", "unknown"}

		mockExitCalled = false
		mockExitCode = 0

		defer func() {
			r := recover()
			if r == nil || r.(string) != "exit called" {
				t.Errorf("expected mockExit to be called")
			}
			if !mockExitCalled || mockExitCode != 42 {
				t.Errorf("expected mockExit to be called with code 42, got %d", mockExitCode)
			}
			w.Close()
			os.Stdout = oldStdoutSubtest
		}()
		runCLI(os.Args, mockExit)

		out, _ := ioutil.ReadAll(r)
		expectedError := `{"error":"unknown command: unknown"}`
		if strings.TrimSpace(string(out)) != expectedError {
			t.Errorf("expected %s, got %s", expectedError, strings.TrimSpace(string(out)))
		}
	})

	t.Run("too few arguments", func(t *testing.T) {
		oldStdoutSubtest := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		os.Args = []string{"todo"}

		mockExitCalled = false
		mockExitCode = 0

		defer func() {
			r := recover()
			if r == nil || r.(string) != "exit called" {
				t.Errorf("expected mockExit to be called")
			}
			if !mockExitCalled || mockExitCode != 42 {
				t.Errorf("expected mockExit to be called with code 42, got %d", mockExitCode)
			}
			w.Close()
			os.Stdout = oldStdoutSubtest
		}()
		runCLI(os.Args, mockExit)

		out, _ := ioutil.ReadAll(r)
		expectedError := `{"error":"invalid arguments"}`
		if strings.TrimSpace(string(out)) != expectedError {
			t.Errorf("expected %s, got %s", expectedError, strings.TrimSpace(string(out)))
		}
	})

	t.Run("add with missing description", func(t *testing.T) {
		oldStdoutSubtest := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		os.Args = []string{"todo", "add"}

		mockExitCalled = false
		mockExitCode = 0

		defer func() {
			r := recover()
			if r == nil || r.(string) != "exit called" {
				t.Errorf("expected mockExit to be called")
			}
			if !mockExitCalled || mockExitCode != 42 {
				t.Errorf("expected mockExit to be called with code 42, got %d", mockExitCode)
			}
			w.Close()
			os.Stdout = oldStdoutSubtest
		}()
		runCLI(os.Args, mockExit)

		out, _ := ioutil.ReadAll(r)
		expectedError := `{"error":"invalid arguments"}`
		if strings.TrimSpace(string(out)) != expectedError {
			t.Errorf("expected %s, got %s", expectedError, strings.TrimSpace(string(out)))
		}
	})

	t.Run("done with missing ID", func(t *testing.T) {
		oldStdoutSubtest := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		os.Args = []string{"todo", "done"}

		mockExitCalled = false
		mockExitCode = 0

		defer func() {
			r := recover()
			if r == nil || r.(string) != "exit called" {
				t.Errorf("expected mockExit to be called")
			}
			if !mockExitCalled || mockExitCode != 42 {
				t.Errorf("expected mockExit to be called with code 42, got %d", mockExitCode)
			}
			w.Close()
			os.Stdout = oldStdoutSubtest
		}()
		runCLI(os.Args, mockExit)

		out, _ := ioutil.ReadAll(r)
		expectedError := `{"error":"invalid arguments"}`
		if strings.TrimSpace(string(out)) != expectedError {
			t.Errorf("expected %s, got %s", expectedError, strings.TrimSpace(string(out)))
		}
	})

	t.Run("done with invalid ID", func(t *testing.T) {
		oldStdoutSubtest := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		os.Args = []string{"todo", "done", "abc"}

		mockExitCalled = false
		mockExitCode = 0

		defer func() {
			r := recover()
			if r == nil || r.(string) != "exit called" {
				t.Errorf("expected mockExit to be called")
			}
			if !mockExitCalled || mockExitCode != 42 {
				t.Errorf("expected mockExit to be called with code 42, got %d", mockExitCode)
			}
			w.Close()
			os.Stdout = oldStdoutSubtest
		}()
		runCLI(os.Args, mockExit)

		out, _ := ioutil.ReadAll(r)
		expectedError := `{"error":"invalid arguments: task ID must be an integer"}`
		if strings.TrimSpace(string(out)) != expectedError {
			t.Errorf("expected %s, got %s", expectedError, strings.TrimSpace(string(out)))
		}
	})
}
