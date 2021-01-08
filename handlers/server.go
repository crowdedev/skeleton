package handlers

import (
	"net/http"

	configs "github.com/crowdeco/todo-service/configs"
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
		configs.Env.User.ID = r.Header.Get("crowde-user-id")
		configs.Env.User.Email = r.Header.Get("crowde-user-email")
		configs.Env.User.Type = r.Header.Get("crowde-user-type")

		s.handler.ServeHTTP(w, r)
	})
}
