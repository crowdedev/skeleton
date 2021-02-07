package main

import (
	"os/exec"
	"strings"

	configs "github.com/crowdeco/skeleton/configs"
	dic "github.com/crowdeco/skeleton/generated/dic"
	"github.com/jinzhu/copier"
	"github.com/vito/go-interact/interact"
)

func main() {
	container, _ := dic.NewContainer()

	util := container.GetCoreUtilCli()
	generator := container.GetCoreModuleGenerator()
	module := container.GetCoreTemplateModule()
	field := container.GetCoreTemplateField()
	word := container.GetCoreUtilWord()

	util.Println(`
     ______                                           __                   ______   __                  __              __
    /      \                                         /  |                 /      \ /  |                /  |            /  |
   /$$$$$$  |  ______    ______   __   __   __   ____$$ |  ______        /$$$$$$  |$$ |   __   ______  $$ |  ______   _$$ |_     ______   _______
   $$ |  $$/  /      \  /      \ /  | /  | /  | /    $$ | /      \       $$ \__$$/ $$ |  /  | /      \ $$ | /      \ / $$   |   /      \ /       \
   $$ |      /$$$$$$  |/$$$$$$  |$$ | $$ | $$ |/$$$$$$$ |/$$$$$$  |      $$      \ $$ |_/$$/ /$$$$$$  |$$ |/$$$$$$  |$$$$$$/   /$$$$$$  |$$$$$$$  |
   $$ |   __ $$ |  $$/ $$ |  $$ |$$ | $$ | $$ |$$ |  $$ |$$    $$ |       $$$$$$  |$$   $$<  $$    $$ |$$ |$$    $$ |  $$ | __ $$ |  $$ |$$ |  $$ |
   $$ \__/  |$$ |      $$ \__$$ |$$ \_$$ \_$$ |$$ \__$$ |$$$$$$$$/       /  \__$$ |$$$$$$  \ $$$$$$$$/ $$ |$$$$$$$$/   $$ |/  |$$ \__$$ |$$ |  $$ |
   $$    $$/ $$ |      $$    $$/ $$   $$   $$/ $$    $$ |$$       |      $$    $$/ $$ | $$  |$$       |$$ |$$       |  $$  $$/ $$    $$/ $$ |  $$ |
    $$$$$$/  $$/        $$$$$$/   $$$$$/$$$$/   $$$$$$$/  $$$$$$$/        $$$$$$/  $$/   $$/  $$$$$$$/ $$/  $$$$$$$/    $$$$/   $$$$$$/  $$/   $$/

`)

	err := interact.NewInteraction("Masukkan Nama Table?").Resolve(&module.Name)
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
			column.NameUnderScore = word.Underscore(column.Name)
			module.Fields = append(module.Fields, &column)

			field.Name = ""
			field.Type = ""

			index++
		}
	}

	generator.Generate(module)
	_, err = exec.Command("sh", "proto_gen.sh").Output()
	if err != nil {
		panic(err)
	}

	_, err = exec.Command("go", "run", "cmd/dic/main.go").Output()
	if err != nil {
		panic(err)
	}

	_, err = exec.Command("go", "mod", "tidy").Output()
	if err != nil {
		panic(err)
	}
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
