package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"

	"rep-2/todo"
)

func main() {
	if len(os.Args) < 2 {
		printErrorAndExit("invalid arguments", 42)
	}

	switch os.Args[1] {
	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		addCmd.Parse(os.Args[2:])
		if addCmd.NArg() == 0 {
			printErrorAndExit("invalid arguments", 42)
		}
		description := addCmd.Arg(0)

		result, err := todo.Add(description)
		if err != nil {
			printErrorAndExit(err.Error(), 1)
		}
		printJSON(result)

	case "list":
		listCmd := flag.NewFlagSet("list", flag.ExitOnError)
		listCmd.Parse(os.Args[2:])

		result, err := todo.List()
		if err != nil {
			printErrorAndExit(err.Error(), 1)
		}
		printJSON(result)

	case "done":
		doneCmd := flag.NewFlagSet("done", flag.ExitOnError)
		doneCmd.Parse(os.Args[2:])
		if doneCmd.NArg() == 0 {
			printErrorAndExit("invalid arguments", 42)
		}
		taskIDStr := doneCmd.Arg(0)
		taskID, err := strconv.Atoi(taskIDStr)
		if err != nil {
			printErrorAndExit("invalid arguments", 42)
		}

		result, err := todo.Complete(taskID)
		if err != nil {
			printErrorAndExit(err.Error(), 1)
		}
		printJSON(result)

	default:
		printErrorAndExit("invalid arguments", 42)
	}
}

func printErrorAndExit(message string, code int) {
	output := map[string]string{"error": message}
	jsonOutput, _ := json.Marshal(output)
	fmt.Println(string(jsonOutput))
	os.Exit(code)
}

func printJSON(data interface{}) {
	jsonOutput, err := json.Marshal(data)
	if err != nil {
		printErrorAndExit(fmt.Sprintf("error marshaling JSON: %v", err), 1)
	}
	fmt.Println(string(jsonOutput))
}
