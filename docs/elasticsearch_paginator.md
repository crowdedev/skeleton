# Use Elasticseearch Paginator

For elasticsearch, we use `https://github.com/olivere/elastic` as library

- Add Elasticsearch config to `dics/container.go`

```go
{
    Name:  "bima:elasticsearch:client",
    Scope: bima.Application,
    Build: func(dsn string) (*elastic.Client, error) {
        client, err := elastic.NewClient(
            elastic.SetURL(dsn),
            elastic.SetSniff(false),
            elastic.SetHealthcheck(false),
            elastic.SetGzip(true),
        )

        if err != nil {
            return nil, nil
        }

        color.New(color.FgCyan, color.Bold).Print("✓ ")
        fmt.Println("Elasticsearch configured")

        return client, nil
    },
    Params: dingo.Params{
        "0": "localhost:9200",
    },
},
{
    Name: "bima:pagination:adapter:elasticsearch",
    Scope: bima.Application,
    Build: func(env *configs.Env, client *elastic.Client, dispatcher *events.Dispatcher) (*adapter.ElasticsearchAdapter, error) {
        return &adapter.ElasticsearchAdapter{
            Debug:      env.Debug,
            Service:    env.Service.ConnonicalName,
            Client:     client,
            Dispatcher: dispatcher,
        }, nil
    },
    Params: dingo.Params{
        "0": dingo.Service("bima:config"),
        "1": dingo.Service("bima:elasticsearch:client"),
        "2": dingo.Service("bima:event:dispatcher"),
    },
},
```

- Change adaptor to Elasticsearch

```go
{
    Name: "bima:handler",
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
        "2": dingo.Service("bima:repository:gorm"),
        "3": dingo.Service("bima:pagination:adapter:elasticsearch"),
    },
},
```

- Implement Elasticsearch filter

```go
{
    Name:  "bima:listener:filter:elasticsearch",
    Scope: bima.Application,
    Build: (*filters.ElasticsearchFilter)(nil),
},
```

- Add to your `configs/listeners.yaml`

```yaml
listeners:
    - filter:elasticsearch
```

- Rerun your service
