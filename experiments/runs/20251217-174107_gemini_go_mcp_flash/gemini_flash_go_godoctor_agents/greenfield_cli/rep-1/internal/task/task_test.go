package task_test

import (
	"os"
	"reflect"
	"testing"

	"todo-cli/internal/task"
)

const testTasksFile = "test_tasks.json"

func setupTest(t *testing.T) func() {
	// Ensure a clean state for each test
	os.Remove(testTasksFile)
	task.SetTasksFile(testTasksFile) // Set the tasks file for testing
	return func() {
		os.Remove(testTasksFile)
		task.SetTasksFile("tasks.json") // Reset to default
	}
}

func TestLoadAndSaveTasks(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	tests := []task.Task{
		{ID: 1, Description: "Test Task 1", Done: false},
		{ID: 2, Description: "Test Task 2", Done: true},
	}

	err := task.SaveTasks(tests)
	if err != nil {
		t.Fatalf("SaveTasks failed: %v", err)
	}

	loadedTasks, err := task.LoadTasks()
	if err != nil {
		t.Fatalf("LoadTasks failed: %v", err)
	}

	if !reflect.DeepEqual(tests, loadedTasks) {
		t.Errorf("Loaded tasks mismatch. Expected %v, got %v", tests, loadedTasks)
	}
}

func TestLoadTasksEmptyFile(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	// Ensure file doesn't exist
	os.Remove(testTasksFile)

	loadedTasks, err := task.LoadTasks()
	if err != nil {
		t.Fatalf("LoadTasks failed on empty file: %v", err)
	}

	if len(loadedTasks) != 0 {
		t.Errorf("Expected empty tasks, got %v", loadedTasks)
	}
}

func TestNextID(t *testing.T) {
	tests := []struct {
		name  string
		tasks []task.Task
		want  int
	}{
		{"empty", []task.Task{}, 1},
		{"one_task", []task.Task{{ID: 1}}, 2},
		{"multiple_tasks", []task.Task{{ID: 1}, {ID: 3}, {ID: 2}}, 4},
		{"no_gaps", []task.Task{{ID: 1}, {ID: 2}, {ID: 3}}, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := task.NextID(tt.tasks)
			if got != tt.want {
				t.Errorf("NextID(%v) = %d; want %d", tt.tasks, got, tt.want)
			}
		})
	}
}

func TestAddTask(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	desc := "New Task"
	addedTask, err := task.AddTask(desc)
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}

	if addedTask.Description != desc || addedTask.ID != 1 || addedTask.Done != false {
		t.Errorf("Added task mismatch. Expected {ID:1, Description: %s, Done:false}, got %v", desc, addedTask)
	}

	loadedTasks, err := task.LoadTasks()
	if err != nil {
		t.Fatalf("LoadTasks after AddTask failed: %v", err)
	}

	if len(loadedTasks) != 1 || loadedTasks[0].Description != desc {
		t.Errorf("Tasks file content mismatch after AddTask. Got %v", loadedTasks)
	}
}

func TestListPendingTasks(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	tasksToSave := []task.Task{
		{ID: 1, Description: "Pending 1", Done: false},
		{ID: 2, Description: "Done 1", Done: true},
		{ID: 3, Description: "Pending 2", Done: false},
	}
	err := task.SaveTasks(tasksToSave)
	if err != nil {
		t.Fatalf("SaveTasks failed: %v", err)
	}

	pending, err := task.ListPendingTasks()
	if err != nil {
		t.Fatalf("ListPendingTasks failed: %v", err)
	}

	expected := []task.Task{
		{ID: 1, Description: "Pending 1", Done: false},
		{ID: 3, Description: "Pending 2", Done: false},
	}

	if !reflect.DeepEqual(pending, expected) {
		t.Errorf("Pending tasks mismatch. Expected %v, got %v", expected, pending)
	}
}

func TestCompleteTask(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	tasksToSave := []task.Task{
		{ID: 1, Description: "Pending Task", Done: false},
	}
	err := task.SaveTasks(tasksToSave)
	if err != nil {
		t.Fatalf("SaveTasks failed: %v", err)
	}

	completedTask, err := task.CompleteTask(1)
	if err != nil {
		t.Fatalf("CompleteTask failed: %v", err)
	}

	if completedTask.ID != 1 || completedTask.Done != true {
		t.Errorf("Completed task mismatch. Expected {ID:1, Done:true}, got %v", completedTask)
	}

	loadedTasks, err := task.LoadTasks()
	if err != nil {
		t.Fatalf("LoadTasks after CompleteTask failed: %v", err)
	}

	if len(loadedTasks) != 1 || !loadedTasks[0].Done {
		t.Errorf("Tasks file content mismatch after CompleteTask. Got %v", loadedTasks)
	}
}

func TestCompleteTaskNotFound(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	tasksToSave := []task.Task{
		{ID: 1, Description: "Pending Task", Done: false},
	}
	err := task.SaveTasks(tasksToSave)
	if err != nil {
		t.Fatalf("SaveTasks failed: %v", err)
	}

	_, err = task.CompleteTask(99)
	if err == nil {
		t.Fatalf("CompleteTask for non-existent ID should return an error, got nil")
	}
	if err.Error() != "task with ID 99 not found" {
		t.Errorf("Expected 'task with ID 99 not found' error, got %v", err)
	}
}