# Cara mendaftarkan middleware

- Buat file `middlewares/rate_limiter.go`

- Buat struct sesuai interface berikut:

```go
Middleware interface {
    Attach(request *http.Request, response http.ResponseWriter) bool
    Priority() int
}
```

- Daftarkan struct pada DIC (selanjutnya disebut **service**) pada folder `dics/<module>.go`, bila bingung bisa baca dokumentasi dari [Dingo](https://github.com/sarulabs/dingo)

- Daftarkan service pada file `middlewares.yaml` 

- Rebuild DIC dengan perintah `go run cmds/dic/main.go`
