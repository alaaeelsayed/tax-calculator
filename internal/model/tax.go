package model

type TaxBracket struct {
	Min  float64 `json:"min"`
	Max  float64 `json:"max,omitempty"`
	Rate float64 `json:"rate"`
}

type TaxBracketResponse struct {
	TaxBrackets []TaxBracket `json:"tax_brackets"`
}

type TaxCalculationResponse struct {
	TotalTax      float64            `json:"total_tax"`
	EffectiveRate float64            `json:"effective_rate"`
	TaxByBracket  []TaxBracketDetail `json:"tax_by_bracket"`
}

type TaxBracketDetail struct {
	Min           float64 `json:"min"`
	Max           float64 `json:"max,omitempty"`
	Rate          float64 `json:"rate"`
	AmountTaxable float64 `json:"amount_taxable"`
	TaxPayable    float64 `json:"tax_payable"`
}
