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
		wantTax float64
	}{
		{"zero income", 0, 0},
		{"negative income", -100, 0},
		{"income in first bracket", 5000, 500},
		{"income at first bracket boundary", 10000, 1000},
		{"income just above first bracket", 10001, 1000 + (1 * 0.20)}, // 1000.2
		{"income in second bracket", 30000, 1000 + (20000 * 0.20)}, // 5000
		{"income at second bracket boundary", 50000, 1000 + (40000 * 0.20)}, // 9000
		{"income just above second bracket", 50001, 1000 + (40000 * 0.20) + (1 * 0.30)}, // 9000.3
		{"income in third bracket", 70000, 1000 + (40000 * 0.20) + (20000 * 0.30)}, // 15000
	}

	const epsilon = 1e-9 // For float comparison

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTax := tc.CalculateTax(tt.income)
			if math.Abs(gotTax-tt.wantTax) > epsilon {
				t.Errorf("CalculateTax(%.2f) = %.2f; want %.2f", tt.income, gotTax, tt.wantTax)
			}
		})
	}
}
