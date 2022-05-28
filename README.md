## Skeleton

![Skeleton](assets/imgs/register.png)

![Run App](assets/imgs/app.png)

![Swagger](assets/imgs/swagger.png)

## Arsitektur

![Architecture](assets/imgs/architecture.png)

## Request Flow

![Request Flow](assets/imgs/flow.png)

![Request Explaination](assets/imgs/explain.png)

## Flow Pagination

![Pagination](assets/imgs/paginated.png)

## Flow Create/Update/Delete

![Create/Update/Delete](assets/imgs/create.png)

## Flow Get

![Get](assets/imgs/get.png)

## Tool

Skeleton menggunakan [gRPC Gateway](https://github.com/grpc-ecosystem/grpc-gateway), untuk menginstall tools gRPC Gateway, kamu dapat mengikuti petunjuk pada dokumentasi resminya atau kamu dapat mengikuti tahapan instalasi Skeleton. Bila kamu berhasil menginstall Skeleton sesuai dengan langkah yang diberikan, maka secara secara otomatis, tools gRPC Gateway pun akan terinstall.

## Perintah

- Build Dependency Graph

```
go run cmds/dic/main.go
```

- Application

```
go run cmds/app/main.go
```

- Generator

```
go run cmds/module/main.go register
go run cmds/module/main.go unregister
```

## Testing

```
$ go test ./... [-v]
```

## Contoh

- [Todo Application](https://github.com/KejawenLab/skeleton-todo)

## Dokumentasi

- [Instalasi](docs/install.md)

- [Awal Memulai](docs/basic_usage.md)

- [Memodifikasi Flow](docs/flow_modification.md)

- [Mendaftarkan Log Extension](docs/log_extension.md)

- [Menggunakan Fitur Pub/Sub](docs/pub_sub.md)

- [HTTP Middleware](docs/http_middleware.md)

- [Security](docs/security.md)

- [Custom Route](docs/custom_route.md)

- [List Dependency Injection](docs/dic.md)
