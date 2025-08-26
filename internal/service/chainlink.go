package service

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/dev-araujo/chainlink-price-feed/contracts"
	"github.com/dev-araujo/chainlink-price-feed/internal/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type PriceData struct {
	Pair      string
	Price     *big.Float
	Timestamp int64
}

type ChainlinkService struct {
	client        *ethclient.Client
	contractAddrs map[string]string
}

func NewChainlinkService(client *ethclient.Client) *ChainlinkService {
	return &ChainlinkService{client: client, contractAddrs: config.ContractAddresses}
}

func (s *ChainlinkService) GetPriceUSD(asset string) (*PriceData, error) {
	return s.fetchPriceFromChainlink(asset)
}

func (s *ChainlinkService) GetPriceBRL(asset string) (*PriceData, error) {
	assetPriceData, err := s.fetchPriceFromChainlink(asset)
	if err != nil {
		return nil, err
	}

	brlRateData := &PriceData{
		Price: new(big.Float).SetFloat64(5.3), // TODO
	}

	priceInBRL := new(big.Float).Mul(assetPriceData.Price, brlRateData.Price)

	return &PriceData{
		Pair:      fmt.Sprintf("%s/BRL", strings.ToUpper(asset)),
		Price:     priceInBRL,
		Timestamp: assetPriceData.Timestamp,
	}, nil
}

func (s *ChainlinkService) fetchPriceFromChainlink(asset string) (*PriceData, error) {
	addressHex, ok := s.contractAddrs[asset]
	if !ok {
		return nil, fmt.Errorf("ativo '%s' não suportado", asset)
	}

	address := common.HexToAddress(addressHex)
	priceFeed, err := contracts.NewAggregatorV3Interface(address, s.client)
	if err != nil {
		return nil, fmt.Errorf("falha ao instanciar contrato para %s: %w", asset, err)
	}

	decimals, err := priceFeed.Decimals(nil)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar decimais para %s: %w", asset, err)
	}

	latestRoundData, err := priceFeed.LatestRoundData(nil)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar dados para %s: %w", asset, err)
	}

	price := new(big.Float).SetInt(latestRoundData.Answer)
	price.Quo(price, new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)))

	pairName := fmt.Sprintf("%s/USD", strings.ToUpper(asset))
	if asset == "brl" {
		pairName = "BRL/USD"
	}

	return &PriceData{
		Pair:      pairName,
		Price:     price,
		Timestamp: latestRoundData.UpdatedAt.Int64(),
	}, nil
}

func (s *ChainlinkService) GetAllPricesUSD() ([]*PriceData, error) {
	var prices []*PriceData
	for asset := range s.contractAddrs {
		priceData, err := s.fetchPriceFromChainlink(asset)
		if err != nil {
			return nil, fmt.Errorf("falha ao buscar preço para %s: %w", asset, err)
		}
		prices = append(prices, priceData)
	}
	return prices, nil
}

func (s *ChainlinkService) GetAllPricesBRL() ([]*PriceData, error) {
	var prices []*PriceData
	for asset := range s.contractAddrs {
		priceData, err := s.GetPriceBRL(asset)
		if err != nil {
			return nil, fmt.Errorf("falha ao buscar preço para %s: %w", asset, err)
		}
		prices = append(prices, priceData)
	}
	return prices, nil
}
