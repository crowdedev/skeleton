package handlers

import (
	"net/http"

	configs "github.com/crowdeco/skeleton/configs"
)

type router struct {
	routes []configs.Router
}

func NewRouter() *router {
	var routes []configs.Router

	return &router{
		routes: routes,
	}
}

func (r *router) Add(route configs.Router) {
	r.routes = append(r.routes, route)
}

func (r *router) Handle(server *http.ServeMux) *http.ServeMux {
	for _, route := range r.routes {
		route.Handle(server)
	}

	return server
}
