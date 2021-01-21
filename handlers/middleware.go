package handlers

import (
	"net/http"

	configs "github.com/crowdeco/skeleton/configs"
)

type Middleware struct {
	middlewares []configs.Middleware
}

func NewMiddleware() *Middleware {
	var middlewares []configs.Middleware

	return &Middleware{
		middlewares: middlewares,
	}
}

func (m *Middleware) Add(middleware configs.Middleware) {
	m.middlewares = append(m.middlewares, middleware)
}

func (m *Middleware) Attach(handler http.Handler) http.Handler {
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
