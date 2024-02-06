package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/modaniru/cards-auth-service/internal/service/auth"
)

type Server struct {
	router       *chi.Mux
	authServices map[string]auth.Auth
}

func NewServer(services ...auth.Auth) *Server {
	r := chi.NewRouter()
	serviceMap := make(map[string]auth.Auth)
	for _, s := range services {
		serviceMap[s.Key()] = s
	}
	server := Server{router: r, authServices: serviceMap}
	server.initRouter()
	return &server
}

func (s *Server) GetRouter() *chi.Mux {
	return s.router
}
