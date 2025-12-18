package main

import (
	"math"
	"testing"
)

func TestCalculateTax(t *testing.T) {
	// Define a standard set of tax brackets for testing.
	// Note: The bug in server.go will cause some of these tests to fail initially.
	standardBrackets := []TaxBracket{
		{Rate: 0.10, Max: 10000.0}, // 10% for income up to 10,000
		{Rate: 0.20, Max: 30000.0}, // 20% for income between 10,000.01 and 30,000
		{Rate: 0.30, Max: 60000.0}, // 30% for income between 30,000.01 and 60,000
		{Rate: 0.40, Max: 0.0},     // 40% for income above 60,000 (Max: 0 indicates no upper limit)
	}

	tc := &TaxCalculator{Brackets: standardBrackets}

	tests := []struct {
		name          string
		income        float64
		expectedTax float64
	}{
		{
			name:          "Zero Income",
			income:        0.0,
			expectedTax: 0.0,
		},
		{
			name:          "Income in first bracket - lower end",
			income:        5000.0,
			expectedTax: 500.0, // 5000 * 0.10
		},
		{
			name:          "Income in first bracket - upper end",
			income:        10000.0,
			expectedTax: 1000.0, // 10000 * 0.10
		},
		{
			name:          "Income crossing first bracket",
			income:        15000.0,
			expectedTax: 2000.0, // (10000 * 0.10) + (5000 * 0.20) = 1000 + 1000
		},
		{
			name:          "Income at second bracket upper end",
			income:        30000.0,
			expectedTax: 5000.0, // (10000 * 0.10) + (20000 * 0.20) = 1000 + 4000
		},
		{
			name:          "Income crossing second bracket",
			income:        40000.0,
			expectedTax: 8000.0, // (10000 * 0.10) + (20000 * 0.20) + (10000 * 0.30) = 1000 + 4000 + 3000
		},
		{
			name:          "Income at third bracket upper end",
			income:        60000.0,
			expectedTax: 14000.0, // (10000 * 0.10) + (20000 * 0.20) + (30000 * 0.30) = 1000 + 4000 + 9000
		},
		{
			name:          "Income in highest bracket",
			income:        70000.0,
			expectedTax: 18000.0, // (10000 * 0.10) + (20000 * 0.20) + (30000 * 0.30) + (10000 * 0.40) = 1000 + 4000 + 9000 + 4000
		},
		{
			name:          "Large Income in highest bracket",
			income:        100000.0,
			expectedTax: 30000.0, // (10000 * 0.10) + (20000 * 0.20) + (30000 * 0.30) + (40000 * 0.40) = 1000 + 4000 + 9000 + 16000
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualTax := tc.CalculateTax(tt.income)
			// Use a small epsilon for float comparison due to potential precision issues
			if math.Abs(actualTax-tt.expectedTax) > 0.001 {
				t.Errorf("For income %f, expected tax %f, but got %f", tt.income, tt.expectedTax, actualTax)
			}
		})
	}
}
