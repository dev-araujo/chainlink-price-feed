package service

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"time"
)

const frankfurterAPIURL = "https://api.frankfurter.app/latest?from=USD&to=BRL"

type ExchangeRateResponse struct {
	Rates map[string]float64 `json:"rates"`
}

type ExchangeService struct {
	httpClient *http.Client
}

func NewExchangeService() *ExchangeService {
	return &ExchangeService{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *ExchangeService) GetBRLRate() (*big.Float, error) {
	resp, err := s.httpClient.Get(frankfurterAPIURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch BRL rate: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch BRL rate, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result ExchangeRateResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	brlRate, ok := result.Rates["BRL"]
	if !ok {
		return nil, fmt.Errorf("BRL rate not found in response")
	}

	return new(big.Float).SetFloat64(brlRate), nil
}
