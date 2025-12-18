package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestLoadAndSaveTasks(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := ioutil.TempDir("", "todo_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up after test

	oldTasksFile := tasksFile
	tasksFile = filepath.Join(tempDir, "tasks.json")
	defer func() { tasksFile = oldTasksFile }() // Restore original tasksFile path

	// Test saving an empty list
	emptyTasks := []Task{}
	err = saveTasks(emptyTasks)
	if err != nil {
		t.Fatalf("Failed to save empty tasks: %v", err)
	}

	loadedTasks, err := loadTasks()
	if err != nil {
		t.Fatalf("Failed to load empty tasks: %v", err)
	}
	if len(loadedTasks) != 0 {
		t.Errorf("Expected 0 tasks, got %d", len(loadedTasks))
	}

	// Test saving some tasks
	someTasks := []Task{
		{ID: 1, Description: "Buy milk", Done: false},
		{ID: 2, Description: "Walk dog", Done: true},
	}
	err = saveTasks(someTasks)
	if err != nil {
		t.Fatalf("Failed to save some tasks: %v", err)
	}

	loadedTasks, err = loadTasks()
	if err != nil {
		t.Fatalf("Failed to load some tasks: %v", err)
	}
	if len(loadedTasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(loadedTasks))
	}
	if loadedTasks[0].ID != 1 || loadedTasks[1].Description != "Walk dog" {
		t.Errorf("Loaded tasks do not match expected: %+v", loadedTasks)
	}
}

// captureOutput redirects os.Stdout and os.Stderr to a buffer and returns the buffer content.
func captureOutput(f func()) (string, string) {
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	os.Stdout = wOut
	os.Stderr = wErr

	outC := make(chan string)
	errC := make(chan string)

	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, rOut)
		outC <- buf.String()
	}()
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, rErr)
		errC <- buf.String()
	}()

	f()

	wOut.Close()
	wErr.Close()

	os.Stdout = oldStdout
	os.Stderr = oldStderr

	return <-outC, <-errC
}

func TestAddCommand(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "todo_add_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	oldTasksFile := tasksFile
	tasksFile = filepath.Join(tempDir, "tasks.json")
	oldArgs := os.Args
	oldExit := exitFunc

	defer func() {
		tasksFile = oldTasksFile
		os.Args = oldArgs
		exitFunc = oldExit
	}()

	exitedWith := -1
	exitFunc = func(code int) {
		exitedWith = code
	}

	// Test adding a task successfully
	os.Args = []string{"todo", "add", "Buy groceries"}
	out, _ := captureOutput(main)

	if exitedWith != 0 {
		t.Errorf("Expected exit code 0, got %d for successful add", exitedWith)
	}

	var addResp struct {ID int `json:"id"`; Status string `json:"status"`}
	err = json.Unmarshal([]byte(out), &addResp)
	if err != nil {
		t.Fatalf("Failed to unmarshal add response: %v\nOutput: %s", err, out)
	}
	if addResp.ID != 1 || addResp.Status != "added" {
		t.Errorf("Expected {ID: 1, Status: \"added\"}, got %+v", addResp)
	}

	// Verify task file content
	loadedTasks, _ := loadTasks()
	if len(loadedTasks) != 1 || loadedTasks[0].Description != "Buy groceries" || loadedTasks[0].ID != 1 {
		t.Errorf("Task file content incorrect after add: %+v", loadedTasks)
	}

	// Test adding another task
	exitedWith = -1 // Reset exit code
	os.Args = []string{"todo", "add", "Pay bills"}
	out, _ = captureOutput(main)

	if exitedWith != 0 {
		t.Errorf("Expected exit code 0, got %d for successful add", exitedWith)
	}
	err = json.Unmarshal([]byte(out), &addResp)
	if err != nil {
		t.Fatalf("Failed to unmarshal add response: %v\nOutput: %s", err, out)
	}
	if addResp.ID != 2 || addResp.Status != "added" {
		t.Errorf("Expected {ID: 2, Status: \"added\"}, got %+v", addResp)
	}

	loadedTasks, _ = loadTasks()
	if len(loadedTasks) != 2 || loadedTasks[1].Description != "Pay bills" || loadedTasks[1].ID != 2 {
		t.Errorf("Task file content incorrect after second add: %+v", loadedTasks)
	}

	// Test invalid arguments for add (missing description)
	exitedWith = -1 // Reset exit code
	os.Args = []string{"todo", "add"}
	out, _ = captureOutput(main)

	if exitedWith != 42 {
		t.Errorf("Expected exit code 42, got %d for invalid add args", exitedWith)
	}
	var errResp ErrorResponse
	json.Unmarshal([]byte(out), &errResp)
	if errResp.Error != "invalid arguments" {
		t.Errorf("Expected error 'invalid arguments', got '%s'", errResp.Error)
	}
}

func TestListCommand(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "todo_list_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	oldTasksFile := tasksFile
	tasksFile = filepath.Join(tempDir, "tasks.json")
	oldArgs := os.Args
	oldExit := exitFunc

	defer func() {
		tasksFile = oldTasksFile
		os.Args = oldArgs
		exitFunc = oldExit
	}()

	exitedWith := -1
	exitFunc = func(code int) {
		exitedWith = code
	}

	// Setup some tasks
	initialTasks := []Task{
		{ID: 1, Description: "Buy milk", Done: false},
		{ID: 2, Description: "Walk dog", Done: true},
		{ID: 3, Description: "Read book", Done: false},
	}
	saveTasks(initialTasks)

	// Test listing tasks
	os.Args = []string{"todo", "list"}
	out, _ := captureOutput(main)

	if exitedWith != 0 {
		t.Errorf("Expected exit code 0, got %d for successful list", exitedWith)
	}

	var listedTasks []Task
	err = json.Unmarshal([]byte(out), &listedTasks)
	if err != nil {
		t.Fatalf("Failed to unmarshal list response: %v\nOutput: %s", err, out)
	}

	if len(listedTasks) != 2 {
		t.Errorf("Expected 2 pending tasks, got %d", len(listedTasks))
	}
	if listedTasks[0].ID != 1 || listedTasks[1].ID != 3 {
		t.Errorf("Listed tasks are not as expected: %+v", listedTasks)
	}
	if listedTasks[0].Done || listedTasks[1].Done {
		t.Errorf("Listed tasks contain completed tasks: %+v", listedTasks)
	}

	// Test listing with no tasks file (should return empty array)
	os.Remove(tasksFile) // Remove the tasks file
	exitedWith = -1
	os.Args = []string{"todo", "list"}
	out, _ = captureOutput(main)

	if exitedWith != 0 {
		t.Errorf("Expected exit code 0, got %d for empty list", exitedWith)
	}
	var emptyList []Task
	json.Unmarshal([]byte(out), &emptyList)
	if len(emptyList) != 0 {
		t.Errorf("Expected empty array for no tasks, got %+v", emptyList)
	}
}

func TestDoneCommand(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "todo_done_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	oldTasksFile := tasksFile
	tasksFile = filepath.Join(tempDir, "tasks.json")
	oldArgs := os.Args
	oldExit := exitFunc

	defer func() {
		tasksFile = oldTasksFile
		os.Args = oldArgs
		exitFunc = oldExit
	}()

	exitedWith := -1
	exitFunc = func(code int) {
		exitedWith = code
	}

	// Setup some tasks
	initialTasks := []Task{
		{ID: 1, Description: "Buy milk", Done: false},
		{ID: 2, Description: "Walk dog", Done: false},
		{ID: 3, Description: "Read book", Done: false},
	}
	saveTasks(initialTasks)

	// Test marking a task as done successfully
	os.Args = []string{"todo", "done", "2"}
	out, _ := captureOutput(main)

	if exitedWith != 0 {
		t.Errorf("Expected exit code 0, got %d for successful done", exitedWith)
	}

	var doneResp struct {ID int `json:"id"`; Status string `json:"status"`}
	err = json.Unmarshal([]byte(out), &doneResp)
	if err != nil {
		t.Fatalf("Failed to unmarshal done response: %v\nOutput: %s", err, out)
	}
	if doneResp.ID != 2 || doneResp.Status != "completed" {
		t.Errorf("Expected {ID: 2, Status: \"completed\"}, got %+v", doneResp)
	}

	// Verify task file content
	loadedTasks, _ := loadTasks()
	if !loadedTasks[1].Done {
		t.Errorf("Task 2 not marked as done: %+v", loadedTasks)
	}

	// Test marking a non-existent task as done
	exitedWith = -1 // Reset exit code
	os.Args = []string{"todo", "done", "99"}
	out, _ = captureOutput(main)

	if exitedWith != 1 {
		t.Errorf("Expected exit code 1, got %d for non-existent task", exitedWith)
	}
	var errResp ErrorResponse
	json.Unmarshal([]byte(out), &errResp)
	if errResp.Error != "task with ID 99 not found" {
		t.Errorf("Expected error 'task with ID 99 not found', got '%s'", errResp.Error)
	}

	// Test invalid arguments for done (missing ID)
	exitedWith = -1 // Reset exit code
	os.Args = []string{"todo", "done"}
	out, _ = captureOutput(main)

	if exitedWith != 42 {
		t.Errorf("Expected exit code 42, got %d for invalid done args (missing ID)", exitedWith)
	}
	json.Unmarshal([]byte(out), &errResp)
	if errResp.Error != "invalid arguments" {
		t.Errorf("Expected error 'invalid arguments', got '%s'", errResp.Error)
	}

	// Test invalid arguments for done (non-numeric ID)
	exitedWith = -1 // Reset exit code
	os.Args = []string{"todo", "done", "abc"}
	out, _ = captureOutput(main)

	if exitedWith != 42 {
		t.Errorf("Expected exit code 42, got %d for invalid done args (non-numeric ID)", exitedWith)
	}
	json.Unmarshal([]byte(out), &errResp)
	if errResp.Error != "invalid arguments" {
		t.Errorf("Expected error 'invalid arguments', got '%s'", errResp.Error)
	}
}

// exitFunc is a variable to allow mocking os.Exit in tests.
var exitFunc = os.Exit

// mainWrapper calls main and captures its exit code.
func mainWrapper() int {
	oldExit := exitFunc
	defer func() { exitFunc = oldExit }()

	exitCode := 0 // Default success
	exitFunc = func(code int) {
		exitCode = code
		panic("os.Exit was called") // Panic to stop execution in test
	}

	defer func() {
		if r := recover(); r != nil {
			if r != "os.Exit was called" {
				panic(r)
			}
		}
	}()

	main()
	return exitCode
}
