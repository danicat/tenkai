package main

import (
	"fmt"
	"math"
)

// TaxBracket defines a tax bracket with a minimum and maximum threshold and a rate.
type TaxBracket struct {
	Min  float64
	Max  float64
	Rate float64
}

// TaxCalculator calculates tax based on predefined brackets.
type TaxCalculator struct {
	Brackets []TaxBracket
}

// NewTaxCalculator creates a new TaxCalculator with default brackets.
func NewTaxCalculator() *TaxCalculator {
	return &TaxCalculator{
		Brackets: []TaxBracket{
			{Min: 0, Max: 10000, Rate: 0.10},       // 10% for income from 0 up to 10,000
			{Min: 10000, Max: 50000, Rate: 0.15},   // 15% for income from 10,000 up to 50,000
			{Min: 50000, Max: math.MaxFloat64, Rate: 0.25}, // 25% for income from 50,000 onwards
		},
	}
}

// CalculateTax calculates the tax for a given income and year.
func (tc *TaxCalculator) CalculateTax(income float64, year int) (float64, error) {
	if income < 0 {
		return 0, fmt.Errorf("income cannot be negative")
	}

	totalTax := 0.0

	// Iterate through brackets, applying rates to the portion of income that falls into each band.
	for _, bracket := range tc.Brackets {
		// If the income is entirely below the current bracket's minimum, then this and higher brackets don't apply.
		if income <= bracket.Min {
			continue
		}

		// Find the effective upper bound for taxing within this bracket's band.
		// It's the smaller of the total income or the bracket's maximum.
		taxableUpperBound := math.Min(income, bracket.Max)

		// The amount of income to tax at *this* bracket's rate is the difference
		// between the effective upper bound and the bracket's minimum.
		// Ensure this amount is not negative.
		amountInThisBand := math.Max(0, taxableUpperBound - bracket.Min)

		totalTax += amountInThisBand * bracket.Rate
		
		// For leap year coverage (no actual tax change)
		if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
			// This block would contain logic if tax rules changed for leap years
		}
	}

	return totalTax, nil
}