package main

import (
	"testing"
)

func TestCalculateTax(t *testing.T) {
	tc := &TaxCalculator{}

	tests := []struct {
		name        string
		income      float64
		year        int
		expectedTax float64
		expectError bool
	}{
		{
			name:        "Income in first bracket",
			income:      5000,
			year:        2023,
			expectedTax: 500, // 5000 * 0.10
			expectError: false,
		},
		{
			name:        "Income at first bracket threshold",
			income:      10000,
			year:        2023,
			expectedTax: 1000, // 10000 * 0.10
			expectError: false,
		},
		{
			name:        "Income in second bracket (just above first)",
			income:      15000,
			year:        2023,
			expectedTax: 1750, // 10000*0.10 + 5000*0.15
			expectError: false,
		},
		{
			name:        "Income at second bracket threshold (intended correct limit of 40000)",
			income:      40000,
			year:        2023,
			expectedTax: 5500, // 10000*0.10 + 30000*0.15
			expectError: false,
		},
		{
			name:        "Income in third bracket (above 40000)",
			income:      50000,
			year:        2023,
			expectedTax: 8000, // 10000*0.10 + 30000*0.15 + 10000*0.25
			expectError: false,
		},
		{
			name:        "Income in third bracket (above 40000, higher value)",
			income:      60000,
			year:        2023,
			expectedTax: 10500, // 10000*0.10 + 30000*0.15 + 20000*0.25
			expectError: false,
		},
		{
			name:        "Negative income should error",
			income:      -100,
			year:        2023,
			expectedTax: 0,
			expectError: true,
		},
		{
			name:        "Leap year 2024 income 30000 (no deduction since income not > 30000)",
			income:      30000,
			year:        2024,
			expectedTax: 4000, // 10000*0.10 + 20000*0.15
			expectError: false,
		},
		{
			name:        "Leap year 2024 income 45000 (with deduction)",
			income:      45000,
			year:        2024,
			expectedTax: 6250, // (10000*0.10 + 30000*0.15 + 5000*0.25) - 500
			expectError: false,
		},
		{
			name:        "Non-leap year 2023 income 45000 (no deduction)",
			income:      45000,
			year:        2023,
			expectedTax: 6750, // 10000*0.10 + 30000*0.15 + 5000*0.25
			expectError: false,
		},
		{
			name:        "Year 1900, not a leap year (divisible by 100 but not 400)",
			income:      45000,
			year:        1900,
			expectedTax: 6750, // 10000*0.10 + 30000*0.15 + 5000*0.25 (no deduction)
			expectError: false,
		},
		{
			name:        "Year 2000, a leap year (divisible by 400)",
			income:      45000,
			year:        2000,
			expectedTax: 6250, // (10000*0.10 + 30000*0.15 + 5000*0.25) - 500
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualTax, err := tc.CalculateTax(tt.income, tt.year)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected an error for income %f, but got none", tt.income)
				}
				return
			}

			if err != nil {
				t.Errorf("Did not expect an error for income %f, but got: %v", tt.income, err)
			}

			if actualTax != tt.expectedTax {
				t.Errorf("For income %f, year %d: expected tax %f, got %f", tt.income, tt.year, tt.expectedTax, actualTax)
			}
		})
	}
}