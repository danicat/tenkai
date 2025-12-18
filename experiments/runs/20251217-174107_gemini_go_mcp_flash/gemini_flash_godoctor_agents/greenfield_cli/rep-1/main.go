package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		printErrorAndExit("no command provided", 42)
	}

	switch os.Args[1] {
	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		description := addCmd.String("description", "", "Description of the task")
		addCmd.Parse(os.Args[2:])

		if *description == "" {
			printErrorAndExit("invalid arguments", 42)
		}

		tasks, err := LoadTasks()
		if err != nil {
			printErrorAndExit(fmt.Sprintf("could not load tasks: %v", err), 1)
		}

		nextID := 1
		if len(tasks) > 0 {
			nextID = tasks[len(tasks)-1].ID + 1
		}

		newTask := Task{
			ID:          nextID,
			Description: *description,
			Done:        false,
		}

		tasks = append(tasks, newTask)

		err = SaveTasks(tasks)
		if err != nil {
			printErrorAndExit(fmt.Sprintf("could not save tasks: %v", err), 1)
		}

		response := map[string]interface{}{"id": newTask.ID, "status": "added"}
		jsonResponse, _ := json.Marshal(response)
		fmt.Println(string(jsonResponse))

	case "list":
		listCmd := flag.NewFlagSet("list", flag.ExitOnError)
		listCmd.Parse(os.Args[2:])

		tasks, err := LoadTasks()
		if err != nil {
			printErrorAndExit(fmt.Sprintf("could not load tasks: %v", err), 1)
		}

		var pendingTasks = make([]Task, 0)
		for _, task := range tasks {
			if !task.Done {
				pendingTasks = append(pendingTasks, task)
			}
		}

		jsonResponse, _ := json.Marshal(pendingTasks)
		fmt.Println(string(jsonResponse))

	case "done":
		doneCmd := flag.NewFlagSet("done", flag.ExitOnError)
		idStr := ""
		if len(os.Args) > 2 {
			idStr = os.Args[2]
		}
		doneCmd.Parse(os.Args[2:])

		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			printErrorAndExit("invalid arguments", 42)
		}

		tasks, err := LoadTasks()
		if err != nil {
			printErrorAndExit(fmt.Sprintf("could not load tasks: %v", err), 1)
		}

		found := false
		for i := range tasks {
			if tasks[i].ID == id {
				tasks[i].Done = true
				found = true
				break
			}
		}

		if !found {
			printErrorAndExit(fmt.Sprintf("task with ID %d not found", id), 42)
		}

		err = SaveTasks(tasks)
		if err != nil {
			printErrorAndExit(fmt.Sprintf("could not save tasks: %v", err), 1)
		}

		response := map[string]interface{}{"id": id, "status": "completed"}
		jsonResponse, _ := json.Marshal(response)
		fmt.Println(string(jsonResponse))

	default:
		printErrorAndExit("unknown command", 42)
	}
}

func printErrorAndExit(message string, code int) {
	output, _ := json.Marshal(map[string]string{"error": message})
	fmt.Println(string(output))
	os.Exit(code)
}
