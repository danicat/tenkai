package main

import (
	"testing"
)

func TestCalculateTax(t *testing.T) {
	tc := &TaxCalculator{}

	tests := []struct {
		name        string
		income      float64
		expectedTax float64
	}{
		{name: "Zero Income", income: 0, expectedTax: 0.0},
		{name: "Negative Income", income: -100, expectedTax: 0.0},
		{name: "Income in first bracket (below 10k)", income: 5000, expectedTax: 500.0},
		{name: "Income at first bracket boundary (10k)", income: 10000, expectedTax: 1000.0},
		{name: "Income just over first bracket (10k + 1)", income: 10001, expectedTax: 1000 + 0.20},
		{name: "Income in second bracket (below 50k)", income: 30000, expectedTax: 1000 + (20000 * 0.20)},    // 1000 + 4000 = 5000
		{name: "Income at second bracket boundary (50k)", income: 50000, expectedTax: 1000 + (40000 * 0.20)}, // 1000 + 8000 = 9000
		{name: "Income just over second bracket (50k + 1)", income: 50001, expectedTax: 9000 + 0.30},         // 9000 + 0.30 = 9000.30
		{name: "Income in third bracket (over 50k)", income: 70000, expectedTax: 9000 + (20000 * 0.30)},      // 9000 + 6000 = 15000
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualTax := tc.CalculateTax(tt.income)
			if actualTax != tt.expectedTax {
				t.Errorf("For income %f, expected tax %f, but got %f", tt.income, tt.expectedTax, actualTax)
			}
		})
	}
}
