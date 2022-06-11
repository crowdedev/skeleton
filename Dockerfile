FROM golang:alpine as builder

RUN apk update && apk add --no-cache git
RUN mkdir -p /go/src/app
COPY . /go/src/app
WORKDIR /go/src/app
RUN go mod tidy
RUN go run cmds/dic/main.go
WORKDIR /go/src/app/cmds/app
RUN go build -o bima .

FROM alpine:latest

COPY --from=builder /go/src/app/cmds/app/bima /usr/local/bin/bima
RUN chmod a+x /usr/local/bin/bima

EXPOSE 7777

CMD ["/usr/local/bin/bima"]
