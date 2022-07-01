package todos

import (
	"net/http"

	"github.com/KejawenLab/bima/v3/middlewares"
	"google.golang.org/grpc"
)

type HelloWorld struct {
}

func (a *HelloWorld) Path() string {
	return "/todos/hello/{name}"
}

func (a *HelloWorld) Method() string {
	return http.MethodGet
}

func (a *HelloWorld) SetClient(client *grpc.ClientConn) {}

func (a *HelloWorld) Middlewares() []middlewares.Middleware {
	return []middlewares.Middleware{&Middleware{}}
}

func (a *HelloWorld) Handle(w http.ResponseWriter, r *http.Request, params map[string]string) {
	w.Write([]byte("Hello " + params["name"]))
}
