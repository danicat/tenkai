package api

import "fmt"

// Data represents a simple structure.
type Data struct {
	ID   int
	Name string
}

// helper is an unexported function.
func helper() {
	fmt.Println("helping")
}

// GetInfo is an exported function that returns a greeting.
func GetInfo() string {
	return "Hello from API"
}

// ProcessData is an exported function.
func ProcessData(d Data) error {
	fmt.Printf("Processing %s\n", d.Name)
	return nil
}

// InternalConfig is an unexported type.
type internalConfig struct {
	secret string
}
