package service

import "github.com/alaaeelsayed/tax-calculator/internal/model"

type TaxBracketProvider interface {
	GetTaxBrackets(year string) ([]model.TaxBracket, error)
}
