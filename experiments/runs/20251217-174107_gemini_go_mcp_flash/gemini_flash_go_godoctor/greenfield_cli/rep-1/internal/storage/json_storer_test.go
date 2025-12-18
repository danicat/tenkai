package storage_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"todo-cli/internal/storage"
	"todo-cli/internal/todo"
)

func TestJSONStorer_SaveAndLoad(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := ioutil.TempFile("", "tasks-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	filePath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(filePath) // Clean up the temporary file

	storer := storage.NewJSONStorer(filePath)

	// 1. Test saving an empty list
	emptyTasks := []*todo.Task{}
	err = storer.Save(emptyTasks)
	if err != nil {
		t.Errorf("Save empty list failed: %v", err)
	}

	loadedEmptyTasks, err := storer.Load()
	if err != nil {
		t.Errorf("Load empty list failed: %v", err)
	}
	if len(loadedEmptyTasks) != 0 {
		t.Errorf("Expected 0 tasks, got %d after loading empty list", len(loadedEmptyTasks))
	}

	// 2. Test saving and loading with tasks
	expectedTasks := []*todo.Task{
		{ID: 1, Description: "Buy milk", Done: false},
		{ID: 2, Description: "Walk dog", Done: true},
	}

	err = storer.Save(expectedTasks)
	if err != nil {
		t.Errorf("Save tasks failed: %v", err)
	}

	loadedTasks, err := storer.Load()
	if err != nil {
		t.Errorf("Load tasks failed: %v", err)
	}

	if !reflect.DeepEqual(loadedTasks, expectedTasks) {
		t.Errorf("Loaded tasks mismatch.\nGot: %+v\nWant: %+v", loadedTasks, expectedTasks)
	}

	// 3. Test loading from a non-existent file (should return empty list, no error)
	nonExistentFile := filepath.Join(os.TempDir(), "non-existent-tasks.json")
	defer os.Remove(nonExistentFile)
	storerNonExistent := storage.NewJSONStorer(nonExistentFile)

	loadedFromNonExistent, err := storerNonExistent.Load()
	if err != nil {
		t.Errorf("Loading from non-existent file should not return an error: %v", err)
	}
	if len(loadedFromNonExistent) != 0 {
		t.Errorf("Expected 0 tasks from non-existent file, got %d", len(loadedFromNonExistent))
	}

	// 4. Test loading from a malformed JSON file
	malformedFile := filepath.Join(os.TempDir(), "malformed.json")
	defer os.Remove(malformedFile)
	err = ioutil.WriteFile(malformedFile, []byte("invalid json"), 0644)
	if err != nil {
		t.Fatalf("Failed to write malformed JSON: %v", err)
	}
	storerMalformed := storage.NewJSONStorer(malformedFile)

	_, err = storerMalformed.Load()
	if err == nil {
		t.Errorf("Expected error when loading malformed JSON, got nil")
	}

	// 5. Test with empty file content
	emptyContentFile := filepath.Join(os.TempDir(), "empty_content.json")
	defer os.Remove(emptyContentFile)
	err = ioutil.WriteFile(emptyContentFile, []byte(""), 0644)
	if err != nil {
		t.Fatalf("Failed to write empty content file: %v", err)
	}
	storerEmptyContent := storage.NewJSONStorer(emptyContentFile)
	loadedEmptyContent, err := storerEmptyContent.Load()
	if err != nil {
		t.Errorf("Loading from empty content file should not return an error: %v", err)
	}
	if len(loadedEmptyContent) != 0 {
		t.Errorf("Expected 0 tasks from empty content file, got %d", len(loadedEmptyContent))
	}
}
