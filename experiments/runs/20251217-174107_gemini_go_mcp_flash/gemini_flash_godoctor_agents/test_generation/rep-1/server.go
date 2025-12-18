package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello, %s!", name)
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	w.Write(body)
}

func setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/echo", echoHandler)
	return mux
}

// CalculateTax calculates the tax based on income using a simplified progressive tax system.
func CalculateTax(income float64, year int) float64 {
	tax := 0.0

	// Bracket 1: 0 - 20000, 10%
	if income <= 20000 {
		tax = income * 0.10
	} else if income > 20000 && income <= 50000 {
		// Bracket 2: 20001 - 50000, 20%
		tax = 20000 * 0.10 // Tax for first bracket
		tax += (income - 20000) * 0.20
	} else {
		// Bracket 3: 50001+, 30%
		tax = 20000 * 0.10 // Tax for first bracket
		tax += 30000 * 0.20 // Tax for second bracket (50000 - 20000)
		tax += (income - 50000) * 0.30
	}
	// The original bug (unconditional +100 for leap year) is removed here.

	return tax
}

func main() {
	mux := setupRoutes()
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
