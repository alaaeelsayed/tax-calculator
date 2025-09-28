package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alaaeelsayed/tax-calculator/internal/model"
	"github.com/alaaeelsayed/tax-calculator/internal/service"
)

type mockTaxBracketProvider struct {
	brackets []model.TaxBracket
	err      error
}

func (m *mockTaxBracketProvider) GetTaxBrackets(year string) ([]model.TaxBracket, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.brackets, nil
}

func TestCalculateTaxesEndpoint(t *testing.T) {
	testBrackets := []model.TaxBracket{
		{Min: 0, Max: 50197, Rate: 0.15},
		{Min: 50197, Max: 100392, Rate: 0.205},
		{Min: 100392, Max: 155625, Rate: 0.26},
		{Min: 155625, Max: 221708, Rate: 0.29},
		{Min: 221708, Rate: 0.33},
	}

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedError  string
		salary         float64
		expectedTotal  float64
	}{
		{
			name:           "Valid request 50k salary",
			path:           "/taxes/2022?salary=50000",
			expectedStatus: http.StatusOK,
			salary:         50000,
			expectedTotal:  7500.00,
		},
		{
			name:           "Valid request 100k salary",
			path:           "/taxes/2022?salary=100000",
			expectedStatus: http.StatusOK,
			salary:         100000,
			expectedTotal:  17739.17,
		},
		{
			name:           "Missing salary parameter",
			path:           "/taxes/2022",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "salary parameter is required",
		},
		{
			name:           "Invalid salary format",
			path:           "/taxes/2022?salary=invalid",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid salary value",
		},
		{
			name:           "Missing year parameter",
			path:           "/taxes/",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "year parameter is required",
		},
		{
			name:           "Zero salary",
			path:           "/taxes/2022?salary=0",
			expectedStatus: http.StatusOK,
			salary:         0,
			expectedTotal:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProvider := &mockTaxBracketProvider{
				brackets: testBrackets,
			}
			calculator := service.NewTaxCalculatorService(mockProvider)
			server := NewServer(calculator)

			req, err := http.NewRequest("GET", tt.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			server.SetupRoutes().ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				var response model.TaxCalculationResponse
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if response.TotalTax != tt.expectedTotal {
					t.Errorf("Expected total tax %f, got %f", tt.expectedTotal, response.TotalTax)
				}

				if tt.salary > 0 && response.EffectiveRate <= 0 {
					t.Errorf("Expected positive effective rate for non-zero salary")
				}

				if len(response.TaxByBracket) == 0 && tt.salary > 0 {
					t.Errorf("Expected tax breakdown for non-zero salary")
				}
			} else if tt.expectedError != "" {
				if !strings.Contains(rr.Body.String(), tt.expectedError) {
					t.Errorf("Expected error message '%s', got '%s'", tt.expectedError, rr.Body.String())
				}
			}
		})
	}
}

func TestCalculateTaxesEndpoint_ExternalAPIError(t *testing.T) {
	mockProvider := &mockTaxBracketProvider{
		err: fmt.Errorf("external API error"),
	}
	calculator := service.NewTaxCalculatorService(mockProvider)
	server := NewServer(calculator)

	req, err := http.NewRequest("GET", "/taxes/2022?salary=50000", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	server.SetupRoutes().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	expectedError := "Unable to calculate tax. An internal error has occured. Please try again later."
	if !strings.Contains(rr.Body.String(), expectedError) {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, rr.Body.String())
	}
}

func TestCalculateTaxesEndpoint_ContentType(t *testing.T) {
	testBrackets := []model.TaxBracket{
		{Min: 0, Max: 50197, Rate: 0.15},
	}

	mockProvider := &mockTaxBracketProvider{
		brackets: testBrackets,
	}
	calculator := service.NewTaxCalculatorService(mockProvider)
	server := NewServer(calculator)

	req, err := http.NewRequest("GET", "/taxes/2022?salary=10000", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	server.SetupRoutes().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}
}
