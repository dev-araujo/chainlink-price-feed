package service

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/dev-araujo/chainlink-price-feed/contracts"
	"github.com/dev-araujo/chainlink-price-feed/internal/config"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/sync/errgroup"
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

func (s *ChainlinkService) GetPriceUSD(ctx context.Context, asset string) (*PriceData, error) {
	return s.fetchPriceFromChainlink(ctx, asset)
}

func (s *ChainlinkService) GetPriceBRL(ctx context.Context, asset string) (*PriceData, error) {
	assetPriceData, err := s.fetchPriceFromChainlink(ctx, asset)
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

func (s *ChainlinkService) fetchPriceFromChainlink(ctx context.Context, asset string) (*PriceData, error) {
	addressHex, ok := s.contractAddrs[asset]
	if !ok {
		return nil, fmt.Errorf("ativo '%s' não suportado", asset)
	}

	address := common.HexToAddress(addressHex)
	priceFeed, err := contracts.NewAggregatorV3Interface(address, s.client)
	if err != nil {
		return nil, fmt.Errorf("falha ao instanciar contrato para %s: %w", asset, err)
	}

	callOpts := &bind.CallOpts{Context: ctx}

	decimals, err := priceFeed.Decimals(callOpts)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar decimais para %s: %w", asset, err)
	}

	latestRoundData, err := priceFeed.LatestRoundData(callOpts)
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

func (s *ChainlinkService) fetchAllPrices(priceFetcher func(ctx context.Context, asset string) (*PriceData, error)) ([]*PriceData, error) {
	prices := make([]*PriceData, 0, len(s.contractAddrs))
	var mu sync.Mutex
	g, ctx := errgroup.WithContext(context.Background())

	for asset := range s.contractAddrs {
		asset := asset
		g.Go(func() error {
			priceData, err := priceFetcher(ctx, asset)
			if err != nil {
				return fmt.Errorf("falha ao buscar preço para %s: %w", asset, err)
			}
			mu.Lock()
			prices = append(prices, priceData)
			mu.Unlock()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return prices, nil
}

func (s *ChainlinkService) GetAllPricesUSD() ([]*PriceData, error) {
	return s.fetchAllPrices(s.GetPriceUSD)
}

func (s *ChainlinkService) GetAllPricesBRL() ([]*PriceData, error) {
	return s.fetchAllPrices(s.GetPriceBRL)
}
