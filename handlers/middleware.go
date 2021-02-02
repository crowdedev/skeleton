package handlers

import (
	"net/http"

	configs "github.com/crowdeco/skeleton/configs"
)

type middleware struct {
	middlewares []configs.Middleware
}

func NewMiddleware() *middleware {
	var middlewares []configs.Middleware

	return &middleware{
		middlewares: middlewares,
	}
}

func (m *middleware) Add(middleware configs.Middleware) {
	m.middlewares = append(m.middlewares, middleware)
}

func (m *middleware) Attach(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		for _, middleware := range m.middlewares {
			stop := middleware.Attach(request, response)
			if stop {
				return
			}
		}

		handler.ServeHTTP(response, request)
	})
}
