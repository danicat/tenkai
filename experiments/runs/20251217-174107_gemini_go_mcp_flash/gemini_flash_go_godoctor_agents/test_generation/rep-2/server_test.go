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
		year   int
		want   float64
	}{
		// Test cases for the first bracket (income <= 10000)
		{"income_0", 0.0, 2023, 0.0},
		{"income_10000_non_leap", 10000.0, 2023, 1000.0}, // 10000 * 0.10
		{"income_10000_leap", 10000.0, 2024, 1100.0},     // 10000 * 0.10 + 100 (leap year bonus)

		// Test cases for the second bracket (10001 <= income <= 50000)
		// Expecting bugs here due to `income-10000` vs `income-10001` and `10000` vs `10001` for the second bracket start
		{"income_10001_non_leap", 10001.0, 2023, 1000.2}, // 10000*0.10 + (1)*0.20 = 1000 + 0.2 = 1000.2 (BUG: should be 1000 + (1)*0.2 = 1000.2, if it starts at 10000. The code has (income - 10000) for the second bracket, so for 10001 it's 1000 + (1)*0.2 = 1000.2)
		{"income_50000_non_leap", 50000.0, 2023, 9000.0}, // 10000*0.10 + 40000*0.20 = 1000 + 8000 = 9000.0
		{"income_50000_leap", 50000.0, 2024, 9100.0},     // 9000.0 + 100 (leap year bonus)

		// Test cases for the third bracket (income > 50000)
		// Expecting bugs here due to `income-50000` boundary
		{"income_50001_non_leap", 50001.0, 2023, 9000.3}, // 1000+8000 + (1)*0.30 = 9000 + 0.3 = 9000.3
		{"income_60000_non_leap", 60000.0, 2023, 12000.0},// 1000+8000 + 10000*0.30 = 9000 + 3000 = 12000.0
		{"income_60000_leap_2000", 60000.0, 2000, 12100.0},// 12000.0 + 100 (2000 is leap year)

		// Leap year specific tests (targeting the `else if year%4 == 0` bug)
		{"leap_year_1900_non_leap_century", 20000.0, 1900, 3000.0}, // 1900 is NOT a leap year. Tax = 3000.0 (10000*0.10 + 10000*0.20)
		{"leap_year_2000_leap_century", 20000.0, 2000, 3100.0},    // 2000 IS a leap year. Tax = 3000.0 + 100 = 3100.0
		{"non_leap_year_2023", 20000.0, 2023, 3000.0},          // 2023 is NOT a leap year. Tax = 3000.0
		{"non_leap_year_2022", 20000.0, 2022, 3000.0},          // 2022 is NOT a leap year. Tax = 3000.0
		{"leap_year_2024", 20000.0, 2024, 3100.0},             // 2024 IS a leap year. Tax = 3000.0 + 100 = 3100.0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tc.CalculateTax(tt.income, tt.year)
			// Using a small epsilon for float comparison
			if math.Abs(got-tt.want) > 0.001 {
				t.Errorf("CalculateTax(income=%.2f, year=%d) got %.2f, want %.2f", tt.income, tt.year, got, tt.want)
			}
		})
	}
}

// Test for the isLeapYear helper function (to ensure it's correct before relying on it)
func TestIsLeapYear(t *testing.T) {
	tests := []struct {
		name string
		year int
		want bool
	}{
		{"2000_leap", 2000, true},
		{"1900_not_leap", 1900, false},
		{"2004_leap", 2004, true},
		{"2023_not_leap", 2023, false},
		{"2022_not_leap", 2022, false},
		{"1600_leap", 1600, true},
		{"1700_not_leap", 1700, false},
		{"2400_leap", 2400, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isLeapYear(tt.year)
			if got != tt.want {
				t.Errorf("isLeapYear(year=%d) got %v, want %v", tt.year, got, tt.want)
			}
		})
	}
}
