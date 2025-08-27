package service

import (
	"fmt"
)

var assetImageURLs = map[string]string{
	"1inch": "https://cryptologos.cc/logos/1inch-1inch-logo.png?v=040",
	"link":  "https://cryptologos.cc/logos/chainlink-link-logo.png?v=040",
	"btc":   "https://cryptologos.cc/logos/bitcoin-btc-logo.png?v=040",
	"eth":   "https://cryptologos.cc/logos/ethereum-eth-logo.png?v=040",
	"paxg":  "https://cryptologos.cc/logos/pax-gold-paxg-logo.png?v=040",
	"stx":   "https://cryptologos.cc/logos/stacks-stx-logo.png?v=040",
	"uni":   "https://cryptologos.cc/logos/uniswap-uni-logo.png?v=040",
}

type AssetService struct{}

func NewAssetService() *AssetService {
	return &AssetService{}
}

func (s *AssetService) GetAssetImageURL(asset string) (string, error) {
	url, found := assetImageURLs[asset]
	if !found {
		return "", fmt.Errorf("imagem para o ativo '%s' não encontrada", asset)
	}
	return url, nil
}
