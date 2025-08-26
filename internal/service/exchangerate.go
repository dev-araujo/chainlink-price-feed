package service

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
)

type ExchangeRateResponse struct {
	Rates map[string]float64 `json:"rates"`
}

type ExchangeService struct{}

func NewExchangeService() *ExchangeService {
	return &ExchangeService{}
}

func (s *ExchangeService) GetBRLRate() (*big.Float, error) {
	resp, err := http.Get("https://api.frankfurter.app/latest?from=USD&to=BRL")
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
