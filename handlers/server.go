package handlers

import (
	"net/http"

	configs "github.com/crowdeco/skeleton/configs"
)

type (
	Server struct {
		handler http.Handler
	}
)

func NewServer(handler http.Handler) *Server {
	return &Server{
		handler: handler,
	}
}

func (s *Server) Serve() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		configs.Env.User.ID = r.Header.Get(configs.Env.HeaderUserId)
		configs.Env.User.Email = r.Header.Get(configs.Env.HeaderUserEmail)
		configs.Env.User.Role = r.Header.Get(configs.Env.HeaderUserRole)

		s.handler.ServeHTTP(w, r)
	})
}
