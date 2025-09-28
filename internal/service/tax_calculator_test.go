package service

import (
	"math"
	"testing"

	"github.com/alaaeelsayed/tax-calculator/internal/model"
)

type fakeTaxBracketProvider struct {
	brackets []model.TaxBracket
	err      error
}

func (f *fakeTaxBracketProvider) GetTaxBrackets(year string) ([]model.TaxBracket, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.brackets, nil
}

func TestCalculateTax(t *testing.T) {
	testBrackets := []model.TaxBracket{
		{Min: 0, Max: 50197, Rate: 0.15},
		{Min: 50197, Max: 100392, Rate: 0.205},
		{Min: 100392, Max: 155625, Rate: 0.26},
		{Min: 155625, Max: 221708, Rate: 0.29},
		{Min: 221708, Rate: 0.33},
	}

	tests := []struct {
		name                  string
		salary                float64
		expectedTotalTaxes    float64
		expectedEffectiveRate float64
		expectError           bool
	}{
		{
			name:                  "Negative salary",
			salary:                -1000,
			expectedTotalTaxes:    0,
			expectedEffectiveRate: 0,
			expectError:           true,
		},
		{
			name:                  "Zero salary",
			salary:                0,
			expectedTotalTaxes:    0,
			expectedEffectiveRate: 0,
			expectError:           false,
		},
		{
			name:                  "First bracket",
			salary:                50_000,
			expectedTotalTaxes:    7_500.00,
			expectedEffectiveRate: 0.15,
			expectError:           false,
		},
		{
			name:                  "Second bracket",
			salary:                100_000,
			expectedTotalTaxes:    17_739.17,
			expectedEffectiveRate: 0.1774,
			expectError:           false,
		},
		{
			name:                  "Third bracket",
			salary:                1_234_567,
			expectedTotalTaxes:    385_587.65,
			expectedEffectiveRate: 0.3123,
			expectError:           false,
		},
	}

	mockTaxBracketProvider := &fakeTaxBracketProvider{
		brackets: testBrackets,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewTaxCalculatorService(mockTaxBracketProvider).CalculateTax(tt.salary, "")

			if tt.expectError {
				if tt.salary >= 0 && err == nil {
					t.Errorf("Expected error for income %f, but got a result", tt.salary)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if math.Abs(result.TotalTax-tt.expectedTotalTaxes) > 0.001 {
				t.Errorf("Expected tax %f, got %f", tt.expectedTotalTaxes, result.TotalTax)
			}

			if math.Abs(result.EffectiveRate-tt.expectedEffectiveRate) > 0.001 {
				t.Errorf("Expected effective rate %f, got %f", tt.expectedEffectiveRate, result.EffectiveRate)
			}
		})
	}
}
