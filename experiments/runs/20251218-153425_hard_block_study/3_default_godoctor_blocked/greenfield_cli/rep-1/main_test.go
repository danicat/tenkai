package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
)

func setupTest(t *testing.T) {
	TasksFile = "test_tasks.json"
	err := ioutil.WriteFile(TasksFile, []byte("[]"), 0644)
	if err != nil {
		t.Fatalf("failed to create test_tasks.json: %v", err)
	}
}

func teardownTest(t *testing.T) {
	if _, err := os.Stat(TasksFile); err == nil {
		err := os.Remove(TasksFile)
		if err != nil {
			t.Fatalf("failed to remove test_tasks.json: %v", err)
		}
	}
}

func TestLoadAndSaveTasks(t *testing.T) {
	setupTest(t)
	defer teardownTest(t)

	tests := []TodoTask{
		{ID: 1, Description: "Task 1", Done: false},
		{ID: 2, Description: "Task 2", Done: true},
	}

	err := SaveTasks(tests)
	if err != nil {
		t.Fatalf("SaveTasks failed: %v", err)
	}

	loadedTasks, err := LoadTasks()
	if err != nil {
		t.Fatalf("LoadTasks failed: %v", err)
	}

	if !reflect.DeepEqual(tests, loadedTasks) {
		t.Errorf("loaded tasks mismatch. Expected %v, got %v", tests, loadedTasks)
	}

	// Test loading from non-existent file
	os.Remove(TasksFile)
	loadedTasks, err = LoadTasks()
	if err != nil {
		t.Fatalf("LoadTasks failed for non-existent file: %v", err)
	}
	if len(loadedTasks) != 0 {
		t.Errorf("Expected empty tasks for non-existent file, got %v", loadedTasks)
	}
}

func TestAddCommand(t *testing.T) {
	setupTest(t)
	defer teardownTest(t)

	var buf bytes.Buffer
	exitCode := runCLI([]string{"todo", "add", "Buy bread"}, &buf)
	if exitCode != 0 {
		t.Fatalf("Add command exited with code %d, output: %s", exitCode, buf.String())
	}

	loadedTasks, err := LoadTasks()
	if err != nil {
		t.Fatalf("LoadTasks failed: %v", err)
	}
	if len(loadedTasks) != 1 || loadedTasks[0].Description != "Buy bread" || loadedTasks[0].ID != 1 {
		t.Errorf("Add command failed. Expected 1 task with description 'Buy bread' and ID 1, got %v", loadedTasks)
	}

	buf.Reset()
	exitCode = runCLI([]string{"todo", "add", "Pay bills"}, &buf)
	if exitCode != 0 {
		t.Fatalf("Add command exited with code %d, output: %s", exitCode, buf.String())
	}
	loadedTasks, err = LoadTasks()
	if err != nil {
		t.Fatalf("LoadTasks failed: %v", err)
	}
	if len(loadedTasks) != 2 || loadedTasks[1].Description != "Pay bills" || loadedTasks[1].ID != 2 {
		t.Errorf("Add command failed with multiple tasks. Expected 2 tasks, second with 'Pay bills' and ID 2, got %v", loadedTasks)
	}

	// Test invalid arguments
	buf.Reset()
	exitCode = runCLI([]string{"todo", "add"}, &buf)
	if exitCode != 42 {
		t.Errorf("Expected exit code 42 for invalid add arguments, got %d", exitCode)
	}
	var output map[string]string
	err = json.Unmarshal(buf.Bytes(), &output)
	if err != nil || output["error"] != "invalid arguments" {
		t.Errorf("Expected \"invalid arguments\" error, got %s", buf.String())
	}
}

func TestListCommand(t *testing.T) {
	setupTest(t)
	defer teardownTest(t)

	tasks := []TodoTask{
		{ID: 1, Description: "Task 1", Done: false},
		{ID: 2, Description: "Task 2", Done: true},
		{ID: 3, Description: "Task 3", Done: false},
	}
	SaveTasks(tasks)

	var buf bytes.Buffer
	exitCode := runCLI([]string{"todo", "list"}, &buf)
	if exitCode != 0 {
		t.Fatalf("List command exited with code %d, output: %s", exitCode, buf.String())
	}

	var listedTasks []TodoTask
	err := json.Unmarshal(buf.Bytes(), &listedTasks)
	if err != nil {
		t.Fatalf("failed to unmarshal list output: %v", err)
	}

	expected := []TodoTask{
		{ID: 1, Description: "Task 1", Done: false},
		{ID: 3, Description: "Task 3", Done: false},
	}

	if !reflect.DeepEqual(listedTasks, expected) {
		t.Errorf("List command failed. Expected %v, got %v", expected, listedTasks)
	}
}

func TestDoneCommand(t *testing.T) {
	setupTest(t)
	defer teardownTest(t)

	tasks := []TodoTask{
		{ID: 1, Description: "Task 1", Done: false},
		{ID: 2, Description: "Task 2", Done: false},
	}
	SaveTasks(tasks)

	var buf bytes.Buffer
	exitCode := runCLI([]string{"todo", "done", "1"}, &buf)
	if exitCode != 0 {
		t.Fatalf("Done command exited with code %d, output: %s", exitCode, buf.String())
	}

	loadedTasks, err := LoadTasks()
	if err != nil {
		t.Fatalf("LoadTasks failed: %v", err)
	}
	if !loadedTasks[0].Done {
		t.Errorf("Done command failed. Task 1 should be done, but is %v", loadedTasks[0].Done)
	}
	if loadedTasks[1].Done {
		t.Errorf("Done command failed. Task 2 should not be done, but is %v", loadedTasks[1].Done)
	}

	// Test non-existent ID
	buf.Reset()
	exitCode = runCLI([]string{"todo", "done", "3"}, &buf)
	if exitCode != 42 {
		t.Errorf("Expected exit code 42 for non-existent ID, got %d", exitCode)
	}
	var output map[string]string
	err = json.Unmarshal(buf.Bytes(), &output)
	if err != nil || output["error"] == "" {
		t.Fatalf("failed to unmarshal error output or error is empty: %v, output: %s", err, buf.String())
	}
	if !strings.Contains(output["error"], "task with ID 3 not found") {
		t.Errorf("Expected \"task with ID 3 not found\" error, got %s", output["error"])
	}

	// Test invalid ID argument
	buf.Reset()
	exitCode = runCLI([]string{"todo", "done", "abc"}, &buf)
	if exitCode != 42 {
		t.Errorf("Expected exit code 42 for invalid ID format, got %d", exitCode)
	}
	err = json.Unmarshal(buf.Bytes(), &output)
	if err != nil || output["error"] != "invalid task ID" {
		t.Errorf("Expected \"invalid task ID\" error, got %s", buf.String())
	}

	// Test missing ID argument
	buf.Reset()
	exitCode = runCLI([]string{"todo", "done"}, &buf)
	if exitCode != 42 {
		t.Errorf("Expected exit code 42 for missing ID argument, got %d", exitCode)
	}
	err = json.Unmarshal(buf.Bytes(), &output)
	if err != nil || output["error"] != "invalid arguments" {
		t.Errorf("Expected \"invalid arguments\" error, got %s", buf.String())
	}
}

func TestUnknownCommand(t *testing.T) {
	setupTest(t)
	defer teardownTest(t)

	var buf bytes.Buffer
	exitCode := runCLI([]string{"todo", "unknown"}, &buf)
	if exitCode != 42 {
		t.Errorf("Expected exit code 42 for unknown command, got %d", exitCode)
	}
	var output map[string]string
	err := json.Unmarshal(buf.Bytes(), &output)
	if err != nil || output["error"] != "unknown command" {
		t.Errorf("Expected \"unknown command\" error, got %s", buf.String())
	}
}

func TestNoArguments(t *testing.T) {
	setupTest(t)
	defer teardownTest(t)

	var buf bytes.Buffer
	exitCode := runCLI([]string{"todo"}, &buf)
	if exitCode != 42 {
		t.Errorf("Expected exit code 42 for no arguments, got %d", exitCode)
	}
	var output map[string]string
	err := json.Unmarshal(buf.Bytes(), &output)
	if err != nil || output["error"] != "invalid arguments" {
		t.Errorf("Expected \"invalid arguments\" error, got %s", buf.String())
	}
}
