package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"fullcycle-go-clima-por-cep/internal/infra/gateway"
	"fullcycle-go-clima-por-cep/internal/usecase"

	"github.com/stretchr/testify/assert"
)

func novoHandlerDeTeste(t *testing.T, statusClima int) *ClimaHandler {
	viaCEP := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/01001000/") {
			w.Write([]byte(`{"cep":"01001-000","localidade":"São Paulo","uf":"SP","estado":"São Paulo"}`))
			return
		}
		w.Write([]byte(`{"erro": "true"}`))
	}))
	t.Cleanup(viaCEP.Close)

	openMeteo := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/search" {
			w.Write([]byte(`{"results":[{"name":"São Paulo","latitude":-23.54,"longitude":-46.63,"admin1":"São Paulo"}]}`))
			return
		}
		w.WriteHeader(statusClima)
		w.Write([]byte(`{"current":{"temperature_2m":28.5}}`))
	}))
	t.Cleanup(openMeteo.Close)

	localizacao := gateway.NovoViaCEPGateway(viaCEP.Client())
	localizacao.URLBase = viaCEP.URL
	clima := gateway.NovoOpenMeteoGateway(openMeteo.Client())
	clima.URLGeocodificacao = openMeteo.URL
	clima.URLPrevisao = openMeteo.URL

	return NovoClimaHandler(usecase.NovoBuscarClimaUseCase(localizacao, clima))
}

func requisitar(h *ClimaHandler, cep string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/?cep="+cep, nil)
	rec := httptest.NewRecorder()
	h.Buscar(rec, req)
	return rec
}

func TestHandlerRetornaTemperaturas(t *testing.T) {
	rec := requisitar(novoHandlerDeTeste(t, http.StatusOK), "01001000")

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	var corpo map[string]float64
	assert.Nil(t, json.Unmarshal(rec.Body.Bytes(), &corpo))
	assert.Equal(t, map[string]float64{"temp_C": 28.5, "temp_F": 83.3, "temp_K": 301.5}, corpo)
}

func TestHandlerComCEPInvalido(t *testing.T) {
	h := novoHandlerDeTeste(t, http.StatusOK)
	for _, cep := range []string{"", "123", "01001-000", "abcdefgh", "010010000"} {
		rec := requisitar(h, cep)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code, "cep %q", cep)
		assert.Equal(t, "invalid zipcode", rec.Body.String())
	}
}

func TestHandlerComCEPNaoEncontrado(t *testing.T) {
	rec := requisitar(novoHandlerDeTeste(t, http.StatusOK), "99999999")

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Equal(t, "can not find zipcode", rec.Body.String())
}

func TestHandlerComFalhaNaOpenMeteo(t *testing.T) {
	rec := requisitar(novoHandlerDeTeste(t, http.StatusBadGateway), "01001000")

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, "internal server error", rec.Body.String())
}
