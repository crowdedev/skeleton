package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	configs "github.com/crowdeco/bima/configs"
	dic "github.com/crowdeco/skeleton/generated/dic"
	"github.com/fatih/color"
	"github.com/jinzhu/copier"
	"github.com/joho/godotenv"
	"github.com/vito/go-interact/interact"
	"golang.org/x/mod/modfile"
)

func main() {
	godotenv.Load()
	container, _ := dic.NewContainer()
	util := container.GetBimaUtilCli()

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

		_, err := exec.Command("sh", "proto_gen.sh").Output()
		if err != nil {
			util.Println(err.Error())
			os.Exit(1)
		}
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

	_, err := exec.Command("go", "mod", "tidy").Output()
	if err != nil {
		util.Println(err.Error())
		os.Exit(1)
	}

	_, err = exec.Command("go", "run", "cmds/dic/main.go").Output()
	if err != nil {
		util.Println(err.Error())
		os.Exit(1)
	}

	util.Println("By:")
	util.Println("𝕒𝕕𝟛𝕟")
}

func unregister(container *dic.Container, util *color.Color, module string) {
	moduleParser := container.GetBimaConfigParserModule()
	word := container.GetBimaUtilWord()
	pluralizer := container.GetBimaUtilPluralizer()
	moduleName := word.Camelcase(pluralizer.Singular(module))
	modulePlural := word.Underscore(pluralizer.Plural(moduleName))
	list := moduleParser.Parse()

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

	jsonModules := fmt.Sprintf("%s/swaggers/modules.json", workDir)
	file, _ := ioutil.ReadFile(jsonModules)
	modulesJson := []configs.ModuleJson{}
	registered := modulesJson
	json.Unmarshal(file, &modulesJson)
	for _, v := range modulesJson {
		if v.Name != moduleName {
			registered = append(registered, v)
		}
	}

	registeredByte, _ := json.Marshal(registered)
	ioutil.WriteFile(jsonModules, registeredByte, 0644)

	packageName := modfile.ModulePath(mod)
	yaml := fmt.Sprintf("%s/configs/modules.yaml", workDir)
	file, _ = ioutil.ReadFile(yaml)
	modules := string(file)

	provider := fmt.Sprintf("%s/configs/provider.go", workDir)
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
	os.Remove(fmt.Sprintf("%s/protos/builds/%s_grpc.pb.go", workDir, word.Underscore(module)))
	os.Remove(fmt.Sprintf("%s/protos/builds/%s.pb.gw.go", workDir, word.Underscore(module)))
	os.Remove(fmt.Sprintf("%s/swaggers/%s.swagger.json", workDir, word.Underscore(module)))

	util.Println("Modul berhasil dihapus")
}

func register(container *dic.Container, util *color.Color) {
	generator := container.GetBimaModuleGenerator()
	module := container.GetBimaTemplateModule()
	field := container.GetBimaTemplateField()
	word := container.GetBimaUtilWord()
	mapType := container.GetBimaConfigType()

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
	moduleName(util, module)

	index := 2
	more := true
	for more {
		err := interact.NewInteraction("Tambah Kolom?").Resolve(&more)
		if err != nil {
			util.Println(err.Error())
			os.Exit(1)
		}

		if more {
			addColumn(util, field, mapType)

			field.Name = strings.Replace(field.Name, " ", "", -1)
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

	if field.Name == "" {
		util.Println("Nama Kolom tidak boleh kosong")
		addColumn(util, field, mapType)
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

func moduleName(util *color.Color, module *configs.ModuleTemplate) {
	err := interact.NewInteraction("Masukkan Nama Modul?").Resolve(&module.Name)
	if err != nil {
		util.Println(err.Error())
		os.Exit(1)
	}

	if strings.HasSuffix(module.Name, "test") {
		util.Println("Modul mengandung kata 'test'")
		moduleName(util, module)
	}

	if module.Name == "" {
		util.Println("Nama Modul tidak boleh kosong")
		moduleName(util, module)
	}
}
