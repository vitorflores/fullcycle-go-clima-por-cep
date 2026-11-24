package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"fullcycle-go-clima-por-cep/internal/entity"
)

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Unidade     string `json:"unidade"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
}

type ViaCEPGateway struct {
	URLBase string
	Cliente *http.Client
}

func NovoViaCEPGateway(cliente *http.Client) *ViaCEPGateway {
	return &ViaCEPGateway{
		URLBase: "https://viacep.com.br/ws",
		Cliente: cliente,
	}
}

func (g *ViaCEPGateway) BuscarLocalizacao(ctx context.Context, cep string) (entity.Localizacao, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s/json/", g.URLBase, cep), nil)
	if err != nil {
		return entity.Localizacao{}, err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := g.Cliente.Do(req)
	if err != nil {
		return entity.Localizacao{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusNotFound {
		return entity.Localizacao{}, entity.ErrCEPNaoEncontrado
	}
	if resp.StatusCode != http.StatusOK {
		return entity.Localizacao{}, fmt.Errorf("viacep respondeu com status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return entity.Localizacao{}, err
	}
	var c ViaCEP
	err = json.Unmarshal(body, &c)
	if err != nil {
		return entity.Localizacao{}, err
	}
	if c.Localidade == "" {
		return entity.Localizacao{}, entity.ErrCEPNaoEncontrado
	}
	return entity.Localizacao{Cidade: c.Localidade, Estado: c.Estado}, nil
}
