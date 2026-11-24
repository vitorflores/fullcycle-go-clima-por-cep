package gateway

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"fullcycle-go-clima-por-cep/internal/entity"

	"github.com/stretchr/testify/assert"
)

func novoViaCEPDeTeste(t *testing.T, status int, corpo string) *ViaCEPGateway {
	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/01001000/json/", r.URL.Path)
		w.WriteHeader(status)
		w.Write([]byte(corpo))
	}))
	t.Cleanup(servidor.Close)
	g := NovoViaCEPGateway(servidor.Client())
	g.URLBase = servidor.URL
	return g
}

func TestViaCEPEncontraCidade(t *testing.T) {
	g := novoViaCEPDeTeste(t, http.StatusOK, `{"cep":"01001-000","localidade":"São Paulo","uf":"SP","estado":"São Paulo"}`)

	localizacao, err := g.BuscarLocalizacao(context.Background(), "01001000")

	assert.Nil(t, err)
	assert.Equal(t, entity.Localizacao{Cidade: "São Paulo", Estado: "São Paulo"}, localizacao)
}

func TestViaCEPNaoEncontraCidade(t *testing.T) {
	for _, corpo := range []string{`{"erro": true}`, `{"erro": "true"}`} {
		g := novoViaCEPDeTeste(t, http.StatusOK, corpo)

		_, err := g.BuscarLocalizacao(context.Background(), "01001000")

		assert.ErrorIs(t, err, entity.ErrCEPNaoEncontrado)
	}
}

func TestViaCEPComStatusInesperado(t *testing.T) {
	g := novoViaCEPDeTeste(t, http.StatusInternalServerError, "")

	_, err := g.BuscarLocalizacao(context.Background(), "01001000")

	assert.EqualError(t, err, "viacep respondeu com status 500")
}
