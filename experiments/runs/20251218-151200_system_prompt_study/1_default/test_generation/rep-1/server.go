package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
)

// TaxBracket defines a tax bracket with a rate, minimum income, and maximum income.
type TaxBracket struct {
	MinIncome float64
	MaxIncome float64 // Upper limit of the bracket (exclusive). Use math.MaxFloat64 for the top bracket.
	Rate      float64
}

// TaxCalculator holds the tax brackets.
type TaxCalculator struct {
	Brackets []TaxBracket
}

// NewTaxCalculator creates a new TaxCalculator with predefined and sorted brackets.
func NewTaxCalculator() *TaxCalculator {
	return &TaxCalculator{
		Brackets: []TaxBracket{
			{MinIncome: 0, MaxIncome: 10000, Rate: 0.10},       // 10% for income from $0 to $10,000 (exclusive)
			{MinIncome: 10000, MaxIncome: 50000, Rate: 0.15},   // 15% for income from $10,000 to $50,000 (exclusive)
			{MinIncome: 50000, MaxIncome: 100000, Rate: 0.20},  // 20% for income from $50,000 to $100,000 (exclusive)
			{MinIncome: 100000, MaxIncome: math.MaxFloat64, Rate: 0.25}, // 25% for income above $100,000
		},
	}
}

// CalculateTax calculates the tax based on the given income.
// This function contains intentional bugs to be caught by tests.
func (tc *TaxCalculator) CalculateTax(income float64) float64 {
	if income <= 0 {
		return 0.0
	}

	totalTax := 0.0
	for _, bracket := range tc.Brackets {
		if income > bracket.MinIncome {
			taxableInBracket := math.Min(income, bracket.MaxIncome) - bracket.MinIncome
			if taxableInBracket > 0 {
				totalTax += taxableInBracket * bracket.Rate
			}
		} else {
			// Income is below or at the current bracket's minimum, so no more tax to calculate
			break
		}
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