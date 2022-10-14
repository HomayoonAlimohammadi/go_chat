# Go Chatting service with gRPC

- Generate protofiles:

```shell
cd proto
protoc --go_out=plugins=grpc:. *.proto # for older versions
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative chatpb/chat.proto # for newer versions
```
