package todo

import (
	"io/ioutil"
	"os"
	"testing"
)

func setupTest(t *testing.T) {
	// Ensure tasks.json doesn't exist before each test
	os.Remove(tasksFile)
}

func teardownTest(t *testing.T) {
	// Clean up tasks.json after each test
	os.Remove(tasksFile)
}

func TestAddTask(t *testing.T) {
	setupTest(t)
	defer teardownTest(t)

	task, err := AddTask("Buy groceries")
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}

	if task.ID != 1 || task.Description != "Buy groceries" || task.Done != false {
		t.Errorf("Unexpected task data: %+v", task)
	}

	tasks, err := loadTasks()
	if err != nil {
		t.Fatalf("loadTasks failed: %v", err)
	}
	if len(tasks) != 1 || tasks[0].ID != 1 {
		t.Errorf("Expected 1 task with ID 1, got %+v", tasks)
	}

	task2, err := AddTask("Walk the dog")
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}

	if task2.ID != 2 || task2.Description != "Walk the dog" || task2.Done != false {
		t.Errorf("Unexpected task data: %+v", task2)
	}

	tasks, err = loadTasks()
	if err != nil {
		t.Fatalf("loadTasks failed: %v", err)
	}
	if len(tasks) != 2 || tasks[1].ID != 2 {
		t.Errorf("Expected 2 tasks with IDs 1 and 2, got %+v", tasks)
	}
}

func TestListTasks(t *testing.T) {
	setupTest(t)
	defer teardownTest(t)

	// Test with no tasks
	tasks, err := ListTasks()
	if err != nil {
		t.Fatalf("ListTasks failed: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks, got %d", len(tasks))
	}

	// Add some tasks
	AddTask("Task 1")
	AddTask("Task 2")

	tasks, err = ListTasks()
	if err != nil {
		t.Fatalf("ListTasks failed: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}
	if tasks[0].Description != "Task 1" || tasks[1].Description != "Task 2" {
		t.Errorf("Unexpected task descriptions: %+v", tasks)
	}
}

func TestCompleteTask(t *testing.T) {
	setupTest(t)
	defer teardownTest(t)

	AddTask("Task to complete") // ID 1
	AddTask("Another task")     // ID 2

	// Complete an existing task
	completedTask, err := CompleteTask(1)
	if err != nil {
		t.Fatalf("CompleteTask failed: %v", err)
	}
	if completedTask.ID != 1 || completedTask.Done != true {
		t.Errorf("Expected task 1 to be done, got %+v", completedTask)
	}

	tasks, err := loadTasks()
	if err != nil {
		t.Fatalf("loadTasks failed: %v", err)
	}
	if !tasks[0].Done {
		t.Errorf("Task 1 not marked as done in storage")
	}

	// Try to complete a non-existent task
	_, err = CompleteTask(99)
	if err == nil || err.Error() != "task not found" {
		t.Errorf("Expected 'task not found' error, got: %v", err)
	}
}

func TestPersistence(t *testing.T) {
	setupTest(t)
	defer teardownTest(t)

	AddTask("Persistent task 1")
	AddTask("Persistent task 2")

	// Simulate closing and reopening the application by clearing in-memory state
	// and then reloading from file.
	// In a real application, this would involve restarting the process.

	tasksAfterRestart, err := loadTasks()
	if err != nil {
		t.Fatalf("loadTasks failed after simulated restart: %v", err)
	}

	if len(tasksAfterRestart) != 2 {
		t.Errorf("Expected 2 tasks after restart, got %d", len(tasksAfterRestart))
	}

	if tasksAfterRestart[0].Description != "Persistent task 1" || tasksAfterRestart[1].Description != "Persistent task 2" {
		t.Errorf("Tasks data corrupted after restart: %+v", tasksAfterRestart)
	}
}
