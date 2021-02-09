package main

import (
	"os"
	"os/exec"
	"strings"

	configs "github.com/crowdeco/skeleton/configs"
	dic "github.com/crowdeco/skeleton/generated/dic"
	"github.com/fatih/color"
	"github.com/jinzhu/copier"
	"github.com/vito/go-interact/interact"
)

func main() {
	container, _ := dic.NewContainer()
	util := container.GetCoreUtilCli()

	if len(os.Args) < 2 {
		util.Println("Cara Penggunaan:")
		util.Println("go run cmds/module/main.go register")
		util.Println("go run cmds/module/main.g unregister")
		os.Exit(1)
	}

	if os.Args[1] != "register" && os.Args[0] != "unregister" {
		util.Println("Perintah tidak diketahui")
		os.Exit(1)
	}

	if os.Args[1] == "register" {
		startGenerator(container, util)
	}

	if os.Args[1] == "unregister" {
		//@todo unregister module
	}

	util.Println("By:")
	util.Println("𝕒𝕕𝟛𝕟")
}

func startGenerator(container *dic.Container, util *color.Color) {
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
		util.Println(err.Error())
		os.Exit(1)
	}

	if strings.HasSuffix(module.Name, "test") {
		util.Println("Module has 'test' as surfix")
		os.Exit(1)
	}

	index := 2
	more := true
	for more {
		err = interact.NewInteraction("Tambah Kolom?").Resolve(&more)
		if err != nil {
			util.Println(err.Error())
			os.Exit(1)
		}

		if more {
			addColumn(util, field, mapType)

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

	if len(module.Fields) < 1 {
		util.Println("Harus ada minimal satu kolom dalam tabel")
		os.Exit(1)
	}

	generator.Generate(module)
	_, err = exec.Command("sh", "proto_gen.sh").Output()
	if err != nil {
		util.Println(err.Error())
		os.Exit(1)
	}

	_, err = exec.Command("go", "run", "cmds/dic/main.go").Output()
	if err != nil {
		util.Println(err.Error())
		os.Exit(1)
	}

	_, err = exec.Command("go", "mod", "tidy").Output()
	if err != nil {
		util.Println(err.Error())
		os.Exit(1)
	}

	workDir, _ := os.Getwd()
	util.Println("Module berhasil didaftarkan pada file: %s/modules.yaml", workDir)
}

func addColumn(util *color.Color, field *configs.FieldTemplate, mapType *configs.Type) {
	err := interact.NewInteraction("Masukkan Nama Kolom?").Resolve(&field.Name)
	if err != nil {
		util.Println(err.Error())
		os.Exit(1)
	}

	var types []interact.Choice
	for k := range mapType.List() {
		types = append(types, interact.Choice{Display: k, Value: k})
	}

	field.ProtobufType = "string"
	err = interact.NewInteraction("Masukkan Tipe Data?", types...).Resolve(&field.ProtobufType)
	if err != nil {
		util.Println(err.Error())
		os.Exit(1)
	}
	field.GolangType = mapType.Value(field.ProtobufType)

	field.IsRequired = true
	err = interact.NewInteraction("Apakah Kolom Wajib Diisi?").Resolve(&field.IsRequired)
	if err != nil {
		util.Println(err.Error())
		os.Exit(1)
	}
}
