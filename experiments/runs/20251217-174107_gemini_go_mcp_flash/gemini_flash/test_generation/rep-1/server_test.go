package main

import (
	"testing"
)

func TestCalculateTax(t *testing.T) {
	tc := NewTaxCalculator()

	tests := []struct {
		name        string
		income      float64
		expectedTax float64
	}{
		{"Zero Income", 0, 0.0},
		{"First Bracket Mid", 5000, 500.0},
		{"First Bracket Max", 10000, 1000.0},
		{"Second Bracket Start", 10001, 1000.0 + (1 * 0.20)}, // 1000 from first bracket + 1 from second
		{"Second Bracket Mid", 30000, 1000.0 + (20000 * 0.20)}, // 1000 from first + 20000 from second
		{"Second Bracket Max", 50000, 1000.0 + (40000 * 0.20)}, // 1000 from first + 40000 from second
		{"Third Bracket Start", 50001, 1000.0 + (40000 * 0.20) + (1 * 0.30)}, // previous + 1 from third
		{"Third Bracket Mid", 75000, 1000.0 + (40000 * 0.20) + (25000 * 0.30)}, // previous + 25000 from third
		{"Third Bracket Max", 100000, 1000.0 + (40000 * 0.20) + (50000 * 0.30)}, // previous + 50000 from third
		{"Fourth Bracket Start", 100001, 1000.0 + (40000 * 0.20) + (50000 * 0.30) + (1 * 0.40)}, // previous + 1 from fourth
		{"Fourth Bracket Mid", 150000, 1000.0 + (40000 * 0.20) + (50000 * 0.30) + (50000 * 0.40)}, // previous + 50000 from fourth
		{"Negative Income", -100, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualTax := tc.CalculateTax(tt.income)
			if actualTax != tt.expectedTax {
				t.Errorf("For income %.2f, expected tax %.2f, but got %.2f", tt.income, tt.expectedTax, actualTax)
			}
		})
	}
}

func TestCalculateTax_NegativeTaxableInBracket(t *testing.T) {
	// Intentionally create a tax calculator with malformed brackets to hit the 'taxableInThisBracket < 0' case
	malformedTc := &TaxCalculator{
		Brackets: []TaxBracket{
			{Limit: 5000, Rate: 0.10},  // Taxes 5000 at 10% = 500
			{Limit: 3000, Rate: 0.20},  // This limit (3000) is less than the previous taxed amount (5000), causing negative taxableInThisBracket
			{Limit: -1, Rate: 0.30},
		},
	}

	// For an income of 6000:
	// 1. First bracket: 5000 * 0.10 = 500. taxedInPreviousBrackets = 5000.
	// 2. Second bracket: income (6000) > taxedInPreviousBrackets (5000).
	//    currentBracketUpperLimit = 3000.
	//    taxableInThisBracket = 3000 - 5000 = -2000.
	//    This should then be set to 0 by the `if taxableInThisBracket < 0` guard.
	//    totalTax remains 500. taxedInPreviousBrackets remains 5000.
	// 3. Third bracket: income (6000) > taxedInPreviousBrackets (5000).
	//    currentBracketUpperLimit = 6000 (since it's the highest bracket in this context and income is 6000).
	//    taxableInThisBracket = 6000 - 5000 = 1000.
	//    totalTax += 1000 * 0.30 = 300. totalTax = 500 + 300 = 800.
	//    taxedInPreviousBrackets += 1000 = 6000.
	// 4. income (6000) <= taxedInPreviousBrackets (6000). Break.
	// Expected tax for 6000 should be 500 (from first bracket) + 1000 * 0.30 (from third bracket) = 800.

	income := 6000.0
	expectedTax := 800.0
	actualTax := malformedTc.CalculateTax(income)

	if actualTax != expectedTax {
		t.Errorf("For income %.2f with malformed brackets, expected tax %.2f, but got %.2f", income, expectedTax, actualTax)
	}
}