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

Berikut adalah daftar event yang terdapat dalam **Skeleton** (terdapat dalam file [`handler.go`](ttps://github.com/crowdeco/skeleton/blob/main/handlers/handler.go))

```go
const PAGINATION_EVENT = "event.pagination"
const BEFORE_CREATE_EVENT = "event.before_create"
const AFTER_CREATE_EVENT = "event.after_create"
const BEFORE_UPDATE_EVENT = "event.before_update"
const AFTER_UPDATE_EVENT = "event.after_update"
const BEFORE_DELETE_EVENT = "event.before_delete"
const AFTER_DELETE_EVENT = "event.after_delete"
```

- Daftarkan struct pada DIC (selanjutnya disebut **service**) pada folder `dics/modules/<module>.go`, bila bingung bisa lihat contoh definisi dari `core:listener:create:elasticsearch` pada file [`dics/core.go`](https://github.com/crowdeco/skeleton/blob/main/dics/core.go#L260)

- Daftarkan service pada file `dics/dispatcher.go` 

- Rebuild DIC dengan perintah `go run cmds/dic/main.go`
