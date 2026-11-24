package usecase

import (
	"context"
	"errors"
	"testing"

	"fullcycle-go-clima-por-cep/internal/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type LocalizacaoGatewayMock struct {
	mock.Mock
}

func (m *LocalizacaoGatewayMock) BuscarLocalizacao(ctx context.Context, cep string) (entity.Localizacao, error) {
	args := m.Called(cep)
	return args.Get(0).(entity.Localizacao), args.Error(1)
}

type ClimaGatewayMock struct {
	mock.Mock
}

func (m *ClimaGatewayMock) BuscarTemperaturaCelsius(ctx context.Context, localizacao entity.Localizacao) (float64, error) {
	args := m.Called(localizacao)
	return args.Get(0).(float64), args.Error(1)
}

var saoPaulo = entity.Localizacao{Cidade: "São Paulo", Estado: "São Paulo"}

func TestBuscarClimaComSucesso(t *testing.T) {
	localizacao := &LocalizacaoGatewayMock{}
	clima := &ClimaGatewayMock{}
	localizacao.On("BuscarLocalizacao", "01001000").Return(saoPaulo, nil)
	clima.On("BuscarTemperaturaCelsius", saoPaulo).Return(28.5, nil)

	output, err := NovoBuscarClimaUseCase(localizacao, clima).Executar(context.Background(), BuscarClimaInputDTO{CEP: "01001000"})

	assert.Nil(t, err)
	assert.Equal(t, BuscarClimaOutputDTO{TempC: 28.5, TempF: 83.3, TempK: 301.5}, output)
	localizacao.AssertExpectations(t)
	clima.AssertExpectations(t)
}

func TestBuscarClimaComCEPInvalido(t *testing.T) {
	localizacao := &LocalizacaoGatewayMock{}
	clima := &ClimaGatewayMock{}

	_, err := NovoBuscarClimaUseCase(localizacao, clima).Executar(context.Background(), BuscarClimaInputDTO{CEP: "123"})

	assert.ErrorIs(t, err, entity.ErrCEPInvalido)
	localizacao.AssertNotCalled(t, "BuscarLocalizacao", mock.Anything)
	clima.AssertNotCalled(t, "BuscarTemperaturaCelsius", mock.Anything)
}

func TestBuscarClimaComCEPNaoEncontrado(t *testing.T) {
	localizacao := &LocalizacaoGatewayMock{}
	clima := &ClimaGatewayMock{}
	localizacao.On("BuscarLocalizacao", "99999999").Return(entity.Localizacao{}, entity.ErrCEPNaoEncontrado)

	_, err := NovoBuscarClimaUseCase(localizacao, clima).Executar(context.Background(), BuscarClimaInputDTO{CEP: "99999999"})

	assert.ErrorIs(t, err, entity.ErrCEPNaoEncontrado)
	clima.AssertNotCalled(t, "BuscarTemperaturaCelsius", mock.Anything)
}

func TestBuscarClimaComFalhaNoClima(t *testing.T) {
	localizacao := &LocalizacaoGatewayMock{}
	clima := &ClimaGatewayMock{}
	localizacao.On("BuscarLocalizacao", "01001000").Return(saoPaulo, nil)
	clima.On("BuscarTemperaturaCelsius", saoPaulo).Return(0.0, errors.New("open-meteo fora do ar"))

	_, err := NovoBuscarClimaUseCase(localizacao, clima).Executar(context.Background(), BuscarClimaInputDTO{CEP: "01001000"})

	assert.EqualError(t, err, "open-meteo fora do ar")
}
