package main

import (
	"math"
	"testing"
)

func TestTaxCalculator_CalculateTax(t *testing.T) {
	tc := &TaxCalculator{}

	tests := []struct {
		name         string
		income       float64
		year         int
		expectedTax  float64
		expectsError bool
		Tolerance    float64
	}{
		{
			name:         "Negative Income",
			income:       -100,
			year:         2023,
			expectedTax:  0,
			expectsError: true,
			Tolerance:    0.001,
		},
		{
			name:         "First Bracket - Below Max",
			income:       5000,
			year:         2023,
			expectedTax:  500, // 5000 * 0.10
			expectsError: false,
			Tolerance:    0.001,
		},
		{
			name:         "First Bracket - At Max",
			income:       10000,
			year:         2023,
			expectedTax:  1000, // 10000 * 0.10
			expectsError: false,
			Tolerance:    0.001,
		},
		{
			name:   "Second Bracket - Just Above First Max (Bug Target)",
			income: 10001,
			year:   2023,
			// Expected tax *if fixed*: (10000 * 0.10) + (1 * 0.20) = 1000 + 0.20 = 1000.20
			// Current buggy code output: 10001 * 0.20 = 2000.20
			expectedTax:  1000.20,
			expectsError: false,
			Tolerance:    0.001,
		},
		{
			name:   "Second Bracket - Mid Range",
			income: 30000,
			year:   2023,
			// Expected tax *if fixed*: (10000 * 0.10) + (20000 * 0.20) = 1000 + 4000 = 5000
			// Current buggy code output: 30000 * 0.20 = 6000
			expectedTax:  5000,
			expectsError: false,
			Tolerance:    0.001,
		},
		{
			name:   "Second Bracket - At Max",
			income: 50000,
			year:   2023,
			// Expected tax *if fixed*: (10000 * 0.10) + (40000 * 0.20) = 1000 + 8000 = 9000
			// Current buggy code output: 50000 * 0.20 = 10000
			expectedTax:  9000,
			expectsError: false,
			Tolerance:    0.001,
		},
		{
			name:   "Third Bracket - Just Above Second Max",
			income: 50001,
			year:   2023,
			// Expected tax *if fixed*: (10000 * 0.10) + (40000 * 0.20) + (1 * 0.30) = 1000 + 8000 + 0.30 = 9000.30
			// Current buggy code output: 50001 * 0.30 = 15000.30
			expectedTax:  9000.30,
			expectsError: false,
			Tolerance:    0.001,
		},
		{
			name:   "Third Bracket - Mid Range",
			income: 70000,
			year:   2023,
			// Expected tax *if fixed*: (10000 * 0.10) + (40000 * 0.20) + (20000 * 0.30) = 1000 + 8000 + 6000 = 15000
			// Current buggy code output: 70000 * 0.30 = 21000
			expectedTax:  15000,
			expectsError: false,
			Tolerance:    0.001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualTax, err := tc.CalculateTax(tt.income, tt.year)

			if tt.expectsError {
				if err == nil {
					t.Errorf("Test %s: Expected an error but got none", tt.name)
				}
			} else {
				if err != nil {
					t.Fatalf("Test %s: Did not expect an error but got: %v", tt.name, err)
				}
				if math.Abs(actualTax-tt.expectedTax) > tt.Tolerance {
					t.Errorf("Test %s: For income %.2f, expected tax %.2f but got %.2f",
						t.name, tt.income, tt.expectedTax, actualTax)
				}
			}
		})
	}
}
