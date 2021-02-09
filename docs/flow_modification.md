# Cara memodifikasi alur menggunakan event listener

- Buat file `listeners/serach.go`

- Buat struct sesuai interface berikut:

```go
Listener interface {
    Handle(event interface{})
    Listen() string
    Priority() int
}
```

- Daftarkan struct pada DIC (selanjutnya disebut **service**) pada folder `dics/modules/<module>.go`, bila bingung bisa lihat contoh definisi dari `listeners/creates/elasticsearch` pada file [`dics/core.go`](https://github.com/crowdeco/skeleton/blob/main/dics/core.go#L253)

- Daftarkan service pada file `dics/dispatcher.go` 

- Rebuild DIC dengan perintah `go run cmds/dic/main.go`
