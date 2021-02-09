# Cara mendaftarkan middleware

- Buat file `middlewares/rate_limiter.go`

- Buat struct sesuai interface berikut:

```go
Middleware interface {
    Attach(request *http.Request, response http.ResponseWriter) bool
    Priority() int
}
```

- Daftarkan struct pada DIC (selanjutnya disebut **service**) pada folder `dics/modules/<module>.go`, bila bingung bisa lihat contoh definisi dari `core:middleware:auth` pada file [`dics/core.go`](https://github.com/crowdeco/skeleton/blob/main/dics/core.go#L391)

- Daftarkan service pada file `dics/middleware.go` 

- Rebuild DIC dengan perintah `go run cmds/dic/main.go`
