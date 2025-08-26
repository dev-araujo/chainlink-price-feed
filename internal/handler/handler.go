package handler

import (
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

func (h *PriceHandler) getAllPricesUsd(c *gin.Context) {
	priceData, err := h.service.GetAllPricesUSD()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []gin.H
	for _, p := range priceData {
		response = append(response, gin.H{
			"pair":      p.Pair,
			"price":     p.Price.Text('f', 2),
			"timestamp": p.Timestamp,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *PriceHandler) getAllPricesBrl(c *gin.Context) {
	priceData, err := h.service.GetAllPricesBRL()
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
	asset := strings.ToLower(c.Param("asset"))

	priceData, err := h.service.GetPriceUSD(asset)
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

func (h *PriceHandler) getPriceBrl(c *gin.Context) {
	asset := strings.ToLower(c.Param("asset"))

	priceData, err := h.service.GetPriceBRL(asset)
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
