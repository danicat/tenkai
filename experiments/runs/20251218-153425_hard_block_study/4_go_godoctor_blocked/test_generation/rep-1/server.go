package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// TaxCalculator provides methods for calculating tax.
type TaxCalculator struct{}

// CalculateTax calculates the tax based on income and year.
func (tc *TaxCalculator) CalculateTax(income float64, year int) (float64, error) {
	var tax float64

	if income < 0 {
		return 0, fmt.Errorf("income cannot be negative")
	}

	// First bracket: 10% on income up to 10,000
	if income > 0 {
		if income <= 10000 {
			tax += income * 0.10
		} else {
			tax += 10000 * 0.10 // Max from first bracket
			remainingIncome := income - 10000

			// Second bracket: 20% on income from 10,001 to 50,000
			if remainingIncome > 0 {
				if remainingIncome <= 40000 { // 50000 - 10000
					tax += remainingIncome * 0.20
				} else {
					tax += 40000 * 0.20 // Max from second bracket
					remainingIncome -= 40000

					// Third bracket: 30% on income above 50,000
					if remainingIncome > 0 {
						tax += remainingIncome * 0.30
					}
				}
			}
		}
	}

	_ = year // The 'year' parameter is currently unused

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
