FROM golang:alpine

RUN apk update && apk add --no-cache git
RUN adduser -D -g '' appuser
WORKDIR $GOPATH/src/app
COPY . .
RUN go mod tidy
RUN go run cmds/dic/main.go
WORKDIR $GOPATH/src/app/cmds/app
RUN go build .

EXPOSE 7777

CMD ["./app"]
