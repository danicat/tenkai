package main

import (
	"math"
	"testing"
)

func TestCalculateTax(t *testing.T) {
	tc := NewTaxCalculator()

	tests := []struct {
		name   string
		income float64
		expect float64
	}{
		{"Zero Income", 0.0, 0.0},
		{"Negative Income", -100.0, 0.0}, // Expect 0 for negative income
		{"Income below first bracket", 5000.0, 500.0}, // 10% of 5000
		{"Income at first bracket boundary", 10000.0, 1000.0}, // 10% of 10000
		{"Income in second bracket low", 15000.0, 1750.0}, // 1000 (10% of 10k) + 750 (15% of 5k)
		{"Income in second bracket high", 40000.0, 5500.0}, // 1000 (10% of 10k) + 4500 (15% of 30k)
		{"Income at second bracket boundary", 50000.0, 7000.0}, // 1000 (10% of 10k) + 6000 (15% of 40k)
		{"Income in third bracket low", 60000.0, 9000.0}, // 1000 (10% of 10k) + 6000 (15% of 40k) + 2000 (20% of 10k)
		{"Income in third bracket high", 90000.0, 15000.0}, // 1000 (10% of 10k) + 6000 (15% of 40k) + 8000 (20% of 40k)
		{"Income at third bracket boundary", 100000.0, 17000.0}, // 1000 (10% of 10k) + 6000 (15% of 40k) + 10000 (20% of 50k)
		{"Income in fourth bracket low (added by problem desc)", 110000.0, 19500.0}, // 1000 (10% of 10k) + 6000 (15% of 40k) + 10000 (20% of 50k) + 2500 (25% of 10k)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tc.CalculateTax(tt.income)
			// Use a small epsilon for float comparison
			if math.Abs(got-tt.expect) > 0.001 {
				t.Errorf("CalculateTax(%.2f) = %.2f; want %.2f", tt.income, got, tt.expect)
			}
		})
	}
}
