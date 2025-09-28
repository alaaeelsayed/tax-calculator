package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alaaeelsayed/tax-calculator/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	calculator *service.TaxCalculatorService
}

func NewServer(calculator *service.TaxCalculatorService) *Server {
	return &Server{
		calculator: calculator,
	}
}

func (s *Server) SetupRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	r.Get("/taxes/{year}", s.calculateTaxes)

	return r
}

func (s *Server) calculateTaxes(w http.ResponseWriter, r *http.Request) {
	year := chi.URLParam(r, "year")
	if year == "" {
		http.Error(w, "year parameter is required", http.StatusBadRequest)
		return
	}

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

	json.NewEncoder(w).Encode(result)
}
