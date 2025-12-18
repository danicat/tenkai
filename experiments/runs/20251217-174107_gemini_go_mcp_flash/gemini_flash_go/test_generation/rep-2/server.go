package main

type TaxCalculator struct{}

// CalculateTax calculates the tax for a given income and year.
// It contains subtle logic bugs for the purpose of this exercise.
func (tc *TaxCalculator) CalculateTax(income float64, year int) float64 {
	tax := 0.0

	if income <= 0 {
		return 0
	}

	// Progressive tax calculation
	if income <= 10000 {
		tax = income * 0.10
	} else if income <= 50000 {
		tax = (10000 * 0.10) + ((income - 10000) * 0.20)
	} else {
		tax = (10000 * 0.10) + (40000 * 0.20) + ((income - 50000) * 0.30)
	}

	// Leap year deduction
	isLeapYear := (year%4 == 0 && year%100 != 0) || (year%400 == 0)
	if isLeapYear {
		tax -= 50.0
	}

	// Ensure tax is not negative
	if tax < 0 {
		return 0
	}

	return tax
}

