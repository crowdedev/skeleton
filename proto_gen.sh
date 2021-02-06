protoc -Iprotos -Ilibs --go_out=plugins=grpc:protos/builds protos/*.proto
protoc -Iprotos -Ilibs --grpc-gateway_out=logtostderr=true:protos/builds protos/*.proto
protoc -Iprotos -Ilibs --openapiv2_out=allow_merge=true:swagger protos/*.proto
