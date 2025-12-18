package main

import (
	"math"
	"testing"
)

func TestCalculateTax(t *testing.T) {
	tests := []struct {
		name        string
		income      float64
		year        int
		expectedTax float64
	}{
		// Basic bracket tests
		{"Income 0", 0, 2023, 0},
		{"Income 5000", 5000, 2023, 500},    // 10% of 5000
		{"Income 10000", 10000, 2023, 1000}, // 10% of 10000

		// Bug 1: Income between 10001 and 50000 - boundary conditions
		{"Income 10001", 10001, 2023, 1000 + 0.20},           // 10% of 10000 + 20% of 1
		{"Income 30000", 30000, 2023, 1000 + (20000 * 0.20)}, // 10% of 10000 + 20% of 20000 = 1000 + 4000 = 5000
		{"Income 49999", 49999, 2023, 1000 + (39999 * 0.20)}, // 10% of 10000 + 20% of 39999 = 1000 + 7999.8 = 8999.8
		{"Income 50000", 50000, 2023, 1000 + (40000 * 0.20)}, // 10% of 10000 + 20% of 40000 = 1000 + 8000 = 9000

		// Income > 50000
		{"Income 50001", 50001, 2023, 9000 + 0.30},             // Tax on 50000 + 30% of 1 = 9000 + 0.30 = 9000.30
		{"Income 100000", 100000, 2023, 9000 + (50000 * 0.30)}, // Tax on 50000 + 30% of 50000 = 9000 + 15000 = 24000

		// Bug 2: Leap year handling
		{"Leap Year 2024, Income 30000", 30000, 2024, 5000 + 100},                              // Should be 5000 + 100 = 5100
		{"Leap Year 2024, Income 30000.01", 30000.01, 2024, 5000.002 + 100},                    // Should be (1000 + (20000.01 * 0.20)) + 100 = 5000.002 + 100 = 5100.002
		{"Leap Year 2020, Income 40000", 40000, 2020, 1000 + (30000 * 0.20) + 100},             // 1000 + 6000 + 100 = 7100
		{"Non-Leap Year 2023, Income 40000", 40000, 2023, 1000 + (30000 * 0.20)},               // 1000 + 6000 = 7000
		{"Leap Year 2024, Income 20000 (below threshold)", 20000, 2024, 1000 + (10000 * 0.20)}, // 1000 + 2000 = 3000 (no extra 100)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTax := CalculateTax(tt.income, tt.year)
			// Using a small epsilon for float comparisons
			if math.Abs(gotTax-tt.expectedTax) > 0.001 {
				t.Errorf("CalculateTax(%f, %d) = %f; want %f", tt.income, tt.year, gotTax, tt.expectedTax)
			}
		})
	}
}

// TestIsLeapYear is a utility test for the isLeapYear function to ensure its correctness.
func TestIsLeapYear(t *testing.T) {
	tests := []struct {
		year     int
		expected bool
	}{
		{2000, true},  // Divisible by 400
		{2004, true},  // Divisible by 4, not by 100
		{1900, false}, // Divisible by 100, not by 400
		{2003, false}, // Not divisible by 4
		{2024, true},
		{2023, false},
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.year)), func(t *testing.T) {
			if got := isLeapYear(tt.year); got != tt.expected {
				t.Errorf("isLeapYear(%d) = %t; want %t", tt.year, got, tt.expected)
			}
		})
	}
}
