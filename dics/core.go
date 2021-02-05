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
	handlers "github.com/crowdeco/skeleton/handlers"
	interfaces "github.com/crowdeco/skeleton/interfaces"
	middlewares "github.com/crowdeco/skeleton/middlewares"
	paginations "github.com/crowdeco/skeleton/paginations"
	routes "github.com/crowdeco/skeleton/routes"
	utils "github.com/crowdeco/skeleton/utils"
	"github.com/gadelkareem/cachita"
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
		Name:  "core:config:user",
		Build: (*configs.User)(nil),
	},
	{
		Name: "core:config:env",
		Build: func(user *configs.User) (*configs.Env, error) {
			godotenv.Load()

			env := configs.Env{}

			env.ServiceName = os.Getenv("APP_NAME")
			env.Version = os.Getenv("APP_VERSION")
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

			return &env, nil
		},
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
		Build: func(env *configs.Env) (*handlers.Logger, error) {
			logger := logrus.New()
			logger.SetFormatter(&logrus.JSONFormatter{})

			mongodb, err := mongodb.NewHooker(fmt.Sprintf("%s:%d", env.MongoDbHost, env.MongoDbPort), env.MongoDbName, "logs")
			if err == nil {
				logger.AddHook(mongodb)
			} else {
				return nil, err
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
		},
	},
	{
		Name: "core:handler:middleware",
		Build: func(
			auth configs.Middleware,
		) (*handlers.Middleware, error) {
			return &handlers.Middleware{
				Middlewares: []configs.Middleware{
					auth,
				},
			}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("core:middleware:auth"),
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
		Name:  "core:router:gateway",
		Build: (*routes.GRpcGateway)(nil),
	},
	{
		Name:  "core:router:mux",
		Build: (*routes.MuxRouter)(nil),
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
		Name: "core:context:background",
		Build: func() (context.Context, error) {
			return context.Background(), nil
		},
	},
	{
		Name: "core:message:config",
		Build: func(env *configs.Env) (amqp.Config, error) {
			address := fmt.Sprintf("amqp://%s:%s@%s:%d/", env.AmqpUser, env.AmqpPassword, env.AmqpHost, env.AmqpPort)

			return amqp.NewDurableQueueConfig(address), nil
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
