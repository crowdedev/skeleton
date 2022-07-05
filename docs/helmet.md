# Enable CORS 

- Add Helmet middleware to `dics/container.go`

```go
{
    Name:  "bima:middleware:helmet",
    Scope: bima.Application,
    Build: (*middlewares.Helmet)(nil),
},
```

- Add to `configs/middlewares.yaml`

```yaml
middlewares:
    - helmet
```
