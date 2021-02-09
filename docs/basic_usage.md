## Memulai Skeleton

- Buat modul baru dengan perintah `go run cmds/module/main.go register`

![Register Module](assets/imgs/register.png)

- Ikuti setiap langkah yang ada, maka **Skeleton** akan membuatkan modul untukmu secara otomatis dan menambahkan modulmu pada file `modules.yaml`

- Jalankan aplikasi `go run cmds/app/main.go`

- Modulmu dapat diakses via `http://localhost:<APP_PORT>/api/[API_VERSION]/<modul-plural>` `modul-plural` adalah bentuk plural dari nama modul yang kamu masukkan, kamu juga dapat melihatnya pada file `protos/<modul>.proto`, bila bingung, bisa melihat contoh pada [skeleton-todo](https://github.com/crowdeco/skeleton-todo/blob/main/protos/todo.proto#L34)
