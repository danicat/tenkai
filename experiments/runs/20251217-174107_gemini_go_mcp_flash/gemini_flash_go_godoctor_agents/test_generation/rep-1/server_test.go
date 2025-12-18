package main

import (
	"fmt"
	"testing"
)

func TestCalculateTax(t *testing.T) {
	tc := &TaxCalculator{}

	tests := []struct {
		name    string
		income  float64
		year    int
		wantTax float64
		wantErr bool
	}{
		{
			name:    "Income below first bracket",
			income:  5000,
			year:    2023,
			wantTax: 500, // 5000 * 0.10
			wantErr: false,
		},
		{
			name:    "Income at first bracket boundary",
			income:  10000,
			year:    2023,
			wantTax: 1000, // 10000 * 0.10
			wantErr: false,
		},
		{
			name:    "Income in second bracket - progressive calculation expected",
			income:  20000,
			year:    2023,
			wantTax: 3000, // (10000 * 0.10) + (10000 * 0.20)
			wantErr: false,
		},
		{
			name:    "Income at second bracket boundary - progressive calculation expected",
			income:  50000,
			year:    2023,
			wantTax: 9000, // (10000 * 0.10) + (40000 * 0.20)
			wantErr: false,
		},
		{
			name:    "Income in third bracket - progressive calculation expected",
			income:  60000,
			year:    2023,
			wantTax: 12000, // (10000 * 0.10) + (40000 * 0.20) + (10000 * 0.30)
			wantErr: false,
		},
		{
			name:    "Zero income",
			income:  0,
			year:    2023,
			wantTax: 0,
			wantErr: false,
		},
		{
			name:    "Negative income should return error",
			income:  -100,
			year:    2023,
			wantTax: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTax, err := tc.CalculateTax(tt.income, tt.year)

			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateTax() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && gotTax != tt.wantTax {
				t.Errorf("CalculateTax() gotTax = %v, want %v", gotTax, tt.wantTax)
			}
		})
	}
}

// Example usage (for documentation/demonstration)
func ExampleTaxCalculator_CalculateTax() {
	tc := &TaxCalculator{}
	tax, err := tc.CalculateTax(75000, 2023)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Tax for 75000: %.2f\n", tax)
	}
	// Output: Tax for 75000: 16500.00
}
