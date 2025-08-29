package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
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

var currencySymbols = map[string]string{
	"brl": "R$",
	"usd": "$",
}

var templates = template.Must(template.ParseFiles("./web/templates/prices.html"))

func pricesHandler(c *gin.Context) {
	currency := c.Param("currency")
	apiURL := getAPIURL(currency)

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

	currencySymbol, ok := currencySymbols[currency]
	if !ok {
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

	err = templates.ExecuteTemplate(c.Writer, "prices.html", viewModels)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error executing template: %v", err)
	}
}

func getAPIURL(currency string) string {
	baseURL := os.Getenv("API_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	return fmt.Sprintf("%s/api/price/all/%s", baseURL, currency)
}

func main() {
	router := gin.Default()

	router.Static("/styles", "./web/styles")
	router.Static("/assets", "./assets")

	router.GET("/", func(c *gin.Context) {
		c.File("./web/index.html")
	})

	router.GET("/prices/all/:currency", pricesHandler)

	port := os.Getenv("WEB_PORT")
	if port == "" {
		port = "8081"
	}
	serverAddr := fmt.Sprintf(":%s", port)

	log.Printf("Iniciando servidor web na porta %s", port)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Falha ao iniciar o servidor web: %v", err)
	}
}
