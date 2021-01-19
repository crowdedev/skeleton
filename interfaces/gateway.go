package interfaces

import (
	"context"
	"net/http"

	grpcs "github.com/crowdeco/skeleton/protos/builds"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func newGateway(ctx context.Context, conn *grpc.ClientConn) (http.Handler, error) {
	mux := runtime.NewServeMux()

	for _, f := range []func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error{
		grpcs.RegisterTodosHandler,
	} {
		if err := f(ctx, mux, conn); err != nil {
			return nil, err
		}
	}
	return mux, nil
}
