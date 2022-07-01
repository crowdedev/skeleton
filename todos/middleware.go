package todos

import (
	"net/http"
)

type Middleware struct {
}

func (a *Middleware) Attach(_ *http.Request, response http.ResponseWriter) bool {
	response.Header().Add("X-Middleware", "My first middleware")

	return false
}

func (a *Middleware) Priority() int {
	return 0
}
