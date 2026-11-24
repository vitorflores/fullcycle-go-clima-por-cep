package usecase

import (
	"context"

	"fullcycle-go-clima-por-cep/internal/entity"
)

type BuscarClimaInputDTO struct {
	CEP string
}

type BuscarClimaOutputDTO struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type BuscarClimaUseCase struct {
	LocalizacaoGateway entity.LocalizacaoGateway
	ClimaGateway       entity.ClimaGateway
}

func NovoBuscarClimaUseCase(localizacao entity.LocalizacaoGateway, clima entity.ClimaGateway) *BuscarClimaUseCase {
	return &BuscarClimaUseCase{
		LocalizacaoGateway: localizacao,
		ClimaGateway:       clima,
	}
}

func (u *BuscarClimaUseCase) Executar(ctx context.Context, input BuscarClimaInputDTO) (BuscarClimaOutputDTO, error) {
	if !entity.CEPValido(input.CEP) {
		return BuscarClimaOutputDTO{}, entity.ErrCEPInvalido
	}
	localizacao, err := u.LocalizacaoGateway.BuscarLocalizacao(ctx, input.CEP)
	if err != nil {
		return BuscarClimaOutputDTO{}, err
	}
	celsius, err := u.ClimaGateway.BuscarTemperaturaCelsius(ctx, localizacao)
	if err != nil {
		return BuscarClimaOutputDTO{}, err
	}
	temperatura := entity.NovaTemperatura(celsius)
	return BuscarClimaOutputDTO{
		TempC: temperatura.Celsius,
		TempF: temperatura.Fahrenheit,
		TempK: temperatura.Kelvin,
	}, nil
}
