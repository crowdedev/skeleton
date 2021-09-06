# Cara mendaftarkan middleware

## Global middleware

- Buat file folder modul misal `middlewares/rate_limiter.go`

- Buat struct sesuai interface berikut:

```go
Middleware interface {
    Attach(request *http.Request, response http.ResponseWriter) bool
    Priority() int
}
```

- Daftarkan struct pada DIC (selanjutnya disebut **service**) pada folder `dics/<module>.go`, bila bingung bisa baca dokumentasi dari [Dingo](https://github.com/sarulabs/dingo)

- Daftarkan service pada file `configs/middlewares.yaml` 

- Rebuild DIC dengan perintah `go run cmds/dic/main.go`

## Middleware per route

- Sebelumnya, baca dulu [cara membuat custom route](./custom_route.md)

- Tambahkan method `Middlewares()`, middleware yang ingin digunakan pada route tersebut

```go
func (r *Route) Middlewares() []Middleware {
    return []Middleware{
        middleware1,
        middleware2,
        middleware3,
    }
}
```

- Perlu dijadikan perhatian, pada route middleware, skeleton tidak melakukan priority order seperti pada global middleware.
