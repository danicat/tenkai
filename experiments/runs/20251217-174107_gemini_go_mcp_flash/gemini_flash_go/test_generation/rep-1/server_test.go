package main

import (
	"errors"
	"testing"
)

// TestCalculateTax_Bugs tests the CalculateTax function with known bugs.
func TestCalculateTax_Bugs(t *testing.T) {
	tc := &TaxCalculator{
		Brackets: []struct {
			Min    float64
			Rate   float64
			IsMax  bool
		}{
			{Min: 0, Rate: 0.10},
			{Min: 10000, Rate: 0.20},
			{Min: 30000, Rate: 0.30, IsMax: true},
		},
	}

	tests := []struct {
		name          string
		income        float64
		year          int
		wantTax       float64
		wantErr       error
	}{
		{
			name:    "Income in first bracket",
			income:  5000,
			year:    2023,
			wantTax: 500.0, // Expected: 5000 * 0.10
			wantErr: nil,
		},
		{
			name:    "Income at first bracket boundary",
			income:  10000,
			year:    2023,
			wantTax: 1000.0, // Expected: 10000 * 0.10
			wantErr: nil,
		},
		{
			name:    "Income across first two brackets",
			income:  20000,
			year:    2023,
			wantTax: 3000.0, // Expected: (10000 * 0.10) + (10000 * 0.20)
			wantErr: nil,
		},
		{
			name:    "Income across all brackets",
			income:  40000,
			year:    2023,
			wantTax: 8000.0, // Expected: (10000 * 0.10) + (20000 * 0.20) + (10000 * 0.30)
			wantErr: nil,
		},
		{
			name:    "Negative income",
			income:  -100,
			year:    2023,
			wantTax: 0.0,
			wantErr: errors.New("income cannot be negative"),
		},
		{
			name:    "Leap year income (2024) - expecting bonus (bug in implementation)",
			income:  10000,
			year:    2024, // 2024 is a leap year
			wantTax: 900.0, // Correctly calculated: 10000 * 0.10 - 100 (bonus)
			wantErr: nil,
		},
		{
			name:    "Non-leap year income (2023)",
			income:  10000,
			year:    2023, // 2023 is not a leap year
			wantTax: 1000.0, // Expected: 10000 * 0.10
			wantErr: nil,
		},
		{
			name:    "Year 2100 (non-leap year, but buggy isLeapYear thinks it's leap)",
			income:  10000,
			year:    2100, // Correctly, 2100 is NOT a leap year
			wantTax: 1000.0, // Expected: 10000 * 0.10 (no bonus)
			wantErr: nil,
		},
		{
			name:    "Year 2000 (leap year, buggy isLeapYear thinks it's leap)",
			income:  10000,
			year:    2000, // Correctly, 2000 IS a leap year
			wantTax: 900.0, // Expected: 10000 * 0.10 - 100 (bonus)
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotTax, gotErr := tc.CalculateTax(tt.income, tt.year)

			if tt.wantErr != nil { // We expect an error
				if gotErr == nil {
					t.Errorf("CalculateTax() got no error, want %v", tt.wantErr)
					return
				}
				if gotErr.Error() != tt.wantErr.Error() { // Compare error messages directly
					t.Errorf("CalculateTax() error = %v, wantErr %v", gotErr, tt.wantErr)
					return
				}
			} else { // We expect no error
				if gotErr != nil {
					t.Errorf("CalculateTax() got error %v, want nil", gotErr)
					return
				}
			}

			if gotTax != tt.wantTax {
				t.Errorf("CalculateTax() gotTax = %v, want %v", gotTax, tt.wantTax)
			}
		})
	}
}
