package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/KejawenLab/bima/v2/configs"
	"github.com/KejawenLab/bima/v2/generators"
	"github.com/KejawenLab/bima/v2/parsers"
	"github.com/KejawenLab/skeleton/generated/dic"
	"github.com/fatih/color"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/copier"
	"github.com/joho/godotenv"
	"github.com/vito/go-interact/interact"
	"golang.org/x/mod/modfile"
)

func main() {
	godotenv.Load()
	container, _ := dic.NewContainer()
	util := color.New(color.FgCyan, color.Bold)

	if len(os.Args) < 2 {
		util.Println("Usage:")
		util.Println("go run cmds/module/main.go register")
		util.Println("go run cmds/module/main.g unregister")
		util.Println("By:")
		util.Println("ad3n")
		os.Exit(1)
	}

	if os.Args[1] != "register" && os.Args[1] != "unregister" {
		util.Println("Command unknown")
		util.Println("By:")
		util.Println("ad3n")
		os.Exit(1)
	}

	if os.Args[1] == "register" {
		register(container, util)

		_, err := exec.Command("sh", "proto_gen.sh").Output()
		if err != nil {
			util.Println(err.Error())
			os.Exit(1)
		}

		_, err = exec.Command("go", "mod", "tidy").Output()
		if err != nil {
			util.Println(err.Error())
			os.Exit(1)
		}

		_, err = exec.Command("go", "run", "cmds/dic/main.go").Output()
		if err != nil {
			util.Println(err.Error())
			os.Exit(1)
		}
	}

	if os.Args[1] == "unregister" {
		if len(os.Args) < 3 {
			util.Println("Module name is required")
			util.Println("By:")
			util.Println("ad3n")
			os.Exit(1)
		}

		unregister(container, util, os.Args[2])

		_, err := exec.Command("go", "run", "cmds/dic/main.go").Output()
		if err != nil {
			util.Println(err.Error())
			os.Exit(1)
		}

		_, err = exec.Command("go", "mod", "tidy").Output()
		if err != nil {
			util.Println(err.Error())
			os.Exit(1)
		}
	}

	util.Println("By:")
	util.Println("ad3n")
}

func unregister(container *dic.Container, util *color.Color, module string) {
	workDir, _ := os.Getwd()
	pluralizer := container.GetBimaUtilPluralizer()
	moduleName := strcase.ToCamel(pluralizer.Singular(module))
	modulePlural := strcase.ToDelimited(pluralizer.Plural(moduleName), '_')
	moduleUnderscore := strcase.ToDelimited(module, '_')
	list := parsers.ParseModule(workDir)

	exist := false
	for _, v := range list {
		if v == fmt.Sprintf("module:%s", moduleUnderscore) {
			exist = true
			break
		}
	}

	if !exist {
		util.Println("Module is not registered")
		return
	}

	mod, err := ioutil.ReadFile(fmt.Sprintf("%s/go.mod", workDir))
	if err != nil {
		panic(err)
	}

	jsonModules := fmt.Sprintf("%s/swaggers/modules.json", workDir)
	file, _ := ioutil.ReadFile(jsonModules)
	modulesJson := []generators.ModuleJson{}
	registered := modulesJson
	json.Unmarshal(file, &modulesJson)
	for _, v := range modulesJson {
		if v.Name != moduleName {
			mUrl, _ := url.Parse(v.Url)
			query := mUrl.Query()

			query.Set("v", strconv.Itoa(int(time.Now().UnixMicro())))
			mUrl.RawQuery = query.Encode()
			v.Url = mUrl.String()
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

	modRegex := regexp.MustCompile(fmt.Sprintf("(?m)[\r\n]+^.*module:%s.*$", moduleUnderscore))
	modules = modRegex.ReplaceAllString(modules, "")
	ioutil.WriteFile(yaml, []byte(modules), 0644)

	regex := regexp.MustCompile(fmt.Sprintf("(?m)[\r\n]+^.*%s.*$", fmt.Sprintf("%s/%s", packageName, modulePlural)))
	codeblock = regex.ReplaceAllString(codeblock, "")

	codeblock = modRegex.ReplaceAllString(codeblock, "")
	ioutil.WriteFile(provider, []byte(codeblock), 0644)

	os.RemoveAll(fmt.Sprintf("%s/%s", workDir, modulePlural))
	os.Remove(fmt.Sprintf("%s/protos/%s.proto", workDir, moduleUnderscore))
	os.Remove(fmt.Sprintf("%s/protos/builds/%s_grpc.pb.go", workDir, moduleUnderscore))
	os.Remove(fmt.Sprintf("%s/protos/builds/%s.pb.go", workDir, moduleUnderscore))
	os.Remove(fmt.Sprintf("%s/protos/builds/%s.pb.gw.go", workDir, moduleUnderscore))
	os.Remove(fmt.Sprintf("%s/swaggers/%s.swagger.json", workDir, moduleUnderscore))

	util.Println("Module deleted")
}

func register(container *dic.Container, util *color.Color) {
	generator := container.GetBimaModuleGenerator()
	module := container.GetBimaTemplateModule()
	field := container.GetBimaTemplateField()
	mapType := container.GetBimaConfigType()

	util.Println("Welcome to Bima Skeleton Module Generator")
	moduleName(util, module)

	index := 2
	more := true
	for more {
		err := interact.NewInteraction("Add new column?").Resolve(&more)
		if err != nil {
			util.Println(err.Error())
			os.Exit(1)
		}

		if more {
			addColumn(util, field, mapType)

			field.Name = strings.Replace(field.Name, " ", "", -1)
			column := generators.FieldTemplate{}

			copier.Copy(&column, field)

			column.Index = index
			column.Name = strings.Title(column.Name)
			column.NameUnderScore = strcase.ToDelimited(column.Name, '_')
			module.Fields = append(module.Fields, &column)

			field.Name = ""
			field.ProtobufType = ""

			index++
		}
	}

	if len(module.Fields) < 1 {
		util.Println("You must have at least one column in table")
		os.Exit(1)
	}

	generator.Generate(module)

	workDir, _ := os.Getwd()
	util.Println(fmt.Sprintf("Module registered in %s/modules.yaml", workDir))
}

func addColumn(util *color.Color, field *generators.FieldTemplate, mapType *configs.Type) {
	err := interact.NewInteraction("Input column name?").Resolve(&field.Name)
	if err != nil {
		util.Println(err.Error())
		os.Exit(1)
	}

	if field.Name == "" {
		util.Println("Column name is required")
		addColumn(util, field, mapType)
	}

	field.ProtobufType = "string"
	err = interact.NewInteraction("Input data type?",
		interact.Choice{Display: "double", Value: "double"},
		interact.Choice{Display: "float", Value: "float"},
		interact.Choice{Display: "int32", Value: "int32"},
		interact.Choice{Display: "int64", Value: "int64"},
		interact.Choice{Display: "uint32", Value: "uint32"},
		interact.Choice{Display: "sint32", Value: "sint32"},
		interact.Choice{Display: "sint64", Value: "sint64"},
		interact.Choice{Display: "fixed32", Value: "fixed32"},
		interact.Choice{Display: "fixed64", Value: "fixed64"},
		interact.Choice{Display: "sfixed32", Value: "sfixed32"},
		interact.Choice{Display: "sfixed64", Value: "sfixed64"},
		interact.Choice{Display: "bool", Value: "bool"},
		interact.Choice{Display: "string", Value: "string"},
		interact.Choice{Display: "bytes", Value: "bytes"},
	).Resolve(&field.ProtobufType)
	if err != nil {
		util.Println(err.Error())
		os.Exit(1)
	}
	field.GolangType = mapType.Value(field.ProtobufType)

	field.IsRequired = true
	err = interact.NewInteraction("Is column required?").Resolve(&field.IsRequired)
	if err != nil {
		util.Println(err.Error())
		os.Exit(1)
	}
}

func moduleName(util *color.Color, module *generators.ModuleTemplate) {
	err := interact.NewInteraction("Input module name?").Resolve(&module.Name)
	if err != nil {
		util.Println(err.Error())
		os.Exit(1)
	}

	if strings.HasSuffix(module.Name, "test") {
		util.Println("Module 'test' is not allowed")
		moduleName(util, module)
	}

	if module.Name == "" {
		util.Println("Module name is required")
		moduleName(util, module)
	}
}
