package main

import (
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
		// Non-leap year tests (2023)
		{"Income_5k_2023", 5000, 2023, 500}, // 5000 * 0.10 = 500
		{"Income_10k_2023", 10000, 2023, 1000}, // 10000 * 0.10 = 1000
		{"Income_10k_1_2023", 10001, 2023, 1000 + 0.20}, // 10000*0.10 + 1*0.20 = 1000.20
		{"Income_20k_2023", 20000, 2023, 1000 + (10000 * 0.20)}, // 10000*0.10 + 10000*0.20 = 1000 + 2000 = 3000
		{"Income_50k_2023", 50000, 2023, 1000 + (40000 * 0.20)}, // 10000*0.10 + 40000*0.20 = 1000 + 8000 = 9000
		{"Income_50k_1_2023", 50001, 2023, 9000 + 0.30}, // 10000*0.10 + 40000*0.20 + 1*0.30 = 9000.30
		{"Income_60k_2023", 60000, 2023, 9000 + (10000 * 0.30)}, // 10000*0.10 + 40000*0.20 + 10000*0.30 = 9000 + 3000 = 12000
		{"Income_0_2023", 0, 2023, 0},
		{"Income_negative_2023", -100, 2023, 0},

		// Leap year tests (2024)
		{"Income_5k_2024_Leap", 5000, 2024, 500 - 50}, // (5000 * 0.10) - 50 = 450
		{"Income_10k_2024_Leap", 10000, 2024, 1000 - 50}, // (10000 * 0.10) - 50 = 950
		{"Income_20k_2024_Leap", 20000, 2024, 3000 - 50}, // Correct marginal tax (3000) - 50 = 2950
		{"Income_60k_2024_Leap", 60000, 2024, 12000 - 50}, // Correct marginal tax (12000) - 50 = 11950
		{"Income_400_2024_Leap_MinTax", 400, 2024, 0}, // (400*0.10) - 50 = 40 - 50 = -10, clamped to 0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tc.CalculateTax(tt.income, tt.year)
			if got != tt.want {
				t.Errorf("CalculateTax(income=%.2f, year=%d) = %.2f; want %.2f", tt.income, tt.year, got, tt.want)
			}
		})
	}
}

// Test for statement coverage (e.g., specific leap year rules, zero/negative income)
func TestCalculateTaxCoverage(t *testing.T) {
	tc := &TaxCalculator{}

	tests := []struct {
		name   string
		income float64
		year   int
		want   float64
	}{
		// Specific leap year checks
		{"Year_1900_NotLeap", 1000, 1900, 100}, // 1900 not a leap year
		{"Year_2000_Leap", 1000, 2000, 100 - 50}, // 2000 is a leap year
		{"Income_just_above_0_Leap", 1, 2024, 0}, // 1*0.1 - 50 = -49.9, clamped to 0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tc.CalculateTax(tt.income, tt.year)
			if got != tt.want {
				t.Errorf("CalculateTax(income=%.2f, year=%d) = %.2f; want %.2f", tt.income, tt.year, got, tt.want)
			}
		})
	}
}
