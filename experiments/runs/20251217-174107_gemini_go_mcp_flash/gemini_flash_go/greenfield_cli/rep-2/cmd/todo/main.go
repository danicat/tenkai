package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/user/todo-cli/internal/task"
)

const (
	errorInvalidArgs    = "invalid arguments"
	exitCodeInvalidArgs = 42
)

type output struct {
	ID     int         `json:"id,omitempty"`
	Status string      `json:"status,omitempty"`
	Error  string      `json:"error,omitempty"`
	Tasks  []task.Task `json:"tasks,omitempty"`
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout))
}

func run(args []string, w io.Writer) int {
	store, err := task.NewStore()
	if err != nil {
		printError(w, err.Error())
		return 1
	}

	if len(args) == 0 {
		printError(w, errorInvalidArgs)
		return exitCodeInvalidArgs
	}

	cmd := args[0]
	subArgs := args[1:]

	switch cmd {
	case "add":
		return handleAdd(subArgs, store, w)
	case "list":
		return handleList(store, w)
	case "done":
		return handleDone(subArgs, store, w)
	default:
		printError(w, fmt.Sprintf("unknown command: %s", cmd))
		return exitCodeInvalidArgs
	}
}

func handleAdd(args []string, store *task.Store, w io.Writer) int {
	if len(args) == 0 {
		printError(w, errorInvalidArgs)
		return exitCodeInvalidArgs
	}
	// Join all arguments to form the description, in case it contains spaces
	description := args[0]

	t, err := store.AddTask(description)
	if err != nil {
		printError(w, err.Error())
		return 1
	}
	printJSON(w, output{ID: t.ID, Status: "added"})
	return 0
}

func handleList(store *task.Store, w io.Writer) int {
	tasks := store.ListTasks(false) // Only list pending tasks
	printJSON(w, output{Tasks: tasks})
	return 0
}

func handleDone(args []string, store *task.Store, w io.Writer) int {
	if len(args) == 0 {
		printError(w, errorInvalidArgs)
		return exitCodeInvalidArgs
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		printError(w, errorInvalidArgs)
		return exitCodeInvalidArgs
	}

	t, err := store.CompleteTask(id)
	if err != nil {
		printError(w, err.Error())
		return 1
	}
	printJSON(w, output{ID: t.ID, Status: "completed"})
	return 0
}

func printJSON(w io.Writer, data interface{}) {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "") // No indent for compact output
	if err := encoder.Encode(data); err != nil {
		// Fallback to plain error if JSON encoding fails
		fmt.Fprintf(w, "{\"error\": \"failed to encode JSON output: %s\"}\n", err)
	}
}

func printError(w io.Writer, errMsg string) {
	printJSON(w, output{Error: errMsg})
}
