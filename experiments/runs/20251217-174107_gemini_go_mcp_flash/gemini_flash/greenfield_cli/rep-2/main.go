package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"experiment/todo"
)

func main() {
	if len(os.Args) < 2 {
		printError("invalid arguments", 42)
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 3 {
			printError("invalid arguments", 42)
			return
		}
		description := strings.Join(os.Args[2:], " ")
		task, err := todo.AddTask(description)
		if err != nil {
			printError(err.Error(), 1)
			return
		}
		jsonOutput, _ := json.Marshal(map[string]interface{}{
			"id": task.ID, "status": "added",
		})
		fmt.Println(string(jsonOutput))

	case "list":
		tasks, err := todo.ListTasks()
		if err != nil {
			printError(err.Error(), 1)
			return
		}
		var pendingTasks []todo.Task
		for _, t := range tasks {
			if !t.Done {
				pendingTasks = append(pendingTasks, t)
			}
		}
		jsonOutput, _ := json.Marshal(pendingTasks)
		fmt.Println(string(jsonOutput))

	case "done":
		if len(os.Args) < 3 {
			printError("invalid arguments", 42)
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			printError("invalid arguments", 42)
			return
		}
		task, err := todo.CompleteTask(id)
		if err != nil {
			printError(err.Error(), 1)
			return
		}
		jsonOutput, _ := json.Marshal(map[string]interface{}{
			"id": task.ID, "status": "completed",
		})
		fmt.Println(string(jsonOutput))

	default:
		printError("invalid arguments", 42)
		return
	}
}

func printError(message string, exitCode int) {
	jsonOutput, _ := json.Marshal(map[string]string{"error": message})
	fmt.Fprintln(os.Stderr, string(jsonOutput))
	os.Exit(exitCode)
}
