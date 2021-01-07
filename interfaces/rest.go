package interfaces

import (
	"context"
	"fmt"
	"net/http"

	configs "github.com/crowdeco/todo-service/configs"
	todos "github.com/crowdeco/todo-service/todos"

	runtime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type (
	Rest struct{}
)

func NewRest() configs.Application {
	return &Rest{}
}

func (g *Rest) Run() {
	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	todos.RegisterRestServer(ctx, mux)

	err := mux.HandlePath("GET", "/health", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		w.Write([]byte("OK"))
	})

	if err != nil {
		panic(err)
	}

	http.ListenAndServe(fmt.Sprintf(":%d", configs.Env.HtppPort), mux)
}
