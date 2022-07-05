# Implement JWT Auth

## Implement JWT Login

- Add jwt login route to `dics/container.go`

```go
{
    Name: "bima:route:jwt:login",
    Scope: bima.Application,
    Build: func(env *configs.Env) (*routes.JwtLogin, error) {
        return routes.DefaultJwtLogin("/api/v1/login", env.Secret, jwt.SigningMethodHS512.Name, true, routes.FindUserByUsernameAndPassword(func(username, password string) jwt.MapClaims {
            return jwt.MapClaims{
                "user": "admin",
            }
        })), nil
    },
    Params: dingo.Params{
        "0": dingo.Service("bima:config"),
    },
}
```

You need to implmement `routes.FindUserByUsernameAndPassword` function with your own logic. If you don't implement refresh token, pass `false` to 4th argument.

- Add to `configs/routes.yaml`

```yaml
routes:
    - jwt:login
```

## Implement JWT Validator (Middleware)

- Add jwt middleware to `dics/container.go`

```go
{
    Name: "bima:middleware:jwt",
    Scope: bima.Application,
    Build: func(env *configs.Env) (*middlewares.Jwt, error) {
        return &middlewares.Jwt{
            Debug:         env.Debug,
            Env:           env,
            Secret:        env.Secret,
            SigningMethod: jwt.SigningMethodHS512.Name,
            Whitelist:     "/health$",
        }, nil
    },
    Params: dingo.Params{
        "0": dingo.Service("bima:config"),
    },
}
```

- Add to `configs/middlewares.yaml`

```yaml
middlewares:
    - jwt
```

You can access user using `configs.Env.User` or via `request.Header.Get("X-Bima-User")`

## Implement Refresh  JWT

- Add refresh jwt route to `dics/container.go`

```go
{
    Name: "bima:route:jwt:refresh",
    Scope: bima.Application,
    Build: func(env *configs.Env) (*routes.JwtRefresh, error) {
        return &routes.JwtRefresh{
            PathUrl:       "/api/v1/token-refresh",
            Secret:        env.Secret,
            SigningMethod: jwt.SigningMethodHS512.Name,
            Expire:        730, //1 month
        }, nil
    },
    Params: dingo.Params{
        "0": dingo.Service("bima:config"),
    },
},
```

- Add to `configs/routes.yaml`

```yaml
routes:
    - jwt:login
    - jwt:refresh
```
