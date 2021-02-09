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
	mapType := container.GetCoreConfigType()

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

	if strings.HasSuffix(module.Name, "test") {
		panic("Module has 'test' as surfix")
	}

	index := 2
	more := true
	for more {
		err = interact.NewInteraction("Tambah Kolom?").Resolve(&more)
		if err != nil {
			panic(err)
		}

		if more {
			addColumn(field, mapType)

			column := configs.FieldTemplate{}

			copier.Copy(&column, field)

			column.Index = index
			column.Name = strings.Title(column.Name)
			column.NameUnderScore = word.Underscore(column.Name)
			module.Fields = append(module.Fields, &column)

			field.Name = ""
			field.ProtobufType = ""

			index++
		}
	}

	generator.Generate(module)
	_, err = exec.Command("sh", "proto_gen.sh").Output()
	if err != nil {
		panic(err)
	}

	_, err = exec.Command("go", "run", "cmds/dic/main.go").Output()
	if err != nil {
		panic(err)
	}

	_, err = exec.Command("go", "mod", "tidy").Output()
	if err != nil {
		panic(err)
	}
}

func addColumn(field *configs.FieldTemplate, mapType *configs.Type) {
	err := interact.NewInteraction("Masukkan Nama Kolom?").Resolve(&field.Name)
	if err != nil {
		panic(err)
	}

	var types []interact.Choice
	for k, _ := range mapType.List() {
		types = append(types, interact.Choice{Display: k, Value: k})
	}

	field.ProtobufType = "string"
	err = interact.NewInteraction("Masukkan Tipe Data?", types...).Resolve(&field.ProtobufType)
	if err != nil {
		panic(err)
	}
	field.GolangType = mapType.Value(field.ProtobufType)

	field.IsRequired = true
	err = interact.NewInteraction("Apakah Kolom Wajib Diisi?").Resolve(&field.IsRequired)
	if err != nil {
		panic(err)
	}
}
