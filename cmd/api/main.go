package main

import (
	"fmt"
	"log"

	"github.com/dev-araujo/chainlink-price-feed/internal/config"
	"github.com/dev-araujo/chainlink-price-feed/internal/handler"
	"github.com/dev-araujo/chainlink-price-feed/internal/service"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	if cfg.RpcURL == "" {
		log.Fatal("SEPOLIA_RPC_URL não pode ser vazia.")
	}
	if cfg.ServerPort == "" {
		cfg.ServerPort = "8080"
	}

	client, err := ethclient.Dial(cfg.RpcURL)
	if err != nil {
		log.Fatalf("Falha ao conectar ao nó da Sepolia: %v", err)
	}
	log.Println("Conectado com sucesso à rede Sepolia!")

	exchangeService := service.NewExchangeService()

	chainlinkService := service.NewChainlinkService(client, exchangeService)

	priceHandler := handler.NewPriceHandler(chainlinkService)

	router := gin.Default()
	router.Use(cors.Default())

	priceHandler.RegisterRoutes(router)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Iniciando servidor na porta %s", cfg.ServerPort)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}
