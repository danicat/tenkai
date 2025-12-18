package main

import (
	"math"
	"testing"
)

func TestCalculateTax(t *testing.T) {
	tc := &TaxCalculator{}

	tests := []struct {
		name   string
		income float64
		expectedTax float64
	}{
		{
			name:   "Income within first bracket",
			income: 5000,
			expectedTax: 500.00,
		},
		{
			name:   "Income at upper bound of first bracket",
			income: 10000,
			expectedTax: 1000.00,
		},
		{
			name:   "Income just above first bracket",
			income: 10001,
			expectedTax: 1000.00 + (1 * 0.20), // 1000 (from 1st bracket) + 0.20 (from 2nd bracket)
		},
		{
			name:   "Income within second bracket",
			income: 15000,
			expectedTax: 2000.00, // 1000 (from 1st) + (5000 * 0.20) = 1000 + 1000 = 2000
		},
		{
			name:   "Income at upper bound of second bracket",
			income: 30000,
			expectedTax: 5000.00, // 1000 (from 1st) + (20000 * 0.20) = 1000 + 4000 = 5000
		},
		{
			name:   "Income just above second bracket",
			income: 30001,
			expectedTax: 5000.00 + (1 * 0.30), // 5000 (from 1st+2nd) + 0.30 (from 3rd)
		},
		{
			name:   "Income within third bracket",
			income: 40000,
			expectedTax: 8000.00, // 5000 (from 1st+2nd) + (10000 * 0.30) = 5000 + 3000 = 8000
		},
		{
			name:   "Income at upper bound of third bracket",
			income: 60000,
			expectedTax: 14000.00, // 5000 (from 1st+2nd) + (30000 * 0.30) = 5000 + 9000 = 14000
		},
		{
			name:   "Income just above third bracket",
			income: 60001,
			expectedTax: 14000.00 + (1 * 0.40), // 14000 (from 1st+2nd+3rd) + 0.40 (from 4th)
		},
		{
			name:   "Income in highest bracket",
			income: 70000,
			expectedTax: 18000.00, // 14000 (from 1st+2nd+3rd) + (10000 * 0.40) = 14000 + 4000 = 18000
		},
		{
			name:   "Zero income",
			income: 0,
			expectedTax: 0.00,
		},
		{
			name:   "Negative income (should be 0 tax)",
			income: -100,
			expectedTax: 0.00,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualTax := tc.CalculateTax(tt.income)
			// Using a small epsilon for float comparison
			if math.Abs(actualTax-tt.expectedTax) > 0.001 {
				t.Errorf("For income %f, expected tax %.2f, but got %.2f", tt.income, tt.expectedTax, actualTax)
			}
		})
	}
}