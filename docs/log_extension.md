# Menambahkan Log Extension

**Skeleton** menggunakan [Logrus](https://github.com/sirupsen/logrus) sebagai logger sehingga memungkinkan semua extension yang ada pada Logrus dapat digunakan pada **Skeleton** dengan mudah.

Untuk mendaftarkan log extension, cukup membuatnya sebagai service pada `/dics/modules/<module>.go`, bila bingung kamu bisa melihat cara mendaftarkan MongoDB Extension pada file [`dics/core.go`](https://github.com/crowdeco/skeleton/blob/main/dics/core.go#L418)

- Daftarkan service pada file `dics/logger.go`, kamu juga bisa menonaktifkan extension dengan menghapusnya dari daftar parameter

- Rebuild DIC dengan perintah `go run cmds/dic/main.go`
