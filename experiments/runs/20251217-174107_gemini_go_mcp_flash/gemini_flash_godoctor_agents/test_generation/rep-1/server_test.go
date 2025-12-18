package main

import (
	"testing"
)

func TestCalculateTax(t *testing.T) {
	tests := []struct {
		name     string
		income   float64
		year     int
		expected float64
	}{
		{
			name:     "Income within first bracket (non-leap year)",
			income:   15000,
			year:     2023,
			expected: 1500, // 15000 * 0.10
		},
		{
			name:     "Income at first bracket boundary (non-leap year)",
			income:   20000,
			year:     2023,
			expected: 2000, // 20000 * 0.10
		},
		{
			name:     "Income slightly above first bracket (non-leap year)",
			income:   20001,
			year:     2023,
			expected: 2000 + (1 * 0.20), // 20000 * 0.10 + 1 * 0.20
		},
		{
			name:     "Income within second bracket (non-leap year)",
			income:   35000,
			year:     2023,
			expected: 2000 + (15000 * 0.20), // 20000 * 0.10 + (35000-20000) * 0.20 = 2000 + 3000 = 5000
		},
		{
			name:     "Income at second bracket boundary (non-leap year)",
			income:   50000,
			year:     2023,
			expected: 2000 + (30000 * 0.20), // 20000 * 0.10 + (50000-20000) * 0.20 = 2000 + 6000 = 8000
		},
		{
			name:     "Income slightly above second bracket (non-leap year)",
			income:   50001,
			year:     2023,
			expected: 8000 + (1 * 0.30), // 8000 + 1 * 0.30
		},
		{
			name:     "Income within third bracket (non-leap year)",
			income:   70000,
			year:     2023,
			expected: 8000 + (20000 * 0.30), // 8000 + (70000-50000) * 0.30 = 8000 + 6000 = 14000
		},
		{
			name:     "Income within first bracket (leap year)",
			income:   15000,
			year:     2024,
			expected: 1500, // Expected *without* +100 bonus, so this should fail
		},
		{
			name:     "Income at second bracket boundary (leap year)",
			income:   50000,
			year:     2024,
			expected: 8000, // Expected *without* +100 bonus, so this should fail
		},
		{
			name:     "Income within third bracket (leap year)",
			income:   70000,
			year:     2024,
			expected: 14000, // Expected *without* +100 bonus, so this should fail
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateTax(tt.income, tt.year)
			if got != tt.expected {
				t.Errorf("CalculateTax(%f, %d) got %f, want %f", tt.income, tt.year, got, tt.expected)
			}
		})
	}
}
