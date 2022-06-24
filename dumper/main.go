package main

import (
	"fmt"
	"os"

	dics "github.com/KejawenLab/skeleton/v3/configs"
	"github.com/sarulabs/dingo/v4"
)

func main() {
	err := dingo.GenerateContainer((*dics.Provider)(nil), "generated")
	if err != nil {
		fmt.Println("Error dumping container: ", err.Error())
		os.Exit(1)
	}
}
