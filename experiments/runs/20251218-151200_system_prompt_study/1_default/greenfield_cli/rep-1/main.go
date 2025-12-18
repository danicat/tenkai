package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
)

var tasksFile = "tasks.json"

// Task represents a single todo item.
type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

// tasks represents the entire list of tasks and a mutex for concurrent access.
type tasks struct {
	sync.Mutex
	List []Task
}

var currentTasks = &tasks{}

func loadTasks() error {
	fmt.Println("DEBUG: Entering loadTasks")
	currentTasks.Lock()
	defer currentTasks.Unlock()

	data, err := ioutil.ReadFile(tasksFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("DEBUG: tasksFile does not exist, initializing empty list")
			currentTasks.List = []Task{}
			return nil
		}
		fmt.Println("DEBUG: Error reading tasks file:", err)
		return fmt.Errorf("error reading tasks file: %w", err)
	}

	if len(data) == 0 {
		fmt.Println("DEBUG: tasksFile is empty, initializing empty list")
		currentTasks.List = []Task{}
		return nil
	}

	err = json.Unmarshal(data, &currentTasks.List)
	if err != nil {
		fmt.Println("DEBUG: Error unmarshaling tasks:", err)
		return fmt.Errorf("error unmarshaling tasks: %w", err)
	}
	fmt.Println("DEBUG: Exiting loadTasks, tasks loaded successfully")
	return nil
}

func saveTasks() error {
	fmt.Println("DEBUG: Entering saveTasks")
	currentTasks.Lock()
	defer currentTasks.Unlock()

	data, err := json.MarshalIndent(currentTasks.List, "", "  ")
	if err != nil {
		fmt.Println("DEBUG: Error marshaling tasks:", err)
		return fmt.Errorf("error marshaling tasks: %w", err)
	}

	err = ioutil.WriteFile(tasksFile, data, 0644)
	if err != nil {
		fmt.Println("DEBUG: Error writing tasks file:", err)
		return fmt.Errorf("error writing tasks file: %w", err)
	}
	fmt.Println("DEBUG: Exiting saveTasks, tasks saved successfully")
	return nil
}

func printError(err error, exitCode int) int {
	fmt.Println("DEBUG: Entering printError with error:", err, "and exit code:", exitCode)
	output, _ := json.Marshal(map[string]string{"error": err.Error()})
	fmt.Println(string(output))
	return exitCode
}

func addCommand(args []string) int {
	fmt.Println("DEBUG: Entering addCommand with args:", args)
	addFlagSet := flag.NewFlagSet("add", flag.ContinueOnError)
	
	// Use a custom usage function to avoid printing to stderr
	addFlagSet.Usage = func() { 
		addFlagSet.SetOutput(ioutil.Discard) // Redirect flag errors to /dev/null
		fmt.Println(`{"error": "invalid arguments: missing task description"}`)
		return
	}

	addFlagSet.Parse(args)

	if addFlagSet.NArg() == 0 {
		return printError(fmt.Errorf("invalid arguments: missing task description"), 42)
	}

	description := addFlagSet.Arg(0)
	fmt.Println("DEBUG: addCommand received description:", description)

	currentTasks.Lock()
	defer currentTasks.Unlock()

	nextID := 1
	if len(currentTasks.List) > 0 {
		nextID = currentTasks.List[len(currentTasks.List)-1].ID + 1
	}

	task := Task{
		ID:          nextID,
		Description: description,
		Done:        false,
	}

	currentTasks.List = append(currentTasks.List, task)
	fmt.Println("DEBUG: Task added to list:", task)

	if err := saveTasks(); err != nil {
		return printError(err, 1)
	}

	output, _ := json.Marshal(map[string]interface{}{"id": task.ID, "status": "added"})
	fmt.Println(string(output))
	fmt.Println("DEBUG: Exiting addCommand successfully")
	return 0
}

func listCommand() int {
	fmt.Println("DEBUG: Entering listCommand")
	currentTasks.Lock()
	defer currentTasks.Unlock()

	var pendingTasks []Task
	for _, task := range currentTasks.List {
		if !task.Done {
			pendingTasks = append(pendingTasks, task)
		}
	}

	output, err := json.Marshal(pendingTasks)
	if err != nil {
		return printError(fmt.Errorf("error marshaling tasks to JSON: %w", err), 1)
	}
	fmt.Println(string(output))
	fmt.Println("DEBUG: Exiting listCommand successfully")
	return 0
}

func doneCommand(args []string) int {
	fmt.Println("DEBUG: Entering doneCommand with args:", args)
	doneFlagSet := flag.NewFlagSet("done", flag.ContinueOnError)
	
	// Use a custom usage function to avoid printing to stderr
	doneFlagSet.Usage = func() {
		doneFlagSet.SetOutput(ioutil.Discard) // Redirect flag errors to /dev/null
		fmt.Println(`{"error": "invalid arguments: missing task ID"}`)
		return
	}

	doneFlagSet.Parse(args)

	if doneFlagSet.NArg() == 0 {
		return printError(fmt.Errorf("invalid arguments: missing task ID"), 42)
	}

	taskIDStr := doneFlagSet.Arg(0)
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		return printError(fmt.Errorf("invalid task ID: %s", taskIDStr), 42)
	}
	fmt.Println("DEBUG: doneCommand received taskID:", taskID)

	currentTasks.Lock()
	defer currentTasks.Unlock()

	found := false
	for i := range currentTasks.List {
		if currentTasks.List[i].ID == taskID {
			currentTasks.List[i].Done = true
			found = true
			break
		}
	}

	if !found {
		return printError(fmt.Errorf("task with ID %d not found", taskID), 42)
	}

	if err := saveTasks(); err != nil {
		return printError(err, 1)
	}

	output, _ := json.Marshal(map[string]interface{}{"id": taskID, "status": "completed"})
	fmt.Println(string(output))
	fmt.Println("DEBUG: Exiting doneCommand successfully")
	return 0
}

func run(args []string) int {
	fmt.Println("DEBUG: Entering run with args:", args)
	if err := loadTasks(); err != nil {
		return printError(err, 1)
	}

	if len(args) < 2 {
		return printError(fmt.Errorf("invalid arguments: please specify a command (add, list, done)"), 42)
	}

	switch args[1] {
	case "add":
		fmt.Println("DEBUG: Calling addCommand")
		return addCommand(args[2:])
	case "list":
		fmt.Println("DEBUG: Calling listCommand")
		return listCommand()
	case "done":
		fmt.Println("DEBUG: Calling doneCommand")
		return doneCommand(args[2:])
	default:
		return printError(fmt.Errorf("unknown command: %s", args[1]), 42)
	}
}

func main() {
	exitCode := run(os.Args)
	os.Exit(exitCode)
}