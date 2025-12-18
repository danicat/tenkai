package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// TaxBracket defines a tax bracket with a rate and a maximum income for that bracket.
type TaxBracket struct {
	Rate float64
	Max  float64 // Maximum income for this bracket. If 0, it's the highest bracket.
}

// TaxCalculator holds the tax brackets.
type TaxCalculator struct {
	Brackets []TaxBracket
}

// CalculateTax calculates the tax for a given income based on the calculator's brackets.
func (tc *TaxCalculator) CalculateTax(income float64) float64 {
	totalTax := 0.0
	remainingIncome := income
	previousBracketMax := 0.0

	for _, bracket := range tc.Brackets {
		if remainingIncome <= 0 {
			break
		}

		currentBracketCeiling := bracket.Max
		if bracket.Max == 0 {
			// This is the highest bracket, tax all remaining income
			totalTax += remainingIncome * bracket.Rate
			remainingIncome = 0
			break
		}

		// The amount of income that falls within the current bracket's scope
		// relative to the previous bracket's max.
		bracketTaxableAmount := currentBracketCeiling - previousBracketMax

		if remainingIncome > bracketTaxableAmount {
			// If remaining income covers the entire current bracket
			totalTax += bracketTaxableAmount * bracket.Rate
			remainingIncome -= bracketTaxableAmount
		} else {
			// If remaining income falls within this bracket (or partially covers it)
			totalTax += remainingIncome * bracket.Rate
			remainingIncome = 0
			break
		}

		previousBracketMax = currentBracketCeiling
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
