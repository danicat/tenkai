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

// TaxCalculator holds tax-related parameters.
type TaxCalculator struct{}

// CalculateTax calculates tax based on income and year.
// It includes simplified tax brackets and a hypothetical leap year adjustment.
// This initial implementation may contain subtle bugs for testing purposes.
func (tc *TaxCalculator) CalculateTax(income float64, year int) float64 {
	var tax float64

	// Simplified tax brackets (with intentional off-by-one or boundary bugs)
	if income <= 10000 {
		tax = income * 0.10
	} else if income <= 50000 {
		tax = 10000*0.10 + (income-10000)*0.20 // Bug: income-10000 should be income-10001 or correct bracket calculation
	} else {
		tax = 10000*0.10 + 40000*0.20 + (income-50000)*0.30 // Bug: boundary
	}

	// Hypothetical leap year effect (now using the correct helper)
	if isLeapYear(year) {
		tax += 100.0 // Flat bonus for a true leap year
	}

	return tax
}

// isLeapYear determines if a year is a leap year.
func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

func main() {
	mux := setupRoutes()
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
