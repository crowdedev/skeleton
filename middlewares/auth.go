package middlewares

import (
	"net/http"

	configs "github.com/crowdeco/skeleton/configs"
)

type Auth struct {
}

func (a *Auth) Attach(request *http.Request, response http.ResponseWriter) bool {
	configs.Env.User.ID = request.Header.Get(configs.Env.HeaderUserId)
	configs.Env.User.Email = request.Header.Get(configs.Env.HeaderUserEmail)
	configs.Env.User.Role = request.Header.Get(configs.Env.HeaderUserRole)

	return false
}
