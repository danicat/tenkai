package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// TaxCalculator provides methods to calculate tax based on income.
type TaxCalculator struct{}

// CalculateTax calculates the income tax based on a progressive tax system.
func (tc *TaxCalculator) CalculateTax(income float64) float64 {
	if income < 0 {
		return 0.00 // No tax on negative income
	}

	var tax float64

	// First bracket: 0 - 10,000 @ 10%
	if income <= 10000 {
		tax = income * 0.10
		return tax
	}
	tax += 10000 * 0.10 // Tax for the first 10,000

	// Second bracket: 10,001 - 30,000 @ 20%
	if income <= 30000 {
		tax += (income - 10000) * 0.20
		return tax
	}
	tax += (30000 - 10000) * 0.20 // Tax for the next 20,000

	// Third bracket: 30,001 - 60,000 @ 30%
	if income <= 60000 {
		tax += (income - 30000) * 0.30
		return tax
	}
	tax += (60000 - 30000) * 0.30 // Tax for the next 30,000

	// Fourth bracket: Above 60,000 @ 40%
	tax += (income - 60000) * 0.40
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
