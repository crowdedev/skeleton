protoc -Iprotos -Ilibs --go_out=:protos/builds --go-grpc_out=:protos/builds protos/*.proto
protoc -Iprotos -Ilibs --grpc-gateway_out=logtostderr=true:protos/builds protos/*.proto
protoc -Iprotos -Ilibs --go_out=:protos/builds --go-grpc_out=:protos/builds libs/bima/*.proto
protoc -Iprotos -Ilibs --grpc-gateway_out=logtostderr=true:protos/builds libs/bima/*.proto
protoc -Iprotos -Ilibs --openapiv2_out=swaggers protos/*.proto
