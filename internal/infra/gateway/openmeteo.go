package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"fullcycle-go-clima-por-cep/internal/entity"
)

type OpenMeteoGeocodificacao struct {
	Results []struct {
		Name      string  `json:"name"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Admin1    string  `json:"admin1"`
	} `json:"results"`
}

type OpenMeteoPrevisao struct {
	Current struct {
		Temperature2m float64 `json:"temperature_2m"`
	} `json:"current"`
}

type OpenMeteoGateway struct {
	URLGeocodificacao string
	URLPrevisao       string
	Cliente           *http.Client
}

func NovoOpenMeteoGateway(cliente *http.Client) *OpenMeteoGateway {
	return &OpenMeteoGateway{
		URLGeocodificacao: "https://geocoding-api.open-meteo.com/v1",
		URLPrevisao:       "https://api.open-meteo.com/v1",
		Cliente:           cliente,
	}
}

func (g *OpenMeteoGateway) BuscarTemperaturaCelsius(ctx context.Context, localizacao entity.Localizacao) (float64, error) {
	latitude, longitude, err := g.buscarCoordenadas(ctx, localizacao)
	if err != nil {
		return 0, err
	}
	parametros := url.Values{}
	parametros.Set("latitude", strconv.FormatFloat(latitude, 'f', -1, 64))
	parametros.Set("longitude", strconv.FormatFloat(longitude, 'f', -1, 64))
	parametros.Set("current", "temperature_2m")
	var previsao OpenMeteoPrevisao
	err = g.buscarJSON(ctx, g.URLPrevisao+"/forecast?"+parametros.Encode(), &previsao)
	if err != nil {
		return 0, err
	}
	return previsao.Current.Temperature2m, nil
}

func (g *OpenMeteoGateway) buscarCoordenadas(ctx context.Context, localizacao entity.Localizacao) (float64, float64, error) {
	parametros := url.Values{}
	parametros.Set("name", localizacao.Cidade)
	parametros.Set("count", "10")
	parametros.Set("language", "pt")
	parametros.Set("countryCode", "BR")
	var geocodificacao OpenMeteoGeocodificacao
	err := g.buscarJSON(ctx, g.URLGeocodificacao+"/search?"+parametros.Encode(), &geocodificacao)
	if err != nil {
		return 0, 0, err
	}
	if len(geocodificacao.Results) == 0 {
		return 0, 0, fmt.Errorf("open-meteo não encontrou coordenadas para %s/%s", localizacao.Cidade, localizacao.Estado)
	}
	escolhido := geocodificacao.Results[0]
	for _, resultado := range geocodificacao.Results {
		if strings.EqualFold(resultado.Admin1, localizacao.Estado) {
			escolhido = resultado
			break
		}
	}
	return escolhido.Latitude, escolhido.Longitude, nil
}

func (g *OpenMeteoGateway) buscarJSON(ctx context.Context, endereco string, destino any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endereco, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := g.Cliente.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("open-meteo respondeu com status %d: %s", resp.StatusCode, string(body))
	}
	return json.Unmarshal(body, destino)
}
