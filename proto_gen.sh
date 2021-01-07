protoc --proto_path=protos --proto_path=libs --go_out=plugins=grpc:protos/builds protos/pagination.proto

protoc --proto_path=protos --proto_path=libs --go_out=plugins=grpc:protos/builds protos/todo.proto
protoc --proto_path=protos --proto_path=libs --grpc-gateway_out=logtostderr=true:protos/builds protos/todo.proto
