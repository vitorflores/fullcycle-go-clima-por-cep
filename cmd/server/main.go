package main

import (
	"log"
	"net/http"
	"time"

	"fullcycle-go-clima-por-cep/configs"
	"fullcycle-go-clima-por-cep/internal/infra/gateway"
	"fullcycle-go-clima-por-cep/internal/infra/web"
	"fullcycle-go-clima-por-cep/internal/infra/web/webserver"
	"fullcycle-go-clima-por-cep/internal/usecase"
)

func main() {
	cfg, err := configs.CarregarConfig()
	if err != nil {
		log.Fatalf("não consegui carregar as configurações: %v", err)
	}

	cliente := &http.Client{Timeout: 10 * time.Second}
	buscarClima := usecase.NovoBuscarClimaUseCase(
		gateway.NovoViaCEPGateway(cliente),
		gateway.NovoOpenMeteoGateway(cliente),
	)
	climaHandler := web.NovoClimaHandler(buscarClima)

	ws := webserver.NovoWebServer(":" + cfg.Porta)
	ws.AdicionarHandler("/", climaHandler.Buscar)

	log.Printf("servidor rodando na porta %s", cfg.Porta)
	if err := ws.Iniciar(); err != nil {
		log.Fatal(err)
	}
}
