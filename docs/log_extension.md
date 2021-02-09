# Menambahkan Log Extension

**Skeleton** menggunakan [Logrus](https://github.com/sirupsen/logrus) sebagai logger sehingga memungkinkan semua extension yang ada pada Logrus dapat digunakan pada **Skeleton** dengan mudah.

Untuk mendaftarkan log extension, cukup membuatnya sebagai service pada `/dics/modules/<module>.go`, bila bingung kamu bisa melihat cara mendaftarkan MongoDB Extension pada file [`dics/core.go`](https://github.com/crowdeco/skeleton/blob/main/dics/core.go#L417)
