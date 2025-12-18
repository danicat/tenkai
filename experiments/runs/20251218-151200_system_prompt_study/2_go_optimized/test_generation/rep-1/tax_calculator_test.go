package main

import (
	"math"
	"testing"
)

func TestGetCurrentYear(t *testing.T) {
	currentYear := GetCurrentYear()
	// Basic check: currentYear should be a reasonable year, e.g., not zero or negative.
	// More robust checks would involve mocking time, but for coverage, this is sufficient.
	if currentYear < 2023 {
		t.Errorf("GetCurrentYear() returned an unexpected year: %d", currentYear)
	}
}

func TestCalculateTax(t *testing.T) {
	tests := []struct {
		name          string
		income        float64
		year          int
		expectedTax   float64
		expectedError string
	}{
		{
			name:        "Income within first bracket - Non-leap year",
			income:      5000,
			year:        2023,
			expectedTax: 500.00,
			expectedError: "",
		},
		{
			name:        "Income at first bracket boundary - Non-leap year",
			income:      10000,
			year:        2023,
			expectedTax: 1000.00,
			expectedError: "",
		},
		{
			name:        "Income slightly above first bracket - Non-leap year",
			income:      10001,
			year:        2023,
			expectedTax: 1000.15,
			expectedError: "",
		},
		{
			name:        "Income within second bracket - Non-leap year",
			income:      20000,
			year:        2023,
			expectedTax: 2500.00,
			expectedError: "",
		},
		{
			name:        "Income at second bracket boundary - Non-leap year",
			income:      30000,
			year:        2023,
			expectedTax: 4000.00,
			expectedError: "",
		},
		{
			name:        "Income slightly above second bracket - Non-leap year",
			income:      30001,
			year:        2023,
			expectedTax: 4000.20,
			expectedError: "",
		},
		{
			name:        "Income within third bracket - Non-leap year",
			income:      40000,
			year:        2023,
			expectedTax: 6000.00,
			expectedError: "",
		},
		{
			name:        "Zero income",
			income:      0,
			year:        2023,
			expectedTax: 0.00,
			expectedError: "",
		},
		{
			name:        "Negative income",
			income:      -100,
			year:        2023,
			expectedTax: 0.00,
			expectedError: "income cannot be negative",
		},
		// Test cases for leap year bug
		{
			name:        "Income within first bracket - Leap year (2020)",
			income:      5000,
			year:        2020, // Leap year
			expectedTax: 500.00, // Should be same as non-leap year
			expectedError: "",
		},
		{
			name:        "Income at first bracket boundary - Leap year (2024)",
			income:      10000,
			year:        2024, // Leap year
			expectedTax: 1000.00, // Should be same as non-leap year
			expectedError: "",
		},
		{
			name:        "Income within second bracket - Leap year (2000)",
			income:      20000,
			year:        2000, // Leap year (divisible by 400)
			expectedTax: 2500.00,
			expectedError: "",
		},
		{
			name:        "Income within third bracket - Leap year (1900, not leap)",
			income:      40000,
			year:        1900, // Not a leap year (divisible by 100 but not 400)
			expectedTax: 6000.00,
			expectedError: "",
		},
	}

	tc := &TaxCalculator{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualTax, err := tc.CalculateTax(tt.income, tt.year)

			if tt.expectedError != "" {
				if err == nil || err.Error() != tt.expectedError {
					t.Errorf("For income %f, year %d: expected error %q, got %v", tt.income, tt.year, tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("For income %f, year %d: unexpected error: %v", tt.income, tt.year, err)
				}
				// Use a small epsilon for float comparison
				if math.Abs(actualTax-tt.expectedTax) > 0.001 {
					t.Errorf("For income %f, year %d: expected tax %f, got %f", tt.income, tt.year, tt.expectedTax, actualTax)
				}
			}
		})
	}
}
