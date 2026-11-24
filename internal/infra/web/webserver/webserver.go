package webserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]http.HandlerFunc
	WebServerPort string
}

func NovoWebServer(porta string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]http.HandlerFunc),
		WebServerPort: porta,
	}
}

func (s *WebServer) AdicionarHandler(caminho string, handler http.HandlerFunc) {
	s.Handlers[caminho] = handler
}

func (s *WebServer) Iniciar() error {
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	for caminho, handler := range s.Handlers {
		s.Router.Get(caminho, handler)
	}
	return http.ListenAndServe(s.WebServerPort, s.Router)
}
