# Cara Install Skeleton

## Kebutuhan Software

- Golang versi 1.14 atau lebih tinggi

- Protobuf

- [Protobuf Compiler](https://grpc.io/docs/protoc-installation)

- Database (`mysql`, `sqlserver` atau `postgresql`)

- AMQP ([RabbitMQ](https://www.rabbitmq.com))

- [Elasticsearch](https://www.elastic.co)

- [MongoDB](https://www.mongodb.com)

## Cara Install

- Clone `git clone https://github.com/crowdeco/skeleton.git project`

- Masuk ke project `cd project`

- Jalankan perintah `go run cmds/dic/main.go`

- Install tool

```bash
go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    github.com/golang/protobuf/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

- Install dependency `go mod tidy`

- Buat konfigurasi `cp .env.example .env` dan ubah sesuai kebutuhan
