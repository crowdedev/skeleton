# Menambahkan Log Extension

**Skeleton** menggunakan [Logrus](https://github.com/sirupsen/logrus) sebagai logger sehingga memungkinkan semua extension yang ada pada Logrus dapat digunakan pada **Skeleton** dengan mudah.

Untuk mendaftarkan log extension, cukup membuatnya sebagai service pada `dics/<module>.go`, bila bingung bisa baca dokumentasi dari [Dingo](https://github.com/sarulabs/dingo)

- Daftarkan service pada file `configs/loggers.yaml`

- Rebuild DIC dengan perintah `go run cmds/dic/main.go`
