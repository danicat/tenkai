package main

import (
	"testing"
)

func TestCalculateTax(t *testing.T) {
	tests := []struct {
		name   string
		income float64
		want   float64
	}{
		{
			name:   "Income Zero",
			income: 0,
			want:   0,
		},
		{
			name:   "First Bracket Max",
			income: 10000,
			want:   1000.0,
		},
		{
			name:   "Just Over First Bracket",
			income: 10001,
			want:   1000.20,
		},
		{
			name:   "Mid Second Bracket",
			income: 20000,
			want:   3000.0,
		},
		{
			name:   "Second Bracket Max",
			income: 30000,
			want:   5000.0,
		},
		{
			name:   "Just Over Second Bracket",
			income: 30001,
			want:   5000.30,
		},
		{
			name:   "Well Into Third Bracket",
			income: 50000,
			want:   11000.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateTax(tt.income)
			if got != tt.want {
				t.Errorf("CalculateTax(%f) got %f, want %f", tt.income, got, tt.want)
			}
		})
	}
}
