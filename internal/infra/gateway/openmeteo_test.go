package gateway

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"fullcycle-go-clima-por-cep/internal/entity"

	"github.com/stretchr/testify/assert"
)

const geocodificacaoBomJesus = `{"results":[
	{"name":"Bom Jesus","latitude":-9.07,"longitude":-44.35,"admin1":"Piauí"},
	{"name":"Bom Jesus","latitude":-28.66,"longitude":-50.41,"admin1":"Rio Grande do Sul"}
]}`

func novoOpenMeteoDeTeste(t *testing.T, statusPrevisao int, corpoGeocodificacao string) (*OpenMeteoGateway, *string) {
	coordenadasRecebidas := new(string)
	geocodificacao := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/search", r.URL.Path)
		assert.Equal(t, "Bom Jesus", r.URL.Query().Get("name"))
		assert.Equal(t, "BR", r.URL.Query().Get("countryCode"))
		w.Write([]byte(corpoGeocodificacao))
	}))
	t.Cleanup(geocodificacao.Close)
	previsao := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/forecast", r.URL.Path)
		assert.Equal(t, "temperature_2m", r.URL.Query().Get("current"))
		*coordenadasRecebidas = r.URL.Query().Get("latitude") + "," + r.URL.Query().Get("longitude")
		w.WriteHeader(statusPrevisao)
		w.Write([]byte(`{"current":{"time":"2026-11-24T02:15","temperature_2m":28.5}}`))
	}))
	t.Cleanup(previsao.Close)
	g := NovoOpenMeteoGateway(http.DefaultClient)
	g.URLGeocodificacao = geocodificacao.URL
	g.URLPrevisao = previsao.URL
	return g, coordenadasRecebidas
}

func TestOpenMeteoRetornaTemperaturaDaCidadeDoEstadoCerto(t *testing.T) {
	g, coordenadas := novoOpenMeteoDeTeste(t, http.StatusOK, geocodificacaoBomJesus)

	celsius, err := g.BuscarTemperaturaCelsius(context.Background(), entity.Localizacao{Cidade: "Bom Jesus", Estado: "Rio Grande do Sul"})

	assert.Nil(t, err)
	assert.Equal(t, 28.5, celsius)
	assert.Equal(t, "-28.66,-50.41", *coordenadas)
}

func TestOpenMeteoUsaPrimeiroResultadoQuandoEstadoNaoBate(t *testing.T) {
	g, coordenadas := novoOpenMeteoDeTeste(t, http.StatusOK, geocodificacaoBomJesus)

	_, err := g.BuscarTemperaturaCelsius(context.Background(), entity.Localizacao{Cidade: "Bom Jesus", Estado: ""})

	assert.Nil(t, err)
	assert.Equal(t, "-9.07,-44.35", *coordenadas)
}

func TestOpenMeteoSemResultadoNaGeocodificacao(t *testing.T) {
	g, _ := novoOpenMeteoDeTeste(t, http.StatusOK, `{"generationtime_ms":0.1}`)

	_, err := g.BuscarTemperaturaCelsius(context.Background(), entity.Localizacao{Cidade: "Bom Jesus", Estado: "Piauí"})

	assert.ErrorContains(t, err, "não encontrou coordenadas")
}

func TestOpenMeteoComErroNaPrevisao(t *testing.T) {
	g, _ := novoOpenMeteoDeTeste(t, http.StatusBadRequest, geocodificacaoBomJesus)

	_, err := g.BuscarTemperaturaCelsius(context.Background(), entity.Localizacao{Cidade: "Bom Jesus", Estado: "Piauí"})

	assert.ErrorContains(t, err, "status 400")
}
