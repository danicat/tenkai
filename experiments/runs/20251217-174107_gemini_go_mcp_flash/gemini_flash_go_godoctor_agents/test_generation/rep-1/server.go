package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// TaxCalculator holds tax calculation logic.
type TaxCalculator struct{}

// CalculateTax calculates tax based on income and year.
// This is a simplified and intentionally buggy implementation for the exercise.
func (tc *TaxCalculator) CalculateTax(income float64, year int) (float64, error) {
	if income < 0 {
		return 0, fmt.Errorf("income cannot be negative")
	}

	var tax float64

	// First bracket: 0-10000 at 10%
	if income > 10000 {
		tax += 10000 * 0.10 // Tax on the full first bracket
		income -= 10000
	} else {
		tax += income * 0.10
		return tax, nil
	}

	// Second bracket: 10001-50000 at 20%
	if income > 40000 { // 50000 - 10000 = 40000
		tax += 40000 * 0.20 // Tax on the full second bracket
		income -= 40000
	} else {
		tax += income * 0.20
		return tax, nil
	}

	// Third bracket: 50001+ at 30%
	tax += income * 0.30

	return tax, nil
}

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
	_, err = w.Write(body)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/echo", echoHandler)
	return mux
}

func main() {
	mux := setupRoutes()
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
