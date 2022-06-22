package skeleton

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/KejawenLab/bima/v2/configs"
	"github.com/KejawenLab/bima/v2/events"
	"github.com/KejawenLab/bima/v2/generators"
	"github.com/KejawenLab/bima/v2/middlewares"
	"github.com/KejawenLab/bima/v2/parsers"
	"github.com/KejawenLab/bima/v2/routes"
	"github.com/KejawenLab/bima/v2/utils"
	"github.com/KejawenLab/skeleton/generated/dic"
	"github.com/fatih/color"
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/copier"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/vito/go-interact/interact"
	"golang.org/x/mod/modfile"
)

type (
	Application string
	Module      string
)

func (_ Application) Run() {
	workDir, _ := os.Getwd()
	godotenv.Load()
	container, _ := dic.NewContainer()
	util := color.New(color.FgCyan, color.Bold)
	env := container.GetBimaConfigEnv()

	var servers []configs.Server
	for _, c := range parsers.ParseModule(workDir) {
		servers = append(servers, container.Get(fmt.Sprintf("%s:server", c)).(configs.Server))
	}

	var listeners []events.Listener
	for _, c := range parsers.ParseListener(workDir) {
		listeners = append(listeners, container.Get(fmt.Sprintf("bima:listener:%s", c)).(events.Listener))
	}

	var hooks []middlewares.Middleware
	for _, c := range parsers.ParseMiddleware(workDir) {
		hooks = append(hooks, container.Get(fmt.Sprintf("bima:middleware:%s", c)).(middlewares.Middleware))
	}

	var extensions []logrus.Hook
	for _, c := range parsers.ParseLogger(workDir) {
		extensions = append(extensions, container.Get(fmt.Sprintf("bima:logger:extension:%s", c)).(logrus.Hook))
	}

	var handlers []routes.Route
	for _, c := range parsers.ParseRoute(workDir) {
		handlers = append(handlers, container.Get(fmt.Sprintf("bima:route:%s", c)).(routes.Route))
	}

	container.GetBimaRouterMux().Register(handlers)
	container.GetBimaLoggerExtension().Register(extensions)
	container.GetBimaMiddlewareFactory().Register(hooks)
	container.GetBimaEventDispatcher().Register(listeners)
	container.GetBimaRouterGateway().Register(servers)

	util.Printf("✓ ")
	fmt.Printf("REST running on %d\n", env.HttpPort)
	if env.Debug {
		util.Printf("✓ ")
		fmt.Println("Api Doc ready on /api/docs")
	}

	container.GetBimaApplication().Run(servers)
}

func (m Module) Run(module string) {
	switch m {
	case "register":
		m.register(module)
	case "remove":
		m.remove(module)
	}
}

func (m Module) register(module string) {
	godotenv.Load()
	container, _ := dic.NewContainer()
	util := color.New(color.FgCyan, color.Bold)

	register(container, util, module)
	_, err := exec.Command("sh", "proto_gen.sh").Output()
	if err != nil {
		util.Println("Error generate code from proto files")
		os.Exit(1)
	}

	_, err = exec.Command("go", "mod", "tidy").Output()
	if err != nil {
		util.Println("Error update dependencies")
		os.Exit(1)
	}

	_, err = exec.Command("go", "run", "cmd/main.go", "dump").Output()
	if err != nil {
		util.Println("Error update DI Container")
		os.Exit(1)
	}

	util.Println("By:")
	util.Println("ad3n")
}

func (m Module) remove(module string) {
	godotenv.Load()
	container, _ := dic.NewContainer()
	util := color.New(color.FgCyan, color.Bold)
	unregister(container, util, module)

	_, err := exec.Command("go", "run", "cmd/main.go", "dump").Output()
	if err != nil {
		util.Println("Error update DI Container")
		os.Exit(1)
	}

	_, err = exec.Command("go", "mod", "tidy").Output()
	if err != nil {
		util.Println("Error update dependencies")
		os.Exit(1)
	}

	util.Println("By:")
	util.Println("ad3n")
}

func unregister(container *dic.Container, util *color.Color, module string) {
	workDir, _ := os.Getwd()
	pluralizer := pluralize.NewClient()
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

	mod, err := os.ReadFile(fmt.Sprintf("%s/go.mod", workDir))
	if err != nil {
		panic(err)
	}

	jsonModules := fmt.Sprintf("%s/swaggers/modules.json", workDir)
	file, _ := os.ReadFile(jsonModules)
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
	os.WriteFile(jsonModules, registeredByte, 0644)

	packageName := modfile.ModulePath(mod)
	yaml := fmt.Sprintf("%s/configs/modules.yaml", workDir)
	file, _ = os.ReadFile(yaml)
	modules := string(file)

	provider := fmt.Sprintf("%s/configs/provider.go", workDir)
	file, _ = os.ReadFile(provider)
	codeblock := string(file)

	modRegex := regexp.MustCompile(fmt.Sprintf("(?m)[\r\n]+^.*module:%s.*$", moduleUnderscore))
	modules = modRegex.ReplaceAllString(modules, "")
	os.WriteFile(yaml, []byte(modules), 0644)

	regex := regexp.MustCompile(fmt.Sprintf("(?m)[\r\n]+^.*%s.*$", fmt.Sprintf("%s/%s", packageName, modulePlural)))
	codeblock = regex.ReplaceAllString(codeblock, "")

	codeblock = modRegex.ReplaceAllString(codeblock, "")
	os.WriteFile(provider, []byte(codeblock), 0644)

	os.RemoveAll(fmt.Sprintf("%s/%s", workDir, modulePlural))
	os.Remove(fmt.Sprintf("%s/protos/%s.proto", workDir, moduleUnderscore))
	os.Remove(fmt.Sprintf("%s/protos/builds/%s_grpc.pb.go", workDir, moduleUnderscore))
	os.Remove(fmt.Sprintf("%s/protos/builds/%s.pb.go", workDir, moduleUnderscore))
	os.Remove(fmt.Sprintf("%s/protos/builds/%s.pb.gw.go", workDir, moduleUnderscore))
	os.Remove(fmt.Sprintf("%s/swaggers/%s.swagger.json", workDir, moduleUnderscore))

	util.Println("Module deleted")
}

func register(container *dic.Container, util *color.Color, name string) {
	generator := container.GetBimaModuleGenerator()
	module := container.GetBimaTemplateModule()
	field := container.GetBimaTemplateField()
	mapType := utils.NewType()

	util.Println("Welcome to Bima Skeleton Module Generator")
	module.Name = name

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

func addColumn(util *color.Color, field *generators.FieldTemplate, mapType utils.Type) {
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
