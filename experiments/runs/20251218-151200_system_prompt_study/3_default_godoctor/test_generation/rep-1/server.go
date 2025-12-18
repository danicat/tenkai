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

// CalculateTax calculates the tax based on a simplified progressive tax system.
func CalculateTax(income float64) float64 {
	if income <= 0 {
		return 0 // No tax for zero or negative income
	}

	var tax float64

	// First bracket: 0-10,000 at 10%
	if income <= 10000 {
		tax = income * 0.10
	} else if income <= 30000 {
		// Second bracket: 10,001-30,000 at 20%
		tax = 10000*0.10 + (income-10000)*0.20
	} else if income <= 60000 {
		// Third bracket: 30,001-60,000 at 30%
		tax = 10000*0.10 + 20000*0.20 + (income-30000)*0.30
	} else {
		// Fourth bracket: 60,001+ at 40%
		tax = 10000*0.10 + 20000*0.20 + 30000*0.30 + (income-60000)*0.40
	}

	return tax
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
