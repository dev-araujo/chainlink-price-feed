package service

import (
	"fmt"
)

var assetImageURLs = map[string]string{
	"1inch": "https://github.com/dev-araujo/chainlink-price-feed/blob/main/assets/tokens/1inch-logo.png?raw=true",
	"link":  "https://github.com/dev-araujo/chainlink-price-feed/blob/main/assets/tokens/link-logo.png?raw=true",
	"btc":   "https://github.com/dev-araujo/chainlink-price-feed/blob/main/assets/tokens/btc-logo.png?raw=true",
	"eth":   "https://github.com/dev-araujo/chainlink-price-feed/blob/main/assets/tokens/ether-logo.png?raw=true",
	"paxg":  "https://github.com/dev-araujo/chainlink-price-feed/blob/main/assets/tokens/paxg-logo.png?raw=true",
	"stx":   "https://github.com/dev-araujo/chainlink-price-feed/blob/main/assets/tokens/stx-logo.png?raw=true",
	"uni":   "https://github.com/dev-araujo/chainlink-price-feed/blob/main/assets/tokens/uni-logo.png?raw=true",
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
