package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/dev-araujo/chainlink-price-feed/internal/service"
	"github.com/gin-gonic/gin"
)

type PriceHandler struct {
	service *service.ChainlinkService
}

func NewPriceHandler(s *service.ChainlinkService) *PriceHandler {
	return &PriceHandler{service: s}
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

func (h *PriceHandler) getPrice(c *gin.Context, priceFetcher func(ctx context.Context, asset string) (*service.PriceData, error)) {
	asset := strings.ToLower(c.Param("asset"))

	priceData, err := priceFetcher(c.Request.Context(), asset)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pair":      priceData.Pair,
		"price":     priceData.Price.Text('f', 2),
		"timestamp": priceData.Timestamp,
	})
}

func (h *PriceHandler) getAllPrices(c *gin.Context, priceFetcher func() ([]*service.PriceData, error)) {
	priceData, err := priceFetcher()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]gin.H, 0, len(priceData))
	for _, p := range priceData {
		response = append(response, gin.H{
			"pair":      p.Pair,
			"price":     p.Price.Text('f', 2),
			"timestamp": p.Timestamp,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *PriceHandler) getPriceUsd(c *gin.Context) {
	h.getPrice(c, h.service.GetPriceUSD)
}

func (h *PriceHandler) getPriceBrl(c *gin.Context) {
	h.getPrice(c, h.service.GetPriceBRL)
}

func (h *PriceHandler) getAllPricesUsd(c *gin.Context) {
	h.getAllPrices(c, h.service.GetAllPricesUSD)
}

func (h *PriceHandler) getAllPricesBrl(c *gin.Context) {
	h.getAllPrices(c, h.service.GetAllPricesBRL)
}
