package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// TaxCalculator holds methods for calculating tax.
type TaxCalculator struct{}

// CalculateTax calculates the tax for a given income based on a simplified bracket system.
func (tc *TaxCalculator) CalculateTax(income float64) float64 {
	var tax float64

	switch {
	case income <= 0:
		tax = 0
	case income <= 10000:
		tax = income * 0.10 // 10% for income up to 10,000
	case income <= 50000:
		// 10% on first 10,000 (1000)
		// FIXING THE BUG: 20% on remaining (income - 10000)
		tax = 1000 + (income-10000)*0.20
	default: // income > 50000
		// 10% on first 10,000 (1000)
		// 20% on next 40,000 (8000)
		// 30% on remaining (income - 50000)
		tax = 1000 + 8000 + (income-50000)*0.30
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
