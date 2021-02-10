package middlewares

import (
	"net/http"
	"strconv"

	configs "github.com/crowdeco/skeleton/configs"
)

type Auth struct {
	Env *configs.Env
}

func (a *Auth) Attach(request *http.Request, response http.ResponseWriter) bool {
	a.Env.User.Id = request.Header.Get(a.Env.HeaderUserId)
	a.Env.User.Email = request.Header.Get(a.Env.HeaderUserEmail)
	a.Env.User.Role, _ = strconv.Atoi(request.Header.Get(a.Env.HeaderUserRole))

	if a.Env.User.Role == 0 || a.Env.User.Role > a.Env.MaximumRole {
		http.Error(response, "Unauthorization", http.StatusUnauthorized)

		return true
	}

	return false
}

func (a *Auth) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
