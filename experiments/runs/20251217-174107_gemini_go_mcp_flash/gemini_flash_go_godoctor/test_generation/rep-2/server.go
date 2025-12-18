package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// TaxCalculator calculates tax based on income brackets.
type TaxCalculator struct{}

// CalculateTax calculates the tax for a given income.
// This implementation contains subtle logic bugs for testing purposes.
func (tc *TaxCalculator) CalculateTax(income float64) float64 {
	if income <= 0 {
		return 0
	}

	var tax float64

	// Bug: Off-by-one or incorrect cumulative calculation
	if income <= 10000 {
		tax = income * 0.10
	} else if income > 10000 && income <= 50000 { // Bug: Should be > 10000 for next bracket, not >=
		tax = (10000 * 0.10) + ((income - 10000) * 0.20)
	} else { // income > 50000
		tax = (10000 * 0.10) + (40000 * 0.20) + ((income - 50000) * 0.30)
	}

	return tax
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
	w.Write(body)
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

