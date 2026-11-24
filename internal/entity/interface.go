package entity

import "context"

type Localizacao struct {
	Cidade string
	Estado string
}

type LocalizacaoGateway interface {
	BuscarLocalizacao(ctx context.Context, cep string) (Localizacao, error)
}

type ClimaGateway interface {
	BuscarTemperaturaCelsius(ctx context.Context, localizacao Localizacao) (float64, error)
}
