package client

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/alaaeelsayed/tax-calculator/internal/model"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) GetTaxBrackets(year string) ([]model.TaxBracket, error) {
	url := fmt.Sprintf("%s/tax-calculator/tax-year/%s", c.baseURL, year)

	resp, err := c.retryHTTPRequest(func() (*http.Response, error) {
		return c.httpClient.Get(url)
	}, "fetch tax brackets")

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response model.TaxBracketResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.TaxBrackets, nil
}

func (c *Client) retryHTTPRequest(requestFunc func() (*http.Response, error), operation string) (*http.Response, error) {
	const maxRetries = 3
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		resp, err := requestFunc()
		if err != nil {
			lastErr = err
			if attempt < maxRetries {
				c.backoff(attempt)
				continue
			}
			return nil, fmt.Errorf("failed to %s after %d attempts: %w", operation, maxRetries, err)
		}

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("API returned status %d", resp.StatusCode)
			if attempt < maxRetries && c.shouldRetryStatus(resp.StatusCode) {
				resp.Body.Close()
				c.backoff(attempt)
				continue
			}
			return resp, lastErr
		}

		return resp, nil
	}

	return nil, lastErr
}

func (c *Client) backoff(attempt int) {
	backoffDuration := time.Duration(math.Pow(2, float64(attempt-1))) * time.Second
	time.Sleep(backoffDuration)
}

func (c *Client) shouldRetryStatus(statusCode int) bool {
	return statusCode >= 500 || statusCode == 429
}
