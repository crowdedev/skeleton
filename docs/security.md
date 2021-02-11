# Security pada Skeleton

**Skeleton** di-design untuk berjalan di belakang Api Gateway sehingga untuk security, **Skeleton** menggunakan paradigma yang berbeda. Karena berjalan dibelakang Api Gateway, maka **Skeleton** menganggap bahwa untuk authentication dihandle oleh Api Gateway dan **Skeleton** menjadikan header penanda bahwa sebuah request telah diautentikasi.

Untuk mensetting header, kamu dapat mengaturnya pada `.env` yaitu `HEADER_USER_ID`, `HEADER_USER_EMAIL`, dan `HEADER_USER_ROLE`. Selain itu, kamu juga bisa mengatur authorization dengan mensetting header `MAXIMUM_ROLE` dimana semakin kecil nilai dari header `MAXIMUM_ROLE`, semakin strict.

## Mengaktifkan Security Middleware

Secara default, fitur security yang dijelaskan di atas tidak diaktifkan (disable) dan untuk mengaktifkannya, kamu cukup menghapus komentar pada file `configs/middlewares.yaml` baris `bima:middleware:auth`.
