

  <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/Gin-0077B5?style=for-the-badge&logo=gin&logoColor=white" alt="Gin">
  <img src="https://img.shields.io/badge/Ethereum-3C3C3D?style=for-the-badge&logo=ethereum&logoColor=white" alt="Go-Ethereum">
  <img src="https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white" alt="Docker">
  <img src="https://img.shields.io/badge/Chainlink-375BD2?style=for-the-badge&logo=chainlink&logoColor=white" alt="Chainlink">

# Chainlink Price Feed com GO

Este projeto é uma API desenvolvida em **Go** que serve como uma ponte entre o mundo da web tradicional e os dados da blockchain Ethereum, utilizando os **Chainlink Data Feeds**.



### O Que São os Chainlink Data Feeds?

A blockchain, por natureza, é um sistema isolado. Ela não tem conhecimento de eventos do mundo real, como o preço atual de uma criptomoeda em dólares. A Chainlink resolve esse problema através de **oráculos**: serviços seguros que coletam informações do mundo real (como dados de preços de várias corretoras), os validam e os publicam de forma confiável e descentralizada na blockchain.

Esses dados são disponibilizados através de **contratos inteligentes** específicos, conhecidos como Data Feeds.

> Para saber mais veja : [Chainlink Data Feeds](https://docs.chain.link/data-feeds)

### Como a Aplicação Funciona?

A API se conecta a um nó da rede Ethereum via RPC e interage diretamente com os contratos inteligentes da Chainlink. O processo envolve:

1.  **Conexão e Instanciação:** Usa a ABI do contrato para criar uma interface em Go.
2.  **Busca de Dados:** Chama funções do contrato (`latestRoundData()`) para obter preços e metadados.
3.  **Serviço via API:** Formata os dados da blockchain e os expõe através de uma API RESTful.

O objetivo é simplificar o acesso a dados on-chain, permitindo que qualquer aplicação consuma preços de criptomoedas de forma segura e sem a complexidade da interação direta com a blockchain.

## 🛠️ Stack

*   [Go](https://golang.org/)
*   [Gin](https://github.com/gin-gonic/gin)
*   [Go-Ethereum](https://github.com/ethereum/go-ethereum)
*   [Docker](https://www.docker.com/)

## 🚀 Executando a aplicação

Siga as instruções abaixo para ter uma cópia do projeto rodando em sua máquina.

### Instalação

1.  Clone o repositório:
    ```sh
    git clone https://github.com/dev-araujo/chainlink-price-feed.git
    cd chainlink-price-feed
    ```

2.  Crie um arquivo `.env` a partir do exemplo e adicione sua URL da Infura:
    ```sh
    cp .env.example .env
    ```
    Edite o arquivo `.env` com suas credenciais:
    
    ```
    RPC_URL="https://mainnet.infura.io/v3/YOUR_INFURA_PROJECT_ID"
    SERVER_PORT="8080"
    GIN_MODE="release"
    ```


   > **💡 A maneira mais fácil e simples de ter acesso a um RPC da Ethereum gratuito é por meio da [Public Node](https://ethereum.publicnode.com/), mas funciona com outras opções como Infura ou Alchemy**
 
---

### Opção 1: Executando com Docker (Recomendado)

Esta é a maneira mais simples de executar a aplicação, pois gerencia todas as dependências para você.

**Pré-requisitos:**
*   [Docker](https://docs.docker.com/get-docker/)

**Comando:**
Para iniciar a aplicação, execute o seguinte comando na raiz do projeto:
```sh
docker-compose up --build
```
A API estará disponível em `http://localhost:8080`.

---

### Opção 2: Executando Localmente

Esta opção é ideal para desenvolvimento e testes diretos no código-fonte.

**Pré-requisitos:**
*   [Go](https://golang.org/doc/install) (versão 1.24.4 ou superior)

**Comando:**
Para iniciar a aplicação localmente, execute o comando abaixo:
```sh
go run ./cmd/api/main.go
```
A API estará disponível em `http://localhost:8080`.

---

## Endpoints da API

A API fornece os seguintes endpoints:

*   `GET /health`: Verifica o status da API.
*   `GET /api/price/:asset/usd`: Retorna o preço do ativo especificado em USD. Substitua `:asset` pelo símbolo do ativo (ex: `btc`, `eth`).
*   `GET /api/price/:asset/brl`: Retorna o preço do ativo especificado em BRL.
*   `GET /api/price/all/usd`: Retorna o preço de todos os ativos suportados em USD.
*   `GET /api/price/all/brl`: Retorna o preço de todos os ativos suportados em BRL.

---

#### Autor 👷

<img src="https://avatars.githubusercontent.com/u/97068163?v=4" width=120 />

[Adriano P Araujo](https://www.linkedin.com/in/araujocode/)