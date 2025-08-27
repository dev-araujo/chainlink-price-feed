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
		return nil, fmt.Errorf("falha ao buscar taxa BRL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("falha ao buscar taxa BRL, código de status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler o corpo da resposta: %w", err)
	}

	var result ExchangeRateResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("falha ao decodificar resposta: %w", err)
	}

	brlRate, ok := result.Rates["BRL"]
	if !ok {
		return nil, fmt.Errorf("taxa BRL não encontrada na resposta")
	}

	return new(big.Float).SetFloat64(brlRate), nil
}
