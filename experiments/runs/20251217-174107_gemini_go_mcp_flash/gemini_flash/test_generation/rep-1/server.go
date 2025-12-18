package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

// TaxCalculator holds the tax bracket information.
type TaxCalculator struct {
	Brackets []TaxBracket
}

// TaxBracket defines a single tax bracket with an upper limit and rate.
type TaxBracket struct {
	Limit float64 // Upper limit of the bracket. -1 indicates no upper limit for the highest bracket.
	Rate  float64
}

// NewTaxCalculator creates a new TaxCalculator with predefined brackets.
func NewTaxCalculator() *TaxCalculator {
	return &TaxCalculator{
		Brackets: []TaxBracket{
			{Limit: 10000, Rate: 0.10},   // 10% for income up to 10,000
			{Limit: 50000, Rate: 0.20},   // 20% for income between 10,001 and 50,000
			{Limit: 100000, Rate: 0.30},  // 30% for income between 50,001 and 100,000
			{Limit: -1, Rate: 0.40},      // 40% for income above 100,000
		},
	}
}

// CalculateTax calculates the tax for a given income using progressive tax brackets.
func (tc *TaxCalculator) CalculateTax(income float64) float64 {
	if income <= 0 {
		return 0.0
	}

	var totalTax float64
	taxedInPreviousBrackets := 0.0

	for _, bracket := range tc.Brackets {
		if income <= taxedInPreviousBrackets {
			break // All income has been accounted for by lower brackets
		}

		// Determine the effective upper limit for the current bracket
		currentBracketUpperLimit := bracket.Limit
		if bracket.Limit == -1 { // For the highest bracket, effectively no upper limit
			currentBracketUpperLimit = income // Limit by the remaining income itself
		}

		// Calculate the amount of income that falls *within the range* of this specific bracket
		// and is also part of the total income.
		taxableInThisBracket := currentBracketUpperLimit - taxedInPreviousBrackets

		// Ensure we don't over-tax if the income is less than the bracket's full range
		if taxableInThisBracket > (income - taxedInPreviousBrackets) {
			taxableInThisBracket = income - taxedInPreviousBrackets
		}

		if taxableInThisBracket < 0 {
			taxableInThisBracket = 0
		}

		totalTax += taxableInThisBracket * bracket.Rate
		taxedInPreviousBrackets += taxableInThisBracket
	}

	return totalTax
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

func taxHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	incomeStr := r.URL.Query().Get("income")
	if incomeStr == "" {
		http.Error(w, "Missing income parameter", http.StatusBadRequest)
		return
	}

	income, err := strconv.ParseFloat(incomeStr, 64)
	if err != nil {
		http.Error(w, "Invalid income parameter", http.StatusBadRequest)
		return
	}

	tc := NewTaxCalculator()
	tax := tc.CalculateTax(income)

	fmt.Fprintf(w, "Income: %.2f, Calculated Tax: %.2f", income, tax)
}

func setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/echo", echoHandler)
	mux.HandleFunc("/tax", taxHandler) // Add tax handler
	return mux
}

func main() {
	mux := setupRoutes()
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
