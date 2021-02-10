package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	configs "github.com/crowdeco/skeleton/configs"
	dic "github.com/crowdeco/skeleton/generated/dic"
	"github.com/fatih/color"
	"github.com/jinzhu/copier"
	"github.com/vito/go-interact/interact"
	"golang.org/x/mod/modfile"
)

func main() {
	container, _ := dic.NewContainer()
	util := container.GetCoreUtilCli()

	if len(os.Args) < 2 {
		util.Println("Cara Penggunaan:")
		util.Println("go run cmds/module/main.go register")
		util.Println("go run cmds/module/main.g unregister")
		util.Println("By:")
		util.Println("𝕒𝕕𝟛𝕟")
		os.Exit(1)
	}

	if os.Args[1] != "register" && os.Args[1] != "unregister" {
		util.Println("Perintah tidak diketahui")
		util.Println("By:")
		util.Println("𝕒𝕕𝟛𝕟")
		os.Exit(1)
	}

	if os.Args[1] == "register" {
		register(container, util)
	}

	if os.Args[1] == "unregister" {
		if len(os.Args) < 3 {
			util.Println("Modul wajib diisi")
			util.Println("By:")
			util.Println("𝕒𝕕𝟛𝕟")
			os.Exit(1)
		}

		unregister(container, util, os.Args[2])
	}

	_, err := exec.Command("sh", "proto_gen.sh").Output()
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

	util.Println("By:")
	util.Println("𝕒𝕕𝟛𝕟")
}

func unregister(container *dic.Container, util *color.Color, module string) {
	config := container.GetCoreConfigParser()
	word := container.GetCoreUtilWord()
	pluralizer := container.GetCoreUtilPluralizer()
	moduleName := word.Camelcase(pluralizer.Singular(module))
	modulePlural := word.Underscore(pluralizer.Plural(moduleName))
	list := config.ParseModules()

	exist := false
	for _, v := range list {
		if v == fmt.Sprintf("module:%s", word.Underscore(module)) {
			exist = true
			break
		}
	}

	if !exist {
		util.Println("Modul tidak terdaftar")
		return
	}

	workDir, _ := os.Getwd()
	mod, err := ioutil.ReadFile(fmt.Sprintf("%s/go.mod", workDir))
	if err != nil {
		panic(err)
	}

	packageName := modfile.ModulePath(mod)
	yaml := fmt.Sprintf("%s/modules.yaml", workDir)
	file, _ := ioutil.ReadFile(yaml)
	modules := string(file)

	provider := fmt.Sprintf("%s/dics/provider.go", workDir)
	file, _ = ioutil.ReadFile(provider)
	codeblock := string(file)

	modRegex := regexp.MustCompile(fmt.Sprintf("(?m)[\r\n]+^.*module:%s.*$", word.Underscore(module)))
	modules = modRegex.ReplaceAllString(modules, "")
	ioutil.WriteFile(yaml, []byte(modules), 0644)

	regex := regexp.MustCompile(fmt.Sprintf("(?m)[\r\n]+^.*%s.*$", fmt.Sprintf("%s/%s", packageName, modulePlural)))
	codeblock = regex.ReplaceAllString(codeblock, "")

	codeblock = modRegex.ReplaceAllString(codeblock, "")
	ioutil.WriteFile(provider, []byte(codeblock), 0644)

	os.RemoveAll(fmt.Sprintf("%s/%s", workDir, modulePlural))
	os.Remove(fmt.Sprintf("%s/protos/%s.proto", workDir, word.Underscore(module)))
	os.Remove(fmt.Sprintf("%s/protos/builds/%s.pb.go", workDir, word.Underscore(module)))
	os.Remove(fmt.Sprintf("%s/protos/builds/%s.pb.gw.go", workDir, word.Underscore(module)))

	util.Println("Modul berhasil dihapus")
}

func register(container *dic.Container, util *color.Color) {
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

	err := interact.NewInteraction("Masukkan Nama Modul?").Resolve(&module.Name)
	if err != nil {
		util.Println(err.Error())
		os.Exit(1)
	}

	if strings.HasSuffix(module.Name, "test") {
		util.Println("Modul mengandung kata 'test'")
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

	workDir, _ := os.Getwd()
	util.Println(fmt.Sprintf("Module berhasil didaftarkan pada file: %s/modules.yaml", workDir))
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
