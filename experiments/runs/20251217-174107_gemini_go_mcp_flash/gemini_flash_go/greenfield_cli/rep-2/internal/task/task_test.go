package task

import (
	"os"
	"testing"
)

var testTasksFile = "test_tasks.json"

func TestMain(m *testing.M) {
	tasksFile = testTasksFile // Redirect to a test-specific file
	code := m.Run()
	os.Remove(testTasksFile) // Clean up after tests
	os.Exit(code)
}

func TestNewStore(t *testing.T) {
	os.Remove(testTasksFile) // Ensure no old test file exists
	s, err := NewStore()
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}

	if len(s.tasks) != 0 {
		t.Errorf("Expected empty tasks on new store, got %d", len(s.tasks))
	}
	if s.nextID != 1 {
		t.Errorf("Expected nextID to be 1, got %d", s.nextID)
	}

	// Test loading from an existing file
	_, _ = s.AddTask("Task 1")
	_, _ = s.AddTask("Task 2")

	s2, err := NewStore()
	if err != nil {
		t.Fatalf("NewStore (existing) failed: %v", err)
	}

	if len(s2.tasks) != 2 {
		t.Errorf("Expected 2 tasks after loading, got %d", len(s2.tasks))
	}
	if s2.nextID != 3 {
		t.Errorf("Expected nextID to be 3 after loading, got %d", s2.nextID)
	}
}

func TestAddTask(t *testing.T) {
	os.Remove(testTasksFile)
	s, _ := NewStore()

	desc1 := "Buy groceries"
	t1, err := s.AddTask(desc1)
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}
	if t1.ID != 1 || t1.Description != desc1 || t1.Done != false {
		t.Errorf("Unexpected task 1: %+v", t1)
	}

	desc2 := "Pay bills"
	t2, err := s.AddTask(desc2)
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}
	if t2.ID != 2 || t2.Description != desc2 || t2.Done != false {
		t.Errorf("Unexpected task 2: %+v", t2)
	}

	if len(s.tasks) != 2 {
		t.Errorf("Expected 2 tasks in store, got %d", len(s.tasks))
	}
}

func TestListTasks(t *testing.T) {
	os.Remove(testTasksFile)
	s, _ := NewStore()

	_, _ = s.AddTask("Task A")
	_, _ = s.AddTask("Task B")
	_, _ = s.CompleteTask(1)

	pending := s.ListTasks(false)
	if len(pending) != 1 || pending[0].ID != 2 {
		t.Errorf("Expected 1 pending task (ID 2), got: %+v", pending)
	}

	all := s.ListTasks(true)
	if len(all) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(all))
	}

	// Test empty list
	os.Remove(testTasksFile)
	s3, _ := NewStore()
	if len(s3.ListTasks(false)) != 0 {
		t.Errorf("Expected empty list for new store, got %d", len(s3.ListTasks(false)))
	}
}

func TestCompleteTask(t *testing.T) {
	os.Remove(testTasksFile)
	s, _ := NewStore()
	_, _ = s.AddTask("Task to complete")

	// Complete existing task
	completedTask, err := s.CompleteTask(1)
	if err != nil {
		t.Fatalf("CompleteTask failed: %v", err)
	}
	if completedTask.ID != 1 || !completedTask.Done {
		t.Errorf("Expected task 1 to be done, got %+v", completedTask)
	}

	// Try to complete already completed task
	_, err = s.CompleteTask(1)
	if err == nil || err.Error() != "task 1 is already completed" {
		t.Errorf("Expected 'already completed' error, got: %v", err)
	}

	// Try to complete non-existent task
	_, err = s.CompleteTask(99)
	if err == nil || err.Error() != "task with ID 99 not found" {
		t.Errorf("Expected 'task not found' error, got: %v", err)
	}
}

func TestPersistence(t *testing.T) {
	os.Remove(testTasksFile)
	s1, _ := NewStore()
	_, _ = s1.AddTask("Persistent task 1")
	_, _ = s1.AddTask("Persistent task 2")
	_, _ = s1.CompleteTask(1)

	s2, err := NewStore() // Load from file
	if err != nil {
		t.Fatalf("NewStore failed during persistence test: %v", err)
	}

	if len(s2.tasks) != 2 {
		t.Errorf("Expected 2 tasks after reload, got %d", len(s2.tasks))
	}

	if !s2.tasks[0].Done {
		t.Errorf("Expected task 1 to be done after reload")
	}
	if s2.nextID != 3 {
		t.Errorf("Expected nextID to be 3 after reload, got %d", s2.nextID)
	}

	// Add another task after reload
	newlyAdded, _ := s2.AddTask("Newly added task")
	if newlyAdded.ID != 3 {
		t.Errorf("Expected newly added task ID to be 3, got %d", newlyAdded.ID)
	}

	// Verify persistence again after another add
	s3, _ := NewStore()
	if len(s3.tasks) != 3 {
		t.Errorf("Expected 3 tasks after second reload, got %d", len(s3.tasks))
	}
}

func TestErrorHandling(t *testing.T) {
	// Test unmarshalling error on load
	os.WriteFile(testTasksFile, []byte("invalid json"), 0644)
	_, err := NewStore()
	if err == nil {
		t.Errorf("Expected error for invalid JSON, got nil")
	}

	// Ensure task file is created if not exists
	os.Remove(testTasksFile)
	store, err := NewStore()
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}
	_, err = store.AddTask("test")
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}

	// Verify file exists
	_, err = os.Stat(testTasksFile)
	if os.IsNotExist(err) {
		t.Errorf("tasks.json was not created")
	}

}
