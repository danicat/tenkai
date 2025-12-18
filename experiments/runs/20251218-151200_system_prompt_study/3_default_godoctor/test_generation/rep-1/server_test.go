package main

import (
	"testing"
)

func TestCalculateTax(t *testing.T) {
	tests := []struct {
		name   string
		income float64
		expectedTax float64
	}{
		// Test cases for Bug 1: Off-by-one error in the first bracket (income <= 10000)
		{"Income 0", 0, 0.0},
		{"Income 1000", 1000, 100.0}, // 1000 * 0.10
		{"Income 10000 (edge of first bracket)", 10000, 1000.0}, // 10000 * 0.10
		{"Income 10001 (start of second bracket)", 10001, 1000 + (10001-10000)*0.20}, // Expected: 1000 + 0.20 = 1000.20. Actual will be 10001 * 0.20 = 2000.2

		// Test cases for Bug 2: Incorrect calculation for the second bracket (income <= 30000)
		// Expected: 1000 (from first bracket) + (income - 10000) * 0.20
		{"Income 20000", 20000, 1000 + 10000*0.20}, // Expected: 1000 + 2000 = 3000. Actual will be 20000 * 0.20 = 4000
		{"Income 30000 (edge of second bracket)", 30000, 1000 + 20000*0.20}, // Expected: 1000 + 4000 = 5000. Actual will be 30000 * 0.20 = 6000
		{"Income 30001 (start of third bracket)", 30001, 5000 + (30001-30000)*0.30}, // Expected: 5000 + 0.30 = 5000.30. Actual will be 30001 * 0.30 = 9000.3

		// Test cases for Bug 3: Incorrect calculation for the third bracket (income <= 60000)
		// Expected: 1000 (first) + 4000 (second) + (income - 30000) * 0.30
		{"Income 45000", 45000, 5000 + 15000*0.30}, // Expected: 5000 + 4500 = 9500. Actual will be 45000 * 0.30 = 13500
		{"Income 60000 (edge of third bracket)", 60000, 5000 + 30000*0.30}, // Expected: 5000 + 9000 = 14000. Actual will be 60000 * 0.30 = 18000
		{"Income 60001 (start of fourth bracket)", 60001, 14000 + (60001-60000)*0.40}, // Expected: 14000 + 0.40 = 14000.40. Actual will be 60001 * 0.40 = 24000.4

		// Test cases for Bug 4: Incorrect calculation for the highest bracket (income > 60000)
		// Expected: 1000 (first) + 4000 (second) + 9000 (third) + (income - 60000) * 0.40
		{"Income 80000", 80000, 14000 + 20000*0.40}, // Expected: 14000 + 8000 = 22000. Actual will be 80000 * 0.40 = 32000
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actualTax := CalculateTax(tt.income)
			if actualTax != tt.expectedTax {
				t.Errorf("For income %.2f, expected tax %.2f, but got %.2f", tt.income, tt.expectedTax, actualTax)
			}
		})
	}
}
