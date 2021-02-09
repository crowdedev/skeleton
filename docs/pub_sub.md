# Cara menggunakan fitur Pub/Sub

## Cara mem-publish pesan

Untuk mem-publish pesan dari file `module.go` kamu cukup memanggil `m.Messenger.Publish("<QueueName>", []byte("Hello from Pub/Sub"))`

## Cara mendapatkan pesan

Untuk mendapatkan pesan, kamu cukup copas potongan code dibawah pada method `Consume()` pada file `module.go`

```go
messages, err := m.Messenger.Consume("<QueueName>")
if err != nil {
    m.Logger.Error(fmt.Sprintf("%+v", err))
}

for message := range messages {
    json.Unmarshal(message.Payload, &v)

    m.Logger.Info(fmt.Sprintf("%+v", v))

    // Your logic here

    message.Ack()
}
```
