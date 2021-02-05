package dic

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/sarulabs/di/v2"
	"github.com/sarulabs/dingo/v4"

	providerPkg "github.com/crowdeco/skeleton/dics"

	context "context"
	aliashttp "net/http"

	amqp "github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
	configs "github.com/crowdeco/skeleton/configs"
	drivers "github.com/crowdeco/skeleton/configs/drivers"
	events "github.com/crowdeco/skeleton/events"
	handlers "github.com/crowdeco/skeleton/handlers"
	interfaces "github.com/crowdeco/skeleton/interfaces"
	middlewares "github.com/crowdeco/skeleton/middlewares"
	paginations "github.com/crowdeco/skeleton/paginations"
	routes "github.com/crowdeco/skeleton/routes"
	utils "github.com/crowdeco/skeleton/utils"
	cachita "github.com/gadelkareem/cachita"
	v "github.com/olivere/elastic/v7"
	grpc "google.golang.org/grpc"
	gorm "gorm.io/gorm"
)

// C retrieves a Container from an interface.
// The function panics if the Container can not be retrieved.
//
// The interface can be :
// - a *Container
// - an *http.Request containing a *Container in its context.Context
//   for the dingo.ContainerKey("dingo") key.
//
// The function can be changed to match the needs of your application.
var C = func(i interface{}) *Container {
	if c, ok := i.(*Container); ok {
		return c
	}
	r, ok := i.(*http.Request)
	if !ok {
		panic("could not get the container with dic.C()")
	}
	c, ok := r.Context().Value(dingo.ContainerKey("dingo")).(*Container)
	if !ok {
		panic("could not get the container from the given *http.Request in dic.C()")
	}
	return c
}

type builder struct {
	builder *di.Builder
}

// NewBuilder creates a builder that can create a Container.
// You should you NewContainer to create the container directly.
// Using NewBuilder allows you to redefine some di services though.
// This could be used for testing.
// But this behaviour is not safe, so be sure to know what you are doing.
func NewBuilder(scopes ...string) (*builder, error) {
	if len(scopes) == 0 {
		scopes = []string{di.App, di.Request, di.SubRequest}
	}
	b, err := di.NewBuilder(scopes...)
	if err != nil {
		return nil, fmt.Errorf("could not create di.Builder: %v", err)
	}
	provider := &providerPkg.Provider{}
	if err := provider.Load(); err != nil {
		return nil, fmt.Errorf("could not load definitions with the Provider (Provider from github.com/crowdeco/skeleton/dics): %v", err)
	}
	for _, d := range getDiDefs(provider) {
		if err := b.Add(d); err != nil {
			return nil, fmt.Errorf("could not add di.Def in di.Builder: %v", err)
		}
	}
	return &builder{builder: b}, nil
}

// Add adds one or more definitions in the Builder.
// It returns an error if a definition can not be added.
func (b *builder) Add(defs ...di.Def) error {
	return b.builder.Add(defs...)
}

// Set is a shortcut to add a definition for an already built object.
func (b *builder) Set(name string, obj interface{}) error {
	return b.builder.Set(name, obj)
}

// Build creates a Container in the most generic scope.
func (b *builder) Build() *Container {
	return &Container{ctn: b.builder.Build()}
}

// NewContainer creates a new Container.
// If no scope is provided, di.App, di.Request and di.SubRequest are used.
// The returned Container has the most generic scope (di.App).
// The SubContainer() method should be called to get a Container in a more specific scope.
func NewContainer(scopes ...string) (*Container, error) {
	b, err := NewBuilder(scopes...)
	if err != nil {
		return nil, err
	}
	return b.Build(), nil
}

// Container represents a generated dependency injection container.
// It is a wrapper around a di.Container.
//
// A Container has a scope and may have a parent in a more generic scope
// and children in a more specific scope.
// Objects can be retrieved from the Container.
// If the requested object does not already exist in the Container,
// it is built thanks to the object definition.
// The following attempts to get this object will return the same object.
type Container struct {
	ctn di.Container
}

// Scope returns the Container scope.
func (c *Container) Scope() string {
	return c.ctn.Scope()
}

// Scopes returns the list of available scopes.
func (c *Container) Scopes() []string {
	return c.ctn.Scopes()
}

// ParentScopes returns the list of scopes wider than the Container scope.
func (c *Container) ParentScopes() []string {
	return c.ctn.ParentScopes()
}

// SubScopes returns the list of scopes that are more specific than the Container scope.
func (c *Container) SubScopes() []string {
	return c.ctn.SubScopes()
}

// Parent returns the parent Container.
func (c *Container) Parent() *Container {
	if p := c.ctn.Parent(); p != nil {
		return &Container{ctn: p}
	}
	return nil
}

// SubContainer creates a new Container in the next sub-scope
// that will have this Container as parent.
func (c *Container) SubContainer() (*Container, error) {
	sub, err := c.ctn.SubContainer()
	if err != nil {
		return nil, err
	}
	return &Container{ctn: sub}, nil
}

// SafeGet retrieves an object from the Container.
// The object has to belong to this scope or a more generic one.
// If the object does not already exist, it is created and saved in the Container.
// If the object can not be created, it returns an error.
func (c *Container) SafeGet(name string) (interface{}, error) {
	return c.ctn.SafeGet(name)
}

// Get is similar to SafeGet but it does not return the error.
// Instead it panics.
func (c *Container) Get(name string) interface{} {
	return c.ctn.Get(name)
}

// Fill is similar to SafeGet but it does not return the object.
// Instead it fills the provided object with the value returned by SafeGet.
// The provided object must be a pointer to the value returned by SafeGet.
func (c *Container) Fill(name string, dst interface{}) error {
	return c.ctn.Fill(name, dst)
}

// UnscopedSafeGet retrieves an object from the Container, like SafeGet.
// The difference is that the object can be retrieved
// even if it belongs to a more specific scope.
// To do so, UnscopedSafeGet creates a sub-container.
// When the created object is no longer needed,
// it is important to use the Clean method to delete this sub-container.
func (c *Container) UnscopedSafeGet(name string) (interface{}, error) {
	return c.ctn.UnscopedSafeGet(name)
}

// UnscopedGet is similar to UnscopedSafeGet but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGet(name string) interface{} {
	return c.ctn.UnscopedGet(name)
}

// UnscopedFill is similar to UnscopedSafeGet but copies the object in dst instead of returning it.
func (c *Container) UnscopedFill(name string, dst interface{}) error {
	return c.ctn.UnscopedFill(name, dst)
}

// Clean deletes the sub-container created by UnscopedSafeGet, UnscopedGet or UnscopedFill.
func (c *Container) Clean() error {
	return c.ctn.Clean()
}

// DeleteWithSubContainers takes all the objects saved in this Container
// and calls the Close function of their Definition on them.
// It will also call DeleteWithSubContainers on each child and remove its reference in the parent Container.
// After deletion, the Container can no longer be used.
// The sub-containers are deleted even if they are still used in other goroutines.
// It can cause errors. You may want to use the Delete method instead.
func (c *Container) DeleteWithSubContainers() error {
	return c.ctn.DeleteWithSubContainers()
}

// Delete works like DeleteWithSubContainers if the Container does not have any child.
// But if the Container has sub-containers, it will not be deleted right away.
// The deletion only occurs when all the sub-containers have been deleted manually.
// So you have to call Delete or DeleteWithSubContainers on all the sub-containers.
func (c *Container) Delete() error {
	return c.ctn.Delete()
}

// IsClosed returns true if the Container has been deleted.
func (c *Container) IsClosed() bool {
	return c.ctn.IsClosed()
}

// SafeGetCoreCacheMemory works like SafeGet but only for CoreCacheMemory.
// It does not return an interface but a *utils.Cache.
func (c *Container) SafeGetCoreCacheMemory() (*utils.Cache, error) {
	i, err := c.ctn.SafeGet("core:cache:memory")
	if err != nil {
		var eo *utils.Cache
		return eo, err
	}
	o, ok := i.(*utils.Cache)
	if !ok {
		return o, errors.New("could get 'core:cache:memory' because the object could not be cast to *utils.Cache")
	}
	return o, nil
}

// GetCoreCacheMemory is similar to SafeGetCoreCacheMemory but it does not return the error.
// Instead it panics.
func (c *Container) GetCoreCacheMemory() *utils.Cache {
	o, err := c.SafeGetCoreCacheMemory()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetCoreCacheMemory works like UnscopedSafeGet but only for CoreCacheMemory.
// It does not return an interface but a *utils.Cache.
func (c *Container) UnscopedSafeGetCoreCacheMemory() (*utils.Cache, error) {
	i, err := c.ctn.UnscopedSafeGet("core:cache:memory")
	if err != nil {
		var eo *utils.Cache
		return eo, err
	}
	o, ok := i.(*utils.Cache)
	if !ok {
		return o, errors.New("could get 'core:cache:memory' because the object could not be cast to *utils.Cache")
	}
	return o, nil
}

// UnscopedGetCoreCacheMemory is similar to UnscopedSafeGetCoreCacheMemory but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetCoreCacheMemory() *utils.Cache {
	o, err := c.UnscopedSafeGetCoreCacheMemory()
	if err != nil {
		panic(err)
	}
	return o
}

// CoreCacheMemory is similar to GetCoreCacheMemory.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetCoreCacheMemory method.
// If the container can not be retrieved, it panics.
func CoreCacheMemory(i interface{}) *utils.Cache {
	return C(i).GetCoreCacheMemory()
}

// SafeGetCoreCachitaCache works like SafeGet but only for CoreCachitaCache.
// It does not return an interface but a cachita.Cache.
func (c *Container) SafeGetCoreCachitaCache() (cachita.Cache, error) {
	i, err := c.ctn.SafeGet("core:cachita:cache")
	if err != nil {
		var eo cachita.Cache
		return eo, err
	}
	o, ok := i.(cachita.Cache)
	if !ok {
		return o, errors.New("could get 'core:cachita:cache' because the object could not be cast to cachita.Cache")
	}
	return o, nil
}

// GetCoreCachitaCache is similar to SafeGetCoreCachitaCache but it does not return the error.
// Instead it panics.
func (c *Container) GetCoreCachitaCache() cachita.Cache {
	o, err := c.SafeGetCoreCachitaCache()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetCoreCachitaCache works like UnscopedSafeGet but only for CoreCachitaCache.
// It does not return an interface but a cachita.Cache.
func (c *Container) UnscopedSafeGetCoreCachitaCache() (cachita.Cache, error) {
	i, err := c.ctn.UnscopedSafeGet("core:cachita:cache")
	if err != nil {
		var eo cachita.Cache
		return eo, err
	}
	o, ok := i.(cachita.Cache)
	if !ok {
		return o, errors.New("could get 'core:cachita:cache' because the object could not be cast to cachita.Cache")
	}
	return o, nil
}

// UnscopedGetCoreCachitaCache is similar to UnscopedSafeGetCoreCachitaCache but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetCoreCachitaCache() cachita.Cache {
	o, err := c.UnscopedSafeGetCoreCachitaCache()
	if err != nil {
		panic(err)
	}
	return o
}

// CoreCachitaCache is similar to GetCoreCachitaCache.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetCoreCachitaCache method.
// If the container can not be retrieved, it panics.
func CoreCachitaCache(i interface{}) cachita.Cache {
	return C(i).GetCoreCachitaCache()
}

// SafeGetCoreConfigEnv works like SafeGet but only for CoreConfigEnv.
// It does not return an interface but a *configs.Env.
func (c *Container) SafeGetCoreConfigEnv() (*configs.Env, error) {
	i, err := c.ctn.SafeGet("core:config:env")
	if err != nil {
		var eo *configs.Env
		return eo, err
	}
	o, ok := i.(*configs.Env)
	if !ok {
		return o, errors.New("could get 'core:config:env' because the object could not be cast to *configs.Env")
	}
	return o, nil
}

// GetCoreConfigEnv is similar to SafeGetCoreConfigEnv but it does not return the error.
// Instead it panics.
func (c *Container) GetCoreConfigEnv() *configs.Env {
	o, err := c.SafeGetCoreConfigEnv()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetCoreConfigEnv works like UnscopedSafeGet but only for CoreConfigEnv.
// It does not return an interface but a *configs.Env.
func (c *Container) UnscopedSafeGetCoreConfigEnv() (*configs.Env, error) {
	i, err := c.ctn.UnscopedSafeGet("core:config:env")
	if err != nil {
		var eo *configs.Env
		return eo, err
	}
	o, ok := i.(*configs.Env)
	if !ok {
		return o, errors.New("could get 'core:config:env' because the object could not be cast to *configs.Env")
	}
	return o, nil
}

// UnscopedGetCoreConfigEnv is similar to UnscopedSafeGetCoreConfigEnv but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetCoreConfigEnv() *configs.Env {
	o, err := c.UnscopedSafeGetCoreConfigEnv()
	if err != nil {
		panic(err)
	}
	return o
}

// CoreConfigEnv is similar to GetCoreConfigEnv.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetCoreConfigEnv method.
// If the container can not be retrieved, it panics.
func CoreConfigEnv(i interface{}) *configs.Env {
	return C(i).GetCoreConfigEnv()
}

// SafeGetCoreConfigUser works like SafeGet but only for CoreConfigUser.
// It does not return an interface but a *configs.User.
func (c *Container) SafeGetCoreConfigUser() (*configs.User, error) {
	i, err := c.ctn.SafeGet("core:config:user")
	if err != nil {
		var eo *configs.User
		return eo, err
	}
	o, ok := i.(*configs.User)
	if !ok {
		return o, errors.New("could get 'core:config:user' because the object could not be cast to *configs.User")
	}
	return o, nil
}

// GetCoreConfigUser is similar to SafeGetCoreConfigUser but it does not return the error.
// Instead it panics.
func (c *Container) GetCoreConfigUser() *configs.User {
	o, err := c.SafeGetCoreConfigUser()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetCoreConfigUser works like UnscopedSafeGet but only for CoreConfigUser.
// It does not return an interface but a *configs.User.
func (c *Container) UnscopedSafeGetCoreConfigUser() (*configs.User, error) {
	i, err := c.ctn.UnscopedSafeGet("core:config:user")
	if err != nil {
		var eo *configs.User
		return eo, err
	}
	o, ok := i.(*configs.User)
	if !ok {
		return o, errors.New("could get 'core:config:user' because the object could not be cast to *configs.User")
	}
	return o, nil
}

// UnscopedGetCoreConfigUser is similar to UnscopedSafeGetCoreConfigUser but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetCoreConfigUser() *configs.User {
	o, err := c.UnscopedSafeGetCoreConfigUser()
	if err != nil {
		panic(err)
	}
	return o
}

// CoreConfigUser is similar to GetCoreConfigUser.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetCoreConfigUser method.
// If the container can not be retrieved, it panics.
func CoreConfigUser(i interface{}) *configs.User {
	return C(i).GetCoreConfigUser()
}

// SafeGetCoreConnectionDatabase works like SafeGet but only for CoreConnectionDatabase.
// It does not return an interface but a *gorm.DB.
func (c *Container) SafeGetCoreConnectionDatabase() (*gorm.DB, error) {
	i, err := c.ctn.SafeGet("core:connection:database")
	if err != nil {
		var eo *gorm.DB
		return eo, err
	}
	o, ok := i.(*gorm.DB)
	if !ok {
		return o, errors.New("could get 'core:connection:database' because the object could not be cast to *gorm.DB")
	}
	return o, nil
}

// GetCoreConnectionDatabase is similar to SafeGetCoreConnectionDatabase but it does not return the error.
// Instead it panics.
func (c *Container) GetCoreConnectionDatabase() *gorm.DB {
	o, err := c.SafeGetCoreConnectionDatabase()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetCoreConnectionDatabase works like UnscopedSafeGet but only for CoreConnectionDatabase.
// It does not return an interface but a *gorm.DB.
func (c *Container) UnscopedSafeGetCoreConnectionDatabase() (*gorm.DB, error) {
	i, err := c.ctn.UnscopedSafeGet("core:connection:database")
	if err != nil {
		var eo *gorm.DB
		return eo, err
	}
	o, ok := i.(*gorm.DB)
	if !ok {
		return o, errors.New("could get 'core:connection:database' because the object could not be cast to *gorm.DB")
	}
	return o, nil
}

// UnscopedGetCoreConnectionDatabase is similar to UnscopedSafeGetCoreConnectionDatabase but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetCoreConnectionDatabase() *gorm.DB {
	o, err := c.UnscopedSafeGetCoreConnectionDatabase()
	if err != nil {
		panic(err)
	}
	return o
}

// CoreConnectionDatabase is similar to GetCoreConnectionDatabase.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetCoreConnectionDatabase method.
// If the container can not be retrieved, it panics.
func CoreConnectionDatabase(i interface{}) *gorm.DB {
	return C(i).GetCoreConnectionDatabase()
}

// SafeGetCoreConnectionElasticsearch works like SafeGet but only for CoreConnectionElasticsearch.
// It does not return an interface but a *v.Client.
func (c *Container) SafeGetCoreConnectionElasticsearch() (*v.Client, error) {
	i, err := c.ctn.SafeGet("core:connection:elasticsearch")
	if err != nil {
		var eo *v.Client
		return eo, err
	}
	o, ok := i.(*v.Client)
	if !ok {
		return o, errors.New("could get 'core:connection:elasticsearch' because the object could not be cast to *v.Client")
	}
	return o, nil
}

// GetCoreConnectionElasticsearch is similar to SafeGetCoreConnectionElasticsearch but it does not return the error.
// Instead it panics.
func (c *Container) GetCoreConnectionElasticsearch() *v.Client {
	o, err := c.SafeGetCoreConnectionElasticsearch()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetCoreConnectionElasticsearch works like UnscopedSafeGet but only for CoreConnectionElasticsearch.
// It does not return an interface but a *v.Client.
func (c *Container) UnscopedSafeGetCoreConnectionElasticsearch() (*v.Client, error) {
	i, err := c.ctn.UnscopedSafeGet("core:connection:elasticsearch")
	if err != nil {
		var eo *v.Client
		return eo, err
	}
	o, ok := i.(*v.Client)
	if !ok {
		return o, errors.New("could get 'core:connection:elasticsearch' because the object could not be cast to *v.Client")
	}
	return o, nil
}

// UnscopedGetCoreConnectionElasticsearch is similar to UnscopedSafeGetCoreConnectionElasticsearch but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetCoreConnectionElasticsearch() *v.Client {
	o, err := c.UnscopedSafeGetCoreConnectionElasticsearch()
	if err != nil {
		panic(err)
	}
	return o
}

// CoreConnectionElasticsearch is similar to GetCoreConnectionElasticsearch.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetCoreConnectionElasticsearch method.
// If the container can not be retrieved, it panics.
func CoreConnectionElasticsearch(i interface{}) *v.Client {
	return C(i).GetCoreConnectionElasticsearch()
}

// SafeGetCoreContextBackground works like SafeGet but only for CoreContextBackground.
// It does not return an interface but a context.Context.
func (c *Container) SafeGetCoreContextBackground() (context.Context, error) {
	i, err := c.ctn.SafeGet("core:context:background")
	if err != nil {
		var eo context.Context
		return eo, err
	}
	o, ok := i.(context.Context)
	if !ok {
		return o, errors.New("could get 'core:context:background' because the object could not be cast to context.Context")
	}
	return o, nil
}

// GetCoreContextBackground is similar to SafeGetCoreContextBackground but it does not return the error.
// Instead it panics.
func (c *Container) GetCoreContextBackground() context.Context {
	o, err := c.SafeGetCoreContextBackground()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetCoreContextBackground works like UnscopedSafeGet but only for CoreContextBackground.
// It does not return an interface but a context.Context.
func (c *Container) UnscopedSafeGetCoreContextBackground() (context.Context, error) {
	i, err := c.ctn.UnscopedSafeGet("core:context:background")
	if err != nil {
		var eo context.Context
		return eo, err
	}
	o, ok := i.(context.Context)
	if !ok {
		return o, errors.New("could get 'core:context:background' because the object could not be cast to context.Context")
	}
	return o, nil
}

// UnscopedGetCoreContextBackground is similar to UnscopedSafeGetCoreContextBackground but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetCoreContextBackground() context.Context {
	o, err := c.UnscopedSafeGetCoreContextBackground()
	if err != nil {
		panic(err)
	}
	return o
}

// CoreContextBackground is similar to GetCoreContextBackground.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetCoreContextBackground method.
// If the container can not be retrieved, it panics.
func CoreContextBackground(i interface{}) context.Context {
	return C(i).GetCoreContextBackground()
}

// SafeGetCoreDatabaseDriverMysql works like SafeGet but only for CoreDatabaseDriverMysql.
// It does not return an interface but a *drivers.Mysql.
func (c *Container) SafeGetCoreDatabaseDriverMysql() (*drivers.Mysql, error) {
	i, err := c.ctn.SafeGet("core:database:driver:mysql")
	if err != nil {
		var eo *drivers.Mysql
		return eo, err
	}
	o, ok := i.(*drivers.Mysql)
	if !ok {
		return o, errors.New("could get 'core:database:driver:mysql' because the object could not be cast to *drivers.Mysql")
	}
	return o, nil
}

// GetCoreDatabaseDriverMysql is similar to SafeGetCoreDatabaseDriverMysql but it does not return the error.
// Instead it panics.
func (c *Container) GetCoreDatabaseDriverMysql() *drivers.Mysql {
	o, err := c.SafeGetCoreDatabaseDriverMysql()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetCoreDatabaseDriverMysql works like UnscopedSafeGet but only for CoreDatabaseDriverMysql.
// It does not return an interface but a *drivers.Mysql.
func (c *Container) UnscopedSafeGetCoreDatabaseDriverMysql() (*drivers.Mysql, error) {
	i, err := c.ctn.UnscopedSafeGet("core:database:driver:mysql")
	if err != nil {
		var eo *drivers.Mysql
		return eo, err
	}
	o, ok := i.(*drivers.Mysql)
	if !ok {
		return o, errors.New("could get 'core:database:driver:mysql' because the object could not be cast to *drivers.Mysql")
	}
	return o, nil
}

// UnscopedGetCoreDatabaseDriverMysql is similar to UnscopedSafeGetCoreDatabaseDriverMysql but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetCoreDatabaseDriverMysql() *drivers.Mysql {
	o, err := c.UnscopedSafeGetCoreDatabaseDriverMysql()
	if err != nil {
		panic(err)
	}
	return o
}

// CoreDatabaseDriverMysql is similar to GetCoreDatabaseDriverMysql.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetCoreDatabaseDriverMysql method.
// If the container can not be retrieved, it panics.
func CoreDatabaseDriverMysql(i interface{}) *drivers.Mysql {
	return C(i).GetCoreDatabaseDriverMysql()
}

// SafeGetCoreDatabaseDriverPostgresql works like SafeGet but only for CoreDatabaseDriverPostgresql.
// It does not return an interface but a *drivers.PostgreSql.
func (c *Container) SafeGetCoreDatabaseDriverPostgresql() (*drivers.PostgreSql, error) {
	i, err := c.ctn.SafeGet("core:database:driver:postgresql")
	if err != nil {
		var eo *drivers.PostgreSql
		return eo, err
	}
	o, ok := i.(*drivers.PostgreSql)
	if !ok {
		return o, errors.New("could get 'core:database:driver:postgresql' because the object could not be cast to *drivers.PostgreSql")
	}
	return o, nil
}

// GetCoreDatabaseDriverPostgresql is similar to SafeGetCoreDatabaseDriverPostgresql but it does not return the error.
// Instead it panics.
func (c *Container) GetCoreDatabaseDriverPostgresql() *drivers.PostgreSql {
	o, err := c.SafeGetCoreDatabaseDriverPostgresql()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetCoreDatabaseDriverPostgresql works like UnscopedSafeGet but only for CoreDatabaseDriverPostgresql.
// It does not return an interface but a *drivers.PostgreSql.
func (c *Container) UnscopedSafeGetCoreDatabaseDriverPostgresql() (*drivers.PostgreSql, error) {
	i, err := c.ctn.UnscopedSafeGet("core:database:driver:postgresql")
	if err != nil {
		var eo *drivers.PostgreSql
		return eo, err
	}
	o, ok := i.(*drivers.PostgreSql)
	if !ok {
		return o, errors.New("could get 'core:database:driver:postgresql' because the object could not be cast to *drivers.PostgreSql")
	}
	return o, nil
}

// UnscopedGetCoreDatabaseDriverPostgresql is similar to UnscopedSafeGetCoreDatabaseDriverPostgresql but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetCoreDatabaseDriverPostgresql() *drivers.PostgreSql {
	o, err := c.UnscopedSafeGetCoreDatabaseDriverPostgresql()
	if err != nil {
		panic(err)
	}
	return o
}

// CoreDatabaseDriverPostgresql is similar to GetCoreDatabaseDriverPostgresql.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetCoreDatabaseDriverPostgresql method.
// If the container can not be retrieved, it panics.
func CoreDatabaseDriverPostgresql(i interface{}) *drivers.PostgreSql {
	return C(i).GetCoreDatabaseDriverPostgresql()
}

// SafeGetCoreEventDispatcher works like SafeGet but only for CoreEventDispatcher.
// It does not return an interface but a *events.Dispatcher.
func (c *Container) SafeGetCoreEventDispatcher() (*events.Dispatcher, error) {
	i, err := c.ctn.SafeGet("core:event:dispatcher")
	if err != nil {
		var eo *events.Dispatcher
		return eo, err
	}
	o, ok := i.(*events.Dispatcher)
	if !ok {
		return o, errors.New("could get 'core:event:dispatcher' because the object could not be cast to *events.Dispatcher")
	}
	return o, nil
}

// GetCoreEventDispatcher is similar to SafeGetCoreEventDispatcher but it does not return the error.
// Instead it panics.
func (c *Container) GetCoreEventDispatcher() *events.Dispatcher {
	o, err := c.SafeGetCoreEventDispatcher()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetCoreEventDispatcher works like UnscopedSafeGet but only for CoreEventDispatcher.
// It does not return an interface but a *events.Dispatcher.
func (c *Container) UnscopedSafeGetCoreEventDispatcher() (*events.Dispatcher, error) {
	i, err := c.ctn.UnscopedSafeGet("core:event:dispatcher")
	if err != nil {
		var eo *events.Dispatcher
		return eo, err
	}
	o, ok := i.(*events.Dispatcher)
	if !ok {
		return o, errors.New("could get 'core:event:dispatcher' because the object could not be cast to *events.Dispatcher")
	}
	return o, nil
}

// UnscopedGetCoreEventDispatcher is similar to UnscopedSafeGetCoreEventDispatcher but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetCoreEventDispatcher() *events.Dispatcher {
	o, err := c.UnscopedSafeGetCoreEventDispatcher()
	if err != nil {
		panic(err)
	}
	return o
}

// CoreEventDispatcher is similar to GetCoreEventDispatcher.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetCoreEventDispatcher method.
// If the container can not be retrieved, it panics.
func CoreEventDispatcher(i interface{}) *events.Dispatcher {
	return C(i).GetCoreEventDispatcher()
}

// SafeGetCoreGrpcServer works like SafeGet but only for CoreGrpcServer.
// It does not return an interface but a *grpc.Server.
func (c *Container) SafeGetCoreGrpcServer() (*grpc.Server, error) {
	i, err := c.ctn.SafeGet("core:grpc:server")
	if err != nil {
		var eo *grpc.Server
		return eo, err
	}
	o, ok := i.(*grpc.Server)
	if !ok {
		return o, errors.New("could get 'core:grpc:server' because the object could not be cast to *grpc.Server")
	}
	return o, nil
}

// GetCoreGrpcServer is similar to SafeGetCoreGrpcServer but it does not return the error.
// Instead it panics.
func (c *Container) GetCoreGrpcServer() *grpc.Server {
	o, err := c.SafeGetCoreGrpcServer()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetCoreGrpcServer works like UnscopedSafeGet but only for CoreGrpcServer.
// It does not return an interface but a *grpc.Server.
func (c *Container) UnscopedSafeGetCoreGrpcServer() (*grpc.Server, error) {
	i, err := c.ctn.UnscopedSafeGet("core:grpc:server")
	if err != nil {
		var eo *grpc.Server
		return eo, err
	}
	o, ok := i.(*grpc.Server)
	if !ok {
		return o, errors.New("could get 'core:grpc:server' because the object could not be cast to *grpc.Server")
	}
	return o, nil
}

// UnscopedGetCoreGrpcServer is similar to UnscopedSafeGetCoreGrpcServer but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetCoreGrpcServer() *grpc.Server {
	o, err := c.UnscopedSafeGetCoreGrpcServer()
	if err != nil {
		panic(err)
	}
	return o
}

// CoreGrpcServer is similar to GetCoreGrpcServer.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetCoreGrpcServer method.
// If the container can not be retrieved, it panics.
func CoreGrpcServer(i interface{}) *grpc.Server {
	return C(i).GetCoreGrpcServer()
}

// SafeGetCoreHandlerHandler works like SafeGet but only for CoreHandlerHandler.
// It does not return an interface but a *handlers.Handler.
func (c *Container) SafeGetCoreHandlerHandler() (*handlers.Handler, error) {
	i, err := c.ctn.SafeGet("core:handler:handler")
	if err != nil {
		var eo *handlers.Handler
		return eo, err
	}
	o, ok := i.(*handlers.Handler)
	if !ok {
		return o, errors.New("could get 'core:handler:handler' because the object could not be cast to *handlers.Handler")
	}
	return o, nil
}

// GetCoreHandlerHandler is similar to SafeGetCoreHandlerHandler but it does not return the error.
// Instead it panics.
func (c *Container) GetCoreHandlerHandler() *handlers.Handler {
	o, err := c.SafeGetCoreHandlerHandler()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetCoreHandlerHandler works like UnscopedSafeGet but only for CoreHandlerHandler.
// It does not return an interface but a *handlers.Handler.
func (c *Container) UnscopedSafeGetCoreHandlerHandler() (*handlers.Handler, error) {
	i, err := c.ctn.UnscopedSafeGet("core:handler:handler")
	if err != nil {
		var eo *handlers.Handler
		return eo, err
	}
	o, ok := i.(*handlers.Handler)
	if !ok {
		return o, errors.New("could get 'core:handler:handler' because the object could not be cast to *handlers.Handler")
	}
	return o, nil
}

// UnscopedGetCoreHandlerHandler is similar to UnscopedSafeGetCoreHandlerHandler but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetCoreHandlerHandler() *handlers.Handler {
	o, err := c.UnscopedSafeGetCoreHandlerHandler()
	if err != nil {
		panic(err)
	}
	return o
}

// CoreHandlerHandler is similar to GetCoreHandlerHandler.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetCoreHandlerHandler method.
// If the container can not be retrieved, it panics.
func CoreHandlerHandler(i interface{}) *handlers.Handler {
	return C(i).GetCoreHandlerHandler()
}

// SafeGetCoreHandlerLogger works like SafeGet but only for CoreHandlerLogger.
// It does not return an interface but a *handlers.Logger.
func (c *Container) SafeGetCoreHandlerLogger() (*handlers.Logger, error) {
	i, err := c.ctn.SafeGet("core:handler:logger")
	if err != nil {
		var eo *handlers.Logger
		return eo, err
	}
	o, ok := i.(*handlers.Logger)
	if !ok {
		return o, errors.New("could get 'core:handler:logger' because the object could not be cast to *handlers.Logger")
	}
	return o, nil
}

// GetCoreHandlerLogger is similar to SafeGetCoreHandlerLogger but it does not return the error.
// Instead it panics.
func (c *Container) GetCoreHandlerLogger() *handlers.Logger {
	o, err := c.SafeGetCoreHandlerLogger()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetCoreHandlerLogger works like UnscopedSafeGet but only for CoreHandlerLogger.
// It does not return an interface but a *handlers.Logger.
func (c *Container) UnscopedSafeGetCoreHandlerLogger() (*handlers.Logger, error) {
	i, err := c.ctn.UnscopedSafeGet("core:handler:logger")
	if err != nil {
		var eo *handlers.Logger
		return eo, err
	}
	o, ok := i.(*handlers.Logger)
	if !ok {
		return o, errors.New("could get 'core:handler:logger' because the object could not be cast to *handlers.Logger")
	}
	return o, nil
}

// UnscopedGetCoreHandlerLogger is similar to UnscopedSafeGetCoreHandlerLogger but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetCoreHandlerLogger() *handlers.Logger {
	o, err := c.UnscopedSafeGetCoreHandlerLogger()
	if err != nil {
		panic(err)
	}
	return o
}

// CoreHandlerLogger is similar to GetCoreHandlerLogger.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetCoreHandlerLogger method.
// If the container can not be retrieved, it panics.
func CoreHandlerLogger(i interface{}) *handlers.Logger {
	return C(i).GetCoreHandlerLogger()
}

// SafeGetCoreHandlerMessager works like SafeGet but only for CoreHandlerMessager.
// It does not return an interface but a *handlers.Messenger.
func (c *Container) SafeGetCoreHandlerMessager() (*handlers.Messenger, error) {
	i, err := c.ctn.SafeGet("core:handler:messager")
	if err != nil {
		var eo *handlers.Messenger
		return eo, err
	}
	o, ok := i.(*handlers.Messenger)
	if !ok {
		return o, errors.New("could get 'core:handler:messager' because the object could not be cast to *handlers.Messenger")
	}
	return o, nil
}

// GetCoreHandlerMessager is similar to SafeGetCoreHandlerMessager but it does not return the error.
// Instead it panics.
func (c *Container) GetCoreHandlerMessager() *handlers.Messenger {
	o, err := c.SafeGetCoreHandlerMessager()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetCoreHandlerMessager works like UnscopedSafeGet but only for CoreHandlerMessager.
// It does not return an interface but a *handlers.Messenger.
func (c *Container) UnscopedSafeGetCoreHandlerMessager() (*handlers.Messenger, error) {
	i, err := c.ctn.UnscopedSafeGet("core:handler:messager")
	if err != nil {
		var eo *handlers.Messenger
		return eo, err
	}
	o, ok := i.(*handlers.Messenger)
	if !ok {
		return o, errors.New("could get 'core:handler:messager' because the object could not be cast to *handlers.Messenger")
	}
	return o, nil
}

// UnscopedGetCoreHandlerMessager is similar to UnscopedSafeGetCoreHandlerMessager but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetCoreHandlerMessager() *handlers.Messenger {
	o, err := c.UnscopedSafeGetCoreHandlerMessager()
	if err != nil {
		panic(err)
	}
	return o
}

// CoreHandlerMessager is similar to GetCoreHandlerMessager.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetCoreHandlerMessager method.
// If the container can not be retrieved, it panics.
func CoreHandlerMessager(i interface{}) *handlers.Messenger {
	return C(i).GetCoreHandlerMessager()
}

// SafeGetCoreHandlerMiddleware works like SafeGet but only for CoreHandlerMiddleware.
// It does not return an interface but a *handlers.Middleware.
func (c *Container) SafeGetCoreHandlerMiddleware() (*handlers.Middleware, error) {
	i, err := c.ctn.SafeGet("core:handler:middleware")
	if err != nil {
		var eo *handlers.Middleware
		return eo, err
	}
	o, ok := i.(*handlers.Middleware)
	if !ok {
		return o, errors.New("could get 'core:handler:middleware' because the object could not be cast to *handlers.Middleware")
	}
	return o, nil
}

// GetCoreHandlerMiddleware is similar to SafeGetCoreHandlerMiddleware but it does not return the error.
// Instead it panics.
func (c *Container) GetCoreHandlerMiddleware() *handlers.Middleware {
	o, err := c.SafeGetCoreHandlerMiddleware()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetCoreHandlerMiddleware works like UnscopedSafeGet but only for CoreHandlerMiddleware.
// It does not return an interface but a *handlers.Middleware.
func (c *Container) UnscopedSafeGetCoreHandlerMiddleware() (*handlers.Middleware, error) {
	i, err := c.ctn.UnscopedSafeGet("core:handler:middleware")
	if err != nil {
		var eo *handlers.Middleware
		return eo, err
	}
	o, ok := i.(*handlers.Middleware)
	if !ok {
		return o, errors.New("could get 'core:handler:middleware' because the object could not be cast to *handlers.Middleware")
	}
	return o, nil
}

// UnscopedGetCoreHandlerMiddleware is similar to UnscopedSafeGetCoreHandlerMiddleware but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetCoreHandlerMiddleware() *handlers.Middleware {
	o, err := c.UnscopedSafeGetCoreHandlerMiddleware()
	if err != nil {
		panic(err)
	}
	return o
}

// CoreHandlerMiddleware is similar to GetCoreHandlerMiddleware.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetCoreHandlerMiddleware method.
// If the container can not be retrieved, it panics.
func CoreHandlerMiddleware(i interface{}) *handlers.Middleware {
	return C(i).GetCoreHandlerMiddleware()
}

// SafeGetCoreHandlerRouter works like SafeGet but only for CoreHandlerRouter.
// It does not return an interface but a *handlers.Router.
func (c *Container) SafeGetCoreHandlerRouter() (*handlers.Router, error) {
	i, err := c.ctn.SafeGet("core:handler:router")
	if err != nil {
		var eo *handlers.Router
		return eo, err
	}
	o, ok := i.(*handlers.Router)
	if !ok {
		return o, errors.New("could get 'core:handler:router' because the object could not be cast to *handlers.Router")
	}
	return o, nil
}

// GetCoreHandlerRouter is similar to SafeGetCoreHandlerRouter but it does not return the error.
// Instead it panics.
func (c *Container) GetCoreHandlerRouter() *handlers.Router {
	o, err := c.SafeGetCoreHandlerRouter()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetCoreHandlerRouter works like UnscopedSafeGet but only for CoreHandlerRouter.
// It does not return an interface but a *handlers.Router.
func (c *Container) UnscopedSafeGetCoreHandlerRouter() (*handlers.Router, error) {
	i, err := c.ctn.UnscopedSafeGet("core:handler:router")
	if err != nil {
		var eo *handlers.Router
		return eo, err
	}
	o, ok := i.(*handlers.Router)
	if !ok {
		return o, errors.New("could get 'core:handler:router' because the object could not be cast to *handlers.Router")
	}
	return o, nil
}

// UnscopedGetCoreHandlerRouter is similar to UnscopedSafeGetCoreHandlerRouter but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetCoreHandlerRouter() *handlers.Router {
	o, err := c.UnscopedSafeGetCoreHandlerRouter()
	if err != nil {
		panic(err)
	}
	return o
}

// CoreHandlerRouter is similar to GetCoreHandlerRouter.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetCoreHandlerRouter method.
// If the container can not be retrieved, it panics.
func CoreHandlerRouter(i interface{}) *handlers.Router {
	return C(i).GetCoreHandlerRouter()
}

// SafeGetCoreHttpMux works like SafeGet but only for CoreHttpMux.
// It does not return an interface but a *aliashttp.ServeMux.
func (c *Container) SafeGetCoreHttpMux() (*aliashttp.ServeMux, error) {
	i, err := c.ctn.SafeGet("core:http:mux")
	if err != nil {
		var eo *aliashttp.ServeMux
		return eo, err
	}
	o, ok := i.(*aliashttp.ServeMux)
	if !ok {
		return o, errors.New("could get 'core:http:mux' because the object could not be cast to *aliashttp.ServeMux")
	}
	return o, nil
}

// GetCoreHttpMux is similar to SafeGetCoreHttpMux but it does not return the error.
// Instead it panics.
func (c *Container) GetCoreHttpMux() *aliashttp.ServeMux {
	o, err := c.SafeGetCoreHttpMux()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetCoreHttpMux works like UnscopedSafeGet but only for CoreHttpMux.
// It does not return an interface but a *aliashttp.ServeMux.
func (c *Container) UnscopedSafeGetCoreHttpMux() (*aliashttp.ServeMux, error) {
	i, err := c.ctn.UnscopedSafeGet("core:http:mux")
	if err != nil {
		var eo *aliashttp.ServeMux
		return eo, err
	}
	o, ok := i.(*aliashttp.ServeMux)
	if !ok {
		return o, errors.New("could get 'core:http:mux' because the object could not be cast to *aliashttp.ServeMux")
	}
	return o, nil
}

// UnscopedGetCoreHttpMux is similar to UnscopedSafeGetCoreHttpMux but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetCoreHttpMux() *aliashttp.ServeMux {
	o, err := c.UnscopedSafeGetCoreHttpMux()
	if err != nil {
		panic(err)
	}
	return o
}

// CoreHttpMux is similar to GetCoreHttpMux.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetCoreHttpMux method.
// If the container can not be retrieved, it panics.
func CoreHttpMux(i interface{}) *aliashttp.ServeMux {
	return C(i).GetCoreHttpMux()
}

// SafeGetCoreInterfaceDatabase works like SafeGet but only for CoreInterfaceDatabase.
// It does not return an interface but a *interfaces.Database.
func (c *Container) SafeGetCoreInterfaceDatabase() (*interfaces.Database, error) {
	i, err := c.ctn.SafeGet("core:interface:database")
	if err != nil {
		var eo *interfaces.Database
		return eo, err
	}
	o, ok := i.(*interfaces.Database)
	if !ok {
		return o, errors.New("could get 'core:interface:database' because the object could not be cast to *interfaces.Database")
	}
	return o, nil
}

// GetCoreInterfaceDatabase is similar to SafeGetCoreInterfaceDatabase but it does not return the error.
// Instead it panics.
func (c *Container) GetCoreInterfaceDatabase() *interfaces.Database {
	o, err := c.SafeGetCoreInterfaceDatabase()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetCoreInterfaceDatabase works like UnscopedSafeGet but only for CoreInterfaceDatabase.
// It does not return an interface but a *interfaces.Database.
func (c *Container) UnscopedSafeGetCoreInterfaceDatabase() (*interfaces.Database, error) {
	i, err := c.ctn.UnscopedSafeGet("core:interface:database")
	if err != nil {
		var eo *interfaces.Database
		return eo, err
	}
	o, ok := i.(*interfaces.Database)
	if !ok {
		return o, errors.New("could get 'core:interface:database' because the object could not be cast to *interfaces.Database")
	}
	return o, nil
}

// UnscopedGetCoreInterfaceDatabase is similar to UnscopedSafeGetCoreInterfaceDatabase but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetCoreInterfaceDatabase() *interfaces.Database {
	o, err := c.UnscopedSafeGetCoreInterfaceDatabase()
	if err != nil {
		panic(err)
	}
	return o
}

// CoreInterfaceDatabase is similar to GetCoreInterfaceDatabase.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetCoreInterfaceDatabase method.
// If the container can not be retrieved, it panics.
func CoreInterfaceDatabase(i interface{}) *interfaces.Database {
	return C(i).GetCoreInterfaceDatabase()
}

// SafeGetCoreInterfaceGrpc works like SafeGet but only for CoreInterfaceGrpc.
// It does not return an interface but a *interfaces.GRpc.
func (c *Container) SafeGetCoreInterfaceGrpc() (*interfaces.GRpc, error) {
	i, err := c.ctn.SafeGet("core:interface:grpc")
	if err != nil {
		var eo *interfaces.GRpc
		return eo, err
	}
	o, ok := i.(*interfaces.GRpc)
	if !ok {
		return o, errors.New("could get 'core:interface:grpc' because the object could not be cast to *interfaces.GRpc")
	}
	return o, nil
}

// GetCoreInterfaceGrpc is similar to SafeGetCoreInterfaceGrpc but it does not return the error.
// Instead it panics.
func (c *Container) GetCoreInterfaceGrpc() *interfaces.GRpc {
	o, err := c.SafeGetCoreInterfaceGrpc()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetCoreInterfaceGrpc works like UnscopedSafeGet but only for CoreInterfaceGrpc.
// It does not return an interface but a *interfaces.GRpc.
func (c *Container) UnscopedSafeGetCoreInterfaceGrpc() (*interfaces.GRpc, error) {
	i, err := c.ctn.UnscopedSafeGet("core:interface:grpc")
	if err != nil {
		var eo *interfaces.GRpc
		return eo, err
	}
	o, ok := i.(*interfaces.GRpc)
	if !ok {
		return o, errors.New("could get 'core:interface:grpc' because the object could not be cast to *interfaces.GRpc")
	}
	return o, nil
}

// UnscopedGetCoreInterfaceGrpc is similar to UnscopedSafeGetCoreInterfaceGrpc but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetCoreInterfaceGrpc() *interfaces.GRpc {
	o, err := c.UnscopedSafeGetCoreInterfaceGrpc()
	if err != nil {
		panic(err)
	}
	return o
}

// CoreInterfaceGrpc is similar to GetCoreInterfaceGrpc.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetCoreInterfaceGrpc method.
// If the container can not be retrieved, it panics.
func CoreInterfaceGrpc(i interface{}) *interfaces.GRpc {
	return C(i).GetCoreInterfaceGrpc()
}

// SafeGetCoreInterfaceQueue works like SafeGet but only for CoreInterfaceQueue.
// It does not return an interface but a *interfaces.Queue.
func (c *Container) SafeGetCoreInterfaceQueue() (*interfaces.Queue, error) {
	i, err := c.ctn.SafeGet("core:interface:queue")
	if err != nil {
		var eo *interfaces.Queue
		return eo, err
	}
	o, ok := i.(*interfaces.Queue)
	if !ok {
		return o, errors.New("could get 'core:interface:queue' because the object could not be cast to *interfaces.Queue")
	}
	return o, nil
}

// GetCoreInterfaceQueue is similar to SafeGetCoreInterfaceQueue but it does not return the error.
// Instead it panics.
func (c *Container) GetCoreInterfaceQueue() *interfaces.Queue {
	o, err := c.SafeGetCoreInterfaceQueue()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetCoreInterfaceQueue works like UnscopedSafeGet but only for CoreInterfaceQueue.
// It does not return an interface but a *interfaces.Queue.
func (c *Container) UnscopedSafeGetCoreInterfaceQueue() (*interfaces.Queue, error) {
	i, err := c.ctn.UnscopedSafeGet("core:interface:queue")
	if err != nil {
		var eo *interfaces.Queue
		return eo, err
	}
	o, ok := i.(*interfaces.Queue)
	if !ok {
		return o, errors.New("could get 'core:interface:queue' because the object could not be cast to *interfaces.Queue")
	}
	return o, nil
}

// UnscopedGetCoreInterfaceQueue is similar to UnscopedSafeGetCoreInterfaceQueue but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetCoreInterfaceQueue() *interfaces.Queue {
	o, err := c.UnscopedSafeGetCoreInterfaceQueue()
	if err != nil {
		panic(err)
	}
	return o
}

// CoreInterfaceQueue is similar to GetCoreInterfaceQueue.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetCoreInterfaceQueue method.
// If the container can not be retrieved, it panics.
func CoreInterfaceQueue(i interface{}) *interfaces.Queue {
	return C(i).GetCoreInterfaceQueue()
}

// SafeGetCoreInterfaceRest works like SafeGet but only for CoreInterfaceRest.
// It does not return an interface but a *interfaces.Rest.
func (c *Container) SafeGetCoreInterfaceRest() (*interfaces.Rest, error) {
	i, err := c.ctn.SafeGet("core:interface:rest")
	if err != nil {
		var eo *interfaces.Rest
		return eo, err
	}
	o, ok := i.(*interfaces.Rest)
	if !ok {
		return o, errors.New("could get 'core:interface:rest' because the object could not be cast to *interfaces.Rest")
	}
	return o, nil
}

// GetCoreInterfaceRest is similar to SafeGetCoreInterfaceRest but it does not return the error.
// Instead it panics.
func (c *Container) GetCoreInterfaceRest() *interfaces.Rest {
	o, err := c.SafeGetCoreInterfaceRest()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetCoreInterfaceRest works like UnscopedSafeGet but only for CoreInterfaceRest.
// It does not return an interface but a *interfaces.Rest.
func (c *Container) UnscopedSafeGetCoreInterfaceRest() (*interfaces.Rest, error) {
	i, err := c.ctn.UnscopedSafeGet("core:interface:rest")
	if err != nil {
		var eo *interfaces.Rest
		return eo, err
	}
	o, ok := i.(*interfaces.Rest)
	if !ok {
		return o, errors.New("could get 'core:interface:rest' because the object could not be cast to *interfaces.Rest")
	}
	return o, nil
}

// UnscopedGetCoreInterfaceRest is similar to UnscopedSafeGetCoreInterfaceRest but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetCoreInterfaceRest() *interfaces.Rest {
	o, err := c.UnscopedSafeGetCoreInterfaceRest()
	if err != nil {
		panic(err)
	}
	return o
}

// CoreInterfaceRest is similar to GetCoreInterfaceRest.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetCoreInterfaceRest method.
// If the container can not be retrieved, it panics.
func CoreInterfaceRest(i interface{}) *interfaces.Rest {
	return C(i).GetCoreInterfaceRest()
}

// SafeGetCoreMessageConfig works like SafeGet but only for CoreMessageConfig.
// It does not return an interface but a amqp.Config.
func (c *Container) SafeGetCoreMessageConfig() (amqp.Config, error) {
	i, err := c.ctn.SafeGet("core:message:config")
	if err != nil {
		var eo amqp.Config
		return eo, err
	}
	o, ok := i.(amqp.Config)
	if !ok {
		return o, errors.New("could get 'core:message:config' because the object could not be cast to amqp.Config")
	}
	return o, nil
}

// GetCoreMessageConfig is similar to SafeGetCoreMessageConfig but it does not return the error.
// Instead it panics.
func (c *Container) GetCoreMessageConfig() amqp.Config {
	o, err := c.SafeGetCoreMessageConfig()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetCoreMessageConfig works like UnscopedSafeGet but only for CoreMessageConfig.
// It does not return an interface but a amqp.Config.
func (c *Container) UnscopedSafeGetCoreMessageConfig() (amqp.Config, error) {
	i, err := c.ctn.UnscopedSafeGet("core:message:config")
	if err != nil {
		var eo amqp.Config
		return eo, err
	}
	o, ok := i.(amqp.Config)
	if !ok {
		return o, errors.New("could get 'core:message:config' because the object could not be cast to amqp.Config")
	}
	return o, nil
}

// UnscopedGetCoreMessageConfig is similar to UnscopedSafeGetCoreMessageConfig but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetCoreMessageConfig() amqp.Config {
	o, err := c.UnscopedSafeGetCoreMessageConfig()
	if err != nil {
		panic(err)
	}
	return o
}

// CoreMessageConfig is similar to GetCoreMessageConfig.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetCoreMessageConfig method.
// If the container can not be retrieved, it panics.
func CoreMessageConfig(i interface{}) amqp.Config {
	return C(i).GetCoreMessageConfig()
}

// SafeGetCoreMessageConsumer works like SafeGet but only for CoreMessageConsumer.
// It does not return an interface but a *amqp.Subscriber.
func (c *Container) SafeGetCoreMessageConsumer() (*amqp.Subscriber, error) {
	i, err := c.ctn.SafeGet("core:message:consumer")
	if err != nil {
		var eo *amqp.Subscriber
		return eo, err
	}
	o, ok := i.(*amqp.Subscriber)
	if !ok {
		return o, errors.New("could get 'core:message:consumer' because the object could not be cast to *amqp.Subscriber")
	}
	return o, nil
}

// GetCoreMessageConsumer is similar to SafeGetCoreMessageConsumer but it does not return the error.
// Instead it panics.
func (c *Container) GetCoreMessageConsumer() *amqp.Subscriber {
	o, err := c.SafeGetCoreMessageConsumer()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetCoreMessageConsumer works like UnscopedSafeGet but only for CoreMessageConsumer.
// It does not return an interface but a *amqp.Subscriber.
func (c *Container) UnscopedSafeGetCoreMessageConsumer() (*amqp.Subscriber, error) {
	i, err := c.ctn.UnscopedSafeGet("core:message:consumer")
	if err != nil {
		var eo *amqp.Subscriber
		return eo, err
	}
	o, ok := i.(*amqp.Subscriber)
	if !ok {
		return o, errors.New("could get 'core:message:consumer' because the object could not be cast to *amqp.Subscriber")
	}
	return o, nil
}

// UnscopedGetCoreMessageConsumer is similar to UnscopedSafeGetCoreMessageConsumer but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetCoreMessageConsumer() *amqp.Subscriber {
	o, err := c.UnscopedSafeGetCoreMessageConsumer()
	if err != nil {
		panic(err)
	}
	return o
}

// CoreMessageConsumer is similar to GetCoreMessageConsumer.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetCoreMessageConsumer method.
// If the container can not be retrieved, it panics.
func CoreMessageConsumer(i interface{}) *amqp.Subscriber {
	return C(i).GetCoreMessageConsumer()
}

// SafeGetCoreMessagePublisher works like SafeGet but only for CoreMessagePublisher.
// It does not return an interface but a *amqp.Publisher.
func (c *Container) SafeGetCoreMessagePublisher() (*amqp.Publisher, error) {
	i, err := c.ctn.SafeGet("core:message:publisher")
	if err != nil {
		var eo *amqp.Publisher
		return eo, err
	}
	o, ok := i.(*amqp.Publisher)
	if !ok {
		return o, errors.New("could get 'core:message:publisher' because the object could not be cast to *amqp.Publisher")
	}
	return o, nil
}

// GetCoreMessagePublisher is similar to SafeGetCoreMessagePublisher but it does not return the error.
// Instead it panics.
func (c *Container) GetCoreMessagePublisher() *amqp.Publisher {
	o, err := c.SafeGetCoreMessagePublisher()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetCoreMessagePublisher works like UnscopedSafeGet but only for CoreMessagePublisher.
// It does not return an interface but a *amqp.Publisher.
func (c *Container) UnscopedSafeGetCoreMessagePublisher() (*amqp.Publisher, error) {
	i, err := c.ctn.UnscopedSafeGet("core:message:publisher")
	if err != nil {
		var eo *amqp.Publisher
		return eo, err
	}
	o, ok := i.(*amqp.Publisher)
	if !ok {
		return o, errors.New("could get 'core:message:publisher' because the object could not be cast to *amqp.Publisher")
	}
	return o, nil
}

// UnscopedGetCoreMessagePublisher is similar to UnscopedSafeGetCoreMessagePublisher but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetCoreMessagePublisher() *amqp.Publisher {
	o, err := c.UnscopedSafeGetCoreMessagePublisher()
	if err != nil {
		panic(err)
	}
	return o
}

// CoreMessagePublisher is similar to GetCoreMessagePublisher.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetCoreMessagePublisher method.
// If the container can not be retrieved, it panics.
func CoreMessagePublisher(i interface{}) *amqp.Publisher {
	return C(i).GetCoreMessagePublisher()
}

// SafeGetCoreMiddlewareAuth works like SafeGet but only for CoreMiddlewareAuth.
// It does not return an interface but a *middlewares.Auth.
func (c *Container) SafeGetCoreMiddlewareAuth() (*middlewares.Auth, error) {
	i, err := c.ctn.SafeGet("core:middleware:auth")
	if err != nil {
		var eo *middlewares.Auth
		return eo, err
	}
	o, ok := i.(*middlewares.Auth)
	if !ok {
		return o, errors.New("could get 'core:middleware:auth' because the object could not be cast to *middlewares.Auth")
	}
	return o, nil
}

// GetCoreMiddlewareAuth is similar to SafeGetCoreMiddlewareAuth but it does not return the error.
// Instead it panics.
func (c *Container) GetCoreMiddlewareAuth() *middlewares.Auth {
	o, err := c.SafeGetCoreMiddlewareAuth()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetCoreMiddlewareAuth works like UnscopedSafeGet but only for CoreMiddlewareAuth.
// It does not return an interface but a *middlewares.Auth.
func (c *Container) UnscopedSafeGetCoreMiddlewareAuth() (*middlewares.Auth, error) {
	i, err := c.ctn.UnscopedSafeGet("core:middleware:auth")
	if err != nil {
		var eo *middlewares.Auth
		return eo, err
	}
	o, ok := i.(*middlewares.Auth)
	if !ok {
		return o, errors.New("could get 'core:middleware:auth' because the object could not be cast to *middlewares.Auth")
	}
	return o, nil
}

// UnscopedGetCoreMiddlewareAuth is similar to UnscopedSafeGetCoreMiddlewareAuth but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetCoreMiddlewareAuth() *middlewares.Auth {
	o, err := c.UnscopedSafeGetCoreMiddlewareAuth()
	if err != nil {
		panic(err)
	}
	return o
}

// CoreMiddlewareAuth is similar to GetCoreMiddlewareAuth.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetCoreMiddlewareAuth method.
// If the container can not be retrieved, it panics.
func CoreMiddlewareAuth(i interface{}) *middlewares.Auth {
	return C(i).GetCoreMiddlewareAuth()
}

// SafeGetCorePaginationPaginator works like SafeGet but only for CorePaginationPaginator.
// It does not return an interface but a *paginations.Pagination.
func (c *Container) SafeGetCorePaginationPaginator() (*paginations.Pagination, error) {
	i, err := c.ctn.SafeGet("core:pagination:paginator")
	if err != nil {
		var eo *paginations.Pagination
		return eo, err
	}
	o, ok := i.(*paginations.Pagination)
	if !ok {
		return o, errors.New("could get 'core:pagination:paginator' because the object could not be cast to *paginations.Pagination")
	}
	return o, nil
}

// GetCorePaginationPaginator is similar to SafeGetCorePaginationPaginator but it does not return the error.
// Instead it panics.
func (c *Container) GetCorePaginationPaginator() *paginations.Pagination {
	o, err := c.SafeGetCorePaginationPaginator()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetCorePaginationPaginator works like UnscopedSafeGet but only for CorePaginationPaginator.
// It does not return an interface but a *paginations.Pagination.
func (c *Container) UnscopedSafeGetCorePaginationPaginator() (*paginations.Pagination, error) {
	i, err := c.ctn.UnscopedSafeGet("core:pagination:paginator")
	if err != nil {
		var eo *paginations.Pagination
		return eo, err
	}
	o, ok := i.(*paginations.Pagination)
	if !ok {
		return o, errors.New("could get 'core:pagination:paginator' because the object could not be cast to *paginations.Pagination")
	}
	return o, nil
}

// UnscopedGetCorePaginationPaginator is similar to UnscopedSafeGetCorePaginationPaginator but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetCorePaginationPaginator() *paginations.Pagination {
	o, err := c.UnscopedSafeGetCorePaginationPaginator()
	if err != nil {
		panic(err)
	}
	return o
}

// CorePaginationPaginator is similar to GetCorePaginationPaginator.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetCorePaginationPaginator method.
// If the container can not be retrieved, it panics.
func CorePaginationPaginator(i interface{}) *paginations.Pagination {
	return C(i).GetCorePaginationPaginator()
}

// SafeGetCoreRouterGateway works like SafeGet but only for CoreRouterGateway.
// It does not return an interface but a *routes.GRpcGateway.
func (c *Container) SafeGetCoreRouterGateway() (*routes.GRpcGateway, error) {
	i, err := c.ctn.SafeGet("core:router:gateway")
	if err != nil {
		var eo *routes.GRpcGateway
		return eo, err
	}
	o, ok := i.(*routes.GRpcGateway)
	if !ok {
		return o, errors.New("could get 'core:router:gateway' because the object could not be cast to *routes.GRpcGateway")
	}
	return o, nil
}

// GetCoreRouterGateway is similar to SafeGetCoreRouterGateway but it does not return the error.
// Instead it panics.
func (c *Container) GetCoreRouterGateway() *routes.GRpcGateway {
	o, err := c.SafeGetCoreRouterGateway()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetCoreRouterGateway works like UnscopedSafeGet but only for CoreRouterGateway.
// It does not return an interface but a *routes.GRpcGateway.
func (c *Container) UnscopedSafeGetCoreRouterGateway() (*routes.GRpcGateway, error) {
	i, err := c.ctn.UnscopedSafeGet("core:router:gateway")
	if err != nil {
		var eo *routes.GRpcGateway
		return eo, err
	}
	o, ok := i.(*routes.GRpcGateway)
	if !ok {
		return o, errors.New("could get 'core:router:gateway' because the object could not be cast to *routes.GRpcGateway")
	}
	return o, nil
}

// UnscopedGetCoreRouterGateway is similar to UnscopedSafeGetCoreRouterGateway but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetCoreRouterGateway() *routes.GRpcGateway {
	o, err := c.UnscopedSafeGetCoreRouterGateway()
	if err != nil {
		panic(err)
	}
	return o
}

// CoreRouterGateway is similar to GetCoreRouterGateway.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetCoreRouterGateway method.
// If the container can not be retrieved, it panics.
func CoreRouterGateway(i interface{}) *routes.GRpcGateway {
	return C(i).GetCoreRouterGateway()
}

// SafeGetCoreRouterMux works like SafeGet but only for CoreRouterMux.
// It does not return an interface but a *routes.MuxRouter.
func (c *Container) SafeGetCoreRouterMux() (*routes.MuxRouter, error) {
	i, err := c.ctn.SafeGet("core:router:mux")
	if err != nil {
		var eo *routes.MuxRouter
		return eo, err
	}
	o, ok := i.(*routes.MuxRouter)
	if !ok {
		return o, errors.New("could get 'core:router:mux' because the object could not be cast to *routes.MuxRouter")
	}
	return o, nil
}

// GetCoreRouterMux is similar to SafeGetCoreRouterMux but it does not return the error.
// Instead it panics.
func (c *Container) GetCoreRouterMux() *routes.MuxRouter {
	o, err := c.SafeGetCoreRouterMux()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetCoreRouterMux works like UnscopedSafeGet but only for CoreRouterMux.
// It does not return an interface but a *routes.MuxRouter.
func (c *Container) UnscopedSafeGetCoreRouterMux() (*routes.MuxRouter, error) {
	i, err := c.ctn.UnscopedSafeGet("core:router:mux")
	if err != nil {
		var eo *routes.MuxRouter
		return eo, err
	}
	o, ok := i.(*routes.MuxRouter)
	if !ok {
		return o, errors.New("could get 'core:router:mux' because the object could not be cast to *routes.MuxRouter")
	}
	return o, nil
}

// UnscopedGetCoreRouterMux is similar to UnscopedSafeGetCoreRouterMux but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetCoreRouterMux() *routes.MuxRouter {
	o, err := c.UnscopedSafeGetCoreRouterMux()
	if err != nil {
		panic(err)
	}
	return o
}

// CoreRouterMux is similar to GetCoreRouterMux.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetCoreRouterMux method.
// If the container can not be retrieved, it panics.
func CoreRouterMux(i interface{}) *routes.MuxRouter {
	return C(i).GetCoreRouterMux()
}

// SafeGetCoreUtilNumber works like SafeGet but only for CoreUtilNumber.
// It does not return an interface but a *utils.Number.
func (c *Container) SafeGetCoreUtilNumber() (*utils.Number, error) {
	i, err := c.ctn.SafeGet("core:util:number")
	if err != nil {
		var eo *utils.Number
		return eo, err
	}
	o, ok := i.(*utils.Number)
	if !ok {
		return o, errors.New("could get 'core:util:number' because the object could not be cast to *utils.Number")
	}
	return o, nil
}

// GetCoreUtilNumber is similar to SafeGetCoreUtilNumber but it does not return the error.
// Instead it panics.
func (c *Container) GetCoreUtilNumber() *utils.Number {
	o, err := c.SafeGetCoreUtilNumber()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetCoreUtilNumber works like UnscopedSafeGet but only for CoreUtilNumber.
// It does not return an interface but a *utils.Number.
func (c *Container) UnscopedSafeGetCoreUtilNumber() (*utils.Number, error) {
	i, err := c.ctn.UnscopedSafeGet("core:util:number")
	if err != nil {
		var eo *utils.Number
		return eo, err
	}
	o, ok := i.(*utils.Number)
	if !ok {
		return o, errors.New("could get 'core:util:number' because the object could not be cast to *utils.Number")
	}
	return o, nil
}

// UnscopedGetCoreUtilNumber is similar to UnscopedSafeGetCoreUtilNumber but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetCoreUtilNumber() *utils.Number {
	o, err := c.UnscopedSafeGetCoreUtilNumber()
	if err != nil {
		panic(err)
	}
	return o
}

// CoreUtilNumber is similar to GetCoreUtilNumber.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetCoreUtilNumber method.
// If the container can not be retrieved, it panics.
func CoreUtilNumber(i interface{}) *utils.Number {
	return C(i).GetCoreUtilNumber()
}

// SafeGetCoreUtilTime works like SafeGet but only for CoreUtilTime.
// It does not return an interface but a *utils.Time.
func (c *Container) SafeGetCoreUtilTime() (*utils.Time, error) {
	i, err := c.ctn.SafeGet("core:util:time")
	if err != nil {
		var eo *utils.Time
		return eo, err
	}
	o, ok := i.(*utils.Time)
	if !ok {
		return o, errors.New("could get 'core:util:time' because the object could not be cast to *utils.Time")
	}
	return o, nil
}

// GetCoreUtilTime is similar to SafeGetCoreUtilTime but it does not return the error.
// Instead it panics.
func (c *Container) GetCoreUtilTime() *utils.Time {
	o, err := c.SafeGetCoreUtilTime()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetCoreUtilTime works like UnscopedSafeGet but only for CoreUtilTime.
// It does not return an interface but a *utils.Time.
func (c *Container) UnscopedSafeGetCoreUtilTime() (*utils.Time, error) {
	i, err := c.ctn.UnscopedSafeGet("core:util:time")
	if err != nil {
		var eo *utils.Time
		return eo, err
	}
	o, ok := i.(*utils.Time)
	if !ok {
		return o, errors.New("could get 'core:util:time' because the object could not be cast to *utils.Time")
	}
	return o, nil
}

// UnscopedGetCoreUtilTime is similar to UnscopedSafeGetCoreUtilTime but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetCoreUtilTime() *utils.Time {
	o, err := c.UnscopedSafeGetCoreUtilTime()
	if err != nil {
		panic(err)
	}
	return o
}

// CoreUtilTime is similar to GetCoreUtilTime.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetCoreUtilTime method.
// If the container can not be retrieved, it panics.
func CoreUtilTime(i interface{}) *utils.Time {
	return C(i).GetCoreUtilTime()
}

// SafeGetCoreUtilWord works like SafeGet but only for CoreUtilWord.
// It does not return an interface but a *utils.Word.
func (c *Container) SafeGetCoreUtilWord() (*utils.Word, error) {
	i, err := c.ctn.SafeGet("core:util:word")
	if err != nil {
		var eo *utils.Word
		return eo, err
	}
	o, ok := i.(*utils.Word)
	if !ok {
		return o, errors.New("could get 'core:util:word' because the object could not be cast to *utils.Word")
	}
	return o, nil
}

// GetCoreUtilWord is similar to SafeGetCoreUtilWord but it does not return the error.
// Instead it panics.
func (c *Container) GetCoreUtilWord() *utils.Word {
	o, err := c.SafeGetCoreUtilWord()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetCoreUtilWord works like UnscopedSafeGet but only for CoreUtilWord.
// It does not return an interface but a *utils.Word.
func (c *Container) UnscopedSafeGetCoreUtilWord() (*utils.Word, error) {
	i, err := c.ctn.UnscopedSafeGet("core:util:word")
	if err != nil {
		var eo *utils.Word
		return eo, err
	}
	o, ok := i.(*utils.Word)
	if !ok {
		return o, errors.New("could get 'core:util:word' because the object could not be cast to *utils.Word")
	}
	return o, nil
}

// UnscopedGetCoreUtilWord is similar to UnscopedSafeGetCoreUtilWord but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetCoreUtilWord() *utils.Word {
	o, err := c.UnscopedSafeGetCoreUtilWord()
	if err != nil {
		panic(err)
	}
	return o
}

// CoreUtilWord is similar to GetCoreUtilWord.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetCoreUtilWord method.
// If the container can not be retrieved, it panics.
func CoreUtilWord(i interface{}) *utils.Word {
	return C(i).GetCoreUtilWord()
}
