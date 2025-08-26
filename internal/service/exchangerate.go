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
		return nil, fmt.Errorf("falha ao buscar taxa BRL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("falha ao buscar taxa BRL, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler o corpo da resposta: %w", err)
	}

	var result ExchangeRateResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("não conseguiu desmantelar a resposta: %w", err)
	}

	brlRate, ok := result.Rates["BRL"]
	if !ok {
		return nil, fmt.Errorf("taxa BRL não encontrada na resposta")
	}

	return new(big.Float).SetFloat64(brlRate), nil
}
