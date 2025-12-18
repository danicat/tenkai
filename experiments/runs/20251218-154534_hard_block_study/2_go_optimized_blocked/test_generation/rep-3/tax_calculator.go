package main

// CalculateTax calculates the tax based on income and year.
func CalculateTax(income float64, year int) float64 {
	tax := 0.0

	// Corrected bracket calculation
	if income <= 10000 {
		tax = income * 0.10
	} else if income <= 50000 {
		tax = 10000*0.10 + (income-10000)*0.20
	} else {
		tax = 10000*0.10 + 40000*0.20 + (income-50000)*0.30
	}

	// Corrected leap year handling: now includes income >= 30000
	if isLeapYear(year) && income >= 30000 { // Applies to any leap year if income >= 30000
		tax += 100.0 // Additional flat tax for leap year
	}

	return tax
}

// isLeapYear determines if a year is a leap year.
func isLeapYear(year int) bool {
	// A year is a leap year if it is divisible by 4,
	// except for end-of-century years which must be divisible by 400.
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}
