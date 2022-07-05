package skeleton

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/KejawenLab/bima/v3"
	"github.com/KejawenLab/bima/v3/configs"
	"github.com/KejawenLab/bima/v3/events"
	"github.com/KejawenLab/bima/v3/generators"
	"github.com/KejawenLab/bima/v3/middlewares"
	"github.com/KejawenLab/bima/v3/parsers"
	"github.com/KejawenLab/bima/v3/routes"
	"github.com/KejawenLab/bima/v3/utils"
	"github.com/KejawenLab/skeleton/v3/generated/engine"
	"github.com/KejawenLab/skeleton/v3/generated/generator"
	"github.com/fatih/color"
	"github.com/gertd/go-pluralize"
	"github.com/goccy/go-json"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/copier"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/vito/go-interact/interact"
	"golang.org/x/mod/modfile"
	"gopkg.in/yaml.v2"
)

type (
	Application string
	Module      string
)

func (_ Application) Run(config string) {
	if config == "" {
		config = ".env"
	}

	container, err := engine.NewContainer(bima.Application)
	if err != nil {
		panic(err)
	}

	env := container.GetBimaConfig()
	loadEnv(env, config, filepath.Ext(config))

	workDir, _ := os.Getwd()
	util := color.New(color.FgCyan, color.Bold)

	var ext []string
	var cName bytes.Buffer

	ext = parsers.ParseModule(workDir)
	servers := make([]configs.Server, len(ext))
	for k, c := range ext {
		cName.Reset()
		cName.WriteString(c)
		cName.WriteString(":server")

		servers[k] = container.Get(cName.String()).(configs.Server)
	}

	ext = parsers.ParseListener(workDir)
	listeners := make([]events.Listener, len(ext))
	for k, c := range ext {
		cName.Reset()
		cName.WriteString("bima:listener:")
		cName.WriteString(c)

		listeners[k] = container.Get(cName.String()).(events.Listener)
	}

	ext = parsers.ParseMiddleware(workDir)
	hooks := make([]middlewares.Middleware, len(ext))
	for k, c := range ext {
		cName.Reset()
		cName.WriteString("bima:middleware:")
		cName.WriteString(c)

		hooks[k] = container.Get(cName.String()).(middlewares.Middleware)
	}

	ext = parsers.ParseLogger(workDir)
	extensions := make([]logrus.Hook, len(ext))
	for k, c := range ext {
		cName.Reset()
		cName.WriteString("bima:logger:extension:")
		cName.WriteString(c)

		extensions[k] = container.Get(cName.String()).(logrus.Hook)
	}

	ext = parsers.ParseRoute(workDir)
	handlers := make([]routes.Route, len(ext))
	for k, c := range ext {
		cName.Reset()
		cName.WriteString("bima:route:")
		cName.WriteString(c)

		handlers[k] = container.Get(cName.String()).(routes.Route)
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

func (m Module) Run(module string, config string) {
	if config == "" {
		config = ".env"
	}

	container, err := engine.NewContainer(bima.Application)
	if err != nil {
		panic(err)
	}

	env := container.GetBimaConfig()
	loadEnv(env, config, filepath.Ext(config))

	switch m {
	case "add":
		m.register(module, env.ApiVersion, env.Db.Driver)
	case "remove":
		m.remove(module)
	}
}

func (m Module) register(module string, apiVersion string, driver string) {
	container, err := generator.NewContainer(bima.Generator)
	if err != nil {
		panic(err)
	}   

	generator := container.GetBimaModuleGenerator()
	util := color.New(color.FgCyan, color.Bold)

	generator.ApiVersion = apiVersion
	generator.Driver = driver
	register(generator, util, module)
	_, err = exec.Command("sh", "proto_gen.sh").Output()
	if err != nil {
		util.Println("Error generate code from proto files")
		os.Exit(1)
	}

	_, err = exec.Command("go", "mod", "tidy").Output()
	if err != nil {
		util.Println("Error update dependencies")
		os.Exit(1)
	}

	_, err = exec.Command("go", "run", "dumper/main.go").Output()
	if err != nil {
		util.Println("Error update DI Container")
		os.Exit(1)
	}

	util.Println("By:")
	util.Println("ad3n")
}

func (m Module) remove(module string) {
	util := color.New(color.FgCyan, color.Bold)
	unregister(util, module)

	_, err := exec.Command("go", "run", "dumper/main.go").Output()
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

func loadEnv(config *configs.Env, filePath string, ext string) {
	switch ext {
	case ".env":
		godotenv.Load()
		processDotEnv(config)
	case ".yaml":
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalln(err.Error())
		}

		err = yaml.Unmarshal(content, config)
		if err != nil {
			log.Fatalln(err.Error())
		}
	case ".json":
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalln(err.Error())
		}

		err = json.Unmarshal(content, config)
		if err != nil {
			log.Fatalln(err.Error())
		}
	}

	if config.Secret == "" {
		hasher := sha256.New()
		hasher.Write([]byte(time.Now().Format(time.RFC3339)))

		config.Secret = base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	}

	if config.RequestIDHeader == "" {
		config.RequestIDHeader = "X-Request-Id"
	}
}

func processDotEnv(config *configs.Env) {
	config.ApiVersion = os.Getenv("API_VERSION")
	config.Secret = os.Getenv("APP_SECRET")
	config.RequestIDHeader = os.Getenv("REQUEST_ID_HEADER")
	config.Debug, _ = strconv.ParseBool(os.Getenv("APP_DEBUG"))
	config.HttpPort, _ = strconv.Atoi(os.Getenv("APP_PORT"))
	config.RpcPort, _ = strconv.Atoi(os.Getenv("GRPC_PORT"))

	sName := os.Getenv("APP_NAME")
	config.Service = configs.Service{
		Name:           sName,
		ConnonicalName: strcase.ToDelimited(sName, '_'),
	}

	dbPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	config.Db = configs.Db{
		Host:     os.Getenv("DB_HOST"),
		Port:     dbPort,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		Driver:   os.Getenv("DB_DRIVER"),
	}

	esPort, _ := strconv.Atoi(os.Getenv("ELASTICSEARCH_PORT"))
	config.Elasticsearch = configs.Elasticsearch{
		Host:  os.Getenv("ELASTICSEARCH_HOST"),
		Port:  esPort,
		Index: config.Db.Name,
	}

	amqpPort, _ := strconv.Atoi(os.Getenv("AMQP_PORT"))
	config.Amqp = configs.Amqp{
		Host:     os.Getenv("AMQP_HOST"),
		Port:     amqpPort,
		User:     os.Getenv("AMQP_USER"),
		Password: os.Getenv("AMQP_PASSWORD"),
	}

	config.CacheLifetime, _ = strconv.Atoi(os.Getenv("CACHE_LIFETIME"))
}

func unregister(util *color.Color, module string) {
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

func register(generator *generators.Factory, util *color.Color, name string) {
	module := generators.ModuleTemplate{}
	field := generators.FieldTemplate{}
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
			addColumn(util, &field, mapType)

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
