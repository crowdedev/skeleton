package main

import (
	"strings"

	"github.com/crowdeco/skeleton/configs"
	dic "github.com/crowdeco/skeleton/generated/dic"
	"github.com/jinzhu/copier"
	"github.com/vito/go-interact/interact"
)

func main() {
	container, _ := dic.NewContainer()
	util, err := container.SafeGetCoreUtilCli()
	if err != nil {
		panic(err)
	}

	util.Println(`
    _________                         .___         ___________          .__          __
    \_   ___ \_______  ______  _  ____| _/____    /   _____/  | __ ____ |  |   _____/  |_  ____   ____
    /    \  \/\_  __ \/  _ \ \/ \/ / __ |/ __ \   \_____  \|  |/ // __ \|  | _/ __ \   __\/  _ \ /    \
    \     \____|  | \(  <_> )     / /_/ \  ___/   /        \    <\  ___/|  |_\  ___/|  | (  <_> )   |  \
     \______  /|__|   \____/ \/\_/\____ |\___  > /_______  /__|_ \\___  >____/\___  >__|  \____/|___|  /
            \/                         \/    \/          \/     \/    \/          \/                 \/
`)

	generator, err := container.SafeGetCoreModuleGenerator()
	if err != nil {
		panic(err)
	}

	module, err := container.SafeGetCoreTemplateModule()
	if err != nil {
		panic(err)
	}

	field, err := container.SafeGetCoreTemplateField()
	if err != nil {
		panic(err)
	}

	err = interact.NewInteraction("Masukkan Nama Table?").Resolve(&module.Name)
	if err != nil {
		panic(err)
	}

	index := 2
	more := true
	for more {
		err = interact.NewInteraction("Tambah Kolom?").Resolve(&more)
		if err != nil {
			panic(err)
		}

		if more {
			addColumn(field)

			column := configs.FieldTemplate{}

			copier.Copy(&column, field)

			column.Index = index
			column.Name = strings.Title(column.Name)
			module.Fields = append(module.Fields, &column)

			field.Name = ""
			field.Type = ""

			index++
		}
	}

	generator.Generate(module)
}

func addColumn(field *configs.FieldTemplate) {
	err := interact.NewInteraction("Masukkan Nama Kolom?").Resolve(&field.Name)
	if err != nil {
		panic(err)
	}

	err = interact.NewInteraction("Masukkan Tipe Data?",
		interact.Choice{Display: "bool", Value: "bool"},
		interact.Choice{Display: "string", Value: "string"},
		interact.Choice{Display: "byte", Value: "byte"},
		interact.Choice{Display: "int", Value: "int"},
		interact.Choice{Display: "int8", Value: "int8"},
		interact.Choice{Display: "int16", Value: "int16"},
		interact.Choice{Display: "int32", Value: "int32"},
		interact.Choice{Display: "int64", Value: "int64"},
		interact.Choice{Display: "uint", Value: "uint"},
		interact.Choice{Display: "uint8", Value: "uint8"},
		interact.Choice{Display: "uint16", Value: "uint16"},
		interact.Choice{Display: "uint32", Value: "uint32"},
		interact.Choice{Display: "uint64", Value: "uint64"},
		interact.Choice{Display: "float32", Value: "float32"},
		interact.Choice{Display: "float64", Value: "float64"},
		interact.Choice{Display: "complex64", Value: "complex64"},
		interact.Choice{Display: "complex128", Value: "complex128"},
		interact.Choice{Display: "rune", Value: "rune"},
		interact.Choice{Display: "uintptr", Value: "uintptr"},
	).Resolve(&field.Type)
	if err != nil {
		panic(err)
	}

	field.IsRequired = true
	err = interact.NewInteraction("Apakah Kolom Wajib Diisi?").Resolve(&field.IsRequired)
	if err != nil {
		panic(err)
	}
}
