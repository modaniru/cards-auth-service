package server

import (
	"github.com/go-chi/chi/v5"
)

type Server struct{
	router *chi.Mux
}

func NewServer() *Server{
	r := chi.NewRouter()
	server := Server{router: r}
	server.initRouter()
	return &server
}

func (s *Server) GetRouter() *chi.Mux{
	return s.router
}