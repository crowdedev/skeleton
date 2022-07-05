# Enable CORS 

- Add CORS middleware to `dics/container.go`

```go
{
    Name:  "bima:middleware:cors",
    Scope: bima.Application,
    Build: (*middlewares.Cors)(nil),
    Params: dingo.Params{
        "Options": cors.Options{},
    },
},
```

You can refer to [cors](github.com/rs/cors) for more options

- Add to `configs/middlewares.yaml`

```yaml
middlewares:
    - cors
```
