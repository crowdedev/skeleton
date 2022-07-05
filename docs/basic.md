# Implement Basic Auth

- Add basic auth middleware to `dics/container.go`

```go
{
    Name:  "bima:middleware:basic-auth",
    Scope: bima.Application,
    Build: (*middlewares.BasicAuth)(nil),
    Params: dingo.Params{
        "Validator": func(username, password string) bool {
			return true
		},
    },
},
```

You need to implement `Validator` with your own logic

- Add to `configs/middlewares.yaml`

```yaml
middlewares:
    - basic-auth
```

