package handlers

import (
	"net/http"

	configs "github.com/crowdeco/skeleton/configs"
)

type (
	Router struct {
		routes []configs.Router
	}
)

func NewRouter() *Router {
	var routes []configs.Router

	return &Router{
		routes: routes,
	}
}

func (r *Router) Add(route configs.Router) {
	r.routes = append(r.routes, route)
}

func (r *Router) Handle(server *http.ServeMux) *http.ServeMux {
	for _, route := range r.routes {
		route.Handle(server)
	}

	return server
}
