package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

const testTasksFile = "test_tasks.json"

// runTodoCommand compiles and runs the todo CLI with given arguments.
// It returns stdout, stderr, and the exit code.
func runTodoCommand(t *testing.T, initialTestDir string, args ...string) (string, string, int, string) {
	var testDir string
	if initialTestDir != "" {
		testDir = initialTestDir
	} else {
		var err error
		testDir, err = ioutil.TempDir("", "todo_test_")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		t.Cleanup(func() {
			os.RemoveAll(testDir)
		})
	}

	execPath := filepath.Join(testDir, "todo")

	// Compile the main package into the temporary directory (if not already compiled)
	// This ensures the executable is present in the target directory for execution.
	// A more efficient approach would compile once outside this function.
	cmd := exec.Command("go", "build", "-o", execPath, ".")
	cmd.Dir = "." // Compile from the original package directory
	buildOut, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to compile todo tool: %v\nOutput: %s", err, buildOut)
	}

	// Run the compiled todo command
	cmd = exec.Command(execPath, args...)
	cmd.Dir = testDir // Run in the temporary directory

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	exitCode := 0
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode = exitErr.ExitCode()
	} else if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}

	return stdout.String(), stderr.String(), exitCode, testDir
}

// loadTasksFromDir reads tasks from a tasks.json file within a specified directory.
func loadTasksFromDir(t *testing.T, dir string) []Task {
	filePath := filepath.Join(dir, tasksFile) // Use the actual tasksFile name
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		// If file doesn't exist, treat as empty list
		if os.IsNotExist(err) {
			return []Task{}
		}
		t.Fatalf("Error reading tasks file %s: %v", filePath, err)
	}

	if len(data) == 0 {
		return []Task{}
	}

	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		t.Fatalf("Error unmarshalling tasks from %s: %v", filePath, err)
	}

	return tasks
}

func TestLoadAndSaveTasks(t *testing.T) {
	// This test still directly manipulates tasksFile, so we need a separate setup.
	tempFile := filepath.Join(t.TempDir(), tasksFile) // Use t.TempDir for this test
	oldTasksFile := tasksFile
	tasksFile = tempFile
	t.Cleanup(func() {
		tasksFile = oldTasksFile
	})

	// Test saving and loading empty list
	emptyTasks := []Task{}
	err := SaveTasks(emptyTasks)
	if err != nil {
		t.Fatalf("SaveTasks failed for empty list: %v", err)
	}

	loadedTasks, err := LoadTasks()
	if err != nil {
		t.Fatalf("LoadTasks failed for empty list: %v", err)
	}
	if len(loadedTasks) != 0 {
		t.Errorf("Expected 0 tasks, got %d", len(loadedTasks))
	}

	// Test saving and loading tasks
	tasks := []Task{
		{ID: 1, Description: "Task 1", Done: false},
		{ID: 2, Description: "Task 2", Done: true},
	}
	err = SaveTasks(tasks)
	if err != nil {
		t.Fatalf("SaveTasks failed: %v", err)
	}

	loadedTasks, err = LoadTasks()
	if err != nil {
		t.Fatalf("LoadTasks failed: %v", err)
	}
	if len(loadedTasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(loadedTasks))
	}
	if loadedTasks[0].ID != 1 || loadedTasks[1].ID != 2 {
		t.Errorf("Loaded tasks do not match expected IDs")
	}
}

func TestAddCommand(t *testing.T) {
	// Use a shared temporary directory for sequential operations
	testDir, err := ioutil.TempDir("", "todo_add_test_")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() {
		os.RemoveAll(testDir)
	})

	// Test adding a single task
	out, _, exitCode, _ := runTodoCommand(t, testDir, "add", "-description", "Buy milk")
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}
	expected := `{"id":1,"status":"added"}\n`
	if out != expected {
		t.Errorf("Expected %q, got %q", expected, out)
	}

	// Verify task is saved
	tasks := loadTasksFromDir(t, testDir)
	if len(tasks) != 1 || tasks[0].Description != "Buy milk" || tasks[0].ID != 1 {
		t.Errorf("Task not saved correctly: %+v", tasks)
	}

	// Test adding another task, checking ID generation in the same directory
	out, _, exitCode, _ = runTodoCommand(t, testDir, "add", "-description", "Do laundry")
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}
	expected = `{"id":2,"status":"added"}\n`
	if out != expected {
		t.Errorf("Expected %q, got %q", expected, out)
	}

	// Verify second task is saved
	tasks = loadTasksFromDir(t, testDir)
	if len(tasks) != 2 || tasks[1].Description != "Do laundry" || tasks[1].ID != 2 {
		t.Errorf("Second task not saved correctly: %+v", tasks)
	}

	// Test invalid arguments (missing description)
	out, _, exitCode, _ = runTodoCommand(t, testDir, "add", "-description", "")
	if exitCode != 42 {
		t.Errorf("Expected exit code 42, got %d", exitCode)
	}
	expected = `{"error":"invalid arguments"}\n`
	if out != expected {
		t.Errorf("Expected %q, got %q", expected, out)
	}
}

func TestListCommand(t *testing.T) {
	// Use a shared temporary directory for sequential operations
	testDir, err := ioutil.TempDir("", "todo_list_test_")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() {
		os.RemoveAll(testDir)
	})

	// Test listing empty tasks
	out, _, exitCode, _ := runTodoCommand(t, testDir, "list")
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}
	expected := `[]\n`
	if out != expected {
		t.Errorf("Expected %q, got %q", expected, out)
	}

	// Add some tasks
	out, _, exitCode, _ = runTodoCommand(t, testDir, "add", "-description", "Task 1")
	if exitCode != 0 {
		t.Fatalf("Failed to add Task 1: %v", exitCode)
	}
	out, _, exitCode, _ = runTodoCommand(t, testDir, "add", "-description", "Task 2")
	if exitCode != 0 {
		t.Fatalf("Failed to add Task 2: %v", exitCode)
	}
	// Mark one task as done to test pending list
	out, _, exitCode, _ = runTodoCommand(t, testDir, "done", "1")
	if exitCode != 0 {
		t.Fatalf("Failed to mark Task 1 as done: %v", exitCode)
	}

	// Test listing pending tasks
	out, _, exitCode, _ = runTodoCommand(t, testDir, "list")
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}
	expected = `[{"id":2,"description":"Task 2","done":false}]\n`
	if out != expected {
		t.Errorf("Expected %q, got %q", expected, out)
	}
}

func TestDoneCommand(t *testing.T) {
	// Use a shared temporary directory for sequential operations
	testDir, err := ioutil.TempDir("", "todo_done_test_")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() {
		os.RemoveAll(testDir)
	})

	// Add tasks for testing done command
	_, _, exitCode, _ := runTodoCommand(t, testDir, "add", "-description", "Task 1")
	if exitCode != 0 {
		t.Fatalf("Failed to add Task 1: %v", exitCode)
	}
	_, _, exitCode, _ = runTodoCommand(t, testDir, "add", "-description", "Task 2")
	if exitCode != 0 {
		t.Fatalf("Failed to add Task 2: %v", exitCode)
	}

	// Test completing an existing task
	out, _, exitCode, _ := runTodoCommand(t, testDir, "done", "1")
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}
	expected := `{"id":1,"status":"completed"}\n`
	if out != expected {
		t.Errorf("Expected %q, got %q", expected, out)
	}

	// Verify task is marked done
	tasks := loadTasksFromDir(t, testDir)
	if len(tasks) != 2 || !tasks[0].Done {
		t.Errorf("Task 1 not marked as done: %+v", tasks)
	}

	// Test completing a non-existent task
	out, _, exitCode, _ = runTodoCommand(t, testDir, "done", "99")
	if exitCode != 42 {
		t.Errorf("Expected exit code 42, got %d", exitCode)
	}
	expected = `{"error":"task with ID 99 not found"}\n`
	if out != expected {
		t.Errorf("Expected %q, got %q", expected, out)
	}

	// Test invalid arguments (non-integer ID)
	out, _, exitCode, _ = runTodoCommand(t, testDir, "done", "abc")
	if exitCode != 42 {
		t.Errorf("Expected exit code 42, got %d", exitCode)
	}
	expected = `{"error":"invalid arguments"}\n`
	if out != expected {
		t.Errorf("Expected %q, got %q", expected, out)
	}
}

// Test for no command provided
func TestNoCommand(t *testing.T) {
	out, _, exitCode, _ := runTodoCommand(t, "")
	if exitCode != 42 {
		t.Errorf("Expected exit code 42, got %d", exitCode)
	}
	expected := `{"error":"no command provided"}\n`
	if out != expected {
		t.Errorf("Expected %q, got %q", expected, out)
	}
}

// Test for unknown command
func TestUnknownCommand(t *testing.T) {
	out, _, exitCode, _ := runTodoCommand(t, "", "foobar")
	if exitCode != 42 {
		t.Errorf("Expected exit code 42, got %d", exitCode)
	}
	expected := `{"error":"unknown command"}\n`
	if out != expected {
		t.Errorf("Expected %q, got %q", expected, out)
	}
}