package main

import (
	"fmt"
	"io"
	"log"
	"math"
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


// TaxCalculator holds tax bracket information
type TaxCalculator struct {
	Brackets []struct {
		Min    float64
		Rate   float64
		IsMax  bool // To indicate the last bracket might not have a max limit
	}
}

// isLeapYear checks if a year is a leap year.
func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// CalculateTax calculates tax based on income and year.
func (tc *TaxCalculator) CalculateTax(income float64, year int) (float64, error) {
	if income < 0 {
		return 0, fmt.Errorf("income cannot be negative")
	}

	totalTax := 0.0
	remainingIncome := income

	for i, bracket := range tc.Brackets {
		if remainingIncome <= 0 {
			break
		}

		currentBracketMin := bracket.Min
		var currentBracketMax float64

		if bracket.IsMax || i == len(tc.Brackets)-1 {
			currentBracketMax = income + 1 // Effectively no upper limit for remaining income
		} else {
			currentBracketMax = tc.Brackets[i+1].Min
		}

		// Amount of income that falls within the current bracket
		taxableInCurrentBracket := math.Min(remainingIncome, currentBracketMax-currentBracketMin)
		
		totalTax += taxableInCurrentBracket * bracket.Rate
		remainingIncome -= taxableInCurrentBracket
	}

	// Apply leap year bonus, ensuring tax doesn't go below zero
	if isLeapYear(year) {
		totalTax = math.Max(0, totalTax-100.0)
	}

	return totalTax, nil
}

func main() {
	// Initialize a dummy TaxCalculator for demonstration
	// This will not be used by the HTTP server, but available for testing
	_ = &TaxCalculator{
		Brackets: []struct {
			Min    float64
			Rate   float64
			IsMax  bool
		}{
			{Min: 0, Rate: 0.10},       // 0-10,000 at 10%
			{Min: 10000, Rate: 0.20},   // 10,000-30,000 at 20%
			{Min: 30000, Rate: 0.30},   // 30,000+ at 30%
		},
	}

	mux := setupRoutes()
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
