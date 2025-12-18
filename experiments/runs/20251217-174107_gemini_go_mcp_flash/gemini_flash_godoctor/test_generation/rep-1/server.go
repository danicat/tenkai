package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// TaxCalculator provides methods for calculating tax.
type TaxCalculator struct{}

// isLeap determines if a year is a leap year.
func isLeap(year int) bool {
	return (year%4 == 0 && year%100 != 0) || year%400 == 0
}

// CalculateTax calculates the tax for a given income and year.
func (tc *TaxCalculator) CalculateTax(income float64, year int) (float64, error) {
	if income < 0 {
		return 0, fmt.Errorf("income cannot be negative")
	}

	var tax float64
	switch {
	case income <= 10000:
		tax = income * 0.10 // 10% for income up to 10,000
	case income <= 40000: // Corrected: 40000 for the next bracket
		tax = 10000*0.10 + (income-10000)*0.15 // 15% for income between 10,001 and 40,000
	default:
		tax = 10000*0.10 + 30000*0.15 + (income-40000)*0.25 // Corrected: 0.25 for the third bracket rate, and base adjusted to 40000
	}

	// Special deduction for leap years if income is over 30000
	if isLeap(year) && income > 30000 {
		tax -= 500
	}

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