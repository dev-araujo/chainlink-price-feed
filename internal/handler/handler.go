package handler

import (
	"context"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/dev-araujo/chainlink-price-feed/internal/service"
	"github.com/gin-gonic/gin"
)

type PriceResponse struct {
	Pair      string `json:"pair"`
	Price     string `json:"price"`
	Timestamp int64  `json:"timestamp"`
	ImageURL  string `json:"imageUrl"`
}

type PriceHandler struct {
	chainlinkService *service.ChainlinkService
	assetService     *service.AssetService
}

func NewPriceHandler(cs *service.ChainlinkService, as *service.AssetService) *PriceHandler {
	return &PriceHandler{
		chainlinkService: cs,
		assetService:     as,
	}
}

func (h *PriceHandler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/price")
	{
		api.GET("/:asset/usd", h.getPriceUsd)
		api.GET("/:asset/brl", h.getPriceBrl)
		api.GET("/all/usd", h.getAllPricesUsd)
		api.GET("/all/brl", h.getAllPricesBrl)
	}
}

func (h *PriceHandler) getPrice(c *gin.Context, getPriceFunc func(ctx context.Context, asset string) (*service.PriceData, error)) {
	asset := strings.ToLower(c.Param("asset"))

	priceData, err := getPriceFunc(c.Request.Context(), asset)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": err.Error()})
		return
	}

	imageURL, err := h.assetService.GetAssetImageURL(asset)
	if err != nil {
		log.Printf("não foi possível obter a URL da imagem para o ativo %s: %v", asset, err)
	}

	c.JSON(http.StatusOK, PriceResponse{
		Pair:      priceData.Pair,
		Price:     priceData.Price.Text('f', 2),
		Timestamp: priceData.Timestamp,
		ImageURL:  imageURL,
	})
}

func (h *PriceHandler) getPriceUsd(c *gin.Context) {
	h.getPrice(c, h.chainlinkService.GetPriceUSD)
}

func (h *PriceHandler) getPriceBrl(c *gin.Context) {
	h.getPrice(c, h.chainlinkService.GetPriceBRL)
}

func (h *PriceHandler) getAllPricesUsd(c *gin.Context) {
	priceData, err := h.chainlinkService.GetAllPricesUSD()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}

	h.buildAndSendAllPricesResponse(c, priceData)
}

func (h *PriceHandler) getAllPricesBrl(c *gin.Context) {
	priceData, err := h.chainlinkService.GetAllPricesBRL()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}

	h.buildAndSendAllPricesResponse(c, priceData)
}

func (h *PriceHandler) buildAndSendAllPricesResponse(c *gin.Context, priceData []*service.PriceData) {
	responses := make([]PriceResponse, len(priceData))
	var wg sync.WaitGroup

	for i, p := range priceData {
		wg.Add(1)
		go func(index int, data *service.PriceData) {
			defer wg.Done()

			assetSymbol := strings.ToLower(strings.Split(data.Pair, "/")[0])
			imageURL, err := h.assetService.GetAssetImageURL(assetSymbol)
			if err != nil {
				log.Printf("não foi possível obter a URL da imagem para o ativo %s: %v", assetSymbol, err)
			}

			responses[index] = PriceResponse{
				Pair:      data.Pair,
				Price:     data.Price.Text('f', 2),
				Timestamp: data.Timestamp,
				ImageURL:  imageURL,
			}
		}(i, p)
	}

	wg.Wait()
	c.JSON(http.StatusOK, responses)
}
