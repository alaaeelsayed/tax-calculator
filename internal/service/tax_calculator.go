package service

import (
	"fmt"
	"math"

	"github.com/alaaeelsayed/tax-calculator/internal/model"
)

type TaxCalculatorService struct {
	taxBracketProvider TaxBracketProvider
}

func NewTaxCalculatorService(taxBracketProvider TaxBracketProvider) *TaxCalculatorService {
	return &TaxCalculatorService{taxBracketProvider: taxBracketProvider}
}

func (s *TaxCalculatorService) CalculateTax(income float64, year string) (*model.TaxCalculationResponse, error) {
	if income < 0 {
		return nil, fmt.Errorf("income cannot be negative")
	}

	brackets, err := s.taxBracketProvider.GetTaxBrackets(year)
	if err != nil {
		return nil, err
	}

	var totalTax float64
	var taxByBracket []model.TaxBracketDetail

	for _, bracket := range brackets {
		var taxableIncome float64
		var tax float64

		if income <= bracket.Min {
			break
		}

		if bracket.Max > 0 && income > bracket.Max {
			taxableIncome = bracket.Max - bracket.Min
		} else {
			taxableIncome = income - bracket.Min // No upper bound or income in range
		}

		tax = taxableIncome * bracket.Rate
		totalTax += tax

		detail := model.TaxBracketDetail{
			Min:           bracket.Min,
			Max:           bracket.Max,
			Rate:          bracket.Rate,
			AmountTaxable: taxableIncome,
			TaxPayable:    tax,
		}
		taxByBracket = append(taxByBracket, detail)
	}

	effectiveRate := 0.0
	if income > 0 {
		effectiveRate = totalTax / income
	}

	return &model.TaxCalculationResponse{
		TotalTax:      math.Round(totalTax*100) / 100,
		EffectiveRate: effectiveRate,
		TaxByBracket:  taxByBracket,
	}, nil
}
