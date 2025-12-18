package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
)

func main() {
	// Initialize a new TodoList
	tl := NewTodoList()

	// Load existing tasks
	err := tl.LoadTasks()
	if err != nil {
		printErrorAndExit(err, 1)
	}

	if len(os.Args) < 2 {
		printErrorAndExit(fmt.Errorf("no command provided"), 42)
	}

	command := os.Args[1]

	switch command {
	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		description := addCmd.String("description", "", "Task description")
		addCmd.Parse(os.Args[2:])

		if *description == "" {
			// If the description flag is not used, try to get it from positional arguments
			if len(addCmd.Args()) > 0 {
				*description = addCmd.Args()[0]
			}
		}

		if *description == "" {
			printErrorAndExit(fmt.Errorf("missing task description"), 42)
		}

		task, err := tl.AddTask(*description)
		if err != nil {
			printErrorAndExit(err, 1)
		}

		err = tl.SaveTasks()
		if err != nil {
			printErrorAndExit(err, 1)
		}

		printJSON("id", task.ID, "status", "added")

	case "list":
		listCmd := flag.NewFlagSet("list", flag.ExitOnError)
		listCmd.Parse(os.Args[2:])

		tasks := tl.ListPendingTasks()
		printJSONTasks(tasks)

	case "done":
		doneCmd := flag.NewFlagSet("done", flag.ExitOnError)
		idStr := doneCmd.String("id", "0", "Task ID")
		doneCmd.Parse(os.Args[2:])

		id := 0
		if *idStr == "0" {
			// If the id flag is not used, try to get it from positional arguments
			if len(doneCmd.Args()) > 0 {
				id, err = strconv.Atoi(doneCmd.Args()[0])
				if err != nil {
					printErrorAndExit(fmt.Errorf("invalid task ID"), 42)
				}
			}
		}

		if id == 0 {
			printErrorAndExit(fmt.Errorf("missing task ID"), 42)
		}

		task, err := tl.CompleteTask(id)
		if err != nil {
			printErrorAndExit(err, 1)
		}

		err = tl.SaveTasks()
		if err != nil {
			printErrorAndExit(err, 1)
		}

		printJSON("id", task.ID, "status", "completed")

	default:
		printErrorAndExit(fmt.Errorf("unknown command: %s", command), 42)
	}
}

func printJSON(args ...interface{}) {
	data := make(map[string]interface{})
	for i := 0; i < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok {
			// This should not happen with correct usage
			continue
		}
		data[key] = args[i+1]
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.Encode(data)
}

func printJSONTasks(tasks []Task) {
	if len(tasks) == 0 {
		fmt.Println("[]")
		return
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.Encode(tasks)
}

func printErrorAndExit(err error, code int) {
	response := map[string]string{"error": err.Error()}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.Encode(response)
	os.Exit(code)
}
