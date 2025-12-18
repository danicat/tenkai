package main

import (
	"math"
	"testing"
)

func TestCalculateTax(t *testing.T) {
	tests := []struct {
		name    string
		income  float64
		expectedTax float64
		expectError bool
	}{
		{
			name:        "Zero income",
			income:      0.0,
			expectedTax: 0.0,
			expectError: false,
		},
		{
			name:        "Income within first bracket (0%)",
			income:      5000.0,
			expectedTax: 0.0,
			expectError: false,
		},
		{
			name:        "Income exactly at first bracket limit (0%)",
			income:      10000.0,
			expectedTax: 0.0,
			expectError: false,
		},
		{
			name:        "Income in second bracket (10%)",			
			income:      20000.0, // 10000 * 0% + (20000-10000) * 10% = 1000
			expectedTax: 1000.0,
			expectError: false,
		},
		{
			name:        "Income at second bracket limit (10%)",
			income:      50000.0, // 10000 * 0% + (50000-10000) * 10% = 4000
			expectedTax: 4000.0,
			expectError: false,
		},
		{
			name:        "Income in third bracket (20%)",
			income:      75000.0, // 10000*0% + 40000*10% + (75000-50000)*20% = 0 + 4000 + 5000 = 9000
			expectedTax: 9000.0,
			expectError: false,
		},
		{
			name:        "Income at third bracket limit (20%)",
			income:      100000.0, // 10000*0% + 40000*10% + 50000*20% = 0 + 4000 + 10000 = 14000
			expectedTax: 14000.0,
			expectError: false,
		},
		{
			name:        "Income in fourth bracket (30%)",
			income:      150000.0, // 14000 + (150000-100000)*30% = 14000 + 15000 = 29000
			expectedTax: 29000.0,
			expectError: false,
		},
		{
			name:        "Negative income",
			income:      -100.0,
			expectedTax: 0.0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			tax, err := CalculateTax(tt.income)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected an error for negative income, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error, but got: %v", err)
				}
				if math.Abs(tax-tt.expectedTax) > 0.001 {
					t.Errorf("For income %.2f, expected tax %.2f, but got %.2f", tt.income, tt.expectedTax, tax)
				}
			}
		})
	}
}
