package web

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"fullcycle-go-clima-por-cep/internal/entity"
	"fullcycle-go-clima-por-cep/internal/usecase"
)

type ClimaHandler struct {
	UseCase *usecase.BuscarClimaUseCase
}

func NovoClimaHandler(useCase *usecase.BuscarClimaUseCase) *ClimaHandler {
	return &ClimaHandler{UseCase: useCase}
}

func (h *ClimaHandler) Buscar(w http.ResponseWriter, r *http.Request) {
	input := usecase.BuscarClimaInputDTO{CEP: r.URL.Query().Get("cep")}
	output, err := h.UseCase.Executar(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrCEPInvalido):
			responderErro(w, http.StatusUnprocessableEntity, err.Error())
		case errors.Is(err, entity.ErrCEPNaoEncontrado):
			responderErro(w, http.StatusNotFound, err.Error())
		default:
			log.Printf("erro ao buscar o clima do cep %q: %v", input.CEP, err)
			responderErro(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}

func responderErro(w http.ResponseWriter, status int, mensagem string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	w.Write([]byte(mensagem))
}
