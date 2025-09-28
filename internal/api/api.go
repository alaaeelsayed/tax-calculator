package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/alaaeelsayed/tax-calculator/internal/service"
)

type Server struct {
	calculator *service.TaxCalculatorService
}

func NewServer(calculator *service.TaxCalculatorService) *Server {
	return &Server{
		calculator: calculator,
	}
}

func (s *Server) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/taxes/", s.calculateTaxes)

	return mux
}

func (s *Server) calculateTaxes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/taxes/")
	parts := strings.Split(strings.Trim(path, "/"), "/")

	if len(parts) == 0 || parts[0] == "" {
		http.Error(w, "year parameter is required", http.StatusBadRequest)
		return
	}

	year := parts[0]

	salaryStr := r.URL.Query().Get("salary")
	if salaryStr == "" {
		http.Error(w, "salary parameter is required", http.StatusBadRequest)
		return
	}

	salary, err := strconv.ParseFloat(salaryStr, 64)
	if err != nil {
		http.Error(w, "invalid salary value", http.StatusBadRequest)
		return
	}

	result, err := s.calculator.CalculateTax(salary, year)
	if err != nil {
		http.Error(w, "Unable to calculate tax. An internal error has occured. Please try again later.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
