package main

import (
	"fmt"
	"math"
	"testing"
)

func TestCalculateTax(t *testing.T) {
	tests := []struct {
		name        string
		income      float64
		year        int
		expectedTax float64
		expectError bool
	}{
		// Negative income test
		{name: "Negative Income", income: -1000, year: 2023, expectedTax: 0, expectError: true},

		// First bracket tests (10%)
		{name: "Income 0", income: 0, year: 2023, expectedTax: 0, expectError: false},
		{name: "Income in first bracket", income: 10000, year: 2023, expectedTax: 1000, expectError: false},
		{name: "Income at first bracket boundary", income: 20000, year: 2023, expectedTax: 2000, expectError: false},

		// Second bracket tests (15%)
		{name: "Income slightly above first bracket", income: 20001, year: 2023, expectedTax: 2000 + (1 * 0.15), expectError: false},
		{name: "Income in second bracket", income: 35000, year: 2023, expectedTax: 2000 + (15000 * 0.15), expectError: false},
		{name: "Income at second bracket boundary", income: 50000, year: 2023, expectedTax: 2000 + (30000 * 0.15), expectError: false},

		// Third bracket tests (20%)
		{name: "Income slightly above second bracket", income: 50001, year: 2023, expectedTax: 2000 + (30000 * 0.15) + (1 * 0.20), expectError: false},
		{name: "Income in third bracket", income: 75000, year: 2023, expectedTax: 2000 + (30000 * 0.15) + (25000 * 0.20), expectError: false},

		// Leap year bonus tests (assuming 500 deduction for income > 20000)
		{name: "Leap Year 2000 (income 30000)", income: 30000, year: 2000, expectedTax: 2000 + (10000 * 0.15) - 500, expectError: false},
		{name: "Leap Year 2004 (income 30000)", income: 30000, year: 2004, expectedTax: 2000 + (10000 * 0.15) - 500, expectError: false},
		{name: "Non-Leap Year 1900 (income 30000)", income: 30000, year: 1900, expectedTax: 2000 + (10000 * 0.15), expectError: false},
		{name: "Non-Leap Year 2100 (income 30000)", income: 30000, year: 2100, expectedTax: 2000 + (10000 * 0.15), expectError: false},
		{name: "Non-Leap Year 2023 (income 30000)", income: 30000, year: 2023, expectedTax: 2000 + (10000 * 0.15), expectError: false},
		{name: "Leap Year 2024 (income 30000)", income: 30000, year: 2024, expectedTax: 2000 + (10000 * 0.15) - 500, expectError: false},
		{name: "Leap Year 2024 (income 15000 - no bonus)", income: 15000, year: 2024, expectedTax: 1500, expectError: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tax, err := CalculateTax(tt.income, tt.year)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected an error for income %f, got nil", tt.income)
				}
			} else {
				if err != nil {
					t.Errorf("did not expect an error for income %f, got %v", tt.income, err)
				}
				// Use a small epsilon for float comparison
				if math.Abs(tax-tt.expectedTax) > 0.001 {
					t.Errorf("For income %f, year %d, expected tax %.2f, got %.2f", tt.income, tt.year, tt.expectedTax, tax)
				}
			}
		})
	}
}

// Test for isLeapYear function
func TestIsLeapYear(t *testing.T) {
	tests := []struct {
		year     int
		expected bool
	}{
		{2000, true},  // Divisible by 400, so it's a leap year
		{2004, true},  // Divisible by 4 but not by 100, so it's a leap year
		{1900, false}, // Divisible by 100 but not by 400, so it's NOT a leap year
		{2100, false}, // Divisible by 100 but not by 400, so it's NOT a leap year
		{2023, false}, // Not divisible by 4, so it's NOT a leap year
		{2024, true},  // Divisible by 4, so it's a leap year
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Year_%d", tt.year), func(t *testing.T) {
			if isLeapYear(tt.year) != tt.expected {
				t.Errorf("isLeapYear(%d): expected %t, got %t", tt.year, tt.expected, isLeapYear(tt.year))
			}
		})
	}
}
