package todo_test

import (
	"reflect"
	"sync"
	"testing"
	"todo-cli/internal/todo"
)

// MockStorer for testing the Manager.
type MockStorer struct {
	mu    sync.Mutex
	tasks []*todo.Task
	loadErr error
	saveErr error
}

func NewMockStorer() *MockStorer {
	return &MockStorer{
		tasks: []*todo.Task{},
	}
}

func (m *MockStorer) Load() ([]*todo.Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.tasks, m.loadErr
}

func (m *MockStorer) Save(tasks []*todo.Task) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tasks = tasks
	return m.saveErr
}

func TestManager_Add(t *testing.T) {

storer := NewMockStorer()
	manager, err := todo.NewManager(storer)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	desc := "Test task"
	task, err := manager.Add(desc)
	if err != nil {
		t.Errorf("Add failed: %v", err)
	}

	if task.ID != 1 {
		t.Errorf("Expected ID 1, got %d", task.ID)
	}
	if task.Description != desc {
		t.Errorf("Expected description %s, got %s", desc, task.Description)
	}
	if task.Done {
		t.Errorf("Expected task to be not done, but it is")
	}

	// Verify it's in the list
	tasks := manager.List()
	if len(tasks) != 1 || tasks[0].ID != 1 {
		t.Errorf("Expected 1 task in list with ID 1, got %+v", tasks)
	}
}

func TestManager_List(t *testing.T) {

storer := NewMockStorer()
	manager, err := todo.NewManager(storer)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	manager.Add("Task A")
	manager.Add("Task B")

	tasks := manager.List()

	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}

	expectedDescriptions := []string{"Task A", "Task B"}
	for i, task := range tasks {
		if task.Description != expectedDescriptions[i] {
			t.Errorf("Expected description %s at index %d, got %s", expectedDescriptions[i], i, task.Description)
		}
	}
}

func TestManager_Complete(t *testing.T) {

storer := NewMockStorer()
	manager, err := todo.NewManager(storer)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	manager.Add("Task to complete")
	manager.Add("Another task")

	completedTask, err := manager.Complete(1)
	if err != nil {
		t.Errorf("Complete failed: %v", err)
	}

	if completedTask.ID != 1 || !completedTask.Done {
		t.Errorf("Expected task 1 to be done, got %+v", completedTask)
	}

	// Verify in the list
	tasks := manager.List()
	if !tasks[0].Done {
		t.Errorf("Task 1 in list should be done, but isn't")
	}
	if tasks[1].Done {
		t.Errorf("Task 2 in list should not be done, but is")
	}

	// Test completing a non-existent task
	_, err = manager.Complete(999)
	if err == nil {
		t.Errorf("Expected error when completing non-existent task, got nil")
	}

	// Test completing an already completed task
	_, err = manager.Complete(1)
	if err != nil {
		t.Errorf("Completing already done task should not return an error: %v", err)
	}

}

func TestManager_LoadAndSave(t *testing.T) {
	// Test initial load with empty storer

storer := NewMockStorer()
	manager, err := todo.NewManager(storer)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	if len(manager.List()) != 0 {
		t.Errorf("Expected empty task list on initial load, got %d tasks", len(manager.List()))
	}

	// Add some tasks and save
	manager.Add("Task 1")
	manager.Add("Task 2")

	// Create a new manager with the same storer to simulate reloading
	manager2, err := todo.NewManager(storer)
	if err != nil {
		t.Fatalf("NewManager (second) failed: %v", err)
	}

	tasks2 := manager2.List()
	if len(tasks2) != 2 {
		t.Errorf("Expected 2 tasks after reload, got %d", len(tasks2))
	}

	expectedTasks := []*todo.Task{
		{ID: 1, Description: "Task 1", Done: false},
		{ID: 2, Description: "Task 2", Done: false},
	}

	// Check if the loaded tasks match the saved tasks
	if !reflect.DeepEqual(tasks2, expectedTasks) {
		t.Errorf("Reloaded tasks mismatch.\nGot: %+v\nWant: %+v", tasks2, expectedTasks)
	}

	// Complete a task in manager2 and check if it persists
	manager2.Complete(1)

	manager3, err := todo.NewManager(storer)
	if err != nil {
		t.Fatalf("NewManager (third) failed: %v", err)
	}

	tasks3 := manager3.List()
	if len(tasks3) != 2 {
		t.Errorf("Expected 2 tasks after third reload, got %d", len(tasks3))
	}

	if !tasks3[0].Done {
		t.Errorf("Task 1 should be done after reload, but isn't")
	}
	if tasks3[1].Done {
		t.Errorf("Task 2 should not be done after reload, but is")
	}


}
