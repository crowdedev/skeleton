# Use MongoDB as Storage

## Basic Usage

- Create MongoDB Database with `bima_skeleton` as name

- Add `DB_DRIVER`, `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, and `DB_PASSWORD` to `.env`

```bash
DB_DRIVER=mongo
DB_HOST=localhost
DB_PORT=27017
DB_NAME=bima_skeleton
DB_USER=root
DB_PASSWORD=aden
```

- Add MongoDB Repository and Paginator Adapter to `dics/containers.go`

```go
{
    Name:  "bima:repository:mongo",
    Scope: bima.Application,
    Build: (*repositories.MongoRepository)(nil),
},
{
    Name:  "bima:pagination:adapter:mongo",
    Scope: bima.Application,
    Build: func(
        env *configs.Env,
        dispatcher *events.Dispatcher,
    ) (*adapters.MongodbAdapter, error) {
        return &adapters.MongodbAdapter{
            Debug:      env.Debug,
            Dispatcher: dispatcher,
        }, nil
    },
    Params: dingo.Params{
        "0": dingo.Service("bima:config"),
        "1": dingo.Service("bima:event:dispatcher"),
    },
},
```

- Change `bima:handler` definition params to mongodb

```go
{
    Name:  "bima:handler",
    Scope: bima.Application,
    Build: func(
        env *configs.Env,
        dispatcher *events.Dispatcher,
        repository repositories.Repository,
        adapter paginations.Adapter,
    ) (*handlers.Handler, error) {
        return &handlers.Handler{
            Debug:      env.Debug,
            Dispatcher: dispatcher,
            Repository: repository,
            Adapter:    adapter,
        }, nil
    },
    Params: dingo.Params{
        "0": dingo.Service("bima:config"),
        "1": dingo.Service("bima:event:dispatcher"),
        "2": dingo.Service("bima:repository:mongo"),
        "3": dingo.Service("bima:pagination:adapter:mongo"),
    },
},
```

## Using Raplica Set or MongoDB Service

To use replica set, just put full dsn in `DB_HOST`

```bash
DB_DRIVER=mongo
DB_HOST=mongodb+srv://user:password@region.mongodb.net
```
