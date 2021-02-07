## Memulai Skeleton

- Buat modul baru dengan perintah `go run cmd/generator/main.go`

- Ikuti setiap langkah yang ada

- Daftarkan modul ke sistem pada file `dics/provider.go`

```go
	if err := p.AddDefSlice(modules.<NamaModul>); err != nil {
		return err
	}
```

`NamaModul` adalah nama variabel pada file `dics/modules/<modul>.go`

- Daftarkan server ke interface pada file `dics/interface.go`

Pada DI Param

```go
"4": dingo.Service("core:interface:rest"),
"5": dingo.Service("module:<modul>:server"),
```

Pada parameter fungsi

```go
rest configs.Application,
<modul> configs.Server,
```

Pada parameter struct

```go
Servers: []configs.Server{<modul>},
```

Bila bingung, bisa melihat contoh pada [skeleton-todo](https://github.com/crowdeco/skeleton-todo/blob/main/dics/interface.go)

- Daftarkan route ke router pada file `dics/routes.go`

Pada DI Param

```go
Params: dingo.Params{
    "0": dingo.Service("module:todo:server"),
},
```

Pada parameter fungsi

```go
<modul> configs.Server,
```

Pada parameter struct

```go
Servers: []configs.Server{<modul>},
```

Bila bingung, bisa melihat contoh pada [skeleton-todo](https://github.com/crowdeco/skeleton-todo/blob/main/dics/routes.go)

- Update DI Container dengan `go run cmd/dic/main.go`

- Jalankan aplikasi `go run cmd/app/main.go`

- Modulmu dapat diakses via `http://localhost:<APP_PORT>/api/[API_VERSION]/<modul-plural>`

`modul-plural` adalah bentuk plural dari nama modul yang kamu masukkan, kamu juga dapat melihatnya pada file `protos/<modul>.proto`
