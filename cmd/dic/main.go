package main

import (
	"fmt"
	"os"

	dics "github.com/crowdeco/skeleton/dics"
	"github.com/sarulabs/dingo/v4"
)

func main() {
	err := dingo.GenerateContainer((*dics.Provider)(nil), "dics/generated")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
