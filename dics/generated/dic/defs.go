package dic

import (
	"errors"

	"github.com/sarulabs/di/v2"
	"github.com/sarulabs/dingo/v4"

	context "context"
	aliashttp "net/http"

	amqp "github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
	configs "github.com/crowdeco/skeleton/configs"
	driver "github.com/crowdeco/skeleton/configs/driver"
	events "github.com/crowdeco/skeleton/events"
	handlers "github.com/crowdeco/skeleton/handlers"
	interfaces "github.com/crowdeco/skeleton/interfaces"
	middlewares "github.com/crowdeco/skeleton/middlewares"
	paginations "github.com/crowdeco/skeleton/paginations"
	routes "github.com/crowdeco/skeleton/routes"
	utils "github.com/crowdeco/skeleton/utils"
	cachita "github.com/gadelkareem/cachita"
	v "github.com/olivere/elastic/v7"
	v1 "github.com/vcraescu/go-paginator/v2"
	grpc "google.golang.org/grpc"
	gorm "gorm.io/gorm"
)

func getDiDefs(provider dingo.Provider) []di.Def {
	return []di.Def{
		{
			Name:  "core:cache:memory",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				pi0, err := ctn.SafeGet("core:config:env")
				if err != nil {
					var eo *utils.Cache
					return eo, err
				}
				p0, ok := pi0.(*configs.Env)
				if !ok {
					var eo *utils.Cache
					return eo, errors.New("could not cast parameter Env to *configs.Env")
				}
				pi1, err := ctn.SafeGet("core:cachita:cache")
				if err != nil {
					var eo *utils.Cache
					return eo, err
				}
				p1, ok := pi1.(cachita.Cache)
				if !ok {
					var eo *utils.Cache
					return eo, errors.New("could not cast parameter Pool to cachita.Cache")
				}
				return &utils.Cache{
					Env:  p0,
					Pool: p1,
				}, nil
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "core:cachita:cache",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("core:cachita:cache")
				if err != nil {
					var eo cachita.Cache
					return eo, err
				}
				b, ok := d.Build.(func() (cachita.Cache, error))
				if !ok {
					var eo cachita.Cache
					return eo, errors.New("could not cast build function to func() (cachita.Cache, error)")
				}
				return b()
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "core:config:env",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("core:config:env")
				if err != nil {
					var eo *configs.Env
					return eo, err
				}
				pi0, err := ctn.SafeGet("core:config:user")
				if err != nil {
					var eo *configs.Env
					return eo, err
				}
				p0, ok := pi0.(*configs.User)
				if !ok {
					var eo *configs.Env
					return eo, errors.New("could not cast parameter 0 to *configs.User")
				}
				b, ok := d.Build.(func(*configs.User) (*configs.Env, error))
				if !ok {
					var eo *configs.Env
					return eo, errors.New("could not cast build function to func(*configs.User) (*configs.Env, error)")
				}
				return b(p0)
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "core:config:user",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				var p1 string
				var p0 string
				var p2 string
				return &configs.User{
					ID:    p0,
					Email: p1,
					Role:  p2,
				}, nil
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "core:connection:database",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("core:connection:database")
				if err != nil {
					var eo *gorm.DB
					return eo, err
				}
				pi0, err := ctn.SafeGet("core:config:env")
				if err != nil {
					var eo *gorm.DB
					return eo, err
				}
				p0, ok := pi0.(*configs.Env)
				if !ok {
					var eo *gorm.DB
					return eo, errors.New("could not cast parameter 0 to *configs.Env")
				}
				pi1, err := ctn.SafeGet("core:database:driver:mysql")
				if err != nil {
					var eo *gorm.DB
					return eo, err
				}
				p1, ok := pi1.(driver.Driver)
				if !ok {
					var eo *gorm.DB
					return eo, errors.New("could not cast parameter 1 to driver.Driver")
				}
				pi2, err := ctn.SafeGet("core:database:driver:postgresql")
				if err != nil {
					var eo *gorm.DB
					return eo, err
				}
				p2, ok := pi2.(driver.Driver)
				if !ok {
					var eo *gorm.DB
					return eo, errors.New("could not cast parameter 2 to driver.Driver")
				}
				b, ok := d.Build.(func(*configs.Env, driver.Driver, driver.Driver) (*gorm.DB, error))
				if !ok {
					var eo *gorm.DB
					return eo, errors.New("could not cast build function to func(*configs.Env, driver.Driver, driver.Driver) (*gorm.DB, error)")
				}
				return b(p0, p1, p2)
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "core:connection:elasticsearch",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("core:connection:elasticsearch")
				if err != nil {
					var eo *v.Client
					return eo, err
				}
				pi0, err := ctn.SafeGet("core:config:env")
				if err != nil {
					var eo *v.Client
					return eo, err
				}
				p0, ok := pi0.(*configs.Env)
				if !ok {
					var eo *v.Client
					return eo, errors.New("could not cast parameter 0 to *configs.Env")
				}
				b, ok := d.Build.(func(*configs.Env) (*v.Client, error))
				if !ok {
					var eo *v.Client
					return eo, errors.New("could not cast build function to func(*configs.Env) (*v.Client, error)")
				}
				return b(p0)
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "core:context:background",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("core:context:background")
				if err != nil {
					var eo context.Context
					return eo, err
				}
				b, ok := d.Build.(func() (context.Context, error))
				if !ok {
					var eo context.Context
					return eo, errors.New("could not cast build function to func() (context.Context, error)")
				}
				return b()
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "core:database:driver:mysql",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				return &driver.Mysql{}, nil
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "core:database:driver:postgresql",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				return &driver.PostgreSql{}, nil
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "core:event:dispatcher",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("core:event:dispatcher")
				if err != nil {
					var eo *events.Dispatcher
					return eo, err
				}
				b, ok := d.Build.(func() (*events.Dispatcher, error))
				if !ok {
					var eo *events.Dispatcher
					return eo, errors.New("could not cast build function to func() (*events.Dispatcher, error)")
				}
				return b()
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "core:grpc:server",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("core:grpc:server")
				if err != nil {
					var eo *grpc.Server
					return eo, err
				}
				b, ok := d.Build.(func() (*grpc.Server, error))
				if !ok {
					var eo *grpc.Server
					return eo, errors.New("could not cast build function to func() (*grpc.Server, error)")
				}
				return b()
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "core:handler:handler",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				pi0, err := ctn.SafeGet("core:context:background")
				if err != nil {
					var eo *handlers.Handler
					return eo, err
				}
				p0, ok := pi0.(context.Context)
				if !ok {
					var eo *handlers.Handler
					return eo, errors.New("could not cast parameter Context to context.Context")
				}
				pi2, err := ctn.SafeGet("core:event:dispatcher")
				if err != nil {
					var eo *handlers.Handler
					return eo, err
				}
				p2, ok := pi2.(*events.Dispatcher)
				if !ok {
					var eo *handlers.Handler
					return eo, errors.New("could not cast parameter Dispatcher to *events.Dispatcher")
				}
				pi1, err := ctn.SafeGet("core:connection:elasticsearch")
				if err != nil {
					var eo *handlers.Handler
					return eo, err
				}
				p1, ok := pi1.(*v.Client)
				if !ok {
					var eo *handlers.Handler
					return eo, errors.New("could not cast parameter Elasticsearch to *v.Client")
				}
				var p3 configs.Service
				return &handlers.Handler{
					Elasticsearch: p1,
					Dispatcher:    p2,
					Service:       p3,
					Context:       p0,
				}, nil
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "core:handler:logger",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("core:handler:logger")
				if err != nil {
					var eo *handlers.Logger
					return eo, err
				}
				pi0, err := ctn.SafeGet("core:config:env")
				if err != nil {
					var eo *handlers.Logger
					return eo, err
				}
				p0, ok := pi0.(*configs.Env)
				if !ok {
					var eo *handlers.Logger
					return eo, errors.New("could not cast parameter 0 to *configs.Env")
				}
				b, ok := d.Build.(func(*configs.Env) (*handlers.Logger, error))
				if !ok {
					var eo *handlers.Logger
					return eo, errors.New("could not cast build function to func(*configs.Env) (*handlers.Logger, error)")
				}
				return b(p0)
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "core:handler:messager",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				pi1, err := ctn.SafeGet("core:message:consumer")
				if err != nil {
					var eo *handlers.Messenger
					return eo, err
				}
				p1, ok := pi1.(*amqp.Subscriber)
				if !ok {
					var eo *handlers.Messenger
					return eo, errors.New("could not cast parameter Consumer to *amqp.Subscriber")
				}
				pi2, err := ctn.SafeGet("core:handler:logger")
				if err != nil {
					var eo *handlers.Messenger
					return eo, err
				}
				p2, ok := pi2.(*handlers.Logger)
				if !ok {
					var eo *handlers.Messenger
					return eo, errors.New("could not cast parameter Logger to *handlers.Logger")
				}
				pi0, err := ctn.SafeGet("core:message:publisher")
				if err != nil {
					var eo *handlers.Messenger
					return eo, err
				}
				p0, ok := pi0.(*amqp.Publisher)
				if !ok {
					var eo *handlers.Messenger
					return eo, errors.New("could not cast parameter Publisher to *amqp.Publisher")
				}
				return &handlers.Messenger{
					Publisher: p0,
					Consumer:  p1,
					Logger:    p2,
				}, nil
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "core:handler:middleware",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("core:handler:middleware")
				if err != nil {
					var eo *handlers.Middleware
					return eo, err
				}
				pi0, err := ctn.SafeGet("core:middleware:auth")
				if err != nil {
					var eo *handlers.Middleware
					return eo, err
				}
				p0, ok := pi0.(configs.Middleware)
				if !ok {
					var eo *handlers.Middleware
					return eo, errors.New("could not cast parameter 0 to configs.Middleware")
				}
				b, ok := d.Build.(func(configs.Middleware) (*handlers.Middleware, error))
				if !ok {
					var eo *handlers.Middleware
					return eo, errors.New("could not cast build function to func(configs.Middleware) (*handlers.Middleware, error)")
				}
				return b(p0)
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "core:handler:router",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("core:handler:router")
				if err != nil {
					var eo *handlers.Router
					return eo, err
				}
				pi0, err := ctn.SafeGet("core:router:gateway")
				if err != nil {
					var eo *handlers.Router
					return eo, err
				}
				p0, ok := pi0.(configs.Router)
				if !ok {
					var eo *handlers.Router
					return eo, errors.New("could not cast parameter 0 to configs.Router")
				}
				pi1, err := ctn.SafeGet("core:router:mux")
				if err != nil {
					var eo *handlers.Router
					return eo, err
				}
				p1, ok := pi1.(configs.Router)
				if !ok {
					var eo *handlers.Router
					return eo, errors.New("could not cast parameter 1 to configs.Router")
				}
				b, ok := d.Build.(func(configs.Router, configs.Router) (*handlers.Router, error))
				if !ok {
					var eo *handlers.Router
					return eo, errors.New("could not cast build function to func(configs.Router, configs.Router) (*handlers.Router, error)")
				}
				return b(p0, p1)
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "core:http:mux",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("core:http:mux")
				if err != nil {
					var eo *aliashttp.ServeMux
					return eo, err
				}
				b, ok := d.Build.(func() (*aliashttp.ServeMux, error))
				if !ok {
					var eo *aliashttp.ServeMux
					return eo, errors.New("could not cast build function to func() (*aliashttp.ServeMux, error)")
				}
				return b()
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "core:interface:database",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("core:interface:database")
				if err != nil {
					var eo *interfaces.Database
					return eo, err
				}
				b, ok := d.Build.(func() (*interfaces.Database, error))
				if !ok {
					var eo *interfaces.Database
					return eo, errors.New("could not cast build function to func() (*interfaces.Database, error)")
				}
				return b()
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "core:interface:grpc",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("core:interface:grpc")
				if err != nil {
					var eo *interfaces.GRpc
					return eo, err
				}
				pi0, err := ctn.SafeGet("core:config:env")
				if err != nil {
					var eo *interfaces.GRpc
					return eo, err
				}
				p0, ok := pi0.(*configs.Env)
				if !ok {
					var eo *interfaces.GRpc
					return eo, errors.New("could not cast parameter 0 to *configs.Env")
				}
				pi1, err := ctn.SafeGet("core:grpc:server")
				if err != nil {
					var eo *interfaces.GRpc
					return eo, err
				}
				p1, ok := pi1.(*grpc.Server)
				if !ok {
					var eo *interfaces.GRpc
					return eo, errors.New("could not cast parameter 1 to *grpc.Server")
				}
				pi2, err := ctn.SafeGet("core:event:dispatcher")
				if err != nil {
					var eo *interfaces.GRpc
					return eo, err
				}
				p2, ok := pi2.(*events.Dispatcher)
				if !ok {
					var eo *interfaces.GRpc
					return eo, errors.New("could not cast parameter 2 to *events.Dispatcher")
				}
				b, ok := d.Build.(func(*configs.Env, *grpc.Server, *events.Dispatcher) (*interfaces.GRpc, error))
				if !ok {
					var eo *interfaces.GRpc
					return eo, errors.New("could not cast build function to func(*configs.Env, *grpc.Server, *events.Dispatcher) (*interfaces.GRpc, error)")
				}
				return b(p0, p1, p2)
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "core:interface:queue",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("core:interface:queue")
				if err != nil {
					var eo *interfaces.Queue
					return eo, err
				}
				b, ok := d.Build.(func() (*interfaces.Queue, error))
				if !ok {
					var eo *interfaces.Queue
					return eo, errors.New("could not cast build function to func() (*interfaces.Queue, error)")
				}
				return b()
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "core:interface:rest",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				pi4, err := ctn.SafeGet("core:context:background")
				if err != nil {
					var eo *interfaces.Rest
					return eo, err
				}
				p4, ok := pi4.(context.Context)
				if !ok {
					var eo *interfaces.Rest
					return eo, errors.New("could not cast parameter Context to context.Context")
				}
				pi0, err := ctn.SafeGet("core:config:env")
				if err != nil {
					var eo *interfaces.Rest
					return eo, err
				}
				p0, ok := pi0.(*configs.Env)
				if !ok {
					var eo *interfaces.Rest
					return eo, errors.New("could not cast parameter Env to *configs.Env")
				}
				pi1, err := ctn.SafeGet("core:handler:middleware")
				if err != nil {
					var eo *interfaces.Rest
					return eo, err
				}
				p1, ok := pi1.(*handlers.Middleware)
				if !ok {
					var eo *interfaces.Rest
					return eo, errors.New("could not cast parameter Middleware to *handlers.Middleware")
				}
				pi2, err := ctn.SafeGet("core:handler:router")
				if err != nil {
					var eo *interfaces.Rest
					return eo, err
				}
				p2, ok := pi2.(*handlers.Router)
				if !ok {
					var eo *interfaces.Rest
					return eo, errors.New("could not cast parameter Router to *handlers.Router")
				}
				pi3, err := ctn.SafeGet("core:http:mux")
				if err != nil {
					var eo *interfaces.Rest
					return eo, err
				}
				p3, ok := pi3.(*aliashttp.ServeMux)
				if !ok {
					var eo *interfaces.Rest
					return eo, errors.New("could not cast parameter Server to *aliashttp.ServeMux")
				}
				return &interfaces.Rest{
					Server:     p3,
					Context:    p4,
					Env:        p0,
					Middleware: p1,
					Router:     p2,
				}, nil
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "core:message:config",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("core:message:config")
				if err != nil {
					var eo amqp.Config
					return eo, err
				}
				pi0, err := ctn.SafeGet("core:config:env")
				if err != nil {
					var eo amqp.Config
					return eo, err
				}
				p0, ok := pi0.(*configs.Env)
				if !ok {
					var eo amqp.Config
					return eo, errors.New("could not cast parameter 0 to *configs.Env")
				}
				b, ok := d.Build.(func(*configs.Env) (amqp.Config, error))
				if !ok {
					var eo amqp.Config
					return eo, errors.New("could not cast build function to func(*configs.Env) (amqp.Config, error)")
				}
				return b(p0)
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "core:message:consumer",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("core:message:consumer")
				if err != nil {
					var eo *amqp.Subscriber
					return eo, err
				}
				pi0, err := ctn.SafeGet("core:message:config")
				if err != nil {
					var eo *amqp.Subscriber
					return eo, err
				}
				p0, ok := pi0.(amqp.Config)
				if !ok {
					var eo *amqp.Subscriber
					return eo, errors.New("could not cast parameter 0 to amqp.Config")
				}
				b, ok := d.Build.(func(amqp.Config) (*amqp.Subscriber, error))
				if !ok {
					var eo *amqp.Subscriber
					return eo, errors.New("could not cast build function to func(amqp.Config) (*amqp.Subscriber, error)")
				}
				return b(p0)
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "core:message:publisher",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("core:message:publisher")
				if err != nil {
					var eo *amqp.Publisher
					return eo, err
				}
				pi0, err := ctn.SafeGet("core:config:env")
				if err != nil {
					var eo *amqp.Publisher
					return eo, err
				}
				p0, ok := pi0.(*configs.Env)
				if !ok {
					var eo *amqp.Publisher
					return eo, errors.New("could not cast parameter 0 to *configs.Env")
				}
				pi1, err := ctn.SafeGet("core:message:config")
				if err != nil {
					var eo *amqp.Publisher
					return eo, err
				}
				p1, ok := pi1.(amqp.Config)
				if !ok {
					var eo *amqp.Publisher
					return eo, errors.New("could not cast parameter 1 to amqp.Config")
				}
				b, ok := d.Build.(func(*configs.Env, amqp.Config) (*amqp.Publisher, error))
				if !ok {
					var eo *amqp.Publisher
					return eo, errors.New("could not cast build function to func(*configs.Env, amqp.Config) (*amqp.Publisher, error)")
				}
				return b(p0, p1)
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "core:middleware:auth",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				pi0, err := ctn.SafeGet("core:config:env")
				if err != nil {
					var eo *middlewares.Auth
					return eo, err
				}
				p0, ok := pi0.(*configs.Env)
				if !ok {
					var eo *middlewares.Auth
					return eo, errors.New("could not cast parameter Env to *configs.Env")
				}
				return &middlewares.Auth{
					Env: p0,
				}, nil
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "core:pagination:paginator",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				var p2 []paginations.Filter
				var p0 int
				var p1 int
				var p4 v1.Paginator
				var p3 string
				return &paginations.Pagination{
					Limit:   p0,
					Page:    p1,
					Filters: p2,
					Search:  p3,
					Pager:   p4,
				}, nil
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "core:router:gateway",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				return &routes.GRpcGateway{}, nil
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "core:router:mux",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				return &routes.MuxRouter{}, nil
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "core:util:number",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				return &utils.Number{}, nil
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "core:util:time",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				return &utils.Time{}, nil
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "core:util:word",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				return &utils.Word{}, nil
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
	}
}
