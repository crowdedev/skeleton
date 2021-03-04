# Cara membuat custom route

- Buat file folder modul misal `controllers/upload.go`

- Buat struct sesuai interface berikut:

```go
Route interface {
    Path() string
    Method() string
    Handle(w http.ResponseWriter, r *http.Request, params map[string]string)
    SetClient(client *grpc.ClientConn)
    Middlewares() []Middleware
}
```
- Daftarkan struct pada DIC (selanjutnya disebut **service**) pada folder `dics/<module>.go`, bila bingung bisa baca dokumentasi dari [Dingo](https://github.com/sarulabs/dingo)

- Daftarkan service pada file `configs/routes.yaml`

- Rebuild DIC dengan perintah `go run cmds/dic/main.go`
