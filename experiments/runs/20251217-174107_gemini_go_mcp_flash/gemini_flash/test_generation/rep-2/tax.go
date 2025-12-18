package main

import "fmt"

// TaxBracket represents a single tax bracket with its upper limit and tax rate.
type TaxBracket struct {
	Limit float64
	Rate  float64
}

// CalculateTax calculates the income tax based on a progressive tax system.
// It takes the total income as input and returns the calculated tax amount.
func CalculateTax(income float64) (float64, error) {
	if income < 0 {
		return 0, fmt.Errorf("income cannot be negative")
	}

	// Define tax brackets (example values, can be modified)
	// Bracket 1: 0% for income up to 10,000
	// Bracket 2: 10% for income from 10,001 to 50,000
	// Bracket 3: 20% for income from 50,001 to 100,000
	// Bracket 4: 30% for income above 100,000
	brackets := []TaxBracket{
		{Limit: 10000, Rate: 0.00},  // 0% for first 10,000
		{Limit: 50000, Rate: 0.10},  // 10% for next 40,000
		{Limit: 100000, Rate: 0.20}, // 20% for next 50,000
		{Limit: -1, Rate: 0.30},     // 30% for anything above 100,000 (-1 indicates no upper limit)
	}

	var totalTax float64
	var previousBracketLimit float64

	for _, bracket := range brackets {
		if income <= previousBracketLimit {
			break
		}

		taxableInThisBracket := 0.0
		currentBracketLimit := bracket.Limit

		if currentBracketLimit == -1 || income < currentBracketLimit {
			taxableInThisBracket = income - previousBracketLimit
		} else {
			taxableInThisBracket = currentBracketLimit - previousBracketLimit
		}

		if taxableInThisBracket > 0 {
			totalTax += taxableInThisBracket * bracket.Rate
		}
		previousBracketLimit = currentBracketLimit
	}

	return totalTax, nil
}
