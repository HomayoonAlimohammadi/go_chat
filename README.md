# Go Chatting service with gRPC

- Generate protofiles:

```shell
cd proto
protoc --go_out=plugins=grpc:. *.proto
```
