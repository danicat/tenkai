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

// TaxCalculator provides methods for calculating tax.
type TaxCalculator struct{}

// CalculateTax calculates the tax based on a progressive tax system.
func (tc *TaxCalculator) CalculateTax(income float64) float64 {
	var tax float64

	switch {
	case income <= 10000:
		tax = 0
	case income <= 40000:
		tax = 0.10 * (income - 10000)
	case income <= 80000:
		tax = (0.10 * (40000 - 10000)) + (0.20 * (income - 40000))
	default:
		tax = (0.10 * (40000 - 10000)) + (0.20 * (80000 - 40000)) + (0.30 * (income - 80000))
	}

	return tax
}

func main() {
	mux := setupRoutes()
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
