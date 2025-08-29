package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type PriceResponse struct {
	Pair      string `json:"pair"`
	Price     string `json:"price"`
	Timestamp int64  `json:"timestamp"`
	ImageURL  string `json:"imageUrl"`
}

type PriceViewModel struct {
	Pair           string
	Price          string
	ImageURL       string
	LastUpdate     string
	CurrencySymbol string
}

func pricesHandler(c *gin.Context) {
	currency := c.Param("currency")
	apiURL := fmt.Sprintf("http://localhost:8080/api/price/all/%s", currency)

	resp, err := http.Get(apiURL)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error fetching prices: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error reading prices: %v", err)
		return
	}

	var prices []PriceResponse
	if err := json.Unmarshal(body, &prices); err != nil {
		c.String(http.StatusInternalServerError, "Error parsing prices: %v", err)
		return
	}

	var currencySymbol string
	if currency == "brl" {
		currencySymbol = "R$"
	} else {
		currencySymbol = "$"
	}

	viewModels := make([]PriceViewModel, len(prices))
	for i, p := range prices {
		viewModels[i] = PriceViewModel{
			Pair:           p.Pair,
			Price:          p.Price,
			ImageURL:       p.ImageURL,
			LastUpdate:     time.Unix(p.Timestamp, 0).Format("15:04:05"),
			CurrencySymbol: currencySymbol,
		}
	}

	tmpl, err := template.ParseFiles("./web/templates/prices.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error parsing template: %v", err)
		return
	}

	tmpl.Execute(c.Writer, viewModels)
}

func main() {
	router := gin.Default()

	router.Static("/styles", "./web/styles")
	router.Static("/assets", "./assets")

	router.GET("/", func(c *gin.Context) {
		c.File("./web/index.html")
	})

	router.GET("/prices/all/:currency", pricesHandler)

	serverAddr := fmt.Sprintf(":%s", "8081")
	log.Printf("Iniciando servidor web na porta %s", "8081")
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Falha ao iniciar o servidor web: %v", err)
	}
}
