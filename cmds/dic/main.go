package main

import (
	"fmt"
	"os"

	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"

	dics "github.com/crowdeco/skeleton/configs"
	"github.com/sarulabs/dingo/v4"
)

func main() {
	err := dingo.GenerateContainer((*dics.Provider)(nil), "generated")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
