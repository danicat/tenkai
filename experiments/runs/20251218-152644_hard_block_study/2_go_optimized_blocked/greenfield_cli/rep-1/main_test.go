package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

// setupTest prepares the environment for a test.
func setupTest() (*bytes.Buffer, func()) {
	// Clear existing tasks and reset nextID
	tasksLock.Lock()
	tasks = []Task{}
	nextID = 1
	tasksLock.Unlock()

	// Remove tasks.json if it exists
	os.Remove(tasksFile)

	// Redirect stdout to capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	returnBuffer := new(bytes.Buffer)
	done := make(chan struct{})
	go func() {
		defer close(done)
		io.Copy(returnBuffer, r)
	}()

	return returnBuffer, func() {
		w.Close()
		<-done
		os.Stdout = oldStdout // Restore original stdout
		os.Remove(tasksFile)
		// Reset tasks after each test to ensure isolation
		tasksLock.Lock()
		tasks = []Task{}
		nextID = 1
		tasksLock.Unlock()
	}
}

// runCLIWrapper provides a consistent way to call runCLI in tests
func runCLIWrapper(args []string) error {
	// Create a new FlagSet for each test run to ensure isolation
	// The first arg is usually the program name, but here we treat it as the command.
	// So, we use "todo" as the program name for the FlagSet.
	fset := flag.NewFlagSet("todo", flag.ContinueOnError)
	fset.SetOutput(ioutil.Discard) // Suppress flag error output

	return runCLI(fset, args)
}

func TestAddTask(t *testing.T) {
	outputBuffer, teardown := setupTest()
	defer teardown()

	err := runCLIWrapper([]string{"add", "Buy groceries"})
	if err != nil {
		t.Fatalf("TestAddTask failed with error: %v", err)
	}

	expected := `{"id":1,"status":"added"}`
	actual := strings.TrimSpace(outputBuffer.String())
	if actual != expected {
		t.Errorf("TestAddTask output failed: Expected %s, Got %s", expected, actual)
	}

	// Verify task is saved
	loadTasks()
	if len(tasks) != 1 || tasks[0].Description != "Buy groceries" || tasks[0].ID != 1 {
		t.Errorf("TestAddTask failed: Task not saved correctly, tasks: %+v", tasks)
	}
}

func TestAddEmptyTask(t *testing.T) {
	outputBuffer, teardown := setupTest()
	defer teardown()

	err := runCLIWrapper([]string{"add", ""})

	var cerr *cliError
	if !errors.As(err, &cerr) {
		t.Fatalf("TestAddEmptyTask did not return cliError, got: %v", err)
	}

	if cerr.exitCode != 42 {
		t.Errorf("TestAddEmptyTask exit code failed: Expected 42, Got %d", cerr.exitCode)
	}

	expectedErrOutput := `{"error": "description cannot be empty"}`
	var simulatedOutputBuffer bytes.Buffer
	simulatedOutputBuffer.WriteString(fmt.Sprintf("{\"error\": %q}\n", cerr.Error()))

	if strings.TrimSpace(simulatedOutputBuffer.String()) != expectedErrOutput {
		t.Errorf("TestAddEmptyTask error output failed: Expected %s, Got %s", expectedErrOutput, strings.TrimSpace(outputBuffer.String()))
	}
}

func TestListTasks(t *testing.T) {
	outputBuffer, teardown := setupTest()
	defer teardown()

	// Add some tasks first (using the direct function, not CLI run)
	err := addTask("Task 1")
	if err != nil {
		t.Fatalf("Failed to add task 1: %v", err)
	}
	err = addTask("Task 2")
	if err != nil {
		t.Fatalf("Failed to add task 2: %v", err)
	}

	outputBuffer.Reset() // Clear output from addTask calls

	err = runCLIWrapper([]string{"list"})
	if err != nil {
		t.Fatalf("TestListTasks failed with error: %v", err)
	}

	expected := `[{"id":1,"description":"Task 1","done":false},{"id":2,"description":"Task 2","done":false}]`
	actual := strings.TrimSpace(outputBuffer.String())

	var expectedTasks, actualTasks []Task
	json.Unmarshal([]byte(expected), &expectedTasks)
	json.Unmarshal([]byte(actual), &actualTasks)

	if len(expectedTasks) != len(actualTasks) {
		t.Errorf("TestListTasks failed: Expected %d tasks, Got %d", len(expectedTasks), len(actualTasks))
	}

	for i := range expectedTasks {
		if expectedTasks[i].ID != actualTasks[i].ID ||
			expectedTasks[i].Description != actualTasks[i].Description ||
			expectedTasks[i].Done != actualTasks[i].Done {
			t.Errorf("TestListTasks failed: Mismatch at index %d. Expected %+v, Got %+v", i, expectedTasks[i], actualTasks[i])
		}
	}

	reMarshaledActual, _ := json.Marshal(actualTasks)
	if string(reMarshaledActual) != expected {
		t.Errorf("TestListTasks failed: JSON output mismatch. Expected %s, Got %s", expected, string(reMarshaledActual))
	}
}

func TestCompleteTask(t *testing.T) {
	outputBuffer, teardown := setupTest()
	defer teardown()

	// Add a task first
	err := addTask("Task to complete")
	if err != nil {
		t.Fatalf("Failed to add task: %v", err)
	}

	outputBuffer.Reset() // Clear output from addTask calls

	err = runCLIWrapper([]string{"done", "1"})
	if err != nil {
		t.Fatalf("TestCompleteTask failed with error: %v", err)
	}

	expected := `{"id":1,"status":"completed"}`
	actual := strings.TrimSpace(outputBuffer.String())
	if actual != expected {
		t.Errorf("TestCompleteTask output failed: Expected %s, Got %s", expected, actual)
	}

	// Verify task is marked as done
	loadTasks()
	if len(tasks) != 1 || !tasks[0].Done || tasks[0].ID != 1 {
		t.Errorf("TestCompleteTask failed: Task not marked as done correctly, tasks: %+v", tasks)
	}
}

func TestCompleteNonExistentTask(t *testing.T) {
	outputBuffer, teardown := setupTest()
	defer teardown()

	// Simulate command line arguments for "done" command with non-existent ID
	err := runCLIWrapper([]string{"done", "99"})

	var cerr *cliError
	if !errors.As(err, &cerr) {
		t.Fatalf("TestCompleteNonExistentTask did not return cliError, got: %v", err)
	}

	if cerr.exitCode != 42 {
		t.Errorf("TestCompleteNonExistentTask exit code failed: Expected 42, Got %d", cerr.exitCode)
	}

	expectedErrOutput := `{"error": "task with ID 99 not found"}`
	var simulatedOutputBuffer bytes.Buffer
	simulatedOutputBuffer.WriteString(fmt.Sprintf("{\"error\": %q}\n", cerr.Error()))

	if strings.TrimSpace(simulatedOutputBuffer.String()) != expectedErrOutput {
		t.Errorf("TestCompleteNonExistentTask error output failed: Expected %s, Got %s", expectedErrOutput, strings.TrimSpace(outputBuffer.String()))
	}
}

func TestInvalidCommand(t *testing.T) {
	outputBuffer, teardown := setupTest()
	defer teardown()

	err := runCLIWrapper([]string{"invalidcommand"})

	var cerr *cliError
	if !errors.As(err, &cerr) {
		t.Fatalf("TestInvalidCommand did not return cliError, got: %v", err)
	}

	if cerr.exitCode != 42 {
		t.Errorf("TestInvalidCommand exit code failed: Expected 42, Got %d", cerr.exitCode)
	}

	expectedErrOutput := `{"error": "unknown command: invalidcommand"}`
	var simulatedOutputBuffer bytes.Buffer
	simulatedOutputBuffer.WriteString(fmt.Sprintf("{\"error\": %q}\n", cerr.Error()))

	if strings.TrimSpace(simulatedOutputBuffer.String()) != expectedErrOutput {
		t.Errorf("TestInvalidCommand error output failed: Expected %s, Got %s", expectedErrOutput, strings.TrimSpace(outputBuffer.String()))
	}
}

func TestMissingCommand(t *testing.T) {
	outputBuffer, teardown := setupTest()
	defer teardown()

	err := runCLIWrapper([]string{}) // No command provided

	var cerr *cliError
	if !errors.As(err, &cerr) {
		t.Fatalf("TestMissingCommand did not return cliError, got: %v", err)
	}

	if cerr.exitCode != 42 {
		t.Errorf("TestMissingCommand exit code failed: Expected 42, Got %d", cerr.exitCode)
	}

	expectedErrOutput := `{"error": "no command provided"}`
	var simulatedOutputBuffer bytes.Buffer
	simulatedOutputBuffer.WriteString(fmt.Sprintf("{\"error\": %q}\n", cerr.Error()))

	if strings.TrimSpace(simulatedOutputBuffer.String()) != expectedErrOutput {
		t.Errorf("TestMissingCommand error output failed: Expected %s, Got %s", expectedErrOutput, strings.TrimSpace(outputBuffer.String()))
	}
}

func TestMissingAddDescription(t *testing.T) {
	outputBuffer, teardown := setupTest()
	defer teardown()

	err := runCLIWrapper([]string{"add"})

	var cerr *cliError
	if !errors.As(err, &cerr) {
		t.Fatalf("TestMissingAddDescription did not return cliError, got: %v", err)
	}

	if cerr.exitCode != 42 {
		t.Errorf("TestMissingAddDescription exit code failed: Expected 42, Got %d", cerr.exitCode)
	}

	expectedErrOutput := `{"error": "missing task description for 'add' command"}`
	var simulatedOutputBuffer bytes.Buffer
	simulatedOutputBuffer.WriteString(fmt.Sprintf("{\"error\": %q}\n", cerr.Error()))

	if strings.TrimSpace(simulatedOutputBuffer.String()) != expectedErrOutput {
		t.Errorf("TestMissingAddDescription error output failed: Expected %s, Got %s", expectedErrOutput, strings.TrimSpace(outputBuffer.String()))
	}
}

func TestMissingDoneID(t *testing.T) {
	outputBuffer, teardown := setupTest()
	defer teardown()

	err := runCLIWrapper([]string{"done"})

	var cerr *cliError
	if !errors.As(err, &cerr) {
		t.Fatalf("TestMissingDoneID did not return cliError, got: %v", err)
	}

	if cerr.exitCode != 42 {
		t.Errorf("TestMissingDoneID exit code failed: Expected 42, Got %d", cerr.exitCode)
	}

	expectedErrOutput := `{"error": "missing task ID for 'done' command"}`
	var simulatedOutputBuffer bytes.Buffer
	simulatedOutputBuffer.WriteString(fmt.Sprintf("{\"error\": %q}\n", cerr.Error()))

	if strings.TrimSpace(simulatedOutputBuffer.String()) != expectedErrOutput {
		t.Errorf("TestMissingDoneID error output failed: Expected %s, Got %s", expectedErrOutput, strings.TrimSpace(outputBuffer.String()))
	}
}

func TestInvalidDoneID(t *testing.T) {
	outputBuffer, teardown := setupTest()
	defer teardown()

	err := runCLIWrapper([]string{"done", "abc"})

	var cerr *cliError
	if !errors.As(err, &cerr) {
		t.Fatalf("TestInvalidDoneID did not return cliError, got: %v", err)
	}

	if cerr.exitCode != 42 {
		t.Errorf("TestInvalidDoneID exit code failed: Expected 42, Got %d", cerr.exitCode)
	}

	expectedErrOutput := `{"error": "invalid task ID: strconv.Atoi: parsing \"abc\": invalid syntax"}`
	var simulatedOutputBuffer bytes.Buffer
	simulatedOutputBuffer.WriteString(fmt.Sprintf("{\"error\": %q}\n", cerr.Error()))

	if strings.TrimSpace(simulatedOutputBuffer.String()) != expectedErrOutput {
		t.Errorf("TestInvalidDoneID error output failed: Expected %s, Got %s", expectedErrOutput, strings.TrimSpace(outputBuffer.String()))
	}
}

func TestListNoTasks(t *testing.T) {
	outputBuffer, teardown := setupTest()
	defer teardown()

	// Ensure no tasks file exists
	os.Remove(tasksFile)

	err := runCLIWrapper([]string{"list"})
	if err != nil {
		t.Fatalf("TestListNoTasks failed with error: %v", err)
	}

	expected := `[]`
	actual := strings.TrimSpace(outputBuffer.String())

	if actual != expected {
		t.Errorf("TestListNoTasks failed: Expected %s, Got %s", expected, actual)
	}
}

func TestListAfterComplete(t *testing.T) {
	outputBuffer, teardown := setupTest()
	defer teardown()

	// Add two tasks
	addTask("Task A")
	addTask("Task B")

	// Complete Task A
	completeTask(1)

	outputBuffer.Reset()

	err := runCLIWrapper([]string{"list"})
	if err != nil {
		t.Fatalf("TestListAfterComplete failed with error: %v", err)
	}

	expected := `[{"id":2,"description":"Task B","done":false}]`
	actual := strings.TrimSpace(outputBuffer.String())

	var expectedTasks, actualTasks []Task
	json.Unmarshal([]byte(expected), &expectedTasks)
	json.Unmarshal([]byte(actual), &actualTasks)

	if len(expectedTasks) != len(actualTasks) {
		t.Errorf("TestListAfterComplete failed: Expected %d tasks, Got %d", len(expectedTasks), len(actualTasks))
	}

	for i := range expectedTasks {
		if expectedTasks[i].ID != actualTasks[i].ID ||
			expectedTasks[i].Description != actualTasks[i].Description ||
			expectedTasks[i].Done != actualTasks[i].Done {
			t.Errorf("TestListAfterComplete failed: Mismatch at index %d. Expected %+v, Got %+v", i, expectedTasks[i], actualTasks[i])
		}
	}

	reMarshaledActual, _ := json.Marshal(actualTasks)
	if string(reMarshaledActual) != expected {
		t.Errorf("TestListAfterComplete failed: JSON output mismatch. Expected %s, Got %s", expected, string(reMarshaledActual))
	}
}

func TestLoadSaveIntegrity(t *testing.T) {
	_, teardown := setupTest()
	defer teardown()

	// Add tasks via direct function calls, not CLI
	addTask("First task")
	addTask("Second task")

	// Manually clear tasks in memory and reload from file
	tasksLock.Lock()
	tasks = []Task{}
	nextID = 1
	tasksLock.Unlock()

	err := loadTasks()
	if err != nil {
		t.Fatalf("Failed to load tasks after saving: %v", err)
	}

	if len(tasks) != 2 || tasks[0].ID != 1 || tasks[1].ID != 2 {
		t.Errorf("Expected 2 tasks after reload with correct IDs, got %d, tasks: %+v", len(tasks), tasks)
	}
	if tasks[0].Description != "First task" || tasks[1].Description != "Second task" {
		t.Errorf("Tasks loaded incorrectly. Got %+v", tasks)
	}
}
