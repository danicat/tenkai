package main

import (
	"fmt"
)

// CalculateTax calculates the tax based on income and year.
// This is a simplified example and does not represent real-world tax laws.
func CalculateTax(income float64, year int) (float64, error) {
	if income < 0 {
		return 0, fmt.Errorf("income cannot be negative")
	}

	// Simplified tax brackets (corrected logic)
	// Bracket 1: 0 - 10,000, Tax: 0%
	// Bracket 2: 10,001 - 50,000, Tax: 10% on amount over 10,000
	// Bracket 3: 50,001+, Tax: 10% on 40,000 (from bracket 2) + 20% on amount over 50,000

	var tax float64

	if income <= 10000 {
		tax = 0
	} else if income <= 50000 {
		// Corrected: 10% on amount over 10,000
		tax = (income - 10000) * 0.10
	} else {
		// Corrected: 10% on 40,000 (which is 4000) + 20% on amount over 50,000
		tax = (50000 - 10000) * 0.10   // Tax from second bracket
		tax += (income - 50000) * 0.20 // Tax from third bracket
	}

	return tax, nil
}
