package main

import (
	"fmt"
	"os"
	"sync"
	"testing"
)

func TestMain(m *testing.M) {
	// Run tests
	exitCode := m.Run()

	// Clean up after tests
	os.Remove(tasksFile)
	os.Exit(exitCode)
}

func TestLoadTasks_FileDoesNotExist(t *testing.T) {
	// Ensure tasks.json does not exist before test
	os.Remove(tasksFile)

	tl := NewTodoList()
	err := tl.LoadTasks()
	if err != nil {
		t.Fatalf("Expected no error for non-existent file, got: %v", err)
	}

	if len(tl.tasks) != 0 {
		t.Errorf("Expected empty tasks, got %d tasks", len(tl.tasks))
	}
}

func TestSaveAndLoadTasks(t *testing.T) {
	// Ensure tasks.json does not exist before test
	os.Remove(tasksFile)

	tl := NewTodoList()

	// Add some tasks
	_, _ = tl.AddTask("Task 1")
	_, _ = tl.AddTask("Task 2")

	// Save tasks
	err := tl.SaveTasks()
	if err != nil {
		t.Fatalf("Failed to save tasks: %v", err)
	}

	// Create a new TodoList and load tasks
	newTl := NewTodoList()
	err = newTl.LoadTasks()
	if err != nil {
		t.Fatalf("Failed to load tasks: %v", err)
	}

	// Verify loaded tasks
	if len(newTl.tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(newTl.tasks))
	}

	if newTl.tasks[0].Description != "Task 1" || newTl.tasks[1].Description != "Task 2" {
		t.Errorf("Loaded tasks do not match original")
	}
}

func TestAddTask(t *testing.T) {
	// Ensure tasks.json does not exist before test
	os.Remove(tasksFile)

	tl := NewTodoList()

	task1, err := tl.AddTask("First task")
	if err != nil {
		t.Fatalf("Failed to add task: %v", err)
	}

	if task1.ID != 1 || task1.Description != "First task" || task1.Done != false {
		t.Errorf("Unexpected task1: %+v", task1)
	}

	task2, err := tl.AddTask("Second task")
	if err != nil {
		t.Fatalf("Failed to add task: %v", err)
	}

	if task2.ID != 2 || task2.Description != "Second task" || task2.Done != false {
		t.Errorf("Unexpected task2: %+v", task2)
	}

	if len(tl.tasks) != 2 {
		t.Errorf("Expected 2 tasks in list, got %d", len(tl.tasks))
	}
}

func TestListPendingTasks(t *testing.T) {
	// Ensure tasks.json does not exist before test
	os.Remove(tasksFile)

	tl := NewTodoList()
	task1, _ := tl.AddTask("Task 1")
	task2, _ := tl.AddTask("Task 2")
	task3, _ := tl.AddTask("Task 3")

	// Complete task2
	_, err := tl.CompleteTask(task2.ID)
	if err != nil {
		t.Fatalf("Failed to complete task: %v", err)
	}

	pendingTasks := tl.ListPendingTasks()

	if len(pendingTasks) != 2 {
		t.Errorf("Expected 2 pending tasks, got %d", len(pendingTasks))
	}

	if pendingTasks[0].ID != task1.ID || pendingTasks[1].ID != task3.ID {
		t.Errorf("Unexpected pending tasks: %+v", pendingTasks)
	}
}

func TestCompleteTask(t *testing.T) {
	// Ensure tasks.json does not exist before test
	os.Remove(tasksFile)

	tl := NewTodoList()
	task1, _ := tl.AddTask("Task to complete")
	_, _ = tl.AddTask("Another task")

	completedTask, err := tl.CompleteTask(task1.ID)
	if err != nil {
		t.Fatalf("Failed to complete task: %v", err)
	}

	if completedTask.ID != task1.ID || completedTask.Done != true {
		t.Errorf("Expected task1 to be completed, got: %+v", completedTask)
	}

	// Verify the task in the internal list is also updated
	found := false
	for _, task := range tl.tasks {
		if task.ID == task1.ID {
			found = true
			if !task.Done {
				t.Errorf("Task in list not marked as done after completion")
			}
			break
		}
	}
	if !found {
		t.Errorf("Completed task not found in the list")
	}

	// Test completing a non-existent task
	_, err = tl.CompleteTask(999)
	if err == nil {
		t.Errorf("Expected error for non-existent task, got nil")
	}
	if err.Error() != "task not found" {
		t.Errorf("Expected 'task not found' error, got: %v", err)
	}
}

func TestConcurrency(t *testing.T) {
	// Ensure tasks.json does not exist before test
	os.Remove(tasksFile)

	tl := NewTodoList()

	// Add initial tasks
	_, _ = tl.AddTask("Initial task 1")
	_, _ = tl.AddTask("Initial task 2")

	// Save initial tasks
	err := tl.SaveTasks()
	if err != nil {
		t.Fatalf("Failed to save initial tasks: %v", err)
	}

	numGoroutines := 10
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			defer wg.Done()

			// Each goroutine tries to add a task, complete a task, list tasks, and save
			_, err := tl.AddTask(fmt.Sprintf("Concurrent task %d", i))
			if err != nil {
				t.Errorf("Goroutine %d: Failed to add task: %v", i, err)
			}

			// Try to complete an existing task (e.g., the first one added, ID 1 or 2)
			// This might fail if another goroutine completes it, but we're testing mutex protection for data consistency
			_, err = tl.CompleteTask(1)
			if err != nil && err.Error() != "task not found" { // task not found is expected if already completed
				t.Errorf("Goroutine %d: Failed to complete task: %v", i, err)
			}

			_ = tl.ListPendingTasks() // Just call to ensure no panic

			err = tl.SaveTasks()
			if err != nil {
				t.Errorf("Goroutine %d: Failed to save tasks: %v", i, err)
			}
		}(i)
	}

	wg.Wait()

	// Load tasks again to check final state consistency
	finalTl := NewTodoList()
	err = finalTl.LoadTasks()
	if err != nil {
		t.Fatalf("Failed to load tasks after concurrency test: %v", err)
	}

	// Check if the number of tasks is as expected (initial 2 + numGoroutines added)
	if len(finalTl.tasks) != 2 + numGoroutines {
		t.Errorf("Expected %d tasks after concurrency, got %d", 2+numGoroutines, len(finalTl.tasks))
	}

	// Verify all tasks have unique IDs and descriptions
	ids := make(map[int]bool)
	descriptions := make(map[string]bool)
	for _, task := range finalTl.tasks {
		if ids[task.ID] {
			t.Errorf("Duplicate task ID found: %d", task.ID)
		}
		ids[task.ID] = true

		if descriptions[task.Description] {
			t.Errorf("Duplicate task description found: %s", task.Description)
		}
		descriptions[task.Description] = true
	}

	// Ensure at least one task is marked done (from completing task ID 1)
	oneTaskDone := false
	for _, task := range finalTl.tasks {
		if task.ID == 1 && task.Done {
			oneTaskDone = true
			break
		}
	}
	if !oneTaskDone {
		t.Errorf("Expected task ID 1 to be marked done by at least one goroutine")
	}
}
