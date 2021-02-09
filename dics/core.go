package dics

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/ThreeDotsLabs/watermill"
	amqp "github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
	configs "github.com/crowdeco/skeleton/configs"
	drivers "github.com/crowdeco/skeleton/configs/drivers"
	generators "github.com/crowdeco/skeleton/generators"
	handlers "github.com/crowdeco/skeleton/handlers"
	interfaces "github.com/crowdeco/skeleton/interfaces"
	creates "github.com/crowdeco/skeleton/listeners/creates"
	deletes "github.com/crowdeco/skeleton/listeners/deletes"
	updates "github.com/crowdeco/skeleton/listeners/updates"
	middlewares "github.com/crowdeco/skeleton/middlewares"
	paginations "github.com/crowdeco/skeleton/paginations"
	routes "github.com/crowdeco/skeleton/routes"
	services "github.com/crowdeco/skeleton/services"
	utils "github.com/crowdeco/skeleton/utils"
	"github.com/fatih/color"
	"github.com/gadelkareem/cachita"
	"github.com/gertd/go-pluralize"
	"github.com/joho/godotenv"
	elastic "github.com/olivere/elastic/v7"
	"github.com/sarulabs/dingo/v4"
	"github.com/sirupsen/logrus"
	mongodb "github.com/weekface/mgorus"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var Core = []dingo.Def{
	{
		Name:  "core:config:parser",
		Build: (*configs.Config)(nil),
	},
	{
		Name:  "core:config:user",
		Build: (*configs.User)(nil),
	},
	{
		Name: "core:config:type",
		Build: func() (*configs.Type, error) {
			return &configs.Type{
				Map: map[string]string{
					"double":   "float64",
					"float":    "float32",
					"int32":    "int32",
					"int64":    "int64",
					"uint32":   "uint32",
					"uint64":   "uint64",
					"sint32":   "int32",
					"sint64":   "int64",
					"fixed32":  "uint32",
					"fixed64":  "uint64",
					"sfixed32": "int32",
					"sfixed64": "int64",
					"bool":     "bool",
					"string":   "string",
					"bytes":    "[]byte",
				},
			}, nil
		},
	},
	{
		Name:  "core:config:template",
		Build: (*configs.Template)(nil),
	},
	{
		Name:  "core:template:module",
		Build: (*configs.ModuleTemplate)(nil),
	},
	{
		Name:  "core:template:field",
		Build: (*configs.FieldTemplate)(nil),
	},
	{
		Name: "core:config:env",
		Build: func(user *configs.User) (*configs.Env, error) {
			godotenv.Load()

			env := configs.Env{}

			env.ServiceName = os.Getenv("APP_NAME")
			env.Version = os.Getenv("APP_VERSION")
			env.ApiVersion = os.Getenv("API_VERSION")
			env.Debug, _ = strconv.ParseBool(os.Getenv("APP_DEBUG"))
			env.HtppPort, _ = strconv.Atoi(os.Getenv("APP_PORT"))
			env.RpcPort, _ = strconv.Atoi(os.Getenv("GRPC_PORT"))

			env.DbDriver = os.Getenv("DB_DRIVER")
			env.DbHost = os.Getenv("DB_HOST")
			env.DbPort, _ = strconv.Atoi(os.Getenv("DB_PORT"))
			env.DbUser = os.Getenv("DB_USER")
			env.DbPassword = os.Getenv("DB_PASSWORD")
			env.DbName = os.Getenv("DB_NAME")
			env.DbAutoMigrate, _ = strconv.ParseBool(os.Getenv("DB_AUTO_CREATE"))

			env.ElasticsearchHost = os.Getenv("ELASTICSEARCH_HOST")
			env.ElasticsearchPort, _ = strconv.Atoi(os.Getenv("ELASTICSEARCH_PORT"))
			env.ElasticsearchIndex = env.DbName

			env.MongoDbHost = os.Getenv("MONGODB_HOST")
			env.MongoDbPort, _ = strconv.Atoi(os.Getenv("MONGODB_PORT"))
			env.MongoDbName = os.Getenv("MONGODB_NAME")

			env.AmqpHost = os.Getenv("AMQP_HOST")
			env.AmqpPort, _ = strconv.Atoi(os.Getenv("AMQP_PORT"))
			env.AmqpUser = os.Getenv("AMQP_USER")
			env.AmqpPassword = os.Getenv("AMQP_PASSWORD")

			env.HeaderUserId = os.Getenv("HEADER_USER_ID")
			env.HeaderUserEmail = os.Getenv("HEADER_USER_EMAIL")
			env.HeaderUserRole = os.Getenv("HEADER_USER_ROLE")

			env.CacheLifetime, _ = strconv.Atoi(os.Getenv("CACHE_LIFETIME"))

			env.User = user

			env.TemplateLocation = generators.TEMPLATE_PATH

			return &env, nil
		},
	},
	{
		Name: "core:module:generator",
		Build: func(
			dic configs.Generator,
			model configs.Generator,
			module configs.Generator,
			proto configs.Generator,
			server configs.Generator,
			validation configs.Generator,
			env *configs.Env,
			pluralizer *pluralize.Client,
			template *configs.Template,
			word *utils.Word,
		) (*generators.Factory, error) {
			return &generators.Factory{
				Env:        env,
				Pluralizer: pluralizer,
				Template:   template,
				Word:       word,
				Generators: []configs.Generator{
					dic,
					model,
					module,
					proto,
					server,
					validation,
				},
			}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("core:generator:dic"),
			"1": dingo.Service("core:generator:model"),
			"2": dingo.Service("core:generator:module"),
			"3": dingo.Service("core:generator:proto"),
			"4": dingo.Service("core:generator:server"),
			"5": dingo.Service("core:generator:validation"),
		},
	},
	{
		Name:  "core:generator:dic",
		Build: (*generators.Dic)(nil),
	},
	{
		Name:  "core:generator:model",
		Build: (*generators.Model)(nil),
	},
	{
		Name:  "core:generator:module",
		Build: (*generators.Module)(nil),
		Params: dingo.Params{
			"Config": dingo.Service("core:config:parser"),
		},
	},
	{
		Name:  "core:generator:proto",
		Build: (*generators.Proto)(nil),
	},
	{
		Name:  "core:generator:server",
		Build: (*generators.Server)(nil),
	},
	{
		Name:  "core:generator:validation",
		Build: (*generators.Validation)(nil),
	},
	{
		Name:  "core:database:driver:mysql",
		Build: (*drivers.Mysql)(nil),
	},
	{
		Name:  "core:database:driver:postgresql",
		Build: (*drivers.PostgreSql)(nil),
	},
	{
		Name: "core:connection:database",
		Build: func(
			env *configs.Env,
			mysql configs.Driver,
			postgresql configs.Driver,
		) (*gorm.DB, error) {
			var db configs.Driver

			switch env.DbDriver {
			case "mysql":
				db = mysql
			case "postgresql":
				db = postgresql
			default:
				return nil, errors.New("Unknown Database Driver")
			}

			fmt.Println("Database configured...")

			return db.Connect(
				env.DbHost,
				env.DbPort,
				env.DbUser,
				env.DbPassword,
				env.DbName,
				env.Debug,
			), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("core:config:env"),
			"1": dingo.Service("core:database:driver:mysql"),
			"2": dingo.Service("core:database:driver:postgresql"),
		},
	},
	{
		Name: "core:connection:elasticsearch",
		Build: func(env *configs.Env) (*elastic.Client, error) {
			client, err := elastic.NewClient(elastic.SetURL(fmt.Sprintf("%s:%d", env.ElasticsearchHost, env.ElasticsearchPort)), elastic.SetSniff(false), elastic.SetHealthcheck(false))
			if err != nil {
				return nil, err
			}

			fmt.Println("Elasticsearch configured...")

			return client, nil
		},
	},
	{
		Name:  "core:listener:create:elasticsearch",
		Build: (*creates.Elasticsearch)(nil),
		Params: dingo.Params{
			"Context":       dingo.Service("core:context:background"),
			"Elasticsearch": dingo.Service("core:connection:elasticsearch"),
		},
	},
	{
		Name:  "core:listener:update:elasticsearch",
		Build: (*updates.Elasticsearch)(nil),
		Params: dingo.Params{
			"Context":       dingo.Service("core:context:background"),
			"Elasticsearch": dingo.Service("core:connection:elasticsearch"),
		},
	},
	{
		Name:  "core:listener:delete:elasticsearch",
		Build: (*deletes.Elasticsearch)(nil),
		Params: dingo.Params{
			"Context":       dingo.Service("core:context:background"),
			"Elasticsearch": dingo.Service("core:connection:elasticsearch"),
		},
	},
	{
		Name:  "core:listener:create:created_by",
		Build: (*creates.CreatedBy)(nil),
		Params: dingo.Params{
			"Env": dingo.Service("core:config:env"),
		},
	},
	{
		Name:  "core:listener:update:updated_by",
		Build: (*updates.UpdatedBy)(nil),
		Params: dingo.Params{
			"Env": dingo.Service("core:config:env"),
		},
	},
	{
		Name:  "core:listener:delete:deleted_by",
		Build: (*deletes.DeletedBy)(nil),
		Params: dingo.Params{
			"Env": dingo.Service("core:config:env"),
		},
	},
	{
		Name:  "core:interface:database",
		Build: (*interfaces.Database)(nil),
	},
	{
		Name:  "core:interface:elasticsearch",
		Build: (*interfaces.Elasticsearch)(nil),
	},
	{
		Name:  "core:interface:grpc",
		Build: (*interfaces.GRpc)(nil),
		Params: dingo.Params{
			"Env":  dingo.Service("core:config:env"),
			"GRpc": dingo.Service("core:grpc:server"),
		},
	},
	{
		Name:  "core:interface:queue",
		Build: (*interfaces.Queue)(nil),
	},
	{
		Name:  "core:interface:rest",
		Build: (*interfaces.Rest)(nil),
		Params: dingo.Params{
			"Middleware": dingo.Service("core:handler:middleware"),
			"Router":     dingo.Service("core:handler:router"),
			"Server":     dingo.Service("core:http:mux"),
			"Context":    dingo.Service("core:context:background"),
		},
	},
	{
		Name: "core:handler:logger",
		Build: func(
			env *configs.Env,
			logger *logrus.Logger,
			extension *configs.LoggerExtension,
		) (*handlers.Logger, error) {
			logger.SetFormatter(&logrus.JSONFormatter{})

			for _, e := range extension.Extensions {
				logger.AddHook(e)
			}

			return &handlers.Logger{
				Env:    env,
				Logger: logger,
			}, nil
		},
	},
	{
		Name:  "core:handler:messager",
		Build: (*handlers.Messenger)(nil),
		Params: dingo.Params{
			"Logger":    dingo.Service("core:handler:logger"),
			"Publisher": dingo.Service("core:message:publisher"),
			"Consumer":  dingo.Service("core:message:consumer"),
		},
	},
	{
		Name:  "core:handler:handler",
		Build: (*handlers.Handler)(nil),
		Params: dingo.Params{
			"Context":       dingo.Service("core:context:background"),
			"Elasticsearch": dingo.Service("core:connection:elasticsearch"),
			"Dispatcher":    dingo.Service("core:event:dispatcher"),
			"Repository":    dingo.Service("core:service:repository"),
		},
	},
	{
		Name: "core:handler:router",
		Build: func(
			gateway configs.Router,
			mux configs.Router,
		) (*handlers.Router, error) {
			return &handlers.Router{
				Routes: []configs.Router{
					gateway,
					mux,
				},
			}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("core:router:gateway"),
			"1": dingo.Service("core:router:mux"),
		},
	},
	{
		Name:  "core:middleware:auth",
		Build: (*middlewares.Auth)(nil),
		Params: dingo.Params{
			"Env": dingo.Service("core:config:env"),
		},
	},
	{
		Name:  "core:router:mux",
		Build: (*routes.MuxRouter)(nil),
	},
	{
		Name:  "core:router:gateway",
		Build: (*routes.GRpcGateway)(nil),
	},
	{
		Name: "core:http:mux",
		Build: func() (*http.ServeMux, error) {
			return http.NewServeMux(), nil
		},
	},
	{
		Name: "core:grpc:server",
		Build: func() (*grpc.Server, error) {
			return grpc.NewServer(), nil
		},
	},
	{
		Name: "core:log:logger",
		Build: func() (*logrus.Logger, error) {
			return logrus.New(), nil
		},
	},
	{
		Name: "core:logger:extension:mongodb",
		Build: func(env *configs.Env) (logrus.Hook, error) {
			mongodb, err := mongodb.NewHooker(fmt.Sprintf("%s:%d", env.MongoDbHost, env.MongoDbPort), env.MongoDbName, "logs")
			if err != nil {
				return nil, err
			}

			return mongodb, nil
		},
	},
	{
		Name: "core:context:background",
		Build: func() (context.Context, error) {
			return context.Background(), nil
		},
	},
	{
		Name: "core:message:config",
		Build: func(env *configs.Env) (amqp.Config, error) {
			return amqp.NewDurableQueueConfig(fmt.Sprintf("amqp://%s:%s@%s:%d/", env.AmqpUser, env.AmqpPassword, env.AmqpHost, env.AmqpPort)), nil
		},
	},
	{
		Name: "core:message:publisher",
		Build: func(env *configs.Env, config amqp.Config) (*amqp.Publisher, error) {
			publisher, err := amqp.NewPublisher(config, watermill.NewStdLogger(env.Debug, env.Debug))
			if err != nil {
				return nil, err
			}

			return publisher, nil
		},
	},
	{
		Name: "core:message:consumer",
		Build: func(config amqp.Config) (*amqp.Subscriber, error) {
			consumer, err := amqp.NewSubscriber(config, watermill.NewStdLogger(false, false))
			if err != nil {
				return nil, err
			}

			return consumer, nil
		},
	},
	{
		Name:  "core:pagination:paginator",
		Build: (*paginations.Pagination)(nil),
	},
	{
		Name:  "core:service:repository",
		Build: (*services.Repository)(nil),
		Params: dingo.Params{
			"Env":      dingo.Service("core:config:env"),
			"Database": dingo.Service("core:connection:database"),
		},
	},
	{
		Name:  "core:cache:memory",
		Build: (*utils.Cache)(nil),
		Params: dingo.Params{
			"Pool": dingo.Service("core:cachita:cache"),
		},
	},
	{
		Name:  "core:util:number",
		Build: (*utils.Number)(nil),
	},
	{
		Name:  "core:util:word",
		Build: (*utils.Word)(nil),
	},
	{
		Name: "core:util:cli",
		Build: func() (*color.Color, error) {
			return color.New(color.FgCyan, color.Bold), nil
		},
	},
	{
		Name: "core:util:pluralizer",
		Build: func() (*pluralize.Client, error) {
			return pluralize.NewClient(), nil
		},
	},
	{
		Name:  "core:util:time",
		Build: (*utils.Time)(nil),
	},
	{
		Name: "core:cachita:cache",
		Build: func() (cachita.Cache, error) {
			return cachita.Memory(), nil
		},
	},
}
