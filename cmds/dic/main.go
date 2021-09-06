package main

import (
	"fmt"
	"os"

	dics "github.com/KejawenLab/skeleton/configs"
	"github.com/sarulabs/dingo/v4"
)

func main() {
	err := dingo.GenerateContainer((*dics.Provider)(nil), "generated")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
