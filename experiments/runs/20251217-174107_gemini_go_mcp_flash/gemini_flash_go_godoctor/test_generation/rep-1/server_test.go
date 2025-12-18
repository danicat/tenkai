package main

import (
	"math"
	"testing"
)

// In reality, you'd define constants for tax brackets.
// For this exercise, we'll hardcode them in the expected values.

func TestTaxCalculator_CalculateTax(t *testing.T) {
	tests := []struct {
		name string
		income float64
		wantTax float64
	}{
		{
			name:    "Zero income",
			income:  0,
			wantTax: 0,
		},
		{
			name:    "Income within 0% bracket",
			income:  5000,
			wantTax: 0,
		},
		{
			name:    "Income at 0% bracket boundary",
			income:  10000,
			wantTax: 0,
		},
		// --- Tests for 10% bracket (income > 10000 and <= 40000) ---
		{
			name:    "Income just above 0% bracket (starts 10% bracket)",
			income:  10001,
			// Expected tax with bug: 0.10 * 10001 = 1000.1
			// Correct tax: 0 + 0.10 * (10001 - 10000) = 0.1
			wantTax: 0.1, // This test expects the CORRECT tax for later verification.
		},
		{
			name:    "Income in middle of 10% bracket",
			income:  25000,
			// Expected tax with bug: 0.10 * 25000 = 2500
			// Correct tax: 0 + 0.10 * (25000 - 10000) = 1500
			wantTax: 1500, // This test expects the CORRECT tax for later verification.
		},
		{
			name:    "Income at 10% bracket boundary",
			income:  40000,
			// Expected tax with bug: 0.10 * 40000 = 4000
			// Correct tax: 0 + 0.10 * (40000 - 10000) = 3000
			wantTax: 3000, // This test expects the CORRECT tax for later verification.
		},
		// --- Tests for 20% bracket (income > 40000 and <= 80000) ---
		{
			name:    "Income just above 10% bracket (starts 20% bracket)",
			income:  40001,
			// Expected tax with bug: 0.20 * 40001 = 8000.2
			// Correct tax: (3000 from previous) + 0.20 * (40001 - 40000) = 3000 + 0.20 = 3000.2
			wantTax: 3000.2, // This test expects the CORRECT tax for later verification.
		},
		{
			name:    "Income in middle of 20% bracket",
			income:  60000,
			// Expected tax with bug: 0.20 * 60000 = 12000
			// Correct tax: (3000 from previous) + 0.20 * (60000 - 40000) = 3000 + 4000 = 7000
			wantTax: 7000, // This test expects the CORRECT tax for later verification.
		},
		{
			name:    "Income at 20% bracket boundary",
			income:  80000,
			// Expected tax with bug: 0.20 * 80000 = 16000
			// Correct tax: (3000 from previous) + 0.20 * (80000 - 40000) = 3000 + 8000 = 11000
			wantTax: 11000, // This test expects the CORRECT tax for later verification.
		},
		// --- Tests for 30% bracket (income > 80000) ---
		{
			name:    "Income just above 20% bracket (starts 30% bracket)",
			income:  80001,
			// Expected tax with bug: 0.30 * 80001 = 24000.3
			// Correct tax: (11000 from previous) + 0.30 * (80001 - 80000) = 11000 + 0.30 = 11000.3
			wantTax: 11000.3, // This test expects the CORRECT tax for later verification.
		},
		{
			name:    "High income in 30% bracket",
			income:  100000,
			// Expected tax with bug: 0.30 * 100000 = 30000
			// Correct tax: (11000 from previous) + 0.30 * (100000 - 80000) = 11000 + 6000 = 17000
			wantTax: 17000, // This test expects the CORRECT tax for later verification.
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // Run tests in parallel if safe
			tc := &TaxCalculator{}
			gotTax := tc.CalculateTax(tt.income)

			// Using a small epsilon for float comparison due to potential precision issues
			if math.Abs(gotTax-tt.wantTax) > 0.001 {
				t.Errorf("CalculateTax(%f) = %f; want %f", tt.income, gotTax, tt.wantTax)
			}
		})
	}
}
