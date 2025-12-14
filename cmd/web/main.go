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

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func pricesHandler(c *gin.Context) {
	currency := c.Param("currency")
	apiURL := getAPIURL(currency)

	resp, err := httpClient.Get(apiURL)
	if err != nil {
		log.Printf("Erro ao conectar na API: %v", err)
		c.String(http.StatusServiceUnavailable, "<div class='error'>Erro ao conectar com o serviço de preços. Tentando novamente...</div>")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Erro ao ler resposta da API: %v", err)
		c.String(http.StatusInternalServerError, "Erro ao ler dados da API")
		return
	}

	contentType := resp.Header.Get("Content-Type")
	if len(body) > 0 && (body[0] == '<' || contentType == "text/html") {
		log.Printf("API retornou HTML inesperado (provavelmente erro 500/502/503 do Render): %s", string(body))
		c.Header("HX-Trigger", "load delay:3s")
		c.String(http.StatusOK, "<div class='loading'>Serviço de preços iniciando... Aguarde.</div>")
		return
	}

	var prices []PriceResponse
	if err := json.Unmarshal(body, &prices); err != nil {
		log.Printf("Erro ao fazer parse do JSON: %v. Body: %s", err, string(body))
		c.String(http.StatusInternalServerError, "Erro ao processar dados de preços")
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

func readyHandler(c *gin.Context) {
	apiURL := getBaseAPIURL() + "/health"

	resp, err := httpClient.Get(apiURL)
	isReady := false

	if err == nil {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			isReady = true
		}
	} else {
		log.Printf("Health check falhou: %v", err)
	}

	if isReady {
		// API está pronta! Retorna o componente original para carregar os preços
		// O atributo hx-get vai disparar imediatamente por causa do hx-trigger="load" implícito ao inserir
		html := `
		<div id="price-list-container" class="price-list" 
			 hx-get="/prices/all/brl" 
			 hx-trigger="load, change from:#currency-select" 
			 hx-swap="innerHTML" 
			 hx-indicator="#loading-indicator">
		</div>`
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, html)
	} else {
		// API ainda não está pronta. Retorna um estado de loading que se auto-atualiza
		// hx-trigger="load delay:2s" faz com que este mesmo endpoint seja chamado novamente em 2s
		// O target deve ser o próprio elemento pai ou substituído corretamente.
		// Vamos retornar um div que substitui o conteúdo atual e tenta de novo.
		html := `
		<div hx-get="/ready" hx-trigger="load delay:2s" hx-swap="outerHTML">
			<div class="skeleton-list" style="opacity: 1; text-align: center; padding: 20px;">
				<p>Iniciando serviços...</p>
				<p><small>(Isso pode levar até 1 minuto no plano gratuito)</small></p>
				<div class="skeleton-card"></div>
				<div class="skeleton-card"></div>
			</div>
		</div>`
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, html)
	}
}

func getAPIURL(currency string) string {
	baseURL := getBaseAPIURL()
	return fmt.Sprintf("%s/api/price/all/%s", baseURL, currency)
}

func getBaseAPIURL() string {
	baseURL := os.Getenv("API_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	return baseURL
}

func main() {
	router := gin.Default()

	router.Static("/styles", "./web/styles")
	router.Static("/assets", "./assets")

	router.GET("/", func(c *gin.Context) {
		c.File("./web/index.html")
	})

	router.GET("/ready", readyHandler)
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
