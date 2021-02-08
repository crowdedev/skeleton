## Memulai Skeleton

- Buat modul baru dengan perintah `go run cmds/generator/main.go`

- Ikuti setiap langkah yang ada, maka **Skeleton** akan membuatkan modul untukmu secara otomatis dan menambahkan modulmu pada file `modules.yaml`

- Daftarkan modul ke sistem pada file `dics/provider.go`

```go
	if err := p.AddDefSlice(modules.<NamaModul>); err != nil {
		return err
	}
```

`NamaModul` adalah nama variabel pada file `dics/modules/<modul>.go`

Bila bingung, bisa melihat contoh pada [skeleton-todo](https://github.com/crowdeco/skeleton-todo/blob/main/dics/provider.go)

- Update DI Container dengan `go run cmds/dic/main.go`

- Jalankan aplikasi `go run cmds/app/main.go`

- Modulmu dapat diakses via `http://localhost:<APP_PORT>/api/[API_VERSION]/<modul-plural>`

`modul-plural` adalah bentuk plural dari nama modul yang kamu masukkan, kamu juga dapat melihatnya pada file `protos/<modul>.proto`

Bila bingung, bisa melihat contoh pada [skeleton-todo](https://github.com/crowdeco/skeleton-todo/blob/main/protos/todo.proto#L34)
