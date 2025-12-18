package main

import (
	"testing"
)

func TestCalculateTax(t *testing.T) {
	tc := NewTaxCalculator()

	tests := []struct {
		name          string
		income        float64
		year          int
		expectedTax   float64
		expectedError bool
	}{
		{
			name:          "Negative income",
			income:        -100,
			year:          2023,
			expectedTax:   0,
			expectedError: true,
		},
		{
			name:          "Zero income",
			income:        0,
			year:          2023,
			expectedTax:   0,
			expectedError: false,
		},
		{
			name:          "Income within first bracket",
			income:        5000,
			year:          2023,
			expectedTax:   500, // 10% of 5000
			expectedError: false,
		},
		{
			name:          "Income at first bracket boundary",
			income:        10000,
			year:          2023,
			expectedTax:   1000, // 10% of 10000
			expectedError: false,
		},
		{
			name:          "Income slightly above first bracket (into second)",
			income:        15000,
			year:          2023,
			expectedTax:   1750, // 10% of 10000 + 15% of 5000
			expectedError: false,
		},
		{
			name:          "Income at second bracket boundary",
			income:        50000,
			year:          2023,
			expectedTax:   7000, // 10% of 10000 + 15% of 40000
			expectedError: false,
		},
		{
			name:          "Income above second bracket (into third)",
			income:        60000,
			year:          2023,
			expectedTax:   9500, // 10% of 10000 + 15% of 40000 + 25% of 10000
			expectedError: false,
		},
		{
			name:          "Income above second bracket (into third) - Leap Year",
			income:        60000,
			year:          2024, // A leap year
			expectedTax:   9500, // Same expected tax as year doesn't affect calculation
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tax, err := tc.CalculateTax(tt.income, tt.year)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected an error for income %f, but got none", tt.income)
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error for income %f, but got: %v", tt.income, err)
				}
				if tax != tt.expectedTax {
					t.Errorf("For income %f, expected tax %f, but got %f", tt.income, tt.expectedTax, tax)
				}
			}
		})
	}
}