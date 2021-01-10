package interfaces

import (
	"context"
	"fmt"
	"log"
	"net/http"

	configs "github.com/crowdeco/skeleton/configs"
	handlers "github.com/crowdeco/skeleton/handlers"
	todos "github.com/crowdeco/skeleton/todos"

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
	todos.NewServer().RegisterRest(ctx, mux)

	err := mux.HandlePath("GET", "/health", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		w.Write([]byte("OK"))
	})

	if err != nil {
		panic(err)
	}

	log.Printf("Starting REST Server on :%d", configs.Env.HtppPort)

	http.ListenAndServe(fmt.Sprintf(":%d", configs.Env.HtppPort), handlers.NewServer(mux).Serve())
}
