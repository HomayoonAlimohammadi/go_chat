# Go Chatting service with gRPC

## Generate protofiles:

```shell
cd proto
protoc --go_out=plugins=grpc:. *.proto # for older versions
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative chatpb/chat.proto # for newer versions
```

## Running the server

- In the root directory (e.g. /path/to/file/go_chat/) where the `main.go` resides, run the following command:

```shell
go run main.go server
```

## Running the client

- Just like the server, in the same directory, run:

```shell
go run main.go client
```

- <b>Note:</b> there are 2 optional flags as well:
  - --sender=homayoon
    - specify `sender_name`, defaults to `default` if omitted.
  - --channel=mychannel
    - specify `channel_name`, defaults to `default` if omitted.

## The below example demonstrates how 3 different users can join the `default` channel:

- Shell #1

```shell
go run main.go server
```

- Shell #2

```shell
go run main.go client
```

- Shell #3

```shell
go run mian.go client --sender=user2
```

- Shell #4

```shell
go run main.go client --sender=user3
```

- 3 users by the names of `default`, `user1` and `user2` will be connected to a channel named `default`.

### Try chatting with other users via `stdin`

## Enjoy!
