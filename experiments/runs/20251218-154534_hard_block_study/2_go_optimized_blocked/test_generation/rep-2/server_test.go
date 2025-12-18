package main

import (
	"testing"
)

func TestCalculateTax(t *testing.T) {
	tests := []struct {
		name    string
		income  float64
		year    int
		wantTax float64
		wantErr bool
	}{
		{
			name:    "Negative income",
			income:  -100,
			year:    2025,
			wantTax: 0,
			wantErr: true,
		},
		{
			name:    "Income in first bracket (0%)",
			income:  5000,
			year:    2025,
			wantTax: 0,
			wantErr: false,
		},
		{
			name:    "Income at first bracket upper bound (0%)",
			income:  10000,
			year:    2025,
			wantTax: 0,
			wantErr: false,
		},
		{
			name:    "Income just above first bracket (10%)",
			income:  10001,
			year:    2025,
			wantTax: (10001 - 10000) * 0.10, // Expected: 0.10
			wantErr: false,
		},
		{
			name:    "Income in second bracket (10%)",
			income:  30000,
			year:    2025,
			wantTax: (30000 - 10000) * 0.10, // Expected: 2000.00
			wantErr: false,
		},
		{
			name:    "Income at second bracket upper bound (10%)",
			income:  50000,
			year:    2025,
			wantTax: (50000 - 10000) * 0.10, // Expected: 4000.00
			wantErr: false,
		},
		{
			name:    "Income just above second bracket (20%)",
			income:  50001,
			year:    2025,
			wantTax: (50000-10000)*0.10 + (50001-50000)*0.20, // Expected: 4000 + 0.20 = 4000.20
			wantErr: false,
		},
		{
			name:    "Income in third bracket (20%)",
			income:  60000,
			year:    2025,
			wantTax: (50000-10000)*0.10 + (60000-50000)*0.20, // Expected: 4000 + 2000 = 6000.00
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tax, err := CalculateTax(tt.income, tt.year)

			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateTax() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tax != tt.wantTax {
				t.Errorf("CalculateTax() got tax = %v, want %v", tax, tt.wantTax)
			}
		})
	}
}
