# Cara memodifikasi alur menggunakan event listener

- Buat file folder modul misal `listeners/serach.go`

- Buat struct sesuai interface berikut:

```go
Listener interface {
    Handle(event interface{})
    Listen() string
    Priority() int
}
```

Berikut adalah daftar event yang terdapat dalam **Skeleton** (terdapat dalam file [`events.go`](https://github.com/KejawenLab/bima/blob/main/events/event.go))

```go
const PAGINATION_EVENT = "event.pagination"
const BEFORE_CREATE_EVENT = "event.before_create"
const AFTER_CREATE_EVENT = "event.after_create"
const BEFORE_UPDATE_EVENT = "event.before_update"
const AFTER_UPDATE_EVENT = "event.after_update"
const BEFORE_DELETE_EVENT = "event.before_delete"
const AFTER_DELETE_EVENT = "event.after_delete"
const REQUEST_EVENT = "event.request"
const RESPONSE_EVENT = "event.response"
```

- Daftarkan struct pada DIC (selanjutnya disebut **service**) pada folder `dics/<module>.go`, bila bingung bisa baca dokumentasi dari [Dingo](https://github.com/sarulabs/dingo)

- Daftarkan service pada file `configs/listeners.yaml` 

- Rebuild DIC dengan perintah `go run cmds/dic/main.go`
