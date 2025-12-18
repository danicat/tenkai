package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// TaxCalculator and related functions
func isLeapYear(year int) bool {
	// Corrected leap year calculation
	return (year%4 == 0 && year%100 != 0) || year%400 == 0
}

func CalculateTax(income float64, year int) (float64, error) {
	if income < 0 {
		return 0, fmt.Errorf("income cannot be negative")
	}

	var tax float64

	switch {
	case income <= 20000:
		tax = income * 0.10
	case income <= 50000:
		tax = 20000*0.10 + (income-20000)*0.15
	default:
		tax = 20000*0.10 + 30000*0.15 + (income-50000)*0.20
	}

	// Corrected leap year bonus condition
	if isLeapYear(year) && income > 20000 { // Bonus applies for income in higher brackets
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
