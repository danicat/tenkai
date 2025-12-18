package main

import (
	"errors"
	"math"
	"time"
)

// TaxCalculator provides methods for calculating tax.
type TaxCalculator struct{}

// CalculateTax calculates the tax based on income and the year.
func (tc *TaxCalculator) CalculateTax(income float64, year int) (float64, error) {
	if income < 0 {
		return 0, errors.New("income cannot be negative")
	}

	tax := 0.0
	remainingIncome := income

	// First bracket: 0 - 10,000 at 10%
	if remainingIncome > 0 {
		taxableInFirstBracket := math.Min(remainingIncome, 10000.0)
		tax += taxableInFirstBracket * 0.10
		remainingIncome -= taxableInFirstBracket
	}

	// Second bracket: 10,001 - 30,000 (i.e., income portion from 10,000.01 to 30,000) at 15%
	// This covers the next 20,000 of income
	if remainingIncome > 0 {
		taxableInSecondBracket := math.Min(remainingIncome, 20000.0) // 30000 - 10000
		tax += taxableInSecondBracket * 0.15
		remainingIncome -= taxableInSecondBracket
	}

	// Third bracket: 30,001+ at 20%
	if remainingIncome > 0 {
		tax += remainingIncome * 0.20
	}

	return tax, nil
}

// GetCurrentYear returns the current year.
// This function is for demonstration and might not be used directly in tax calculation
// but could be part of a larger system using the calculator.
func GetCurrentYear() int {
	return time.Now().Year()
}